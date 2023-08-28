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

func (s IntegrationTestSuite) TestGetBlocksByHash() {
	val := s.network.Validators[0]
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	height := int64(1)
	blockResult, err := tmservice.GetBlock(ctx.Context(), val.ClientCtx, &height)
	s.Require().NoError(err)

	blockHash := blockResult.Block.Hash()
	blockResult2, err := tmservice.GetBlockByHash(val.ClientCtx, blockHash)
	s.Require().NoError(err)
	s.Require().Equal(blockResult2.Block.Height, blockResult.Block.Height)
}

func (s IntegrationTestSuite) TestGetBlockResultsByHeight() {
	val := s.network.Validators[0]

	height := int64(1)
	blockResult, err := tmservice.GetBlockResultsByHeight(val.ClientCtx, &height)
	s.Require().NoError(err)
	s.Require().Equal(height, blockResult.Height)
	s.Require().NotNil(blockResult.BeginBlockEvents)
	s.Require().Nil(blockResult.EndBlockEvents)
	s.Require().Nil(blockResult.ValidatorUpdates)
}
