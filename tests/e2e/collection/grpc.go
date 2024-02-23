package collection

import (
	"fmt"

	"github.com/cosmos/gogoproto/proto"

	cmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Finschia/finschia-sdk/x/collection"
)

const NotExistContractID = "aaaabbbb"

func (s *E2ETestSuite) TestBalanceGRPC() {
	val := s.network.Validators[0]
	tokenID := s.mintNFT(s.contractID, s.vendor, s.customer, s.nftClassID)

	testCases := []struct {
		name       string
		url        string
		expErr     bool
		respType   proto.Message
		expectResp proto.Message
	}{
		{
			"valid request",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/balances/%s/%s", val.APIAddress, s.contractID, s.customer.String(), tokenID),
			false,
			&collection.QueryBalanceResponse{},
			&collection.QueryBalanceResponse{
				Balance: collection.NewCoin(tokenID, cmath.OneInt()),
			},
		},
		{
			"not own NFT",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/balances/%s/%s", val.APIAddress, s.contractID, s.stranger.String(), tokenID),
			false,
			&collection.QueryBalanceResponse{},
			&collection.QueryBalanceResponse{
				Balance: collection.Coin{TokenId: tokenID, Amount: cmath.ZeroInt()},
			},
		},
		{
			"invalid contract ID",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/balances/%s/%s", val.APIAddress, "wrong id", s.customer.String(), tokenID),
			true,
			&collection.QueryBalanceResponse{},
			&collection.QueryBalanceResponse{},
		},
		{
			"invalid token ID",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/balances/%s/%s", val.APIAddress, s.contractID, s.customer.String(), "wrong id"),
			true,
			&collection.QueryBalanceResponse{},
			&collection.QueryBalanceResponse{},
		},
		{
			"invalid address",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/balances/%s/%s", val.APIAddress, s.contractID, "wrong address", tokenID),
			true,
			&collection.QueryBalanceResponse{},
			&collection.QueryBalanceResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequest(tc.url)
			s.Require().NoError(err)

			err = s.cfg.Codec.UnmarshalJSON(resp, tc.respType)

			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectResp.String(), tc.respType.String())
			}
		})
	}
}

func (s *E2ETestSuite) TestBalancesGRPC() {
	val := s.network.Validators[0]
	tokenID := s.mintNFT(s.contractID, s.vendor, s.customer, s.nftClassID)

	testCases := []struct {
		name     string
		url      string
		expErr   bool
		respType proto.Message
	}{
		{
			"valid request",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/balances/%s", val.APIAddress, s.contractID, s.vendor.String()),
			false,
			&collection.QueryAllBalancesResponse{},
		},
		{
			"invalid contract ID",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/balances/%s/%s", val.APIAddress, "wrong id", s.vendor.String(), tokenID),
			true,
			&collection.QueryAllBalancesResponse{},
		},
		{
			"invalid address",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/balances/%s/%s", val.APIAddress, s.contractID, "wrong address", tokenID),
			true,
			&collection.QueryAllBalancesResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequest(tc.url)
			s.Require().NoError(err)

			err = s.cfg.Codec.UnmarshalJSON(resp, tc.respType)

			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().GreaterOrEqual(len(tc.respType.(*collection.QueryAllBalancesResponse).Balances), 1)
			}
		})
	}
}

