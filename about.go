package main

import (
    "encoding/json"
    "io"
    "net/http"
    "log"
)

type AboutResponse struct {
    BuildNumber string  `json:"buildNumber"`
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("Handling /about")
    var build string

    // Get build number (if it was compiled in)
    if bambooBuildNumber != "" {
        build = bambooBuildNumber
    } else {
        build = "unknown"
    }

    resp := AboutResponse{
        BuildNumber: build,
    }

    b, _ := json.Marshal(resp)
    io.WriteString(w, string(b))
}
