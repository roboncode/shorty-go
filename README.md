# Url Shortener built with Go

A simple URL shortener using Go and Mongo.

This project was built using [Echo](https://echo.labstack.com/) and offers two data stores:

* [Badger](https://github.com/dgraph-io/badger) - Embedded Go Key/Value Database for a simple standalone exec using 
* [Mongo](https://github.com/mongodb/mongo-go-driver) - To handle multiple services such as in a Kubernetes environment


## Develoment

Use the Makefile to run docker

```
make start
```

Local development using Badger

```
make mongo
go run main.go
```

## API

The API is pretty simple.

```
Authentication required - uri?key=:authKey

POST    /shorten        body{ url:String }
GET     /links?l=:limit&s=:skip
GET     /links/newest 
GET     /links/:code
DELETE  /links/:code

No Authentication required

GET     /               Landing page
GET     /:code          Redirect to long url
GET     /*              404 page
```

Feel free to fork it, hack it and use it any way you please.

**MIT License**