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

    db.LogMode(true)

    // Create tables and automigrate
    log.Println("Creating/migrating tables")
    db.AutoMigrate(&User{})
    db.AutoMigrate(&UserFriend{})

    testUser := User{
        Uid:       12345,
        Name:      "Jayden Smith",
        FirstName: "Jayden",
        LastName:  "Smith",
        Email:     "poop@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 1")
    db.Create(&testUser)

    testUser = User{
        Uid:       12346,
        Name:      "Will Smith",
        FirstName: "Will",
        LastName:  "Smith",
        Email:     "doody@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&testUser)

    log.Println("Accessing test user 1")
    user1 := User{}
    db.Where(&User{Uid: 12345}).First(&user1)
    log.Printf("User's name: %v\n", user1.Name)

    log.Println("Accessing test user 2")
    user2 := User{}
    db.Where(&User{Uid: 12346}).First(&user2)
    log.Printf("User's name: %v\n", user2.Name)

    log.Println("Adding test user 2 as friend to test user 1")
    user1.addFriend(db, user2)
    log.Println(user1)
    log.Println(user2)

    log.Println("Get the users from the DB again")
    db.Where(&User{Uid: 12345}).First(&user1)
    db.Where(&User{Uid: 12346}).First(&user2)
    log.Println(user1)
    log.Println(user2)

    log.Println("Adding test user 1 as friend to test user 2")
    user2.addFriend(db, user1)
    log.Println(user1)
    log.Println(user2)

    log.Println("Get the users from the DB again")
    db.Where(&User{Uid: 12345}).First(&user1)
    db.Where(&User{Uid: 12346}).First(&user2)
    log.Println(user1)
    log.Println(user2)

    log.Println("Get friends of user 1")
    friends := user1.getFriends(db)
    log.Println(friends)
    log.Println("Get friends of user 2")
    friends = user2.getFriends(db)
    log.Println(friends)

    testUser = User{
        Uid:       12347,
        Name:      "Kanye West",
        FirstName: "Kanye",
        LastName:  "West",
        Email:     "shit@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 3")
    db.Create(&testUser)

    log.Println("Get friends of user 1")
    friends = user1.getFriends(db)
    log.Println(friends)
    log.Println("Get friends of user 2")
    friends = user2.getFriends(db)
    log.Println(friends)
    log.Println("Get friends of user 3")
    user3 := User{}
    db.Where(&User{Uid: 12347}).First(&user3)
    friends = user3.getFriends(db)
    log.Println(friends)

    log.Println("Adding test user 3 as friend to test user 1")
    user1.addFriend(db, user3)

    log.Println("Get friends of user 1")
    friends = user1.getFriends(db)
    log.Println(friends)
    log.Println("Get friends of user 2")
    friends = user2.getFriends(db)
    log.Println(friends)
    log.Println("Get friends of user 3")
    friends = user3.getFriends(db)
    log.Println(friends)

    log.Println("Adding test user 2 as friend to test user 3")
    user3.addFriend(db, user2)

    log.Println("Get friends of user 1")
    friends = user1.getFriends(db)
    log.Println(friends)
    log.Println("Get friends of user 2")
    friends = user2.getFriends(db)
    log.Println(friends)
    log.Println("Get friends of user 3")
    friends = user3.getFriends(db)
    log.Println(friends)

    // Set up HTTP handlers
    log.Println("Starting HTTP server")
    http.HandleFunc("/about", aboutHandler)
    http.HandleFunc("/verify", verifyHandler)
    http.ListenAndServe("127.0.0.1:8000", nil)
}
