package tmservice_test

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/client/grpc/tmservice"
	"github.com/Finschia/finschia-sdk/simapp"
)

func (s IntegrationTestSuite) TestGetProtoBlock() {
	val := s.network.Validators[0]
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	height := int64(-1)
	blockID, block, err := tmservice.GetProtoBlock(ctx.Context(), val.ClientCtx, &height)
	s.Require().Equal(tmproto.BlockID{}, blockID)
	s.Require().Nil(block)
	s.Require().Error(err)

	height = int64(1)
	_, _, err = tmservice.GetProtoBlock(ctx.Context(), val.ClientCtx, &height)
	s.Require().NoError(err)
}
