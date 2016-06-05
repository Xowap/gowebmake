package main

import (
    "log"
    "net/http"
    "github.com/xowap/gowebmake/handlers"
    "github.com/xowap/gowebmake/common"
)

func main() {
    opts := common.ParseOpts()
    config := common.ParseConfig(opts)

    http.HandleFunc("/github", func(w http.ResponseWriter, r *http.Request) {
        handlers.GitHubHandler(config, w, r)
    })

    log.Println("Binding and listening to " + config.Bind)

    err := http.ListenAndServe(config.Bind, nil)
    if err != nil {
        log.Fatal(err)
    }
}
