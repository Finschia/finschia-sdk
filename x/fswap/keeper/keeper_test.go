package keeper_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/simapp"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	bank "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/fswap/keeper"
	"github.com/Finschia/finschia-sdk/x/fswap/testutil"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
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
	pebSwapRateForCony, err := sdk.NewDecFromStr("148079656000000")
	s.Require().NoError(err)
	swapCap := sdk.NewIntFromBigInt(pebSwapRateForCony.Mul(s.initBalance.ToDec()).BigInt())
	swapCap = swapCap.Mul(sdk.NewInt(numAcc))
	s.Require().NoError(err)
	s.swap = types.Swap{
		FromDenom:           "fromdenom",
		ToDenom:             "todenom",
		AmountCapForToDenom: swapCap,
		SwapRate:            pebSwapRateForCony,
	}
	s.toDenomMetadata = bank.Metadata{
		Description: "This is metadata for to-coin",
		DenomUnits: []*bank.DenomUnit{
			{Denom: s.swap.ToDenom, Exponent: 0},
		},
		Base:    "todenom",
		Display: "todenomcoin",
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

func TestNewKeeper(t *testing.T) {
	app := simapp.Setup(false)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authKeeper := testutil.NewMockAccountKeeper(ctrl)
	testCases := map[string]struct {
		malleate func()
		isPanic  bool
	}{
		"fswap module account has not been set": {
			malleate: func() {
				authKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(nil).Times(1)
				keeper.NewKeeper(app.AppCodec(), sdk.NewKVStoreKey(types.StoreKey), types.DefaultConfig(), types.DefaultAuthority().String(), authKeeper, app.BankKeeper)
			},
			isPanic: true,
		},
		"fswap authority must be the gov or foundation module account": {
			malleate: func() {
				authKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(authtypes.NewModuleAddress(types.ModuleName)).Times(1)
				keeper.NewKeeper(app.AppCodec(), sdk.NewKVStoreKey(types.StoreKey), types.DefaultConfig(), authtypes.NewModuleAddress("invalid").String(), authKeeper, app.BankKeeper)
			},
			isPanic: true,
		},
		"success - gov authority": {
			malleate: func() {
				authKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(authtypes.NewModuleAddress(types.ModuleName)).Times(1)
				keeper.NewKeeper(app.AppCodec(), sdk.NewKVStoreKey(types.StoreKey), types.DefaultConfig(), authtypes.NewModuleAddress(govtypes.ModuleName).String(), authKeeper, app.BankKeeper)
			},
			isPanic: false,
		},
		"success - foundation authority": {
			malleate: func() {
				authKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(authtypes.NewModuleAddress(types.ModuleName)).Times(1)
				keeper.NewKeeper(app.AppCodec(), sdk.NewKVStoreKey(types.StoreKey), types.DefaultConfig(), authtypes.NewModuleAddress(foundation.ModuleName).String(), authKeeper, app.BankKeeper)
			},
			isPanic: false,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.isPanic {
				require.Panics(t, tc.malleate)
			} else {
				tc.malleate()
			}
		})
	}
}

func (s *KeeperTestSuite) TestSwap() {
	swap2ExpectedAmount, ok := sdk.NewIntFromString("296159312000000")
	s.Require().True(ok)
	swap100ExpectedAmount, ok := sdk.NewIntFromString("14807965600000000")
	s.Require().True(ok)
	swapAllExpectedBalance, ok := sdk.NewIntFromString("18281438845984584000000")
	s.Require().True(ok)
	testCases := map[string]struct {
		from             sdk.AccAddress
		amountToSwap     sdk.Coin
		toDenom          string
		expectedAmount   sdk.Int
		shouldThrowError bool
		expectedError    error
	}{
		"swap 2 from-denom": {
			s.accWithFromCoin,
			sdk.NewCoin(s.swap.GetFromDenom(), sdk.NewInt(2)),
			s.swap.GetToDenom(),
			swap2ExpectedAmount,
			false,
			nil,
		},
		"swap some": {
			s.accWithFromCoin,
			sdk.NewCoin(s.swap.GetFromDenom(), sdk.NewInt(100)),
			s.swap.GetToDenom(),
			swap100ExpectedAmount,
			false,
			nil,
		},
		"swap all the balance": {
			s.accWithFromCoin,
			sdk.NewCoin(s.swap.GetFromDenom(), s.initBalance),
			s.swap.GetToDenom(),
			swapAllExpectedBalance,
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
			err := s.keeper.SetSwap(ctx, s.swap, s.toDenomMetadata)
			s.Require().NoError(err)

			err = s.keeper.Swap(ctx, tc.from, tc.amountToSwap, tc.toDenom)
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)

			actualAmount := s.keeper.GetBalance(ctx, tc.from, s.swap.GetToDenom()).Amount
			s.Require().Equal(tc.expectedAmount, actualAmount)
		})
	}
}

