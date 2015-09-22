package main

import (
    "net/http"
)

/*
 * Helper functions
 */

// Gets GoogleInfo (and whether the user is authenticated)
func getAuthInfo(r *http.Request) (info GoogleInfo, authenticated bool) {
    c := newContext(r)

    // Get token from header
    if tokens, ok := r.Header["X-Session-Token"]; ok {
        info, err := verifyIDToken(c, tokens[0])
        if err == nil {
            return info, true
        }
    }

    return info, false
}
