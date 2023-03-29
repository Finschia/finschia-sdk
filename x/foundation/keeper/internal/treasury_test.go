package internal_test

import (
	sdk "github.com/line/lbm-sdk/types"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	"github.com/line/lbm-sdk/x/foundation"
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

			before := s.impl.GetTreasury(ctx)

			amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount))
			err := s.impl.FundTreasury(ctx, s.stranger, amount)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			expectedAfter := before.Add(sdk.NewDecCoinsFromCoins(amount...)...)
			poolAfter := s.impl.GetTreasury(ctx)
			s.Require().Equal(sdk.NewDecCoins(expectedAfter...), sdk.NewDecCoins(poolAfter...))
			balanceAfter := sdk.NewDecCoinsFromCoins(s.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(foundation.TreasuryName))...)
			s.Require().Equal(sdk.NewDecCoins(expectedAfter...), sdk.NewDecCoins(balanceAfter...))
		})
	}
}

func (s *KeeperTestSuite) TestWithdrawFromTreasury() {
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

			before := s.impl.GetTreasury(ctx)

			amount := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, tc.amount))
			err := s.impl.WithdrawFromTreasury(ctx, s.stranger, amount)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			expectedAfter := before.Sub(sdk.NewDecCoinsFromCoins(amount...))
			poolAfter := s.impl.GetTreasury(ctx)
			s.Require().Equal(sdk.NewDecCoins(expectedAfter...), sdk.NewDecCoins(poolAfter...))
			balanceAfter := sdk.NewDecCoinsFromCoins(s.bankKeeper.GetAllBalances(ctx, authtypes.NewModuleAddress(foundation.TreasuryName))...)
			s.Require().Equal(sdk.NewDecCoins(expectedAfter...), sdk.NewDecCoins(balanceAfter...))
		})
	}
}
