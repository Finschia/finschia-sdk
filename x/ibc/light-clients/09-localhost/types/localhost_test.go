package types_test

import (
	"testing"

	ostproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/simapp"
	sdk "github.com/line/lbm-sdk/v2/types"
	clienttypes "github.com/line/lbm-sdk/v2/x/ibc/core/02-client/types"
	"github.com/line/lbm-sdk/v2/x/ibc/core/exported"
)

const (
	height = 4
)

var (
	clientHeight = clienttypes.NewHeight(0, 10)
)

type LocalhostTestSuite struct {
	suite.Suite

	cdc   codec.Marshaler
	ctx   sdk.Context
	store sdk.KVStore
}

func (suite *LocalhostTestSuite) SetupTest() {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	suite.cdc = app.AppCodec()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, ostproto.Header{Height: 1, ChainID: "ibc-chain"})
	suite.store = app.IBCKeeper.ClientKeeper.ClientStore(suite.ctx, exported.Localhost)
}

func TestLocalhostTestSuite(t *testing.T) {
	suite.Run(t, new(LocalhostTestSuite))
}
