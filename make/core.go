package make

import (
    "os/exec"
    "errors"
    "github.com/xowap/gowebmake/common"
    "os"
    "fmt"
)

func RunTarget(config common.Config, wd string, ghb common.GitHubBranch) error {
    if ghb.Target == "" {
        return errors.New("Missing make target!")
    }

    env := os.Environ()

    for k, v := range ghb.Env {
        env = append(env, fmt.Sprintf("%s=%s", k, v))
    }

    cmd := exec.Command(config.MakeBin, ghb.Target)
    cmd.Dir = wd
    cmd.Env = env
    if err := cmd.Run(); err != nil {
        return errors.New("make target failed")
    }

    return nil
}
