package keeper_test

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection/keeper"
)

func (s *KeeperTestSuite) TestTotalSupplyInvariant() {
	testCases := map[string]struct {
		malleate func(ctx sdk.Context)
		valid    bool
	}{
		"invariant not broken": {
			valid: true,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			invariant := keeper.TotalFTSupplyInvariant(s.keeper)
			_, broken := invariant(ctx)
			s.Require().Equal(!tc.valid, broken)
		})
	}
}
