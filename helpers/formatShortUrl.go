package helpers

import (
	"github.com/roboncode/go-urlshortener/consts"
	"github.com/roboncode/go-urlshortener/models"
	"github.com/spf13/viper"
)

func FormatShortUrl(link *models.Link) {
	link.ShortUrl = viper.GetString(consts.BaseUrl) + "/" + link.Code
}
