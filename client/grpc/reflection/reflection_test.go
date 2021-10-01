package reflection_test

import (
	"context"
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/client/grpc/reflection"
	"github.com/line/lbm-sdk/simapp"
)

type IntegrationTestSuite struct {
	suite.Suite

	queryClient reflection.ReflectionServiceClient
}

func (s *IntegrationTestSuite) SetupSuite() {
	app := simapp.Setup(false)

	sdkCtx := app.BaseApp.NewContext(false, ocproto.Header{})
	queryHelper := baseapp.NewQueryServerTestHelper(sdkCtx, app.InterfaceRegistry())
	queryClient := reflection.NewReflectionServiceClient(queryHelper)
	s.queryClient = queryClient
}

func (s IntegrationTestSuite) TestSimulateService() {
	// We will test the following interface for testing.
	var iface = "lbm.evidence.v1.Evidence"

	// Test that "lbm.evidence.v1.Evidence" is included in the
	// interfaces.
	resIface, err := s.queryClient.ListAllInterfaces(
		context.Background(),
		&reflection.ListAllInterfacesRequest{},
	)
	s.Require().NoError(err)
	s.Require().Contains(resIface.GetInterfaceNames(), iface)

	// Test that "lbm.evidence.v1.Evidence" has at least the
	// Equivocation implementations.
	resImpl, err := s.queryClient.ListImplementations(
		context.Background(),
		&reflection.ListImplementationsRequest{InterfaceName: iface},
	)
	s.Require().NoError(err)
	s.Require().Contains(resImpl.GetImplementationMessageNames(), "/lbm.evidence.v1.Equivocation")
}

func TestSimulateTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
