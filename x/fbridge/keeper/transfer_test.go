package keeper_test

import (
	"encoding/binary"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/keeper"
	"github.com/Finschia/finschia-sdk/x/fbridge/module"
	"github.com/Finschia/finschia-sdk/x/fbridge/testutil"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func TestHandleBridgeTransfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	key := storetypes.NewKVStoreKey(types.StoreKey)
	ctx := testutil.DefaultContextWithDB(t, key, sdk.NewTransientStoreKey("transient_test"))
	encCfg := testutil.MakeTestEncodingConfig(module.AppModuleBasic{})

	authKeeper := testutil.NewMockAccountKeeper(ctrl)
	bankKeeper := testutil.NewMockBankKeeper(ctrl)

	sender := sdk.AccAddress("test")
	amt := sdk.NewInt(1000000)
	denom := "stake"
	token := sdk.Coins{sdk.Coin{Denom: denom, Amount: amt}}

	bridge := authtypes.NewEmptyModuleAccount("bridge")
	authKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(bridge.GetAddress()).AnyTimes()
	authKeeper.EXPECT().GetModuleAccount(ctx, types.ModuleName).Return(bridge).AnyTimes()
	bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, token).Return(nil)
	bankKeeper.EXPECT().BurnCoins(ctx, types.ModuleName, token).Return(nil)

	k := keeper.NewKeeper(encCfg.Codec, key, authKeeper, bankKeeper, denom, "gov")
	beforeSeq := uint64(2)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, beforeSeq)
	ctx.KVStore(key).Set(types.KeyNextSeqSend, bz)

	afterSeq, err := k.HandleBridgeTransfer(ctx, sender, amt)
	require.NoError(t, err)
	require.Equal(t, beforeSeq+1, afterSeq)
}

func TestIsValidEthereumAddress(t *testing.T) {
	tcs := map[string]struct {
		valid   bool
		address string
	}{
		"valid": {
			valid:   true,
			address: "0xf7bAc63fc7CEaCf0589F25454Ecf5C2ce904997c",
		},
		"invalid": {
			valid:   false,
			address: "0xf7bAc63fc7CEaCf0589F25454Ecf5C2ce905997c",
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.valid, keeper.IsValidEthereumAddress(tc.address))
		})
	}
}
