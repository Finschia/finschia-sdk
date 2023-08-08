package keeper_test

import (
	"time"

	"github.com/Finschia/finschia-rdk/types/query"
	"github.com/Finschia/finschia-rdk/x/or/da/types"
)

func (s *KeeperTestSuite) TestParams() {
	testCases := map[string]struct {
		req      *types.QueryParamsRequest
		expError string
	}{
		"invalid request": {
			nil,
			"invalid request",
		},
		"valid request": {
			req: &types.QueryParamsRequest{},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			resp, err := s.queryServer.Params(s.goCtx, tc.req)
			if tc.expError != "" {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expError)
				s.Require().Nil(resp)
				return
			}
			s.Require().NoError(err)
			s.Require().Equal(types.DefaultParams(), resp.Params)
		})
	}
}

func (s *KeeperTestSuite) TestCCState() {
	testCases := map[string]struct {
		name     string
		expError string
		postTest func(res *types.QueryCCStateResponse)
	}{
		"invalid request": {
			"null",
			"invalid request",
			nil,
		},
		"empty rollup name": {
			"",
			"empty rollup name",
			nil,
		},
		"valid request": {
			"default-rollup1",
			"",
			func(res *types.QueryCCStateResponse) {
				s.Require().Equal(&types.CCState{
					Base:             1,
					Height:           2,
					L1Height:         23,
					ProcessedL2Block: 4,
					Timestamp:        time.Unix(2500, 0).UTC(),
				}, res.State)
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			var req *types.QueryCCStateRequest
			if tc.name != "null" {
				req = &types.QueryCCStateRequest{
					RollupName: tc.name,
				}
			}
			resp, err := s.queryServer.CCState(s.goCtx, req)
			if tc.expError != "" {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expError)
				s.Require().Nil(resp)
				return
			}
			s.Require().NoError(err)
			tc.postTest(resp)
		})
	}
}

func (s *KeeperTestSuite) TestCCRef() {
	testCases := map[string]struct {
		name     string
		height   uint64
		expError string
		postTest func(res *types.QueryCCRefResponse)
	}{
		"invalid request": {
			"null",
			0,
			"invalid request",
			nil,
		},
		"empty rollup name": {
			"",
			0,
			"empty rollup name",
			nil,
		},
		"valid request": {
			"default-rollup1",
			2,
			"",
			func(res *types.QueryCCRefResponse) {
				s.Require().Equal(&types.CCRef{
					TxHash:      []uint8{0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55},
					BatchSize:   2,
					TotalFrames: 4,
					MsgIndex:    0xffffffff,
					BatchRoot:   []uint8{0xad, 0xde, 0x32, 0xa4, 0x40, 0xe8, 0x78, 0x9c, 0x82, 0x1, 0x26, 0x50, 0x7f, 0xa9, 0xaf, 0xb8, 0x9d, 0xbd, 0xaa, 0x68, 0xde, 0xf3, 0x3f, 0x4b, 0x0, 0x19, 0xc1, 0xbd, 0x31, 0xa1, 0x18, 0x9c},
				}, res.Ref)
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			var req *types.QueryCCRefRequest
			if tc.name != "null" {
				req = &types.QueryCCRefRequest{
					RollupName:  tc.name,
					BatchHeight: tc.height,
				}
			}
			resp, err := s.queryServer.CCRef(s.goCtx, req)
			if tc.expError != "" {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expError)
				s.Require().Nil(resp)
				return
			}
			s.Require().NoError(err)
			tc.postTest(resp)
		})
	}
}

