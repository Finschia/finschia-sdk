package rpc_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	ctypes "github.com/line/ostracon/rpc/core/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/client/rpc"
	"github.com/line/lbm-sdk/codec/legacy"
	clitestutil "github.com/line/lbm-sdk/testutil/cli"
	"github.com/line/lbm-sdk/testutil/network"
	"github.com/line/lbm-sdk/types/rest"
)

type IntegrationTestSuite struct {
	suite.Suite

	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = network.New(s.T(), network.DefaultConfig())
	s.Require().NotNil(s.network)

	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestStatusCommand() {
	val0 := s.network.Validators[0]
	cmd := rpc.StatusCommand()

	out, err := clitestutil.ExecTestCLICmd(val0.ClientCtx, cmd, []string{})
	s.Require().NoError(err)

	// Make sure the output has the validator moniker.
	s.Require().Contains(out.String(), fmt.Sprintf("\"moniker\":\"%s\"", val0.Moniker))
}

func (s *IntegrationTestSuite) TestLatestBlocks() {
	val0 := s.network.Validators[0]

	res, err := rest.GetRequest(fmt.Sprintf("%s/blocks/latest", val0.APIAddress))
	s.Require().NoError(err)

	var result ctypes.ResultBlock
	err = legacy.Cdc.UnmarshalJSON(res, &result)
	s.Require().NoError(err)

	{
		var result2 ctypes.ResultBlock
		res, err := rest.GetRequest(fmt.Sprintf("%s/blocks/%d", val0.APIAddress, result.Block.Height))
		s.Require().NoError(err)
		err = legacy.Cdc.UnmarshalJSON(res, &result2)
		s.Require().NoError(err)
		s.Require().Equal(result, result2)
	}
	{
		var result3 ctypes.ResultBlock
		hash64 := base64.URLEncoding.EncodeToString(result.Block.Hash())
		res, err := rest.GetRequest(fmt.Sprintf("%s/block/%s", val0.APIAddress, hash64))
		s.Require().NoError(err)
		err = legacy.Cdc.UnmarshalJSON(res, &result3)
		s.Require().NoError(err)
		s.Require().Equal(result, result3)
	}
}

func (s *IntegrationTestSuite) TestBlockWithFailure() {
	val0 := s.network.Validators[0]

	tcs := []struct {
		name   string
		height string
		errRes string
	}{
		{
			name:   "parse error",
			height: "a",
			errRes: "{\"error\":\"couldn't parse block height. Assumed format is '/blocks/{height}'.\"}",
		},
		{
			name:   "bigger height error",
			height: "1234567890",
			errRes: "{\"error\":\"requested block height is bigger then the chain length\"}",
		},
	}
	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			res, err := rest.GetRequest(fmt.Sprintf("%s/blocks/%s", val0.APIAddress, tc.height))
			s.Require().NoError(err)
			s.Require().Equal(tc.errRes, string(res))
		})
	}
}

func (s *IntegrationTestSuite) TestBlockByHashWithFailure() {
	val0 := s.network.Validators[0]

	tcs := []struct {
		name   string
		hash   string
		errRes string
	}{
		{
			name:   "base64 error",
			hash:   "wrong hash",
			errRes: "{\"error\":\"couldn't decode block hash by Base64URLDecode.\"}",
		},
		{
			name:   "size error",
			hash:   base64.URLEncoding.EncodeToString([]byte{0}),
			errRes: "{\"error\":\"the length of block hash must be 32\"}",
		},
	}
	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			res, err := rest.GetRequest(fmt.Sprintf("%s/block/%s", val0.APIAddress, tc.hash))
			s.Require().NoError(err)
			s.Require().Equal(tc.errRes, string(res))
		})
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
