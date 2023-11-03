package ocservice_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Finschia/ostracon/libs/bytes"

	"github.com/Finschia/finschia-sdk/client/grpc/ocservice"
	codectypes "github.com/Finschia/finschia-sdk/codec/types"
	cryptotypes "github.com/Finschia/finschia-sdk/crypto/types"
	"github.com/Finschia/finschia-sdk/testutil/network"
	"github.com/Finschia/finschia-sdk/testutil/rest"
	qtypes "github.com/Finschia/finschia-sdk/types/query"
	"github.com/Finschia/finschia-sdk/version"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	queryClient ocservice.ServiceClient
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	s.Require().NotNil(s.network)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.queryClient = ocservice.NewServiceClient(s.network.Validators[0].ClientCtx)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s IntegrationTestSuite) TestQueryNodeInfo() {
	val := s.network.Validators[0]

	res, err := s.queryClient.GetNodeInfo(context.Background(), &ocservice.GetNodeInfoRequest{})
	s.Require().NoError(err)
	s.Require().Equal(res.ApplicationVersion.AppName, version.NewInfo().AppName)

	restRes, err := rest.GetRequest(fmt.Sprintf("%s/lbm/base/ostracon/v1/node_info", val.APIAddress))
	s.Require().NoError(err)
	var getInfoRes ocservice.GetNodeInfoResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(restRes, &getInfoRes))
	s.Require().Equal(getInfoRes.ApplicationVersion.AppName, version.NewInfo().AppName)
}

func (s IntegrationTestSuite) TestQuerySyncing() {
	val := s.network.Validators[0]

	_, err := s.queryClient.GetSyncing(context.Background(), &ocservice.GetSyncingRequest{})
	s.Require().NoError(err)

	restRes, err := rest.GetRequest(fmt.Sprintf("%s/lbm/base/ostracon/v1/syncing", val.APIAddress))
	s.Require().NoError(err)
	var syncingRes ocservice.GetSyncingResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(restRes, &syncingRes))
}

func (s IntegrationTestSuite) TestQueryLatestBlock() {
	val := s.network.Validators[0]
	_, err := s.queryClient.GetLatestBlock(context.Background(), &ocservice.GetLatestBlockRequest{})
	s.Require().NoError(err)

	restRes, err := rest.GetRequest(fmt.Sprintf("%s/lbm/base/ostracon/v1/blocks/latest", val.APIAddress))
	s.Require().NoError(err)
	var blockInfoRes ocservice.GetLatestBlockResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(restRes, &blockInfoRes))
}

func (s IntegrationTestSuite) TestQueryBlockByHash() {
	val := s.network.Validators[0]
	node, _ := val.ClientCtx.GetNode()
	blk, _ := node.Block(context.Background(), nil)
	blkhash := blk.BlockID.Hash

	tcs := []struct {
		hash  bytes.HexBytes
		isErr bool
		err   string
	}{
		{blkhash, false, ""},
		{bytes.HexBytes("wrong hash"), true, "the length of block hash must be 32: invalid request"},
		{bytes.HexBytes(""), true, "block hash cannot be empty"},
	}

	for _, tc := range tcs {
		_, err := s.queryClient.GetBlockByHash(context.Background(), &ocservice.GetBlockByHashRequest{Hash: tc.hash})
		if tc.isErr {
			s.Require().Error(err)
			s.Require().Contains(err.Error(), tc.err)
		} else {
			s.Require().NoError(err)
		}
	}

	restRes, err := rest.GetRequest(fmt.Sprintf("%s/lbm/base/ostracon/v1/block/%s", val.APIAddress, base64.URLEncoding.EncodeToString(blkhash)))
	s.Require().NoError(err)
	var blockInfoRes ocservice.GetBlockByHashResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(restRes, &blockInfoRes))
	blockId := blockInfoRes.GetBlockId()
	s.Require().Equal(blkhash, bytes.HexBytes(blockId.Hash))

	block := blockInfoRes.GetBlock()
	s.Require().Equal(val.ClientCtx.ChainID, block.Header.ChainID)
}

func (s IntegrationTestSuite) TestGetLatestBlock() {
	val := s.network.Validators[0]
	latestHeight, err := s.network.LatestHeight()
	s.Require().NoError(err)

	blockRes, err := s.queryClient.GetLatestBlock(context.Background(), &ocservice.GetLatestBlockRequest{})
	s.Require().NoError(err)
	s.Require().NotNil(blockRes)
	s.Require().Equal(latestHeight, blockRes.Block.Header.Height)

	restRes, err := rest.GetRequest(fmt.Sprintf("%s/lbm/base/ostracon/v1/blocks/latest", val.APIAddress))
	s.Require().NoError(err)
	var blockInfoRes ocservice.GetLatestBlockResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(restRes, &blockInfoRes))
	s.Require().Equal(blockRes.Block.Header.Height, blockInfoRes.Block.Header.Height)
}

