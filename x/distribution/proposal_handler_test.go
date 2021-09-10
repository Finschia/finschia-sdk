package distribution_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/crypto/keys/ed25519"
	"github.com/line/lfb-sdk/simapp"
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/distribution"
	"github.com/line/lfb-sdk/x/distribution/types"
)

var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.BytesToAccAddress(delPk1.Address())

	amount = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1)))
)

func testProposal(recipient sdk.AccAddress, amount sdk.Coins) *types.CommunityPoolSpendProposal {
	return types.NewCommunityPoolSpendProposal("Test", "description", recipient, amount)
}

func TestProposalHandlerPassed(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	recipient := delAddr1

	// add coins to the module account
	macc := app.DistrKeeper.GetDistributionAccount(ctx)
	balances := app.BankKeeper.GetAllBalances(ctx, macc.GetAddress())
	err := app.BankKeeper.SetBalances(ctx, macc.GetAddress(), balances.Add(amount...))
	require.NoError(t, err)

	app.AccountKeeper.SetModuleAccount(ctx, macc)

	account := app.AccountKeeper.NewAccountWithAddress(ctx, recipient)
	app.AccountKeeper.SetAccount(ctx, account)
	require.True(t, app.BankKeeper.GetAllBalances(ctx, account.GetAddress()).IsZero())

	feePool := app.DistrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = sdk.NewDecCoinsFromCoins(amount...)
	app.DistrKeeper.SetFeePool(ctx, feePool)

	tp := testProposal(recipient, amount)
	hdlr := distribution.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)
	require.NoError(t, hdlr(ctx, tp))

	balances = app.BankKeeper.GetAllBalances(ctx, recipient)
	require.Equal(t, balances, amount)
}

func TestProposalHandlerFailed(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	recipient := delAddr1

	account := app.AccountKeeper.NewAccountWithAddress(ctx, recipient)
	app.AccountKeeper.SetAccount(ctx, account)
	require.True(t, app.BankKeeper.GetAllBalances(ctx, account.GetAddress()).IsZero())

	tp := testProposal(recipient, amount)
	hdlr := distribution.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)
	require.Error(t, hdlr(ctx, tp))

	balances := app.BankKeeper.GetAllBalances(ctx, recipient)
	require.True(t, balances.IsZero())
}
