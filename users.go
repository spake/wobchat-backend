package main

import (
    "os/exec"
    "log"
    "encoding/json"
)

type User struct {
    ID        string  `gorm:"primary_key"`
    Name      string
    FirstName string
    LastName  string
    Email     string
    Picture   string
    Friends   []User  `gorm:"many2many:user_friends;"`
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

func verifyIdToken(token string) *User {
    verifierCmd := exec.Command("./verify_token.py", token)
    idinfo, err := verifierCmd.Output()
    if err != nil {
        // unable to verify the token
        log.Println("verify_token.py failed: " + string(idinfo))
        return nil
    }

    // extract the json into an untyped map
    var mapIdinfo map[string]interface{}
    if err := json.Unmarshal(idinfo, &mapIdinfo); err != nil {
        panic(err)
    }

    // turn the map into a user struct
    return &User{
        ID: mapIdinfo["sub"].(string),
        Name: mapIdinfo["name"].(string),
        FirstName: mapIdinfo["given_name"].(string),
        LastName: mapIdinfo["family_name"].(string),
        Email: mapIdinfo["email"].(string),
        Picture: mapIdinfo["picture"].(string)}

}