func (s IntegrationTestSuite) TestQueryBlockByHeight() {
	val := s.network.Validators[0]
	_, err := s.queryClient.GetBlockByHeight(context.Background(), &ocservice.GetBlockByHeightRequest{})
	s.Require().Error(err)

	_, err = s.queryClient.GetBlockByHeight(context.Background(), &ocservice.GetBlockByHeightRequest{Height: 1})
	s.Require().NoError(err)

	restRes, err := rest.GetRequest(fmt.Sprintf("%s/lbm/base/ostracon/v1/blocks/%d", val.APIAddress, 1))
	s.Require().NoError(err)
	var blockInfoRes ocservice.GetBlockByHeightResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(restRes, &blockInfoRes))

	block := blockInfoRes.GetBlock()
	s.Require().Equal(int64(1), block.Header.Height)
}

func (s IntegrationTestSuite) TestQueryBlockResultsByHeight() {
	val := s.network.Validators[0]
	_, err := s.queryClient.GetBlockResultsByHeight(context.Background(), &ocservice.GetBlockResultsByHeightRequest{Height: 1})
	s.Require().NoError(err)

	restRes, err := rest.GetRequest(fmt.Sprintf("%s/lbm/base/ostracon/v1/blockresults/%d", val.APIAddress, 1))
	s.Require().NoError(err)
	var blockResultsRes ocservice.GetBlockResultsByHeightResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(restRes, &blockResultsRes))

	txResult := blockResultsRes.GetTxsResults()
	s.Require().Equal(0, len(txResult))

	beginBlock := blockResultsRes.GetResBeginBlock()
	s.Require().Equal(6, len(beginBlock.Events)) // coinbase event (2) + transfer mintModule to feeCollectorName(3) + mint event

	endBlock := blockResultsRes.GetResEndBlock()
	s.Require().Equal(0, len(endBlock.Events))
}

func (s IntegrationTestSuite) TestQueryLatestValidatorSet() {
	val := s.network.Validators[0]

	// nil pagination
	res, err := s.queryClient.GetLatestValidatorSet(context.Background(), &ocservice.GetLatestValidatorSetRequest{
		Pagination: nil,
	})
	s.Require().NoError(err)
	s.Require().Equal(1, len(res.Validators))
	content, ok := res.Validators[0].PubKey.GetCachedValue().(cryptotypes.PubKey)
	s.Require().Equal(true, ok)
	s.Require().Equal(content, val.PubKey)

	// with pagination
	_, err = s.queryClient.GetLatestValidatorSet(context.Background(), &ocservice.GetLatestValidatorSetRequest{Pagination: &qtypes.PageRequest{
		Offset: 0,
		Limit:  10,
	}})
	s.Require().NoError(err)

	// rest request without pagination
	_, err = rest.GetRequest(fmt.Sprintf("%s/lbm/base/ostracon/v1/validatorsets/latest", val.APIAddress))
	s.Require().NoError(err)

	// rest request with pagination
	restRes, err := rest.GetRequest(fmt.Sprintf("%s/lbm/base/ostracon/v1/validatorsets/latest?pagination.offset=%d&pagination.limit=%d", val.APIAddress, 0, 1))
	s.Require().NoError(err)
	var validatorSetRes ocservice.GetLatestValidatorSetResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(restRes, &validatorSetRes))
	s.Require().Equal(1, len(validatorSetRes.Validators))
	anyPub, err := codectypes.NewAnyWithValue(val.PubKey)
	s.Require().NoError(err)
	s.Require().Equal(validatorSetRes.Validators[0].PubKey, anyPub)
}

func (s IntegrationTestSuite) TestLatestValidatorSet_GRPC() {
	vals := s.network.Validators
	testCases := []struct {
		name      string
		req       *ocservice.GetLatestValidatorSetRequest
		expErr    bool
		expErrMsg string
	}{
		{"nil request", nil, true, "cannot be nil"},
		{"no pagination", &ocservice.GetLatestValidatorSetRequest{}, false, ""},
		{"with pagination", &ocservice.GetLatestValidatorSetRequest{Pagination: &qtypes.PageRequest{Offset: 0, Limit: uint64(len(vals))}}, false, ""},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			grpcRes, err := s.queryClient.GetLatestValidatorSet(context.Background(), tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				s.Require().NoError(err)
				s.Require().Len(grpcRes.Validators, len(vals))
				s.Require().Equal(grpcRes.Pagination.Total, uint64(len(vals)))
				content, ok := grpcRes.Validators[0].PubKey.GetCachedValue().(cryptotypes.PubKey)
				s.Require().Equal(true, ok)
				s.Require().Equal(content, vals[0].PubKey)
			}
		})
	}
}

