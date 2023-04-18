package tx2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/testutil/network"
	"github.com/Finschia/finschia-sdk/testutil/rest"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/query"
	"github.com/Finschia/finschia-sdk/types/tx2"
	bankcli "github.com/Finschia/finschia-sdk/x/bank/client/testutil"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	txHeight    int64
	queryClient tx2.ServiceClient
	txRes       sdk.TxResponse
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)
	s.Require().NotNil(s.network)

	val := s.network.Validators[0]

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.queryClient = tx2.NewServiceClient(val.ClientCtx)

	// Create a new MsgSend tx from val to itself.
	out, err := bankcli.MsgSendExec(
		val.ClientCtx,
		val.Address,
		val.Address,
		sdk.NewCoins(
			sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10)),
		),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		fmt.Sprintf("--gas=%d", flags.DefaultGasLimit),
		fmt.Sprintf("--%s=foobar", flags.FlagNote),
	)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &s.txRes))
	s.Require().Equal(uint32(0), s.txRes.Code)

	out, err = bankcli.MsgSendExec(
		val.ClientCtx,
		val.Address,
		val.Address,
		sdk.NewCoins(
			sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(1)),
		),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=2", flags.FlagSequence),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		fmt.Sprintf("--gas=%d", flags.DefaultGasLimit),
		fmt.Sprintf("--%s=foobar", flags.FlagNote),
	)
	s.Require().NoError(err)
	var tr sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &tr))
	s.Require().Equal(uint32(0), tr.Code)

	s.Require().NoError(s.network.WaitForNextBlock())
	height, err := s.network.LatestHeight()
	s.Require().NoError(err)
	s.txHeight = height
	fmt.Printf("s.txHeight: %d\n", height)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s IntegrationTestSuite) TestGetBlockWithTxs_GRPC() {
	testCases := []struct {
		name      string
		req       *tx2.GetBlockWithTxsRequest
		expErr    bool
		expErrMsg string
	}{
		{"nil request", nil, true, "request cannot be nil"},
		{"empty request", &tx2.GetBlockWithTxsRequest{}, true, "height must not be less than 1 or greater than the current height"},
		{"bad height", &tx2.GetBlockWithTxsRequest{Height: 99999999}, true, "height must not be less than 1 or greater than the current height"},
		{"bad pagination", &tx2.GetBlockWithTxsRequest{Height: s.txHeight, Pagination: &query.PageRequest{Offset: 1000, Limit: 100}}, true, "out of range"},
		{"good request", &tx2.GetBlockWithTxsRequest{Height: s.txHeight}, false, ""},
		{"with pagination request", &tx2.GetBlockWithTxsRequest{Height: s.txHeight, Pagination: &query.PageRequest{Offset: 0, Limit: 1}}, false, ""},
		{"page all request", &tx2.GetBlockWithTxsRequest{Height: s.txHeight, Pagination: &query.PageRequest{Offset: 0, Limit: 100}}, false, ""},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Query the tx via gRPC.
			grpcRes, err := s.queryClient.GetBlockWithTxs(context.Background(), tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				s.Require().NoError(err)
				s.Require().Equal("foobar", grpcRes.Txs[0].Body.Memo)
				s.Require().Equal(grpcRes.Block.Header.Height, tc.req.Height)
				if tc.req.Pagination != nil {
					s.Require().LessOrEqual(len(grpcRes.Txs), int(tc.req.Pagination.Limit))
				}
			}
		})
	}
}

func (s IntegrationTestSuite) TestGetBlockWithTxs_GRPCGateway() {
	val := s.network.Validators[0]
	testCases := []struct {
		name      string
		url       string
		expErr    bool
		expErrMsg string
	}{
		{
			"empty params",
			fmt.Sprintf("%s/lbm/tx/v1beta1/txs/block/0", val.APIAddress),
			true, "height must not be less than 1 or greater than the current height",
		},
		{
			"bad height",
			fmt.Sprintf("%s/lbm/tx/v1beta1/txs/block/%d", val.APIAddress, 9999999),
			true, "height must not be less than 1 or greater than the current height",
		},
		{
			"good request",
			fmt.Sprintf("%s/lbm/tx/v1beta1/txs/block/%d", val.APIAddress, s.txHeight),
			false, "",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			res, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)
			if tc.expErr {
				s.Require().Contains(string(res), tc.expErrMsg)
			} else {
				var result tx2.GetBlockWithTxsResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(res, &result)
				s.Require().NoError(err)
				s.Require().Equal("foobar", result.Txs[0].Body.Memo)
				s.Require().Equal(result.Block.Header.Height, s.txHeight)
			}
		})
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
