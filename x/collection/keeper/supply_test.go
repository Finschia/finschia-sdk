package keeper_test

import (
	"github.com/line/lbm-sdk/x/collection"
)

func (s *KeeperTestSuite) TestCreateContract() {
	ctx, _ := s.ctx.CacheContext()

	input := collection.Contract{
		Name: "tibetian fox",
		Meta: "Tibetian Fox",
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
		class collection.TokenClass
		valid  bool
	}{
		"valid fungible token class": {
			contractID: s.contractID,
			class: &collection.FTClass{},
			valid: true,
		},
		"valid non-fungible token class": {
			contractID: s.contractID,
			class: &collection.NFTClass{},
			valid: true,
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
