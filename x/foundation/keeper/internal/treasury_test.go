package internal_test

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestCollectFoundationTax() {
	ctx, _ := s.ctx.CacheContext()

	// empty fee collector first
	// send the fee to the stranger
	// and get it back later if the test case requires
	collector := authtypes.NewModuleAddress(authtypes.FeeCollectorName)
	fees := s.bankKeeper.GetAllBalances(ctx, collector)
	err := s.bankKeeper.SendCoinsFromModuleToAccount(ctx, authtypes.FeeCollectorName, s.stranger, fees)
	s.Require().NoError(err)

	for name, tc := range map[string]struct {
		fee      math.Int
		taxRatio math.LegacyDec
		tax      math.Int
		valid    bool
	}{
		"common": {
			fee:      fees[0].Amount,
			taxRatio: math.LegacyMustNewDecFromStr("0.123456789"),
			tax:      math.NewInt(121932631),
			valid:    true,
		},
		"zero fee": {
			fee:      math.ZeroInt(),
			taxRatio: math.LegacyMustNewDecFromStr("0.123456789"),
			tax:      math.ZeroInt(),
			valid:    true,
		},
		"zero ratio": {
			fee:      fees[0].Amount,
			taxRatio: math.LegacyZeroDec(),
			tax:      math.ZeroInt(),
			valid:    true,
		},
	} {
		s.Run(name, func() {
			ctx, _ := ctx.CacheContext()

			// set fee
			err := s.bankKeeper.SendCoinsFromAccountToModule(ctx, s.stranger, authtypes.FeeCollectorName, sdk.NewCoins(sdk.NewCoin(fees[0].Denom, tc.fee)))
			s.Require().NoError(err)

			// set tax ratio
			s.impl.SetParams(ctx, foundation.Params{
				FoundationTax: tc.taxRatio,
			})

			before := s.impl.GetTreasury(ctx)
			s.Require().Equal(1, len(before))
			s.Require().Equal(math.LegacyNewDecFromInt(s.balance), before[0].Amount)

			tax := math.LegacyNewDecFromInt(tc.fee).MulTruncate(tc.taxRatio).TruncateInt()
			// ensure the behavior does not change
			s.Require().Equal(tc.tax, tax)

			err = s.impl.CollectFoundationTax(ctx)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)

			expectedAfter := s.balance.Add(tax)
			after := s.impl.GetTreasury(ctx)
			s.Require().Equal(1, len(after))
			s.Require().Equal(math.LegacyNewDecFromInt(expectedAfter), after[0].Amount)
		})
	}
}

func (s *KeeperTestSuite) TestFundTreasury() {
	testCases := map[string]struct {
		amount math.Int
		valid  bool
	}{
		"valid amount": {
			amount: s.balance,
			valid:  true,
		},
		"insufficient coins": {
			amount: s.balance.Add(math.OneInt()),
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
		amount math.Int
		valid  bool
	}{
		"valid amount": {
			amount: s.balance,
			valid:  true,
		},
		"insufficient coins": {
			amount: s.balance.Add(math.OneInt()),
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
