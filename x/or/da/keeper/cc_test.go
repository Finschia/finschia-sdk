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
		"batch size exceeded - zlib": {
			compressor: func(src []byte) []byte {
				var b bytes.Buffer
				w := zlib.NewWriter(&b)
				_, err := w.Write(serializedSrc)
				s.Require().NoError(err)
				err = w.Close()
				s.Require().NoError(err)
				return append(b.Bytes(), make([]byte, types.DefaultCCBatchMaxBytes)...)
			},
			opt: types.OptionZLIB,
			err: true,
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

func (s *KeeperTestSuite) TestSaveCCBatch() {
	rollupName := "rollup1"

	testCases := map[string]struct {
		src       *types.CCBatch
		malleate  func(b *types.CCBatch)
		postCheck func()
		isErr     bool
	}{
		"wrong ShouldStartAtFrame": {
			src: &types.CCBatch{
				ShouldStartAtFrame: 1000,
			},
			malleate: func(b *types.CCBatch) {
				_ = prepareCCPreset(s, rollupName)
			},
			isErr: true,
		},
		"empty frame": {
			src: &types.CCBatch{
				Frames: []*types.CCBatchFrame{
					{
						Elements: nil,
					},
				},
			},
			malleate: func(b *types.CCBatch) {
				_ = prepareCCPreset(s, rollupName)
				b.ShouldStartAtFrame = calShouldStartAtFrame(s, rollupName)
			},
			isErr: true,
		},
		"wrong frame header l2 height": {
			src: &types.CCBatch{
				Frames: []*types.CCBatchFrame{
					{
						Header: &types.CCBatchHeader{
							ParentHash: []byte("parent_hash"),
							L2Height:   10,
						},
						Elements: []*types.CCBatchElement{
							{
								Txraw: []byte("txraw"),
							},
						},
					},
				},
			},
			malleate: func(b *types.CCBatch) {
				_ = prepareCCPreset(s, rollupName)
				b.ShouldStartAtFrame = calShouldStartAtFrame(s, rollupName)
			},
			isErr: true,
		},
		"empty parent hash": {
			src: &types.CCBatch{
				Frames: []*types.CCBatchFrame{
					{
						Header: &types.CCBatchHeader{
							ParentHash: nil,
							L1Height:   10,
							Timestamp:  time.Unix(2000, 0).UTC(),
						},
						Elements: []*types.CCBatchElement{
							{
								Txraw: []byte("txraw"),
							},
						},
					},
				},
			},
			malleate: func(b *types.CCBatch) {
				_ = prepareCCPreset(s, rollupName)
				b.ShouldStartAtFrame = calShouldStartAtFrame(s, rollupName)
				b.Frames[0].Header.L2Height = calL2Height(s, rollupName)
			},
			isErr: true,
		},
		"outdated frame - timestamp": {
			src: &types.CCBatch{
				Frames: []*types.CCBatchFrame{
					{
						Header: &types.CCBatchHeader{
							ParentHash: []byte("parent_hash2"),
							L1Height:   10,
							Timestamp:  time.Unix(500, 0).UTC(),
						},
						Elements: []*types.CCBatchElement{
							{
								Txraw: []byte("txraw"),
							},
						},
					},
				},
			},
			malleate: func(b *types.CCBatch) {
				_ = prepareCCPreset(s, rollupName)
				b.ShouldStartAtFrame = calShouldStartAtFrame(s, rollupName)
				b.Frames[0].Header.L2Height = calL2Height(s, rollupName)
			},
			isErr: true,
		},
		"outdated frame - l1 block height": {
			src: &types.CCBatch{
				Frames: []*types.CCBatchFrame{
					{
						Header: &types.CCBatchHeader{
							ParentHash: []byte("parent_hash2"),
							L1Height:   8,
							Timestamp:  time.Unix(1100, 0).UTC(),
						},
						Elements: []*types.CCBatchElement{
							{
								Txraw: []byte("txraw"),
							},
						},
					},
				},
			},
			malleate: func(b *types.CCBatch) {
				_ = prepareCCPreset(s, rollupName)
				b.ShouldStartAtFrame = calShouldStartAtFrame(s, rollupName)
				b.Frames[0].Header.L2Height = calL2Height(s, rollupName)
			},
			isErr: true,
		},
		"process queue txs in the wrong order": {
			src: &types.CCBatch{
				Frames: []*types.CCBatchFrame{
					{
						Header: &types.CCBatchHeader{
							ParentHash: []byte("parent_hash2"),
							L1Height:   11,
							Timestamp:  time.Unix(1200, 0).UTC(),
						},
						Elements: []*types.CCBatchElement{
							{
								Txraw:      nil,
								QueueIndex: 3,
							},
						},
					},
				},
			},
			malleate: func(b *types.CCBatch) {
				saveQueueTx(s.storeKey, s.ctx, rollupName, 1, &types.L1ToL2Queue{Txraw: []byte("txraw"), Status: types.QUEUE_TX_PENDING})
				saveQueueTxState(s.storeKey, s.ctx, rollupName, &types.QueueTxState{ProcessedQueueIndex: 0, NextQueueIndex: 2})
				err := prepareCCPreset(s, rollupName)
				b.ShouldStartAtFrame = calShouldStartAtFrame(s, rollupName)
				b.Frames[0].Header.L2Height = calL2Height(s, rollupName)

				if err != nil {
					_ = s.keeper.SaveCCBatch(s.ctx, rollupName, &types.CCBatch{
						ShouldStartAtFrame: calShouldStartAtFrame(s, rollupName),
						Frames: []*types.CCBatchFrame{
							{
								Header: &types.CCBatchHeader{
									ParentHash: []byte("parent_hash2"),
									L2Height:   2,
									L1Height:   10,
									Timestamp:  time.Unix(1100, 0).UTC(),
								},
								Elements: []*types.CCBatchElement{
									{
										Txraw: []byte("txraw"),
									},
									{
										Txraw:      nil,
										QueueIndex: 1,
									},
								},
							},
						},
					})
					b.ShouldStartAtFrame = calShouldStartAtFrame(s, rollupName)
					b.Frames[0].Header.L2Height = calL2Height(s, rollupName)
				}
			},
			isErr: true,
		},
		"valid batch": {
			src: &types.CCBatch{
				Frames: []*types.CCBatchFrame{
					{
						Header: &types.CCBatchHeader{
							ParentHash: []byte("parent_hash"),
							L1Height:   21,
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
							Timestamp:  time.Unix(2300, 0).UTC(),
						},
						Elements: []*types.CCBatchElement{
							{
								Txraw: []byte("txraw"),
							},
						},
					},
				},
			},
			malleate: func(b *types.CCBatch) {
				_ = prepareCCPreset(s, rollupName)
				b.ShouldStartAtFrame = calShouldStartAtFrame(s, rollupName)
				b.Frames[0].Header.L2Height = calL2Height(s, rollupName)
				b.Frames[1].Header.L2Height = calL2Height(s, rollupName) + 1
			},
			isErr: false,
			postCheck: func() {
				ccState, err := s.keeper.GetCCState(s.ctx, rollupName)
				s.Require().NoError(err)
				_, err = s.keeper.GetQueueTxState(s.ctx, rollupName)
				s.Require().NoError(err)
				ccRef, err := s.keeper.GetCCRef(s.ctx, rollupName, ccState.Height)
				s.Require().NoError(err)
				for i := uint64(1); i <= ccRef.TotalFrames; i++ {
					_, err = s.keeper.GetL2HeightBatchMap(s.ctx, rollupName, i)
					s.Require().NoError(err)
				}
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			if tc.malleate != nil {
				tc.malleate(tc.src)
			}
			err := s.keeper.SaveCCBatch(s.ctx, rollupName, tc.src)
			if tc.isErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				if tc.postCheck != nil {
					tc.postCheck()
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestSaveQueueTx() {
	s.ctx = s.ctx.WithBlockHeight(10)
	t := time.Now().UTC()
	s.ctx = s.ctx.WithBlockTime(t)

	err := s.keeper.SaveQueueTx(s.ctx, "rollup1", []byte("qtx1"), 90000, 10)
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
	evt, err := sdktypes.ParseTypedEvent(s.ctx.EventManager().ABCIEvents()[0])
	s.Require().NoError(err)
	parsedEvt := evt.(*types.EventSaveQueueTx)
	s.Require().Equal("rollup1", parsedEvt.RollupName)
	s.Require().Equal(uint64(2), parsedEvt.NextQueueIndex)
	s.Require().Equal(uint64(8966), parsedEvt.ExtraConsumedGas)
	s.Require().Equal(uint64(90000), parsedEvt.L2GasLimit)

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

func calShouldStartAtFrame(s *KeeperTestSuite, rollupName string) (ShouldStartAtFrame uint64) {
	state, err := s.keeper.GetCCState(s.ctx, rollupName)
	if err == nil {
		ref, _ := s.keeper.GetCCRef(s.ctx, rollupName, state.Height)
		ShouldStartAtFrame = ref.TotalFrames
	}
	return
}

func calL2Height(s *KeeperTestSuite, rollupName string) (height uint64) {
	state, err := s.keeper.GetCCState(s.ctx, rollupName)
	if err == nil {
		height = state.ProcessedL2Block + 1
	}
	return
}

func prepareCCPreset(s *KeeperTestSuite, rollupName string) (err error) {
	err = s.keeper.SaveCCBatch(s.ctx, rollupName, &types.CCBatch{
		ShouldStartAtFrame: calShouldStartAtFrame(s, rollupName),
		Frames: []*types.CCBatchFrame{
			{
				Header: &types.CCBatchHeader{
					ParentHash: nil,
					L2Height:   1,
					L1Height:   10,
					Timestamp:  time.Unix(1000, 0).UTC(),
				},
				Elements: []*types.CCBatchElement{
					{
						Txraw: []byte("txraw"),
					},
				},
			},
		},
	})
	return
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
