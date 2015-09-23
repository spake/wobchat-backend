package main

import (
    "encoding/json"
    "errors"
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/gorilla/mux"
)

/*
 * DB data types
 */
// Represents a user in the database
type User struct {
    Id        int       `gorm:"primary_key" sql:"auto_increment"`
    Uid       string    `sql:"unique"`
    Name      string
    FirstName string
    LastName  string
    Email     string
    Picture   string
}

// Represents one-way friendship in the database
// Note that two rows are required to represent a reflexive friendship
type UserFriend struct {
    UserId      int `gorm:"primary_key"`
    FriendId    int `gorm:"primary_key"`
}

type Users []User

// Smaller version of User without sensitive/unnecessary info, for sending to
// third parties, like a user's friends
type PublicUser struct {
    Id          int     `json:"id"`
    Uid         string  `json:"uid"`
    Name        string  `json:"name"`
    FirstName   string  `json:"firstName"`
    LastName    string  `json:"lastName"`
    Picture     string  `json:"picture"`
}

func (user *User) getFriends() Users {
    friends := []User{}
    db.Joins("inner join user_friends on friend_id = id").Where(&UserFriend{UserId: user.Id}).Find(&friends)
    return friends
}

func (user *User) addFriend(friend User) error {
    if user.Id != friend.Id {
        tx := db.Begin()

        userFriend := UserFriend{UserId: user.Id, FriendId: friend.Id}
        if err := tx.Create(&userFriend).Error; err != nil {
            tx.Rollback()
            return err
        }

        userFriend = UserFriend{UserId: friend.Id, FriendId: user.Id}
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

func (user *User) isFriend(friend User) bool {
    var uf UserFriend
    if err := db.Where(&UserFriend{UserId: user.Id, FriendId: friend.Id}).Find(&uf).Error; err != nil {
        return false
    }
    return true
}

func (user *User) toPublic() PublicUser {
    return PublicUser{
        Id:         user.Id,
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

func (user *User) getMessagesWithUser(otherUser User) (msgs Messages) {
    db.Where("(sender_id = ? and recipient_id = ?) or (sender_id = ? and recipient_id = ?)", user.Id, otherUser.Id, otherUser.Id, user.Id).Find(&msgs)
    return msgs
}

func (user *User) addMessageToUser(otherUser User, content string, contentType ContentType) (id int, err error) {
    if !contentType.valid() {
        return 0, errors.New("Invalid content type")
    }

    msg := Message{
        Content:        content,
        ContentType:    contentType,
        SenderId:       user.Id,
        RecipientId:    otherUser.Id,
        RecipientType:  RecipientTypeUser,
        Timestamp:      time.Now(),
    }

    db.Create(&msg)

    return msg.Id, nil
}

/*
 * DB manipulation functions
 */
func getUserFromInfo(info GoogleInfo) (user User) {
    log.Printf("Getting user %v\n", info.ID)

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
        var req AddFriendRequest
        err := decoder.Decode(&req)
        if err != nil {
            return http.StatusBadRequest
        }
        if req.Id <= 0 {
            log.Println("Friend ID not positive integer")
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
    Success bool            `json:"success"`
    Friends []PublicUser    `json:"friends"`
}

func listFriendsEndpoint(user User) ListFriendsResponse {
    var friends Users
    friends = user.getFriends()

    resp := ListFriendsResponse{
        Success:    true,
        Friends:    friends.toPublic(),
    }

    return resp
}

/*
 * POST /friends
 * Adds a friend to the current user.
 */
type AddFriendRequest struct {
    Id  int `json:"id"`
}

type AddFriendResponse struct {
    Success bool        `json:"success"`
    Error   string      `json:"error"`
    Friend  PublicUser  `json:"friend"`
}

func addFriendEndpoint(user User, req AddFriendRequest) AddFriendResponse {
    var friend User
    dbErr := db.Where(&User{Id: req.Id}).First(&friend).Error

    if dbErr != nil {
        // friend they are trying to add not found
        return AddFriendResponse{
                Success: false,
                Error:   "Friend not found"}
    }
    
    addErr := user.addFriend(friend)

    if addErr != nil {
        return AddFriendResponse{
            Success: false,
            Error:   addErr.Error()}
    }

    return AddFriendResponse{
        Success: true,
        Friend: friend.toPublic()}
}

/*
 * /friends/{friendId} endpoint
 */

func friendHandler(w http.ResponseWriter, r *http.Request) int {
    log.Println("Handling /friend/{friendId}")
    user, ok := getCurrentUser(r)
    if !ok {
        return http.StatusUnauthorized
    }

    vars := mux.Vars(r)
    friendId, err := strconv.Atoi(vars["friendId"])
    if err != nil || friendId <= 0 {
        log.Println("Friend ID not positive integer")
        return http.StatusBadRequest
    }

    var resp interface{}

    switch r.Method {
    case "GET":
        resp = getFriendEndpoint(user, friendId)
    default:
        return http.StatusMethodNotAllowed
    }

    sendJSONResponse(w, resp)
    return http.StatusOK
}

/*
 * GET /friends/{friendId}
 * Gets a list of the messages between the current user and the specified friend.
 */
type GetFriendResponse struct {
    Success bool        `json:"success"`
    Error   string      `json:"error"`
    Friend  PublicUser  `json:"friend"`
}

func getFriendEndpoint(user User, friendId int) GetFriendResponse {
    if friendId == user.Id {
        return GetFriendResponse{
            Success:    false,
            Error:      "Friend ID cannot be your own",
        }
    }

    var friend User
    dbErr := db.Where(&User{Id: friendId}).First(&friend).Error

    if dbErr != nil {
        // friend not found
        return GetFriendResponse{
            Success:    false,
            Error:      "Friend not found",
        }
    }

    if !user.isFriend(friend) {
        return GetFriendResponse{
            Success:    false,
            Error:      "User is not your friend",
        }
    }

    return GetFriendResponse{
        Success:    true,
        Friend:     friend.toPublic(),
    }
}
