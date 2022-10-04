package short

import (
	"fmt"
	"time"
)

type UrlConfig interface {
	getConfig() *urlConfig

	// WithAlias sets a short url alias instead of generating a random one.
	// E.g.: if the alias is `tastypizzas` the shortened url could be https://link.com/tastypizzas
	WithAlias(alias string) UrlConfig

	// WithOverrideAlias set the override configuration.
	// When override is `true` it will insert a new or override an existing shortened url.
	// This field is ignored when there is no alias.
	WithOverrideAlias(override bool) UrlConfig

	// WithExpirationDate sets an expiration date for the shortened url.
	// Once the expiration date has expired the url becomes invalid or allocated for other urls.
	WithExpirationDate(expriationDate time.Time) UrlConfig
}

type urlConfig struct {
	alias          string
	overrideAlias  bool
	expirationDate *time.Time

	err error
}

// DefaultConfig returns a configuration with default values.
// default alias: "" (empty string).
// default overrideAlias: false.
// default expirationDate: no expiration.
func DefaultUrlConfig() UrlConfig {
	return &urlConfig{}
}

func (u *urlConfig) getConfig() *urlConfig {
	return u
}

func (u urlConfig) WithAlias(alias string) UrlConfig {
	if !isAlphaNumeric(alias) {
		u.err = fmt.Errorf("alias %s contains non-alphanumeric characters", alias)
	} else {
		u.alias = alias
	}

	return &u
}

func (u urlConfig) WithOverrideAlias(override bool) UrlConfig {
	u.overrideAlias = override
	return &u
}

func (u urlConfig) WithExpirationDate(expriationDate time.Time) UrlConfig {
	u.expirationDate = &expriationDate
	return &u
}