func (s *KeeperTestSuite) TestCCRefs() {
	testCases := map[string]struct {
		name     string
		height   uint64
		expError string
		postTest func(res *types.QueryCCRefsResponse)
	}{
		"invalid request": {
			"null",
			0,
			"invalid request",
			nil,
		},
		"empty rollup name": {
			"",
			0,
			"empty rollup name",
			nil,
		},
		"valid request": {
			"default-rollup1",
			2,
			"",
			func(res *types.QueryCCRefsResponse) {
				expected := []*types.CCRef{{
					TxHash:      []uint8{0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55},
					BatchSize:   2,
					TotalFrames: 2,
					MsgIndex:    0xffffffff,
					BatchRoot:   []uint8{0x20, 0x2f, 0xda, 0x6f, 0x11, 0x55, 0xa8, 0xad, 0xca, 0x8e, 0xf3, 0xb1, 0x59, 0x65, 0xae, 0x3c, 0xb1, 0x74, 0x6f, 0xc0, 0x4b, 0x29, 0x41, 0xa2, 0xb8, 0x46, 0xb3, 0x2c, 0x16, 0xa, 0x4b, 0x3},
				},
					{
						TxHash:      []uint8{0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55},
						BatchSize:   2,
						TotalFrames: 4,
						MsgIndex:    0xffffffff,
						BatchRoot:   []uint8{0xad, 0xde, 0x32, 0xa4, 0x40, 0xe8, 0x78, 0x9c, 0x82, 0x1, 0x26, 0x50, 0x7f, 0xa9, 0xaf, 0xb8, 0x9d, 0xbd, 0xaa, 0x68, 0xde, 0xf3, 0x3f, 0x4b, 0x0, 0x19, 0xc1, 0xbd, 0x31, 0xa1, 0x18, 0x9c},
					},
				}

				for i, ref := range res.Refs {
					s.Require().Equal(expected[i], ref)
				}
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			var req *types.QueryCCRefsRequest
			if tc.name != "null" {
				req = &types.QueryCCRefsRequest{
					RollupName: tc.name,
					Pagination: &query.PageRequest{},
				}
			}
			resp, err := s.queryServer.CCRefs(s.goCtx, req)
			if tc.expError != "" {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expError)
				s.Require().Nil(resp)
				return
			}
			s.Require().NoError(err)
			tc.postTest(resp)
		})
	}
}

func (s *KeeperTestSuite) TestQueueTxState() {
	testCases := map[string]struct {
		name     string
		expError string
		postTest func(res *types.QueryQueueTxStateResponse)
	}{
		"invalid request": {
			"null",
			"invalid request",
			nil,
		},
		"empty rollup name": {
			"",
			"empty rollup name",
			nil,
		},
		"valid request": {
			"default-rollup1",
			"",
			func(res *types.QueryQueueTxStateResponse) {
				s.Require().Equal(&types.QueueTxState{
					ProcessedQueueIndex: 1,
					NextQueueIndex:      3,
				}, res.State)
			},
		},
		"valid request2": {
			"default-rollup2",
			"",
			func(res *types.QueryQueueTxStateResponse) {
				s.Require().Equal(&types.QueueTxState{
					ProcessedQueueIndex: 0,
					NextQueueIndex:      2,
				}, res.State)
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			var req *types.QueryQueueTxStateRequest
			if tc.name != "null" {
				req = &types.QueryQueueTxStateRequest{
					RollupName: tc.name,
				}
			}
			resp, err := s.queryServer.QueueTxState(s.goCtx, req)
			if tc.expError != "" {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expError)
				s.Require().Nil(resp)
				return
			}
			s.Require().NoError(err)
			tc.postTest(resp)
		})
	}
}

