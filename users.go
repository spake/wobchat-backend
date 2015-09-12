package main

type User struct {
    ID      string  `gorm:"primary_key"`
    Name    string
    Email   string
    Picture string
}
