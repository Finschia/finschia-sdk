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
	stranger sdk.AccAddress

	contractID string

	balance sdk.Int
}

func createRandomAccounts(accNum int) []sdk.AccAddress {
	seenAddresses := make(map[sdk.AccAddress]bool, accNum)
	addresses := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		var address sdk.AccAddress
		for {
			pk := secp256k1.GenPrivKey().PubKey()
			address = sdk.BytesToAccAddress(pk.Address())
			if !seenAddresses[address] {
				seenAddresses[address] = true
				break
			}
		}
		addresses[i] = address
	}
	return addresses
}

func (s *KeeperTestSuite) SetupTest() {
	checkTx := false
	app := simapp.Setup(checkTx)
	s.ctx = app.BaseApp.NewContext(checkTx, ocproto.Header{})
	s.goCtx = sdk.WrapSDKContext(s.ctx)
	s.keeper = app.TokenKeeper

	s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	addresses := []*sdk.AccAddress{
		&s.vendor,
		&s.operator,
		&s.customer,
		&s.stranger,
	}
	for i, address := range createRandomAccounts(len(addresses)) {
		*addresses[i] = address
	}

	s.balance = sdk.NewInt(1000)

	// create a mintable class
	s.contractID = "f00dbabe"
	class := token.TokenClass{
		ContractId: s.contractID,
		Name:       "Mintable",
		Symbol:     "OK",
		Mintable:   true,
	}
	s.keeper.Issue(s.ctx, class, s.vendor, s.vendor, s.balance)
	err := s.keeper.Burn(s.ctx, s.contractID, s.vendor, s.balance)
	s.Require().NoError(err)

	// create another class for the query test
	class.ContractId = "deadbeef"
	s.keeper.Issue(s.ctx, class, s.vendor, s.vendor, s.balance)

	// mint to the others
	for _, to := range []sdk.AccAddress{s.vendor, s.operator, s.customer} {
		err = s.keeper.Mint(s.ctx, s.contractID, s.vendor, to, s.balance)
		s.Require().NoError(err)
	}

	// grant operator
	for _, permission := range []token.Permission{
		token.PermissionMint,
		token.PermissionBurn,
	} {
		s.keeper.Grant(s.ctx, s.contractID, s.vendor, s.operator, permission)
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
