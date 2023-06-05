package keeper_test

import (
	"bytes"
	"compress/zlib"
	"time"

	"github.com/klauspost/compress/zstd"

	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

func (s *KeeperTestSuite) TestDecompressCCBatch() {
	src := types.CCBatch{
		ShouldStartAtFrame: 30,
		Frames: []*types.CCBatchFrame{
			{
				Hedaer: &types.CCBatchHeader{
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
				batch, err := s.keeper.DecompressCCBatch(cb)
				s.Require().NoError(err)
				s.Require().Equal(src.ShouldStartAtFrame, batch.ShouldStartAtFrame)
				s.Require().Equal(src.Frames, batch.Frames)
			} else {
				if tc.compressor == nil {
					cb := types.CompressedCCBatch{
						Compression: types.OptionEmpty,
						Data:        serializedSrc,
					}
					_, err := s.keeper.DecompressCCBatch(cb)
					s.Require().Error(err)
				} else {
					cb := types.CompressedCCBatch{
						Compression: tc.opt,
						Data:        tc.compressor(serializedSrc),
					}
					_, err := s.keeper.DecompressCCBatch(cb)
					s.Require().Error(err)
				}
			}
		})
	}
}
