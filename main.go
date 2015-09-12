package main

import (
    "log"
    "net/http"

    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)

// Should be set by Bamboo during build
var bambooBuildNumber string

func main() {
    log.Println("Opening DB connection")
    db, _ := gorm.Open("postgres", "dbname=backend sslmode=disable")
    db.LogMode(false)

    // Create tables and automigrate
    log.Println("Creating/migrating tables")
    db.CreateTable(&User{})
    db.AutoMigrate(&User{})

    // Set up HTTP handlers
    log.Println("Starting HTTP server")
    http.HandleFunc("/about", aboutHandler)
    http.ListenAndServe("127.0.0.1:8000", nil)
}
