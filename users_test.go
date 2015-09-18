package main

import (
    "testing"
    "log"
    "github.com/jinzhu/gorm"
    "os"
    "flag"
)

var printQueries = flag.Bool("printqueries", false, "Print all queries run through the database")

func TestMain(m *testing.M) {
    flag.Parse()
    log.Println("Opening DB connection")
    dbTmp, err := gorm.Open("postgres", "host=/var/run/postgresql dbname=backendtest sslmode=disable")
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

    log.Println("Creating/migrating tables")
    db.AutoMigrate(&User{})
    db.AutoMigrate(&UserFriend{})

    result := m.Run()

    db.DropTable(&User{})
    db.DropTable(&UserFriend{})

    os.Exit(result)
}

func TestCreatingUsers(t *testing.T) {
    testUser1 := User{
        Uid:       "12345",
        Name:      "Jayden Smith",
        FirstName: "Jayden",
        LastName:  "Smith",
        Email:     "poop@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 1")
    db.Create(&testUser1)

    testUser2 := User{
        Uid:       "12346",
        Name:      "Will Smith",
        FirstName: "Will",
        LastName:  "Smith",
        Email:     "doody@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&testUser2)

    log.Println("Accessing test user 1")
    user1 := User{}
    db.Where(&User{Uid: "12345"}).First(&user1)
    if user1 != testUser1 {
        t.Error("User accessed not the same as user inserted")
    }

    log.Println("Accessing test user 2")
    user2 := User{}
    db.Where(&User{Uid: "12346"}).First(&user2)
    if user2 != testUser2 {
        t.Error("User accessed not the same as user inserted")
    }
}

func TestDeletingUsers(t *testing.T) {
    log.Println("Deleting test user 1")
    db.Where(&User{Uid: "12345"}).Delete(User{})

    log.Println("Accessing test user 1")
    user1 := User{}
    db.Where(&User{Uid: "12345"}).First(&user1)
    if user1.Uid != "" {
        t.Error("Deleted user still exists")
    }

    log.Println("Deleting test user 2")
    db.Where(&User{Uid: "12345"}).Delete(User{})

    log.Println("Accessing test user 2")
    user2 := User{}
    db.Where(&User{Uid: "12345"}).First(&user1)
    if user2.Uid != "" {
        t.Error("Deleted user still exists")
    }
}

func TestAddingFriends(t *testing.T) {
    testUser1 := User{
        Uid:       "12345",
        Name:      "Jayden Smith",
        FirstName: "Jayden",
        LastName:  "Smith",
        Email:     "poop@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 1")
    db.Create(&testUser1)

    testUser2 := User{
        Uid:       "12346",
        Name:      "Will Smith",
        FirstName: "Will",
        LastName:  "Smith",
        Email:     "doody@gmail.com",
        Picture:   "someurl"}

    log.Println("Get the friends of test user 1 - should be empty")
    friends := testUser1.getFriends(db)
    if len(friends) != 0 {
        t.Error("Friends found for user with no friends")
    }

    log.Println("Get the friends of test user 2 - should be empty")
    friends = testUser2.getFriends(db)
    if len(friends) != 0 {
        t.Error("Friends found for user with no friends")
    }

    log.Println("Add test user 2 as a friend to test user 1")
    testUser1.addFriend(db, testUser2)

    log.Println("Get friends of user 1")
    friends = testUser1.getFriends(db)
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser2 {
        t.Errorf("Friend not equal to test user 2")
    }

    log.Println("Get friends of user 2")
    friends = testUser2.getFriends(db)
    if len(friends) != 0 {
        t.Error("Friends found for user with no friends")
    }

    log.Println("Add test user 2 as a friend to test user 1")
    testUser2.addFriend(db, testUser1)

    log.Println("Get friends of user 2")
    friends = testUser2.getFriends(db)
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser1 {
        t.Errorf("Friend not equal to test user 1")
    }

    log.Println("Get friends of user 1")
    friends = testUser1.getFriends(db)
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser2 {
        t.Errorf("Friend not equal to test user 2")
    }

    testUser3 := User{
        Uid:       "12347",
        Name:      "Kanye West",
        FirstName: "Kanye",
        LastName:  "West",
        Email:     "shit@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 3")
    db.Create(&testUser3)

    log.Println("Get friends of user 3")
    friends = testUser3.getFriends(db)
    if len(friends) != 0 {
        t.Error("Friends found for user with no friends")
    }

    log.Println("Adding test user 3 as friend to test user 1")
    testUser1.addFriend(db, testUser3)

    log.Println("Get friends of user 1")
    friends = testUser1.getFriends(db)
    if len(friends) != 2 {
        t.Errorf("2 friends should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser2 {
        t.Errorf("test user 2 not found in friends")
    }
    if friends[1] != testUser3 {
        t.Errorf("test user 3 not found in friends")
    }
    log.Println("Get friends of user 2")
    friends = testUser2.getFriends(db)
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser1 {
        t.Errorf("test user 1 not found in friends")
    }
    log.Println("Get friends of user 3")
    friends = testUser3.getFriends(db)
    if len(friends) != 0 {
        t.Error("Friends found for user with no friends")
    }

    log.Println("Adding test user 2 as friend to test user 3")
    testUser3.addFriend(db, testUser2)

    log.Println("Get friends of user 1")
    friends = testUser1.getFriends(db)
    if len(friends) != 2 {
        t.Errorf("2 friends should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser2 {
        t.Errorf("test user 2 not found in friends")
    }
    if friends[1] != testUser3 {
        t.Errorf("test user 3 not found in friends")
    }
    log.Println("Get friends of user 2")
    friends = testUser2.getFriends(db)
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser1 {
        t.Errorf("test user 1 not found in friends")
    }
    log.Println("Get friends of user 3")
    friends = testUser3.getFriends(db)
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser2 {
        t.Errorf("test user 2 not found in friends")
    }
}

func comparePublicUser(t *testing.T, testUser User, publicUser PublicUser) {
    log.Printf("Comparing user %v\n", testUser.Uid)

    if testUser.Uid != publicUser.Uid {
        t.Errorf("Uid not the same (%v -> %v)\n", testUser.Uid, publicUser.Uid)
    }
    if testUser.Name != publicUser.Name {
        t.Errorf("Name not the same (%v -> %v)\n", testUser.Name, publicUser.Name)
    }
    if testUser.FirstName != publicUser.FirstName {
        t.Errorf("FirstName not the same (%v -> %v)\n", testUser.FirstName, publicUser.FirstName)
    }
    if testUser.LastName != publicUser.LastName {
        t.Errorf("LastName not the same (%v -> %v)\n", testUser.LastName, publicUser.LastName)
    }
    if testUser.Picture != publicUser.Picture {
        t.Errorf("Picture not the same (%v -> %v)\n", testUser.Picture, publicUser.Picture)
    }
}


func TestPublicUser(t *testing.T) {
    testUser := User{
        Uid:       "1337",
        Name:      "John Smith",
        FirstName: "John",
        LastName:  "Smith",
        Email:     "poopmaster@gmail.com",
        Picture:   "something",
    }

    log.Println("Creating test user")
    db.Create(&testUser)

    log.Println("Converting to public")
    publicUser := testUser.toPublic()

    comparePublicUser(t, testUser, publicUser)
}

func TestPublicUsers(t *testing.T) {
    testUser1 := User{
        Uid:       "1338",
        Name:      "John Smith",
        FirstName: "John",
        LastName:  "Smith",
        Email:     "poopmaster@gmail.com",
        Picture:   "something1",
    }
    testUser2 := User{
        Uid:       "1339",
        Name:      "Jane Smith",
        FirstName: "Jane",
        LastName:  "Smith",
        Email:     "shitking@gmail.com",
        Picture:   "something2",
    }
    testUser3 := User{
        Uid:       "1340",
        Name:      "Jake Smith",
        FirstName: "Jake",
        LastName:  "Smith",
        Email:     "scrumlord@gmail.com",
        Picture:   "something3",
    }

    log.Println("Creating test users")
    db.Create(&testUser1)
    db.Create(&testUser2)
    db.Create(&testUser3)

    log.Println("Creating slice of test users")
    var testUsers Users
    testUsers = []User{testUser1, testUser2, testUser3}

    publicUsers := testUsers.toPublic()

    if len(testUsers) != len(publicUsers) {
        t.Errorf("Number of users not the same (%v -> %v)", len(testUsers), len(publicUsers))
    }

    for i := 0; i < len(testUsers); i++ {
        comparePublicUser(t, testUsers[i], publicUsers[i])
    }
}
