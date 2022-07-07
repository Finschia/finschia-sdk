package bank_test

import (
	"strings"
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/testutil/testdata"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	"github.com/line/lbm-sdk/x/bank"
	bankkeeper "github.com/line/lbm-sdk/x/bank/keeper"
	"github.com/line/lbm-sdk/x/bank/types"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

func TestInvalidMsg(t *testing.T) {
	h := bank.NewHandler(nil)

	res, err := h(sdk.NewContext(nil, ocproto.Header{}, false, nil), testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)

	_, _, log := sdkerrors.ABCIInfo(err, false)
	require.True(t, strings.Contains(log, "unrecognized bank message type"))
}

// A module account cannot be the recipient of bank sends unless it has been marked as such
func TestSendToModuleAccount(t *testing.T) {
	priv1 := secp256k1.GenPrivKey()
	addr1 := sdk.AccAddress(priv1.PubKey().Address())
	moduleAccAddr := authtypes.NewModuleAddress(stakingtypes.BondedPoolName)
	coins := sdk.Coins{sdk.NewInt64Coin("foocoin", 10)}

	tests := []struct {
		name          string
		expectedError error
		msg           *types.MsgSend
	}{
		{
			name:          "not allowed module account",
			msg:           types.NewMsgSend(addr1, moduleAccAddr, coins),
			expectedError: sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", moduleAccAddr),
		},
		{
			name:          "allowed module account",
			msg:           types.NewMsgSend(addr1, authtypes.NewModuleAddress(stakingtypes.ModuleName), coins),
			expectedError: nil,
		},
	}

	acc1 := &authtypes.BaseAccount{
		Address: addr1.String(),
	}
	accs := authtypes.GenesisAccounts{acc1}
	balances := []types.Balance{
		{
			Address: addr1.String(),
			Coins:   coins,
		},
	}

	app := simapp.SetupWithGenesisAccounts(accs, balances...)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	app.BankKeeper = bankkeeper.NewBaseKeeper(
		app.AppCodec(), app.GetKey(types.StoreKey), app.AccountKeeper, app.GetSubspace(types.ModuleName), map[string]bool{
			moduleAccAddr.String(): true,
		},
	)
	handler := bank.NewHandler(app.BankKeeper)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := handler(ctx, tc.msg)
			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
