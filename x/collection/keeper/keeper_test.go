package keeper_test

import (
	"fmt"
	"context"
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"

	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
	"github.com/line/lbm-sdk/x/collection/keeper"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx         sdk.Context
	goCtx       context.Context
	keeper      keeper.Keeper
	queryServer collection.QueryServer
	msgServer   collection.MsgServer

	vendor   sdk.AccAddress
	operator sdk.AccAddress
	customer sdk.AccAddress

	contractID string
	fungibleTokenClassID string
	nonFungibleTokenClassID string

	balance sdk.Int
}

func (s *KeeperTestSuite) SetupTest() {
	checkTx := false
	app := simapp.Setup(checkTx)
	s.ctx = app.BaseApp.NewContext(checkTx, ocproto.Header{})
	s.goCtx = sdk.WrapSDKContext(s.ctx)
	s.keeper = app.CollectionKeeper

	// s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	createAddress := func() sdk.AccAddress {
		return sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}
	s.vendor = createAddress()
	s.operator = createAddress()
	s.customer = createAddress()

	s.balance = sdk.OneInt()

	// create a contract
	s.contractID = "f00dbabe"

	// create token classes
	s.fungibleTokenClassID = fmt.Sprintf("%08d", 1)
	s.nonFungibleTokenClassID = "10000001"

	// set balances
	// TODO: replace with mint
	for _, to := range []sdk.AccAddress{s.vendor, s.operator, s.customer} {
		err := s.keeper.SetBalance(s.ctx, s.contractID, to, s.fungibleTokenClassID + fmt.Sprintf("%08d", 0), s.balance)
		s.Require().NoError(err)
	}
	err := s.keeper.SetBalance(s.ctx, s.contractID, s.vendor, s.nonFungibleTokenClassID + fmt.Sprintf("%08d", 1), s.balance)
	s.Require().NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
