package main

import (
    "encoding/json"
    "io"
    "net/http"
)

// Set standard response headers
func setHTTPHeaders(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "https://wob.chat")
}

// Sends a JSON response, and sets up necessary headers
func sendJSONResponse(w http.ResponseWriter, resp interface{}) {
    setHTTPHeaders(w)

    b, _ := json.Marshal(resp)
    io.WriteString(w, string(b))
}

// Custom HTTP handler type
type APIHandler func(http.ResponseWriter, *http.Request) int
func (handler APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if status := handler(w, r); status != http.StatusOK {
        http.Error(w, http.StatusText(status), status)
    }
}

func setupAPIHandlers() {
    http.HandleFunc("/about", aboutHandler)
    http.HandleFunc("/verify", verifyHandler)
    http.Handle("/friends", APIHandler(listFriendsHandler))
}
