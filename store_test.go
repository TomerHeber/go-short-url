package short

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
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
	getDatabaseName := func(t *testing.T, uri string) string {
		cs, err := connstring.Parse(uri)
		require.Nil(t, err)
		return cs.Database
	}

	t.Run("NewStore", func(t *testing.T) {
		uri := getRandomMongoURIForTesting()
		_, err := NewStore(uri, "short.link")
		require.Nil(t, err)
		_, err = NewStore(uri, "short.com")
		require.Nil(t, err)
		_, err = NewStore(uri, "short.link")
		require.Nil(t, err)

		client := mongoDbClientMap[uri]
		databaseName := getDatabaseName(t, uri)
		database := client.Database(databaseName)
		collections, _ := database.ListCollectionNames(context.Background(), bson.D{{}})
		require.Len(t, collections, 3)

		collectionsMapCollection := database.Collection(collectionsMapName)
		documentsCount, _ := collectionsMapCollection.CountDocuments(context.Background(), bson.D{{}})
		require.Equal(t, int64(2), documentsCount)
	})

	t.Run("Insert", func(t *testing.T) {
		collectionName := "col1"

		getStoreHelper := func(t *testing.T) Store {
			t.Helper()
			uri := getRandomMongoURIForTesting()
			s, err := NewStore(uri, collectionName)
			require.Nil(t, err)
			return s
		}

		validateDocumentHelper := func(t *testing.T, s Store, id string, url string) {
			t.Helper()
			si := s.(*store)
			res := si.collection.FindOne(context.Background(), bson.M{"id": id})
			require.Nil(t, res.Err())

			var payload struct {
				Url string `bson:"url"`
			}

			err := res.Decode(&payload)
			require.Nil(t, err)
			require.Equal(t, url, payload.Url)
		}

		t.Run("Insert new", func(t *testing.T) {
			s := getStoreHelper(t)
			tm := time.Now().Add(time.Hour)
			err := s.Insert(context.Background(), &insertConfig{
				url:        "https://test.com",
				id:         "id",
				expiration: &tm,
			})
			require.Nil(t, err)
			validateDocumentHelper(t, s, "id", "https://test.com")
		})

		t.Run("Insert already exist", func(t *testing.T) {
			s := getStoreHelper(t)
			err := s.Insert(context.Background(), &insertConfig{
				url: "https://test.com",
				id:  "id",
			})
			require.Nil(t, err)
			err = s.Insert(context.Background(), &insertConfig{
				url: "https://test222.com",
				id:  "id",
			})
			require.ErrorIs(t, err, &ConflictError{})
			validateDocumentHelper(t, s, "id", "https://test.com")
		})

		t.Run("Override new", func(t *testing.T) {
			s := getStoreHelper(t)
			err := s.Insert(context.Background(), &insertConfig{
				url:      "https://test.com",
				id:       "id",
				override: true,
			})
			require.Nil(t, err)
			validateDocumentHelper(t, s, "id", "https://test.com")
		})

		t.Run("Override already exist", func(t *testing.T) {
			s := getStoreHelper(t)
			err := s.Insert(context.Background(), &insertConfig{
				url:      "https://test.com",
				id:       "id",
				override: true,
			})
			require.Nil(t, err)
			err = s.Insert(context.Background(), &insertConfig{
				url:      "https://test222.com",
				id:       "id",
				override: true,
			})
			require.Nil(t, err)
			validateDocumentHelper(t, s, "id", "https://test222.com")
		})
	})
}
