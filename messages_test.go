package main

import (
    "testing"
    "log"
)

//things to fix:
// need to finish TestSendMessageEndpoint

func TestGetSender(t *testing.T) {
    defer resetTables()
    user1 := User{
        Uid:        "1",
        Name:       "Snoop Doge",
        FirstName:  "Snoop",
        LastName:   "Doge",
        Email:      "poop@gmail.com",
        Picture:    "blah",
    }
    db.Create(&user1)

    user2 := User{
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
    if sender, _ := messages[0].getSender(); sender.Uid != user1.Uid {
        t.Error("Wrong sender returned")
    }
}

func TestGetRecipientUser(t *testing.T) {
    defer resetTables()
    user1 := User{
        Uid:        "1",
        Name:       "Snoop Doge",
        FirstName:  "Snoop",
        LastName:   "Doge",
        Email:      "poop@gmail.com",
        Picture:    "blah",
    }
    db.Create(&user1)

    user2 := User{
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
    if recipient, _ := messages[0].getRecipientUser(); recipient.Uid != user2.Uid {
        t.Error("Wrong recipient returned")
    }
}

func TestListMessagesEndpoint(t *testing.T) {
    defer resetTables()

    user1 := User{
        Uid:        "1",
        Name:       "Snoop Doge",
        FirstName:  "Snoop",
        LastName:   "Doge",
        Email:      "poop@gmail.com",
        Picture:    "blah",
    }
    db.Create(&user1)

    user2 := User{
        Uid:        "2",
        Name:       "Malcolm Turnbull",
        FirstName:  "Malcolm",
        LastName:   "Turnbull",
        Email:      "pm@gmail.com",
        Picture:    "hehe",
    }
    db.Create(&user2)

    user3 := User{
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
    response := listMessagesEndpoint(user2, "1")

    if response.Error != "" {
        t.Error("Response had error when it shouldn't have")
    }

    if len(response.Messages) != 1 {
        t.Errorf("1 message expected, found %v\n", len(response.Messages))
    } else {
        if response.Messages[0].Content != "this is a message from user1 to user2" {
            t.Errorf("Message returned had the wrong content: %v\n", response.Messages[0].Content)
        }
        if sender, _ := response.Messages[0].getSender(); sender.Uid != "1" {
            t.Errorf("Message returned had the wrong senderid: %v\n", response.Messages[0].SenderId)
        }
        if recipient, _ := response.Messages[0].getRecipientUser(); recipient.Uid != "2" {
            t.Errorf("Message returned had the wrong recipientid: %v\n", response.Messages[0].RecipientId)
        }
    }

    log.Println("List the messages between user1 and user2")
    response = listMessagesEndpoint(user1, "2")

    if response.Error != "" {
        t.Error("Response had error when it shouldn't have")
    }

    if len(response.Messages) != 1 {
        t.Errorf("1 message expected, found %v\n", len(response.Messages))
    } else {
        if response.Messages[0].Content != "this is a message from user1 to user2" {
            t.Errorf("Message returned had the wrong content: %v\n", response.Messages[0].Content)
        }
        if sender, _ := response.Messages[0].getSender(); sender.Uid != "1" {
            t.Errorf("Message returned had the wrong senderid: %v\n", response.Messages[0].SenderId)
        }
        if recipient, _ := response.Messages[0].getRecipientUser(); recipient.Uid != "2" {
            t.Errorf("Message returned had the wrong recipientid: %v\n", response.Messages[0].RecipientId)
        }
    }

    log.Println("Add a message from user2 to user3")
    user2.addMessageToUser(user3, "this is a message from user2 to user3", 1)

    log.Println("List the messages between user3 and user2")
    response = listMessagesEndpoint(user3, "2")

    if response.Error != "" {
        t.Error("Response had error when it shouldn't have")
    }

    if len(response.Messages) != 1 {
        t.Errorf("1 message expected, found %v\n", len(response.Messages))
    } else {
        if response.Messages[0].Content != "this is a message from user2 to user3" {
            t.Errorf("Message returned had the wrong content: %v\n", response.Messages[0].Content)
        }
        if sender, _ := response.Messages[0].getSender(); sender.Uid != "2" {
            t.Errorf("Message returned had the wrong senderid: %v\n", response.Messages[0].SenderId)
        }
        if recipient, _ := response.Messages[0].getRecipientUser(); recipient.Uid != "3" {
            t.Errorf("Message returned had the wrong recipientid: %v\n", response.Messages[0].RecipientId)
        }
    }

    log.Println("Add a message from user1 to user2")
    user1.addMessageToUser(user2, "this is another message from user1 to user2", 1)

    log.Println("List the messages between user2 and user1")
    response = listMessagesEndpoint(user2, "1")

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
    response = listMessagesEndpoint(user1, "2")

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
    response = listMessagesEndpoint(user1, "123")
    if response.Error != "Friend not found" {
        t.Errorf("Should have returned an error that friend couldn't be found")
    }
}

func TestSendMessageEndpoint(t *testing.T) {
    defer resetTables()
}