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

    TimeoutChan *chan bool
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
    listener.TimeoutChan = nil

    listener.Lock.Unlock()
}

// Waits until an event is received (or timeout).
func waitForEvent(userId int) (eventType int, event interface{}, timedOut bool) {
    listener := getListener(userId)

    // spin off goroutine for wait loop
    done := make(chan bool)
    go func() {
        for {
            listener.Lock.Lock()
            listener.Cond.Wait()

            if listener.TimeoutChan == nil {
                // legit signal!
                eventType = listener.EventType
                event = listener.Event

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

        return eventType, event, false
    case <-time.After(time.Second * EventTimeout):
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
            return eventType, event, false
        } else {
            log.Println("Timeout successful")
            return 0, nil, true
        }
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
