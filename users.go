package main

import (
    "log"
    "net/http"
    "encoding/json"
    "errors"
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

func (user *User) getFriends() Users {
    friends := []User{}
    db.Joins("inner join user_friends on friend_uid = uid").Where(&UserFriend{UserUid: user.Uid}).Find(&friends)
    return friends
}

func (user *User) addFriend(friend User) error {
    if user.Uid != friend.Uid {
        tx := db.Begin()

        userFriend := UserFriend{UserUid: user.Uid,FriendUid: friend.Uid}
        if err := tx.Create(&userFriend).Error; err != nil {
            tx.Rollback()
            return err
        }

        userFriend = UserFriend{UserUid: friend.Uid,FriendUid: user.Uid}
        if err := tx.Create(&userFriend).Error; err != nil {
            tx.Rollback()
            return err
        }

        tx.Commit()
        return nil
    } else {
        return errors.New("Cannot add yourself as a friend")
    }
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
func getUserFromInfo(info GoogleInfo) (user User) {
    log.Printf("Getting user %v\n", info.ID)

    log.Println(db)

    // check if user already exists
    if err := db.Where(&User{Uid: info.ID}).First(&user).Error; err != nil {
        // create user
        log.Println("Creating new user in db")
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
        log.Println("Updating existing user in db")
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

func getCurrentUser(r *http.Request) (user User, ok bool) {
    info, authenticated := getAuthInfo(r)
    if !authenticated {
        log.Println("Not authenticated")
        return user, false
    }

    user = getUserFromInfo(info)
    return user, true
}

/*
 * API endpoints
 */

/*
 * /friends endpoint
 */

func friendsHandler(w http.ResponseWriter, r *http.Request) int {
    log.Println("Handling /friends")
    user, ok := getCurrentUser(r)
    if !ok {
        return http.StatusUnauthorized
    }

    var resp interface{}

    switch r.Method {
    case "GET":
        resp = listFriendsEndpoint(user)
    case "POST":
        decoder := json.NewDecoder(r.Body)
        var req AddFriendsRequest
        err := decoder.Decode(&req)
        if err != nil {
            return http.StatusBadRequest
        }
        resp = addFriendEndpoint(user, req)
    default:
        return http.StatusMethodNotAllowed
    }

    sendJSONResponse(w, resp)
    return http.StatusOK
}

/*
 * GET /friends
 * Gets a list of the current user's friends.
 */
type ListFriendsResponse struct {
    Friends []PublicUser    `json:"friends"`
}

func listFriendsEndpoint(user User) ListFriendsResponse {
    var friends Users
    friends = user.getFriends()

    resp := ListFriendsResponse{
        Friends: friends.toPublic(),
    }

    return resp
}

/*
 * POST /friends
 * Adds a friend to the current user.
 */
type AddFriendsRequest struct {
    Uid     string    `json:"uid"`
}

type AddFriendsResponse struct {
    Success bool      `json:"success"`
    Error   string    `json:"error"`
}

func addFriendEndpoint(user User, req AddFriendsRequest) AddFriendsResponse {
    var friend User
    dbErr := db.Where(&User{Uid: req.Uid}).First(&friend).Error

    var resp AddFriendsResponse

    log.Println("DB Error: ", dbErr)

    if dbErr == nil {
        addErr := user.addFriend(friend)

        if addErr != nil {
            resp = AddFriendsResponse{
                Success: false,
                Error:   addErr.Error()}
        } else {
            resp = AddFriendsResponse{
                Success: true,
                Error:   ""}
        }
    } else {
        // friend they are trying to add not found
        resp = AddFriendsResponse{
                Success: false,
                Error:   "Friend not found"}
    }
    return resp
}
