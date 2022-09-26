package short

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGSUConfig(t *testing.T) {
	t.Run("with empty alias", func(t *testing.T) {
		c := DefaultGSUConfig().WithAlias("")
		require.Nil(t, c.getConfig().err)
	})

	t.Run("with alias", func(t *testing.T) {
		alias := "asdas21312FGSDFccxs"
		c := DefaultGSUConfig().WithAlias(alias)
		require.Nil(t, c.getConfig().err)
		require.Equal(t, alias, c.getConfig().alias)
	})

	t.Run("with invalid alias", func(t *testing.T) {
		alias := "asdas21312F$GSDFccxs"
		c := DefaultGSUConfig().WithAlias(alias)
		require.NotNil(t, c.getConfig().err)
		require.Equal(t, "", c.getConfig().alias)
	})
}
