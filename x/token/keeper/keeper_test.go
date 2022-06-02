package keeper_test

import (
	"context"
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"

	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/keeper"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx         sdk.Context
	goCtx       context.Context
	keeper      keeper.Keeper
	queryServer token.QueryServer
	msgServer   token.MsgServer

	vendor   sdk.AccAddress
	operator sdk.AccAddress
	customer sdk.AccAddress

	contractID string

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

	createAddress := func() sdk.AccAddress {
		return sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}
	s.vendor = createAddress()
	s.operator = createAddress()
	s.customer = createAddress()

	s.balance = sdk.NewInt(1000)

	// create a mintable class
	s.contractID = "f00dbabe"
	class := token.TokenClass{
		ContractId:       s.contractID,
		Name:     "Mintable",
		Symbol:   "OK",
		Mintable: true,
	}
	err := s.keeper.Issue(s.ctx, class, s.vendor, s.vendor, s.balance)
	s.Require().NoError(err)
	err = s.keeper.Burn(s.ctx, s.contractID, s.vendor, s.balance)
	s.Require().NoError(err)

	// create another class for the query test
	class.ContractId = "deadbeef"
	err = s.keeper.Issue(s.ctx, class, s.vendor, s.vendor, s.balance)
	s.Require().NoError(err)

	// mint to the others
	for _, to := range []sdk.AccAddress{s.vendor, s.operator, s.customer} {
		err = s.keeper.Mint(s.ctx, s.contractID, s.vendor, to, s.balance)
		s.Require().NoError(err)
	}

	// grant operator
	for _, permission := range []token.Permission{
		token.Permission_Mint,
		token.Permission_Burn,
	} {
		err = s.keeper.Grant(s.ctx, s.contractID, s.vendor, s.operator, permission)
		s.Require().NoError(err)
	}

	// authorize operator
	for _, holder := range []sdk.AccAddress{s.vendor, s.customer} {
		err = s.keeper.AuthorizeOperator(s.ctx, s.contractID, holder, s.operator)
		s.Require().NoError(err)
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
