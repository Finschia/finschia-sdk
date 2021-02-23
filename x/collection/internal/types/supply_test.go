package types

import (
	"encoding/json"
	"fmt"
	"testing"

	sdk "github.com/line/lbm-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestSupply(t *testing.T) {
	var supply Supply
	supply = DefaultSupply(defaultTokenIDFT)

	// create default
	require.Equal(t, defaultTokenIDFT, supply.GetContractID())
	require.Equal(t, NewCoins(), supply.GetTotalSupply())
	require.Equal(t, NewCoins(), supply.GetTotalMint())
	require.Equal(t, NewCoins(), supply.GetTotalBurn())

	// set total supply
	initialSupply := NewCoins(NewCoin(defaultTokenIDFT, sdk.NewInt(3)))
	supply = supply.SetTotalSupply(initialSupply)
	require.Equal(t, initialSupply, supply.GetTotalSupply())
	require.Equal(t, initialSupply, supply.GetTotalMint())
	require.Equal(t, NewCoins(), supply.GetTotalBurn())

	// inflate
	toInflate := NewCoins(NewCoin(defaultTokenIDFT, sdk.NewInt(2)))
	supply = supply.Inflate(toInflate)
	require.Equal(t, initialSupply.Add(toInflate...), supply.GetTotalSupply())
	require.Equal(t, initialSupply.Add(toInflate...), supply.GetTotalMint())
	require.Equal(t, NewCoins(), supply.GetTotalBurn())

	// deflate
	toDeflate := NewCoins(NewCoin(defaultTokenIDFT, sdk.NewInt(4)))
	supply = supply.Deflate(toDeflate)
	require.Equal(t, initialSupply.Add(toInflate...).Sub(toDeflate), supply.GetTotalSupply())
	require.Equal(t, initialSupply.Add(toInflate...), supply.GetTotalMint())
	require.Equal(t, toDeflate, supply.GetTotalBurn())

	// total
	ts, err1 := json.Marshal(NewCoins(NewCoin(defaultTokenIDFT, sdk.NewInt(1))))
	require.NoError(t, err1)
	tm, err2 := json.Marshal(NewCoins(NewCoin(defaultTokenIDFT, sdk.NewInt(5))))
	require.NoError(t, err2)
	tb, err3 := json.Marshal(NewCoins(NewCoin(defaultTokenIDFT, sdk.NewInt(4))))
	require.NoError(t, err3)
	expected := fmt.Sprintf(
		`{"contract_id":"%s","total_supply":%v,"total_mint":%v,"total_burn":%v}`,
		defaultTokenIDFT,
		string(ts),
		string(tm),
		string(tb),
	)
	require.Equal(t, expected, supply.String())
}

func TestSupplyMarshalYAML(t *testing.T) {
	supply := DefaultSupply(defaultContractID)
	coins := NewCoins(NewCoin(defaultTokenIDFT, sdk.OneInt()))
	supply = supply.Inflate(coins)

	bzCoins, err := json.Marshal(coins)
	require.NoError(t, err)

	expected := fmt.Sprintf(
		`{"contract_id":"%s","total_supply":%s,"total_mint":%s,"total_burn":%s}`,
		defaultContractID,
		string(bzCoins),
		string(bzCoins),
		[]Coin{},
	)

	require.Equal(t, expected, supply.String())
}
