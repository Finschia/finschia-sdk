package keeper_test

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *KeeperTestSuite) TestMsgSwap() {
	testCases := map[string]struct {
		request                        *types.MsgSwap
		expectedBalanceWithoutMultiply sdk.Int
		shouldThrowError               bool
		expectedError                  error
	}{
		"swap some": {
			&types.MsgSwap{
				FromAddress:    s.accWithFromCoin.String(),
				FromCoinAmount: sdk.NewCoin(s.swapInit.GetFromDenom(), sdk.NewInt(100)),
				ToDenom:        s.swapInit.GetToDenom(),
			},
			sdk.NewInt(100),
			false,
			nil,
		},
		"swap all the balance": {
			&types.MsgSwap{
				FromAddress:    s.accWithFromCoin.String(),
				FromCoinAmount: sdk.NewCoin(s.swapInit.GetFromDenom(), s.initBalance),
				ToDenom:        s.swapInit.GetToDenom(),
			},
			s.initBalance,
			false,
			nil,
		},
		"account holding new coin only": {
			&types.MsgSwap{
				FromAddress:    s.accWithToCoin.String(),
				FromCoinAmount: sdk.NewCoin(s.swapInit.GetFromDenom(), sdk.NewInt(100)),
				ToDenom:        s.swapInit.GetToDenom(),
			},
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
			expectedAmount := tc.expectedBalanceWithoutMultiply.Mul(s.swapInit.SwapMultiple)
			s.Require().Equal(expectedAmount, actualAmount)
		})
	}
}

func (s *KeeperTestSuite) TestMsgSwapAll() {
	testCases := map[string]struct {
		request                        *types.MsgSwapAll
		expectedBalanceWithoutMultiply sdk.Int
		shouldThrowError               bool
		expectedError                  error
	}{
		"swapAll": {
			&types.MsgSwapAll{
				FromAddress: s.accWithFromCoin.String(),
				FromDenom:   s.swapInit.GetFromDenom(),
				ToDenom:     s.swapInit.GetToDenom(),
			},
			s.initBalance,
			false,
			nil,
		},
		"account holding new coin only": {
			&types.MsgSwapAll{
				FromAddress: s.accWithToCoin.String(),
				FromDenom:   s.swapInit.GetFromDenom(),
				ToDenom:     s.swapInit.GetToDenom(),
			},
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
			expectedAmount := tc.expectedBalanceWithoutMultiply.Mul(s.swapInit.SwapMultiple)
			s.Require().Equal(expectedAmount, actualAmount)
		})
	}
}
