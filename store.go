package short

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Store interface {
	Insert(ctx context.Context, ic *insertConfig) error
}

type insertConfig struct {
	url    string
	id     string
	upsert bool
}

type store struct {
	name       string
	collection *mongo.Collection
}

const databaseName = "short"
const collectionsMapName = "collections_map"

// used as a cache to store MongoDB clients.
var mongoDbClientMap = map[string]*mongo.Client{}
var mongoDbClientMapLock sync.Mutex

func getMongoClient(ctx context.Context, mongoUri string) (*mongo.Client, error) {
	mongoDbClientMapLock.Lock()
	defer mongoDbClientMapLock.Unlock()

	if c, ok := mongoDbClientMap[mongoUri]; ok {
		return c, nil
	}

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, err
	}

	mongoDbClientMap[mongoUri] = c

	return c, nil
}

func getMongoCollection(ctx context.Context, client *mongo.Client, name string) (*mongo.Collection, error) {
	collectionsMap := client.Database(databaseName).Collection(collectionsMapName)

	// Index the collectionsMap collection.
	if _, err := collectionsMap.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}

	id := uuid.New().String()

	// Find the collection name in the collections map.
	res := collectionsMap.FindOneAndUpdate(
		ctx,
		bson.M{"name": name},
		bson.M{"$setOnInsert": bson.M{"collectionName": id, "name": name}},
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var payload struct {
		CollectionName string `bson:"collectionName"`
	}

	if err := res.Decode(&payload); err != nil {
		return nil, err
	}

	collection := client.Database(databaseName).Collection(payload.CollectionName)
	// Index the collection.
	if _, err := collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.M{"url": 1},
		},
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"expireAt": 1},
			Options: options.Index().SetExpireAfterSeconds(0),
		}}); err != nil {
		return nil, err
	}

	return collection, nil
}

func NewStore(mongoUri string, name string) (Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := getMongoClient(ctx, mongoUri)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.PrimaryPreferred()); err != nil {
		return nil, err
	}

	collection, err := getMongoCollection(ctx, client, name)
	if err != nil {
		return nil, err
	}

	return &store{
		name:       name,
		collection: collection,
	}, nil
}

func (s *store) Insert(ctx context.Context, ic *insertConfig) error {
	return nil
}
