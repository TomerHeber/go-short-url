package short

import "fmt"

type UrlConfig interface {
	getConfig() *urlConfig

	WithAlias(alias string) UrlConfig
	WithUpsertAlias(upsert bool) UrlConfig
}

type urlConfig struct {
	alias       string
	upsertAlias bool

	err error
}

// DefaultConfig returns a configuration with default values.
// default alias: "" (empty string).
// default upsert: false.
func DefaultUrlConfig() UrlConfig {
	return &urlConfig{}
}

func (u *urlConfig) getConfig() *urlConfig {
	return u
}

// WithAlias sets a short url alias instead of generating a random one.
// E.g.: if the alias is `tastypizzas` the shortened url could be https://link.com/tastypizzas
func (u urlConfig) WithAlias(alias string) UrlConfig {
	if !isAlphaNumeric(alias) {
		u.err = fmt.Errorf("alias %s contains non-alphanumeric characters", alias)
	} else {
		u.alias = alias
	}

	return &u
}

// WithUpsertAlias set the upsert configuration.
// When upsert is `true` it will insert a new or override an existing shortened url.
// This field is ignored when there is no alias.
func (u urlConfig) WithUpsertAlias(upsert bool) UrlConfig {
	u.upsertAlias = upsert
	return &u
}
