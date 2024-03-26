package keeper_test

import (
	"testing"
	"time"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	bankkeeper "github.com/Finschia/finschia-sdk/x/bank/keeper"
	banktestutil "github.com/Finschia/finschia-sdk/x/bank/testutil"
	distributionkeeper "github.com/Finschia/finschia-sdk/x/distribution/keeper"
	stakingkeeper "github.com/Finschia/finschia-sdk/x/staking/keeper"
	stakingtypes "github.com/Finschia/finschia-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestSlashRedelegation(t *testing.T) {
	// setting up
	app := simapp.Setup(false)

	stakingKeeper := app.StakingKeeper
	bankKeeper := app.BankKeeper
	slashKeeper := app.SlashingKeeper
	distrKeeper := app.DistrKeeper

	// get sdk context, staking msg server and bond denom
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Height: app.LastBlockHeight() + 1})
	stakingMsgServer := stakingkeeper.NewMsgServerImpl(stakingKeeper)
	bondDenom := stakingKeeper.BondDenom(ctx)

	// evilVal will be slashed, goodVal won't be slashed
	evilValPubKey := secp256k1.GenPrivKey().PubKey()
	goodValPubKey := secp256k1.GenPrivKey().PubKey()

	// both test acc 1 and 2 delegated to evil val, both acc should be slashed when evil val is slashed
	// test acc 1 use the "undelegation after redelegation" trick (redelegate to good val and then undelegate) to avoid slashing
	// test acc 2 only undelegate from evil val
	testAcc1 := sdk.AccAddress([]byte("addr1_______________"))
	testAcc2 := sdk.AccAddress([]byte("addr2_______________"))

	// fund acc 1 and acc 2
	testCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, stakingKeeper.TokensFromConsensusPower(ctx, 10)))
	err := banktestutil.FundAccount(bankKeeper, ctx, testAcc1, testCoins)
	require.NoError(t, err)
	err = banktestutil.FundAccount(bankKeeper, ctx, testAcc2, testCoins)
	require.NoError(t, err)

	balance1Before := bankKeeper.GetBalance(ctx, testAcc1, bondDenom)
	balance2Before := bankKeeper.GetBalance(ctx, testAcc2, bondDenom)

	// assert acc 1 and acc 2 balance
	require.Equal(t, balance1Before.Amount.String(), testCoins[0].Amount.String())
	require.Equal(t, balance2Before.Amount.String(), testCoins[0].Amount.String())

	// creating evil val
	evilValAddr := sdk.ValAddress(evilValPubKey.Address())
	err = banktestutil.FundAccount(bankKeeper, ctx, sdk.AccAddress(evilValAddr), testCoins)
	require.NoError(t, err)
	createValMsg1, _ := stakingtypes.NewMsgCreateValidator(
		evilValAddr, evilValPubKey, testCoins[0], stakingtypes.Description{Details: "test"}, stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0)), sdk.OneInt())
	_, err = stakingMsgServer.CreateValidator(sdk.WrapSDKContext(ctx), createValMsg1)
	require.NoError(t, err)

	// creating good val
	goodValAddr := sdk.ValAddress(goodValPubKey.Address())
	err = banktestutil.FundAccount(bankKeeper, ctx, sdk.AccAddress(goodValAddr), testCoins)
	require.NoError(t, err)
	createValMsg2, _ := stakingtypes.NewMsgCreateValidator(
		goodValAddr, goodValPubKey, testCoins[0], stakingtypes.Description{Details: "test"}, stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0)), sdk.OneInt())
	_, err = stakingMsgServer.CreateValidator(sdk.WrapSDKContext(ctx), createValMsg2)
	require.NoError(t, err)

	// next block, commit height 1, move to height 2
	// acc 1 and acc 2 delegate to evil val
	ctx = app.NextBlock(ctx, time.Duration(1))
	require.NoError(t, err)

	// Acc 2 delegate
	delMsg := stakingtypes.NewMsgDelegate(testAcc2, evilValAddr, testCoins[0])
	_, err = stakingMsgServer.Delegate(sdk.WrapSDKContext(ctx), delMsg)
	require.NoError(t, err)

	// Acc 1 delegate
	delMsg = stakingtypes.NewMsgDelegate(testAcc1, evilValAddr, testCoins[0])
	_, err = stakingMsgServer.Delegate(sdk.WrapSDKContext(ctx), delMsg)
	require.NoError(t, err)

	// next block, commit height 2, move to height 3
	// with the new delegations, evil val increases in voting power and commit byzantine behavior at height 3 consensus
	// at the same time, acc 1 and acc 2 withdraw delegation from evil val
	ctx = app.NextBlock(ctx, time.Duration(1))
	require.NoError(t, err)

	evilVal, found := stakingKeeper.GetValidator(ctx, evilValAddr)
	require.True(t, found)

	evilPower := stakingKeeper.TokensToConsensusPower(ctx, evilVal.Tokens)

	// Acc 1 redelegate from evil val to good val
	redelMsg := stakingtypes.NewMsgBeginRedelegate(testAcc1, evilValAddr, goodValAddr, testCoins[0])
	_, err = stakingMsgServer.BeginRedelegate(sdk.WrapSDKContext(ctx), redelMsg)
	require.NoError(t, err)

	// Acc 1 undelegate from good val
	undelMsg := stakingtypes.NewMsgUndelegate(testAcc1, goodValAddr, testCoins[0])
	_, err = stakingMsgServer.Undelegate(sdk.WrapSDKContext(ctx), undelMsg)
	require.NoError(t, err)

	// Acc 2 undelegate from evil val
	undelMsg = stakingtypes.NewMsgUndelegate(testAcc2, evilValAddr, testCoins[0])
	_, err = stakingMsgServer.Undelegate(sdk.WrapSDKContext(ctx), undelMsg)
	require.NoError(t, err)

	// next block, commit height 3, move to height 4
	// Slash evil val for byzantine behavior at height 3 consensus,
	// at which acc 1 and acc 2 still contributed to evil val voting power
	// even tho they undelegate at block 3, the valset update is applied after committed block 3 when height 3 consensus already passes
	ctx = app.NextBlock(ctx, time.Duration(1))
	require.NoError(t, err)

	// slash evil val with slash factor = 0.9, leaving only 10% of stake after slashing
	evilVal, _ = stakingKeeper.GetValidator(ctx, evilValAddr)
	evilValConsAddr, err := evilVal.GetConsAddr()
	require.NoError(t, err)

	slashKeeper.Slash(ctx, evilValConsAddr, sdk.MustNewDecFromStr("0.9"), evilPower, 3)

	// assert invariant to make sure we conduct slashing correctly
	_, stop := stakingkeeper.AllInvariants(stakingKeeper)(ctx)
	require.False(t, stop)

	_, stop = bankkeeper.AllInvariants(bankKeeper)(ctx)
	require.False(t, stop)

	_, stop = distributionkeeper.AllInvariants(distrKeeper)(ctx)
	require.False(t, stop)

	// one eternity later
	ctx = app.NextBlock(ctx, time.Duration(1000000000000000000))
	ctx = app.NextBlock(ctx, time.Duration(1))

	// confirm that account 1 and account 2 has been slashed, and the slash amount is correct
	balance1AfterSlashing := bankKeeper.GetBalance(ctx, testAcc1, bondDenom)
	balance2AfterSlashing := bankKeeper.GetBalance(ctx, testAcc2, bondDenom)

	require.Equal(t, balance1AfterSlashing.Amount.Mul(sdk.NewIntFromUint64(10)).String(), balance1Before.Amount.String())
	require.Equal(t, balance2AfterSlashing.Amount.Mul(sdk.NewIntFromUint64(10)).String(), balance2Before.Amount.String())
}
