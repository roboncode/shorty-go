package stores

import (
	"context"
	"github.com/labstack/gommon/color"
	"github.com/patrickmn/go-cache"
	"github.com/roboncode/go-urlshortener/helpers"
	"github.com/roboncode/go-urlshortener/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"time"
)

// :: Config ::
const MongoCacheCleanup = "MONGO_CACHE_CLEANUP"
const MongoCacheExp = "MONGO_CACHE_EXP"
const MongoCounterCollection = "MONGO_COL_COUNTER"
const MongoLinksCollection = "MONGO_COL_LINKS"
const MongoUrl = "MONGO_URL"
const MongoDb = "MONGO_DB"

// :: Internal ::
const ConnectingMsg = "Connecting to Mongo database"

type MongoStore struct {
	db *mongo.Database
	c  *cache.Cache
}

func NewMongoStore() *MongoStore {
	viper.SetDefault(MongoCacheCleanup, 60)
	viper.SetDefault(MongoCacheExp, 15)
	viper.SetDefault(MongoCounterCollection, "counter")
	viper.SetDefault(MongoLinksCollection, "links")
	viper.SetDefault(MongoUrl, "mongodb://localhost:27017")
	viper.SetDefault(MongoDb, "shorturls")

	_ = viper.BindEnv(MongoCacheCleanup)
	_ = viper.BindEnv(MongoCacheExp)
	_ = viper.BindEnv(MongoCounterCollection)
	_ = viper.BindEnv(MongoLinksCollection)
	_ = viper.BindEnv(MongoUrl)
	_ = viper.BindEnv(MongoDb)

	m := MongoStore{}
	m.db = m.connect()
	m.ensureIndexes()

	m.c = cache.New(viper.GetDuration(MongoCacheExp)*time.Minute, viper.GetDuration(MongoCacheCleanup)*time.Minute)

	return &m
}

func (m *MongoStore) connect() *mongo.Database {
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString(MongoUrl)))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(color.Blue(ConnectingMsg))
	dbName := viper.GetString(MongoDb)
	return client.Database(dbName)
}

func (m *MongoStore) ensureIndexes() {
	collection := m.db.Collection(MongoLinksCollection)
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

func (m *MongoStore) IncCount() int64 {
	var counter models.Counter
	collection := m.db.Collection(MongoCounterCollection)
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

func (m *MongoStore) Create(code string, longUrl string) (*models.Link, error) {
	var link models.Link
	collection := m.db.Collection(MongoLinksCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	if err := collection.FindOne(ctx, bson.M{
		"longUrl": longUrl,
	}).Decode(&link); err != nil {
		link = models.Link{
			LongUrl: longUrl,
			Code:    code,
			Created: time.Now(),
		}
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
		res, err := collection.InsertOne(ctx, link)
		if err != nil {
			return nil, err
		}
		link.ID = res.InsertedID
	}
	helpers.FormatShortUrl(&link)
	return &link, nil
}

func (m *MongoStore) Read(code string) (*models.Link, error) {
	var link *models.Link
	var err error
	cachedItem, found := m.c.Get(code)
	if found {
		link = cachedItem.(*models.Link)
	} else {
		collection := m.db.Collection(MongoLinksCollection)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err = collection.FindOne(ctx, bson.M{
			"code": code,
		}).Decode(&link)
		if err != nil {
			return nil, err
		}
	}
	if link != nil {
		m.c.Set(code, link, cache.DefaultExpiration)
		helpers.FormatShortUrl(link)
	}
	return link, nil
}

func (m *MongoStore) Delete(code string) int64 {
	collection := m.db.Collection(MongoLinksCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.DeleteOne(ctx, bson.M{
		"code": code,
	})
	m.c.Delete(code)
	return result.DeletedCount
}

func (m *MongoStore) List(limit int64, skip int64) []models.Link {
	// https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
	links := make([]models.Link, 0) // Do this to ensure empty array
	collection := m.db.Collection(MongoLinksCollection)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := collection.Find(ctx, bson.M{}, &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	})

	defer cursor.Close(ctx)

	if err == nil {
		for cursor.Next(ctx) {
			var link models.Link
			_ = cursor.Decode(&link)
			helpers.FormatShortUrl(&link)
			links = append(links, link)
		}
	}

	return links
}
