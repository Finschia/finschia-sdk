package internal_test

import (
	sdk "github.com/line/lbm-sdk/types"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestCollectFoundationTax() {
	for name, tc := range map[string]struct {
		taxRatio sdk.Dec
		tax      sdk.Int
		valid    bool
	}{
		"common": {
			taxRatio: sdk.MustNewDecFromStr("0.123456789"),
			tax:      sdk.NewInt(121932631),
			valid:    true,
		},
		"zero ratio": {
			taxRatio: sdk.ZeroDec(),
			tax:      sdk.ZeroInt(),
			valid:    true,
		},
		"send fails": {
			taxRatio: sdk.MustNewDecFromStr("1.00000001"),
			tax:      sdk.NewInt(987654330),
		},
	} {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			s.impl.SetParams(ctx, foundation.Params{
				FoundationTax: tc.taxRatio,
			})

			before := s.impl.GetTreasury(ctx)
			s.Require().Equal(1, len(before))
			s.Require().Equal(sdk.NewDecFromInt(s.balance), before[0].Amount)

			tax := sdk.NewDecFromInt(s.balance).MulTruncate(tc.taxRatio).TruncateInt()
			// ensure the behavior does not change
			s.Require().Equal(tc.tax, tax)

			err := s.impl.CollectFoundationTax(ctx)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			expectedAfter := s.balance.Add(tax)
			after := s.impl.GetTreasury(ctx)
			s.Require().Equal(1, len(after))
			s.Require().Equal(sdk.NewDecFromInt(expectedAfter), after[0].Amount)
		})
	}
}

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
