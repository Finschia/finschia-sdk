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

	mintableClass string
	notMintableClass string

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

	s.balance = sdk.NewInt(10000)

	// create a mintable class
	s.mintableClass = "foodbabe"
	mintableClass := token.Token{
		Id: s.mintableClass,
		Name: "Mintable",
		Symbol: "OK",
		Mintable: true,
	}
	err := s.keeper.Issue(s.ctx, mintableClass, s.vendor, s.vendor, s.balance)
	s.Require().NoError(err)

	// create a not mintable class
	s.notMintableClass = "fee1dead"
	notMintableClass := token.Token{
		Id: s.notMintableClass,
		Name: "NOT Mintable",
		Symbol: "NO",
		Mintable: false,
	}
	err = s.keeper.Issue(s.ctx, notMintableClass, s.vendor, s.vendor, s.balance)
	s.Require().NoError(err)

	// mint to the others
	for _, to := range []sdk.AccAddress{s.operator, s.customer} {
		amounts := []token.FT{
			{
				ClassId: s.mintableClass,
				Amount: s.balance,
			},
		}
		err = s.keeper.Mint(s.ctx, s.vendor, to, amounts)
		s.Require().NoError(err)
	}

	// grant operator
	for _, class := range []string{s.mintableClass, s.notMintableClass} {
		for _, action := range []string{"mint", "burn", "modify"} {
			err = s.keeper.Grant(s.ctx, s.vendor, s.operator, class, action)
			s.Require().NoError(err)
		}
	}

	// approve operator
	for _, class := range []string{s.mintableClass, s.notMintableClass} {
		err = s.keeper.Approve(s.ctx, s.customer, s.operator, class)
		s.Require().NoError(err)
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
