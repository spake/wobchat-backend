wobchat-backend
===============

You should have a config file in `/etc/wobchat-backend.conf` specifying things like your database settings. You can probably just use `wobchat-backend-example.conf` as-is, unless your dev environment is weird.


API Documentation
=================

Enum Values
===========
##`ContentType`
>A field in the Message struct. Defines what type of message the message is.
>
| Name            | Value | Description                   |
| --------------- |:-----:|:----------------------------- |
| ContentTypeText |   1   | The message is a text message |

##`RecipientType`
>A field in the Message struct. Defines what type of entity the message is being sent to.
>
| Name              | Value | Description                    |
| ----------------- |:-----:|:------------------------------ |
| RecipientTypeText |   1   | The recipient is a single user |



Endpoints
=========
Friends
-------

##`/friends`

###`GET`

>Gets a list of the current user's friends.
>
####Response Format:
    {
      "success": true,
      "friends": [
        {
          "id": 1,
          "uid": "123456789",
          "name": "Wayne Wobcke",
          "firstName": "Wayne",
          "lastName": "Wobcke",
          "picture": "https://lh6.googleusercontent.com/something/photo.jpg"
        }
      ]
    }

###`POST`

>Adds a user as a friend of the current user.
>
####Request Format:
    {
      "id":1
    }
>
####Response Format:
    {
      "success": true,
      "error": "",
      "friend": {
        "id": 2,
        "uid": "123456788",
        "name": "Shrek The Ogre",
        "firstName": "Shrek",
        "lastName": "The Ogre",
        "picture": "https://lh6.googleusercontent.com/something/photo.jpg"
      }
    }


##`/friends/{friendId}`

###`GET`

>Gets a friend of the current user by their Id.
>
####Response Format:
    {
      "success": true,
      "error": "",
      "friend": {
        "id": 3,
        "uid": "123456787",
        "name": "Snoop Dogg",
        "firstName": "Snoop",
        "lastName": "Dogg",
        "picture": "https://lh6.googleusercontent.com/something/photo.jpg"
      }
    }

###`DELETE`

>Removes a user from the current user's friends list.
>
####Response Format:
    {
      "success": true,
      "error": ""
    }

##`/friends/{friendId}/messages`

###`GET`

>Gets a list of the messages between the current user and their friend specified by the Id.
>
####Response Format:
    {
      "success": true,
      "error": "",
      "messages": [
        {
          "id": 2,
          "content": "Hey now, you're an all star.",
          "contentType": 1,
          "senderId": 2,
          "recipientId": 1,
          "recipientType": 1,
          "timestamp": "2015-09-23T02:14:29.945951+10:00"
        }
      ]
    }

###`POST`

>Sends a message from the current user to their friend specified by the Id.
>
####Request Format:
    {
      "content":"That's some good stuff right there.",
      "contentType":1
    }
>
####Response Format:
    {
      "success": true,
      "error": "",
      "id": 1
    }



Users
-----

##`/me`

###`GET`

> Gets information about the current user.
>
####Response Format:
    {
      "success": true,
      "error": "",
      "user": {
        "id": 1,
        "uid": "123456789",
        "name": "Wayne Wobcke",
        "firstName": "Wayne",
        "lastName": "Wobcke",
        "picture": "https://lh6.googleusercontent.com/something/photo.jpg"
      }
    }


##`/users[?q=partialname]`

###`GET`

>Gets a list of all users whose names match the given query.
>
####Response Format:
    {
      "success": true,
      "users": [
        {
          "id": 1,
          "uid": "123456789",
          "name": "Wayne Wobcke",
          "firstName": "Wayne",
          "lastName": "Wobcke",
          "picture": "https://lh6.googleusercontent.com/something/photo.jpg"
        }
      ]
    }

Events
------

##`/nextMessage?after={messageID}`

###`GET`

>Given the ID of the last message the client has seen (`after`), gets
>the next message that the client has not yet seen.
>If a new message already exists, then it is returned immediately;
>otherwise, the endpoint waits (for up to 60 seconds) for a message to
>be received.
>
####Response Format:
    {
      "success": true,
      "error": "",
      "message": ...
    }

    OR

    {
      "success": false,
      "error": "Timed out",
      "message": ...
    }
