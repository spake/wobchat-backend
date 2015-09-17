package main

import (
    "encoding/json"
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

/*
 * API endpoints
 */

/*
 * /verify
 * Verifies a Google token. (More of a dummy endpoint at this stage.)
 */
type VerifyRequest struct {
    Token   string  `json:"token"`
}

type VerifyResponse struct {
    OK      bool    `json:"ok"`
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
    c := newContext(r)

    resp := VerifyResponse{}

    _, ok := func() (GoogleInfo, bool) {
        info := GoogleInfo{}

        // decode json request
        decoder := json.NewDecoder(r.Body)
        var req VerifyRequest
        err := decoder.Decode(&req)
        if err != nil {
            return info, false
        }

        // verify token using the google JWT stuff, and get their info
        info, err = verifyIDToken(c, req.Token)
        if err != nil {
            return info, false
        }

        // success!
        return info, true
    }();

    resp.OK = ok

    if ok {
        // TODO: save shit into db
    }

    sendJSONResponse(w, resp)
}
