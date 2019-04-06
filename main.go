package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
	c "github.com/roboncode/go-urlshortener/consts"
	"github.com/roboncode/go-urlshortener/stores"
	"github.com/speps/go-hashids"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strconv"
)

// :: Internal ::
const ConnectedMongoMsg = "Connected Mongo database"
const ConnectedBadgerMsg = "Connected Badger database"
const NoConfigMsg = "No config file found. Using default settings"
const MissingRequiredUrlMsg = `Missing required property "url"`

var store stores.Store
var h *hashids.HashID

func setupHashIds() *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = viper.GetString(c.HashSalt)
	hd.MinLength = viper.GetInt(c.HashMin)
	h, _ := hashids.NewWithData(hd)
	return h
}

func readConfig() {
	_ = viper.BindEnv(c.AuthKey)
	_ = viper.BindEnv(c.BaseUrl)
	_ = viper.BindEnv(c.Env)
	_ = viper.BindEnv(c.HashMin)
	_ = viper.BindEnv(c.HashSalt)
	_ = viper.BindEnv(c.ServerAddr)
	_ = viper.BindEnv(c.Store)

	viper.SetDefault(c.AuthKey, "")
	viper.SetDefault(c.BaseUrl, "")
	viper.SetDefault(c.ServerAddr, ":8080")
	viper.SetDefault(c.Store, "badger")
	viper.SetDefault(c.HashSalt, "$h0rtur1$")
	viper.SetDefault(c.HashMin, 5)
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Println(color.Red(NoConfigMsg))
	}
}

func main() {
	readConfig()
	h = setupHashIds()
	switch viper.GetString(c.Store) {
	case "mongo":
		store = stores.NewMongoStore()
		log.Println(color.Green(ConnectedMongoMsg))
	default:
		store = stores.NewBadgerStore()
		log.Println(color.Green(ConnectedBadgerMsg))
	}

	// Echo instance
	e := echo.New()

	// Middleware
	if os.Getenv(c.Env) != "prod" {
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
			return key == viper.GetString(c.AuthKey), nil
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
	e.Logger.Fatal(e.Start(viper.GetString(c.ServerAddr)))
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
		return c.JSON(http.StatusBadRequest, MissingRequiredUrlMsg)
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
