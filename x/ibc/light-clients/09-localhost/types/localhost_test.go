package types_test

import (
	"testing"

	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/line/ostracon/proto/ostracon/types"

	clienttypes "github.com/line/lbm-sdk/x/ibc/core/02-client/types"
	"github.com/line/lbm-sdk/x/ibc/core/exported"
	"github.com/line/lbm-sdk/x/ibc/testing/simapp"
)

const (
	height = 4
)

var (
	clientHeight = clienttypes.NewHeight(0, 10)
)

type LocalhostTestSuite struct {
	suite.Suite

	cdc   codec.Codec
	ctx   sdk.Context
	store sdk.KVStore
}

func (suite *LocalhostTestSuite) SetupTest() {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.AppCodec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{Height: 1, ChainID: "ibc-chain"})
	suite.store = app.IBCKeeper.ClientKeeper.ClientStore(suite.ctx, exported.Localhost)
}

func TestLocalhostTestSuite(t *testing.T) {
	suite.Run(t, new(LocalhostTestSuite))
}
