package short

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/bson"
)

var mongoServer *memongo.Server

func TestMain(m *testing.M) {
	var err error

	mongoServer, err = memongo.Start("4.4.16")
	if err != nil {
		log.Fatal(err)
	}
	defer mongoServer.Stop()

	os.Exit(m.Run())
}

func getRandomMongoURIForTesting() string {
	return mongoServer.URIWithRandomDB()
}

func TestStore(t *testing.T) {
	t.Run("NewStore", func(t *testing.T) {
		uri := getRandomMongoURIForTesting()
		_, err := NewStore(uri, "short.link")
		require.Nil(t, err)
		_, err = NewStore(uri, "short.com")
		require.Nil(t, err)
		_, err = NewStore(uri, "short.link")
		require.Nil(t, err)

		client := mongoDbClientMap[uri]
		database := client.Database(databaseName)
		collections, _ := database.ListCollectionNames(context.Background(), bson.D{{}})
		require.Len(t, collections, 3)

		collectionsMapCollection := database.Collection(collectionsMapName)
		documentsCount, _ := collectionsMapCollection.CountDocuments(context.Background(), bson.D{{}})
		require.Equal(t, int64(2), documentsCount)
	})
}
