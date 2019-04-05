package models

import (
	"encoding/json"
	"time"
)

type Link struct {
	ID       interface{} `json:"id,omitempty" bson:"_id,omitempty"`
	Code     string      `json:"code" bson:"code"`
	LongUrl  string      `json:"longUrl" bson:"longUrl"`
	ShortUrl string      `json:"shortUrl,omitempty" bson:"shortUrl,omitempty"`
	Created  time.Time   `json:"created" bson:"created"`
}

func (link Link) EncodeLink() []byte {
	data, err := json.Marshal(link)
	if err != nil {
		panic(err)
	}
	return data
}

func DecodeLink(data []byte) (Link, error) {
	var link Link
	err := json.Unmarshal(data, &link)
	return link, err
}
