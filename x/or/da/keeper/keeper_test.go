package keeper_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/simapp"
	simappparams "github.com/Finschia/finschia-sdk/simapp/params"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdktypes "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/da/keeper"
	"github.com/Finschia/finschia-sdk/x/or/da/testutil"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

type KeeperTestSuite struct {
	suite.Suite

	storeKey    sdk.StoreKey
	ctx         sdk.Context
	goCtx       context.Context
	keeper      keeper.Keeper
	queryServer types.QueryServer
	msgServer   types.MsgServer
	encCfg      simappparams.EncodingConfig

	addrs    []sdk.AccAddress
	initQtxs [][]byte
}

func (s *KeeperTestSuite) SetupTest() {
	s.encCfg = simapp.MakeTestEncodingConfig()
	testdata.RegisterInterfaces(s.encCfg.InterfaceRegistry)
	k, ctx, skey := testutil.DaKeeper(s.T(), s.encCfg)
	s.storeKey = skey
	s.ctx = ctx
	s.goCtx = sdk.WrapSDKContext(s.ctx)
	s.keeper = k
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
	s.queryServer = s.keeper

	for i := 0; i < 3; i++ {
		s.addrs = append(s.addrs, testutil.AccAddress())
	}

	err := s.keeper.SetParams(s.ctx, types.DefaultParams())
	s.Require().NoError(err)

	mockTxs := s.genMockTxs(3)
	s.initQtxs = mockTxs
	err = s.keeper.SaveQueueTx(s.ctx, "default-rollup1", mockTxs[0], 100000, 5)
	s.Require().NoError(err)
	err = s.keeper.SaveQueueTx(s.ctx, "default-rollup1", mockTxs[1], 100000, 5)
	s.Require().NoError(err)
	err = s.keeper.SaveQueueTx(s.ctx, "default-rollup2", mockTxs[2], 100000, 5)
	s.Require().NoError(err)

	batches := dummyBatches()
	err = s.keeper.SaveCCBatch(ctx, "default-rollup1", batches[0])
	s.Require().NoError(err)
	err = s.keeper.SaveCCBatch(ctx, "default-rollup1", batches[1])
	s.Require().NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func dummyBatches() []*types.CCBatch {
	return []*types.CCBatch{{
		ShouldStartAtFrame: 0,
		Frames: []*types.CCBatchFrame{
			{
				Header: &types.CCBatchHeader{
					ParentHash: []byte("parent_hash"),
					L1Height:   21,
					L2Height:   1,
					Timestamp:  time.Unix(2200, 0).UTC(),
				},
				Elements: []*types.CCBatchElement{
					{
						Txraw: []byte("txraw"),
					},
				},
			},
			{
				Header: &types.CCBatchHeader{
					ParentHash: []byte("parent_hash2"),
					L1Height:   22,
					L2Height:   2,
					Timestamp:  time.Unix(2300, 0).UTC(),
				},
				Elements: []*types.CCBatchElement{
					{
						Txraw: []byte("txraw2"),
					},
				},
			},
		}},
		{
			ShouldStartAtFrame: 2,
			Frames: []*types.CCBatchFrame{
				{
					Header: &types.CCBatchHeader{
						ParentHash: []byte("parent_hash3"),
						L1Height:   23,
						L2Height:   3,
						Timestamp:  time.Unix(2400, 0).UTC(),
					},
					Elements: []*types.CCBatchElement{
						{
							Txraw: []byte("txraw3"),
						},
					},
				},
				{
					Header: &types.CCBatchHeader{
						ParentHash: []byte("parent_hash4"),
						L1Height:   23,
						L2Height:   4,
						Timestamp:  time.Unix(2500, 0).UTC(),
					},
					Elements: []*types.CCBatchElement{
						{
							Txraw:      nil,
							QueueIndex: 1,
						},
					},
				},
			}},
	}
}

func (s *KeeperTestSuite) genMockTxs(numToGenerate int) [][]byte {
	msg := testdata.NewTestMsg(s.addrs[0])
	txCfg := simappparams.MakeTestEncodingConfig().TxConfig
	txs, err := simapp.GenSequenceOfTxs(txCfg, []sdktypes.Msg{msg}, []uint64{0}, []uint64{uint64(0)}, numToGenerate, secp256k1.GenPrivKey())
	s.Require().NoError(err)

	serializedTxs := make([][]byte, numToGenerate)
	for i := 0; i < numToGenerate; i++ {
		serializedTxs[i], err = txCfg.TxEncoder()(txs[i])
		s.Require().NoError(err)
	}

	return serializedTxs
}
