package short

import "errors"

type ShortenedURL interface {
	GetURL() string
}

type Shortener interface {
	GetShortenedURL(url string, config ...GSUConfig) (ShortenedURL, error)
}

type shortner struct {
	host  string
	store Store
}

// NewShortener creates a new shortener.
// If no configuration is passed uses the default configuration (see: `DefaultConfig()`)
func NewShortener(config ...Config) (Shortener, error) {
	var c Config
	var err error

	if len(config) > 1 {
		return nil, errors.New("passed multiple configuration instances")
	}

	if len(config) == 0 {
		c = DefaultConfig()
	} else {
		c = config[0]
	}

	ci := c.getConfig()
	if ci.err != nil {
		return nil, ci.err
	}

	var s shortner

	s.host = ci.host
	s.store, err = NewStore(ci.mongoUri, ci.host)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (s *shortner) GetShortenedURL(url string, config ...GSUConfig) (ShortenedURL, error) {
	return nil, nil
}
