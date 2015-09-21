package main

import (
    "errors"
    "time"
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
