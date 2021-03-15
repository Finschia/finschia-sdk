package tx

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/v2/codec"
	codectypes "github.com/line/lbm-sdk/v2/codec/types"
	"github.com/line/lbm-sdk/v2/std"
	"github.com/line/lbm-sdk/v2/testutil/testdata"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/auth/testutil"
)

func TestGenerator(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &testdata.TestMsg{})
	protoCodec := codec.NewProtoCodec(interfaceRegistry)
	suite.Run(t, testutil.NewTxConfigTestSuite(NewTxConfig(protoCodec, DefaultSignModes)))
}