func (s IntegrationTestSuite) TestLatestValidatorSet_GRPCGateway() {
	vals := s.network.Validators
	testCases := []struct {
		name      string
		url       string
		expErr    bool
		expErrMsg string
	}{
		{"no pagination", fmt.Sprintf("%s/lbm/base/ostracon/v1/validatorsets/latest", vals[0].APIAddress), false, ""},
		{"pagination invalid fields", fmt.Sprintf("%s/lbm/base/ostracon/v1/validatorsets/latest?pagination.offset=-1&pagination.limit=-2", vals[0].APIAddress), true, "strconv.ParseUint"},
		{"with pagination", fmt.Sprintf("%s/lbm/base/ostracon/v1/validatorsets/latest?pagination.offset=0&pagination.limit=2", vals[0].APIAddress), false, ""},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			res, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)
			if tc.expErr {
				s.Require().Contains(string(res), tc.expErrMsg)
			} else {
				var result ocservice.GetLatestValidatorSetResponse
				err = vals[0].ClientCtx.Codec.UnmarshalJSON(res, &result)
				s.Require().NoError(err)
				s.Require().Equal(uint64(len(vals)), result.Pagination.Total)
				anyPub, err := codectypes.NewAnyWithValue(vals[0].PubKey)
				s.Require().NoError(err)
				s.Require().Equal(result.Validators[0].PubKey, anyPub)
			}
		})
	}
}

func (s IntegrationTestSuite) TestValidatorSetByHeight_GRPC() {
	vals := s.network.Validators
	testCases := []struct {
		name      string
		req       *ocservice.GetValidatorSetByHeightRequest
		expErr    bool
		expErrMsg string
	}{
		{"nil request", nil, true, "request cannot be nil"},
		{"empty request", &ocservice.GetValidatorSetByHeightRequest{}, true, "height must be greater than 0"},
		{"no pagination", &ocservice.GetValidatorSetByHeightRequest{Height: 1}, false, ""},
		{"with pagination", &ocservice.GetValidatorSetByHeightRequest{Height: 1, Pagination: &qtypes.PageRequest{Offset: 0, Limit: 1}}, false, ""},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			grpcRes, err := s.queryClient.GetValidatorSetByHeight(context.Background(), tc.req)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expErrMsg)
			} else {
				s.Require().NoError(err)
				s.Require().Len(grpcRes.Validators, len(vals))
				s.Require().Equal(grpcRes.Pagination.Total, uint64(len(vals)))
			}
		})
	}
}

func (s IntegrationTestSuite) TestValidatorSetByHeight_GRPCGateway() {
	vals := s.network.Validators
	testCases := []struct {
		name      string
		url       string
		expErr    bool
		expErrMsg string
	}{
		{"invalid height", fmt.Sprintf("%s/lbm/base/ostracon/v1/validatorsets/%d", vals[0].APIAddress, -1), true, "height must be greater than 0"},
		{"no pagination", fmt.Sprintf("%s/lbm/base/ostracon/v1/validatorsets/%d", vals[0].APIAddress, 1), false, ""},
		{"pagination invalid fields", fmt.Sprintf("%s/lbm/base/ostracon/v1/validatorsets/%d?pagination.offset=-1&pagination.limit=-2", vals[0].APIAddress, 1), true, "strconv.ParseUint"},
		{"with pagination", fmt.Sprintf("%s/lbm/base/ostracon/v1/validatorsets/%d?pagination.offset=0&pagination.limit=2", vals[0].APIAddress, 1), false, ""},
	}
	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			res, err := rest.GetRequest(tc.url)
			s.Require().NoError(err)
			if tc.expErr {
				s.Require().Contains(string(res), tc.expErrMsg)
			} else {
				var result ocservice.GetValidatorSetByHeightResponse
				err = vals[0].ClientCtx.Codec.UnmarshalJSON(res, &result)
				s.Require().NoError(err)
				s.Require().Equal(uint64(len(vals)), result.Pagination.Total)
			}
		})
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
