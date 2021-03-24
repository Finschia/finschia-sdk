package legacytx_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/v2/codec"
	cryptoAmino "github.com/line/lbm-sdk/v2/crypto/codec"
	"github.com/line/lbm-sdk/v2/testutil/testdata"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/auth/legacy/legacytx"
	"github.com/line/lbm-sdk/v2/x/auth/testutil"
)

func testCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)
	cryptoAmino.RegisterCrypto(cdc)
	cdc.RegisterConcrete(&testdata.TestMsg{}, "lbm-sdk/Test", nil)
	return cdc
}

func TestStdTxConfig(t *testing.T) {
	cdc := testCodec()
	txGen := legacytx.StdTxConfig{Cdc: cdc}
	suite.Run(t, testutil.NewTxConfigTestSuite(txGen))
}
