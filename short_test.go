package short

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestShort(t *testing.T) {
	t.Run("CreateShortenedUrl and GetUrlFromShortenedUrl", func(t *testing.T) {
		host := "host.com:12345"
		url := "https://test.com/?sdfsdfsd"

		shortner, err := NewShortener(DefaultConfig().WithHost(host).WithMongoUri(getRandomMongoURIForTesting()))
		require.Nil(t, err)

		t.Run("valid url", func(t *testing.T) {
			surl, err := shortner.CreateShortenedUrl(context.Background(), url)
			require.Nil(t, err)
			require.Contains(t, surl.GetUrl(), "https://host.com:12345/")

			rurl, err := shortner.GetUrlFromShortenedUrl(context.Background(), surl.GetUrl())
			require.Nil(t, err)
			require.Equal(t, url, rurl)
		})

		t.Run("invalid url", func(t *testing.T) {
			_, err := shortner.CreateShortenedUrl(context.Background(), "http1s://test.com/?sdfsdfsd")
			require.Error(t, err)
		})

		t.Run("invalid short url", func(t *testing.T) {
			surl, err := shortner.CreateShortenedUrl(context.Background(), url)
			require.Nil(t, err)

			_, err = shortner.GetUrlFromShortenedUrl(context.Background(), surl.GetUrl()+"!")
			require.Error(t, err)
			require.Contains(t, err.Error(), "invalid short url path")
		})

		t.Run("id not found", func(t *testing.T) {
			_, err = shortner.GetUrlFromShortenedUrl(context.Background(), "https://host.com:12345/aaaaaaa")
			require.Error(t, err)
			var perr *IdNotFoundError
			require.ErrorAs(t, err, &perr)
			require.Contains(t, err.Error(), "aaaaaaa")
		})

		t.Run("expired short url", func(t *testing.T) {
			surl, err := shortner.CreateShortenedUrl(context.Background(), url, DefaultUrlConfig().WithExpirationDate(time.Now().Add(-time.Hour)))
			require.Nil(t, err)

			_, err = shortner.GetUrlFromShortenedUrl(context.Background(), surl.GetUrl())
			require.Error(t, err)
			var perr *IdNotFoundError
			require.ErrorAs(t, err, &perr)
		})
	})
}
