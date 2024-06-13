package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/simapp"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	bank "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
	fkeeper "github.com/Finschia/finschia-sdk/x/fswap/keeper"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *KeeperTestSuite) TestMsgSwap() {
	swap2ExpectedAmount, ok := sdk.NewIntFromString("296159312000000")
	s.Require().True(ok)
	swap100ExpectedAmount, ok := sdk.NewIntFromString("14807965600000000")
	s.Require().True(ok)
	swapAllExpectedBalance, ok := sdk.NewIntFromString("18281438845984584000000")
	s.Require().True(ok)
	testCases := map[string]struct {
		request          *types.MsgSwap
		expectedAmount   sdk.Int
		shouldThrowError bool
		expectedError    error
	}{
		"swap 2 from-denom": {
			&types.MsgSwap{
				FromAddress:    s.accWithFromCoin.String(),
				FromCoinAmount: sdk.NewCoin(s.swap.GetFromDenom(), sdk.NewInt(2)),
				ToDenom:        s.swap.GetToDenom(),
			},
			swap2ExpectedAmount,
			false,
			nil,
		},
		"swap some": {
			&types.MsgSwap{
				FromAddress:    s.accWithFromCoin.String(),
				FromCoinAmount: sdk.NewCoin(s.swap.GetFromDenom(), sdk.NewInt(100)),
				ToDenom:        s.swap.GetToDenom(),
			},
			swap100ExpectedAmount,
			false,
			nil,
		},
		"swap all the balance": {
			&types.MsgSwap{
				FromAddress:    s.accWithFromCoin.String(),
				FromCoinAmount: sdk.NewCoin(s.swap.GetFromDenom(), s.initBalance),
				ToDenom:        s.swap.GetToDenom(),
			},
			swapAllExpectedBalance,
			false,
			nil,
		},
		"account holding to-coin only": {
			&types.MsgSwap{
				FromAddress:    s.accWithToCoin.String(),
				FromCoinAmount: sdk.NewCoin(s.swap.GetFromDenom(), sdk.NewInt(100)),
				ToDenom:        s.swap.GetToDenom(),
			},
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

			swapResponse, err := s.msgServer.Swap(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NotNil(swapResponse)
			s.Require().NoError(err)

			from, err := sdk.AccAddressFromBech32(tc.request.FromAddress)
			s.Require().NoError(err)
			actualAmount := s.keeper.GetBalance(ctx, from, tc.request.GetToDenom()).Amount
			s.Require().Equal(tc.expectedAmount, actualAmount)
		})
	}
}

func (s *KeeperTestSuite) TestMsgSwapAll() {
	swapAllExpectedBalance, ok := sdk.NewIntFromString("18281438845984584000000")
	s.Require().True(ok)
	testCases := map[string]struct {
		request          *types.MsgSwapAll
		expectedAmount   sdk.Int
		shouldThrowError bool
		expectedError    error
	}{
		"swapAll": {
			&types.MsgSwapAll{
				FromAddress: s.accWithFromCoin.String(),
				FromDenom:   s.swap.GetFromDenom(),
				ToDenom:     s.swap.GetToDenom(),
			},
			swapAllExpectedBalance,
			false,
			nil,
		},
		"account holding to-coin only": {
			&types.MsgSwapAll{
				FromAddress: s.accWithToCoin.String(),
				FromDenom:   s.swap.GetFromDenom(),
				ToDenom:     s.swap.GetToDenom(),
			},
			s.initBalance,
			true,
			sdkerrors.ErrInsufficientFunds,
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			err := s.keeper.SetSwap(ctx, s.swap, s.toDenomMetadata)
			s.Require().NoError(err)

			swapResponse, err := s.msgServer.SwapAll(sdk.WrapSDKContext(ctx), tc.request)
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NotNil(swapResponse)
			s.Require().NoError(err)

			from, err := sdk.AccAddressFromBech32(tc.request.FromAddress)
			s.Require().NoError(err)
			actualAmount := s.keeper.GetBalance(ctx, from, tc.request.GetToDenom()).Amount
			s.Require().Equal(tc.expectedAmount, actualAmount)
		})
	}
}

func TestMsgSetSwap(t *testing.T) {
	authority := authtypes.NewModuleAddress(foundation.ModuleName)
	checkTx := false
	app := simapp.Setup(checkTx)
	testdata.RegisterInterfaces(app.InterfaceRegistry())
	testdata.RegisterMsgServer(app.MsgServiceRouter(), testdata.MsgServerImpl{})
	ctx := app.BaseApp.NewContext(checkTx, tmproto.Header{})
	keeper := app.FswapKeeper
	msgServer := fkeeper.NewMsgServer(keeper)
	fromDenomStr := "cony"
	fromDenom := bank.Metadata{
		Description: "This is metadata for from-coin",
		DenomUnits: []*bank.DenomUnit{
			{Denom: fromDenomStr, Exponent: 0},
		},
		Base:    fromDenomStr,
		Display: fromDenomStr,
		Name:    "FROM",
		Symbol:  "FROM",
	}
	app.BankKeeper.SetDenomMetaData(ctx, fromDenom)

	testCases := map[string]struct {
		request       *types.MsgSetSwap
		expectedError error
	}{
		"valid": {
			request: &types.MsgSetSwap{
				Authority: authority.String(),
				Swap: types.Swap{
					FromDenom:           fromDenomStr,
					ToDenom:             "kei",
					AmountCapForToDenom: sdk.OneInt(),
					SwapRate:            sdk.NewDec(1),
				},
				ToDenomMetadata: bank.Metadata{
					Description: "desc",
					DenomUnits: []*bank.DenomUnit{
						{
							Denom:    "kei",
							Exponent: 0,
							Aliases:  nil,
						},
					},
					Base:    "kei",
					Display: "kei",
					Name:    "kei",
					Symbol:  "KAIA",
				},
			},
			expectedError: nil,
		},
		"invalid: authority": {
			request: &types.MsgSetSwap{
				Authority: "invalid-authority",
				Swap: types.Swap{
					FromDenom:           fromDenomStr,
					ToDenom:             "kei",
					AmountCapForToDenom: sdk.OneInt(),
					SwapRate:            sdk.NewDec(1),
				},
				ToDenomMetadata: bank.Metadata{
					Description: "desc",
					DenomUnits: []*bank.DenomUnit{
						{
							Denom:    "kei",
							Exponent: 0,
							Aliases:  nil,
						},
					},
					Base:    "kei",
					Display: "kei",
					Name:    "kei",
					Symbol:  "KAIA",
				},
			},
			expectedError: sdkerrors.ErrUnauthorized,
		},
		"invalid: Swap.ToDenom": {
			request: &types.MsgSetSwap{
				Authority: authority.String(),
				Swap: types.Swap{
					FromDenom:           fromDenomStr,
					ToDenom:             fromDenomStr,
					AmountCapForToDenom: sdk.OneInt(),
					SwapRate:            sdk.NewDec(1),
				},
				ToDenomMetadata: bank.Metadata{
					Description: "desc",
					DenomUnits: []*bank.DenomUnit{
						{
							Denom:    "kei",
							Exponent: 0,
							Aliases:  nil,
						},
					},
					Base:    "kei",
					Display: "kei",
					Name:    "kei",
					Symbol:  "KAIA",
				},
			},
			expectedError: sdkerrors.ErrInvalidRequest,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			context, _ := ctx.CacheContext()
			_, err := msgServer.SetSwap(sdk.WrapSDKContext(context), tc.request)
			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}
