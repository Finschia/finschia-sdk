package keeper_test

import (
	gocontext "context"
	"testing"
	"time"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	banktypes "github.com/line/lbm-sdk/x/bank/types"
	"github.com/line/lbm-sdk/x/slashing/testslashing"
	"github.com/line/lbm-sdk/x/slashing/types"
)

type SlashingTestSuite struct {
	suite.Suite

	app         *simapp.SimApp
	ctx         sdk.Context
	queryClient types.QueryClient
	addrDels    []sdk.AccAddress
}

func (suite *SlashingTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())
	app.BankKeeper.SetParams(ctx, banktypes.DefaultParams())
	app.SlashingKeeper.SetParams(ctx, testslashing.TestParams())

	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(200))

	info1 := types.NewValidatorSigningInfo(addrDels[0].ToConsAddress(),
		time.Unix(2, 0), false, int64(10), int64(3))
	info2 := types.NewValidatorSigningInfo(addrDels[1].ToConsAddress(),
		time.Unix(2, 0), false, int64(10), int64(4))

	app.SlashingKeeper.SetValidatorSigningInfo(ctx, addrDels[0].ToConsAddress(), info1)
	app.SlashingKeeper.SetValidatorSigningInfo(ctx, addrDels[1].ToConsAddress(), info2)

	suite.app = app
	suite.ctx = ctx
	suite.addrDels = addrDels

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.SlashingKeeper)
	queryClient := types.NewQueryClient(queryHelper)
	suite.queryClient = queryClient
}

func (suite *SlashingTestSuite) TestGRPCQueryParams() {
	queryClient := suite.queryClient
	paramsResp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})

	suite.NoError(err)
	suite.Equal(testslashing.TestParams(), paramsResp.Params)
}

func (suite *SlashingTestSuite) TestGRPCSigningInfo() {
	queryClient := suite.queryClient

	infoResp, err := queryClient.SigningInfo(gocontext.Background(), &types.QuerySigningInfoRequest{ConsAddress: ""})
	suite.Error(err)
	suite.Nil(infoResp)

	consAddr := suite.addrDels[0].ToConsAddress()
	info, found := suite.app.SlashingKeeper.GetValidatorSigningInfo(suite.ctx, consAddr)
	suite.True(found)

	infoResp, err = queryClient.SigningInfo(gocontext.Background(),
		&types.QuerySigningInfoRequest{ConsAddress: consAddr.String()})
	suite.NoError(err)
	suite.Equal(info, infoResp.ValSigningInfo)
}

func (suite *SlashingTestSuite) TestGRPCSigningInfos() {
	queryClient := suite.queryClient

	var signingInfos []types.ValidatorSigningInfo

	suite.app.SlashingKeeper.IterateValidatorSigningInfos(suite.ctx, func(consAddr sdk.ConsAddress, info types.ValidatorSigningInfo) (stop bool) {
		signingInfos = append(signingInfos, info)
		return false
	})

	// verify all values are returned without pagination
	var infoResp, err = queryClient.SigningInfos(gocontext.Background(),
		&types.QuerySigningInfosRequest{Pagination: nil})
	suite.NoError(err)
	suite.Equal(signingInfos, infoResp.Info)

	infoResp, err = queryClient.SigningInfos(gocontext.Background(),
		&types.QuerySigningInfosRequest{Pagination: &query.PageRequest{Limit: 1, CountTotal: true}})
	suite.NoError(err)
	suite.Len(infoResp.Info, 1)
	suite.Equal(signingInfos[0], infoResp.Info[0])
	suite.NotNil(infoResp.Pagination.NextKey)
	suite.Equal(uint64(2), infoResp.Pagination.Total)
}

func TestSlashingTestSuite(t *testing.T) {
	suite.Run(t, new(SlashingTestSuite))
}
