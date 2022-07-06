package keeper_test

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

func (s *KeeperTestSuite) TestCreateContract() {
	ctx, _ := s.ctx.CacheContext()

	input := collection.Contract{
		Name:       "tibetian fox",
		Meta:       "Tibetian Fox",
		BaseImgUri: "file:///tibetian_fox.png",
	}
	id := s.keeper.CreateContract(ctx, s.vendor, input)
	s.Require().NotEmpty(id)

	output, err := s.keeper.GetContract(ctx, id)
	s.Require().NoError(err)
	s.Require().NotNil(output)

	s.Require().Equal(id, output.ContractId)
	s.Require().Equal(input.Name, output.Name)
	s.Require().Equal(input.Meta, output.Meta)
	s.Require().Equal(input.BaseImgUri, output.BaseImgUri)
}

func (s *KeeperTestSuite) TestCreateTokenClass() {
	testCases := map[string]struct {
		contractID string
		class      collection.TokenClass
		valid      bool
	}{
		"valid fungible token class": {
			contractID: s.contractID,
			class:      &collection.FTClass{},
			valid:      true,
		},
		"valid non-fungible token class": {
			contractID: s.contractID,
			class:      &collection.NFTClass{},
			valid:      true,
		},
		"invalid contract id": {
			class: &collection.FTClass{},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			id, err := s.keeper.CreateTokenClass(s.ctx, tc.contractID, tc.class)
			if !tc.valid {
				s.Require().Error(err)
				s.Require().Nil(id)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(id)

			class, err := s.keeper.GetTokenClass(s.ctx, tc.contractID, *id)
			s.Require().NoError(err)
			s.Require().NoError(class.ValidateBasic())
		})
	}
}

func (s *KeeperTestSuite) TestMintFT() {
	testCases := map[string]struct {
		contractID string
		amount     collection.Coin
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			amount:     collection.NewFTCoin(s.ftClassID, sdk.OneInt()),
			valid:      true,
		},
		"invalid token id": {
			contractID: s.contractID,
			amount:     collection.NewNFTCoin(s.ftClassID, 1),
		},
		"class not found": {
			contractID: s.contractID,
			amount:     collection.NewFTCoin("00bab10c", sdk.OneInt()),
		},
		"not a class id of ft": {
			contractID: s.contractID,
			amount:     collection.NewFTCoin(s.nftClassID, sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			err := s.keeper.MintFT(s.ctx, tc.contractID, s.stranger, collection.NewCoins(tc.amount))
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
		})
	}
}

func (s *KeeperTestSuite) TestMintNFT() {
	testCases := map[string]struct {
		contractID string
		params     []collection.MintNFTParam
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			params:     []collection.MintNFTParam{{TokenType: s.nftClassID}},
			valid:      true,
		},
		"class not found": {
			contractID: s.contractID,
			params:     []collection.MintNFTParam{{TokenType: "deadbeef"}},
		},
		"not a class id of nft": {
			contractID: s.contractID,
			params:     []collection.MintNFTParam{{TokenType: s.ftClassID}},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			_, err := s.keeper.MintNFT(s.ctx, tc.contractID, s.stranger, tc.params)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
		})
	}
}

func (s *KeeperTestSuite) TestModifyContract() {
	contractDescriptions := map[string]string{
		s.contractID: "valid",
		"deadbeef":   "not-exist",
	}
	changes := []collection.Attribute{
		{Key: collection.AttributeKey_Name.String(), Value: "fox"},
		{Key: collection.AttributeKey_BaseImgURI.String(), Value: "file:///fox.png"},
		{Key: collection.AttributeKey_Meta.String(), Value: "Fox"},
	}

	for contractID, contractDesc := range contractDescriptions {
		name := fmt.Sprintf("Contract: %s", contractDesc)
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.ModifyContract(ctx, contractID, s.vendor, changes)
			if contractID == s.contractID {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestModifyTokenClass() {
	contractDescriptions := map[string]string{
		s.contractID: "valid",
		"deadbeef":   "not-exist",
	}
	classDescriptions := map[string]string{
		s.nftClassID: "valid",
		"deadbeef":   "not-exist",
	}
	changes := []collection.Attribute{
		{Key: collection.AttributeKey_Name.String(), Value: "arctic fox"},
		{Key: collection.AttributeKey_Meta.String(), Value: "Arctic Fox"},
	}

	for contractID, contractDesc := range contractDescriptions {
		for classID, classDesc := range classDescriptions {
			name := fmt.Sprintf("Contract: %s, Class: %s", contractDesc, classDesc)
			s.Run(name, func() {
				ctx, _ := s.ctx.CacheContext()

				err := s.keeper.ModifyTokenClass(ctx, contractID, classID, s.vendor, changes)
				if contractID == s.contractID && classID == s.nftClassID {
					s.Require().NoError(err)
				} else {
					s.Require().Error(err)
				}
			})
		}
	}
}

func (s *KeeperTestSuite) TestModifyNFT() {
	contractDescriptions := map[string]string{
		s.contractID: "valid",
		"deadbeef":   "not-exist",
	}
	validTokenID := collection.NewNFTID(s.nftClassID, 1)
	tokenDescriptions := map[string]string{
		validTokenID:                       "valid",
		collection.NewNFTID("deadbeef", 1): "not-exist",
	}
	changes := []collection.Attribute{
		{Key: collection.AttributeKey_Name.String(), Value: "fennec fox 1"},
		{Key: collection.AttributeKey_Meta.String(), Value: "Fennec Fox 1"},
	}

	for contractID, contractDesc := range contractDescriptions {
		for tokenID, tokenDesc := range tokenDescriptions {
			name := fmt.Sprintf("Contract: %s, Token: %s", contractDesc, tokenDesc)
			s.Run(name, func() {
				ctx, _ := s.ctx.CacheContext()

				err := s.keeper.ModifyNFT(ctx, contractID, tokenID, s.vendor, changes)
				if contractID == s.contractID && tokenID == validTokenID {
					s.Require().NoError(err)
				} else {
					s.Require().Error(err)
				}
			})
		}
	}
}
