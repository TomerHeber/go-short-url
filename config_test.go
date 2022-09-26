package short

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("WithHost", func(t *testing.T) {
		t.Run("default", func(t *testing.T) {
			c := DefaultConfig()
			require.Equal(t, "localhost:8080", c.getConfig().host)
		})

		helperValid := func(t *testing.T, host string) {
			t.Helper()
			c := DefaultConfig().WithHost(host)
			require.Nil(t, c.getConfig().err)
			require.Equal(t, host, c.getConfig().host)
		}

		helperInvalid := func(t *testing.T, host string) {
			t.Helper()
			c := DefaultConfig().WithHost(host)
			require.NotNil(t, c.getConfig().err)
		}

		t.Run("valid english", func(t *testing.T) {
			helperValid(t, "url.com")
		})

		t.Run("valid english + port", func(t *testing.T) {
			helperValid(t, "url.com:1234")
		})

		t.Run("valid non-english", func(t *testing.T) {
			helperValid(t, "ヒキワリ.ナットウ.ニホン")
		})

		t.Run("valid non-english + port", func(t *testing.T) {
			helperValid(t, "ヒキワリ.ナットウ.ニホン:80")
		})

		t.Run("invalid hostname", func(t *testing.T) {
			helperInvalid(t, "myurl%.com")
		})

		t.Run("invalid port", func(t *testing.T) {
			helperInvalid(t, "myurl.com:fsdfsdf")
		})

		t.Run("valid hostname with https schema", func(t *testing.T) {
			c := DefaultConfig().WithHost("https://myurl.com")
			require.Nil(t, c.getConfig().err)
			require.Equal(t, "myurl.com", c.getConfig().host)
		})

		t.Run("valid hostname with http schema", func(t *testing.T) {
			c := DefaultConfig().WithHost("http://myurl.com")
			require.Nil(t, c.getConfig().err)
			require.Equal(t, "myurl.com", c.getConfig().host)
		})
	})
}
