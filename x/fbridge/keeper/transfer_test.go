package keeper

import (
	"encoding/binary"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/testutil"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func TestHandleBridgeTransfer(t *testing.T) {
	key, memKey, ctx, encCfg, authKeeper, bankKeeper, _ := testutil.PrepareFbridgeTest(t, 0)

	sender := sdk.AccAddress("test")
	amt := sdk.NewInt(1000000)
	denom := "stake"
	token := sdk.Coins{sdk.Coin{Denom: denom, Amount: amt}}
	bankKeeper.EXPECT().IsSendEnabledCoins(ctx, token).Return(nil).Times(1)
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, token).Return(nil).Times(1)
	bankKeeper.EXPECT().BurnCoins(ctx, types.ModuleName, token).Return(nil).Times(1)

	k := NewKeeper(encCfg.Codec, key, memKey, authKeeper, bankKeeper, types.DefaultAuthority().String())
	params := types.DefaultParams()
	params.TargetDenom = denom
	err := k.SetParams(ctx, params)
	require.NoError(t, err)
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

	// test error cases
	bankKeeper.EXPECT().IsSendEnabledCoins(ctx, token).Return(nil).Times(1)
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, token).Return(errors.New("insufficient funds")).Times(1)
	handledSeq, err = k.handleBridgeTransfer(ctx, sender, amt)
	require.Error(t, err)

	bankKeeper.EXPECT().IsSendEnabledCoins(ctx, token).Return(nil).Times(1)
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, token).Return(nil).Times(1)
	bankKeeper.EXPECT().BurnCoins(ctx, types.ModuleName, token).Return(errors.New("failed to burn coins")).Times(1)
	require.Panics(t, func() { _, _ = k.handleBridgeTransfer(ctx, sender, amt) }, "cannot burn coins after a successful send to a module account: failed to burn coins")
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
