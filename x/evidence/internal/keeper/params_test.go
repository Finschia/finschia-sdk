package keeper_test

import (
	"github.com/line/lbm-sdk/x/evidence/internal/types"
)

func (suite *KeeperTestSuite) TestParams() {
	ctx := suite.ctx.WithIsCheckTx(false)
	suite.Equal(types.DefaultParams(), suite.keeper.GetParams(ctx))
	suite.Equal(types.DefaultMaxEvidenceAge, suite.keeper.MaxEvidenceAge(ctx))
}
