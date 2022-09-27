package short

import (
	"context"
	"errors"
)

type ShortenedURL interface {
	GetURL() string
}

type Shortener interface {
	CreateShortenedURL(ctx context.Context, url string, config ...UrlConfig) (ShortenedURL, error)
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

func (s *shortner) insert(ctx context.Context, url string, id string, upsert bool) (ShortenedURL, error) {
	return nil, nil
}

// CreateShortenedURL creates a shortened url.
// If no configuration is passed uses the default configuration (see: `DefaultUrlConfig()`)
func (s *shortner) CreateShortenedURL(ctx context.Context, url string, config ...UrlConfig) (ShortenedURL, error) {
	if err := validateUrl(url); err != nil {
		return nil, err
	}

	var uc UrlConfig

	if len(config) > 1 {
		return nil, errors.New("passed multiple configuration instances")
	}

	if len(config) == 0 {
		uc = DefaultUrlConfig()
	} else {
		uc = config[0]
	}

	uci := uc.getConfig()
	if uci.err != nil {
		return nil, uci.err
	}

	if len(uci.alias) > 0 {
		return s.insert(ctx, url, uci.alias, uci.upsertAlias)
	}

	for {
		id, err := generateRandomId()
		if err != nil {
			return nil, err
		}

		shortenedUrl, err := s.insert(ctx, url, id, false)
		if err != nil {
			// In rare cases a conflict can occur. Generate a new random id.
			if errors.Is(err, &ConflictError{}) {
				continue
			}
			return nil, err
		}
		return shortenedUrl, nil
	}
}
