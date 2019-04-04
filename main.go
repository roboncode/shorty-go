package main

import (
	"context"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
	"github.com/speps/go-hashids"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Counter struct {
	Value int `bson:"value"`
}

type Link struct {
	ID       interface{} `json:"id,omitempty" bson:"_id,omitempty"`
	Url      string      `json:"url" bson:"url"`
	Code     string      `json:"code" bson:"code"`
	ShortUrl string      `json:"shortUrl,omitempty" bson:"shortUrl,omitempty"`
	Created  time.Time   `json:"created" bson:"created"`
}

var collectionName = "links"
var db *mongo.Database
var h *hashids.HashID

func getCounter() int {
	var counter Counter
	collection := db.Collection("counter")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.FindOneAndUpdateOptions{}
	opts.SetUpsert(true)
	opts.SetReturnDocument(options.ReturnDocument(options.After))
	err := collection.FindOneAndUpdate(ctx, bson.M{}, bson.M{"$inc": bson.M{"value": 1}}, &opts).Decode(&counter)
	if err != nil {
		return 0
	}
	return counter.Value
}

func connectToDatabase() *mongo.Database {
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("mongoUrl")))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(color.Green("Successfully connected to database"))
	dbName := viper.GetString("database")
	return client.Database(dbName)
}

func setupHashIds() *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = viper.GetString("hashSalt")
	hd.MinLength = viper.GetInt("hashMin")
	h, _ := hashids.NewWithData(hd)
	return h
}

func readConfig() {

	viper.SetDefault("mongoUrl", "mongodb://localhost:27017")
	viper.SetDefault("database", "shorturls")
	viper.SetDefault("hashSalt", "shorturls")
	viper.SetDefault("hashMin", 5)
	viper.SetDefault("address", ":1323")
	viper.SetDefault("baseUrl", "")
	viper.SetDefault("authKey", "")
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

func populateShortUrl(link *Link) {
	link.ShortUrl = viper.GetString("baseUrl") + "/" + link.Code
}

func ensureIndexes() {
	collection := db.Collection(collectionName)
	_, err := collection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"code", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true).SetBackground(true),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	readConfig()
	h = setupHashIds()
	db = connectToDatabase()
	ensureIndexes()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
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
	e.POST("/shorten", ShortenUrl)
	e.GET("/urls", GetUrls)
	e.GET("/urls/newest", GetNewestCode)
	e.GET("/urls/:code", GetCode)
	e.DELETE("/urls/:code", DeleteCode)
	e.File("/", "public/index.html")
	e.File("/404", "public/404.html")
	e.GET("/:code", RedirectToUrl)
	e.File("/*", "public/404.html")

	// Start server
	e.Logger.Fatal(e.Start(viper.GetString("address")))
}

// Handler
func ShortenUrl(c echo.Context) error {
	var body = new(struct {
		Url string `json:"url"`
	})
	if err := c.Bind(&body); err != nil {
		return err
	}

	if body.Url == "" {
		return c.JSON(http.StatusBadRequest, `Missing required property "url"`)
	}

	var link Link
	collection := db.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := collection.FindOne(ctx, bson.M{
		"url": body.Url,
	}).Decode(&link); err != nil {
		counter := getCounter()
		code, _ := h.Encode([]int{counter})
		link = Link{
			Url:     body.Url,
			Code:    code,
			Created: time.Now(),
		}
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		res, err := collection.InsertOne(ctx, link)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		link.ID = res.InsertedID
		populateShortUrl(&link)
		return c.JSON(http.StatusOK, link)
	}

	return c.JSON(http.StatusOK, link)

}

func GetUrls(c echo.Context) error {
	// https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
	links := make([]Link, 0) // Do this to ensure empty array
	collection := db.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	skip, _ := strconv.ParseInt(c.QueryParam("s"), 10, 64)
	limit, _ := strconv.ParseInt(c.QueryParam("l"), 10, 64)
	cursor, err := collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var link Link
		_ = cursor.Decode(&link)
		populateShortUrl(&link)
		links = append(links, link)
	}
	if err := cursor.Err(); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, links)
}

func GetNewestCode(c echo.Context) error {
	var link Link
	collection := db.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.M{}, &options.FindOneOptions{
		Sort: map[string]int{"created": -1},
	}).Decode(&link)
	if err != nil {
		var empty interface{}
		_ = c.JSON(http.StatusOK, empty)
		return err
	}
	populateShortUrl(&link)
	return c.JSON(http.StatusOK, link)
}

func GetCode(c echo.Context) error {
	var link Link
	collection := db.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.M{
		"code": c.Param("code"),
	}).Decode(&link)
	if err != nil {
		var empty interface{}
		_ = c.JSON(http.StatusOK, empty)
		return err
	}
	return c.JSON(http.StatusOK, link)
}

func DeleteCode(c echo.Context) error {
	collection := db.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, _ = collection.DeleteOne(ctx, bson.M{
		"code": c.Param("code"),
	})
	return c.NoContent(http.StatusOK)
}

func RedirectToUrl(c echo.Context) error {
	var link Link
	collection := db.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.M{
		"code": c.Param("code"),
	}).Decode(&link)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/404")
	}
	return c.Redirect(http.StatusMovedPermanently, link.Url)
}
