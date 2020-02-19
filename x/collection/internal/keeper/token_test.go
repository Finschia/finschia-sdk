package keeper

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_ModifyTokenURI(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Prepare Collection and Token. Add Modify Permission to Addr1")
	{
		store := ctx.KVStore(keeper.storeKey)
		collection := types.NewCollection(defaultSymbol, defaultName)
		store.Set(types.CollectionKey(collection.GetSymbol()), keeper.cdc.MustMarshalBinaryBare(collection))
		token := types.NewFT(collection, defaultName, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
		err := keeper.SetToken(ctx, token)
		require.NoError(t, err)
		keeper.AddPermission(ctx, addr1, types.NewModifyTokenURIPermission(defaultSymbol+defaultTokenIDFT))
	}

	const modifiedTokenURI = "modifiedtokenuri"
	t.Log("Modify Token URI")
	{
		require.NoError(t, keeper.ModifyTokenURI(ctx, addr1, defaultSymbol, defaultTokenIDFT, modifiedTokenURI))
	}
	t.Log("Compare Token")
	{
		actual, err := keeper.GetFT(ctx, defaultSymbol, defaultTokenIDFT)
		require.NoError(t, err)
		require.Equal(t, modifiedTokenURI, actual.GetTokenURI())
	}
}
