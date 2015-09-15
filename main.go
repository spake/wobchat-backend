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
    log.Println("Creating certificate cache")
    cache = newMemoryCache()

    log.Println("Opening DB connection")
    db, err := gorm.Open("postgres", "host=/var/run/postgresql dbname=backend sslmode=disable")
    if err != nil {
        log.Println("Failed to open DB connection")
        panic(err)
    }
    defer db.Close()

    // Create tables and automigrate
    log.Println("Creating/migrating tables")
    db.CreateTable(&User{})
    db.AutoMigrate(&User{})

    // Set up HTTP handlers
    log.Println("Starting HTTP server")
    http.HandleFunc("/about", aboutHandler)
    http.HandleFunc("/verify", verifyHandler)
    http.ListenAndServe("127.0.0.1:8000", nil)
}