func (s *E2ETestSuite) TestStatisticsGRPC() {
	val := s.network.Validators[0]
	tokenID := s.mintNFT(s.contractID, s.vendor, s.vendor, s.nftClassID)
	_ = s.burnNFT(s.contractID, s.vendor, tokenID)

	testCases := []struct {
		name     string
		url      string
		expErr   bool
		respType proto.Message
	}{
		{
			"valid request - Supply",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s/supply", val.APIAddress, s.contractID, s.nftClassID),
			false,
			&collection.QueryNFTSupplyResponse{},
		},
		{
			"valid request - Minted",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s/minted", val.APIAddress, s.contractID, s.nftClassID),
			false,
			&collection.QueryNFTMintedResponse{},
		},
		{
			"valid request - Burnt",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s/burnt", val.APIAddress, s.contractID, s.nftClassID),
			false,
			&collection.QueryNFTBurntResponse{},
		},
		{
			"invalid request (wrong token type) - Supply",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s/supply", val.APIAddress, s.contractID, "wrong ID"),
			true,
			&collection.QueryNFTSupplyResponse{},
		},
		{
			"invalid request (wrong token type) - Minted",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s/minted", val.APIAddress, s.contractID, "wrong ID"),
			true,
			&collection.QueryNFTMintedResponse{},
		},
		{
			"invalid request (wrong token type) - Burnt",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s/burnt", val.APIAddress, s.contractID, "wrong ID"),
			true,
			&collection.QueryNFTBurntResponse{},
		},
		{
			"invalid request (wrong contract ID) - Supply",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s/supply", val.APIAddress, "wrong ID", s.nftClassID),
			true,
			&collection.QueryNFTSupplyResponse{},
		},
		{
			"invalid request (wrong contract ID) - Minted",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s/supply", val.APIAddress, "wrong ID", s.nftClassID),
			true,
			&collection.QueryNFTMintedResponse{},
		},
		{
			"invalid request (wrong contract ID) - Burnt",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s/supply", val.APIAddress, "wrong ID", s.nftClassID),
			true,
			&collection.QueryNFTBurntResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequest(tc.url)
			s.Require().NoError(err)

			err = s.cfg.Codec.UnmarshalJSON(resp, tc.respType)

			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				switch r := tc.respType.(type) {
				case *collection.QueryNFTSupplyResponse:
					s.Require().True(r.Supply.GTE(cmath.OneInt()))
				case *collection.QueryNFTMintedResponse:
					s.Require().True(r.Minted.GTE(cmath.OneInt()))
				case *collection.QueryNFTBurntResponse:
					s.Require().True(r.Burnt.GTE(cmath.OneInt()))
				default:
					s.Require().Fail("unexpected response type")
				}
			}
		})
	}
}

func (s *E2ETestSuite) TestContractGRPC() {
	val := s.network.Validators[0]

	testCases := []struct {
		name       string
		url        string
		expErr     bool
		respType   proto.Message
		expectResp proto.Message
	}{
		{
			"valid request",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s", val.APIAddress, s.contractID),
			false,
			&collection.QueryContractResponse{},
			&collection.QueryContractResponse{
				Contract: collection.Contract{
					Name: "",
					Id:   s.contractID,
					Meta: "",
					Uri:  "",
				},
			},
		},
		{
			"invalid request - wrong contract ID",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s", val.APIAddress, "wrong ID"),
			true,
			&collection.QueryContractResponse{},
			&collection.QueryContractResponse{},
		},
		{
			"invalid request - not found",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s", val.APIAddress, NotExistContractID),
			true,
			&collection.QueryContractResponse{},
			&collection.QueryContractResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequest(tc.url)
			s.Require().NoError(err)

			err = s.cfg.Codec.UnmarshalJSON(resp, tc.respType)

			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectResp.String(), tc.respType.String())
			}
		})
	}
}

func (s *E2ETestSuite) TestTokenClassTypeNameGRPC() {
	val := s.network.Validators[0]

	testCases := []struct {
		name       string
		url        string
		expErr     bool
		respType   proto.Message
		expectResp proto.Message
	}{
		{
			"valid request",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_classes/%s/type_name", val.APIAddress, s.contractID, s.nftClassID),
			false,
			&collection.QueryTokenClassTypeNameResponse{},
			&collection.QueryTokenClassTypeNameResponse{
				Name: proto.MessageName(&collection.NFTClass{}),
			},
		},
		{
			"invalid request - wrong contract ID",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_classes/%s/type_name", val.APIAddress, s.contractID, "wrong ID"),
			true,
			&collection.QueryTokenClassTypeNameResponse{},
			&collection.QueryTokenClassTypeNameResponse{},
		},
		{
			"invalid request - wrong class ID",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_classes/%s/type_name", val.APIAddress, "wrong ID", s.nftClassID),
			true,
			&collection.QueryTokenClassTypeNameResponse{},
			&collection.QueryTokenClassTypeNameResponse{},
		},
		{
			"invalid request - not found",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_classes/%s/type_name", val.APIAddress, NotExistContractID, s.nftClassID),
			true,
			&collection.QueryTokenClassTypeNameResponse{},
			&collection.QueryTokenClassTypeNameResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequest(tc.url)
			s.Require().NoError(err)

			err = s.cfg.Codec.UnmarshalJSON(resp, tc.respType)

			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectResp.String(), tc.respType.String())
			}
		})
	}
}

