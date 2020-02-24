package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestSupply(t *testing.T) {
	var supply Supply
	supply = DefaultSupply(defaultSymbol)

	require.Equal(t, defaultSymbol, supply.GetSymbol())
	require.Equal(t, sdk.ZeroInt(), supply.GetTotal())
	supply = supply.SetTotal(sdk.OneInt())
	require.Equal(t, sdk.OneInt(), supply.GetTotal())
	supply = supply.Inflate(sdk.OneInt())
	require.Equal(t, sdk.OneInt().Add(sdk.OneInt()), supply.GetTotal())
	supply = supply.Deflate(sdk.OneInt())
	require.Equal(t, sdk.OneInt(), supply.GetTotal())

	require.True(t, len(supply.String()) > 0)
}
