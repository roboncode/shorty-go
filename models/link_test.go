package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLink_EncodeLink(t *testing.T) {
	link := Link{
		ID:       1,
		Code:     "abcde",
		LongUrl:  "https://roboncode.com",
		ShortUrl: "http://ac.me/abcde",
	}
	encoded := link.EncodeLink()
	assert.Equal(t, string(encoded), `{"id":1,"code":"abcde","longUrl":"https://roboncode.com","shortUrl":"http://ac.me/abcde","created":"0001-01-01T00:00:00Z"}`, "EncodeLink failed to match")
}

func TestDecodeLink(t *testing.T) {
	link := Link{
		ID:       1,
		Code:     "abcde",
		LongUrl:  "https://roboncode.com",
		ShortUrl: "http://ac.me/abcde",
	}
	encoded := link.EncodeLink()
	decoded, _ := DecodeLink(encoded)
	assert.Equal(t, decoded.Code, "abcde", "DecodeLink failed to code")
	assert.Equal(t, decoded.LongUrl, "https://roboncode.com", "DecodeLink failed to long url")
	assert.Equal(t, decoded.ShortUrl, "http://ac.me/abcde", "DecodeLink failed to short url")
}
