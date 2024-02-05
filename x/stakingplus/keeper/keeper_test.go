package keeper_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/runtime"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Finschia/finschia-sdk/x/stakingplus/keeper"
	"github.com/Finschia/finschia-sdk/x/stakingplus/testutil"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx sdk.Context

	app              *runtime.App
	accountKeeper    authkeeper.AccountKeeper
	bankKeeper       bankkeeper.Keeper
	foundationKeeper *testutil.MockFoundationKeeper
	stakingKeeper    *stakingkeeper.Keeper
	msgServer        stakingtypes.MsgServer
}

func (s *KeeperTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.foundationKeeper = testutil.NewMockFoundationKeeper(ctrl)

	app, err := simtestutil.Setup(
		depinject.Configs(
			testutil.AppConfig,
			depinject.Supply(log.NewNopLogger()),
		),
		&s.accountKeeper,
		&s.bankKeeper,
		&s.stakingKeeper,
	)
	s.Require().NoError(err)

	s.app = app
	s.ctx = s.app.BaseApp.NewContext(false)

	s.msgServer = keeper.NewMsgServerImpl(s.stakingKeeper, s.foundationKeeper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
