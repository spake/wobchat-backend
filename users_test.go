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

func TestCreatingUsers(t *testing.T) {
    defer resetTables()

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
    defer resetTables()

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
    defer resetTables()

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

    log.Println("Get the friends of test user 1 - should be empty")
    friends := testUser1.getFriends()
    if len(friends) != 0 {
        t.Error("Friends found for user with no friends")
    }

    log.Println("Get the friends of test user 2 - should be empty")
    friends = testUser2.getFriends()
    if len(friends) != 0 {
        t.Error("Friends found for user with no friends")
    }

    log.Println("Make test user 2 and test user 1 friends")
    testUser1.addFriend(testUser2)

    log.Println("Get friends of user 1")
    friends = testUser1.getFriends()
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser2 {
        t.Errorf("Friend not equal to test user 2")
    }

    log.Println("Get friends of user 2")
    friends = testUser2.getFriends()
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser1 {
        t.Errorf("Friend not equal to test user 1")
    }

    log.Println("Make test user 2 and test user 1 friends")
    testUser2.addFriend(testUser1)

    log.Println("Get friends of user 1")
    friends = testUser1.getFriends()
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser2 {
        t.Errorf("Friend not equal to test user 2")
    }

    log.Println("Get friends of user 2")
    friends = testUser2.getFriends()
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser1 {
        t.Errorf("Friend not equal to test user 1")
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
    friends = testUser3.getFriends()
    if len(friends) != 0 {
        t.Error("Friends found for user with no friends")
    }

    log.Println("Make test user 3 and test user 1 friends")
    testUser1.addFriend(testUser3)

    log.Println("Get friends of user 1")
    friends = testUser1.getFriends()
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
    friends = testUser2.getFriends()
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser1 {
        t.Errorf("Friend not equal to test user 1")
    }
    log.Println("Get friends of user 3")
    friends = testUser3.getFriends()
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser1 {
        t.Errorf("Friend not equal to test user 1")
    }

    log.Println("Make test user 2 and test user 3 friends")
    testUser3.addFriend(testUser2)

    log.Println("Get friends of user 1")
    friends = testUser1.getFriends()
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
    friends = testUser2.getFriends()
    if len(friends) != 2 {
        t.Errorf("2 friends should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser1 {
        t.Errorf("test user 1 not found in friends")
    }
    if friends[1] != testUser3 {
        t.Errorf("test user 3 not found in friends")
    }
    log.Println("Get friends of user 3")
    friends = testUser3.getFriends()
    if len(friends) != 2 {
        t.Errorf("2 friends should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser1 {
        t.Errorf("test user 1 not found in friends")
    }
    if friends[1] != testUser2 {
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
    defer resetTables()

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
    defer resetTables()

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
    defer resetTables()

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
    gotUserBefore := getUserFromInfo(testInfoBefore)

    if gotUserBefore != testUser {
        t.Errorf("User from before info is not the same")
    }

    testInfoAfter := testInfoBefore
    testInfoAfter.Picture = newPictureURL

    log.Println("Getting user from after info")
    gotUserAfter := getUserFromInfo(testInfoAfter)

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
    defer resetTables()

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
    gotUser := getUserFromInfo(testInfo)

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

func TestAddFriendEndpoint(t *testing.T) {
    defer resetTables()

    testUser4 := User{
        Uid:       "420",
        Name:      "Snoop Dogg",
        FirstName: "Snoop",
        LastName:  "Dogg",
        Email:     "blazeit@gmail.com",
        Picture:   "40keks"}

    log.Println("Creating test user 4")
    db.Create(&testUser4)

    testUser1 := User{
        Uid:       "12345",
        Name:      "Jayden Smith",
        FirstName: "Jayden",
        LastName:  "Smith",
        Email:     "poop@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 1")
    db.Create(&testUser1)

    log.Println("Get the friends of test user 4 - should be empty")
    friends := testUser4.getFriends()
    if len(friends) != 0 {
        t.Error("Friends found for user with no friends")
    }

    log.Println("Make test user 4 and test user 1 friends")
    response := addFriendEndpoint(testUser1, AddFriendRequest{Uid: "420"})

    if response.Success == false {
        t.Errorf("Adding friends didn't succeed when it should have. Error: %v\n", response.Error)
    }

    log.Println("Get the friends of test user 4")
    friends = testUser4.getFriends()
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser1 {
        t.Errorf("Friend not equal to test user 1")
    }

    log.Println("Get the friends of test user 1")
    friends = testUser1.getFriends()
    if len(friends) != 1 {
        t.Errorf("1 friends should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser4 {
        t.Errorf("test user 4 not found in friends")
    }

    log.Println("Make test user 4 and test user 1 friends again")
    response = addFriendEndpoint(testUser1, AddFriendRequest{Uid: "420"})

    if response.Success == true {
        t.Error("Adding friends succeeded when it shouldn't have.")
    }

    log.Println("Get the friends of test user 4")
    friends = testUser4.getFriends()
    if len(friends) != 1 {
        t.Errorf("1 friend should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser1 {
        t.Errorf("Friend not equal to test user 1")
    }

    log.Println("Get the friends of test user 1")
    friends = testUser1.getFriends()
    if len(friends) != 1 {
        t.Errorf("1 friends should have been found, found %v\n", len(friends))
    }
    if friends[0] != testUser4 {
        t.Errorf("test user 4 not found in friends")
    }
}

func TestAddAndGetMessagesWithUser(t *testing.T) {
    defer resetTables()

    user1 := User{
        Uid:       "420",
        Name:      "Snoop Doge",
        FirstName: "Snoop",
        LastName:  "Doge",
        Email:     "higher@gmail.com",
        Picture:   "42keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    user2 := User{
        Uid:       "421",
        Name:      "Peppa Pig",
        FirstName: "Peppa",
        LastName:  "Pig",
        Email:     "p.pig@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    log.Println("Getting messages (should be none)")
    messages1 := user1.getMessagesWithUser(user2)
    messages2 := user2.getMessagesWithUser(user1)
    if len(messages1) != 0 {
        t.Errorf("Should have found 0 messages, found %v\n", len(messages1))
    }
    if len(messages2) != 0 {
        t.Errorf("Should have found 0 messages, found %v\n", len(messages2))
    }

    log.Println("Adding invalid messages")
    if err := user1.addMessageToUser(user2, "this is messed up", -1); err == nil {
        t.Errorf("Should have failed with an invalid content type")
    }

    log.Println("Adding empty message")
    if err := user1.addMessageToUser(user2, "", ContentTypeText); err != nil {
        t.Errorf("Shouldn't have failed on empty message")
    }

    messages1 = user1.getMessagesWithUser(user2)
    messages2 = user2.getMessagesWithUser(user1)
    if len(messages1) != 1 {
        t.Errorf("Should have found 1 message, found %v\n", len(messages1))
    }
    if len(messages2) != 1 {
        t.Errorf("Should have found 1 message, found %v\n", len(messages2))
    }

    text1 := "u wot snoop?"
    text2 := "top kek"

    log.Println("Adding normal messages")
    if err := user2.addMessageToUser(user1, text1, ContentTypeText); err != nil {
        t.Errorf("Shouldn't have failed on normal message")
    }
    if err := user1.addMessageToUser(user2, text2, ContentTypeText); err != nil {
        t.Errorf("Shouldn't have failed on normal message")
    }

    messages1 = user1.getMessagesWithUser(user2)
    messages2 = user2.getMessagesWithUser(user1)
    if len(messages1) != 3 {
        t.Errorf("Should have found 3 messages, found %v\n", len(messages1))
    }
    if len(messages2) != 3 {
        t.Errorf("Should have found 3 messages, found %v\n", len(messages2))
    }

    if messages1[0].Content != "" {
        t.Errorf("Invalid message content; wanted %v, found %v\n", "", messages1[0].Content)
    }
    if messages1[1].Content != text1 {
        t.Errorf("Invalid message content; wanted %v, found %v\n", text1, messages1[1].Content)
    }
    if messages1[2].Content != text2 {
        t.Errorf("Invalid message content; wanted %v, found %v\n", text2, messages1[2].Content)
    }
    if messages2[0].Content != "" {
        t.Errorf("Invalid message content; wanted %v, found %v\n", "", messages2[0].Content)
    }
    if messages2[1].Content != text1 {
        t.Errorf("Invalid message content; wanted %v, found %v\n", text1, messages2[1].Content)
    }
    if messages2[2].Content != text2 {
        t.Errorf("Invalid message content; wanted %v, found %v\n", text2, messages2[2].Content)
    }

    if sender, err := messages1[0].getSender(); err != nil || sender.Id != messages1[0].SenderId {
        t.Errorf("Invalid sender ID")
    }
    if recipient, err := messages1[0].getRecipientUser(); err != nil || recipient.Id != messages1[0].RecipientId {
        t.Errorf("Invalid recipient ID")
    }
}

func TestListMessagesEndpoint(t *testing.T) {
    defer resetTables()

    user1 := User{
        Uid:       "420",
        Name:      "Snoop Doge",
        FirstName: "Snoop",
        LastName:  "Doge",
        Email:     "higher@gmail.com",
        Picture:   "42keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    user2 := User{
        Uid:       "421",
        Name:      "Peppa Pig",
        FirstName: "Peppa",
        LastName:  "Pig",
        Email:     "p.pig@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    log.Println("Getting lists of messages (expecting it to be empty)")
    resp1 := listMessagesEndpoint(user1, user2.Uid)
    resp2 := listMessagesEndpoint(user2, user1.Uid)
    if len(resp1.Messages) != 0 {
        t.Errorf("Wrong number of messages; expected 0, found %v\n", len(resp1.Messages))
    }
    if len(resp2.Messages) != 0 {
        t.Errorf("Wrong number of messages; expected 0, found %v\n", len(resp2.Messages))
    }

    log.Println("Adding message")
    if err := user1.addMessageToUser(user2, "hello", ContentTypeText); err != nil {
        t.Errorf("Failed to add message")
    }

    log.Println("Getting lists of messages (expecting it to have 1)")
    resp1 = listMessagesEndpoint(user1, user2.Uid)
    resp2 = listMessagesEndpoint(user2, user1.Uid)
    if len(resp1.Messages) != 1 {
        t.Errorf("Wrong number of messages; expected 1, found %v\n", len(resp1.Messages))
    }
    if len(resp2.Messages) != 1 {
        t.Errorf("Wrong number of messages; expected 1, found %v\n", len(resp2.Messages))
    }

    log.Println("Adding another message")
    if err := user2.addMessageToUser(user1, "yo", ContentTypeText); err != nil {
        t.Errorf("Failed to add message")
    }

    log.Println("Getting lists of messages (expecting it to have 2)")
    resp1 = listMessagesEndpoint(user1, user2.Uid)
    resp2 = listMessagesEndpoint(user2, user1.Uid)
    if len(resp1.Messages) != 2 {
        t.Errorf("Wrong number of messages; expected 2, found %v\n", len(resp1.Messages))
    }
    if len(resp2.Messages) != 2 {
        t.Errorf("Wrong number of messages; expected 2, found %v\n", len(resp2.Messages))
    }
}

