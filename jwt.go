// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "bytes"
    "encoding/gob"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "time"

    jwt "github.com/dgrijalva/jwt-go"
    "golang.org/x/net/context"
)

// Stores info obtained from a Google JWT token
type GoogleInfo struct {
    ID          string
    DisplayName string
    FirstName   string
    LastName    string
    Email       string
    Picture     string
}

/*
    idinfo looks like this when it's returned

    {"picture": "<url>", 
    "aud": "<client id>.apps.googleusercontent.com", 
    "family_name": "Smith", 
    "iss": "accounts.google.com", 
    "email_verified": true, 
    "name": "Jayden Smith", 
    "at_hash": "<some hash>", 
    "given_name": "Jayden", 
    "exp": <some number>, 
    "azp": "<client id>.apps.googleusercontent.com", 
    "iat": <some number>, 
    "locale": "en", 
    "email": "jaydensmith@gmail.com", 
    "sub": "<user id>"}
*/

// verifyIDToken verifies Google ID token, which heavily based on JWT.
// It returns user ID of the pricipal who granted an authorization.
func verifyIDToken(c context.Context, t string) (GoogleInfo, error) {
    token, err := jwt.Parse(t, func(j *jwt.Token) (interface{}, error) {
        kid, _ := j.Header["kid"].(string)
        keys, err := idTokenCerts(c)
        if err != nil {
            return nil, err
        }
        cert, ok := keys[kid]
        if !ok {
            return nil, fmt.Errorf("verifyIDToken: keys[%q] = nil", kid)
        }
        return cert, nil
    })

    info := GoogleInfo{}

    if err != nil {
        return info, err
    }

    info.ID = token.Claims["sub"].(string)
    info.DisplayName = token.Claims["name"].(string)
    info.FirstName = token.Claims["given_name"].(string)
    info.LastName = token.Claims["family_name"].(string)
    info.Email = token.Claims["email"].(string)
    // some users may not have a picture
    picUrl, ok := token.Claims["picture"].(string)
    if !ok {
        info.Picture = ""
    } else {
        info.Picture = picUrl
    }

    return info, nil
}

// idTokenCerts returns public certificates used to encrypt ID tokens.
// It returns a cached copy, if available, or fetches from a known URL otherwise.
// The returnd map is keyed after the cert IDs.
func idTokenCerts(c context.Context) (map[string][]byte, error) {
    //certURL := config.Google.CertURL
    certURL := "https://www.googleapis.com/oauth2/v1/certs"
    // try cache first
    keys, err := certsFromCache(c, certURL)
    if err == nil {
        return keys, nil
    }
    // fetch from public endpoint otherwise
    var exp time.Duration
    keys, exp, err = fetchPublicKeys(c, certURL)
    if err != nil {
        return nil, err
    }
    if exp <= 0 {
        return keys, nil
    }
    // cache the result for duration exp
    var data bytes.Buffer
    if err := gob.NewEncoder(&data).Encode(keys); err != nil {
        fmt.Errorf("idTokenCerts: %v", err)
    } else if err := cache.set(c, certURL, data.Bytes(), exp); err != nil {
        fmt.Errorf("idTokenCerts: cache.set(%q): %v", certURL, err)
    }
    // return the result anyway, even on cache errors
    return keys, nil
}

// certsFromCache returns cached public keys.
// See idTokenCerts func.
func certsFromCache(c context.Context, k string) (map[string][]byte, error) {
    data, err := cache.get(c, k)
    if err != nil {
        return nil, err
    }
    var keys map[string][]byte
    return keys, gob.NewDecoder(bytes.NewReader(data)).Decode(&keys)
}

// httpTransport returns a suitable HTTP transport for current backend
// hosting environment.
// It uses http.DefaultTransport by default.
var httpTransport = func(_ context.Context) http.RoundTripper {
    return http.DefaultTransport
}

// httpClient create a new HTTP client using httpTransport(),
// setting request timeout to 10 seconds if supported.
func httpClient(c context.Context) *http.Client {
    cl := &http.Client{Transport: httpTransport(c)}
    type canceler interface {
        CancelRequest(*http.Request)
    }
    if _, ok := cl.Transport.(canceler); ok {
        // NB: this got commented out because it didn't work
        //cl.Timeout = 10 * time.Second
    }
    return cl
}

// certsFromCache fetches public keys from the network.
// See idTokenCerts func.
func fetchPublicKeys(c context.Context, url string) (map[string][]byte, time.Duration, error) {
    res, err := httpClient(c).Get(url)
    if err != nil {
        return nil, 0, err
    }
    defer res.Body.Close()
    if res.StatusCode != http.StatusOK {
        return nil, 0, fmt.Errorf("fetchPublicKeys: %s: %v", url, res.Status)
    }
    var body map[string]string
    if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
        return nil, 0, err
    }
    keys := make(map[string][]byte)
    for k, v := range body {
        keys[k] = []byte(v)
    }
    return keys, resourceExpiry(res.Header), nil
}

// resourceExpiry returns the remaining life of a resource
// based on Cache-Control and Age headers.
func resourceExpiry(h http.Header) time.Duration {
    var max int64
    for _, c := range strings.Split(h.Get("cache-control"), ",") {
        c = strings.ToLower(strings.TrimSpace(c))
        if !strings.HasPrefix(c, "max-age=") {
            continue
        }
        var err error
        if max, err = strconv.ParseInt(c[8:], 10, 64); err != nil {
            max = 0
        }
        break
    }
    age, err := strconv.ParseInt(h.Get("age"), 10, 64)
    if err != nil {
        age = 0
    }
    r := max - age
    if r < 0 {
        return 0
    }
    return time.Duration(r) * time.Second
}

func newContext(r *http.Request) context.Context {
    return context.Background()
}
