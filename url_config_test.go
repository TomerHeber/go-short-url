package short

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUrlConfig(t *testing.T) {
	t.Run("with empty alias", func(t *testing.T) {
		c := DefaultUrlConfig().WithAlias("")
		require.Nil(t, c.getConfig().err)
	})

	t.Run("with alias", func(t *testing.T) {
		alias := "asdas21312FGSDFccxs"
		c := DefaultUrlConfig().WithAlias(alias)
		require.Nil(t, c.getConfig().err)
		require.Equal(t, alias, c.getConfig().alias)
	})

	t.Run("with invalid alias", func(t *testing.T) {
		alias := "asdas21312F$GSDFccxs"
		c := DefaultUrlConfig().WithAlias(alias)
		require.NotNil(t, c.getConfig().err)
		require.Equal(t, "", c.getConfig().alias)
	})
}
