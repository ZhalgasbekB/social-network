#  Social-Network-Test

## Introduction

This is a simple social network API that allows users to create profiles, create posts, like posts, and get notifications.
Where I create API using Golang and MongoDB.

## Curl Commands for API Endpoints

1. Get User Profile
```bash
    curl -X GET http://localhost:8080/profile \
    -H "Content-Type: application/json" \
    -d '{"user_id": 1}'
```   

2. Create User Profile
```bash
    curl -X POST http://localhost:8080/profile \
    -H "Content-Type: application/json" \
    -d '{"name": "John Doe", "avatar": "url"}'
```
 
3. Create a New Post
```bash
    curl -X POST http://localhost:8080/posts \
    -H "Content-Type: application/json" \
    -d '{"content": "New Post", "author": 1}'
``` 

4. Get User's Posts
```bash
    curl -X GET http://localhost:8080/posts \
    -H "Content-Type: application/json" \
    -d '{"user_id": 1}'
```

5. Get Liked Posts
```bash
    curl -X GET http://localhost:8080/posts/liked \
    -H "Content-Type: application/json" \
    -d '{"user_id": 1}'
```

6. Like a Post
```bash
    curl -X POST http://localhost:8080/posts/2/like \
    -H "Content-Type: application/json" \
    -d '{"user_id": 1}'
```

7. Get Notifications
```bash
    curl -X GET http://localhost:8080/notifications\
    -H "Content-Type: application/json" \
    -d '{"user_id": 1}'
```

All the above commands are tested using curl and Postman. 
And you see a many times "user_id" for validate user because we can user another ways like JWT 
tokens or cookies, but I guess this is the simplest way to validate user. 

And as previously seen, I use CLI like curl for testing the API endpoints.
Curl I installed by homebrew on my Mac. 

```bash
    brew install curl
```
But Postman I also use for testing the endpoints.
Postman I use for GUI test. 

## How to Run the Code

If you want, run the code by local change the database configuration by Mongo URL 
and run the below command to start the server. 
```bash
    go run ./cmd
```
