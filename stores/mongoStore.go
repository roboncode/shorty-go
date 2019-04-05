package stores

import (
	"context"
	"github.com/labstack/gommon/color"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"roboncode.com/go-urlshortener/models"
	"time"
)

var collectionName = "links"

type MongoStore struct {
	db *mongo.Database
}

func NewMongoStore() *MongoStore {
	m := MongoStore{}
	m.db = m.connect()
	m.ensureIndexes()
	return &m
}

func (m *MongoStore) connect() *mongo.Database {
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

func (m *MongoStore) ensureIndexes() {
	collection := m.db.Collection(collectionName)
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

func (m *MongoStore) formatShortUrl(link *models.Link) {
	link.ShortUrl = viper.GetString("baseUrl") + "/" + link.Code
}

func (m *MongoStore) IncCount() int64 {
	var counter models.Counter
	collection := m.db.Collection("counter")
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
	collection := m.db.Collection(collectionName)
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
	m.formatShortUrl(&link)
	return &link, nil
}

func (m *MongoStore) Read(code string) (*models.Link, error) {
	var link models.Link
	collection := m.db.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, bson.M{
		"code": code,
	}).Decode(&link)
	if err != nil {
		return nil, err
	}
	m.formatShortUrl(&link)
	return &link, nil
}

func (m *MongoStore) Delete(code string) int64 {
	collection := m.db.Collection(collectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.DeleteOne(ctx, bson.M{
		"code": code,
	})
	return result.DeletedCount
}

func (m *MongoStore) List(limit int64, skip int64) []models.Link {
	// https://danott.co/posts/json-marshalling-empty-slices-to-empty-arrays-in-go.html
	links := make([]models.Link, 0) // Do this to ensure empty array
	collection := m.db.Collection(collectionName)
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
			m.formatShortUrl(&link)
			links = append(links, link)
		}
	}

	return links
}