func (s *E2ETestSuite) TestTokenTypeNameGRPC() {
	val := s.network.Validators[0]
	classID := s.createNFTClass(s.contractID, s.vendor)

	testCases := []struct {
		name       string
		url        string
		expErr     bool
		respType   proto.Message
		expectResp proto.Message
	}{
		{
			"valid request",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s", val.APIAddress, s.contractID, classID),
			false,
			&collection.QueryTokenTypeResponse{},
			&collection.QueryTokenTypeResponse{
				TokenType: collection.TokenType{
					ContractId: s.contractID,
					TokenType:  classID,
					Name:       "",
					Meta:       "",
				},
			},
		},
		{
			"invalid request - wrong contract ID",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s", val.APIAddress, "wrong ID", classID),
			true,
			&collection.QueryTokenTypeResponse{},
			&collection.QueryTokenTypeResponse{},
		},
		{
			"invalid request - wrong class ID",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s", val.APIAddress, s.contractID, "wrong ID"),
			true,
			&collection.QueryTokenTypeResponse{},
			&collection.QueryTokenTypeResponse{},
		},
		{
			"invalid request - not found",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/token_types/%s", val.APIAddress, NotExistContractID, classID),
			true,
			&collection.QueryTokenTypeResponse{},
			&collection.QueryTokenTypeResponse{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequest(tc.url)
			s.Require().NoError(err)

			err = s.cfg.Codec.UnmarshalJSON(resp, tc.respType)

			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectResp.String(), tc.respType.String())
			}
		})
	}
}

func (s *E2ETestSuite) TestTokenGRPC() {
	val := s.network.Validators[0]
	tokenID := s.mintNFT(s.contractID, s.vendor, s.vendor, s.nftClassID)

	testCases := []struct {
		name       string
		url        string
		expErr     bool
		respType   proto.Message
		expectResp proto.Message
	}{
		{
			"valid request",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/tokens/%s", val.APIAddress, s.contractID, tokenID),
			false,
			&collection.QueryTokenResponse{},
			&collection.OwnerNFT{
				ContractId: s.contractID,
				TokenId:    tokenID,
				Owner:      s.vendor.String(),
				Name:       "arctic fox",
				Meta:       "",
			},
		},
		{
			"invalid request - wrong contract ID",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/tokens/%s", val.APIAddress, "wrong ID", tokenID),
			true,
			&collection.QueryTokenResponse{},
			&collection.OwnerNFT{},
		},
		{
			"invalid request - wrong class ID",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/tokens/%s", val.APIAddress, s.contractID, "wrong ID"),
			true,
			&collection.QueryTokenResponse{},
			&collection.OwnerNFT{},
		},
		{
			"invalid request - not found",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/tokens/%s", val.APIAddress, NotExistContractID, tokenID),
			true,
			&collection.QueryTokenResponse{},
			&collection.OwnerNFT{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequest(tc.url)
			s.Require().NoError(err)

			err = s.cfg.Codec.UnmarshalJSON(resp, tc.respType)

			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				res := tc.respType.(*collection.QueryTokenResponse)
				s.Require().Equal("/lbm.collection.v1.OwnerNFT", res.Token.TypeUrl)
				var token collection.Token
				err = s.cfg.InterfaceRegistry.UnpackAny(&res.Token, &token)
				s.Require().NoError(err)
				nft, ok := token.(*collection.OwnerNFT)
				s.Require().True(ok)
				s.Require().Equal(tc.expectResp.String(), nft.String())
			}
		})
	}
}

func (s *E2ETestSuite) TestGranteeGrantsGRPC() {
	val := s.network.Validators[0]
	s.grant(s.contractID, s.vendor, s.stranger, collection.PermissionIssue)
	dummyAddr := "link1hcpqj6w2eq30jcdggs7892lmask0cacvynqg7d"

	testCases := []struct {
		name       string
		url        string
		expErr     bool
		respType   proto.Message
		expectResp proto.Message
	}{
		{
			"valid request",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/grants/%s", val.APIAddress, s.contractID, s.stranger),
			false,
			&collection.QueryGranteeGrantsResponse{},
			&collection.QueryGranteeGrantsResponse{
				Grants: []collection.Grant{
					{
						Grantee:    s.stranger.String(),
						Permission: collection.PermissionIssue,
					},
				},
				Pagination: &query.PageResponse{
					Total: 1,
				},
			},
		},
		{
			"invalid request - wrong contract ID",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/grants/%s", val.APIAddress, "wrong ID", s.stranger),
			true,
			&collection.QueryGranteeGrantsResponse{},
			&collection.QueryGranteeGrantsResponse{},
		},
		{
			"invalid request - not found",
			fmt.Sprintf("%s/lbm/collection/v1/contracts/%s/grants/%s", val.APIAddress, s.contractID, dummyAddr),
			false,
			&collection.QueryGranteeGrantsResponse{},
			&collection.QueryGranteeGrantsResponse{
				Pagination: &query.PageResponse{
					Total: 0,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			resp, err := testutil.GetRequest(tc.url)
			s.Require().NoError(err)

			err = s.cfg.Codec.UnmarshalJSON(resp, tc.respType)

			if tc.expErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectResp.String(), tc.respType.String())
			}
		})
	}
}
