package keeper_test

import (
	"bytes"
	"compress/zlib"
	"time"

	"github.com/klauspost/compress/zstd"

	sdktypes "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

func (s *KeeperTestSuite) TestDecompressCCBatch() {
	src := types.CCBatch{
		ShouldStartAtFrame: 30,
		Frames: []*types.CCBatchFrame{
			{
				Header: &types.CCBatchHeader{
					ParentHash: []byte("parent_hash"),
					Timestamp:  time.Now().UTC(),
					L2Height:   10,
					L1Height:   100,
				},
				Elements: []*types.CCBatchElement{
					{
						Txraw:      []byte("Transactions"),
						QueueIndex: 10,
					},
				},
			},
		},
	}

	serializedSrc := s.encCfg.Marshaler.MustMarshal(&src)

	testCases := map[string]struct {
		compressor func(src []byte) []byte
		opt        types.CompressionOption
		err        bool
	}{
		"success - zlib": {
			compressor: func(src []byte) []byte {
				var b bytes.Buffer
				w := zlib.NewWriter(&b)
				_, err := w.Write(serializedSrc)
				s.Require().NoError(err)
				err = w.Close()
				s.Require().NoError(err)
				return b.Bytes()
			},
			opt: types.OptionZLIB,
			err: false,
		},
		"wrong compressed - zstd": {
			compressor: func(src []byte) []byte {
				var b bytes.Buffer
				w, _ := zstd.NewWriter(&b)
				_, err := w.Write(serializedSrc)
				s.Require().NoError(err)
				err = w.Close()
				s.Require().NoError(err)
				return b.Bytes()
			},
			opt: types.OptionZSTD,
			err: true,
		},
		"option empty": {
			compressor: nil,
			opt:        types.OptionEmpty,
			err:        true,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			if !tc.err {
				cb := types.CompressedCCBatch{
					Compression: tc.opt,
					Data:        tc.compressor(serializedSrc),
				}
				batch, err := s.keeper.DecompressCCBatch(s.ctx, cb)
				s.Require().NoError(err)
				s.Require().Equal(src.ShouldStartAtFrame, batch.ShouldStartAtFrame)
				s.Require().Equal(src.Frames, batch.Frames)
			} else {
				if tc.compressor == nil {
					cb := types.CompressedCCBatch{
						Compression: types.OptionEmpty,
						Data:        serializedSrc,
					}
					_, err := s.keeper.DecompressCCBatch(s.ctx, cb)
					s.Require().Error(err)
				} else {
					cb := types.CompressedCCBatch{
						Compression: tc.opt,
						Data:        tc.compressor(serializedSrc),
					}
					_, err := s.keeper.DecompressCCBatch(s.ctx, cb)
					s.Require().Error(err)
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestSaveQueueTx() {
	s.ctx = s.ctx.WithBlockHeight(10)
	t := time.Now().UTC()
	s.ctx = s.ctx.WithBlockTime(t)

	err := s.keeper.SaveQueueTx(s.ctx, "rollup1", []byte("qtx1"), 90000)
	s.Require().NoError(err)
	state, err := s.keeper.GetQueueTxState(s.ctx, "rollup1")
	s.Require().NoError(err)
	s.Require().Equal(uint64(2), state.NextQueueIndex)
	s.Require().Equal(uint64(0), state.ProcessedQueueIndex)
	qtx, err := s.keeper.GetQueueTx(s.ctx, "rollup1", 1)
	s.Require().NoError(err)
	s.Require().Equal([]byte("qtx1"), qtx.Txraw)
	s.Require().Equal(types.QUEUE_TX_PENDING, qtx.Status)
	s.Require().Equal(int64(10), qtx.L1Height)
	s.Require().Equal(t, qtx.Timestamp)
}

func (s *KeeperTestSuite) TestUpdateQueueTxsStatus() {
	sampleIdx := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s.ctx = s.ctx.WithBlockHeight(604)
	for _, i := range sampleIdx {
		saveQueueTx(s.storeKey, s.ctx, "rollup1", i, &types.L1ToL2Queue{
			Timestamp: time.Now(),
			Status:    types.QUEUE_TX_PENDING,
			Txraw:     nil,
			L1Height:  int64(i),
		})
	}

	testCases := map[string]struct {
		beforeState *types.QueueTxState
		afterState  *types.QueueTxState
		isErr       bool
		expired     int
	}{
		"queue tx never registered": {
			beforeState: new(types.QueueTxState),
			isErr:       false,
		},
		"queue tx not found": {
			beforeState: &types.QueueTxState{
				ProcessedQueueIndex: 9,
				NextQueueIndex:      13,
			},
			isErr: true,
		},
		"valid case": {
			beforeState: &types.QueueTxState{
				ProcessedQueueIndex: 0,
				NextQueueIndex:      10,
			},
			afterState: &types.QueueTxState{
				ProcessedQueueIndex: 4,
				NextQueueIndex:      10,
			},
			isErr:   false,
			expired: 4,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			saveQueueTxState(s.storeKey, s.ctx, "rollup1", tc.beforeState)
			err := s.keeper.UpdateQueueTxsStatus(s.ctx)
			if tc.isErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)

				if tc.beforeState.NextQueueIndex != 0 {
					store := s.ctx.KVStore(s.storeKey)
					bz := store.Get(types.GetQueueTxStateStoreKey("rollup1"))
					var state types.QueueTxState
					s.Require().NoError(state.Unmarshal(bz))
					s.Require().Equal(tc.afterState.ProcessedQueueIndex, state.ProcessedQueueIndex)
					s.Require().Equal(tc.afterState.NextQueueIndex, state.NextQueueIndex)

					actualExpired := 0
					for _, i := range sampleIdx {
						bz := store.Get(types.GetCCQueueTxKey("rollup1", i))
						var tx types.L1ToL2Queue
						s.Require().NoError(tx.Unmarshal(bz))
						if tx.Status == types.QUEUE_TX_EXPIRED {
							actualExpired++
						}
					}
					s.Require().Equal(tc.expired, actualExpired)
				}
			}
		})
	}
}

func saveQueueTx(skey sdktypes.StoreKey, ctx sdktypes.Context, rollupName string, idx uint64, elem *types.L1ToL2Queue) {
	store := ctx.KVStore(skey)
	bz, _ := elem.Marshal()
	store.Set(types.GetCCQueueTxKey(rollupName, idx), bz)
}

func saveQueueTxState(skey sdktypes.StoreKey, ctx sdktypes.Context, rollupName string, state *types.QueueTxState) {
	store := ctx.KVStore(skey)
	bz, _ := state.Marshal()
	store.Set(types.GetQueueTxStateStoreKey(rollupName), bz)
}
