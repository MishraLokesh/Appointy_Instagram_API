# Appointy_Instagram_API

# Simple GO Lang REST API

> Simple Instagram API to create, read posts using MongoDB for the database implementation

## Quick Start

```bash
go build
./MishraLokesh
```

## Endpoints

### Get All users

```bash
GET /all_users
```

### Get Single user

```bash
POST /users/{id}
#pass the password in req body for verification
```

### Delete User

```bash
POST api/books/{id}
```

### Create new Post
``` bash
POST /users
# Request sample
# {
# Name: "Lokesh",
# Email: "user_One@gmail.com",
# Password: "yoyo"
# }
```

### Create new post
```bash
POST /posts
```

### Get all posts
```bash
POST /posts/{id}
```

### Create new Post
``` bash
POST /posts
# Request sample
# {
#  ID: 456
# Caption: "Hehe",
# ImageURL: "www.image_url.com",
# Timestamp: "123545680"
# }
```

### Get all Posts
``` bash
POST /posts/{id}

```

## App Info

### Author

Lokesh Mishra

### Version

1.0.0

### License

This project is licensed under the MIT License
```
