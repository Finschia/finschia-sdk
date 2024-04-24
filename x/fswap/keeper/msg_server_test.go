package keeper_test

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Finschia/finschia-sdk/testutil"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *KeeperTestSuite) TestMsgSwap() {
	testCases := map[string]struct {
		req                            *types.MsgSwapRequest
		expectedBalanceWithoutMultiply sdk.Int
		shouldThrowError               bool
		expectedError                  error
		expectedEvents                 sdk.Events
	}{
		"valid request": {
			req: &types.MsgSwapRequest{
				FromAddress: s.accWithOldCoin.String(),
				Amount:      sdk.NewCoin(s.keeper.OldDenom(), sdk.NewInt(100)),
			},
			expectedBalanceWithoutMultiply: sdk.NewInt(100),
			expectedEvents: sdk.Events{
				sdk.Event{
					Type: "lbm.fswap.v1.EventSwapCoins",
					Attributes: []abci.EventAttribute{
						{Key: []byte("address"), Value: testutil.W(s.accWithOldCoin.String()), Index: false},
						{Key: []byte("new_coin_amount"), Value: testutil.W(sdk.NewInt(100).Mul(s.keeper.SwapMultiple()).String()), Index: false},
						{Key: []byte("old_coin_amount"), Value: testutil.W(sdk.NewInt(100).String()), Index: false},
					},
				},
			},
		},
		"invalid request(try with swap with newCoin)": {
			req: &types.MsgSwapRequest{
				FromAddress: s.accWithOldCoin.String(),
				Amount:      sdk.NewCoin(s.keeper.NewDenom(), sdk.NewInt(100)),
			},
			expectedBalanceWithoutMultiply: sdk.NewInt(100),
			shouldThrowError:               true,
			expectedError:                  sdkerrors.ErrInvalidCoins,
			expectedEvents:                 sdk.Events{},
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			s.Require().NoError(tc.req.ValidateBasic())
			from, err := sdk.AccAddressFromBech32(tc.req.FromAddress)
			s.Require().NoError(err)
			ctx, _ := s.ctx.CacheContext()

			res, err := s.msgServer.Swap(sdk.WrapSDKContext(ctx), tc.req)
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)

			expectedAmount := tc.expectedBalanceWithoutMultiply.Mul(s.keeper.SwapMultiple())
			actualAmount := s.keeper.GetBalance(ctx, from, s.keeper.NewDenom()).Amount
			s.Require().Equal(expectedAmount, actualAmount)

			events := ctx.EventManager().Events()
			s.Require().Contains(events, tc.expectedEvents[0])
		})
	}
}

func (s *KeeperTestSuite) TestMsgSwapAll() {
	testCases := map[string]struct {
		req                            *types.MsgSwapAllRequest
		expectedBalanceWithoutMultiply sdk.Int
		shouldThrowError               bool
		expectedError                  error
		expectedEvents                 sdk.Events
	}{
		"valid request": {
			req: &types.MsgSwapAllRequest{
				FromAddress: s.accWithOldCoin.String(),
			},
			expectedBalanceWithoutMultiply: s.initBalance,
			expectedEvents: sdk.Events{
				sdk.Event{
					Type: "lbm.fswap.v1.EventSwapCoins",
					Attributes: []abci.EventAttribute{
						{Key: []byte("address"), Value: testutil.W(s.accWithOldCoin.String()), Index: false},
						{Key: []byte("new_coin_amount"), Value: testutil.W(s.initBalance.Mul(s.keeper.SwapMultiple()).String()), Index: false},
						{Key: []byte("old_coin_amount"), Value: testutil.W(s.initBalance.String()), Index: false},
					},
				},
			},
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			s.Require().NoError(tc.req.ValidateBasic())
			from, err := sdk.AccAddressFromBech32(tc.req.FromAddress)
			s.Require().NoError(err)
			ctx, _ := s.ctx.CacheContext()

			res, err := s.msgServer.SwapAll(sdk.WrapSDKContext(ctx), tc.req)
			if tc.shouldThrowError {
				s.Require().ErrorIs(err, tc.expectedError)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)

			expectedAmount := tc.expectedBalanceWithoutMultiply.Mul(s.keeper.SwapMultiple())
			actualAmount := s.keeper.GetBalance(ctx, from, s.keeper.NewDenom()).Amount
			s.Require().Equal(expectedAmount, actualAmount)

			events := ctx.EventManager().Events()
			s.Require().Contains(events, tc.expectedEvents[0])
		})
	}
}