func (s *KeeperTestSuite) TestSetSwap() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	bankKeeper.EXPECT().HasSupply(gomock.Any(), "fromdenom").Return(true).AnyTimes()
	bankKeeper.EXPECT().HasSupply(gomock.Any(), gomock.Any()).Return(false).AnyTimes()
	bankKeeper.EXPECT().GetDenomMetaData(gomock.Any(), "todenom").Return(bank.Metadata{}, false).AnyTimes()
	bankKeeper.EXPECT().GetDenomMetaData(gomock.Any(), gomock.Any()).Return(s.toDenomMetadata, true).AnyTimes()
	bankKeeper.EXPECT().SetDenomMetaData(gomock.Any(), s.toDenomMetadata).Times(1)
	s.keeper.BankKeeper = bankKeeper

	testCases := map[string]struct {
		swap          types.Swap
		toDenomMeta   bank.Metadata
		expectedError error
	}{
		"valid": {
			types.Swap{
				FromDenom:           "fromdenom",
				ToDenom:             "todenom",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			s.toDenomMetadata,
			nil,
		},
		"from-denom does not exist": {
			types.Swap{
				FromDenom:           "fakedenom",
				ToDenom:             "todenom",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			s.toDenomMetadata,
			sdkerrors.ErrInvalidRequest,
		},
		"to-denom does not equal with metadata": {
			types.Swap{
				FromDenom:           "fromdenom",
				ToDenom:             "fakedenom",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			s.toDenomMetadata,
			sdkerrors.ErrInvalidRequest,
		},
		"to-denom metadata change not allowed": {
			types.Swap{
				FromDenom:           "fromdenom",
				ToDenom:             "change",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			bank.Metadata{
				Description: s.toDenomMetadata.Description,
				DenomUnits:  s.toDenomMetadata.DenomUnits,
				Base:        "change",
				Display:     s.toDenomMetadata.Display,
				Name:        s.toDenomMetadata.Name,
				Symbol:      s.toDenomMetadata.Symbol,
			},
			sdkerrors.ErrInvalidRequest,
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			err := s.keeper.SetSwap(ctx, tc.swap, s.toDenomMetadata)
			s.Require().ErrorIs(err, tc.expectedError)
		})
	}
}

func (s *KeeperTestSuite) TestSwapValidateBasic() {
	testCases := map[string]struct {
		swap             types.Swap
		shouldThrowError bool
		expectedError    error
	}{
		"valid": {
			types.Swap{
				FromDenom:           "fromD",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			false,
			nil,
		},
		"invalid empty from-denom": {
			types.Swap{
				FromDenom:           "",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
		"invalid empty to-denom": {
			types.Swap{
				FromDenom:           "fromD",
				ToDenom:             "",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
		"invalid zero amount cap for to-denom": {
			types.Swap{
				FromDenom:           "fromD",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.ZeroInt(),
				SwapRate:            sdk.OneDec(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
		"invalid zero swap-rate": {
			types.Swap{
				FromDenom:           "fromD",
				ToDenom:             "toD",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.ZeroDec(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
		"invalid the same from-denom and to-denom": {
			types.Swap{
				FromDenom:           "same",
				ToDenom:             "same",
				AmountCapForToDenom: sdk.OneInt(),
				SwapRate:            sdk.OneDec(),
			},
			true,
			sdkerrors.ErrInvalidRequest,
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			err := tc.swap.ValidateBasic()
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)
		})
	}
}
