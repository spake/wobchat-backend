package main

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
)

type User struct {
    UID       int       `gorm:"primary_key"`
    Name      string
    FirstName string
    LastName  string
    Email     string
    Picture   string
}

type VerifyRequest struct {
    Token   string
}

type VerifyResponse struct {
    OK      bool
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
    c := newContext(r)

    resp := VerifyResponse{}

    info, ok := func() (GoogleInfo, bool) {
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

    log.Printf("info: %v\n", info)

    b, _ := json.Marshal(resp)
    io.WriteString(w, string(b))
}
