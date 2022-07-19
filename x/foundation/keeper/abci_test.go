package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/keeper"
)

func (s *KeeperTestSuite) TestBeginBlocker() {
	ctx, _ := s.ctx.CacheContext()

	s.keeper.SetParams(ctx, &foundation.Params{
		Enabled:       true,
		FoundationTax: sdk.MustNewDecFromStr("0.5"),
	})

	before := s.keeper.GetTreasury(ctx)
	s.Require().Equal(1, len(before))
	s.Require().Equal(s.balance, before[0].Amount)

	// collect
	keeper.BeginBlocker(ctx, s.keeper)

	after := s.keeper.GetTreasury(ctx)
	s.Require().Equal(1, len(after))
	// s.balance + s.balance * 0.5
	s.Require().Equal(s.balance.Add(s.balance.Quo(sdk.NewInt(2))), after[0].Amount)
}
