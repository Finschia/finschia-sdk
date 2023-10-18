package keeper_test

import (
	"fmt"
	"math/rand"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection"
)

func (s *KeeperTestSuite) TestCreateContract() {
	ctx, _ := s.ctx.CacheContext()

	input := collection.Contract{
		Name: "tibetian fox",
		Meta: "Tibetian Fox",
		Uri:  "file:///tibetian_fox.png",
	}
	id := s.keeper.CreateContract(ctx, s.vendor, input)
	s.Require().NotEmpty(id)

	output, err := s.keeper.GetContract(ctx, id)
	s.Require().NoError(err)
	s.Require().NotNil(output)

	s.Require().Equal(id, output.Id)
	s.Require().Equal(input.Name, output.Name)
	s.Require().Equal(input.Meta, output.Meta)
	s.Require().Equal(input.Uri, output.Uri)
}

func (s *KeeperTestSuite) TestCreateTokenClass() {
	testCases := map[string]struct {
		contractID string
		class      collection.TokenClass
		err        error
	}{
		"valid fungible token class": {
			contractID: s.contractID,
			class:      &collection.FTClass{},
		},
		"valid non-fungible token class": {
			contractID: s.contractID,
			class:      &collection.NFTClass{},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			id, err := s.keeper.CreateTokenClass(ctx, tc.contractID, tc.class)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				s.Require().Nil(id)
				return
			}
			s.Require().NotNil(id)

			class, err := s.keeper.GetTokenClass(ctx, tc.contractID, *id)
			s.Require().NoError(err)
			s.Require().NoError(class.ValidateBasic())
		})
	}
}

func (s *KeeperTestSuite) TestMintFT() {
	testCases := map[string]struct {
		contractID string
		amount     collection.Coin
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			amount:     collection.NewFTCoin(s.ftClassID, sdk.OneInt()),
		},
		"invalid token id": {
			contractID: s.contractID,
			amount:     collection.NewNFTCoin(s.ftClassID, 1),
			err:        collection.ErrTokenNotExist,
		},
		"class not found": {
			contractID: s.contractID,
			amount:     collection.NewFTCoin("00bab10c", sdk.OneInt()),
			err:        collection.ErrTokenNotExist,
		},
		"not a class id of ft": {
			contractID: s.contractID,
			amount:     collection.NewFTCoin(s.nftClassID, sdk.OneInt()),
			err:        collection.ErrTokenNotMintable,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			// gather state
			classID := collection.SplitTokenID(tc.amount.TokenId)
			balanceBefore := s.keeper.GetBalance(ctx, tc.contractID, s.stranger, collection.NewFTID(classID))
			supplyBefore := s.keeper.GetSupply(ctx, tc.contractID, classID)
			mintedBefore := s.keeper.GetMinted(ctx, tc.contractID, classID)
			burntBefore := s.keeper.GetBurnt(ctx, tc.contractID, classID)

			err := s.keeper.MintFT(ctx, tc.contractID, s.stranger, collection.NewCoins(tc.amount))
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			amount := tc.amount.Amount
			balanceAfter := s.keeper.GetBalance(ctx, tc.contractID, s.stranger, collection.NewFTID(classID))
			s.Require().Equal(balanceBefore.Add(amount), balanceAfter)
			supplyAfter := s.keeper.GetSupply(ctx, tc.contractID, classID)
			s.Require().Equal(supplyBefore.Add(amount), supplyAfter)
			mintedAfter := s.keeper.GetMinted(ctx, tc.contractID, classID)
			s.Require().Equal(mintedBefore.Add(amount), mintedAfter)
			burntAfter := s.keeper.GetBurnt(ctx, tc.contractID, classID)
			s.Require().Equal(burntBefore, burntAfter)
		})
	}

	// accumulation test
	s.Run("accumulation test", func() {
		ctx, _ := s.ctx.CacheContext()
		numMints := int64(100)
		contractID := s.contractID
		classID := s.ftClassID
		for i := int64(1); i <= numMints; i++ {
			amount := sdk.NewInt(rand.Int63())

			// gather state
			balanceBefore := s.keeper.GetBalance(ctx, s.contractID, s.stranger, collection.NewFTID(s.ftClassID))
			supplyBefore := s.keeper.GetSupply(ctx, contractID, classID)
			mintedBefore := s.keeper.GetMinted(ctx, contractID, classID)
			burntBefore := s.keeper.GetBurnt(ctx, contractID, classID)

			s.keeper.MintFT(ctx, s.contractID, s.stranger, collection.NewCoins(collection.NewFTCoin(s.ftClassID, amount)))

			balanceAfter := s.keeper.GetBalance(ctx, s.contractID, s.stranger, collection.NewFTID(s.ftClassID))
			s.Require().Equal(balanceBefore.Add(amount), balanceAfter)
			supplyAfter := s.keeper.GetSupply(ctx, contractID, classID)
			s.Require().Equal(supplyBefore.Add(amount), supplyAfter)
			mintedAfter := s.keeper.GetMinted(ctx, contractID, classID)
			s.Require().Equal(mintedBefore.Add(amount), mintedAfter)
			burntAfter := s.keeper.GetBurnt(ctx, contractID, classID)
			s.Require().Equal(burntBefore, burntAfter)
		}
	})
}

