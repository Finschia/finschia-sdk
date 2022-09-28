package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
)

func (s *KeeperTestSuite) TestFundTreasury() {
	testCases := map[string]struct {
		amount sdk.Int
		valid  bool
	}{
		"valid amount": {
			amount: s.balance,
			valid:  true,
		},
		"insufficient coins": {
			amount: s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount))
			err := s.keeper.FundTreasury(ctx, s.stranger, amount)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestWithDrawFromTreasury() {
	testCases := map[string]struct {
		amount sdk.Int
		valid  bool
	}{
		"valid amount": {
			amount: s.balance,
			valid:  true,
		},
		"insufficient coins": {
			amount: s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount))
			err := s.keeper.WithdrawFromTreasury(ctx, s.stranger, amount)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
