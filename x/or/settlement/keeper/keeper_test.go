package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/baseapp"
	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/settlement/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app *simapp.SimApp
	ctx sdk.Context

	queryClient types.QueryClient
}

func (s *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	app.SettlementKeeper.SetParams(ctx, types.DefaultParams())

	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.SettlementKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	s.app = app
	s.ctx = ctx
	s.queryClient = queryClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
