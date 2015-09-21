package main

const (
    ContentTypeText = 1
)

const (
    RecipientTypeUser = 1
)

type Message struct {
    Mid                 int         `gorm:"primary_key"`
    Content             string
    ContentType         int
    SenderUid           int         `sql:"not null"`
    RecipientId         int         `sql:"not null"`
    RecipientType       int
    Timestamp           time.Time
}
