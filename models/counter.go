package models

import "encoding/json"

type Counter struct {
	Value int `json:"value" bson:"value"`
}

func (c *Counter) EncodeCounter() []byte {
	data, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}
	return data
}

func DecodeCounter(data []byte) (Counter, error) {
	var c Counter
	err := json.Unmarshal(data, &c)
	return c, err
}
