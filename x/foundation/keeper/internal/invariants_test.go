package internal_test

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/keeper/internal"
)

func (s *KeeperTestSuite) TestModuleAccountInvariant() {
	testCases := map[string]struct {
		malleate func(ctx sdk.Context)
		valid    bool
	}{
		"invariant not broken": {
			valid: true,
		},
		"treasury differs from the balance": {
			malleate: func(ctx sdk.Context) {
				balance := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, s.balance.Add(sdk.OneInt())))
				s.impl.SetPool(ctx, foundation.Pool{
					Treasury: sdk.NewDecCoinsFromCoins(balance...),
				})
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			invariant := internal.ModuleAccountInvariant(s.impl)
			_, broken := invariant(ctx)
			s.Require().Equal(!tc.valid, broken)
		})
	}
}

func (s *KeeperTestSuite) TestTotalWeightInvariant() {
	testCases := map[string]struct {
		malleate func(ctx sdk.Context)
		valid    bool
	}{
		"invariant not broken": {
			valid: true,
		},
		"total weight differs from the number of foundation members": {
			malleate: func(ctx sdk.Context) {
				info := s.impl.GetFoundationInfo(ctx)
				numMembers := len(s.impl.GetMembers(ctx))
				info.TotalWeight = sdk.NewDec(int64(numMembers)).Add(sdk.OneDec())
				s.impl.SetFoundationInfo(ctx, info)
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			invariant := internal.TotalWeightInvariant(s.impl)
			_, broken := invariant(ctx)
			s.Require().Equal(!tc.valid, broken)
		})
	}
}

func (s *KeeperTestSuite) TestProposalInvariant() {
	testCases := map[string]struct {
		malleate func(ctx sdk.Context)
		valid    bool
	}{
		"invariant not broken": {
			valid: true,
		},
		"active old proposal exists": {
			malleate: func(ctx sdk.Context) {
				info := s.impl.GetFoundationInfo(ctx)
				info.Version--
				s.impl.SetFoundationInfo(ctx, info)
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			invariant := internal.ProposalInvariant(s.impl)
			_, broken := invariant(ctx)
			s.Require().Equal(!tc.valid, broken)
		})
	}
}
