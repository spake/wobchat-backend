package main

import (
    "log"
    "net/http"
    "strconv"
    "sync"
    "time"
)

const MessageEventTimeout = 60

type MessageEventListener struct {
    Lock        sync.Mutex
    Cond        *sync.Cond
    Message     Message

    TimeoutChan *chan bool
}

var listenersLock   sync.Mutex
var listeners       map[int]*MessageEventListener

func getListener(userId int) *MessageEventListener {
    listenersLock.Lock()

    if listeners == nil {
        listeners = make(map[int]*MessageEventListener)
    }

    listener := listeners[userId]
    if listener == nil {
        listener = new(MessageEventListener)
        listener.Cond = sync.NewCond(&listener.Lock)
        listeners[userId] = listener
    }

    listenersLock.Unlock()

    return listener
}

// Sends an event to the given user.
func sendMessageEvent(userId int, message Message) {
    listener := getListener(userId)

    listener.Lock.Lock()
    listener.Cond.Broadcast()

    listener.Message = message
    listener.TimeoutChan = nil

    listener.Lock.Unlock()
}

// Waits until an event is received (or timeout).
func waitForMessageEvent(userId int) (message Message, timedOut bool) {
    listener := getListener(userId)

    // spin off goroutine for wait loop
    done := make(chan bool)
    go func() {
        for {
            listener.Lock.Lock()
            listener.Cond.Wait()

            if listener.TimeoutChan == nil {
                // legit signal!
                message = listener.Message

                listener.Lock.Unlock()
                done <- true

                break
            } else if listener.TimeoutChan == &done {
                // timeout was for us :(
                
                listener.Lock.Unlock()
                done <- false

                break
            }

            listener.Lock.Unlock()
        }
    }()

    select {
    case <-done:
        log.Println("Received message while waiting")

        return message, false
    case <-time.After(time.Second * MessageEventTimeout):
        log.Println("Timed out while waiting")
        log.Println("Broadcasting timeout signal")

        listener.Lock.Lock()

        listener.TimeoutChan = &done

        listener.Cond.Broadcast()
        listener.Lock.Unlock()

        log.Println("Waiting for response from listener")
        // this will either mean we successfully timed out or we
        // were receiving something while timing out
        if <-done {
            log.Println("Surprise! Received data")
            return message, false
        } else {
            log.Println("Timeout successful")
            return Message{}, true
        }
    }
}

func nextMessageHandler(w http.ResponseWriter, r *http.Request) int {
    log.Println("Handling /nextMessage")
    user, ok := getCurrentUser(r)
    if !ok {
        return http.StatusUnauthorized
    }
    
    var resp interface{}

    switch r.Method {
    case "GET":
        // default to ID of 0
        afterId := 0

        afterIdStr := r.FormValue("after")
        if afterIdStr != "" {
            var err error
            afterId, err = strconv.Atoi(afterIdStr)
            if err != nil || afterId <= 0 {
                log.Println("After ID not positive integer")
                return http.StatusBadRequest
            }
        }

        log.Printf("After ID: %v\n", afterId)
        resp = getNextMessageEndpoint(user, afterId)
    default:
        return http.StatusMethodNotAllowed
    }

    sendJSONResponse(w, resp)
    return http.StatusOK
}

type GetNextMessageResponse struct {
    Success     bool        `json:"success"`
    Error       string      `json:"error"`
    Message     Message     `json:"message"`
}

func getNextMessageEndpoint(user User, afterId int) GetNextMessageResponse {
    // is there already a new message?
    if afterId > 0 {
        message, ok := user.getNextMessageAfterId(afterId)
        if ok {
            log.Printf("Found existing message: %v\n", message.Id)
            return GetNextMessageResponse{
                Success:    true,
                Message:    message,
            }
        }
    }

    // no new messages: long-poll and wait
    log.Println("Waiting for message")

    message, timedOut := waitForMessageEvent(user.Id)

    if timedOut {
        return GetNextMessageResponse{
            Success:    false,
            Error:      "Timed out",
        }
    } else {
        return GetNextMessageResponse{
            Success:    true,
            Message:    message,
        }
    }
}
