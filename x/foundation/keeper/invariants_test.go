package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/keeper"
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
				s.keeper.SetPool(ctx, foundation.Pool{
					Treasury: sdk.NewDecCoinsFromCoins(balance...),
				})
			},
		},
	}

	for name, tc := range testCases {
		ctx, _ := s.ctx.CacheContext()
		if tc.malleate != nil {
			tc.malleate(ctx)
		}

		invariant := keeper.ModuleAccountInvariant(s.keeper)
		_, broken := invariant(ctx)
		s.Require().Equal(!tc.valid, broken, name)
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
				info := s.keeper.GetFoundationInfo(ctx)
				info.Version--
				s.keeper.SetFoundationInfo(ctx, info)
			},
		},
	}

	for name, tc := range testCases {
		ctx, _ := s.ctx.CacheContext()
		if tc.malleate != nil {
			tc.malleate(ctx)
		}

		invariant := keeper.ProposalInvariant(s.keeper)
		_, broken := invariant(ctx)
		s.Require().Equal(!tc.valid, broken, name)
	}
}