func (s *KeeperTestSuite) TestMintNFT() {
	testCases := map[string]struct {
		contractID string
		param      collection.MintNFTParam
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			param:      collection.MintNFTParam{TokenType: s.nftClassID},
		},
		"class not found": {
			contractID: s.contractID,
			param:      collection.MintNFTParam{TokenType: "deadbeef"},
			err:        collection.ErrTokenTypeNotExist,
		},
		"not a class id of nft": {
			contractID: s.contractID,
			param:      collection.MintNFTParam{TokenType: s.ftClassID},
			err:        collection.ErrTokenTypeNotExist,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			// gather state
			classID := tc.param.TokenType
			supplyBefore := s.keeper.GetSupply(ctx, tc.contractID, classID)
			mintedBefore := s.keeper.GetMinted(ctx, tc.contractID, classID)
			burntBefore := s.keeper.GetBurnt(ctx, tc.contractID, classID)

			tokens, err := s.keeper.MintNFT(ctx, tc.contractID, s.stranger, []collection.MintNFTParam{tc.param})
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			amount := sdk.OneInt()
			s.Require().Len(tokens, 1)
			tokenID := tokens[0].TokenId
			balanceAfter := s.keeper.GetBalance(ctx, tc.contractID, s.stranger, tokenID)
			s.Require().Equal(amount, balanceAfter)
			supplyAfter := s.keeper.GetSupply(ctx, tc.contractID, classID)
			s.Require().Equal(supplyBefore.Add(amount), supplyAfter)
			mintedAfter := s.keeper.GetMinted(ctx, tc.contractID, classID)
			s.Require().Equal(mintedBefore.Add(amount), mintedAfter)
			burntAfter := s.keeper.GetBurnt(ctx, tc.contractID, classID)
			s.Require().Equal(burntBefore, burntAfter)
		})
	}
}

func (s *KeeperTestSuite) TestBurnCoins() {
	testCases := map[string]struct {
		contractID string
		amount     collection.Coin
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			amount:     collection.NewFTCoin(s.ftClassID, sdk.OneInt()),
		},
		"insufficient tokens": {
			contractID: s.contractID,
			amount:     collection.NewFTCoin("00bab10c", sdk.OneInt()),
			err:        collection.ErrInsufficientToken,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			// gather state
			classID := collection.SplitTokenID(tc.amount.TokenId)
			balanceBefore := s.keeper.GetBalance(ctx, tc.contractID, s.vendor, collection.NewFTID(classID))
			supplyBefore := s.keeper.GetSupply(ctx, tc.contractID, classID)
			mintedBefore := s.keeper.GetMinted(ctx, tc.contractID, classID)
			burntBefore := s.keeper.GetBurnt(ctx, tc.contractID, classID)

			_, err := s.keeper.BurnCoins(ctx, tc.contractID, s.vendor, collection.NewCoins(tc.amount))
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			amount := tc.amount.Amount
			balanceAfter := s.keeper.GetBalance(ctx, tc.contractID, s.vendor, collection.NewFTID(classID))
			s.Require().Equal(balanceBefore.Sub(amount), balanceAfter)
			supplyAfter := s.keeper.GetSupply(ctx, tc.contractID, classID)
			s.Require().Equal(supplyBefore.Sub(amount), supplyAfter)
			mintedAfter := s.keeper.GetMinted(ctx, tc.contractID, classID)
			s.Require().Equal(mintedBefore, mintedAfter)
			burntAfter := s.keeper.GetBurnt(ctx, tc.contractID, classID)
			s.Require().Equal(burntBefore.Add(amount), burntAfter)
		})
	}
}

