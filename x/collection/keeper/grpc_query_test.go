package keeper_test

import (
	"github.com/gogo/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		code       codes.Code
		postTest   func(res *collection.QueryBalanceResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			address:    s.vendor,
			tokenID:    tokenID,
			postTest: func(res *collection.QueryBalanceResponse) {
				expected := collection.NewCoin(tokenID, s.balance)
				s.Require().Equal(expected, res.Balance)
			},
		},
		"invalid contract id": {
			address: s.vendor,
			tokenID: tokenID,
			code:    codes.InvalidArgument,
		},
		"invalid address": {
			contractID: s.contractID,
			tokenID:    tokenID,
			code:       codes.InvalidArgument,
		},
		"invalid token id": {
			contractID: s.contractID,
			address:    s.vendor,
			code:       codes.InvalidArgument,
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
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		count      uint64
		postTest   func(res *collection.QueryAllBalancesResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			address:    s.customer,
			postTest: func(res *collection.QueryAllBalancesResponse) {
				s.Require().Equal(s.numNFTs+1, len(res.Balances))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			address:    s.customer,
			count:      1,
			postTest: func(res *collection.QueryAllBalancesResponse) {
				s.Require().Equal(1, len(res.Balances))
			},
		},
		"invalid contract id": {
			address: s.customer,
			code:    codes.InvalidArgument,
		},
		"invalid address": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
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
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryFTSupplyResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			postTest: func(res *collection.QueryFTSupplyResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(3)), res.Supply)
			},
		},
		"invalid contract id": {
			tokenID: tokenID,
			code:    codes.InvalidArgument,
		},
		"invalid token id": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
		"no such a token": {
			contractID: s.contractID,
			tokenID:    collection.NewFTID("00bab10c"),
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryFTSupplyRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.FTSupply(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryFTMintedResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			postTest: func(res *collection.QueryFTMintedResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(6)), res.Minted)
			},
		},
		"invalid contract id": {
			tokenID: tokenID,
			code:    codes.InvalidArgument,
		},
		"invalid token id": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
		"no such a token": {
			contractID: s.contractID,
			tokenID:    collection.NewFTID("00bab10c"),
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryFTMintedRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.FTMinted(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryFTBurntResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			postTest: func(res *collection.QueryFTBurntResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(3)), res.Burnt)
			},
		},
		"invalid contract id": {
			tokenID: tokenID,
			code:    codes.InvalidArgument,
		},
		"invalid token id": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
		"no such a token": {
			contractID: s.contractID,
			tokenID:    collection.NewFTID("00bab10c"),
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryFTBurntRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.FTBurnt(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryNFTSupplyResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenType:  s.nftClassID,
			postTest: func(res *collection.QueryNFTSupplyResponse) {
				s.Require().EqualValues(s.numNFTs*3, res.Supply.Int64())
			},
		},
		"invalid contract id": {
			tokenType: s.nftClassID,
			code:      codes.InvalidArgument,
		},
		"invalid token type": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
		"no such a token type": {
			contractID: s.contractID,
			tokenType:  "deadbeef",
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryNFTSupplyRequest{
				ContractId: tc.contractID,
				TokenType:  tc.tokenType,
			}
			res, err := s.queryServer.NFTSupply(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryNFTMintedResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenType:  s.nftClassID,
			postTest: func(res *collection.QueryNFTMintedResponse) {
				s.Require().EqualValues(s.numNFTs*3, res.Minted.Int64())
			},
		},
		"invalid contract id": {
			tokenType: s.nftClassID,
			code:      codes.InvalidArgument,
		},
		"invalid token type": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
		"no such a token type": {
			contractID: s.contractID,
			tokenType:  "deadbeef",
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryNFTMintedRequest{
				ContractId: tc.contractID,
				TokenType:  tc.tokenType,
			}
			res, err := s.queryServer.NFTMinted(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryNFTBurntResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenType:  s.nftClassID,
			postTest: func(res *collection.QueryNFTBurntResponse) {
				s.Require().Equal(sdk.ZeroInt(), res.Burnt)
			},
		},
		"invalid contract id": {
			tokenType: s.nftClassID,
			code:      codes.InvalidArgument,
		},
		"invalid token type": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
		"no such a token type": {
			contractID: s.contractID,
			tokenType:  "deadbeef",
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryNFTBurntRequest{
				ContractId: tc.contractID,
				TokenType:  tc.tokenType,
			}
			res, err := s.queryServer.NFTBurnt(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryContractResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			postTest: func(res *collection.QueryContractResponse) {
				s.Require().Equal(s.contractID, res.Contract.ContractId)
			},
		},
		"invalid contract id": {
			code: codes.InvalidArgument,
		},
		"no such an id": {
			contractID: "deadbeef",
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryContractRequest{
				ContractId: tc.contractID,
			}
			res, err := s.queryServer.Contract(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
			s.Require().NotNil(res)

			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryTokenClassTypeName() {
	// empty request
	_, err := s.queryServer.TokenClassTypeName(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		classID    string
		code       codes.Code
		postTest   func(res *collection.QueryTokenClassTypeNameResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			classID:    s.ftClassID,
			postTest: func(res *collection.QueryTokenClassTypeNameResponse) {
				s.Require().Equal(proto.MessageName(&collection.FTClass{}), res.Name)
			},
		},
		"invalid contract id": {
			classID: s.ftClassID,
			code:    codes.InvalidArgument,
		},
		"invalid class id": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
		"no such a class": {
			contractID: s.contractID,
			classID:    "00bab10c",
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryTokenClassTypeNameRequest{
				ContractId: tc.contractID,
				ClassId:    tc.classID,
			}
			res, err := s.queryServer.TokenClassTypeName(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryTokenTypeResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenType:  s.nftClassID,
			postTest: func(res *collection.QueryTokenTypeResponse) {
				s.Require().Equal(s.contractID, res.TokenType.ContractId)
				s.Require().Equal(s.nftClassID, res.TokenType.TokenType)
			},
		},
		"invalid contract id": {
			tokenType: s.nftClassID,
			code:      codes.InvalidArgument,
		},
		"invalid token type": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
		"no such a token type": {
			contractID: s.contractID,
			tokenType:  "deadbeef",
			code:       codes.NotFound,
		},
		"not a class of nft": {
			contractID: s.contractID,
			tokenType:  s.ftClassID,
			code:       codes.InvalidArgument,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryTokenTypeRequest{
				ContractId: tc.contractID,
				TokenType:  tc.tokenType,
			}
			res, err := s.queryServer.TokenType(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		count      uint64
		postTest   func(res *collection.QueryTokenTypesResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			postTest: func(res *collection.QueryTokenTypesResponse) {
				s.Require().Equal(1, len(res.TokenTypes))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			count:      1,
			postTest: func(res *collection.QueryTokenTypesResponse) {
				s.Require().Equal(1, len(res.TokenTypes))
			},
		},
		"invalid contract id": {
			code: codes.InvalidArgument,
		},
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
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryTokenResponse)
	}{
		"valid ft request": {
			contractID: s.contractID,
			tokenID:    ftTokenID,
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
			code:    codes.InvalidArgument,
		},
		"invalid token id": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
		"no such a fungible token": {
			contractID: s.contractID,
			tokenID:    collection.NewFTID("00bab10c"),
			code:       codes.NotFound,
		},
		"no such a non-fungible token": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID("deadbeef", 1),
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryTokenRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.Token(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		count      uint64
		postTest   func(res *collection.QueryTokensWithTokenTypeResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenType:  s.nftClassID,
			count:      1000000,
			postTest: func(res *collection.QueryTokensWithTokenTypeResponse) {
				s.Require().Equal(s.numNFTs*3, len(res.Tokens))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			tokenType:  s.nftClassID,
			count:      1,
			postTest: func(res *collection.QueryTokensWithTokenTypeResponse) {
				s.Require().Equal(1, len(res.Tokens))
			},
		},
		"invalid contract id": {
			tokenType: s.nftClassID,
			code:      codes.InvalidArgument,
		},
		"invalid token type": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
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
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		count      uint64
		postTest   func(res *collection.QueryTokensResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			count:      1000000,
			postTest: func(res *collection.QueryTokensResponse) {
				s.Require().Equal(s.numNFTs*3+1, len(res.Tokens))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			count:      1,
			postTest: func(res *collection.QueryTokensResponse) {
				s.Require().Equal(1, len(res.Tokens))
			},
		},
		"invalid contract id": {
			code: codes.InvalidArgument,
		},
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
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryRootResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			postTest: func(res *collection.QueryRootResponse) {
				s.Require().Equal(collection.NewNFTID(s.nftClassID, 1), res.Root.Id)
			},
		},
		"invalid contract id": {
			tokenID: tokenID,
			code:    codes.InvalidArgument,
		},
		"invalid token id": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
		"no such a token": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID("deadbeef", 1),
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryRootRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.Root(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryParentResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			postTest: func(res *collection.QueryParentResponse) {
				s.Require().Equal(collection.NewNFTID(s.nftClassID, 1), res.Parent.Id)
			},
		},
		"invalid contract id": {
			tokenID: tokenID,
			code:    codes.InvalidArgument,
		},
		"invalid token id": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
		"no such a token": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID("deadbeef", 1),
			code:       codes.NotFound,
		},
		"no parent": {
			contractID: s.contractID,
			tokenID:    collection.NewNFTID(s.nftClassID, 1),
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryParentRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.Parent(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		count      uint64
		postTest   func(res *collection.QueryChildrenResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			postTest: func(res *collection.QueryChildrenResponse) {
				s.Require().Equal(1, len(res.Children))
				s.Require().Equal(collection.NewNFTID(s.nftClassID, 2), res.Children[0].Id)
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			tokenID:    tokenID,
			count:      1,
			postTest: func(res *collection.QueryChildrenResponse) {
				s.Require().Equal(1, len(res.Children))
				s.Require().Equal(collection.NewNFTID(s.nftClassID, 2), res.Children[0].Id)
			},
		},
		"invalid contract id": {
			tokenID: tokenID,
			code:    codes.InvalidArgument,
		},
		"invalid token id": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
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
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryGranteeGrantsResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			grantee:    s.vendor,
			postTest: func(res *collection.QueryGranteeGrantsResponse) {
				s.Require().Equal(4, len(res.Grants))
			},
		},
		"invalid contract id": {
			grantee: s.vendor,
			code:    codes.InvalidArgument,
		},
		"invalid grantee": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryGranteeGrantsRequest{
				ContractId: tc.contractID,
				Grantee:    tc.grantee.String(),
			}
			res, err := s.queryServer.GranteeGrants(s.goCtx, req)
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		postTest   func(res *collection.QueryApprovedResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			address:    s.operator,
			approver:   s.customer,
			postTest: func(res *collection.QueryApprovedResponse) {
				s.Require().True(res.Approved)
			},
		},
		"invalid contract id": {
			address:  s.operator,
			approver: s.customer,
			code:     codes.InvalidArgument,
		},
		"invalid address": {
			contractID: s.contractID,
			approver:   s.customer,
			code:       codes.InvalidArgument,
		},
		"invalid approver": {
			contractID: s.contractID,
			address:    s.operator,
			code:       codes.InvalidArgument,
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
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
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
		code       codes.Code
		count      uint64
		postTest   func(res *collection.QueryApproversResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			address:    s.operator,
			postTest: func(res *collection.QueryApproversResponse) {
				s.Require().Equal(1, len(res.Approvers))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			address:    s.operator,
			count:      1,
			postTest: func(res *collection.QueryApproversResponse) {
				s.Require().Equal(1, len(res.Approvers))
			},
		},
		"invalid contract id": {
			address: s.operator,
			code:    codes.InvalidArgument,
		},
		"invalid address": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
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
			grpcstatus, ok := status.FromError(err)
			s.Require().True(ok)
			s.Require().Equal(tc.code, grpcstatus.Code())
			if tc.code != codes.OK {
				s.Require().Nil(res)
				return
			}
			s.Require().NotNil(res)

			tc.postTest(res)
		})
	}
}
