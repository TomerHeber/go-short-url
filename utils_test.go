package short

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsAlphaNumeric(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		require.True(t, isAlphaNumeric(""))
	})

	t.Run("alphanumeric", func(t *testing.T) {
		require.True(t, isAlphaNumeric("aBsds323aaAAZZ"))
	})

	t.Run("non-alphanumeric", func(t *testing.T) {
		require.False(t, isAlphaNumeric(" 1123aavvbb"))
	})
}

func TestValidateUrl(t *testing.T) {
	t.Run("no scheme", func(t *testing.T) {
		err := validateUrl("www.google.com?aaa=bb")
		require.NotNil(t, err)
		require.Contains(t, err.Error(), "the url must have an 'https' or an 'http' scheme")
	})

	t.Run("invalid", func(t *testing.T) {
		err := validateUrl("https://www com")
		require.NotNil(t, err)
	})

	t.Run("valid", func(t *testing.T) {
		err := validateUrl("https://www.com:4563/?aaa=fvvv#hash")
		require.Nil(t, err)
	})
}

func TestGeneratorRandomId(t *testing.T) {
	for i := 0; i < 100; i++ {
		id, err := generateRandomId()
		require.Nil(t, err)
		require.NotEmpty(t, id)
		require.LessOrEqual(t, len(id), 7)
		require.True(t, isAlphaNumeric(id))
	}
}
