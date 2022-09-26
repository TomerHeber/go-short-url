package short

import (
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type Config interface {
	getConfig() *config

	// WithHost sets the hosts of the shortened url.
	// E.g. if host is `my.url`` the shortened url could be `https://my.url/eRt35df`.
	WithHost(host string) Config

	// WithMongo sets the URI for connecting to Mongo.
	// https://www.mongodb.com/docs/manual/reference/connection-string/
	// Example: `mongodb://root:password123@198.174.21.23:27017/databasename`
	WithMongoUri(mongoUri string) Config
}

type config struct {
	host     string
	mongoUri string

	err error
}

// DefaultConfig returns a configuration with default values.
// default host: `localhost:8080`.
// default mongo URI: `mongodb://localhost:27017`.
func DefaultConfig() Config {
	var c config

	c.host = "localhost:8080"
	c.mongoUri = "mongodb://localhost:27017"

	return &c
}

func (c *config) getConfig() *config {
	return c
}

// WithHost set the short link host.
func (c config) WithHost(host string) Config {
	if !strings.HasPrefix(host, "https://") && !strings.HasPrefix(host, "http://") {
		host = "https://" + host
	}
	u, err := url.ParseRequestURI(host)

	if err != nil {
		c.err = err
	} else {
		c.host = u.Host
	}

	return &c
}

// WithMongoUri set the URI for connecting to MongoDB.
func (c config) WithMongoUri(mongoUri string) Config {
	if _, err := connstring.ParseAndValidate(mongoUri); err != nil {
		c.err = err
	} else {
		c.mongoUri = mongoUri
	}

	return &c
}
