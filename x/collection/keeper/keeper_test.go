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
	ftClassID string
	nftClassID string

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
	contractID, err := s.keeper.CreateContract(s.ctx, collection.Contract{
		Name: "test contract",
	})
	s.Require().NoError(err)
	s.contractID = *contractID

	// create a fungible token class
	ftClassID, err := s.keeper.CreateTokenClass(s.ctx, &collection.FTClass{
		ContractId: *contractID,
		Name: "test ft class",
	})
	s.Require().NoError(err)
	s.ftClassID = *ftClassID

	// create a non-fungible token class
	nftClassID, err := s.keeper.CreateTokenClass(s.ctx, &collection.NFTClass{
		ContractId: *contractID,
		Name: "test ft class",
	})
	s.Require().NoError(err)
	s.nftClassID = *nftClassID

	// set balances
	// TODO: replace with mint
	for _, to := range []sdk.AccAddress{s.vendor, s.operator, s.customer} {
		s.keeper.SetBalance(s.ctx, s.contractID, to, s.ftClassID + fmt.Sprintf("%08x", 0), s.balance)
	}
	s.keeper.SetBalance(s.ctx, s.contractID, s.customer, s.nftClassID + fmt.Sprintf("%08x", 1), s.balance)

	// authorize
	err = s.keeper.AuthorizeOperator(s.ctx, s.contractID, s.customer, s.operator)
	s.Require().NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
