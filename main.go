package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
	"github.com/roboncode/go-urlshortener/stores"
	"github.com/speps/go-hashids"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strconv"
)

var store stores.Store
var h *hashids.HashID

func setupHashIds() *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = viper.GetString("hashSalt")
	hd.MinLength = viper.GetInt("hashMin")
	h, _ := hashids.NewWithData(hd)
	return h
}

func readConfig() {
	viper.SetDefault("store", "badger")
	viper.SetDefault("mongoUrl", "mongodb://localhost:27017")
	viper.SetDefault("database", "shorturls")
	viper.SetDefault("hashSalt", "shorturls")
	viper.SetDefault("hashMin", 5)
	viper.SetDefault("address", ":1323")
	viper.SetDefault("baseUrl", "")
	viper.SetDefault("authKey", "")
	viper.SetDefault("cacheExpMin", 15)
	viper.SetDefault("cacheCleanupMin", 60)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Println(color.Red("No config file found. Using default settings"))
	} else {
		log.Println(color.Green("Config found -- overriding defaults"))
	}

	mongoUrl := os.Getenv("MONGO_DB")
	if mongoUrl != "" {
		viper.Set("mongoUrl", mongoUrl)
	}
}

func main() {
	readConfig()
	h = setupHashIds()
	storeType := os.Getenv("STORE")
	if storeType == "" {
		storeType = viper.GetString("store")
	}
	switch storeType {
	case "mongo":
		store = stores.NewMongoStore()
		log.Println(color.Green("Connected Mongo database"))
	default:
		store = stores.NewBadgerStore()
		log.Println(color.Green("Connected Badger database"))
	}

	// Echo instance
	e := echo.New()

	// Middleware
	if os.Getenv("ENV") != "prod" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:key",
		Skipper: func(e echo.Context) bool {
			switch e.Path() {
			case "/", "/404", "/:code", "/*":
				return true
			}
			return false
		},
		Validator: func(key string, e echo.Context) (bool, error) {
			return key == viper.GetString("authKey"), nil
		},
	}))

	// Routes
	e.POST("/shorten", CreateLink)
	e.GET("/links", GetLinks)
	e.GET("/links/:code", GetLink)
	e.DELETE("/links/:code", DeleteLink)
	e.File("/", "public/index.html")
	e.File("/404", "public/404.html")
	e.GET("/:code", RedirectToUrl)
	e.File("/*", "public/404.html")

	// Start server
	e.Logger.Fatal(e.Start(viper.GetString("address")))
}

// Handlers
func CreateLink(c echo.Context) error {
	var body = new(struct {
		Url string `json:"url"`
	})

	if err := c.Bind(&body); err != nil {
		return err
	}

	if body.Url == "" {
		return c.JSON(http.StatusBadRequest, `Missing required property "url"`)
	}

	counter := int(store.IncCount())
	if code, err := h.Encode([]int{counter}); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else if link, err := store.Create(code, body.Url); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		return c.JSON(http.StatusOK, link)
	}
}

func GetLinks(c echo.Context) error {
	skip, _ := strconv.ParseInt(c.QueryParam("s"), 10, 64)
	limit, _ := strconv.ParseInt(c.QueryParam("l"), 10, 64)
	links := store.List(limit, skip)
	return c.JSON(http.StatusOK, links)
}

func GetLink(c echo.Context) error {
	if link, err := store.Read(c.Param("code")); err != nil {
		return c.NoContent(http.StatusNotFound)
	} else {
		return c.JSON(http.StatusOK, link)
	}
}

func DeleteLink(c echo.Context) error {
	if count := store.Delete(c.Param("code")); count == 0 {
		return c.NoContent(http.StatusNotFound)
	} else {
		return c.NoContent(http.StatusOK)
	}
}

func RedirectToUrl(c echo.Context) error {
	if link, err := store.Read(c.Param("code")); err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/404")
	} else {
		return c.Redirect(http.StatusMovedPermanently, link.LongUrl)
	}
}
