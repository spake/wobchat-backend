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
    ContentTypeVideo = 2
    ContentTypeShake = 3
)

func (ct *ContentType) valid() bool {
    return ct != nil && *ct >= 1 && *ct <= 3
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
 * /friends/{friendId}/messages endpoint
 */

func messagesHandler(w http.ResponseWriter, r *http.Request) int {
    log.Println("Handling /friends/{friendId}/messages")
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
 * GET /friends/{friendId}/messages
 * Gets a list of the messages between the current user and their friend specified
 * by the Id.
 */
type ListMessagesResponse struct {
    Success     bool        `json:"success"`
    Error       string      `json:"error"`
    Messages    Messages    `json:"messages"`
}

func listMessagesEndpoint(user User, friendId int) ListMessagesResponse {
    if friendId == user.Id {
        return ListMessagesResponse{
            Success:    false,
            Error:      "You can't list messages from yourself",
        }
    }

    var friend User
    dbErr := db.Where(&User{Id: friendId}).First(&friend).Error

    if dbErr != nil {
        // friend they are trying to list messages between not found
        return ListMessagesResponse{
            Success:    false,
            Error:      "Friend not found",
        }
    }

    if !user.isFriend(friend) {
        return ListMessagesResponse{
            Success:    false,
            Error:      "User is not your friend",
        }
    }

    var messages Messages
    messages = user.getMessagesWithUser(friend)

    return ListMessagesResponse{
        Success:    true,
        Messages:   messages,
    }
}

/*
 * POST /friends/{friendId}/messages
 * Sends a message from the current user to their friend specified by the Id.
 */
type SendMessageRequest struct {
    Content     string      `json:"content"`
    ContentType ContentType `json:"contentType"`
}

type SendMessageResponse struct {
    Success bool        `json:"success"`
    Error   string      `json:"error"`
    Id      int         `json:"id"`
}

func sendMessageEndpoint(user User, friendId int, req SendMessageRequest) SendMessageResponse {
    if friendId == user.Id {
        return SendMessageResponse{
            Success:    false,
            Error:      "You can't send messages to yourself",
        }
    }

    var friend User
    dbErr := db.Where(&User{Id: friendId}).First(&friend).Error

    if dbErr != nil {
        // friend they are trying to send message to not found
        return SendMessageResponse{
            Success:    false,
            Error:      "Friend not found",
        }
    }

    if !user.isFriend(friend) {
        // users are not friends
        return SendMessageResponse{
            Success:    false,
            Error:      "User is not your friend",
        }
    }

    msg, sendErr := user.addMessageToUser(friend, req.Content, req.ContentType)

    if sendErr != nil {
        return SendMessageResponse{
            Success:    false,
            Error:      sendErr.Error(),
        }
    }

    // send event, in case the friend is currently long-polling
    sendMessageEvent(friend.Id, msg)

    return SendMessageResponse{
        Success:    true,
        Id:         msg.Id,
    }
}
