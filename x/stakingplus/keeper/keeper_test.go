package keeper_test

import (
	"testing"
	"time"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/stakingplus/keeper"
	"github.com/Finschia/finschia-sdk/x/stakingplus/module"
	"github.com/Finschia/finschia-sdk/x/stakingplus/testutil"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx sdk.Context

	foundationKeeper *testutil.MockFoundationKeeper
	stakingMsgServer *testutil.MockStakingMsgServer
	msgServer        stakingtypes.MsgServer
}

func (s *KeeperTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.stakingMsgServer = testutil.NewMockStakingMsgServer(ctrl)
	s.foundationKeeper = testutil.NewMockFoundationKeeper(ctrl)

	key := storetypes.NewKVStoreKey(foundation.StoreKey)
	tkey := storetypes.NewTransientStoreKey("transient_test")
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	cms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, storetypes.StoreTypeTransient, db)
	err := cms.LoadLatestVersion()
	s.Require().NoError(err)
	s.ctx = sdk.NewContext(cms, cmtproto.Header{Time: time.Now()}, false, log.NewNopLogger())

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})
	s.msgServer = keeper.NewMsgServerImpl(s.stakingMsgServer, s.foundationKeeper, encCfg.InterfaceRegistry.SigningContext().ValidatorAddressCodec())
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
