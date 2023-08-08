package tx

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Finschia/finschia-rdk/codec"
	codectypes "github.com/Finschia/finschia-rdk/codec/types"
	"github.com/Finschia/finschia-rdk/std"
	"github.com/Finschia/finschia-rdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-rdk/types"
	"github.com/Finschia/finschia-rdk/x/auth/testutil"
)

func TestGenerator(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &testdata.TestMsg{})
	protoCodec := codec.NewProtoCodec(interfaceRegistry)
	suite.Run(t, testutil.NewTxConfigTestSuite(NewTxConfig(protoCodec, DefaultSignModes)))
}
