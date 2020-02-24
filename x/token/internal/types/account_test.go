package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	var acc Account
	acc = NewBaseAccountWithAddress(defaultSymbol, addr1)

	require.Equal(t, defaultSymbol, acc.GetSymbol())
	require.Equal(t, addr1, acc.GetAddress())
	require.Equal(t, sdk.ZeroInt(), acc.GetBalance())

	acc = acc.SetBalance(sdk.OneInt())

	require.Equal(t, sdk.OneInt(), acc.GetBalance())

	require.True(t, len(acc.String()) > 0)
}
