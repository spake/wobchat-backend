package main

import (
    "log"
    "net/http"

    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
)

// store DB in a global :(
var db gorm.DB

var cfg Config

func main() {
    log.Println("Creating certificate cache")
    cache = newMemoryCache()

    log.Println("Reading config file")
    cfg = setupConfig()

    log.Println("Opening DB connection")

    var err error
    db, err = gorm.Open(cfg.Database.Type, cfg.Database.ConnectionString)

    if err != nil {
        log.Println("Failed to open DB connection")
        panic(err)
    }
    defer db.Close()

    // Create tables and automigrate
    log.Println("Creating/migrating tables")
    db.AutoMigrate(&User{})
    db.AutoMigrate(&UserFriend{})
    db.AutoMigrate(&Message{})

    // Set up HTTP handlers
    log.Println("Starting HTTP server")
    router := setupAPIHandlers()
    http.ListenAndServe("127.0.0.1:8000", router)
}
