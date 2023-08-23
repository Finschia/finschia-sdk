package testutil

import (
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	dacli "github.com/Finschia/finschia-sdk/x/or/da/client/cli"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

func (s *IntegrationTestSuite) TestCmdQueryParams() {
	val := s.network.Validators[0]

	tcs := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query params": {
			[]string{},
			true,
		},
	}

	for name, tc := range tcs {
		s.Run(name, func() {
			cmd := dacli.CmdQueryParams()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)
			if tc.valid {
				s.Require().NoError(err)
				var res types.QueryParamsResponse
				s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &res))
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCmdQueryCCState() {
	val := s.network.Validators[0]

	tcs := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				rollupName,
			},
			true,
		},
		"invalid query - wrong rollup name": {
			[]string{
				"invalid-rollup",
			},
			false,
		},
	}

	for name, tc := range tcs {
		s.Run(name, func() {
			cmd := dacli.CmdQueryCCState()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)
			if tc.valid {
				s.Require().NoError(err)
				var res types.QueryCCStateResponse
				s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &res))
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCmdCCRef() {
	val := s.network.Validators[0]

	tcs := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				rollupName,
				"1",
			},
			true,
		},
		"invalid query - wrong rollup name": {
			[]string{
				"invalid-rollup",
				"1",
			},
			false,
		},
		"invalid query - wrong index": {
			[]string{
				rollupName,
				"10",
			},
			false,
		},
	}

	for name, tc := range tcs {
		s.Run(name, func() {
			cmd := dacli.CmdCCRef()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)
			if tc.valid {
				s.Require().NoError(err)
				var res types.QueryCCRefResponse
				s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &res))
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCmdCCRefs() {
	val := s.network.Validators[0]

	tcs := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				rollupName,
			},
			true,
		},
		"invalid query - wrong rollup name": {
			[]string{
				"invalid-rollup",
			},
			false,
		},
	}

	for name, tc := range tcs {
		s.Run(name, func() {
			cmd := dacli.CmdCCRefs()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)
			s.Require().NoError(err)
			var res types.QueryCCRefsResponse
			s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &res))
			if tc.valid {
				s.Require().NotNil(res.Refs)
			} else {
				s.Require().Equal(0, len(res.Refs))
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueueTxState() {
	val := s.network.Validators[0]

	tcs := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				rollupName,
			},
			true,
		},
		"invalid query - wrong rollup name": {
			[]string{
				"invalid-rollup",
			},
			false,
		},
	}

	for name, tc := range tcs {
		s.Run(name, func() {
			cmd := dacli.CmdQueueTxState()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)
			if tc.valid {
				s.Require().NoError(err)
				var res types.QueryQueueTxStateResponse
				s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &res))
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueueTx() {
	val := s.network.Validators[0]

	tcs := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				rollupName,
				"1",
			},
			true,
		},
		"invalid query - wrong rollup name": {
			[]string{
				"invalid-rollup",
				"1",
			},
			false,
		},
		"invalid query - wrong queue index": {
			[]string{
				"invalid-rollup",
				"10",
			},
			false,
		},
	}

	for name, tc := range tcs {
		s.Run(name, func() {
			cmd := dacli.CmdQueueTx()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)

			if tc.valid {
				s.Require().NoError(err)
				var res types.QueryQueueTxResponse
				s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &res))
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueueTxs() {
	val := s.network.Validators[0]

	tcs := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				rollupName,
			},
			true,
		},
		"invalid query - wrong rollup name": {
			[]string{
				"invalid-rollup",
			},
			false,
		},
	}

	for name, tc := range tcs {
		s.Run(name, func() {
			cmd := dacli.CmdQueueTxs()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)
			s.Require().NoError(err)
			var res types.QueryQueueTxsResponse
			s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &res))

			if tc.valid {
				s.Require().NotEqual(0, len(res.Txs))
			} else {
				s.Require().Equal(0, len(res.Txs))
			}

		})
	}
}

func (s *IntegrationTestSuite) TestMappedBatch() {
	val := s.network.Validators[0]

	tcs := map[string]struct {
		args  []string
		valid bool
	}{
		"valid query": {
			[]string{
				rollupName,
				"1",
			},
			true,
		},
		"invalid query - wrong rollup name": {
			[]string{
				"invalid-rollup",
				"1",
			},
			false,
		},
		"invalid query - wrong rollup height": {
			[]string{
				rollupName,
				"1000",
			},
			false,
		},
	}

	for name, tc := range tcs {
		s.Run(name, func() {
			cmd := dacli.CmdMappedBatch()
			out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, cmd, tc.args)
			if tc.valid {
				s.Require().NoError(err)
				var res types.QueryMappedBatchResponse
				s.Require().NoError(val.ClientCtx.JSONCodec.UnmarshalJSON(out.Bytes(), &res))
			} else {
				s.Require().Error(err)
			}
		})
	}
}
