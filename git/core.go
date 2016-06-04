package git

import (
    "path"
    "fmt"
    "os"
    "errors"
    "os/exec"
    "github.com/xowap/gowebmake/common"
)

func Clone(config common.Config, repo, target string) error {
    p := path.Dir(target)

    stat, err := os.Stat(target)

    if err != nil {
        if !os.IsNotExist(err) {
            return err
        }

        if err := os.MkdirAll(p, 0755); err != nil {
            return err
        }

        cmd := exec.Command(config.GitBin, "clone", repo, target)
        if err := cmd.Run(); err != nil {
            return errors.New("Could not clone the repo")
        }

        cmd = exec.Command(config.GitBin, "submodule", "update", "--init")
        cmd.Dir = target
        if err := cmd.Run(); err != nil {
            return errors.New("Could not init submodules")
        }
    } else if !stat.IsDir() {
        return errors.New(fmt.Sprintf("`%s` is not a directory", target))
    } else {
        stat, err := os.Stat(path.Join(target, ".git"))

        if err != nil {
            if os.IsNotExist(err) {
                return errors.New(fmt.Sprintf("`%s` is not a git repository"))
            } else {
                return err
            }
        }

        if !stat.IsDir() {
            return errors.New(fmt.Sprint("`%s` is not a git repository"))
        }
    }

    return nil
}

func Fetch(config common.Config, target string, remote string) error {
    cmd := exec.Command(config.GitBin, "fetch", remote, "-p")
    cmd.Dir = target
    if err := cmd.Run(); err != nil {
        return errors.New("Could not fetch remote")
    }

    return nil
}

func Reset(config common.Config, target string, commit string) error {
    cmd := exec.Command(config.GitBin, "reset", "--hard", commit)
    cmd.Dir = target
    if err := cmd.Run(); err != nil {
        return errors.New("Could not reset")
    }

    cmd = exec.Command(config.GitBin, "submodule", "update", "--init")
    cmd.Dir = target
    if err := cmd.Run(); err != nil {
        return errors.New("Error while updating submodules")
    }

    return nil
}
