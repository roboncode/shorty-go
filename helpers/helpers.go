package helpers

import (
	"crypto/md5"
	"encoding/hex"
	c "github.com/roboncode/shorty-go/consts"
	"github.com/speps/go-hashids"
	"github.com/spf13/viper"
)

func GetShortUrl(code string) string {
	return viper.GetString(c.BaseUrl) + "/" + code
}

func MD5(text string) string {
	md5Hash := md5.New()
	md5Hash.Write([]byte(text))
	return hex.EncodeToString(md5Hash.Sum(nil))
}

func NewHashID() *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = viper.GetString(c.HashSalt)
	hd.MinLength = viper.GetInt(c.HashMin)
	h, _ := hashids.NewWithData(hd)
	return h
}
