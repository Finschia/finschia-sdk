package types

import (
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	var acc Account
	acc = NewBaseAccountWithAddress(defaultContractID, addr1)

	require.Equal(t, defaultContractID, acc.GetContractID())
	require.Equal(t, addr1, acc.GetAddress())
	require.Equal(t, sdk.ZeroInt(), acc.GetBalance())

	acc = acc.SetBalance(sdk.OneInt())

	require.Equal(t, sdk.OneInt(), acc.GetBalance())

	require.True(t, len(acc.String()) > 0)
}
