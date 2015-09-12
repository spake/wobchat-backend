package main

import (
    "encoding/json"
    "io"
    "net/http"
)

type AboutResponse struct {
    BuildNumber string
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
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
