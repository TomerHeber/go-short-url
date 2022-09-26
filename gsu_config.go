package short

import "fmt"

type GSUConfig interface {
	getConfig() *gsuConfig

	WithAlias(alias string) GSUConfig
}

type gsuConfig struct {
	alias string

	err error
}

func DefaultGSUConfig() GSUConfig {
	return &gsuConfig{}
}

func (g *gsuConfig) getConfig() *gsuConfig {
	return g
}

func (g gsuConfig) WithAlias(alias string) GSUConfig {
	if !isAlphaNumeric(alias) {
		g.err = fmt.Errorf("alias %s contains non-alphanumeric characters", alias)
	} else {
		g.alias = alias
	}

	return &g
}
