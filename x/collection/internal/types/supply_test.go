package types

import (
	"encoding/json"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestSupplyMarshalYAML(t *testing.T) {
	supply := DefaultSupply(defaultSymbol)
	coins := NewCoins(NewCoin(defaultTokenIDFT, sdk.OneInt()))
	supply = supply.Inflate(coins)

	bzCoins, err := json.Marshal(coins)
	require.NoError(t, err)

	want := fmt.Sprintf(`{"symbol":"%s","total":%s}`, defaultSymbol, string(bzCoins))

	require.Equal(t, want, supply.String())
}
