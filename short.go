package short

type ShortenedURL interface {
	GetURL() string
}

type Shortener interface {
	GetShortenedURL(url string, config ...GSUConfig) (ShortenedURL, error)
}
