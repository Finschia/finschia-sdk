package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/collection"
)

func (s *KeeperTestSuite) TestQueryBalance() {
	// empty request
	_, err := s.queryServer.Balance(s.goCtx, nil)
	s.Require().Error(err)

	tokenID := collection.NewFTID(s.ftClassID)
	testCases := map[string]struct {
		contractID string
		address    sdk.AccAddress
		tokenID    string
		valid      bool
		postTest   func(res *collection.QueryBalanceResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			address:    s.vendor,
			tokenID:    tokenID,
			valid:      true,
			postTest: func(res *collection.QueryBalanceResponse) {
				expected := collection.NewCoin(tokenID, s.balance)
				s.Require().Equal(expected, res.Balance)
			},
		},
		"invalid contract id": {
			address: s.vendor,
			tokenID: tokenID,
		},
		"invalid address": {
			contractID: s.contractID,
			tokenID:    tokenID,
		},
		"valid token id": {
			contractID: s.contractID,
			address:    s.vendor,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryBalanceRequest{
				ContractId: tc.contractID,
				Address:    tc.address.String(),
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.Balance(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryAllBalances() {
	// empty request
	_, err := s.queryServer.AllBalances(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		address    sdk.AccAddress
		valid      bool
		count      uint64
		postTest   func(res *collection.QueryAllBalancesResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			address:    s.customer,
			valid:      true,
			postTest: func(res *collection.QueryAllBalancesResponse) {
				s.Require().Equal(s.numNFTs+1, len(res.Balances))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			address:    s.customer,
			valid:      true,
			count:      1,
			postTest: func(res *collection.QueryAllBalancesResponse) {
				s.Require().Equal(1, len(res.Balances))
			},
		},
		"invalid contract id": {
			address: s.customer,
		},
		"invalid address": {
			contractID: s.contractID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			pageReq := &query.PageRequest{}
			if tc.count != 0 {
				pageReq.Limit = tc.count
			}
			req := &collection.QueryAllBalancesRequest{
				ContractId: tc.contractID,
				Address:    tc.address.String(),
				Pagination: pageReq,
			}
			res, err := s.queryServer.AllBalances(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryFTSupply() {
	// empty request
	_, err := s.queryServer.FTSupply(s.goCtx, nil)
	s.Require().Error(err)

	tokenID := collection.NewFTID(s.ftClassID)
	testCases := map[string]struct {
		contractID string
		tokenID    string
		valid      bool
		postTest   func(res *collection.QueryFTSupplyResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			valid:      true,
			postTest: func(res *collection.QueryFTSupplyResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(3)), res.Supply)
			},
		},
		"invalid contract id": {
			tokenID: tokenID,
		},
		"invalid token id": {
			contractID: s.contractID,
		},
		"no such a token": {
			contractID: s.contractID,
			tokenID:    collection.NewFTID("00bab10c"),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryFTSupplyRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.FTSupply(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryFTMinted() {
	// empty request
	_, err := s.queryServer.FTMinted(s.goCtx, nil)
	s.Require().Error(err)

	tokenID := collection.NewFTID(s.ftClassID)
	testCases := map[string]struct {
		contractID string
		tokenID    string
		valid      bool
		postTest   func(res *collection.QueryFTMintedResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			valid:      true,
			postTest: func(res *collection.QueryFTMintedResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(6)), res.Minted)
			},
		},
		"invalid contract id": {
			tokenID: tokenID,
		},
		"invalid token id": {
			contractID: s.contractID,
		},
		"no such a token": {
			contractID: s.contractID,
			tokenID:    collection.NewFTID("00bab10c"),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryFTMintedRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.FTMinted(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryFTBurnt() {
	// empty request
	_, err := s.queryServer.FTBurnt(s.goCtx, nil)
	s.Require().Error(err)

	tokenID := collection.NewFTID(s.ftClassID)
	testCases := map[string]struct {
		contractID string
		tokenID    string
		valid      bool
		postTest   func(res *collection.QueryFTBurntResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			valid:      true,
			postTest: func(res *collection.QueryFTBurntResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(3)), res.Burnt)
			},
		},
		"invalid contract id": {
			tokenID: tokenID,
		},
		"invalid token id": {
			contractID: s.contractID,
		},
		"no such a token": {
			contractID: s.contractID,
			tokenID:    collection.NewFTID("00bab10c"),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryFTBurntRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.FTBurnt(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryNFTSupply() {
	// empty request
	_, err := s.queryServer.NFTSupply(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		tokenType  string
		valid      bool
		postTest   func(res *collection.QueryNFTSupplyResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenType:  s.nftClassID,
			valid:      true,
			postTest: func(res *collection.QueryNFTSupplyResponse) {
				s.Require().EqualValues(s.numNFTs*3, res.Supply.Int64())
			},
		},
		"invalid contract id": {
			tokenType: s.nftClassID,
		},
		"invalid token type": {
			contractID: s.contractID,
		},
		"no such a token type": {
			contractID: s.contractID,
			tokenType:  "deadbeef",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryNFTSupplyRequest{
				ContractId: tc.contractID,
				TokenType:  tc.tokenType,
			}
			res, err := s.queryServer.NFTSupply(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryNFTMinted() {
	// empty request
	_, err := s.queryServer.NFTMinted(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		tokenType  string
		valid      bool
		postTest   func(res *collection.QueryNFTMintedResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenType:  s.nftClassID,
			valid:      true,
			postTest: func(res *collection.QueryNFTMintedResponse) {
				s.Require().EqualValues(s.numNFTs*3, res.Minted.Int64())
			},
		},
		"invalid contract id": {
			tokenType: s.nftClassID,
		},
		"invalid token type": {
			contractID: s.contractID,
		},
		"no such a token type": {
			contractID: s.contractID,
			tokenType:  "deadbeef",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryNFTMintedRequest{
				ContractId: tc.contractID,
				TokenType:  tc.tokenType,
			}
			res, err := s.queryServer.NFTMinted(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryNFTBurnt() {
	// empty request
	_, err := s.queryServer.NFTBurnt(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		tokenType  string
		valid      bool
		postTest   func(res *collection.QueryNFTBurntResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenType:  s.nftClassID,
			valid:      true,
			postTest: func(res *collection.QueryNFTBurntResponse) {
				s.Require().Equal(sdk.ZeroInt(), res.Burnt)
			},
		},
		"invalid contract id": {
			tokenType: s.nftClassID,
		},
		"invalid token type": {
			contractID: s.contractID,
		},
		"no such a token type": {
			contractID: s.contractID,
			tokenType:  "deadbeef",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryNFTBurntRequest{
				ContractId: tc.contractID,
				TokenType:  tc.tokenType,
			}
			res, err := s.queryServer.NFTBurnt(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryContract() {
	// empty request
	_, err := s.queryServer.Contract(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		valid      bool
		postTest   func(res *collection.QueryContractResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			valid:      true,
			postTest: func(res *collection.QueryContractResponse) {
				s.Require().Equal(s.contractID, res.Contract.ContractId)
			},
		},
		"invalid contract id": {},
		"no such an id": {
			contractID: "deadbeef",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryContractRequest{
				ContractId: tc.contractID,
			}
			res, err := s.queryServer.Contract(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryTokenType() {
	// empty request
	_, err := s.queryServer.TokenType(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		tokenType  string
		valid      bool
		postTest   func(res *collection.QueryTokenTypeResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenType:  s.nftClassID,
			valid:      true,
			postTest: func(res *collection.QueryTokenTypeResponse) {
				s.Require().Equal(s.contractID, res.TokenType.ContractId)
				s.Require().Equal(s.nftClassID, res.TokenType.TokenType)
			},
		},
		"invalid contract id": {
			tokenType: s.nftClassID,
		},
		"invalid token type": {
			contractID: s.contractID,
		},
		"no such a token type": {
			contractID: s.contractID,
			tokenType:  "deadbeef",
		},
		"not a class of nft": {
			contractID: s.contractID,
			tokenType:  s.ftClassID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryTokenTypeRequest{
				ContractId: tc.contractID,
				TokenType:  tc.tokenType,
			}
			res, err := s.queryServer.TokenType(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryTokenTypes() {
	// empty request
	_, err := s.queryServer.TokenTypes(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		valid      bool
		count      uint64
		postTest   func(res *collection.QueryTokenTypesResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			valid:      true,
			postTest: func(res *collection.QueryTokenTypesResponse) {
				s.Require().Equal(1, len(res.TokenTypes))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			valid:      true,
			count:      1,
			postTest: func(res *collection.QueryTokenTypesResponse) {
				s.Require().Equal(1, len(res.TokenTypes))
			},
		},
		"invalid contract id": {},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			pageReq := &query.PageRequest{}
			if tc.count != 0 {
				pageReq.Limit = tc.count
			}
			req := &collection.QueryTokenTypesRequest{
				ContractId: tc.contractID,
				Pagination: pageReq,
			}
			res, err := s.queryServer.TokenTypes(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryToken() {
	// empty request
	_, err := s.queryServer.Token(s.goCtx, nil)
	s.Require().Error(err)

	ftTokenID := collection.NewFTID(s.ftClassID)
	nftTokenID := collection.NewNFTID(s.nftClassID, 1)
	testCases := map[string]struct {
		contractID string
		tokenID    string
		valid      bool
		postTest   func(res *collection.QueryTokenResponse)
	}{
		"valid ft request": {
			contractID: s.contractID,
			tokenID:    ftTokenID,
			valid:      true,
			postTest: func(res *collection.QueryTokenResponse) {
				s.Require().Equal("/lbm.collection.v1.FT", res.Token.TypeUrl)
				token := collection.TokenFromAny(&res.Token)
				ft, ok := token.(*collection.FT)
				s.Require().True(ok)
				s.Require().Equal(s.contractID, ft.ContractId)
				s.Require().Equal(ftTokenID, ft.TokenId)
			},
		},
		"valid nft request": {
			contractID: s.contractID,
			tokenID:    nftTokenID,
			valid:      true,
			postTest: func(res *collection.QueryTokenResponse) {
				s.Require().Equal("/lbm.collection.v1.OwnerNFT", res.Token.TypeUrl)
				token := collection.TokenFromAny(&res.Token)
				nft, ok := token.(*collection.OwnerNFT)
				s.Require().True(ok)
				s.Require().Equal(s.contractID, nft.ContractId)
				s.Require().Equal(nftTokenID, nft.TokenId)
			},
		},
		"invalid contract id": {
			tokenID: ftTokenID,
		},
		"invalid token id": {
			contractID: s.contractID,
		},
		"no such a fungible token": {
			contractID: s.contractID,
			tokenID:    collection.NewFTID("00bab10c"),
		},
		"no such a non-fungible token": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID("deadbeef", 1),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryTokenRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.Token(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryTokensWithTokenType() {
	// empty request
	_, err := s.queryServer.TokensWithTokenType(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		tokenType  string
		valid      bool
		count      uint64
		postTest   func(res *collection.QueryTokensWithTokenTypeResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenType:  s.nftClassID,
			valid:      true,
			count:      1000000,
			postTest: func(res *collection.QueryTokensWithTokenTypeResponse) {
				s.Require().Equal(s.numNFTs*3, len(res.Tokens))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			tokenType:  s.nftClassID,
			valid:      true,
			count:      1,
			postTest: func(res *collection.QueryTokensWithTokenTypeResponse) {
				s.Require().Equal(1, len(res.Tokens))
			},
		},
		"invalid contract id": {
			tokenType: s.nftClassID,
		},
		"invalid token type": {
			contractID: s.contractID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			pageReq := &query.PageRequest{}
			if tc.count != 0 {
				pageReq.Limit = tc.count
			}
			req := &collection.QueryTokensWithTokenTypeRequest{
				ContractId: tc.contractID,
				TokenType:  tc.tokenType,
				Pagination: pageReq,
			}
			res, err := s.queryServer.TokensWithTokenType(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryTokens() {
	// empty request
	_, err := s.queryServer.Tokens(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		valid      bool
		count      uint64
		postTest   func(res *collection.QueryTokensResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			valid:      true,
			count:      1000000,
			postTest: func(res *collection.QueryTokensResponse) {
				s.Require().Equal(s.numNFTs*3+1, len(res.Tokens))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			valid:      true,
			count:      1,
			postTest: func(res *collection.QueryTokensResponse) {
				s.Require().Equal(1, len(res.Tokens))
			},
		},
		"invalid contract id": {},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			pageReq := &query.PageRequest{}
			if tc.count != 0 {
				pageReq.Limit = tc.count
			}
			req := &collection.QueryTokensRequest{
				ContractId: tc.contractID,
				Pagination: pageReq,
			}
			res, err := s.queryServer.Tokens(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryRoot() {
	// empty request
	_, err := s.queryServer.Root(s.goCtx, nil)
	s.Require().Error(err)

	tokenID := collection.NewNFTID(s.nftClassID, 2)
	testCases := map[string]struct {
		contractID string
		tokenID    string
		valid      bool
		postTest   func(res *collection.QueryRootResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			valid:      true,
			postTest: func(res *collection.QueryRootResponse) {
				s.Require().Equal(collection.NewNFTID(s.nftClassID, 1), res.Root.Id)
			},
		},
		"invalid contract id": {
			tokenID: tokenID,
		},
		"invalid token id": {
			contractID: s.contractID,
		},
		"no such a token": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID("deadbeef", 1),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryRootRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.Root(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryParent() {
	// empty request
	_, err := s.queryServer.Parent(s.goCtx, nil)
	s.Require().Error(err)

	tokenID := collection.NewNFTID(s.nftClassID, 2)
	testCases := map[string]struct {
		contractID string
		tokenID    string
		valid      bool
		postTest   func(res *collection.QueryParentResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			valid:      true,
			postTest: func(res *collection.QueryParentResponse) {
				s.Require().Equal(collection.NewNFTID(s.nftClassID, 1), res.Parent.Id)
			},
		},
		"invalid contract id": {
			tokenID: tokenID,
		},
		"invalid token id": {
			contractID: s.contractID,
		},
		"no such a token": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID("deadbeef", 1),
		},
		"no parent": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID(s.nftClassID, 1),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryParentRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.Parent(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryChildren() {
	// empty request
	_, err := s.queryServer.Children(s.goCtx, nil)
	s.Require().Error(err)

	tokenID := collection.NewNFTID(s.nftClassID, 1)
	testCases := map[string]struct {
		contractID string
		tokenID    string
		valid      bool
		count      uint64
		postTest   func(res *collection.QueryChildrenResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			valid:      true,
			postTest: func(res *collection.QueryChildrenResponse) {
				s.Require().Equal(1, len(res.Children))
				s.Require().Equal(collection.NewNFTID(s.nftClassID, 2), res.Children[0].Id)
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			tokenID:    tokenID,
			valid:      true,
			count:      1,
			postTest: func(res *collection.QueryChildrenResponse) {
				s.Require().Equal(1, len(res.Children))
				s.Require().Equal(collection.NewNFTID(s.nftClassID, 2), res.Children[0].Id)
			},
		},
		"invalid contract id": {
			tokenID: tokenID,
		},
		"invalid token id": {
			contractID: s.contractID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			pageReq := &query.PageRequest{}
			if tc.count != 0 {
				pageReq.Limit = tc.count
			}
			req := &collection.QueryChildrenRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
				Pagination: pageReq,
			}
			res, err := s.queryServer.Children(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryGranteeGrants() {
	// empty request
	_, err := s.queryServer.GranteeGrants(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		valid      bool
		postTest   func(res *collection.QueryGranteeGrantsResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			grantee:    s.vendor,
			valid:      true,
			postTest: func(res *collection.QueryGranteeGrantsResponse) {
				s.Require().Equal(4, len(res.Grants))
			},
		},
		"invalid contract id": {
			grantee: s.vendor,
		},
		"invalid grantee": {
			contractID: s.contractID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryGranteeGrantsRequest{
				ContractId: tc.contractID,
				Grantee:    tc.grantee.String(),
			}
			res, err := s.queryServer.GranteeGrants(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryApproved() {
	// empty request
	_, err := s.queryServer.Approved(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		address    sdk.AccAddress
		approver   sdk.AccAddress
		valid      bool
		postTest   func(res *collection.QueryApprovedResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			address:    s.operator,
			approver:   s.customer,
			valid:      true,
			postTest: func(res *collection.QueryApprovedResponse) {
				s.Require().True(res.Approved)
			},
		},
		"invalid contract id": {
			address:  s.operator,
			approver: s.customer,
		},
		"invalid address": {
			contractID: s.contractID,
			approver:   s.customer,
		},
		"invalid approver": {
			contractID: s.contractID,
			address:    s.operator,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryApprovedRequest{
				ContractId: tc.contractID,
				Address:    tc.address.String(),
				Approver:   tc.approver.String(),
			}
			res, err := s.queryServer.Approved(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryApprovers() {
	// empty request
	_, err := s.queryServer.Approvers(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		address    sdk.AccAddress
		valid      bool
		count      uint64
		postTest   func(res *collection.QueryApproversResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			address:    s.operator,
			valid:      true,
			postTest: func(res *collection.QueryApproversResponse) {
				s.Require().Equal(1, len(res.Approvers))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			address:    s.operator,
			valid:      true,
			count:      1,
			postTest: func(res *collection.QueryApproversResponse) {
				s.Require().Equal(1, len(res.Approvers))
			},
		},
		"invalid contract id": {
			address: s.operator,
		},
		"invalid address": {
			contractID: s.contractID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			pageReq := &query.PageRequest{}
			if tc.count != 0 {
				pageReq.Limit = tc.count
			}
			req := &collection.QueryApproversRequest{
				ContractId: tc.contractID,
				Address:    tc.address.String(),
				Pagination: pageReq,
			}
			res, err := s.queryServer.Approvers(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}
