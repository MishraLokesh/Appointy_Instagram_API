# Appointy_Instagram_API

# Simple GO Lang REST API

> Simple RESTful API to create, read, update and delete books. No database implementation yet

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
pass the password in req body for verification
```

### Delete Book

```bash
POST api/books/{id}
```

### Create new user

````bash
POST /users
````

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
# Name: "Lokesh",
# Email: "user_One@gmail.com",
# Password: "yoyo"
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
