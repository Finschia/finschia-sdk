package ante_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestAccountLock(t *testing.T) {
	ctx, tx := prepare(t)

	ald := ante.NewAccountLockDecorator()
	antehandler := sdk.ChainAnteDecorators(ald)

	// no errors
	_, err := antehandler(ctx, tx, false)
	require.Nil(t, err)

	// no errors as well when second try
	_, err = antehandler(ctx, tx, false)
	require.Nil(t, err)
}

func TestAccountLockWithPanic(t *testing.T) {
	ctx, tx := prepare(t)

	ald := ante.NewAccountLockDecorator()
	panicantehandler := sdk.ChainAnteDecorators(ald, PanicDecorator{})

	require.Panics(t, func() { panicantehandler(ctx, tx, false) })

	antehandler := sdk.ChainAnteDecorators(ald)
	// no errors even though after panic
	_, err := antehandler(ctx, tx, false)
	require.Nil(t, err)
}

func prepare(t *testing.T) (sdk.Context, sdk.Tx) {
	// setup
	app, ctx := createTestApp(true)

	// keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()
	priv2, _, addr2 := types.KeyTestPubAddr()
	priv3, _, addr3 := types.KeyTestPubAddr()

	addrs := []sdk.AccAddress{addr1, addr2, addr3}
	msgs := make([]sdk.Msg, len(addrs))
	// set accounts and create msg for each address
	for i, addr := range addrs {
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
		require.NoError(t, acc.SetAccountNumber(uint64(i)))
		app.AccountKeeper.SetAccount(ctx, acc)
		msgs[i] = types.NewTestMsg(addr)
	}

	fee := types.NewTestStdFee()

	privs, accNums, seqs := []crypto.PrivKey{priv1, priv2, priv3}, []uint64{0, 1, 2}, []uint64{0, 0, 0}
	tx := types.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)

	return ctx, tx
}
