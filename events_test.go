package main

import (
    "testing"
    "log"
    "time"
)

func TestEvents(t *testing.T) {
    defer resetTables()

    user1 := User{
        Id:         1000,
        Uid:        "1000",
        Name:       "Tony Abbott",
        FirstName:  "Tony",
        LastName:   "Abbott",
        Email:      "xXx_0n10n_fan_xXx@hotmail.com",
        Picture:    "tone.jpg",
    }
    db.Create(&user1)

    user2 := User{
        Id:         1001,
        Uid:        "1001",
        Name:       "Malcolm Turnbull",
        FirstName:  "Malcolm",
        LastName:   "Turnbull",
        Email:      "pm@gmail.com",
        Picture:    "hehe",
    }
    db.Create(&user2)

    user1.addFriend(user2)
    
    msg1A, _ := user2.addMessageToUser(user1, "soz", ContentTypeText)
    msg1B, _ := user1.addMessageToUser(user2, "malcom pls", ContentTypeText)
    msg2A, _ := user2.addMessageToUser(user1, "lel", ContentTypeText)
    msg2B, _ := user1.addMessageToUser(user2, "not lel", ContentTypeText)
    msg3A, _ := user2.addMessageToUser(user1, "idc", ContentTypeText)
    msg3B, _ := user1.addMessageToUser(user2, "i h8 u", ContentTypeText)

    msg4, _ := user2.addMessageToUser(user1, "top kek", ContentTypeText)
    msg5, _ := user2.addMessageToUser(user1, "bye", ContentTypeText)

    timeout := 100 * time.Millisecond
    sendWait := 100 * time.Millisecond

    log.Println("** Testing sequential message events for user1 and user2")

    done1 := make(chan bool)
    done2 := make(chan bool)

    go func() {
        msg, timedOut := waitForMessageEvent(user1.Id)
        log.Printf("[1] A: msg %v, timedOut %v\n", msg.Id, timedOut)
        done1 <- (msg.Id == msg1A.Id && msg.RecipientId == user1.Id && !timedOut)
    }()
    go func() {
        msg, timedOut := waitForMessageEvent(user2.Id)
        log.Printf("[2] A: msg %v, timedOut %v\n", msg.Id, timedOut)
        done2 <- (msg.Id == msg1B.Id && msg.RecipientId == user2.Id && !timedOut)
    }()

    log.Println("Awaiting user1 response (should time out)")
    select {
    case <-done1:
        t.Errorf("waitForMessageEvent returned too quickly for user1")
    case <-time.After(timeout):
        log.Println("Timed out successfully for user1")
    }

    log.Println("Awaiting user2 response (should time out)")
    select {
    case <-done2:
        t.Errorf("waitForMessageEvent returned too quickly for user2")
    case <-time.After(timeout):
        log.Println("Timed out successfully for user2")
    }

    time.Sleep(sendWait)
    sendMessageEvent(user1.Id, msg1A)

    log.Println("Awaiting user1 response (should receive)")
    select {
    case <-done1:
        log.Println("Received successfully for user1")
    case <-time.After(timeout):
        t.Errorf("waitForMessageEvent returned too slowly for user1")
    }

    log.Println("Awaiting user2 response (should time out)")
    select {
    case <-done2:
        t.Errorf("waitForMessageEvent returned too quickly for user2")
    case <-time.After(timeout):
        log.Println("Timed out successfully for user2")
    }

    time.Sleep(sendWait)
    sendMessageEvent(user2.Id, msg1B)

    log.Println("Awaiting user2 response (should receive)")
    select {
    case <-done2:
        log.Println("Received successfully for user2")
    case <-time.After(timeout):
        t.Errorf("waitForMessageEvent returned too slowly for user2")
    }

    close(done1)
    close(done2)

    log.Println("** Testing message event to user1 and user2 simultaneously")

    done1 = make(chan bool)
    done2 = make(chan bool)

    go func() {
        msg, timedOut := waitForMessageEvent(user1.Id)
        log.Printf("[1] B: msg %v, timedOut %v\n", msg.Id, timedOut)
        done1 <- (msg.Id == msg2A.Id && msg.RecipientId == user1.Id && !timedOut)
    }()
    go func() {
        msg, timedOut := waitForMessageEvent(user2.Id)
        log.Printf("[2] B: msg %v, timedOut %v\n", msg.Id, timedOut)
        done2 <- (msg.Id == msg2B.Id && msg.RecipientId == user2.Id && !timedOut)
    }()

    time.Sleep(sendWait)
    sendMessageEvent(user1.Id, msg2A)
    sendMessageEvent(user2.Id, msg2B)

    log.Println("Awaiting user1 response (should receive)")
    select {
    case ok := <-done1:
        if ok {
            log.Println("Received successfully for user1")
        } else {
            t.Errorf("Message event wasn't received correctly for user1")
        }
    case <-time.After(timeout):
        t.Errorf("waitForMessageEvent returned too slowly for user1")
    }

    log.Println("Awaiting user2 response (should receive)")
    select {
    case ok := <-done2:
        if ok {
            log.Println("Received successfully for user2")
        } else {
            t.Errorf("Message event wasn't received correctly for user2")
        }
    case <-time.After(timeout):
        t.Errorf("waitForMessageEvent returned too slowly for user2")
    }

    close(done1)
    close(done2)

    log.Println("** Testing multiple listeners for user1")

    done1A := make(chan bool)
    done1B := make(chan bool)
    done2 = make(chan bool)

    go func() {
        msg, timedOut := waitForMessageEvent(user1.Id)
        log.Printf("[1] C1: msg %v, timedOut %v\n", msg.Id, timedOut)
        done1A <- (msg.Id == msg3A.Id && msg.RecipientId == user1.Id && !timedOut)
    }()
    go func() {
        msg, timedOut := waitForMessageEvent(user1.Id)
        log.Printf("[1] C2: msg %v, timedOut %v\n", msg.Id, timedOut)
        done1B <- (msg.Id == msg3A.Id && msg.RecipientId == user1.Id && !timedOut)
    }()
    go func() {
        msg, timedOut := waitForMessageEvent(user2.Id)
        log.Printf("[2] C: msg %v, timedOut %v\n", msg.Id, timedOut)
        done2 <- (msg.Id == msg3B.Id && msg.RecipientId == user2.Id && !timedOut)
    }()

    log.Println("Awaiting user1 response (A) (should time out)")
    select {
    case <-done1A:
        t.Errorf("Shouldn't have received an event for user1 (A)")
    case <-time.After(timeout):
        log.Println("Timed out successfully for user1 (A)")
    }

    log.Println("Awaiting user1 response (B) (should time out)")
    select {
    case <-done1B:
        t.Errorf("Shouldn't have received an event for user1 (B)")
    case <-time.After(timeout):
        log.Println("Timed out successfully for user1 (B)")
    }

    log.Println("Awaiting user2 response (should time out)")
    select {
    case <-done2:
        t.Errorf("Shouldn't have received an event for user2")
    case <-time.After(timeout):
        log.Println("Timed out successfully for user2")
    }

    time.Sleep(sendWait)
    sendMessageEvent(user1.Id, msg3A)

    log.Println("Awaiting user1 (A) response (should receive)")
    select {
    case ok := <-done1A:
        if ok {
            log.Println("Done successfully for user1 (A)")
        } else {
            t.Errorf("Message event wasn't received correctly for user1 (A)")
        }
    case <-time.After(timeout):
        t.Errorf("waitForMessageEvent returned too slowly for user1 (A)")
    }

    log.Println("Awaiting user1 (B) response (should receive)")
    select {
    case ok := <-done1B:
        if ok {
            log.Println("Done successfully for user1 (B)")
        } else {
            t.Errorf("Message event wasn't received correctly for user1 (B)")
        }
    case <-time.After(timeout):
        t.Errorf("waitForMessageEvent returned too slowly for user1 (B)")
    }

    log.Println("Awaiting user2 response (should time out)")
    select {
    case <-done2:
        t.Errorf("Shouldn't have received an event for user2")
    case <-time.After(timeout):
        log.Println("Timed out successfully for user2")
    }

    time.Sleep(sendWait)
    sendMessageEvent(user2.Id, msg3B)

    log.Println("Awaiting user2 response (should receive)")
    select {
    case ok := <-done2:
        if ok {
            log.Println("Done successfully for user2")
        } else {
            t.Errorf("Message event wasn't received correctly for user2")
        }
    case <-time.After(timeout):
        t.Errorf("waitForMessageEvent returned too slowly for user2")
    }

    close(done1A)
    close(done1B)
    close(done2)

    log.Println("** Testing getNextMessageEndpoint")

    done1 = make(chan bool)

    log.Println("Testing afterId = 0")
    go func() {
        resp := getNextMessageEndpoint(user1, 0)
        log.Printf("[1] D1: success %v, error %v, msg %v\n", resp.Success, resp.Error, resp.Message.Id)
        done1 <- (resp.Success && resp.Message.Id == msg4.Id)
    }()

    log.Println("Awaiting response (should time out)")
    select {
    case <-done1:
        t.Errorf("Shouldn't have received a response")
    case <-time.After(timeout):
        log.Println("Timed out successfully")
    }

    time.Sleep(sendWait)
    sendMessageEvent(user1.Id, msg4)

    log.Println("Awaiting response (should receive)")
    select {
    case ok := <-done1:
        if ok {
            log.Println("Response successful")
        } else {
            t.Errorf("Response was unsuccessful")
        }
    case <-time.After(timeout):
        t.Errorf("Response wasn't received in time")
    }

    close(done1)
    done1 = make(chan bool)

    log.Printf("Testing afterId = msg4.Id (%v)\n", msg4.Id)
    go func() {
        resp := getNextMessageEndpoint(user1, msg4.Id)
        log.Printf("[1] D2: success %v, error %v, msg %v\n", resp.Success, resp.Error, resp.Message.Id)
        done1 <- (resp.Success && resp.Message.Id == msg5.Id)
    }()

    log.Println("Awaiting response (should receive)")
    select {
    case ok := <-done1:
        if ok {
            log.Println("Response successful")
        } else {
            t.Errorf("Response was unsuccessful")
        }
    case <-time.After(timeout):
        t.Errorf("Response wasn't received in time")
    }

    close(done1)
    done1 = make(chan bool)

    log.Println("Adding msg6")
    msg6, _ := user2.addMessageToUser(user1, "glhf", ContentTypeText)

    log.Printf("Testing afterId = msg5.Id (%v)\n", msg5.Id)
    go func() {
        resp := getNextMessageEndpoint(user1, msg5.Id)
        log.Printf("[1] D3: success %v, error %v, msg %v\n", resp.Success, resp.Error, resp.Message.Id)
        done1 <- (resp.Success && resp.Message.Id == msg6.Id)
    }()

    log.Println("Awaiting response (should receive)")
    select {
    case ok := <-done1:
        if ok {
            log.Println("Response successful")
        } else {
            t.Errorf("Response was unsuccessful")
        }
    case <-time.After(timeout):
        t.Errorf("Response wasn't received in time")
    }

    close(done1)
    done1 = make(chan bool)

    nextId := msg6.Id + 1

    log.Printf("Testing afterId = msg6.Id (%v)\n", msg6.Id)
    go func() {
        resp := getNextMessageEndpoint(user1, msg6.Id)
        log.Printf("[1] D4: success %v, error %v, msg %v\n", resp.Success, resp.Error, resp.Message.Id)
        done1 <- (resp.Success && resp.Message.Id == nextId)
    }()

    time.Sleep(sendWait)
    log.Println("Adding msg7")
    msg7, _ := user2.addMessageToUser(user1, "gg", ContentTypeText)
    if msg7.Id != nextId {
        t.Errorf("msg7.Id wasn't what we expected: expected %v, got %v; correct assumption!", nextId, msg7.Id)
    }

    log.Println("Awaiting response (should time out)")
    select {
    case <-done1:
        t.Errorf("Shouldn't have received a response")
    case <-time.After(timeout):
        log.Println("Timed out successfully")
    }

    time.Sleep(sendWait)
    sendMessageEvent(user1.Id, msg7)

    log.Println("Awaiting response (should receive)")
    select {
    case ok := <-done1:
        if ok {
            log.Println("Response successful")
        } else {
            t.Errorf("Response was unsuccessful")
        }
    case <-time.After(timeout):
        t.Errorf("Response wasn't received in time")
    }

    close(done1)

    log.Println("** Testing sendMessageEndpoint sending events")

    done1 = make(chan bool)
    
    nextId = msg7.Id + 1
    req := SendMessageRequest{
        Content:        "sorry not sorry",
        ContentType:    ContentTypeText,
    }

    log.Println("Testing afterId = 0")
    go func() {
        resp := getNextMessageEndpoint(user1, 0)
        log.Printf("[1] E: success %v, error %v, msg %v\n", resp.Success, resp.Error, resp.Message.Id)
        done1 <- (resp.Success && resp.Message.Id == nextId)
    }()

    log.Println("Awaiting response (should time out)")
    select {
    case <-done1:
        t.Errorf("Shouldn't have received a response")
    case <-time.After(timeout):
        log.Println("Timed out successfully")
    }

    time.Sleep(sendWait)
    resp := sendMessageEndpoint(user2, user1.Id, req)
    if !resp.Success {
        t.Errorf("Send message failed: %v", resp.Error)
    }
    if resp.Id != nextId {
        t.Errorf("resp.Id wasn't what we expected: expected %v, got %v; correct assumption!", nextId, resp.Id)
    }

    log.Println("Awaiting response (should receive)")
    select {
    case ok := <-done1:
        if ok {
            log.Println("Response successful")
        } else {
            t.Errorf("Response was unsuccessful")
        }
    case <-time.After(timeout):
        t.Errorf("Response wasn't received in time")
    }

    close(done1)

    log.Println("** Event tests should be done now")
}
