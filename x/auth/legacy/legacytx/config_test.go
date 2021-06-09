package legacytx_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/line/lfb-sdk/codec"
	cryptoAmino "github.com/line/lfb-sdk/crypto/codec"
	"github.com/line/lfb-sdk/testutil/testdata"
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/auth/legacy/legacytx"
	"github.com/line/lfb-sdk/x/auth/testutil"
)

func testCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)
	cryptoAmino.RegisterCrypto(cdc)
	cdc.RegisterConcrete(&testdata.TestMsg{}, "lfb-sdk/Test", nil)
	return cdc
}

func TestStdTxConfig(t *testing.T) {
	cdc := testCodec()
	txGen := legacytx.StdTxConfig{Cdc: cdc}
	suite.Run(t, testutil.NewTxConfigTestSuite(txGen))
}
