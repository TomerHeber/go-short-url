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

func (g *gsuConfig) WithAlias(alias string) GSUConfig {
	ng := *g
	ng.alias = alias
	if !isAlphaNumeric(alias) {
		ng.err = fmt.Errorf("alias %s contains non-alphanumeric characters", alias)
	}
	return &ng
}
