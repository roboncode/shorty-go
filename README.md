# Url Shortener built with Go

A simple URL shortener using Go and Mongo.

This project was built using [Echo](https://echo.labstack.com/) and the official [Go MongoDb Driver](https://github.com/mongodb/mongo-go-driver).

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
POST    /shorten        body{ url:String }
GET     /urls
GET     /urls/newest   
GET     /urls/:code
DELETE  /urls/:code?apiKey=:apikey
GET     /               Landing page
GET     /:code          Redirect to long url
GET     /*              404 page
```

Feel free to fork and build on it. It works for my purposes.

**MIT License**