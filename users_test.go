package main

import (
    "fmt"
    "testing"
    "time"
    "log"
)

func TestCreatingUsers(t *testing.T) {
    defer resetTables()

    testUser1 := User{
        Id:        12345,
        Uid:       "12345",
        Name:      "Jayden Smith",
        FirstName: "Jayden",
        LastName:  "Smith",
        Email:     "poop@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 1")
    db.Create(&testUser1)

    testUser2 := User{
        Id:        12346,
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
    db.Where(&User{Id: 12345}).First(&user1)
    if user1 != testUser1 {
        t.Error("User accessed not the same as user inserted")
    }

    log.Println("Accessing test user 2")
    user2 := User{}
    db.Where(&User{Id: 12346}).First(&user2)
    if user2 != testUser2 {
        t.Error("User accessed not the same as user inserted")
    }
}

func TestDeletingUsers(t *testing.T) {
    defer resetTables()

    log.Println("Deleting test user 1")
    db.Where(&User{Id: 12345}).Delete(User{})

    log.Println("Accessing test user 1")
    user1 := User{}
    db.Where(&User{Id: 12345}).First(&user1)
    if user1.Uid != "" {
        t.Error("Deleted user still exists")
    }

    log.Println("Deleting test user 2")
    db.Where(&User{Id: 12345}).Delete(User{})

    log.Println("Accessing test user 2")
    user2 := User{}
    db.Where(&User{Id: 12345}).First(&user1)
    if user2.Uid != "" {
        t.Error("Deleted user still exists")
    }
}

func TestIsFriend(t *testing.T) {
    defer resetTables()

    testUser1 := User{
        Id:        12345,
        Uid:       "12345",
        Name:      "Jayden Smith",
        FirstName: "Jayden",
        LastName:  "Smith",
        Email:     "poop@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 1")
    db.Create(&testUser1)

    testUser2 := User{
        Id:        12346,
        Uid:       "12346",
        Name:      "Will Smith",
        FirstName: "Will",
        LastName:  "Smith",
        Email:     "doody@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&testUser2)

    log.Println("Check they are not friends")
    if testUser1.isFriend(testUser2) {
        t.Error("User 1 is friends with user 2")
    }
    if testUser2.isFriend(testUser1) {
        t.Error("User 2 is friends with user 1")
    }

    log.Println("Adding them as friends")
    testUser1.addFriend(testUser2)

    log.Println("Checking they are friends")
    if !testUser1.isFriend(testUser2) {
        t.Error("User 1 is not friends with user 2")
    }
    if !testUser2.isFriend(testUser1) {
        t.Error("User 2 is not friends with user 1")
    }
}

func TestAddingFriends(t *testing.T) {
    defer resetTables()

    testUser1 := User{
        Id:        12345,
        Uid:       "12345",
        Name:      "Jayden Smith",
        FirstName: "Jayden",
        LastName:  "Smith",
        Email:     "poop@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 1")
    db.Create(&testUser1)

    testUser2 := User{
        Id:        12346,
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
        Id:        12347,
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
    if friends[1] != testUser2 {
        t.Errorf("test user 2 not found in friends")
    }
    if friends[0] != testUser3 {
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
    if friends[1] != testUser2 {
        t.Errorf("test user 2 not found in friends")
    }
    if friends[0] != testUser3 {
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

func TestDeletingFriends(t *testing.T) {
    defer resetTables()

    user1 := User{
        Id:        12345,
        Uid:       "12345",
        Name:      "Jayden Smith",
        FirstName: "Jayden",
        LastName:  "Smith",
        Email:     "poop@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    user2 := User{
        Id:        12346,
        Uid:       "12346",
        Name:      "Will Smith",
        FirstName: "Will",
        LastName:  "Smith",
        Email:     "doody@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    user3 := User{
        Id:        12347,
        Uid:       "12347",
        Name:      "Hayden Smith",
        FirstName: "Hayden",
        LastName:  "Smith",
        Email:     "dootdoot@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 3")
    db.Create(&user3)

    log.Println("Trying to delete non-existent friendship")
    if err := user1.deleteFriend(user2); err == nil {
        t.Error("Succeeded in deleting non-existent friendship")
    }
    if err := user2.deleteFriend(user1); err == nil {
        t.Error("Succeeded in deleting non-existent friendship")
    }

    log.Println("Adding friendship")
    user1.addFriend(user2)

    log.Println("Trying to delete friendship (1)")
    if err := user1.deleteFriend(user2); err != nil {
        t.Errorf("Failed to delete friendship: %v", err)
    }
    if user1.isFriend(user2) {
        t.Error("Friendship wasn't actually deleted")
    }
   
    log.Println("Trying to delete friendship that's already deleted")
    if err := user1.deleteFriend(user2); err == nil {
        t.Error("Succeeded in deleting non-existent friendship")
    }
    if err := user2.deleteFriend(user1); err == nil {
        t.Error("Succeeded in deleting non-existent friendship")
    }

    log.Println("Adding friendship")
    user1.addFriend(user2)

    log.Println("Trying to delete friendship (2)")
    if err := user2.deleteFriend(user1); err != nil {
        t.Errorf("Failed to delete friendship: %v", err)
    }
    if user2.isFriend(user1) {
        t.Error("Friendship wasn't actually deleted")
    }

    log.Println("Trying to delete friendship that's already deleted")
    if err := user1.deleteFriend(user2); err == nil {
        t.Error("Succeeded in deleting non-existent friendship")
    }
    if err := user2.deleteFriend(user1); err == nil {
        t.Error("Succeeded in deleting non-existent friendship")
    }

    log.Println("Trying to delete self as friend")
    if err := user1.deleteFriend(user1); err == nil {
        t.Error("Succeeded in deleting self as friend")
    }

    log.Println("Ensuring deletion doesn't affect other friendships")
    user1.addFriend(user2)
    user1.addFriend(user3)
    if err := user1.deleteFriend(user2); err != nil {
        t.Errorf("Failed to delete friendship: %v", err)
    }
    if user1.isFriend(user2) {
        t.Error("Failed to delete friendship")
    }
    if !user1.isFriend(user3) {
        t.Error("Deletion affected the wrong friendship")
    }
}

func comparePublicUser(t *testing.T, testUser User, publicUser PublicUser) {
    log.Printf("Comparing user %v\n", testUser.Id)

    if testUser.Id != publicUser.Id {
        t.Errorf("Id not the same (%v -> %v)\n", testUser.Id, publicUser.Id)
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
        Id:        1337,
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
        Id:        1338,
        Uid:       "1338",
        Name:      "John Smith",
        FirstName: "John",
        LastName:  "Smith",
        Email:     "poopmaster@gmail.com",
        Picture:   "something1",
    }
    testUser2 := User{
        Id:        1339,
        Uid:       "1339",
        Name:      "Jane Smith",
        FirstName: "Jane",
        LastName:  "Smith",
        Email:     "shitking@gmail.com",
        Picture:   "something2",
    }
    testUser3 := User{
        Id:        1340,
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
        Id:        10000,
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
    if err := db.Where("id = ?", testUser.Id).Find(&testUserAfter).Error; err != nil {
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
        Id:        420,
        Uid:       "420",
        Name:      "Snoop Dogg",
        FirstName: "Snoop",
        LastName:  "Dogg",
        Email:     "blazeit@gmail.com",
        Picture:   "40keks"}

    log.Println("Creating test user 4")
    db.Create(&testUser4)

    testUser1 := User{
        Id:        12345,
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
    response := addFriendEndpoint(testUser1, AddFriendRequest{Id: 420})

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
    response = addFriendEndpoint(testUser1, AddFriendRequest{Id: 420})

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

func TestDeleteFriendEndpoint(t *testing.T) {
    defer resetTables()

    user1 := User{
        Id:        420,
        Uid:       "420",
        Name:      "Snoop Dogg",
        FirstName: "Snoop",
        LastName:  "Dogg",
        Email:     "blazeit@gmail.com",
        Picture:   "40keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    user2 := User{
        Id:        12345,
        Uid:       "12345",
        Name:      "Jayden Smith",
        FirstName: "Jayden",
        LastName:  "Smith",
        Email:     "poop@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    log.Println("Trying to delete friends (not friends yet)")
    resp := deleteFriendEndpoint(user1, user2.Id)
    if resp.Success || resp.Error == "" {
        t.Error("Succeeded in deleting friendship that didn't exist")
    }

    log.Println("Adding users as friends")
    user1.addFriend(user2)

    log.Println("Trying to delete friends (users are friends)")
    resp = deleteFriendEndpoint(user1, user2.Id)
    if !resp.Success || resp.Error != "" {
        t.Errorf("Failed to delete friendship that existed (%v)", resp.Error)
    }

    log.Println("Trying to delete friends again (users no longer friends)")
    resp = deleteFriendEndpoint(user1, user2.Id)
    if resp.Success || resp.Error == "" {
        t.Error("Succeeded in deleting friendship that didn't exist")
    }
}

func TestGetFriendEndpoint(t *testing.T) {
    defer resetTables()

    testUser1 := User{
        Id:        1,
        Uid:       "1",
        Name:      "Snoop Dogg",
        FirstName: "Snoop",
        LastName:  "Dogg",
        Email:     "blazeit@gmail.com",
        Picture:   "40keks"}

    log.Println("Creating test user 1")
    db.Create(&testUser1)

    testUser2 := User{
        Id:        2,
        Uid:       "2",
        Name:      "Jayden Smith",
        FirstName: "Jayden",
        LastName:  "Smith",
        Email:     "poop@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&testUser2)

    log.Println("User 1 get friend user 1 (same user)")
    resp := getFriendEndpoint(testUser1, 1)
    if resp.Error == "" {
        t.Error("Users shouldn't be able to do get friend on their own ID")
    }

    log.Println("User 1 get friend user 2 (not friends)")
    resp = getFriendEndpoint(testUser1, 2)
    if resp.Error == "" {
        t.Error("Users shouldn't be able to do get friend on users that aren't their friends")
    }

    log.Println("Adding users as friends")
    testUser1.addFriend(testUser2)

    log.Println("User 1 get friend user 2")
    resp = getFriendEndpoint(testUser1, 2)
    if resp.Error != "" {
        t.Error("Users should be able to do get friend on users that are their friends")
    }
    if resp.Friend.Id != 2 {
        t.Error("Get friend returned bad ID")
    }

    log.Println("User 2 get friend user 1")
    resp = getFriendEndpoint(testUser2, 1)
    if resp.Error != "" {
        t.Error("Users should be able to do get friend on users that are their friends")
    }
    if resp.Friend.Id != 1 {
        t.Error("Get friend returned bad ID")
    }
}

func TestAddAndGetMessagesWithUser(t *testing.T) {
    defer resetTables()

    user1 := User{
        Id:        420,
        Uid:       "420",
        Name:      "Snoop Doge",
        FirstName: "Snoop",
        LastName:  "Doge",
        Email:     "higher@gmail.com",
        Picture:   "42keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    user2 := User{
        Id:        421,
        Uid:       "421",
        Name:      "Peppa Pig",
        FirstName: "Peppa",
        LastName:  "Pig",
        Email:     "p.pig@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    log.Println("Getting messages (should be none)")
    messages1 := user1.getMessagesWithUser(user2, -1, 100)
    messages2 := user2.getMessagesWithUser(user1, -1, 100)
    if len(messages1) != 0 {
        t.Errorf("Should have found 0 messages, found %v\n", len(messages1))
    }
    if len(messages2) != 0 {
        t.Errorf("Should have found 0 messages, found %v\n", len(messages2))
    }

    log.Println("Adding invalid messages")
    if _, err := user1.addMessageToUser(user2, "this is messed up", -1); err == nil {
        t.Errorf("Should have failed with an invalid content type")
    }

    log.Println("Adding empty message")
    if _, err := user1.addMessageToUser(user2, "", ContentTypeText); err != nil {
        t.Errorf("Shouldn't have failed on empty message")
    }

    messages1 = user1.getMessagesWithUser(user2, -1, 100)
    messages2 = user2.getMessagesWithUser(user1, -1, 100)
    if len(messages1) != 1 {
        t.Errorf("Should have found 1 message, found %v\n", len(messages1))
    }
    if len(messages2) != 1 {
        t.Errorf("Should have found 1 message, found %v\n", len(messages2))
    }

    text1 := "u wot snoop?"
    text2 := "top kek"

    log.Println("Adding normal messages")
    if _, err := user2.addMessageToUser(user1, text1, ContentTypeText); err != nil {
        t.Errorf("Shouldn't have failed on normal message")
    }
    if _, err := user1.addMessageToUser(user2, text2, ContentTypeText); err != nil {
        t.Errorf("Shouldn't have failed on normal message")
    }

    messages1 = user1.getMessagesWithUser(user2, -1, 100)
    messages2 = user2.getMessagesWithUser(user1, -1, 100)
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

func TestGetMessagesWithUserWithLastAndAmount(t *testing.T) {
    defer resetTables()

    user1 := User{
        Id:        420,
        Uid:       "420",
        Name:      "Snoop Doge",
        FirstName: "Snoop",
        LastName:  "Doge",
        Email:     "higher@gmail.com",
        Picture:   "42keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    user2 := User{
        Id:        421,
        Uid:       "421",
        Name:      "Peppa Pig",
        FirstName: "Peppa",
        LastName:  "Pig",
        Email:     "p.pig@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    user3 := User{
        Id:        422,
        Uid:       "422",
        Name:      "Brad Heap",
        FirstName: "Brad",
        LastName:  "Heap",
        Email:     "kiwisarecool@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 3")
    db.Create(&user3)

    log.Println("Adding 100 messages from user1 to user2")
    for i := 0; i < 100; i++ {
        user2.addMessageToUser(user1, fmt.Sprintf("Hello user2 from user1 %v", i), ContentTypeText)
    }

    log.Println("Adding 5 messages from user1 to user3")
    for i := 0; i < 5; i++ {
        user3.addMessageToUser(user1, fmt.Sprintf("Hello user3 from user1 %v", i), ContentTypeText)
    }

    log.Println("Adding 100 messages from user2 to user1")
    for i := 0; i < 100; i++ {
        user1.addMessageToUser(user2, fmt.Sprintf("Hello user1 from user2 %v", i), ContentTypeText)
    }

    log.Println("Getting last 100 messages between user2 and user1")
    messages := user1.getMessagesWithUser(user2, -1, 100)
    if len(messages) != 100 {
        t.Errorf("Should have returned only 100 messages, returned %v\n", len(messages))
    } else {
        for i := 0; i < 100; i++ {
            if messages[i].Content != fmt.Sprintf("Hello user1 from user2 %v", i) {
                t.Errorf("%vth message wrong", i)
            }
        }
    }

    log.Println("Getting last 100 messages between user1 and user2")
    messages = user2.getMessagesWithUser(user1, -1, 100)
    if len(messages) != 100 {
        t.Errorf("Should have returned only 100 messages, returned %v\n", len(messages))
    } else {
        for i := 0; i < 100; i++ {
            if messages[i].Content != fmt.Sprintf("Hello user1 from user2 %v", i) {
                t.Errorf("%vth message wrong", i)
            }
        }
    }

    last_message_id := messages[0].Id

    log.Println("Getting first 100 messages between user1 and user2")
    messages = user2.getMessagesWithUser(user1, last_message_id, 100)
    if len(messages) != 100 {
        t.Errorf("Should have returned only 100 messages, returned %v\n", len(messages))
    } else {
        for i := 0; i < 100; i++ {
            if messages[i].Content != fmt.Sprintf("Hello user2 from user1 %v", i) {
                t.Errorf("%vth message wrong", i)
            }
        }
    }

    log.Println("Getting first 10 messages between user1 and user2")
    messages = user2.getMessagesWithUser(user1, -1, 10)
    if len(messages) != 10 {
        t.Errorf("Should have returned only 10 messages, returned %v\n", len(messages))
    } else {
        for i := 0; i < 10; i++ {
            if messages[i].Content != fmt.Sprintf("Hello user1 from user2 %v", i+90) {
                t.Errorf("%vth message wrong", i)
            }
        }
    }

    last_message_id = messages[0].Id

    log.Println("Getting next 10 messages between user1 and user2")
    messages = user2.getMessagesWithUser(user1, last_message_id, 10)
    if len(messages) != 10 {
        t.Errorf("Should have returned only 10 messages, returned %v\n", len(messages))
    } else {
        for i := 0; i < 10; i++ {
            if messages[i].Content != fmt.Sprintf("Hello user1 from user2 %v", i+80) {
                t.Errorf("%vth message wrong", i)
            }
        }
    }

    log.Println("Getting first 10 messages between user1 and user3 but only 5 exist")
    messages = user3.getMessagesWithUser(user1, -1, 10)
    if len(messages) != 5 {
        t.Errorf("Should have returned only 5 messages, returned %v\n", len(messages))
    } else {
        for i := 0; i < 5; i++ {
            if messages[i].Content != fmt.Sprintf("Hello user3 from user1 %v", i) {
                t.Errorf("%vth message wrong", i)
            }
        }
    }

    third_message_id := messages[2].Id

    log.Println("Getting 10 messages before the third message between user1 and user3")
    messages = user3.getMessagesWithUser(user1, third_message_id, 10)
    if len(messages) != 2 {
        t.Errorf("Should have returned only 2 messages, returned %v\n", len(messages))
    } else {
        for i := 0; i < 2; i++ {
            if messages[i].Content != fmt.Sprintf("Hello user3 from user1 %v", i) {
                t.Errorf("%vth message wrong", i)
            }
        }
    }

    first_message_id := messages[0].Id
    log.Println("Getting 10 messages before the first message between user1 and user3")
    messages = user3.getMessagesWithUser(user1, first_message_id, 10)
    if len(messages) != 0 {
        t.Errorf("Should have returned no messages, returned %v\n", len(messages))
    }

    log.Println("Getting 1 message before the first message between user1 and user3")
    messages = user3.getMessagesWithUser(user1, first_message_id, 1)
    if len(messages) != 0 {
        t.Errorf("Should have returned no messages, returned %v\n", len(messages))
    }
}

func TestSearchUsernames(t *testing.T) {
    defer resetTables()

    log.Println("Listing all users - there should be none")
    users := searchUsernames("", 101)
    if len(users) != 0 {
        t.Errorf("0 users should have been found, found %v\n", len(users))
    }

    user1 := User{
        Id:        420,
        Uid:       "420",
        Name:      "Snoop Doge",
        FirstName: "Snoop",
        LastName:  "Doge",
        Email:     "higher@gmail.com",
        Picture:   "42keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    log.Println("Listing all users as id 421")
    users = searchUsernames("", 421)
    if len(users) != 1 {
        t.Errorf("1 user should have been found, found %v\n", len(users))
    } else {   
        if users[0].Id != 420 {
            t.Errorf("User had wrong Id. Expected 420, found %v\n", users[0].Id)
        }
    }

    log.Println("Listing all users as id 420")
    users = searchUsernames("", 420)
    if len(users) != 0 {
        t.Errorf("0 users should have been found, found %v\n", len(users))
    }

    user2 := User{
        Id:        421,
        Uid:       "421",
        Name:      "Peppa Pig",
        FirstName: "Peppa",
        LastName:  "Pig",
        Email:     "p.pig@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    log.Println("Listing all users as id 422")
    users = searchUsernames("", 422)
    if len(users) != 2 {
        t.Errorf("2 users should have been found, found %v\n", len(users))
    } else {
        if users[0].Id != 420 {
            t.Errorf("User had wrong Id. Expected 420, found %v\n", users[0].Id)
        }
        if users[1].Id != 421 {
            t.Errorf("User had wrong Id. Expected 421, found %v\n", users[1].Id)
        }
    }

    log.Println("Listing all users as id 420")
    users = searchUsernames("", 420)
    if len(users) != 1 {
        t.Errorf("1 user should have been found, found %v\n", len(users))
    } else {
        if users[0].Id != 421 {
            t.Errorf("User had wrong Id. Expected 421, found %v\n", users[0].Id)
        }
    }

    log.Println("Listing users matching 'Peppa Pig'")
    users = searchUsernames("Peppa Pig", 101)
    if len(users) != 1 {
        t.Errorf("1 user should have been found, found %v\n", len(users))
    } else {
        if users[0].Id != 421 {
            t.Errorf("User had wrong Id. Expected 421, found %v\n", users[0].Id)
        }
    }

    log.Println("Listing users matching 'Snoop'")
    users = searchUsernames("Snoop", 101)
    if len(users) != 1 {
        t.Errorf("1 user should have been found, found %v\n", len(users))
    } else {
        if users[0].Id != 420 {
            t.Errorf("User had wrong Id. Expected 420, found %v\n", users[0].Id)
        }
    }

    log.Println("Listing users matching 'p' as id 101")
    users = searchUsernames("p", 101)
    if len(users) != 2 {
        t.Errorf("2 users should have been found, found %v\n", len(users))
    } else {
        if users[0].Id != 420 {
            t.Errorf("User had wrong Id. Expected 420, found %v\n", users[0].Id)
        }
        if users[1].Id != 421 {
            t.Errorf("User had wrong Id. Expected 421, found %v\n", users[1].Id)
        }
    }

    log.Println("Listing users matching 'p' as id 420")
    users = searchUsernames("p", 420)
    if len(users) != 1{
        t.Errorf("1 user should have been found, found %v\n", len(users))
    } else {
        if users[0].Id != 421 {
            t.Errorf("User had wrong Id. Expected 421, found %v\n", users[0].Id)
        }
    }

    log.Println("Listing users matching 'sno'")
    users = searchUsernames("sno", 101)
    if len(users) != 1 {
        t.Errorf("1 user should have been found, found %v\n", len(users))
    } else {
        if users[0].Id != 420 {
            t.Errorf("User had wrong Id. Expected 420, found %v\n", users[0].Id)
        }
    }

    log.Println("Listing users matching 'ig'")
    users = searchUsernames("ig", 101)
    if len(users) != 1 {
        t.Errorf("1 user should have been found, found %v\n", len(users))
    } else {
        if users[0].Id != 421 {
            t.Errorf("User had wrong Id. Expected 421, found %v\n", users[0].Id)
        }
    }

    log.Println("Listing users matching 'Snooo'")
    users = searchUsernames("Snooo", 101)
    if len(users) != 0 {
        t.Errorf("0 users should have been found, found %v\n", len(users))
    }

    log.Println("Listing users matching 'higher@gmail.com'")
    users = searchUsernames("higher@gmail.com", 101)
    if len(users) != 1 {
        t.Errorf("1 user should have been found, found %v\n", len(users))
    } else {
        if users[0].Id != 420 {
            t.Errorf("User had wrong Id. Expected 420, found %v\n", users[0].Id)
        }
    }

    log.Println("Listing users matching 'hIGHer@gmAIl.com'")
    users = searchUsernames("hIGHer@gmAIl.com", 101)
    if len(users) != 1 {
        t.Errorf("1 user should have been found, found %v\n", len(users))
    } else {
        if users[0].Id != 420 {
            t.Errorf("User had wrong Id. Expected 420, found %v\n", users[0].Id)
        }
    }

    log.Println("Listing users matching 'p.pig@gmail.com'")
    users = searchUsernames("p.pig@gmail.com", 101)
    if len(users) != 1 {
        t.Errorf("1 user should have been found, found %v\n", len(users))
    } else {
        if users[0].Id != 421 {
            t.Errorf("User had wrong Id. Expected 421, found %v\n", users[0].Id)
        }
    }

    log.Println("Listing users matching 'p.pig@gmail.com' as user 420")
    users = searchUsernames("p.pig@gmail.com", 420)
    if len(users) != 1 {
        t.Errorf("1 user should have been found, found %v\n", len(users))
    } else {
        if users[0].Id != 421 {
            t.Errorf("User had wrong Id. Expected 421, found %v\n", users[0].Id)
        }
    }

    log.Println("Listing users matching 'p.pig@gmail.com' as user 421")
    users = searchUsernames("p.pig@gmail.com", 421)
    if len(users) != 0 {
        t.Errorf("0 users should have been found, found %v\n", len(users))
    }

    log.Println("Listing users matching 'p.pig@gmail.co'")
    users = searchUsernames("p.pig@gmail.co", 101)
    if len(users) != 0 {
        t.Errorf("0 users should have been found, found %v\n", len(users))
    }

    log.Println("Listing users matching 'p.pig@gmail.com.'")
    users = searchUsernames("p.pig@gmail.com.", 101)
    if len(users) != 0 {
        t.Errorf("0 users should have been found, found %v\n", len(users))
    }

    log.Println("Listing users matching 'a@b.c")
    users = searchUsernames("a@b.c", 101)
    if len(users) != 0 {
        t.Errorf("0 users should have been found, found %v\n", len(users))
    }
}

func TestListUsersEndpoint(t *testing.T) {
    defer resetTables()

    user1 := User{
        Id:        420,
        Uid:       "420",
        Name:      "Snoop Doge",
        FirstName: "Snoop",
        LastName:  "Doge",
        Email:     "higher@gmail.com",
        Picture:   "42keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    user2 := User{
        Id:        421,
        Uid:       "421",
        Name:      "Peppa Pig",
        FirstName: "Peppa",
        LastName:  "Pig",
        Email:     "p.pig@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    log.Println("Listing users matching 'DOG'")
    resp := listUsersEndpoint("DOG", 421)

    if !resp.Success {
        t.Error("Listing users didn't succeed when it should have.")
    }

    if len(resp.Users) != 1 {
        t.Errorf("1 user should have been found, found %v\n", len(resp.Users))
    } else {
        if resp.Users[0].Id != 420 {
            t.Errorf("User had wrong Id. Expected 420, found %v\n", resp.Users[0].Id)
        }
    }
}

func TestGetMeEndpoint(t *testing.T) {
    defer resetTables()

    user1 := User{
        Id:        420,
        Uid:       "420",
        Name:      "Snoop Doge",
        FirstName: "Snoop",
        LastName:  "Doge",
        Email:     "higher@gmail.com",
        Picture:   "42keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)
    
    resp := getMeEndpoint(user1)

    // this should always work...
    if !resp.Success || resp.Error != "" {
        t.Errorf("Get me endpoint shouldn't fail")
    }
    if resp.User != user1.toPublic() {
        t.Errorf("Get me endpoint's public user didn't match")
    }
}

func TestFriendRequests(t *testing.T) {
    defer resetTables()
    user1 := User{
        Id:        420,
        Uid:       "420",
        Name:      "Snoop Doge",
        FirstName: "Snoop",
        LastName:  "Doge",
        Email:     "higher@gmail.com",
        Picture:   "42keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    log.Println("Requesting a friend request to yourself")
    err := user1.addFriendRequest(user1)
    if err.Error() != "Cannot request to be your own friend" {
        t.Error("Friend requests to yourself should fail")
    }

    user2 := User{
        Id:        421,
        Uid:       "421",
        Name:      "Peppa Pig",
        FirstName: "Peppa",
        LastName:  "Pig",
        Email:     "p.pig@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    log.Println("Getting user 1's friend requests - should be empty")
    friendrequests := user1.getFriendRequests()
    if len(friendrequests) != 0 {
        t.Error("Friend requests found for user with no friend requests")
    }

    log.Println("Getting user 2's friend requests - should be empty")
    friendrequests = user2.getFriendRequests()
    if len(friendrequests) != 0 {
        t.Error("Friend requests found for user with no friend requests")
    }

    log.Println("Testing if user 1 has a friend request from user 2")
    if user1.hasFriendRequest(user2) {
        t.Error("User 1 shouldn't have a friend request from user 2")
    }

    log.Println("Testing if user 2 has a friend request from user 1")
    if user2.hasFriendRequest(user1) {
        t.Error("User 2 shouldn't have a friend request from user 1")
    }

    log.Println("Add a friend request from user2 to user1")
    user1.addFriendRequest(user2)

    log.Println("Getting user 1's friend requests")
    friendrequests = user1.getFriendRequests()
    if len(friendrequests) != 1 {
        t.Errorf("1 friend request should have been found, found %v\n", len(friendrequests))
    } else {
        if friendrequests[0].Id != 421 {
            t.Errorf("Friend request had wrong user Id. Expected 421, found %v\n", friendrequests[0].Id)
        }
    }

    log.Println("Getting user 2's friend requests - should be empty")
    friendrequests = user2.getFriendRequests()
    if len(friendrequests) != 0 {
        t.Error("Friend requests found for user with no friend requests")
    }

    log.Println("Testing if user 1 has a friend request from user 2")
    if !user1.hasFriendRequest(user2) {
        t.Error("User 1 should have a friend request from user 2")
    }

    log.Println("Testing if user 2 has a friend request from user 1")
    if user2.hasFriendRequest(user1) {
        t.Error("User 2 shouldn't have a friend request from user 1")
    }

    log.Println("Add a friend request from user1 to user2")
    user2.addFriendRequest(user1)

    log.Println("Getting user 1's friend requests")
    friendrequests = user1.getFriendRequests()
    if len(friendrequests) != 1 {
        t.Errorf("1 friend request should have been found, found %v\n", len(friendrequests))
    } else {
        if friendrequests[0].Id != 421 {
            t.Errorf("Friend request had wrong user Id. Expected 421, found %v\n", friendrequests[0].Id)
        }
    }

    log.Println("Getting user 2's friend requests")
    friendrequests = user2.getFriendRequests()
    if len(friendrequests) != 1 {
        t.Errorf("1 friend request should have been found, found %v\n", len(friendrequests))
    } else {
        if friendrequests[0].Id != 420 {
            t.Errorf("Friend request had wrong user Id. Expected 420, found %v\n", friendrequests[0].Id)
        }
    }

    log.Println("Testing if user 1 has a friend request from user 2")
    if !user1.hasFriendRequest(user2) {
        t.Error("User 1 should have a friend request from user 2")
    }

    log.Println("Testing if user 2 has a friend request from user 1")
    if !user2.hasFriendRequest(user1) {
        t.Error("User 2 should have a friend request from user 1")
    }

    user3 := User{
        Id:        422,
        Uid:       "422",
        Name:      "Yo Mum",
        FirstName: "Yo",
        LastName:  "Mum",
        Email:     "top.kek@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 3")
    db.Create(&user3)

    log.Println("Testing if user 1 has a friend request from user 3")
    if user1.hasFriendRequest(user3) {
        t.Error("User 1 shouldn't have a friend request from user 3")
    }

    log.Println("Add a friend request from user3 to user1")
    user1.addFriendRequest(user3)

    log.Println("Getting user 1's friend requests")
    friendrequests = user1.getFriendRequests()
    if len(friendrequests) != 2 {
        t.Errorf("2 friend requests should have been found, found %v\n", len(friendrequests))
    } else {
        if friendrequests[0].Id != 421 {
            t.Errorf("Friend request had wrong user Id. Expected 421, found %v\n", friendrequests[0].Id)
        }
        if friendrequests[1].Id != 422 {
            t.Errorf("Friend request had wrong user Id. Expected 422, found %v\n", friendrequests[1].Id)
        }
    }

    log.Println("Testing if user 1 has a friend request from user 3")
    if !user1.hasFriendRequest(user3) {
        t.Error("User 1 should have a friend request from user 3")
    }

    log.Println("Testing if user 3 has a friend request from user 1")
    if user3.hasFriendRequest(user1) {
        t.Error("User 3 shouldn't have a friend request from user 1")
    }
}

func TestListMyFriendRequestsEndpoint(t *testing.T) {
    defer resetTables()
    user1 := User{
        Id:        420,
        Uid:       "420",
        Name:      "Snoop Doge",
        FirstName: "Snoop",
        LastName:  "Doge",
        Email:     "higher@gmail.com",
        Picture:   "42keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    user2 := User{
        Id:        421,
        Uid:       "421",
        Name:      "Peppa Pig",
        FirstName: "Peppa",
        LastName:  "Pig",
        Email:     "p.pig@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    user3 := User{
        Id:        422,
        Uid:       "422",
        Name:      "Yo Mum",
        FirstName: "Yo",
        LastName:  "Mum",
        Email:     "top.kek@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 3")
    db.Create(&user3)

    log.Println("Add a friend request from user2 to user1")
    user1.addFriendRequest(user2)

    log.Println("Add a friend request from user3 to user1")
    user1.addFriendRequest(user3)

    log.Println("Listing friend requests")
    resp := listMyFriendRequestsEndpoint(user1)

    if !resp.Success {
        t.Error("Listing failed when it should have succeeded")
    }

    if len(resp.Requestors) != 2 {
        t.Errorf("2 friend requests should have been found, found %v\n", len(resp.Requestors))
    } else {
        if resp.Requestors[0].Id != 421 {
            t.Errorf("Friend request had wrong user Id. Expected 421, found %v\n", resp.Requestors[0].Id)
        }
        if resp.Requestors[1].Id != 422 {
            t.Errorf("Friend request had wrong user Id. Expected 422, found %v\n", resp.Requestors[1].Id)
        }
    }
}

func TestModifyMyFriendRequestEndpoint(t *testing.T) {
    defer resetTables()
    user1 := User{
        Id:        420,
        Uid:       "420",
        Name:      "Snoop Doge",
        FirstName: "Snoop",
        LastName:  "Doge",
        Email:     "higher@gmail.com",
        Picture:   "42keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    user2 := User{
        Id:        421,
        Uid:       "421",
        Name:      "Peppa Pig",
        FirstName: "Peppa",
        LastName:  "Pig",
        Email:     "p.pig@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    user3 := User{
        Id:        422,
        Uid:       "422",
        Name:      "Yo Mum",
        FirstName: "Yo",
        LastName:  "Mum",
        Email:     "top.kek@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 3")
    db.Create(&user3)

    log.Println("Add a friend request from user2 to user1")
    user1.addFriendRequest(user2)

    log.Println("Add a friend request from user3 to user1")
    user1.addFriendRequest(user3)

    log.Println("Trying to modify a request from yourself")
    resp := modifyMyFriendRequestEndpoint(user1, 420, "accept")
    if resp.Success || resp.Error == "" {
        t.Error("Succeeded in modifying reflective friend request")
    }

    log.Println("Trying to modify a request from a user that doesn't exist")
    resp = modifyMyFriendRequestEndpoint(user1, 1337, "accept")
    if resp.Success || resp.Error == "" {
        t.Error("Succeeded in modifying friend request from non existent user")
    }

    log.Println("Trying to modify a request that doesn't exist")
    resp = modifyMyFriendRequestEndpoint(user2, 420, "accept")
    if resp.Success || resp.Error == "" {
        t.Error("Succeeded in modifying friend request that doesn't exist")
    }

    log.Println("Accepting a friend request")
    resp = modifyMyFriendRequestEndpoint(user1, 421, "accept")
    if !resp.Success || resp.Error != "" {
        t.Error("Didn't succeed in modifying a valid friend request")
    } else {
        log.Println("Testing the request was deleted")
        if user1.hasFriendRequest(user2) {
            t.Error("Friend request wasn't deleted")
        }
        log.Println("Testing the users are now friends")
        if !user1.isFriend(user2) {
            t.Error("Users aren't now friends")
        }
    }

    log.Println("Declining a friend request")
    resp = modifyMyFriendRequestEndpoint(user1, 422, "decline")
    if !resp.Success || resp.Error != "" {
        t.Error("Didn't succeed in modifying a valid friend request")
    } else {
        log.Println("Testing the request was deleted")
        if user1.hasFriendRequest(user3) {
            t.Error("Friend request wasn't deleted")
        }
        log.Println("Testing the users are now not friends")
        if user1.isFriend(user3) {
            t.Error("Users are now friends when they shouldn't be")
        }
    }
}

func TestAddOthersFriendRequestEndpoint(t *testing.T) {
    defer resetTables()
    user1 := User{
        Id:        420,
        Uid:       "420",
        Name:      "Snoop Doge",
        FirstName: "Snoop",
        LastName:  "Doge",
        Email:     "higher@gmail.com",
        Picture:   "42keks"}

    log.Println("Creating test user 1")
    db.Create(&user1)

    user2 := User{
        Id:        421,
        Uid:       "421",
        Name:      "Peppa Pig",
        FirstName: "Peppa",
        LastName:  "Pig",
        Email:     "p.pig@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 2")
    db.Create(&user2)

    user3 := User{
        Id:        422,
        Uid:       "422",
        Name:      "Yo Mum",
        FirstName: "Yo",
        LastName:  "Mum",
        Email:     "top.kek@gmail.com",
        Picture:   "someurl"}

    log.Println("Creating test user 3")
    db.Create(&user3)

    log.Println("Add a friend request from user2 to user1")
    user1.addFriendRequest(user2)

    log.Println("Make user2 and user3 friends")
    user2.addFriend(user3)

    log.Println("Adding a friend request to a non-existent user")
    resp := addOthersFriendRequestEndpoint(user1, 24601)
    if resp.Success || resp.Error == "" {
        t.Error("Succeeded in adding a friend request to a non-existent user")
    }

    log.Println("Adding a friend request to a friend")
    resp = addOthersFriendRequestEndpoint(user2, 422)
    if resp.Success || resp.Error == "" {
        t.Error("Succeeded in adding a friend request to an existing friend")
    }

    log.Println("Adding a friend request that already exists")
    resp = addOthersFriendRequestEndpoint(user2, 420)
    if resp.Success || resp.Error == "" {
        t.Error("Succeeded in adding a friend request that already exists")
    }

    log.Println("Adding a friend request that already exists in the opposite direction")
    resp = addOthersFriendRequestEndpoint(user1, 421)
    if resp.Success || resp.Error == "" {
        t.Error("Succeeded in adding a friend request that already exists in the opposite direction")
    }

    log.Println("Adding a friend request to yourself")
    resp = addOthersFriendRequestEndpoint(user1, 420)
    if resp.Success || resp.Error == "" {
        t.Error("Succeeded in adding a friend request to yourself")
    }

    log.Println("Adding a valid friend request")
    resp = addOthersFriendRequestEndpoint(user3, 420)
    if !resp.Success || resp.Error != "" {
        t.Error("Didn't succeed in adding a valid friend request")
    }
    friendrequests := user1.getFriendRequests()
    if len(friendrequests) != 2 {
        t.Errorf("2 friend requests should have been found, found %v\n", len(friendrequests))
    } else {
        if friendrequests[0].Id != 421 {
            t.Errorf("Friend request had wrong user Id. Expected 421, found %v\n", friendrequests[0].Id)
        }
        if friendrequests[1].Id != 422 {
            t.Errorf("Friend request had wrong user Id. Expected 422, found %v\n", friendrequests[1].Id)
        }
    }
}

func TestGetNextMessageEndpoint(t *testing.T) {
    defer resetTables()
    user1 := User{
        Id:         1,
        Uid:        "1",
        Name:       "Tony Abbott",
        FirstName:  "Tony",
        LastName:   "Abbott",
        Email:      "xXx_0n10n_fan_xXx@hotmail.com",
        Picture:    "tone.jpg",
    }
    db.Create(&user1)

    user2 := User{
        Id:         2,
        Uid:        "2",
        Name:       "Malcolm Turnbull",
        FirstName:  "Malcolm",
        LastName:   "Turnbull",
        Email:      "pm@gmail.com",
        Picture:    "hehe",
    }
    db.Create(&user2)

    log.Println("Adding messages")
    msg1, _ := user1.addMessageToUser(user2, "malcom pls", ContentTypeText)
    msg2, _ := user2.addMessageToUser(user1, "lel", ContentTypeText)
    msg3, _ := user1.addMessageToUser(user2, "y u do dis", ContentTypeText)
    msg4, _ := user2.addMessageToUser(user1, "get rekt", ContentTypeText)

    log.Println("Checking getNextMessageAfterId for user1")
    if msg, ok := user1.getNextMessageAfterId(0); !ok || msg.Id != msg2.Id {
        t.Errorf("Message not found/wrong ID. Expected %v/%v, found %v/%v\n", true, msg2.Id, ok, msg.Id)
    }
    if msg, ok := user1.getNextMessageAfterId(msg2.Id); !ok || msg.Id != msg4.Id {
        t.Errorf("Message not found/wrong ID. Expected %v/%v, found %v/%v\n", true, msg4.Id, ok, msg.Id)
    }
    if msg, ok := user1.getNextMessageAfterId(msg4.Id); ok {
        t.Errorf("Message found when there shouldn't have been one. Found %v/%v\n", ok, msg.Id)
    }

    log.Println("Checking getNextMessageAfterId for user2")
    if msg, ok := user2.getNextMessageAfterId(0); !ok || msg.Id != msg1.Id {
        t.Errorf("Message not found/wrong ID. Expected %v/%v, found %v/%v\n", true, msg1.Id, ok, msg.Id)
    }
    if msg, ok := user2.getNextMessageAfterId(msg1.Id); !ok || msg.Id != msg3.Id {
        t.Errorf("Message not found/wrong ID. Expected %v/%v, found %v/%v\n", true, msg3.Id, ok, msg.Id)
    }
    if msg, ok := user2.getNextMessageAfterId(msg3.Id); ok {
        t.Errorf("Message found when there shouldn't have been one. Found %v/%v\n", ok, msg.Id)
    }
}

func TestTimeOfLastMessage(t *testing.T) {
    defer resetTables()
    user1 := User{
        Id:         1,
        Uid:        "1",
        Name:       "Tony Abbott",
        FirstName:  "Tony",
        LastName:   "Abbott",
        Email:      "xXx_0n10n_fan_xXx@hotmail.com",
        Picture:    "tone.jpg",
    }
    db.Create(&user1)

    user2 := User{
        Id:         2,
        Uid:        "2",
        Name:       "Malcolm Turnbull",
        FirstName:  "Malcolm",
        LastName:   "Turnbull",
        Email:      "pm@gmail.com",
        Picture:    "hehe",
    }
    db.Create(&user2)

    defaultTimestamp := time.Time{}

    log.Println("Testing time of last message with no messages")
    if ts := user1.timeOfLastMessageWithUser(user2); !ts.Equal(defaultTimestamp) {
        t.Errorf("Time of last message was not correct: expected %v, got %v", defaultTimestamp, ts)
    }
    if ts := user2.timeOfLastMessageWithUser(user1); !ts.Equal(defaultTimestamp) {
        t.Errorf("Time of last message was not correct: expected %v, got %v", defaultTimestamp, ts)
    }

    log.Println("Adding a message")
    msg1, _ := user1.addMessageToUser(user2, "malcom pls", ContentTypeText)
    timestamp1 := msg1.Timestamp.Round(time.Millisecond)

    log.Println("Testing time of last message after adding 1 message")
    if ts := user1.timeOfLastMessageWithUser(user2).Round(time.Millisecond); !ts.Equal(timestamp1) {
        t.Errorf("Time of last message was not correct: expected %v, got %v", timestamp1, ts)
    }
    if ts := user2.timeOfLastMessageWithUser(user1).Round(time.Millisecond); !ts.Equal(timestamp1) {
        t.Errorf("Time of last message was not correct: expected %v, got %v", timestamp1, ts)
    }

    time.Sleep(10 * time.Millisecond)

    log.Println("Adding another message")
    msg2, _ := user2.addMessageToUser(user1, "lel", ContentTypeText)
    timestamp2 := msg2.Timestamp.Round(time.Millisecond)

    log.Println("Testing time of last message after adding 2 messages")
    if ts := user1.timeOfLastMessageWithUser(user2).Round(time.Millisecond); !ts.Equal(timestamp2) {
        t.Errorf("Time of last message was not correct: expected %v, got %v", timestamp2, ts)
    }
    if ts := user2.timeOfLastMessageWithUser(user1).Round(time.Millisecond); !ts.Equal(timestamp2) {
        t.Errorf("Time of last message was not correct: expected %v, got %v", timestamp2, ts)
    }
}

func TestGetFriendsSorting(t *testing.T) {
    defer resetTables()
    user1 := User{
        Id:         1,
        Uid:        "1",
        Name:       "Tony Abbott",
        FirstName:  "Tony",
        LastName:   "Abbott",
        Email:      "xXx_0n10n_fan_xXx@hotmail.com",
        Picture:    "tone.jpg",
    }
    db.Create(&user1)

    user2 := User{
        Id:         2,
        Uid:        "2",
        Name:       "Malcolm Turnbull",
        FirstName:  "Malcolm",
        LastName:   "Turnbull",
        Email:      "pm@gmail.com",
        Picture:    "hehe",
    }
    db.Create(&user2)

    user3 := User{
        Id:        3,
        Uid:       "3",
        Name:      "Bill Shorten",
        FirstName: "Bill",
        LastName:  "Shorten",
        Email:     "billyboy@gmail.com",
        Picture:   "someurl",
    }
    db.Create(&user3)

    user4 := User{
        Id:         4,
        Uid:        "4",
        Name:       "Clive Palmer",
        FirstName:  "Clive",
        LastName:   "Palmer",
        Email:      "mining_is_gr8@gmail.com",
        Picture:    "coal.jpg",
    }
    db.Create(&user4)

    log.Println("Adding friends")
    user1.addFriend(user2)
    user1.addFriend(user3) // yeah right
    user1.addFriend(user4)

    log.Println("Testing ordering with no messages")
    friends := user1.getFriends()
    // should be alphabetical: bill (3), clive (4), malcolm (2)
    if friends[0].Id != user3.Id {
        t.Errorf("Friends weren't in alphabetical order: expected %v, got %v", user3.FirstName, friends[0].FirstName)
    }
    if friends[1].Id != user4.Id {
        t.Errorf("Friends weren't in alphabetical order: expected %v, got %v", user4.FirstName, friends[1].FirstName)
    }
    if friends[2].Id != user2.Id {
        t.Errorf("Friends weren't in alphabetical order: expected %v, got %v", user2.FirstName, friends[2].FirstName)
    }

    // send message from tony to malcolm
    user1.addMessageToUser(user2, "we're friends...right?", ContentTypeText)
    log.Println("Testing ordering with message to 1 friend")
    friends = user1.getFriends()
    // should be: malcolm (2), bill(3), clive(4)
    if friends[0].Id != user2.Id {
        t.Errorf("Friends weren't in order: expected %v, got %v", user2.FirstName, friends[0].FirstName)
    }
    if friends[1].Id != user3.Id {
        t.Errorf("Friends weren't in order: expected %v, got %v", user3.FirstName, friends[1].FirstName)
    }
    if friends[2].Id != user4.Id {
        t.Errorf("Friends weren't in order: expected %v, got %v", user4.FirstName, friends[2].FirstName)
    }

    // send message from clive to tony
    user4.addMessageToUser(user1, "bye bye tony, bye bye", ContentTypeText)
    log.Println("Testing ordering with message to 2 friends")
    friends = user1.getFriends()
    // should be: clive(4), malcolm(2), bill(3)
    if friends[0].Id != user4.Id {
        t.Errorf("Friends weren't in order: expected %v, got %v", user4.FirstName, friends[0].FirstName)
    }
    if friends[1].Id != user2.Id {
        t.Errorf("Friends weren't in order: expected %v, got %v", user2.FirstName, friends[1].FirstName)
    }
    if friends[2].Id != user3.Id {
        t.Errorf("Friends weren't in order: expected %v, got %v", user3.FirstName, friends[2].FirstName)
    }

    // send message from bill to malcolm (shouldn't affect ordering)
    user3.addMessageToUser(user2, "fuk", ContentTypeText)
    log.Println("Testing ordering with message to 2 friends, with 1 irrelevant message")
    friends = user1.getFriends()
    // should be: clive(4), malcolm(2), bill(3)
    if friends[0].Id != user4.Id {
        t.Errorf("Friends weren't in order: expected %v, got %v", user4.FirstName, friends[0].FirstName)
    }
    if friends[1].Id != user2.Id {
        t.Errorf("Friends weren't in order: expected %v, got %v", user2.FirstName, friends[1].FirstName)
    }
    if friends[2].Id != user3.Id {
        t.Errorf("Friends weren't in order: expected %v, got %v", user3.FirstName, friends[2].FirstName)
    }

    // send message from bill to tony (should affect ordering)
    user3.addMessageToUser(user1, "rip", ContentTypeText)
    log.Println("Testing ordering with message to 3 friends, with 1 irrelevant message")
    friends = user1.getFriends()
    // should be: bill(3), clive(4), malcolm(2)
    if friends[0].Id != user3.Id {
        t.Errorf("Friends weren't in order: expected %v, got %v", user3.FirstName, friends[0].FirstName)
    }
    if friends[1].Id != user4.Id {
        t.Errorf("Friends weren't in order: expected %v, got %v", user4.FirstName, friends[1].FirstName)
    }
    if friends[2].Id != user2.Id {
        t.Errorf("Friends weren't in order: expected %v, got %v", user2.FirstName, friends[2].FirstName)
    }
}
