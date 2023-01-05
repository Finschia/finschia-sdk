package keeper_test

import (
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	authkeeper "github.com/line/lbm-sdk/x/auth/keeper"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	"github.com/line/lbm-sdk/x/distribution/keeper"
	mocktypes "github.com/line/lbm-sdk/x/distribution/mocks/types"
	"github.com/line/lbm-sdk/x/distribution/types"
	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/suite"
	"testing"
)

func TestFeePoolSuite(t *testing.T) {
	suite.Run(t, new(feePoolTestSuite))
}

type feePoolTestSuite struct {
	suite.Suite
	mockBank *mocktypes.BankKeeper
	app      *simapp.SimApp
	cut      keeper.Keeper
}

func (s *feePoolTestSuite) SetupTest() {
	s.app = simapp.Setup(false)
	cdc := s.app.AppCodec()
	s.mockBank = &mocktypes.BankKeeper{}

	maccPerms := map[string][]string{types.ModuleName: nil}
	s.app.AccountKeeper = authkeeper.NewAccountKeeper(cdc, s.app.GetKey(authtypes.StoreKey), s.app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms)
	s.cut = keeper.NewKeeper(cdc, s.app.GetKey(types.StoreKey), s.app.GetSubspace(types.ModuleName), s.app.AccountKeeper, s.mockBank, nil, authtypes.FeeCollectorName, nil)
}

func (s *feePoolTestSuite) TearDownTest() {
	s.mockBank.AssertExpectations(s.T())
}

func (s *feePoolTestSuite) TestFeePoolHappyPath() {
	// Arrange
	ctx := s.app.BaseApp.NewContext(false, ocproto.Header{})
	addr := simapp.AddTestAddrs(s.app, ctx, 1, sdk.NewInt(1000))
	feePool := types.InitialFeePool()
	goodAmount := sdk.NewCoin("mytoken", sdk.NewInt(950))
	totalAmount := sdk.NewDecCoin("mytoken", sdk.NewInt(1000))
	feePool.CommunityPool = feePool.CommunityPool.Add(totalAmount)
	s.cut.SetFeePool(ctx, feePool)
	s.mockBank.On("SendCoinsFromModuleToAccount", ctx, types.ModuleName, addr[0], sdk.NewCoins(goodAmount)).Return(nil)

	// Act
	err := s.cut.DistributeFromFeePool(ctx, sdk.NewCoins(goodAmount), addr[0])

	// Assert
	require.NoError(s.T(), err)
	actual := s.cut.GetFeePool(ctx)
	expected := sdk.NewDecCoins(sdk.NewDecCoin("mytoken", sdk.NewInt(50)))
	require.True(s.T(), actual.CommunityPool.IsEqual(expected))
}

func (s *feePoolTestSuite) TestFeePoolThrowBadDistributionErrWhenTryToDistributeMoreThenItHas() {
	// Arrange
	ctx := s.app.BaseApp.NewContext(false, ocproto.Header{})
	addr := simapp.AddTestAddrs(s.app, ctx, 1, sdk.NewInt(1000))
	feePool := types.InitialFeePool()
	badAmount := sdk.NewCoin("mytoken", sdk.NewInt(1500))
	totalAmount := sdk.NewDecCoin("mytoken", sdk.NewInt(1000))
	feePool.CommunityPool = feePool.CommunityPool.Add(totalAmount)
	s.cut.SetFeePool(ctx, feePool)

	// Act
	err := s.cut.DistributeFromFeePool(ctx, sdk.NewCoins(badAmount), addr[0])

	// Assert
	require.ErrorIs(s.T(), err, types.ErrBadDistribution)
}
