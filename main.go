package main

import (
    "log"
    "net/http"

    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)

// Should be set by Bamboo during build
var bambooBuildNumber string

// store DB in a global :(
var db gorm.DB

func main() {
    log.Println("Creating certificate cache")
    cache = newMemoryCache()

    log.Println("Opening DB connection")
    var err error

    db, err := gorm.Open("postgres", "host=/var/run/postgresql dbname=backend sslmode=disable")

    // test configuration: leave commented out in production
    //db, err = gorm.Open("postgres", "dbname=backend sslmode=disable")

    if err != nil {
        log.Println("Failed to open DB connection")
        panic(err)
    }
    defer db.Close()

    // Create tables and automigrate
    log.Println("Creating/migrating tables")
    db.AutoMigrate(&User{})
    db.AutoMigrate(&UserFriend{})

    // Set up HTTP handlers
    log.Println("Starting HTTP server")
    setupAPIHandlers()
    http.ListenAndServe("127.0.0.1:8000", nil)
}
