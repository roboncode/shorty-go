package models

import "time"

type Link struct {
	ID       interface{} `json:"id,omitempty" bson:"_id,omitempty"`
	Code     string      `json:"code" bson:"code"`
	LongUrl  string      `json:"longUrl" bson:"longUrl"`
	ShortUrl string      `json:"shortUrl,omitempty" bson:"shortUrl,omitempty"`
	Created  time.Time   `json:"created" bson:"created"`
}
