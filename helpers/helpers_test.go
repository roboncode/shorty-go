package helpers

import (
	"github.com/roboncode/shorty-go/consts"
	"github.com/roboncode/shorty-go/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetShortUrl(t *testing.T) {
	viper.Set(consts.BaseUrl, "https://ac.me")

	link := models.Link{
		Code: "abcde",
	}
	link.ShortUrl = GetShortUrl(link.Code)
	if link.ShortUrl != "https://ac.me/abcde" {
		t.Errorf("GetShortUrl(%s) = %s; want https://ac.me/abcde", link.Code, link.ShortUrl)
	}
}

func TestMD5(t *testing.T) {
	hash := MD5("Hello, world!")
	assert.Equal(t, hash, "6cd3556deb0da54bca060b4c39479839", "MD5 hash failed to match")
}
