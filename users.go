package main

import (
    "encoding/json"
    "errors"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
    "regexp"

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

// Represents a 
type FriendRequest struct {
    UserId         int `gorm:"primary_key"`
    RequestorId    int `gorm:"primary_key"`
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
    if user.Id == friend.Id {
        return errors.New("Cannot add yourself as a friend")
    }

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
}

func (user *User) deleteFriend(friend User) error {
    var uf1, uf2 UserFriend

    // check two-way friendship exists
    if err := db.Where(&UserFriend{UserId: user.Id, FriendId: friend.Id}).First(&uf1).Error; err != nil {
        return err
    }
    if err := db.Where(&UserFriend{UserId: friend.Id, FriendId: user.Id}).First(&uf2).Error; err != nil {
        return err
    }

    tx := db.Begin()

    // do actual deleting
    if err := db.Where("user_id = ? and friend_id = ?", user.Id, friend.Id).Delete(UserFriend{}).Error; err != nil {
        tx.Rollback()
        return err
    }
    if err := db.Where("friend_id = ? and user_id = ?", user.Id, friend.Id).Delete(UserFriend{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    tx.Commit()
    return nil
}

func (user *User) isFriend(friend User) bool {
    var uf UserFriend
    if err := db.Where(&UserFriend{UserId: user.Id, FriendId: friend.Id}).Find(&uf).Error; err != nil {
        return false
    }
    return true
}

// get friend requests sent to a user
func (user *User) getFriendRequests() Users {
    requestors := []User{}
    db.Joins("inner join friend_requests on requestor_id = id").Where(&UserFriend{UserId: user.Id}).Find(&requestors)
    return requestors
}

// add a friend request to user from requestor
func (user *User) addFriendRequest(requestor User) error {
    if user.Id == requestor.Id {
        return errors.New("Cannot request to be your own friend")
    }

    err := db.Create(&FriendRequest{UserId: user.Id, RequestorId: requestor.Id}).Error

    return err
}

// get whether a friend request has been sent from requestor to user
func (user *User) hasFriendRequest(requestor User) bool {
    var friendRequest FriendRequest
    if err := db.Where(&FriendRequest{UserId: user.Id, RequestorId: requestor.Id}).Find(&friendRequest).Error; err != nil {
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

func (m Messages) reverse() {
    for i := 0; i < len(m)/2; i++ {
        m[i], m[len(m)-i-1] = m[len(m)-i-1], m[i]
    }
}

func (user *User) getMessagesWithUser(otherUser User, last int, amount int) (msgs Messages) {
    if last == -1 {
        db.Where("((sender_id = ? and recipient_id = ?) or (sender_id = ? and recipient_id = ?))", user.Id, otherUser.Id, otherUser.Id, user.Id).Order("id desc").Limit(amount).Find(&msgs)
    } else {
        db.Where("((sender_id = ? and recipient_id = ?) or (sender_id = ? and recipient_id = ?)) and id < ?", user.Id, otherUser.Id, otherUser.Id, user.Id, last).Order("id desc").Limit(amount).Find(&msgs)
    }
    msgs.reverse()
    
    return msgs
}

// Gets next message (i.e. with a greater ID than afterId) that the user has received
func (user *User) getNextMessageAfterId(afterId int) (msg Message, ok bool) {
    if err := db.Where("recipient_id = ? and id > ?", user.Id, afterId).First(&msg).Error; err == nil {
        return msg, true
    }
    return Message{}, false
}

func (user *User) addMessageToUser(otherUser User, content string, contentType ContentType) (msg Message, err error) {
    if !contentType.valid() {
        return msg, errors.New("Invalid content type")
    }

    msg = Message{
        Content:        content,
        ContentType:    contentType,
        SenderId:       user.Id,
        RecipientId:    otherUser.Id,
        RecipientType:  RecipientTypeUser,
        Timestamp:      time.Now(),
    }

    db.Create(&msg)

    return msg, nil
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

// search for users by name or by email
func searchUsernames(q string, userid int) (users Users) {
    // check if q looks like an email
    if match, _ := regexp.MatchString(".+@.+\\..+", q); match {
        // search by email
        db.Where("upper(email) = ? and id != ?", strings.ToUpper(q), userid).Find(&users)
    } else {
        // search by name
        db.Where("upper(name) LIKE ? and id != ?", "%%"+strings.ToUpper(q)+"%%", userid).Find(&users)
    }
    return users
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
 * Adds a user as a friend of the current user.
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
    case "DELETE":
        resp = deleteFriendEndpoint(user, friendId)
    default:
        return http.StatusMethodNotAllowed
    }

    sendJSONResponse(w, resp)
    return http.StatusOK
}

/*
 * GET /friends/{friendId}
 * Gets a friend of the current user by their Id.
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

/*
 * DELETE /friends/{friendId}
 * Removes a user from the current user's friends list.
 */
type DeleteFriendResponse struct {
    Success bool    `json:"success"`
    Error   string  `json:"error"`
}

func deleteFriendEndpoint(user User, friendId int) DeleteFriendResponse {
    if friendId == user.Id {
        return DeleteFriendResponse{
            Success:    false,
            Error:      "Friend ID cannot be your own",
        }
    }

    var friend User
    dbErr := db.Where(&User{Id: friendId}).First(&friend).Error

    if dbErr != nil {
        return DeleteFriendResponse{
            Success:    false,
            Error:      "Friend not found",
        }
    }

    if !user.isFriend(friend) {
        return DeleteFriendResponse{
            Success:    false,
            Error:      "User is not your friend",
        }
    }

    // actually delete the friend
    if err := user.deleteFriend(friend); err != nil {
        return DeleteFriendResponse{
            Success:    false,
            Error:      err.Error(),
        }
    }

    return DeleteFriendResponse{
        Success:    true,
    }
}

/*
 * /users endpoint
 */

func usersHandler(w http.ResponseWriter, r *http.Request) int {
    log.Println("Handling /users")
    user, ok := getCurrentUser(r)
    if !ok {
        return http.StatusUnauthorized
    }

    var resp interface{}

    switch r.Method {
    case "GET":
        q := r.FormValue("q")
        resp = listUsersEndpoint(q, user.Id)
    default:
        return http.StatusMethodNotAllowed
    }

    sendJSONResponse(w, resp)
    return http.StatusOK
}

/*
 * /me endpoint
 */

func meHandler(w http.ResponseWriter, r *http.Request) int {
    log.Println("Handling /me")
    user, ok := getCurrentUser(r)
    if !ok {
        return http.StatusUnauthorized
    }

    var resp interface{}

    switch r.Method {
    case "GET":
        resp = getMeEndpoint(user)
    default:
        return http.StatusMethodNotAllowed
    }

    sendJSONResponse(w, resp)
    return http.StatusOK
}

/*
 * GET /users
 * Gets a list of all users (except the current user) whose names match the given query.
 */
type ListUsersResponse struct {
    Success bool            `json:"success"`
    Users []PublicUser      `json:"users"`
}

func listUsersEndpoint(q string, userid int) ListUsersResponse {
    var users Users
    users = searchUsernames(q, userid)

    resp := ListUsersResponse{
        Success:    true,
        Users:    users.toPublic(),
    }

    return resp
}

/*
 * GET /me
 * Gets information about the current user.
 */
type GetMeResponse struct {
    Success bool        `json:"success"`
    Error   string      `json:"error"`
    User    PublicUser  `json:"user"`
}

func getMeEndpoint(currentUser User) GetMeResponse {
    return GetMeResponse{
        Success:    true,
        Error:      "",
        User:       currentUser.toPublic(),
    }
}

/*
 * /friendrequests endpoint
 */

func myFriendRequestsHandler(w http.ResponseWriter, r *http.Request) int {
    log.Println("Handling /friendrequests")
    user, ok := getCurrentUser(r)
    if !ok {
        return http.StatusUnauthorized
    }

    var resp interface{}

    switch r.Method {
    case "GET":
        resp = listMyFriendRequestsEndpoint(user)
    default:
        return http.StatusMethodNotAllowed
    }

    sendJSONResponse(w, resp)
    return http.StatusOK
}

/*
 * GET /friendrequests
 * Gets a list of friend requests made to the current user.
 */
type ListMyFriendRequestsResponse struct {
    Success bool               `json:"success"`
    Requestors []PublicUser    `json:"requestors"`
}

func listMyFriendRequestsEndpoint(user User) ListMyFriendRequestsResponse {
    var requestors Users
    requestors = user.getFriendRequests()

    resp := ListMyFriendRequestsResponse{
        Success:    true,
        Requestors: requestors.toPublic(),
    }

    return resp
}

/*
 * /friendrequests/{requestorId} endpoint
 */

func myFriendRequestHandler(w http.ResponseWriter, r *http.Request) int {
    log.Println("Handling /friendrequests/{requestorId}")
    user, ok := getCurrentUser(r)
    if !ok {
        return http.StatusUnauthorized
    }

    vars := mux.Vars(r)
    requestorId, err := strconv.Atoi(vars["requestorId"])
    if err != nil || requestorId <= 0 {
        log.Println("Requestor User ID not positive integer")
        return http.StatusBadRequest
    }

    var resp interface{}

    switch r.Method {
    case "DELETE":
        resp = modifyMyFriendRequestEndpoint(user, requestorId, "decline")
    case "PUT":
        resp = modifyMyFriendRequestEndpoint(user, requestorId, "accept")
    default:
        return http.StatusMethodNotAllowed
    }

    sendJSONResponse(w, resp)
    return http.StatusOK
}

/*
 * PUT /friendrequests/{requestorId}
 * Accepts a friend request from the supplied user to the current user.
 */

 /*
 * DELETE /friendrequests/{requestorId}
 * Declines a friend request from the supplied user to the current user.
 */
type ModifyMyFriendRequestResponse struct {
    Success bool    `json:"success"`
    Error   string  `json:"error"`
}

func modifyMyFriendRequestEndpoint(user User, requestorId int, action string) ModifyMyFriendRequestResponse {
    // check if the current user and the specified user are the same
    if requestorId == user.Id {
        return ModifyMyFriendRequestResponse{
            Success:    false,
            Error:      "Requestor User ID cannot be your own",
        }
    }

    // get the user from the ID
    var requestor User
    dbErr := db.Where(&User{Id: requestorId}).First(&requestor).Error

    // check if the user exists
    if dbErr != nil {
        return ModifyMyFriendRequestResponse{
            Success:    false,
            Error:      "User not found",
        }
    }

    // check if the request exists
    if !user.hasFriendRequest(requestor) {
        return ModifyMyFriendRequestResponse{
            Success:    false,
            Error:      "User has not requested to be your friend",
        }
    }

    if action == "accept" {
        // add the friend
        if err := user.addFriend(requestor); err != nil {
            return ModifyMyFriendRequestResponse{
                Success:    false,
                Error:      err.Error(),
            }
        }
    }

    // delete the request
    db.Where("user_id = ? and requestor_id = ?", user.Id, requestor.Id).Delete(FriendRequest{})

    return ModifyMyFriendRequestResponse{
        Success:    true,
    }
}

/*
 * /users/{userId}/friendrequests endpoint
 */

func othersFriendRequestHandler(w http.ResponseWriter, r *http.Request) int {
    log.Println("Handling /users/{userId}/friendrequests")
    user, ok := getCurrentUser(r)
    if !ok {
        return http.StatusUnauthorized
    }

    vars := mux.Vars(r)
    userId, err := strconv.Atoi(vars["userId"])
    if err != nil || userId <= 0 {
        log.Println("User ID not positive integer")
        return http.StatusBadRequest
    }

    var resp interface{}

    switch r.Method {
    case "POST":
        resp = addOthersFriendRequestEndpoint(user, userId)
    default:
        return http.StatusMethodNotAllowed
    }

    sendJSONResponse(w, resp)
    return http.StatusOK
}

/*
 * POST /users/{userId}/friendrequests
 * Sends a friend request from the current user to the supplied user.
 */

type AddOthersFriendRequestResponse struct {
    Success bool        `json:"success"`
    Error   string      `json:"error"`
}

func addOthersFriendRequestEndpoint(user User, requestedId int) AddOthersFriendRequestResponse {
    var requestedFriend User
    dbErr := db.Where(&User{Id: requestedId}).First(&requestedFriend).Error

    if dbErr != nil {
        // friend they are requesting not found
        return AddOthersFriendRequestResponse{
            Success: false,
            Error:   "User not found"}
    }

    // check if they are already friends
    if requestedFriend.isFriend(user) {
        return AddOthersFriendRequestResponse{
            Success:    false,
            Error:      "User is already your friend",
        }
    }

    // check if the request exists
    if requestedFriend.hasFriendRequest(user) {
        return AddOthersFriendRequestResponse{
            Success:    false,
            Error:      "User already has a friend request from you",
        }
    }

    // check if the opposite request exists
    if user.hasFriendRequest(requestedFriend) {
        return AddOthersFriendRequestResponse{
            Success:    false,
            Error:      "You already have a friend request from that user",
        }
    }
    
    addErr := requestedFriend.addFriendRequest(user)

    if addErr != nil {
        return AddOthersFriendRequestResponse{
            Success: false,
            Error:   addErr.Error()}
    }

    return AddOthersFriendRequestResponse{Success: true}
}
