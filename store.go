package short

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Store interface {
	// Insert adds a record to the storage.
	// url, id and override are passed via an insertConfig struct.
	Insert(ctx context.Context, ic *insertConfig) error
	// GetUrl returns the url given an id.
	GetUrl(ctx context.Context, id string) (string, error)
}

type insertConfig struct {
	url      string
	id       string
	override bool
}

type store struct {
	name       string
	collection *mongo.Collection
}

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
		return nil, fmt.Errorf("failed to connect to %s: %w", mongoUri, err)
	}

	mongoDbClientMap[mongoUri] = c

	return c, nil
}

func getMongoCollection(ctx context.Context, database *mongo.Database, name string) (*mongo.Collection, error) {
	collectionsMap := database.Collection(collectionsMapName)

	// Index the collectionsMap collection.
	if _, err := collectionsMap.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, fmt.Errorf("failed to create an index for collection %s: %w", collectionsMapName, err)
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
		return nil, fmt.Errorf("failed to find and/or insert a collection mapping for %s: %w", name, res.Err())
	}

	var payload struct {
		CollectionName string `bson:"collectionName"`
	}

	if err := res.Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode a collection mapping document for %s: %w", name, err)
	}

	collection := database.Collection(payload.CollectionName)
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
		return nil, fmt.Errorf("failed to create an index for collection %s: %w", payload.CollectionName, err)
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

	cs, err := connstring.Parse(mongoUri)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", mongoUri, err)
	}
	if cs.Database == "" {
		cs.Database = "short"
	}

	database := client.Database(cs.Database)

	collection, err := getMongoCollection(ctx, database, name)
	if err != nil {
		return nil, err
	}

	return &store{
		name:       name,
		collection: collection,
	}, nil
}

func (s *store) Insert(ctx context.Context, ic *insertConfig) error {
	if ic.override {
		if _, err := s.collection.UpdateOne(
			ctx,
			bson.M{"id": ic.id},
			bson.M{"$set": bson.M{"id": ic.id, "url": ic.url}},
			options.Update().SetUpsert(true),
		); err != nil {
			return fmt.Errorf("failed to update or insert id %s: %w", ic.id, err)
		}
		return nil
	}

	if _, err := s.collection.InsertOne(ctx, bson.M{"id": ic.id, "url": ic.url}); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &ConflictError{}
		}
		return fmt.Errorf("failed to insert id %s: %w", ic.id, err)
	}

	return nil
}

func (s *store) GetUrl(ctx context.Context, id string) (string, error) {
	res := s.collection.FindOne(ctx, bson.M{"id": id})
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return "", &IdNotFoundError{id: id}
		}

		return "", fmt.Errorf("error when calling FindOne in the store %s: %w", s.name, res.Err())
	}

	var payload struct {
		Url string `json:"url"`
	}

	if err := res.Decode(&payload); err != nil {
		return "", fmt.Errorf("failed to decode record: %w", err)
	}

	return payload.Url, nil
}
