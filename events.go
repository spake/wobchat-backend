package main

import (
    "log"
    "net/http"
    "sync"
)

type EventListener struct {
    Lock        sync.Mutex
    Cond        *sync.Cond
    EventType   int
    Event       interface{}
}

var listeners map[int]*EventListener

func eventHandler(w http.ResponseWriter, r *http.Request) int {
    log.Println("Handling /nextEvent")
    user, ok := getCurrentUser(r)
    if !ok {
        return http.StatusUnauthorized
    }

    var resp interface{}

    if listeners == nil {
        listeners = make(map[int]*EventListener)
    }

    listener := listeners[user.Id]
    if listener == nil {
        listener = new(EventListener)
        listener.Cond = sync.NewCond(&listener.Lock)
        listeners[user.Id] = listener
    }

    switch r.Method {
    case "GET":
        resp = getEventEndpoint(user, listener)
    case "POST":
        resp = postEventEndpoint(user, listener)
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

func getEventEndpoint(user User, listener *EventListener) GetEventResponse {
    log.Println("Waiting for event")

    listener.Lock.Lock()
    listener.Cond.Wait()
    eventType := listener.EventType
    event := listener.Event
    listener.Lock.Unlock()

    return GetEventResponse{
        Success:    true,
        Error:      "",
        EventType:  eventType,
        Event:      event,
    }
}

type PostEventResponse struct {
    Success bool    `json:"success"`
    Error   string  `json:"error"`
}

func postEventEndpoint(user User, listener *EventListener) PostEventResponse {
    log.Println("Sending message")

    msg := Message{
        Content:        "condition variables are gr8",
        ContentType:    ContentTypeText,
        SenderId:       0,
        RecipientId:    user.Id,
        RecipientType:  RecipientTypeUser,
    }

    event := EventNewMessage{
        Message:    msg,
    }

    listener.Lock.Lock()
    listener.Cond.Broadcast()
    listener.EventType = EventTypeNewMessage
    listener.Event = event
    listener.Lock.Unlock()

    return PostEventResponse{
        Success:    true,
        Error:      "",
    }
}
