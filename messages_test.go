package main

import (
    "testing"
    "log"
)

func TestGetSender(t *testing.T) {
    defer resetTables()
    user1 := User{
        Id:         1,
        Uid:        "1",
        Name:       "Snoop Doge",
        FirstName:  "Snoop",
        LastName:   "Doge",
        Email:      "poop@gmail.com",
        Picture:    "blah",
    }
    db.Create(&user1)

    user2 := User{
        Id:         2,
        Uid:        "2",
        Name:       "Malcolm Turnbull",
        FirstName:  "Malcolm",
        LastName:   "Turnbull",
        Email:      "pm@gmail.com",
        Picture:    "hehe",
    }
    db.Create(&user2)

    log.Println("Add a message from user1 to user2")
    user1.addMessageToUser(user2, "this is a message from user1 to user2", 1)
    
    log.Println("Get messages between user1 and user2")
    messages := user1.getMessagesWithUser(user2)

    log.Println("Check sender")
    if sender, _ := messages[0].getSender(); sender.Id != user1.Id {
        t.Error("Wrong sender returned")
    }
}

func TestGetRecipientUser(t *testing.T) {
    defer resetTables()
    user1 := User{
        Id:         1,
        Uid:        "1",
        Name:       "Snoop Doge",
        FirstName:  "Snoop",
        LastName:   "Doge",
        Email:      "poop@gmail.com",
        Picture:    "blah",
    }
    db.Create(&user1)

    user2 := User{
        Id:         2,
        Uid:        "2",
        Name:       "Malcolm Turnbull",
        FirstName:  "Malcolm",
        LastName:   "Turnbull",
        Email:      "pm@gmail.com",
        Picture:    "hehe",
    }
    db.Create(&user2)

    log.Println("Add a message from user1 to user2")
    user1.addMessageToUser(user2, "this is a message from user1 to user2", 1)
    
    log.Println("Get messages between user1 and user2")
    messages := user1.getMessagesWithUser(user2)

    log.Println("Check recipient")
    if recipient, _ := messages[0].getRecipientUser(); recipient.Id != user2.Id {
        t.Error("Wrong recipient returned")
    }
}