func (s *KeeperTestSuite) TestQueueTx() {
	testCases := map[string]struct {
		name     string
		index    uint64
		expError string
		postTest func(res *types.QueryQueueTxResponse)
	}{
		"invalid request": {
			"null",
			0,
			"invalid request",
			nil,
		},
		"empty rollup name": {
			"",
			0,
			"empty rollup name",
			nil,
		},
		"queue tx non exist": {
			"default-rollup2",
			2,
			"queue tx not found",
			nil,
		},
		"valid request": {
			"default-rollup1",
			1,
			"",
			func(res *types.QueryQueueTxResponse) {
				s.Require().Equal(&types.L1ToL2Queue{
					Txraw:  []byte("first qtx"),
					Status: types.QUEUE_TX_SUBMITTED,
				}, res.Tx)
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			var req *types.QueryQueueTxRequest
			if tc.name != "null" {
				req = &types.QueryQueueTxRequest{
					RollupName: tc.name,
					QueueIndex: tc.index,
				}
			}
			resp, err := s.queryServer.QueueTx(s.goCtx, req)
			if tc.expError != "" {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expError)
				s.Require().Nil(resp)
				return
			}
			s.Require().NoError(err)
			tc.postTest(resp)
		})
	}
}

func (s *KeeperTestSuite) TestQueueTxs() {
	testCases := map[string]struct {
		name     string
		expError string
		postTest func(res *types.QueryQueueTxsResponse)
	}{
		"invalid request": {
			"null",
			"invalid request",
			nil,
		},
		"empty rollup name": {
			"",
			"empty rollup name",
			nil,
		},
		"valid request": {
			"default-rollup1",
			"",
			func(res *types.QueryQueueTxsResponse) {
				expected := []*types.L1ToL2Queue{
					{
						Txraw:  []byte("first qtx"),
						Status: types.QUEUE_TX_SUBMITTED,
					},
					{
						Txraw:  []byte("second qtx"),
						Status: types.QUEUE_TX_PENDING,
					},
				}
				for i, tx := range res.Txs {
					s.Require().Equal(expected[i], tx)
				}
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			var req *types.QueryQueueTxsRequest
			if tc.name != "null" {
				req = &types.QueryQueueTxsRequest{
					RollupName: tc.name,
					Pagination: &query.PageRequest{},
				}
			}
			resp, err := s.queryServer.QueueTxs(s.goCtx, req)
			if tc.expError != "" {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expError)
				s.Require().Nil(resp)
				return
			}
			s.Require().NoError(err)
			tc.postTest(resp)
		})
	}
}

func (s *KeeperTestSuite) TestMappedBatch() {
	testCases := map[string]struct {
		name     string
		l2height uint64
		expError string
		postTest func(res *types.QueryMappedBatchResponse)
	}{
		"invalid request": {
			"null",
			0,
			"invalid request",
			nil,
		},
		"empty rollup name": {
			"",
			0,
			"empty rollup name",
			nil,
		},
		"invalid l2 height": {
			"default-rollup1",
			0,
			"invalid rollup height",
			nil,
		},
		"valid request": {
			"default-rollup1",
			4,
			"",
			func(res *types.QueryMappedBatchResponse) {
				expected := &types.CCRef{
					TxHash:      []uint8{0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24, 0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55},
					BatchSize:   2,
					TotalFrames: 4,
					MsgIndex:    0xffffffff,
					BatchRoot:   []uint8{0xad, 0xde, 0x32, 0xa4, 0x40, 0xe8, 0x78, 0x9c, 0x82, 0x1, 0x26, 0x50, 0x7f, 0xa9, 0xaf, 0xb8, 0x9d, 0xbd, 0xaa, 0x68, 0xde, 0xf3, 0x3f, 0x4b, 0x0, 0x19, 0xc1, 0xbd, 0x31, 0xa1, 0x18, 0x9c},
				}
				s.Require().Equal(expected, res.Ref)
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			var req *types.QueryMappedBatchRequest
			if tc.name != "null" {
				req = &types.QueryMappedBatchRequest{
					RollupName: tc.name,
					L2Height:   tc.l2height,
				}
			}
			resp, err := s.queryServer.MappedBatch(s.goCtx, req)
			if tc.expError != "" {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expError)
				s.Require().Nil(resp)
				return
			}
			s.Require().NoError(err)
			tc.postTest(resp)
		})
	}
}
