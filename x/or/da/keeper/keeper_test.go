package keeper_test

import (
	"context"
	"github.com/Finschia/finschia-sdk/simapp"
	"github.com/Finschia/finschia-sdk/simapp/params"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/da/keeper"
	"github.com/Finschia/finschia-sdk/x/or/da/testutil"
	datypes "github.com/Finschia/finschia-sdk/x/or/da/types"
	"github.com/stretchr/testify/suite"
	"testing"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	goCtx       context.Context
	keeper      keeper.Keeper
	queryServer datypes.QueryServer
	msgServer   datypes.MsgServer
	encCfg      params.EncodingConfig

	addrs []sdk.AccAddress
}

func (s *KeeperTestSuite) SetupTest() {
	k, ctx := testutil.DaKeeper(s.T())
	s.ctx = ctx
	s.goCtx = sdk.WrapSDKContext(s.ctx)
	s.keeper = k
	s.encCfg = simapp.MakeTestEncodingConfig()

	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
	s.queryServer = s.keeper

	for i := 0; i < 3; i++ {
		s.addrs = append(s.addrs, testutil.AccAddress())
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
