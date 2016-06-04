package main

import (
    "log"
    "net/http"
    "fmt"
    "github.com/xowap/gowebmake/handlers"
    "github.com/xowap/gowebmake/common"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
    opts := common.ParseOpts()
    config := common.ParseConfig(opts)

    http.HandleFunc("/", handler)
    http.HandleFunc("/github", func(w http.ResponseWriter, r *http.Request) {
        handlers.GitHubHandler(config, w, r)
    })

    log.Println("Binding and listening to " + config.Bind)

    err := http.ListenAndServe(config.Bind, nil)
    if err != nil {
        log.Fatal(err)
    }
}
