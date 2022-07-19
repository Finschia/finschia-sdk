package keeper_test

import (
	"github.com/gogo/protobuf/proto"
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

func (s *KeeperTestSuite) TestQuerySupply() {
	// empty request
	_, err := s.queryServer.Supply(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		classID    string
		valid      bool
		postTest   func(res *collection.QuerySupplyResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			classID:    s.ftClassID,
			valid:      true,
			postTest: func(res *collection.QuerySupplyResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(3)), res.Supply)
			},
		},
		"invalid contract id": {
			classID: s.ftClassID,
		},
		"invalid class id": {
			contractID: s.contractID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QuerySupplyRequest{
				ContractId: tc.contractID,
				ClassId:    tc.classID,
			}
			res, err := s.queryServer.Supply(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryMinted() {
	// empty request
	_, err := s.queryServer.Minted(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		classID    string
		valid      bool
		postTest   func(res *collection.QueryMintedResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			classID:    s.ftClassID,
			valid:      true,
			postTest: func(res *collection.QueryMintedResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(6)), res.Minted)
			},
		},
		"invalid contract id": {
			classID: s.ftClassID,
		},
		"invalid class id": {
			contractID: s.contractID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryMintedRequest{
				ContractId: tc.contractID,
				ClassId:    tc.classID,
			}
			res, err := s.queryServer.Minted(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryBurnt() {
	// empty request
	_, err := s.queryServer.Burnt(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		classID    string
		valid      bool
		postTest   func(res *collection.QueryBurntResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			classID:    s.ftClassID,
			valid:      true,
			postTest: func(res *collection.QueryBurntResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(3)), res.Burnt)
			},
		},
		"invalid contract id": {
			classID: s.ftClassID,
		},
		"invalid class id": {
			contractID: s.contractID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryBurntRequest{
				ContractId: tc.contractID,
				ClassId:    tc.classID,
			}
			res, err := s.queryServer.Burnt(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryContracts() {
	// empty request
	_, err := s.queryServer.Contracts(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		valid    bool
		count    uint64
		postTest func(res *collection.QueryContractsResponse)
	}{
		"valid request": {
			valid: true,
			postTest: func(res *collection.QueryContractsResponse) {
				s.Require().Equal(1, len(res.Contracts))
			},
		},
		"valid request with limit": {
			valid: true,
			count: 1,
			postTest: func(res *collection.QueryContractsResponse) {
				s.Require().Equal(1, len(res.Contracts))
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			pageReq := &query.PageRequest{}
			if tc.count != 0 {
				pageReq.Limit = tc.count
			}
			req := &collection.QueryContractsRequest{
				Pagination: pageReq,
			}
			res, err := s.queryServer.Contracts(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryFTClass() {
	// empty request
	_, err := s.queryServer.FTClass(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		classID    string
		valid      bool
		postTest   func(res *collection.QueryFTClassResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			classID:    s.ftClassID,
			valid:      true,
			postTest: func(res *collection.QueryFTClassResponse) {
				s.Require().Equal(s.ftClassID, res.Class.GetId())
			},
		},
		"invalid contract id": {
			classID: s.ftClassID,
		},
		"invalid class id": {
			contractID: s.contractID,
		},
		"no such a class": {
			contractID: s.contractID,
			classID:    "deadbeef",
		},
		"not a class of ft": {
			contractID: s.contractID,
			classID:    s.nftClassID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryFTClassRequest{
				ContractId: tc.contractID,
				ClassId:    tc.classID,
			}
			res, err := s.queryServer.FTClass(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryNFTClass() {
	// empty request
	_, err := s.queryServer.NFTClass(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		classID    string
		valid      bool
		postTest   func(res *collection.QueryNFTClassResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			classID:    s.nftClassID,
			valid:      true,
			postTest: func(res *collection.QueryNFTClassResponse) {
				s.Require().Equal(s.nftClassID, res.Class.GetId())
			},
		},
		"invalid contract id": {
			classID: s.ftClassID,
		},
		"invalid class id": {
			contractID: s.contractID,
		},
		"no such a class": {
			contractID: s.contractID,
			classID:    "deadbeef",
		},
		"not a class of nft": {
			contractID: s.contractID,
			classID:    s.ftClassID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryNFTClassRequest{
				ContractId: tc.contractID,
				ClassId:    tc.classID,
			}
			res, err := s.queryServer.NFTClass(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryTokenClassTypeName() {
	// empty request
	_, err := s.queryServer.TokenClassTypeName(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		classID    string
		valid      bool
		postTest   func(res *collection.QueryTokenClassTypeNameResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			classID:    s.ftClassID,
			valid:      true,
			postTest: func(res *collection.QueryTokenClassTypeNameResponse) {
				s.Require().Equal(proto.MessageName(&collection.FTClass{}), res.Name)
			},
		},
		"invalid contract id": {
			classID: s.ftClassID,
		},
		"invalid class id": {
			contractID: s.contractID,
		},
		"no such a class": {
			contractID: s.contractID,
			classID:    "00bab10c",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryTokenClassTypeNameRequest{
				ContractId: tc.contractID,
				ClassId:    tc.classID,
			}
			res, err := s.queryServer.TokenClassTypeName(s.goCtx, req)
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

// func (s *KeeperTestSuite) TestQueryTokenClasses() {
// 	// empty request
// 	_, err := s.queryServer.TokenClasses(s.goCtx, nil)
// 	s.Require().Error(err)

// 	testCases := map[string]struct {
// 		contractID string
// 		valid      bool
// 		count      uint64
// 		postTest   func(res *collection.QueryTokenClassesResponse)
// 	}{
// 		"valid request": {
// 			contractID: s.contractID,
// 			valid:      true,
// 			postTest: func(res *collection.QueryTokenClassesResponse) {
// 				s.Require().Equal(2, len(res.Classes))
// 			},
// 		},
// 		"valid request with limit": {
// 			contractID: s.contractID,
// 			valid:      true,
// 			count:      1,
// 			postTest: func(res *collection.QueryTokenClassesResponse) {
// 				s.Require().Equal(1, len(res.Classes))
// 			},
// 		},
// 		"invalid contract id": {},
// 	}

// 	for name, tc := range testCases {
// 		s.Run(name, func() {
// 			pageReq := &query.PageRequest{}
// 			if tc.count != 0 {
// 				pageReq.Limit = tc.count
// 			}
// 			req := &collection.QueryTokenClassesRequest{
// 				ContractId: tc.contractID,
// 				Pagination: pageReq,
// 			}
// 			res, err := s.queryServer.TokenClasses(s.goCtx, req)
// 			if !tc.valid {
// 				s.Require().Error(err)
// 				return
// 			}
// 			s.Require().NoError(err)
// 			s.Require().NotNil(res)
// 			tc.postTest(res)
// 		})
// 	}
// }

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
			},
		},
		"valid nft request": {
			contractID: s.contractID,
			tokenID:    nftTokenID,
			valid:      true,
			postTest: func(res *collection.QueryTokenResponse) {
				s.Require().Equal("/lbm.collection.v1.OwnerNFT", res.Token.TypeUrl)
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

func (s *KeeperTestSuite) TestQueryNFT() {
	// empty request
	_, err := s.queryServer.NFT(s.goCtx, nil)
	s.Require().Error(err)

	tokenID := collection.NewNFTID(s.nftClassID, 1)
	testCases := map[string]struct {
		contractID string
		tokenID    string
		valid      bool
		postTest   func(res *collection.QueryNFTResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			valid:      true,
			postTest: func(res *collection.QueryNFTResponse) {
				s.Require().Equal(tokenID, res.Token.Id)
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
			req := &collection.QueryNFTRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.NFT(s.goCtx, req)
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

// func (s *KeeperTestSuite) TestQueryNFTs() {
// 	// empty request
// 	_, err := s.queryServer.NFTs(s.goCtx, nil)
// 	s.Require().Error(err)

// 	testCases := map[string]struct {
// 		contractID string
// 		valid      bool
// 		count      uint64
// 		postTest   func(res *collection.QueryNFTsResponse)
// 	}{
// 		"valid request": {
// 			contractID: s.contractID,
// 			valid:      true,
// 			count:      1000000,
// 			postTest: func(res *collection.QueryNFTsResponse) {
// 				s.Require().Equal(s.lenChain*6, len(res.Tokens))
// 			},
// 		},
// 		"valid request with limit": {
// 			contractID: s.contractID,
// 			valid:      true,
// 			count:      1,
// 			postTest: func(res *collection.QueryNFTsResponse) {
// 				s.Require().Equal(1, len(res.Tokens))
// 			},
// 		},
// 		"invalid contract id": {},
// 	}

// 	for name, tc := range testCases {
// 		s.Run(name, func() {
// 			pageReq := &query.PageRequest{}
// 			if tc.count != 0 {
// 				pageReq.Limit = tc.count
// 			}
// 			req := &collection.QueryNFTsRequest{
// 				ContractId: tc.contractID,
// 				Pagination: pageReq,
// 			}
// 			res, err := s.queryServer.NFTs(s.goCtx, req)
// 			if !tc.valid {
// 				s.Require().Error(err)
// 				return
// 			}
// 			s.Require().NoError(err)
// 			s.Require().NotNil(res)
// 			tc.postTest(res)
// 		})
// 	}
// }

func (s *KeeperTestSuite) TestQueryOwner() {
	// empty request
	_, err := s.queryServer.Owner(s.goCtx, nil)
	s.Require().Error(err)

	tokenID := collection.NewNFTID(s.nftClassID, 1)
	testCases := map[string]struct {
		contractID string
		tokenID    string
		valid      bool
		postTest   func(res *collection.QueryOwnerResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			tokenID:    tokenID,
			valid:      true,
			postTest: func(res *collection.QueryOwnerResponse) {
				s.Require().Equal(s.customer.String(), res.Owner)
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
			req := &collection.QueryOwnerRequest{
				ContractId: tc.contractID,
				TokenId:    tc.tokenID,
			}
			res, err := s.queryServer.Owner(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryGrant() {
	// empty request
	_, err := s.queryServer.Grant(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		permission collection.Permission
		valid      bool
		postTest   func(res *collection.QueryGrantResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			grantee:    s.vendor,
			permission: collection.PermissionModify,
			valid:      true,
			postTest: func(res *collection.QueryGrantResponse) {
				s.Require().Equal(s.vendor.String(), res.Grant.Grantee)
				s.Require().Equal(collection.PermissionModify, res.Grant.Permission)
			},
		},
		"invalid contract id": {
			grantee:    s.vendor,
			permission: collection.PermissionModify,
		},
		"invalid grantee": {
			contractID: s.contractID,
			permission: collection.PermissionModify,
		},
		"invalid permission": {
			contractID: s.contractID,
			grantee:    s.vendor,
		},
		"no permission": {
			contractID: s.contractID,
			grantee:    s.customer,
			permission: collection.PermissionModify,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryGrantRequest{
				ContractId: tc.contractID,
				Grantee:    tc.grantee.String(),
				Permission: tc.permission,
			}
			res, err := s.queryServer.Grant(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryAuthorization() {
	// empty request
	_, err := s.queryServer.Authorization(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		holder     sdk.AccAddress
		valid      bool
		postTest   func(res *collection.QueryAuthorizationResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			holder:     s.customer,
			valid:      true,
			postTest: func(res *collection.QueryAuthorizationResponse) {
				expected := collection.Authorization{
					Holder:   s.customer.String(),
					Operator: s.operator.String(),
				}
				s.Require().Equal(expected, res.Authorization)
			},
		},
		"invalid contract id": {
			operator: s.operator,
			holder:   s.customer,
		},
		"invalid operator": {
			contractID: s.contractID,
			holder:     s.customer,
		},
		"invalid holder": {
			contractID: s.contractID,
			operator:   s.operator,
		},
		"no authorization found": {
			contractID: s.contractID,
			operator:   s.vendor,
			holder:     s.customer,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &collection.QueryAuthorizationRequest{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				Holder:     tc.holder.String(),
			}
			res, err := s.queryServer.Authorization(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryOperatorAuthorizations() {
	// empty request
	_, err := s.queryServer.OperatorAuthorizations(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		valid      bool
		count      uint64
		postTest   func(res *collection.QueryOperatorAuthorizationsResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			valid:      true,
			postTest: func(res *collection.QueryOperatorAuthorizationsResponse) {
				s.Require().Equal(1, len(res.Authorizations))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			operator:   s.operator,
			valid:      true,
			count:      1,
			postTest: func(res *collection.QueryOperatorAuthorizationsResponse) {
				s.Require().Equal(1, len(res.Authorizations))
			},
		},
		"invalid contract id": {
			operator: s.operator,
		},
		"invalid operator": {
			contractID: s.contractID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			pageReq := &query.PageRequest{}
			if tc.count != 0 {
				pageReq.Limit = tc.count
			}
			req := &collection.QueryOperatorAuthorizationsRequest{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				Pagination: pageReq,
			}
			res, err := s.queryServer.OperatorAuthorizations(s.goCtx, req)
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
