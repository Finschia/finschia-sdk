package types

import (
	"fmt"
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestSupply(t *testing.T) {
	var supply Supply
	supply = DefaultSupply(defaultContractID)

	// create default
	require.Equal(t, defaultContractID, supply.GetContractID())
	require.Equal(t, sdk.ZeroInt(), supply.GetTotalSupply())
	require.Equal(t, sdk.ZeroInt(), supply.GetTotalMint())
	require.Equal(t, sdk.ZeroInt(), supply.GetTotalBurn())

	// set total supply
	initialSupply := sdk.NewInt(3)
	supply = supply.SetTotalSupply(initialSupply)
	require.Equal(t, initialSupply, supply.GetTotalSupply())
	require.Equal(t, initialSupply, supply.GetTotalMint())
	require.Equal(t, sdk.ZeroInt(), supply.GetTotalBurn())

	// inflate
	toInflate := sdk.NewInt(2)
	supply = supply.Inflate(toInflate)
	require.Equal(t, initialSupply.Add(toInflate), supply.GetTotalSupply())
	require.Equal(t, initialSupply.Add(toInflate), supply.GetTotalMint())
	require.Equal(t, sdk.ZeroInt(), supply.GetTotalBurn())

	// deflate
	toDeflate := sdk.NewInt(4)
	supply = supply.Deflate(toDeflate)
	require.Equal(t, initialSupply.Add(toInflate).Sub(toDeflate), supply.GetTotalSupply())
	require.Equal(t, initialSupply.Add(toInflate), supply.GetTotalMint())
	require.Equal(t, toDeflate, supply.GetTotalBurn())

	// total
	expected := fmt.Sprintf(
		`{"contract_id":"%s","total_supply":"%s","total_mint":"%s","total_burn":"%s"}`,
		defaultContractID,
		initialSupply.Add(toInflate).Sub(toDeflate),
		initialSupply.Add(toInflate),
		toDeflate,
	)
	require.Equal(t, expected, supply.String())
}
