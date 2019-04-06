package helpers

import (
	"github.com/roboncode/go-urlshortener/consts"
	"github.com/roboncode/go-urlshortener/models"
	"github.com/spf13/viper"
	"testing"
)

func Test(t *testing.T) {
	viper.Set(consts.BaseUrl, "http://ac.me")

	link := models.Link{
		Code: "abcde",
	}
	FormatShortUrl(&link)
	if link.ShortUrl != "http://ac.me/abcde" {
		t.Errorf("FormatShortUrl(link) = %s; want http://ac.me/abcde", link.ShortUrl)
	}
}
