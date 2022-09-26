package short

import (
	"net/url"
	"strings"
)

type Config interface {
	getConfig() *config
	// WithHost sets the hosts of the shortened url.
	// E.g. if host is `my.url`` the shortened url could be `https://my.url/eRt35df`.
	WithHost(host string) Config
}

type config struct {
	host string

	err error
}

// DefaultConfig returns a configuration with default values.
// default host: `localhost:8080`.
func DefaultConfig() Config {
	var c config

	c.host = "localhost:8080"

	return &c
}

func (c *config) getConfig() *config {
	return c
}

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
