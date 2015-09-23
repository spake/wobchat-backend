package main

import (
    "encoding/json"
    "io"
    "net/http"

    "github.com/gorilla/mux"
)

// Sends a JSON response, and sets up necessary headers
func sendJSONResponse(w http.ResponseWriter, resp interface{}) {
    w.Header().Set("Content-Type", "application/json")

    b, _ := json.Marshal(resp)
    io.WriteString(w, string(b))
}

// Custom HTTP handler type
type APIHandler func(http.ResponseWriter, *http.Request) int
func (handler APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Handle origin stuff, otherwise cross-domain frontend requests will fail
    if origin := r.Header.Get("Origin"); origin != "" {
        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers",
            "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Session-Token")
    }
    if r.Method == "OPTIONS" {
        return
    }

    // Call the actual handler
    if status := handler(w, r); status != http.StatusOK {
        http.Error(w, http.StatusText(status), status)
    }
}

func setupAPIHandlers() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)

    router.Handle("/friends", APIHandler(friendsHandler))
    router.Handle("/friends/{friendId:[0-9]+}", APIHandler(friendHandler))
    router.Handle("/friends/{friendId:[0-9]+}/messages", APIHandler(messagesHandler))
    router.Handle("/me", APIHandler(meHandler))

    return router
}
