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

func TestGetUserFromInfo(t *testing.T) {
    oldPictureURL := "old_picture"
    newPictureURL := "new_picture"

    testUser := User{
        Uid:       "10000",
        Name:      "Wanye West",
        FirstName: "Wanye",
        LastName:  "West",
        Email:     "uhh_hello@gmail.com",
        Picture:   oldPictureURL,
    }

    log.Println("Creating test user")
    db.Create(&testUser)

    testInfoBefore := GoogleInfo{
        ID:             testUser.Uid,
        DisplayName:    testUser.Name,
        FirstName:      testUser.FirstName,
        LastName:       testUser.LastName,
        Email:          testUser.Email,
        Picture:        oldPictureURL,
    }

    log.Println("Getting user from before info")
    gotUserBefore := getUserFromInfo(db, testInfoBefore)

    if gotUserBefore != testUser {
        t.Errorf("User from before info is not the same")
    }

    testInfoAfter := testInfoBefore
    testInfoAfter.Picture = newPictureURL

    log.Println("Getting user from after info")
    gotUserAfter := getUserFromInfo(db, testInfoAfter)

    if gotUserAfter.Picture != newPictureURL {
        t.Errorf("User was not updated from GoogleInfo")
    }

    // check user is still in the db, and has been updated
    var testUserAfter User
    if err := db.Where("uid = ?", testUser.Uid).Find(&testUserAfter).Error; err != nil {
        t.Errorf("User is no longer in the database")
    }
    if testUserAfter.Picture != newPictureURL {
        t.Errorf("User was not updated from GoogleInfo in the database")
    }

    // check everything else was the same
    gotUserAfter.Picture = oldPictureURL
    testUserAfter.Picture = oldPictureURL
    if gotUserAfter != testUser {
        t.Errorf("Something wrong was changed after updating user from GoogleInfo")
    }
    if testUserAfter != testUser {
        t.Errorf("Something wrong was changed in the database after updating user from GoogleInfo")
    }
}

func TestGetUserFromInfoNew(t *testing.T) {
    testInfo := GoogleInfo{
        ID:             "10001",
        DisplayName:    "Wanye Test",
        FirstName:      "Wanye",
        LastName:       "Test",
        Email:          "uhh_sorry@gmail.com",
        Picture:        "something",
    }

    // ensure the user doesn't exist yet
    var testUser User
    if err := db.Where("uid = ?", testInfo.ID).Find(&testUser).Error; err == nil {
        t.Errorf("User already existed before signing in")
    }

    log.Println("Getting user from info")
    gotUser := getUserFromInfo(db, testInfo)

    // check all the fields are correct
    if testInfo.ID != gotUser.Uid {
        t.Errorf("Uid is not correct (%v -> %v)\n", testInfo.ID, gotUser.Uid)
    }
    if testInfo.DisplayName != gotUser.Name {
        t.Errorf("Name is not correct (%v -> %v)\n", testInfo.DisplayName, gotUser.Name)
    }
    if testInfo.FirstName != gotUser.FirstName {
        t.Errorf("FirstName is not correct (%v -> %v)\n", testInfo.FirstName, gotUser.FirstName)
    }
    if testInfo.LastName != gotUser.LastName {
        t.Errorf("LastName is not correct (%v -> %v)\n", testInfo.LastName, gotUser.LastName)
    }
    if testInfo.Email != gotUser.Email {
        t.Errorf("Email is not correct (%v -> %v)\n", testInfo.Email, gotUser.Email)
    }
    if testInfo.Picture != gotUser.Picture {
        t.Errorf("Picture is not correct (%v -> %v)\n", testInfo.Picture, gotUser.Picture)
    }

    // ensure the user exists now
    if err := db.Where("uid = ?", testInfo.ID).Find(&testUser).Error; err != nil {
        t.Errorf("User doesn't exist in the database after calling getUserFromInfo")
    }

    // ensure this user is the same as what getUserFromInfo returned
    if testUser != gotUser {
        t.Errorf("User in database is different to what getUserFromInfo returned")
    }
}
