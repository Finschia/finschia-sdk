package keeper

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/testutil"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func TestHandleBridgeTransfer(t *testing.T) {
	key, memKey, ctx, encCfg, authKeeper, bankKeeper := testutil.PrepareFbridgeTest(t)

	sender := sdk.AccAddress("test")
	amt := sdk.NewInt(1000000)
	denom := "stake"
	token := sdk.Coins{sdk.Coin{Denom: denom, Amount: amt}}

	bridge := authtypes.NewEmptyModuleAccount("fbridge")
	authKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(bridge.GetAddress()).AnyTimes()
	bankKeeper.EXPECT().IsSendEnabledCoins(ctx, token).Return(nil)
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, token).Return(nil)
	bankKeeper.EXPECT().BurnCoins(ctx, types.ModuleName, token).Return(nil)

	k := NewKeeper(encCfg.Codec, key, memKey, authKeeper, bankKeeper, denom, types.DefaultAuthority().String())
	targetSeq := uint64(2)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, targetSeq)
	ctx.KVStore(key).Set(types.KeyNextSeqSend, bz)

	handledSeq, err := k.handleBridgeTransfer(ctx, sender, amt)
	require.NoError(t, err)
	require.Equal(t, targetSeq, handledSeq)
	afterSeq := k.GetNextSequence(ctx)
	require.Equal(t, targetSeq+1, afterSeq)
	h, err := k.GetSeqToBlocknum(ctx, handledSeq)
	require.NoError(t, err)
	require.Equal(t, uint64(ctx.BlockHeight()), h)
}

func TestIsValidEthereumAddress(t *testing.T) {
	tcs := map[string]struct {
		isErr   bool
		address string
	}{
		"valid": {
			isErr:   true,
			address: "0xf7bAc63fc7CEaCf0589F25454Ecf5C2ce904997c",
		},
		"invalid": {
			isErr:   false,
			address: "0xf7bAc63fc7CEaCf0589F25454Ecf5C2ce905997c",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.isErr, IsValidEthereumAddress(tc.address) == nil)
		})
	}
}
