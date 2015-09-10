package main

import (
    "io"
    "net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Uhh, hello?")
}

func main() {
    http.HandleFunc("/", hello)
    http.ListenAndServe("127.0.0.1:8000", nil)
}
