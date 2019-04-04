# Url Shortener built with Go

A simple URL shortener using Go and Mongo.

This project was built using [Echo](https://echo.labstack.com/) and the official [Go MongoDb Driver](https://github.com/mongodb/mongo-go-driver). Since the app is pretty simple, <300 lines, it exists as a single file.

## TLDR;

Use the Makefile to run docker

```
make start
```

Local development

```
make mongo
go run main.go
```

## API

The API is pretty simple.

```
Authentication required - uri?key=:authKey

POST    /shorten        body{ url:String }
GET     /urls
GET     /urls/newest 
GET     /urls/:code
DELETE  /urls/:code

No Authentication required

GET     /               Landing page
GET     /:code          Redirect to long url
GET     /*              404 page
```

Feel free to fork it, hack it and use it any way you please.

**MIT License**