func TestListMessagesEndpoint(t *testing.T) {
    defer resetTables()

    user1 := User{
        Id:         1,
        Uid:        "1",
        Name:       "Snoop Doge",
        FirstName:  "Snoop",
        LastName:   "Doge",
        Email:      "poop@gmail.com",
        Picture:    "blah",
    }
    db.Create(&user1)

    user2 := User{
        Id:         2,
        Uid:        "2",
        Name:       "Malcolm Turnbull",
        FirstName:  "Malcolm",
        LastName:   "Turnbull",
        Email:      "pm@gmail.com",
        Picture:    "hehe",
    }
    db.Create(&user2)

    user3 := User{
        Id:         3,
        Uid:        "3",
        Name:       "Shrek",
        FirstName:  "Shrek",
        LastName:   "The Ogre",
        Email:      "swamp@gmail.com",
        Picture:    "40keks",
    }
    db.Create(&user3)

    log.Println("Add a message from user1 to user2")
    user1.addMessageToUser(user2, "this is a message from user1 to user2", 1)

    log.Println("List the messages between user2 and user1")
    response := listMessagesEndpoint(user2, 1)

    if response.Error != "" {
        t.Error("Response had error when it shouldn't have")
    }

    if len(response.Messages) != 1 {
        t.Errorf("1 message expected, found %v\n", len(response.Messages))
    } else {
        if response.Messages[0].Content != "this is a message from user1 to user2" {
            t.Errorf("Message returned had the wrong content: %v\n", response.Messages[0].Content)
        }
        if sender, _ := response.Messages[0].getSender(); sender.Id != 1 {
            t.Errorf("Message returned had the wrong senderid: %v\n", response.Messages[0].SenderId)
        }
        if recipient, _ := response.Messages[0].getRecipientUser(); recipient.Id != 2 {
            t.Errorf("Message returned had the wrong recipientid: %v\n", response.Messages[0].RecipientId)
        }
    }

    log.Println("List the messages between user1 and user2")
    response = listMessagesEndpoint(user1, 2)

    if response.Error != "" {
        t.Error("Response had error when it shouldn't have")
    }

    if len(response.Messages) != 1 {
        t.Errorf("1 message expected, found %v\n", len(response.Messages))
    } else {
        if response.Messages[0].Content != "this is a message from user1 to user2" {
            t.Errorf("Message returned had the wrong content: %v\n", response.Messages[0].Content)
        }
        if sender, _ := response.Messages[0].getSender(); sender.Id != 1 {
            t.Errorf("Message returned had the wrong senderid: %v\n", response.Messages[0].SenderId)
        }
        if recipient, _ := response.Messages[0].getRecipientUser(); recipient.Id != 2 {
            t.Errorf("Message returned had the wrong recipientid: %v\n", response.Messages[0].RecipientId)
        }
    }

    log.Println("Add a message from user2 to user3")
    user2.addMessageToUser(user3, "this is a message from user2 to user3", 1)

    log.Println("List the messages between user3 and user2")
    response = listMessagesEndpoint(user3, 2)

    if response.Error != "" {
        t.Error("Response had error when it shouldn't have")
    }

    if len(response.Messages) != 1 {
        t.Errorf("1 message expected, found %v\n", len(response.Messages))
    } else {
        if response.Messages[0].Content != "this is a message from user2 to user3" {
            t.Errorf("Message returned had the wrong content: %v\n", response.Messages[0].Content)
        }
        if sender, _ := response.Messages[0].getSender(); sender.Id != 2 {
            t.Errorf("Message returned had the wrong senderid: %v\n", response.Messages[0].SenderId)
        }
        if recipient, _ := response.Messages[0].getRecipientUser(); recipient.Id != 3 {
            t.Errorf("Message returned had the wrong recipientid: %v\n", response.Messages[0].RecipientId)
        }
    }

    log.Println("Add a message from user1 to user2")
    user1.addMessageToUser(user2, "this is another message from user1 to user2", 1)

    log.Println("List the messages between user2 and user1")
    response = listMessagesEndpoint(user2, 1)

    if response.Error != "" {
        t.Error("Response had error when it shouldn't have")
    }

    if len(response.Messages) != 2 {
        t.Errorf("2 messages expected, found %v\n", len(response.Messages))
    } else {
        if response.Messages[0].Content != "this is a message from user1 to user2" {
            t.Errorf("Message returned had the wrong content: %v\n", response.Messages[0].Content)
        }
        if response.Messages[1].Content != "this is another message from user1 to user2" {
            t.Errorf("Message returned had the wrong content: %v\n", response.Messages[1].Content)
        }
    }

    log.Println("Add a message from user2 to user1")
    user1.addMessageToUser(user2, "this is a message from user2 to user1", 1)

    log.Println("List the messages between user1 and user2")
    response = listMessagesEndpoint(user1, 2)

    if response.Error != "" {
        t.Error("Response had error when it shouldn't have")
    }

    if len(response.Messages) != 3 {
        t.Errorf("1 message expected, found %v\n", len(response.Messages))
    } else {
        if response.Messages[0].Content != "this is a message from user1 to user2" {
            t.Errorf("Message returned had the wrong content: %v\n", response.Messages[0].Content)
        }
        if response.Messages[1].Content != "this is another message from user1 to user2" {
            t.Errorf("Message returned had the wrong content: %v\n", response.Messages[1].Content)
        }
        if response.Messages[2].Content != "this is a message from user2 to user1" {
            t.Errorf("Message returned had the wrong content: %v\n", response.Messages[2].Content)
        }
    }

    log.Println("List the messages between user1 and a non existent user")
    response = listMessagesEndpoint(user1, 123)
    if response.Error != "Friend not found" {
        t.Errorf("Response returned the wrong error. Got error %v\n", response.Error)
    }
}

