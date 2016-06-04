package handlers

import (
    "log"
    "net/http"
    "github.com/xowap/gowebmake/common"
    "github.com/jeffail/gabs"
    "io/ioutil"
    "crypto/hmac"
    "crypto/sha1"
    "strings"
    "encoding/hex"
    "github.com/xowap/gowebmake/git"
    "github.com/xowap/gowebmake/make"
    "path"
    "errors"
)

func parseGHSig(sig string) (buf []byte, ok bool) {
    ok = false

    parts := strings.Split(sig, "=")
    if len(parts) != 2 {
        return
    }

    msg, err := hex.DecodeString(parts[1])
    if err != nil {
        return
    }

    return msg, true
}

func checkGHSig(body []byte, secret string, sig string) bool {
    provided, ok := parseGHSig(sig)
    if !ok {
        return false
    }

    mac := hmac.New(sha1.New, []byte(secret))
    mac.Write(body)
    expected := mac.Sum(nil)

    return hmac.Equal(provided, expected)
}

func guessRepoAddress(gh common.GitHub, body *gabs.Container) (addr string, err error) {
    private, ok := body.Path("repository.private").Data().(bool)

    if !ok {
        return "", errors.New("Impossible to determine privacy of repository")
    }

    if gh.Protocol == "ssh" || gh.Protocol == "" && private {
        if addr, ok = body.Path("repository.ssh_url").Data().(string); !ok {
            return "", errors.New("Missing SSH address")
        }
    } else {
        if addr, ok = body.Path("repository.clone_url").Data().(string); !ok {
            return "", errors.New("Missing HTTP address")
        }
    }

    return
}

func guessBranch(body *gabs.Container) (branch string, err error) {
    ref, ok := body.Path("ref").Data().(string)

    if !ok {
        return "", errors.New("Missing ref from input")
    }

    parts := strings.SplitN(ref, "/", 3)

    if len(parts) != 3 {
        return "", errors.New("Incorrectly formatted ref")
    }

    branch = parts[2]
    return
}

func GitHubHandler(config common.Config, w http.ResponseWriter, r *http.Request) {
    log.Print("Incoming GitHub request... ")

    if r.ContentLength > 10485760 {
        log.Println("Too long")
        return
    }

    bodyRaw, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Println("Could not read body")
    }

    body, err := gabs.ParseJSON(bodyRaw)
    if err != nil {
        log.Println("Could not decode body")
        return
    }

    repo, ok := body.Path("repository.full_name").Data().(string)
    if !ok {
        log.Println("Doesn't have repo name")
        return
    }

    gh, ok := config.GitHub[repo]
    if !ok {
        log.Printf("Repo `%s` is not known", repo)
        return
    }

    v := checkGHSig(bodyRaw, gh.Secret, r.Header.Get("X-Hub-Signature"))

    if !v {
        log.Println("Invalid signature")
        return
    }

    target := path.Join(config.WorkDir, repo)
    repoAddr, err := guessRepoAddress(gh, body)

    if err != nil {
        log.Println(err)
        return
    }

    if err := git.Clone(config, repoAddr, target); err != nil {
        log.Println(err)
        return
    }

    if err := git.Fetch(config, target, "origin"); err != nil {
        log.Println(err)
        return
    }

    branch, err := guessBranch(body)
    if err != nil {
        log.Println(err)
        return
    }

    ghb, ok := gh.Branches[branch]

    if !ok {
        log.Println("Unkown branch got pushed, ignoring")
        return
    }

    commit, ok := body.Path("after").Data().(string)
    if !ok {
        log.Println("Commit was not provided")
        return
    }

    if err := git.Reset(config, target, commit); err != nil {
        log.Println(err)
        return
    }

    log.Println("Git repository was updated")

    if err := make.RunTarget(config, target, ghb); err != nil {
        log.Println(err)
        return
    }

    log.Printf("make `%s` ran all right!", ghb.Target)
}
