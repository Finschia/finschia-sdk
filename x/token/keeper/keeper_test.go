package keeper_test

import (
	"context"
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/keeper"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx sdk.Context
	goCtx context.Context
	keeper keeper.Keeper
	queryServer token.QueryServer
	msgServer token.MsgServer

	vendor sdk.AccAddress
	operator sdk.AccAddress
	customer sdk.AccAddress

	classID string

	balance sdk.Int
}

func (s *KeeperTestSuite) SetupTest() {
	checkTx := false
	app := simapp.Setup(checkTx)
	s.ctx = app.BaseApp.NewContext(checkTx, ocproto.Header{})
	s.goCtx = sdk.WrapSDKContext(s.ctx)
	s.keeper = app.TokenKeeper

	s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	s.vendor = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	s.operator = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	s.customer = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())

	s.balance = sdk.NewInt(1000000)

	// create a mintable class
	s.classID = "foodbabe"
	class := token.Token{
		Id: s.classID,
		Name: "Mintable",
		Symbol: "OK",
		Mintable: true,
	}
	err := s.keeper.Issue(s.ctx, class, s.vendor, s.vendor, s.balance)
	s.Require().NoError(err)
	amount := []token.FT{{ClassId: s.classID, Amount: s.balance}}
	err = s.keeper.Burn(s.ctx, s.vendor, amount)
	s.Require().NoError(err)

	// create another class for the query test
	class.Id = "deadbeef"
	err = s.keeper.Issue(s.ctx, class, s.vendor, s.vendor, s.balance)
	s.Require().NoError(err)

	// mint to the others
	for _, to := range []sdk.AccAddress{s.vendor, s.operator, s.customer} {
		err = s.keeper.Mint(s.ctx, s.vendor, to, amount)
		s.Require().NoError(err)
	}

	// grant operator
	for _, action := range []string{token.ActionMint, token.ActionBurn} {
		err = s.keeper.Grant(s.ctx, s.vendor, s.operator, s.classID, action)
		s.Require().NoError(err)
	}

	// approve operator
	err = s.keeper.Approve(s.ctx, s.customer, s.operator, s.classID)
	s.Require().NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
