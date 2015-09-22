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

type ContentType int
const (
    ContentTypeText = 1
)

func (ct *ContentType) valid() bool {
    return ct != nil && *ct >= 1 && *ct <= 1
}

type RecipientType int
const (
    RecipientTypeUser = 1
)

func (rt *RecipientType) valid() bool {
    return rt != nil && *rt >= 1 && *rt <= 1
}

type Message struct {
    Id                  int             `json:"id" gorm:"primary_key" sql:"auto_increment"`
    Content             string          `json:"content" sql:"type:varchar(1024)"`
    ContentType         ContentType     `json:"contentType" sql:"not null"`
    SenderId            int             `json:"senderId" sql:"not null"`
    RecipientId         int             `json:"recipientId" sql:"not null"`
    RecipientType       RecipientType   `json:"recipientType" sql:"not null"`
    Timestamp           time.Time       `json:"timestamp" sql:"not null"`
}

type Messages []Message

func (msg *Message) getSender() (sender User, err error) {
    err = db.Where(&User{Id: msg.SenderId}).First(&sender).Error
    return
}

func (msg *Message) getRecipientUser() (recipient User, err error) {
    if msg.RecipientType != RecipientTypeUser {
        err = errors.New("Invalid recipient type")
    } else {
        err = db.Where(&User{Id: msg.RecipientId}).First(&recipient).Error
    }
    return
}

/*
 * API endpoints
 */

/*
 * /messages endpoint
 */

func messagesHandler(w http.ResponseWriter, r *http.Request) int {
    log.Println("Handling /messages")
    user, ok := getCurrentUser(r)
    if !ok {
        return http.StatusUnauthorized
    }

    vars := mux.Vars(r)
    friendId, err := strconv.Atoi(vars["friendId"])
    if err != nil {
        log.Println("Friend ID not integer")
        return http.StatusBadRequest
    }

    var resp interface{}

    switch r.Method {
    case "GET":
        resp = listMessagesEndpoint(user, friendId)
    case "POST":
        decoder := json.NewDecoder(r.Body)
        var req SendMessageRequest
        err := decoder.Decode(&req)
        if err != nil {
            log.Println("JSON decoding failed")
            return http.StatusBadRequest
        }
        resp = sendMessageEndpoint(user, friendId, req)
    default:
        return http.StatusMethodNotAllowed
    }

    sendJSONResponse(w, resp)
    return http.StatusOK
}

/*
 * GET /messages/{friendId}
 * Gets a list of the messages between the current user and the specified friend.
 */
type ListMessagesResponse struct {
    Messages Messages    `json:"messages"`
    Error    string      `json:"error"`
}

func listMessagesEndpoint(user User, friendId int) ListMessagesResponse {
    var friend User
    dbErr := db.Where(&User{Id: friendId}).First(&friend).Error
    

    if dbErr != nil {
        // friend they are trying to list messages between not found'
        return ListMessagesResponse{
                  Error:   "Friend not found"}
    }

    var messages Messages
    messages = user.getMessagesWithUser(friend)

    return ListMessagesResponse{
        Messages: messages,
    }
}

/*
 * POST /messages/{friendId}
 * Sends a message from the current user to the specified friend
 */
type SendMessageRequest struct {
    Content     string      `json:"content"`
    ContentType ContentType `json:"contentType"`
}

type SendMessageResponse struct {
    Success bool        `json:"success"`
    Error   string      `json:"error"`
}

func sendMessageEndpoint(user User, friendId int, req SendMessageRequest) SendMessageResponse {
    var friend User
    dbErr := db.Where(&User{Id: friendId}).First(&friend).Error

    if dbErr != nil {
        // friend they are trying to send message to not found
        // TODO: check if they are also a friend of that user
        return SendMessageResponse{
            Success: false,
            Error:   "Friend not found"}
    }

    sendErr := user.addMessageToUser(friend, req.Content, req.ContentType)

    if sendErr != nil {
        return SendMessageResponse{
            Success: false,
            Error:   sendErr.Error()}
    }

    return SendMessageResponse{Success: true}
}
