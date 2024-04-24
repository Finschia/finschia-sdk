package keeper_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/simapp"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
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

	accWithOldCoin sdk.AccAddress
	accWithNewCoin sdk.AccAddress
	initBalance    sdk.Int
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

	s.createAccountsWithInitBalance(app)
	app.AccountKeeper.GetModuleAccount(s.ctx, types.ModuleName)
}

func (s *KeeperTestSuite) createAccountsWithInitBalance(app *simapp.SimApp) {
	addresses := []*sdk.AccAddress{
		&s.accWithOldCoin,
		&s.accWithNewCoin,
	}
	for i, address := range s.createRandomAccounts(len(addresses)) {
		*addresses[i] = address
	}
	s.initBalance = sdk.NewInt(123456789)
	minter := app.AccountKeeper.GetModuleAccount(s.ctx, minttypes.ModuleName).GetAddress()
	oldAmount := sdk.NewCoins(sdk.NewCoin(s.keeper.OldDenom(), s.initBalance))
	s.Require().NoError(app.BankKeeper.MintCoins(s.ctx, minttypes.ModuleName, oldAmount))
	s.Require().NoError(app.BankKeeper.SendCoins(s.ctx, minter, s.accWithOldCoin, oldAmount))

	newAmount := sdk.NewCoins(sdk.NewCoin(s.keeper.NewDenom(), s.initBalance))
	s.Require().NoError(app.BankKeeper.MintCoins(s.ctx, minttypes.ModuleName, newAmount))
	s.Require().NoError(app.BankKeeper.SendCoins(s.ctx, minter, s.accWithNewCoin, newAmount))
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, &KeeperTestSuite{})
}

func (s *KeeperTestSuite) TestSwap() {
	testCases := map[string]struct {
		from                           sdk.AccAddress
		amountToSwap                   sdk.Int
		expectedBalanceWithoutMultiply sdk.Int
		shouldThrowError               bool
		expectedError                  error
	}{
		"swap some": {
			s.accWithOldCoin,
			sdk.NewInt(100),
			sdk.NewInt(100),
			false,
			nil,
		},
		"swap all the balance": {
			s.accWithOldCoin,
			s.initBalance,
			s.initBalance,
			false,
			nil,
		},
		"account holding new coin only": {
			s.accWithNewCoin,
			sdk.NewInt(100),
			s.initBalance,
			true,
			sdkerrors.ErrInsufficientFunds,
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			fmt.Printf("tc name=%s, %s\n", name, s.keeper.GetBalance(ctx, tc.from, s.keeper.OldDenom()))
			err := s.keeper.Swap(ctx, tc.from, sdk.NewCoin(s.keeper.OldDenom(), tc.amountToSwap))
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)
			actualAmount := s.keeper.GetBalance(ctx, tc.from, s.keeper.NewDenom()).Amount
			expectedAmount := tc.expectedBalanceWithoutMultiply.Mul(s.keeper.SwapMultiple())
			s.Require().Equal(expectedAmount, actualAmount)
		})
	}
}

func (s *KeeperTestSuite) TestSwapAll() {
	testCases := map[string]struct {
		from                           sdk.AccAddress
		expectedBalanceWithoutMultiply sdk.Int
		shouldThrowError               bool
		expectedError                  error
	}{
		"account holding old coin": {
			s.accWithOldCoin,
			s.initBalance,
			false,
			nil,
		},
		"account holding new coin only": {
			s.accWithNewCoin,
			s.initBalance,
			true,
			sdkerrors.ErrInsufficientFunds,
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			err := s.keeper.SwapAll(ctx, tc.from)
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)
			actualAmount := s.keeper.GetBalance(ctx, tc.from, s.keeper.NewDenom()).Amount
			expectedAmount := tc.expectedBalanceWithoutMultiply.Mul(s.keeper.SwapMultiple())
			s.Require().Equal(expectedAmount, actualAmount)
		})
	}
}
