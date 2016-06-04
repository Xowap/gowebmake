package common

import (
    "log"
    "github.com/jessevdk/go-flags"
    "os"
    "github.com/burntsushi/toml"
)

type Options struct {
    ConfigPath string `short:"c" long:"config" description:"Path to the configuration file" required:"true"`
}

type GitHubBranch struct {
    Target string
    Env map[string]string
}

type GitHub struct {
    Protocol string
    Secret string
    Branches map[string]GitHubBranch
}

type Config struct {
    Bind string
    WorkDir string
    GitHub map[string]GitHub
    GitBin string
    MakeBin string
}

func ParseOpts() (opts Options) {
    _, err := flags.Parse(&opts)
    ferr, ok := err.(*flags.Error)
    if err != nil {
        if ok && ferr.Type != flags.ErrHelp {
            log.Fatal(err)
        }

        os.Exit(1)
    }

    return
}

func ParseConfig(opts Options) (config Config) {
    if _, err := toml.DecodeFile(opts.ConfigPath, &config); err != nil {
        log.Fatal(err)
    }

    if config.Bind == "" {
        config.Bind = "[::1]:8777"
    }

    if config.WorkDir == "" {
        log.Fatal("Option \"workdir\" is required")
    }

    if config.GitBin == "" {
        config.GitBin = "git"
    }

    if config.MakeBin == "" {
        config.MakeBin = "make"
    }

    return
}
