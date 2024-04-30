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

	swapInit     types.SwapInit
	fromDenom    string
	toDenom      string
	swapMultiple sdk.Int
	swapCap      sdk.Int
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

	s.fromDenom = "fromdenom"
	s.toDenom = "todenom"
	s.swapMultiple = sdk.NewInt(1000)
	s.initBalance = sdk.NewInt(123456789)
	numAcc := int64(2)
	s.swapCap = s.initBalance.Mul(s.swapMultiple.Mul(sdk.NewInt(numAcc)))
	s.swapInit = types.SwapInit{
		FromDenom:           s.fromDenom,
		ToDenom:             s.toDenom,
		AmountCapForToDenom: s.swapCap,
		SwapMultiple:        s.swapMultiple,
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
	fromAmount := sdk.NewCoins(sdk.NewCoin(s.fromDenom, s.initBalance))
	s.Require().NoError(app.BankKeeper.MintCoins(s.ctx, minttypes.ModuleName, fromAmount))
	s.Require().NoError(app.BankKeeper.SendCoins(s.ctx, minter, s.accWithFromCoin, fromAmount))

	toAmount := sdk.NewCoins(sdk.NewCoin(s.toDenom, s.initBalance))
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
			sdk.NewCoin(s.fromDenom, sdk.NewInt(100)),
			s.toDenom,
			sdk.NewInt(100),
			false,
			nil,
		},
		"swap all the balance": {
			s.accWithFromCoin,
			sdk.NewCoin(s.fromDenom, s.initBalance),
			s.toDenom,
			s.initBalance,
			false,
			nil,
		},
		"swap without holding enough balance": {
			s.accWithFromCoin,
			sdk.NewCoin(s.fromDenom, sdk.OneInt().Add(s.initBalance)),
			s.toDenom,
			sdk.ZeroInt(),
			true,
			sdkerrors.ErrInsufficientFunds,
		},
		"account holding new coin only": {
			s.accWithToCoin,
			sdk.NewCoin(s.fromDenom, sdk.NewInt(100)),
			s.toDenom,
			sdk.ZeroInt(),
			true,
			sdkerrors.ErrInsufficientFunds,
		},
		"swap with the same from-denom and to-denom": {
			s.accWithFromCoin,
			sdk.NewCoin(s.toDenom, s.initBalance),
			s.toDenom,
			sdk.ZeroInt(),
			true,
			sdkerrors.ErrInvalidRequest,
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			err := s.keeper.SwapInit(ctx, s.swapInit)
			s.Require().NoError(err)

			err = s.keeper.Swap(ctx, tc.from, tc.amountToSwap, tc.toDenom)
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)

			actualAmount := s.keeper.GetBalance(ctx, tc.from, s.toDenom).Amount
			expectedAmount := tc.expectedBalanceWithoutMultiply.Mul(s.swapMultiple)
			s.Require().Equal(expectedAmount, actualAmount)
		})
	}
}

func (s *KeeperTestSuite) TestSwapAll() {
	testCases := map[string]struct {
		from                           sdk.AccAddress
		fromDenom                      string
		toDenom                        string
		expectedBalanceWithoutMultiply sdk.Int
		shouldThrowError               bool
		expectedError                  error
	}{
		"account holding from coin": {
			s.accWithFromCoin,
			s.fromDenom,
			s.toDenom,
			s.initBalance,
			false,
			nil,
		},
		"account holding to coin only": {
			s.accWithToCoin,
			s.fromDenom,
			s.toDenom,
			s.initBalance,
			true,
			sdkerrors.ErrInsufficientFunds,
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			err := s.keeper.SwapInit(ctx, s.swapInit)
			s.Require().NoError(err)

			err = s.keeper.SwapAll(ctx, tc.from, tc.fromDenom, tc.toDenom)
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)

			actualAmount := s.keeper.GetBalance(ctx, tc.from, s.toDenom).Amount
			expectedAmount := tc.expectedBalanceWithoutMultiply.Mul(s.swapMultiple)
			s.Require().Equal(expectedAmount, actualAmount)
		})
	}
}

func (s *KeeperTestSuite) TestSwapInit() {
	testCases := map[string]struct {
		req              types.SwapInit
		shouldThrowError bool
		expectedError    error
	}{
		"valid swapInit": {
			types.SwapInit{
				FromDenom:           "fromD",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.OneInt(),
				SwapMultiple:        sdk.OneInt(),
			},
			false,
			nil,
		},
		"invalid empty from-denom": {
			types.SwapInit{
				FromDenom:           "",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.OneInt(),
				SwapMultiple:        sdk.OneInt(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
		"invalid empty to-denom": {
			types.SwapInit{
				FromDenom:           "fromD",
				ToDenom:             "",
				AmountCapForToDenom: sdk.OneInt(),
				SwapMultiple:        sdk.OneInt(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
		"invalid zero amount cap for to-denom": {
			types.SwapInit{
				FromDenom:           "fromD",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.ZeroInt(),
				SwapMultiple:        sdk.OneInt(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
		"invalid zero swap-rate": {
			types.SwapInit{
				FromDenom:           "fromD",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.OneInt(),
				SwapMultiple:        sdk.ZeroInt(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
		"invalid the same from-denom and to-denom": {
			types.SwapInit{
				FromDenom:           "same",
				ToDenom:             "same",
				AmountCapForToDenom: sdk.OneInt(),
				SwapMultiple:        sdk.OneInt(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			err := s.keeper.SwapInit(ctx, tc.req)
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)
		})
	}
}
