package rpc_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Finschia/finschia-sdk/client/rpc"
	clitestutil "github.com/Finschia/finschia-sdk/testutil/cli"
	"github.com/Finschia/finschia-sdk/testutil/network"
)

type IntegrationTestSuite struct {
	suite.Suite

	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")
	s.T().Skipf("🔬 To work with minimal modifications.")

	s.network = network.New(s.T(), network.DefaultConfig())
	s.Require().NotNil(s.network)

	s.Require().NoError(s.network.WaitForNextBlock())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestStatusCommand() {
	s.T().Skipf("🔬 The rollkit/cosmos-sdk also remains faulty.")
	val0 := s.network.Validators[0]
	cmd := rpc.StatusCommand()

	out, err := clitestutil.ExecTestCLICmd(val0.ClientCtx, cmd, []string{})
	s.Require().NoError(err)

	// Make sure the output has the validator moniker.
	s.Require().Contains(out.String(), fmt.Sprintf("\"moniker\":\"%s\"", val0.Moniker))
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
