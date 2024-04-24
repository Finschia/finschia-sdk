package keeper_test

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

func (s *KeeperTestSuite) TestQuerySwap() {
	testCases := map[string]struct {
		swapReq                        *types.MsgSwapRequest
		expectedBalanceWithoutMultiply sdk.Int
	}{
		"valid request": {
			swapReq: &types.MsgSwapRequest{
				FromAddress: s.accWithOldCoin.String(),
				Amount:      sdk.NewCoin(s.keeper.OldDenom(), sdk.NewInt(100)),
			},
			expectedBalanceWithoutMultiply: sdk.NewInt(100),
		},
	}
	for name, tc := range testCases {
		s.Run(name, func() {
			s.Require().NoError(tc.swapReq.ValidateBasic())
			ctx, _ := s.ctx.CacheContext()
			res, err := s.msgServer.Swap(sdk.WrapSDKContext(ctx), tc.swapReq)
			s.Require().NoError(err)
			s.Require().NotNil(res)
			swapped, err := s.queryServer.Swapped(sdk.WrapSDKContext(ctx), &types.QuerySwappedRequest{})
			s.Require().NoError(err)

			expectedOldAmount := tc.expectedBalanceWithoutMultiply
			expectedNewAmount := tc.expectedBalanceWithoutMultiply.Mul(s.keeper.SwapMultiple())
			actualOldCoinAmount := swapped.GetSwapped().OldCoinAmount
			actualNewCoinAmount := swapped.GetSwapped().NewCoinAmount
			s.Require().Equal(expectedOldAmount, actualOldCoinAmount)
			s.Require().Equal(expectedNewAmount, actualNewCoinAmount)
		})
	}
}

func (s *KeeperTestSuite) TestQueryTotalNewCurrencySwapLimit() {
	// TODO: Need to confirm, it may not necessary
	// Can be calculated by query Params.SwappableNewCoinAmount and query Swapped.NewCoinAmount

	// SwappableNewCoinAmount << as param and after first set it'll become constant value
	// Why user want to know constant?? use may want to remaining swappable balance.

	// s.Require().NoError(err)

	//s.Require().Equal(expectedLimit, res.SwappableNewCoinAmount)
}
