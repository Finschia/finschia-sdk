package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_SendCoins(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	require.EqualError(t, keeper.SendCoins(ctx, defaultContractID, addr3, addr1, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), types.ErrInsufficientToken(types.DefaultCodespace, "insufficient account funds[abcdef01]; account has no coin").Error())
}

func TestKeeper_SetCoins(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	coins := types.Coins{types.Coin{Denom: defaultTokenIDFT, Amount: sdk.NewInt(-1)}}
	require.EqualError(t, keeper.SetCoins(ctx, defaultContractID, addr1, coins), types.ErrInvalidCoin(types.DefaultCodespace, coins.String()).Error())
}