func TestListMessagesEndpoint2(t *testing.T) {
    defer resetTables()

    user1 := User{
        Id:        420,
        Uid:       "420",
        Name:      "Snoop Doge",
        FirstName: "Snoop",
        LastName:  "Doge",
        Email:     "higher@gmail.com",
        Picture:   "42keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    user2 := User{
        Id:        421,
        Uid:       "421",
        Name:      "Peppa Pig",
        FirstName: "Peppa",
        LastName:  "Pig",
        Email:     "p.pig@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    log.Println("Getting lists of messages (expecting it to be empty)")
    resp1 := listMessagesEndpoint(user1, user2.Id)
    resp2 := listMessagesEndpoint(user2, user1.Id)
    if len(resp1.Messages) != 0 {
        t.Errorf("Wrong number of messages; expected 0, found %v\n", len(resp1.Messages))
    }
    if len(resp2.Messages) != 0 {
        t.Errorf("Wrong number of messages; expected 0, found %v\n", len(resp2.Messages))
    }

    log.Println("Adding message")
    if err := user1.addMessageToUser(user2, "hello", ContentTypeText); err != nil {
        t.Errorf("Failed to add message")
    }

    log.Println("Getting lists of messages (expecting it to have 1)")
    resp1 = listMessagesEndpoint(user1, user2.Id)
    resp2 = listMessagesEndpoint(user2, user1.Id)
    if len(resp1.Messages) != 1 {
        t.Errorf("Wrong number of messages; expected 1, found %v\n", len(resp1.Messages))
    }
    if len(resp2.Messages) != 1 {
        t.Errorf("Wrong number of messages; expected 1, found %v\n", len(resp2.Messages))
    }

    log.Println("Adding another message")
    if err := user2.addMessageToUser(user1, "yo", ContentTypeText); err != nil {
        t.Errorf("Failed to add message")
    }

    log.Println("Getting lists of messages (expecting it to have 2)")
    resp1 = listMessagesEndpoint(user1, user2.Id)
    resp2 = listMessagesEndpoint(user2, user1.Id)
    if len(resp1.Messages) != 2 {
        t.Errorf("Wrong number of messages; expected 2, found %v\n", len(resp1.Messages))
    }
    if len(resp2.Messages) != 2 {
        t.Errorf("Wrong number of messages; expected 2, found %v\n", len(resp2.Messages))
    }
}

func TestSendMessageEndpoint(t *testing.T) {
    defer resetTables()

    user1 := User{
        Id:         1,
        Uid:        "1",
        Name:       "Snoop Doge",
        FirstName:  "Snoop",
        LastName:   "Doge",
        Email:      "poop@gmail.com",
        Picture:    "blah",
    }
    db.Create(&user1)

    user2 := User{
        Id:         2,
        Uid:        "2",
        Name:       "Malcolm Turnbull",
        FirstName:  "Malcolm",
        LastName:   "Turnbull",
        Email:      "pm@gmail.com",
        Picture:    "hehe",
    }
    db.Create(&user2)

    log.Println("Send a message from user1 to user2")
    req := SendMessageRequest{
        Content:     "You are a nice person",
        ContentType: ContentTypeText,
    }
    resp := sendMessageEndpoint(user1, 2, req)

    if !resp.Success {
        t.Error("Response returned not success when it should have been successful")
    }

    if resp.Error != "" {
        t.Error("Response returned an error when it shouldn't have")
    }

    messages := user1.getMessagesWithUser(user2)

    if len(messages) != 1 {
        t.Errorf("1 message expected, found %v\n", len(messages))
    } else {
        if messages[0].Content != "You are a nice person" {
            t.Errorf("Message returned had the wrong content: %v\n", messages[0].Content)
        }
        if sender, _ := messages[0].getSender(); sender.Id != 1 {
            t.Errorf("Message returned had the wrong senderid: %v\n", messages[0].SenderId)
        }
        if recipient, _ := messages[0].getRecipientUser(); recipient.Id != 2 {
            t.Errorf("Message returned had the wrong recipientid: %v\n", messages[0].RecipientId)
        }
    }

    log.Println("Send a message from user1 to a non existent user")
    req = SendMessageRequest{
        Content:     "You are a nice person",
        ContentType: ContentTypeText,
    }
    resp = sendMessageEndpoint(user1, 1234, req)

    if resp.Success {
        t.Error("Response returned success when it should have been unsuccessful")
    }

    if resp.Error != "Friend not found" {
        t.Errorf("Response returned the wrong error. Got error %v\n", resp.Error)
    }

    log.Println("Send a message from user1 to user2 with an invalid content type")
    req = SendMessageRequest{
        Content:     "You are a nice person",
        ContentType: 2,
    }
    resp = sendMessageEndpoint(user1, 2, req)

    if resp.Success {
        t.Error("Response returned success when it should have been unsuccessful")
    }

    if resp.Error != "Invalid content type" {
        t.Errorf("Response returned the wrong error. Got error %v\n", resp.Error)
    }
}
