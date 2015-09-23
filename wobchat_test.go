package main

import (
    "testing"
    "log"
    "github.com/jinzhu/gorm"
    "os"
    "flag" // TW
)

var printQueries = flag.Bool("printqueries", false, "Print all queries run through the database")

func TestMain(m *testing.M) {
    cfg = setupConfig()

    log.Println("Opening DB connection")
    dbTmp, err := gorm.Open(cfg.Database.Type, cfg.Database.TestConnectionString)
    if err != nil {
        log.Println("Failed to open DB connection")
        panic(err)
    }
    db = dbTmp
    defer db.Close()

    if *printQueries {
        db.LogMode(true)
    }

    // drop the tables in case the last test run didn't drop them
    db.DropTable(&User{})
    db.DropTable(&UserFriend{})
    db.DropTable(&Message{})

    log.Println("Creating/migrating tables")
    db.AutoMigrate(&User{})
    db.AutoMigrate(&UserFriend{})
    db.AutoMigrate(&Message{})

    result := m.Run()

    db.DropTable(&User{})
    db.DropTable(&UserFriend{})
    db.DropTable(&Message{})

    os.Exit(result)
}

func resetTables() {
    log.Println("Resetting tables")
    db.Exec("DELETE FROM users;")
    db.Exec("DELETE FROM user_friends;")
    db.Exec("DELETE FROM messages;")
}
