package tx

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/line/lfb-sdk/codec"
	codectypes "github.com/line/lfb-sdk/codec/types"
	"github.com/line/lfb-sdk/std"
	"github.com/line/lfb-sdk/testutil/testdata"
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/auth/testutil"
)

func TestGenerator(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &testdata.TestMsg{})
	protoCodec := codec.NewProtoCodec(interfaceRegistry)
	suite.Run(t, testutil.NewTxConfigTestSuite(NewTxConfig(protoCodec, DefaultSignModes)))
}
