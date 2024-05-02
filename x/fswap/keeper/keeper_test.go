package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/simapp"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	bank "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/fswap/keeper"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
	minttypes "github.com/Finschia/finschia-sdk/x/mint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	goCtx       context.Context
	keeper      keeper.Keeper
	queryServer types.QueryServer
	msgServer   types.MsgServer

	accWithFromCoin sdk.AccAddress
	accWithToCoin   sdk.AccAddress
	initBalance     sdk.Int

	swap            types.Swap
	toDenomMetadata bank.Metadata
}

func (s *KeeperTestSuite) createRandomAccounts(n int) []sdk.AccAddress {
	seenAddresses := make(map[string]bool, n)
	addresses := make([]sdk.AccAddress, n)
	for i := range addresses {
		var address sdk.AccAddress
		for {
			pk := secp256k1.GenPrivKey().PubKey()
			address = sdk.AccAddress(pk.Address())
			if !seenAddresses[address.String()] {
				seenAddresses[address.String()] = true
				break
			}
		}
		addresses[i] = address
	}
	return addresses
}

func (s *KeeperTestSuite) SetupTest() {
	checkTx := false
	app := simapp.Setup(checkTx)
	testdata.RegisterInterfaces(app.InterfaceRegistry())
	testdata.RegisterMsgServer(app.MsgServiceRouter(), testdata.MsgServerImpl{})
	s.ctx = app.BaseApp.NewContext(checkTx, tmproto.Header{})
	s.goCtx = sdk.WrapSDKContext(s.ctx)
	s.keeper = app.FswapKeeper
	s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	numAcc := int64(2)
	s.initBalance = sdk.NewInt(123456789)
	swapRateForCony, err := sdk.NewDecFromStr("148079656000000")
	s.Require().NoError(err)
	swapCap := sdk.NewIntFromBigInt(swapRateForCony.Mul(s.initBalance.ToDec()).BigInt())
	swapCap = swapCap.Mul(sdk.NewInt(numAcc))
	s.Require().NoError(err)
	s.swap = types.Swap{
		FromDenom:           "fromdenom",
		ToDenom:             "todenom",
		AmountCapForToDenom: swapCap,
		SwapRate:            swapRateForCony,
	}
	s.toDenomMetadata = bank.Metadata{
		Description: "This is metadata for to-coin",
		DenomUnits: []*bank.DenomUnit{
			{Denom: s.swap.ToDenom, Exponent: 0},
		},
		Base:    "dummy",
		Display: "dummycoin",
		Name:    "DUMMY",
		Symbol:  "DUM",
	}
	s.createAccountsWithInitBalance(app)
	app.AccountKeeper.GetModuleAccount(s.ctx, types.ModuleName)
}

func (s *KeeperTestSuite) createAccountsWithInitBalance(app *simapp.SimApp) {
	addresses := []*sdk.AccAddress{
		&s.accWithFromCoin,
		&s.accWithToCoin,
	}
	for i, address := range s.createRandomAccounts(len(addresses)) {
		*addresses[i] = address
	}
	minter := app.AccountKeeper.GetModuleAccount(s.ctx, minttypes.ModuleName).GetAddress()
	fromAmount := sdk.NewCoins(sdk.NewCoin(s.swap.GetFromDenom(), s.initBalance))
	s.Require().NoError(app.BankKeeper.MintCoins(s.ctx, minttypes.ModuleName, fromAmount))
	s.Require().NoError(app.BankKeeper.SendCoins(s.ctx, minter, s.accWithFromCoin, fromAmount))

	toAmount := sdk.NewCoins(sdk.NewCoin(s.swap.GetToDenom(), s.initBalance))
	s.Require().NoError(app.BankKeeper.MintCoins(s.ctx, minttypes.ModuleName, toAmount))
	s.Require().NoError(app.BankKeeper.SendCoins(s.ctx, minter, s.accWithToCoin, toAmount))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, &KeeperTestSuite{})
}

func (s *KeeperTestSuite) TestSwap() {
	testCases := map[string]struct {
		from                           sdk.AccAddress
		amountToSwap                   sdk.Coin
		toDenom                        string
		expectedBalanceWithoutMultiply sdk.Int
		shouldThrowError               bool
		expectedError                  error
	}{
		"swap some": {
			s.accWithFromCoin,
			sdk.NewCoin(s.swap.GetFromDenom(), sdk.NewInt(100)),
			s.swap.GetToDenom(),
			sdk.NewInt(100),
			false,
			nil,
		},
		"swap all the balance": {
			s.accWithFromCoin,
			sdk.NewCoin(s.swap.GetFromDenom(), s.initBalance),
			s.swap.GetToDenom(),
			s.initBalance,
			false,
			nil,
		},
		"swap without holding enough balance": {
			s.accWithFromCoin,
			sdk.NewCoin(s.swap.GetFromDenom(), sdk.OneInt().Add(s.initBalance)),
			s.swap.GetToDenom(),
			sdk.ZeroInt(),
			true,
			sdkerrors.ErrInsufficientFunds,
		},
		"account holding new coin only": {
			s.accWithToCoin,
			sdk.NewCoin(s.swap.GetFromDenom(), sdk.NewInt(100)),
			s.swap.GetToDenom(),
			sdk.ZeroInt(),
			true,
			sdkerrors.ErrInsufficientFunds,
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			err := s.keeper.MakeSwap(ctx, s.swap, s.toDenomMetadata)
			s.Require().NoError(err)

			err = s.keeper.Swap(ctx, tc.from, tc.amountToSwap, tc.toDenom)
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)

			actualAmount := s.keeper.GetBalance(ctx, tc.from, s.swap.GetToDenom()).Amount
			multipliedAmountDec := s.swap.SwapRate.Mul(sdk.NewDecFromBigInt(tc.expectedBalanceWithoutMultiply.BigInt()))
			expectedAmount := sdk.NewIntFromBigInt(multipliedAmountDec.BigInt())
			s.Require().Equal(expectedAmount, actualAmount)
		})
	}
}
