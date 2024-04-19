package keeper_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/keeper"
	"github.com/Finschia/finschia-sdk/x/fswap/testutil"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

type KeeperTestSuite struct {
	suite.Suite
	sdkCtx     sdk.Context
	goCtx      context.Context
	keeper     keeper.Keeper
	bankKeeper types.BankKeeper

	msgServer types.MsgServer
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	bankKeeper := testutil.NewMockBankKeeper(ctrl)
	bankKeeper.EXPECT().GetSupply(gomock.Any(), "cony").Return(testutil.MockOldCoin)
	s.bankKeeper = bankKeeper
	checkTx := false
	app := simapp.Setup(checkTx)
	s.sdkCtx = app.BaseApp.NewContext(checkTx, tmproto.Header{})
	s.goCtx = sdk.WrapSDKContext(s.sdkCtx)
	s.keeper = app.FswapKeeper

	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
}
