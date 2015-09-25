package main

import (
    "log"
    "net/http"
    "sync"
    "time"
)

const EventTimeout = 60

type EventListener struct {
    Lock        sync.Mutex
    Cond        *sync.Cond
    EventType   int
    Event       interface{}
}

var listenersLock   sync.Mutex
var listeners       map[int]*EventListener

func getListener(userId int) *EventListener {
    listenersLock.Lock()

    if listeners == nil {
        listeners = make(map[int]*EventListener)
    }

    listener := listeners[userId]
    if listener == nil {
        listener = new(EventListener)
        listener.Cond = sync.NewCond(&listener.Lock)
        listeners[userId] = listener
    }

    listenersLock.Unlock()

    return listener
}

// Sends an event to the given user.
func sendEvent(userId int, eventType int, event interface{}) {
    listener := getListener(userId)

    listener.Lock.Lock()
    listener.Cond.Broadcast()
    listener.EventType = eventType
    listener.Event = event
    listener.Lock.Unlock()
}

// Waits until an event is received (or timeout).
func waitForEvent(userId int) (eventType int, event interface{}, timedOut bool) {
    listener := getListener(userId)

    c := make(chan int)

    go func() {
        listener.Lock.Lock()
        listener.Cond.Wait()
        eventType = listener.EventType
        event = listener.Event
        listener.Lock.Unlock()

        c <- 1
    }()

    select {
    case <-c:
        return eventType, event, false
    case <-time.After(time.Second * EventTimeout):
        return eventType, event, true
    }
}

func eventHandler(w http.ResponseWriter, r *http.Request) int {
    log.Println("Handling /nextEvent")
    user, ok := getCurrentUser(r)
    if !ok {
        return http.StatusUnauthorized
    }

    var resp interface{}

    switch r.Method {
    case "GET":
        resp = getEventEndpoint(user)
    case "POST":
        resp = postEventEndpoint(user)
    default:
        return http.StatusMethodNotAllowed
    }

    sendJSONResponse(w, resp)
    return http.StatusOK
}

const (
    EventTypeNewMessage = 1
)

type GetEventResponse struct {
    Success     bool        `json:"success"`
    Error       string      `json:"error"`
    EventType   int         `json:"eventType"`
    Event       interface{} `json:"event"`
}

type EventNewMessage struct {
    Message Message `json:"message"`
}

func getEventEndpoint(user User) GetEventResponse {
    log.Println("Waiting for event")

    eventType, event, timedOut := waitForEvent(user.Id)

    if timedOut {
        return GetEventResponse{
            Success:    false,
            Error:      "Timed out",
        }
    } else {
        return GetEventResponse{
            Success:    true,
            Error:      "",
            EventType:  eventType,
            Event:      event,
        }
    }
}

type PostEventResponse struct {
    Success bool    `json:"success"`
    Error   string  `json:"error"`
}

func postEventEndpoint(user User) PostEventResponse {
    log.Println("Sending message")

    msg := Message{
        Content:        "condition variables are gr8",
        ContentType:    ContentTypeText,
        SenderId:       420,
        RecipientId:    user.Id,
        RecipientType:  RecipientTypeUser,
    }

    event := EventNewMessage{
        Message:    msg,
    }

    sendEvent(user.Id, EventTypeNewMessage, event)

    return PostEventResponse{
        Success:    true,
        Error:      "",
    }
}