func (s *KeeperTestSuite) TestModifyContract() {
	contractDescriptions := map[string]string{
		s.contractID: "valid",
		"deadbeef":   "not-exist",
	}
	changes := []collection.Attribute{
		{Key: collection.AttributeKeyName.String(), Value: "fox"},
		{Key: collection.AttributeKeyURI.String(), Value: "file:///fox.png"},
		{Key: collection.AttributeKeyMeta.String(), Value: "Fox"},
	}

	for contractID, contractDesc := range contractDescriptions {
		name := fmt.Sprintf("Contract: %s", contractDesc)
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			call := func() {
				s.keeper.ModifyContract(ctx, contractID, s.vendor, changes)
			}

			if contractID != s.contractID {
				s.Require().Panics(call)
				return
			}
			call()

			contract, err := s.keeper.GetContract(ctx, contractID)
			s.Require().NoError(err)

			s.Require().Equal(changes[0].Value, contract.Name)
			s.Require().Equal(changes[1].Value, contract.Uri)
			s.Require().Equal(changes[2].Value, contract.Meta)
		})
	}
}

func (s *KeeperTestSuite) TestModifyTokenClass() {
	classDescriptions := map[string]string{
		s.nftClassID: "valid",
		"deadbeef":   "not-exist",
	}
	changes := []collection.Attribute{
		{Key: collection.AttributeKeyName.String(), Value: "arctic fox"},
		{Key: collection.AttributeKeyMeta.String(), Value: "Arctic Fox"},
	}

	for classID, classDesc := range classDescriptions {
		name := fmt.Sprintf("Class: %s", classDesc)
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.ModifyTokenClass(ctx, s.contractID, classID, s.vendor, changes)
			if classID != s.nftClassID {
				s.Require().ErrorIs(err, collection.ErrTokenTypeNotExist)
				return
			}
			s.Require().NoError(err)

			class, err := s.keeper.GetTokenClass(ctx, s.contractID, classID)
			s.Require().NoError(err)

			nftClass, ok := class.(*collection.NFTClass)
			s.Require().True(ok)

			s.Require().Equal(changes[0].Value, nftClass.Name)
			s.Require().Equal(changes[1].Value, nftClass.Meta)
		})
	}
}

func (s *KeeperTestSuite) TestModifyNFT() {
	validTokenID := collection.NewNFTID(s.nftClassID, 1)
	tokenDescriptions := map[string]string{
		validTokenID:                       "valid",
		collection.NewNFTID("deadbeef", 1): "not-exist",
	}
	changes := []collection.Attribute{
		{Key: collection.AttributeKeyName.String(), Value: "fennec fox 1"},
		{Key: collection.AttributeKeyMeta.String(), Value: "Fennec Fox 1"},
	}

	for tokenID, tokenDesc := range tokenDescriptions {
		name := fmt.Sprintf("Token: %s", tokenDesc)
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.ModifyNFT(ctx, s.contractID, tokenID, s.vendor, changes)
			if tokenID != validTokenID {
				s.Require().ErrorIs(err, collection.ErrTokenNotExist)
				return
			}
			s.Require().NoError(err)

			nft, err := s.keeper.GetNFT(ctx, s.contractID, tokenID)
			s.Require().NoError(err)

			s.Require().Equal(changes[0].Value, nft.Name)
			s.Require().Equal(changes[1].Value, nft.Meta)
		})
	}
}
