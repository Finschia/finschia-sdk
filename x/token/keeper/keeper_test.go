package keeper_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/token"
	"github.com/Finschia/finschia-sdk/x/token/keeper"
)

type KeeperTestSuite struct {
	suite.Suite

	deterministic bool

	ctx         sdk.Context
	goCtx       context.Context
	keeper      keeper.Keeper
	queryServer token.QueryServer
	msgServer   token.MsgServer

	vendor   sdk.AccAddress
	operator sdk.AccAddress
	customer sdk.AccAddress
	stranger sdk.AccAddress

	contractID           string
	unmintableContractId string

	balance sdk.Int
}

func createRandomAccounts(accNum int, deterministic bool) []sdk.AccAddress {
	if deterministic {
		addresses := make([]sdk.AccAddress, accNum)
		for i := range addresses {
			addresses[i] = sdk.AccAddress(fmt.Sprintf("address%d", i))
		}
		return addresses
	} else {
		seenAddresses := make(map[string]bool, accNum)
		addresses := make([]sdk.AccAddress, accNum)
		for i := range addresses {
			var address sdk.AccAddress
			for {
				pk := secp256k1.GenPrivKey().PubKey()
				address = sdk.AccAddress(pk.Address())
				if !seenAddresses[address.String()] {
					seenAddresses[address.String()] = true
					break
				}
			}
			addresses[i] = address
		}
		return addresses
	}
}

func (s *KeeperTestSuite) SetupTest() {
	checkTx := false
	app := simapp.Setup(checkTx)
	s.ctx = app.BaseApp.NewContext(checkTx, tmproto.Header{})
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
	for i, address := range createRandomAccounts(len(addresses), s.deterministic) {
		*addresses[i] = address
	}

	s.balance = sdk.NewInt(1000)

	// create a mintable class
	class := token.Contract{
		Name:     "Mintable",
		Symbol:   "OK",
		Mintable: true,
	}
	s.contractID = s.keeper.Issue(s.ctx, class, s.vendor, s.vendor, s.balance)

	err := s.keeper.Burn(s.ctx, s.contractID, s.vendor, s.balance)
	s.Require().NoError(err)

	// create another class for the query test
	s.keeper.Issue(s.ctx, class, s.vendor, s.vendor, s.balance)

	// mint to the others
	for _, to := range []sdk.AccAddress{s.vendor, s.operator, s.customer} {
		err := s.keeper.Mint(s.ctx, s.contractID, s.vendor, to, s.balance)
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
		err := s.keeper.AuthorizeOperator(s.ctx, s.contractID, holder, s.operator)
		s.Require().NoError(err)
	}

	// not token contract
	notTokenContractID := app.ClassKeeper.NewID(s.ctx)
	err = keeper.ValidateLegacyContract(s.keeper, s.ctx, notTokenContractID)
	s.Require().ErrorIs(err, token.ErrTokenNotExist)

	s.unmintableContractId = s.keeper.Issue(s.ctx, token.Contract{
		Name:     "Unmintable",
		Symbol:   "UMT",
		Mintable: false,
	}, s.vendor, s.vendor, s.balance)
}

func TestKeeperTestSuite(t *testing.T) {
	for _, deterministic := range []bool{
		false,
		true,
	} {
		suite.Run(t, &KeeperTestSuite{deterministic: deterministic})
	}
}
