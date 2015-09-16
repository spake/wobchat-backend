package main

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    "github.com/jinzhu/gorm"
)

type User struct {
    Uid       int       `gorm:"primary_key"`
    Name      string
    FirstName string
    LastName  string
    Email     string
    Picture   string
}

type UserFriend struct {
    UserUid   int       `gorm:"primary_key"`
    FriendUid int       `gorm:"primary_key"`
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
    w.Header().Set("Access-Control-Allow-Origin", "https://wob.chat")
    io.WriteString(w, string(b))
}

func (user *User) getFriends(db gorm.DB) []User {
    friends := []User{}
    db.Joins("inner join user_friends on friend_uid = uid").Where(&UserFriend{UserUid: user.Uid}).Find(&friends)
    return friends
}

func (user *User) addFriend(db gorm.DB, friend User) {
    userFriend := UserFriend{UserUid: user.Uid,FriendUid: friend.Uid}
    db.Create(&userFriend)
}