package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
	c "github.com/roboncode/shorty-go/consts"
	"github.com/roboncode/shorty-go/handlers"
	"github.com/roboncode/shorty-go/helpers"
	"github.com/roboncode/shorty-go/stores"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	// :: Internal ::
	NoConfigMsg        = "No config file found. Using default settings"
	ConnectedMongoMsg  = "Connected Mongo database"
	ConnectedBadgerMsg = "Connected Badger database"
)

func initConfig() {
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

func NewStore(storeType string) stores.Store {
	var store stores.Store
	switch storeType {
	case "mongo":
		store = stores.NewMongoStore()
		log.Println(color.Green(ConnectedMongoMsg))
	default:
		store = stores.NewBadgerStore()
		log.Println(color.Green(ConnectedBadgerMsg))
	}
	return store
}

func main() {
	initConfig()

	// Echo instance
	e := echo.New()

	store := NewStore(viper.GetString(c.Store))
	hashID := helpers.NewHashID()

	h := handlers.Handler{
		Store:  store,
		HashID: hashID,
	}

	// Middleware
	if os.Getenv(c.Env) != "prod" {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.Recover())
	e.Use(middleware.Static("/public"))
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:key",
		Skipper: func(e echo.Context) bool {
			switch e.Path() {
			case "/shorten", "/links", "/links/:code":
				return false
			}
			return true
		},
		Validator: func(key string, e echo.Context) (bool, error) {
			return key == viper.GetString(c.AuthKey), nil
		},
	}))

	// Routes
	e.File("/admin", "public/admin/index.html")
	e.File("/404", "public/404.html")
	e.POST("/shorten", h.CreateLink)
	e.GET("/links", h.GetLinks)
	e.GET("/links/:code", h.GetLink)
	e.DELETE("/links/:code", h.DeleteLink)
	e.GET("/:code", h.RedirectToUrl)
	e.Static("/*", "public")

	// Start server
	e.Logger.Fatal(e.Start(viper.GetString(c.ServerAddr)))
}
