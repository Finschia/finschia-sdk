package ante_test

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/codec/types"
	"github.com/line/lbm-sdk/testutil/testdata"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/auth/ante"
	"github.com/line/lbm-sdk/x/auth/tx"
)

type setFeeGranter interface {
	SetFeeGranter(feeGranter sdk.AccAddress)
}

func (suite *AnteTestSuite) TestRejectFeeGranter() {
	suite.SetupTest(true) // setup
	txConfig := tx.NewTxConfig(codec.NewProtoCodec(types.NewInterfaceRegistry()), tx.DefaultSignModes)
	txBuilder := txConfig.NewTxBuilder()
	d := ante.NewRejectFeeGranterDecorator()
	antehandler := sdk.ChainAnteDecorators(d)

	_, err := antehandler(suite.ctx, txBuilder.GetTx(), false)
	suite.Require().NoError(err)

	setGranterTx := txBuilder.(setFeeGranter)
	_, _, addr := testdata.KeyTestPubAddr()
	setGranterTx.SetFeeGranter(addr)

	_, err = antehandler(suite.ctx, txBuilder.GetTx(), false)
	suite.Require().Error(err)
}
