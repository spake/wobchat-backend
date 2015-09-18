package main

import (
    "log"
    "net/http"
    "github.com/jinzhu/gorm"
)

/*
 * DB data types
 */
// Represents a user in the database
type User struct {
    Uid       string    `gorm:"primary_key"`
    Name      string
    FirstName string
    LastName  string
    Email     string
    Picture   string
}

// Represents one-way friendship in the database
// Note that two rows are required to represent a reflexive friendship
type UserFriend struct {
    UserUid   string    `gorm:"primary_key"`
    FriendUid string    `gorm:"primary_key"`
}

type Users []User

// Smaller version of User without sensitive/unnecessary info, for sending to
// third parties, like a user's friends
type PublicUser struct {
    Uid         string  `json:"uid"`
    Name        string  `json:"name"`
    FirstName   string  `json:"firstName"`
    LastName    string  `json:"lastName"`
    Picture     string  `json:"picture"`
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

func (user *User) toPublic() PublicUser {
    return PublicUser{
        Uid:        user.Uid,
        Name:       user.Name,
        FirstName:  user.FirstName,
        LastName:   user.LastName,
        Picture:    user.Picture,
    }
}

func (users *Users) toPublic() (publicUsers []PublicUser) {
    for _, user := range *users {
        publicUsers = append(publicUsers, user.toPublic())
    }
    return
}

/*
 * DB manipulation functions
 */
func getUserFromInfo(db gorm.DB, info GoogleInfo) (user User) {
    log.Printf("Getting user %v\n", info.ID)

    // check if user already exists
    if err := db.Where("uid = ?", info.ID).First(&user).Error; err != nil {
        // create user
        user = User{
            Uid:        info.ID,
            Name:       info.DisplayName,
            FirstName:  info.FirstName,
            LastName:   info.LastName,
            Email:      info.Email,
            Picture:    info.Picture,
        }
        db.Create(&user)
    } else {
        // update things from the info, in case they've changed
        user.Uid = info.ID
        user.Name = info.DisplayName
        user.FirstName = info.FirstName
        user.LastName = info.LastName
        user.Email = info.Email
        user.Picture = info.Picture
        db.Save(&user)
    }

    return user
}

func getCurrentUser(db gorm.DB, r *http.Request) (user User, ok bool) {
    info, authenticated := getAuthInfo(r)
    if !authenticated {
        log.Println("Not authenticated")
        return user, false
    }

    user = getUserFromInfo(db, info)
    return user, true
}

/*
 * API endpoints
 */

/*
 * /friends
 * Gets a list of the current user's friends.
 */
type ListFriendsResponse struct {
    Friends []PublicUser    `json:"friends"`
}

func listFriendsHandler(w http.ResponseWriter, r *http.Request) int {
    user, ok := getCurrentUser(db, r)
    if !ok {
        return http.StatusUnauthorized
    }

    var friends Users
    friends = user.getFriends(db)

    resp := ListFriendsResponse{
        Friends: friends.toPublic(),
    }

    sendJSONResponse(w, resp)
    return http.StatusOK
}
