package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_SendCoins(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	require.EqualError(t, keeper.SendCoins(ctx, addr3, addr1, types.NewCoins(types.NewCoin(defaultTokenIDFT, sdk.NewInt(1)))), sdkerrors.Wrap(types.ErrInsufficientToken, "insufficient account funds[abcdef01]; account has no coin").Error())
}

func TestKeeper_SetCoins(t *testing.T) {
	ctx := cacheKeeper()
	prepareCollectionTokens(ctx, t)

	coins := types.Coins{types.Coin{Denom: defaultTokenIDFT, Amount: sdk.NewInt(-1)}}
	require.EqualError(t, keeper.SetCoins(ctx, addr1, coins), sdkerrors.Wrapf(types.ErrInvalidCoin, "invalid amount: %s", coins.String()).Error())
}
