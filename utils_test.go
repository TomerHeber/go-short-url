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
