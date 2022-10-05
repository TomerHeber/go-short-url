package short

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type ShortenedURL interface {
	GetUrl() string
}

type Shortener interface {
	// CreateShortenedUrl creates a shortened url for `url`.
	CreateShortenedUrl(ctx context.Context, url string, config ...UrlConfig) (ShortenedURL, error)
	// GetUrlFromShortenedUrl receives a shortened url `surl` and returns the original url.
	// E.g.: https://short.com/abCD123
	GetUrlFromShortenedUrl(ctx context.Context, surl string) (string, error)
	// GetUrlFromShortenedUrl receives a shortened url `id` and returns the original url.
	// E.g.: abCD123
	GetUrlFromShortenedUrlId(ctx context.Context, id string) (string, error)
}

type shortner struct {
	host  string
	store Store
}

type shortenedUrl struct {
	url string
}

func (s *shortenedUrl) GetUrl() string {
	return s.url
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

func newShortenedUrl(id string, host string) ShortenedURL {
	scheme := "https://"
	if strings.HasPrefix(host, "localhost") {
		scheme = "http://"
	}

	return &shortenedUrl{
		url: scheme + host + "/" + id,
	}
}

func (s *shortner) insert(ctx context.Context, ic *insertConfig) (ShortenedURL, error) {
	if err := s.store.Insert(ctx, ic); err != nil {
		return nil, fmt.Errorf("failed to insert an entry for a shortened url: %w", err)
	}

	return newShortenedUrl(ic.id, s.host), nil
}

// CreateShortenedUrl creates a shortened url.
// If no configuration is passed uses the default configuration (see: `DefaultUrlConfig()`)
func (s *shortner) CreateShortenedUrl(ctx context.Context, url string, config ...UrlConfig) (ShortenedURL, error) {
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
		return s.insert(ctx, &insertConfig{
			url: url, id: uci.alias, override: uci.overrideAlias,
		})
	}

	for {
		id, err := generateRandomId()
		if err != nil {
			return nil, err
		}

		shortenedUrl, err := s.insert(ctx, &insertConfig{
			url: url, id: id, override: false,
		})
		if err != nil {
			// In rare cases (statistically) a conflict may occur. Generate a new random id.
			if errors.Is(err, &ConflictError{}) {
				continue
			}
			return nil, err
		}

		return shortenedUrl, nil
	}
}

func (s *shortner) GetUrlFromShortenedUrl(ctx context.Context, surl string) (string, error) {
	su, err := url.ParseRequestURI(surl)
	if err != nil {
		return "", fmt.Errorf("invalid short url %s: %w", surl, err)
	}

	id := strings.Trim(su.Path, "/")

	return s.GetUrlFromShortenedUrlId(ctx, id)
}

func (s *shortner) GetUrlFromShortenedUrlId(ctx context.Context, id string) (string, error) {
	if len(id) == 0 || !isAlphaNumeric(id) {
		return "", fmt.Errorf("invalid short url path %s", id)
	}

	return s.store.GetUrl(ctx, id)
}
