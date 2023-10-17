package keeper_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection/keeper"
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
	stranger sdk.AccAddress

	contractID string
	ftClassID  string
	nftClassID string

	balance sdk.Int

	depthLimit int

	numNFTs  int
	numRoots int
}

func (s *KeeperTestSuite) createRandomAccounts(accNum int) []sdk.AccAddress {
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

func (s *KeeperTestSuite) SetupTest() {
	checkTx := false
	app := simapp.Setup(checkTx)
	s.ctx = app.BaseApp.NewContext(checkTx, tmproto.Header{})
	s.goCtx = sdk.WrapSDKContext(s.ctx)
	s.keeper = app.CollectionKeeper

	s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	s.depthLimit = 4
	s.keeper.SetParams(s.ctx, collection.Params{
		DepthLimit: uint32(s.depthLimit),
		WidthLimit: 4,
	})

	addresses := []*sdk.AccAddress{
		&s.vendor,
		&s.operator,
		&s.customer,
		&s.stranger,
	}
	for i, address := range s.createRandomAccounts(len(addresses)) {
		*addresses[i] = address
	}

	s.balance = sdk.NewInt(1000000)

	// create a contract
	s.contractID = s.keeper.CreateContract(s.ctx, s.vendor, collection.Contract{
		Name: "fox",
	})

	for _, permission := range []collection.Permission{
		collection.PermissionMint,
		collection.PermissionBurn,
	} {
		s.keeper.Grant(s.ctx, s.contractID, s.vendor, s.operator, permission)
	}

	// create a fungible token class
	ftClassID, err := s.keeper.CreateTokenClass(s.ctx, s.contractID, &collection.FTClass{
		Name:     "tibetian fox",
		Mintable: true,
	})
	s.Require().NoError(err)
	s.ftClassID = *ftClassID

	// create a non-fungible token class
	nftClassID, err := s.keeper.CreateTokenClass(s.ctx, s.contractID, &collection.NFTClass{
		Name: "fennec fox",
	})
	s.Require().NoError(err)
	s.nftClassID = *nftClassID

	// mint & burn fts
	for _, to := range []sdk.AccAddress{s.customer, s.operator, s.vendor} {
		tokenID := collection.NewFTID(s.ftClassID)
		amount := collection.NewCoins(collection.NewCoin(tokenID, s.balance))

		err := s.keeper.MintFT(s.ctx, s.contractID, to, amount)
		s.Require().NoError(err)

		_, err = s.keeper.BurnCoins(s.ctx, s.contractID, to, amount)
		s.Require().NoError(err)
		err = s.keeper.MintFT(s.ctx, s.contractID, to, amount)
		s.Require().NoError(err)
	}

	// mint nfts
	newParams := func(classID string, size int) []collection.MintNFTParam {
		res := make([]collection.MintNFTParam, size)
		for i := range res {
			res[i] = collection.MintNFTParam{
				TokenType: s.nftClassID,
			}
		}
		return res
	}
	// 1 for the successful attach, 2 for the failure
	remainders := 1 + 2
	s.numNFTs = s.depthLimit + remainders
	// 3 chains, and each chain has depth_limit, 1 and 2 of its length.
	s.numRoots = 3
	for _, to := range []sdk.AccAddress{s.customer, s.operator, s.vendor} {
		tokens, err := s.keeper.MintNFT(s.ctx, s.contractID, to, newParams(s.nftClassID, s.depthLimit))
		s.Require().NoError(err)

		// create a chain of its length depth_limit
		for i := range tokens[1:] {
			r := len(tokens) - 1 - i
			subject := tokens[r].TokenId
			target := tokens[r-1].TokenId
			err := s.keeper.Attach(s.ctx, s.contractID, to, subject, target)
			s.Require().NoError(err)
		}

		tokens, err = s.keeper.MintNFT(s.ctx, s.contractID, to, newParams(s.nftClassID, remainders))
		s.Require().NoError(err)

		// a chain of length 2
		err = s.keeper.Attach(s.ctx, s.contractID, to, tokens[remainders-1].TokenId, tokens[remainders-2].TokenId)
		s.Require().NoError(err)

	}

	// authorize
	err = s.keeper.AuthorizeOperator(s.ctx, s.contractID, s.customer, s.operator)
	s.Require().NoError(err)
	err = s.keeper.AuthorizeOperator(s.ctx, s.contractID, s.customer, s.stranger)
	s.Require().NoError(err)

	// not token contract
	notTokenContractID := app.ClassKeeper.NewID(s.ctx)
	err = keeper.ValidateLegacyContract(s.keeper, s.ctx, notTokenContractID)
	s.Require().ErrorIs(err, collection.ErrCollectionNotExist)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
