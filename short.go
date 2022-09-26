package short

type ShortenedURL interface {
}

type Shortener interface {
	GetShortenedURL(url string, config ...GSUConfig) (ShortenedURL, error)
}
