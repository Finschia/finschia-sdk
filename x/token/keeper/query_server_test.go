package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestQueryBalance() {
	// empty request
	_, err := s.queryServer.Balance(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		classId  string
		address  sdk.AccAddress
		valid    bool
		postTest func(res *token.QueryBalanceResponse)
	}{
		"valid request": {
			classId: s.classID,
			address: s.vendor,
			valid:   true,
			postTest: func(res *token.QueryBalanceResponse) {
				s.Require().Equal(s.balance, res.Amount)
			},
		},
		"invalid class id": {
			address: s.vendor,
		},
		"invalid address": {
			classId: s.classID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryBalanceRequest{
				ContractId: tc.classId,
				Address: tc.address.String(),
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

func (s *KeeperTestSuite) TestQuerySupply() {
	// empty request
	_, err := s.queryServer.Supply(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		classId  string
		reqType  string
		valid    bool
		postTest func(res *token.QuerySupplyResponse)
	}{
		"valid supply request": {
			classId: s.classID,
			valid:   true,
			postTest: func(res *token.QuerySupplyResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(3)), res.Amount)
			},
		},
		"invalid class id": {
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QuerySupplyRequest{
				ContractId: tc.classId,
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
		classId  string
		valid    bool
		postTest func(res *token.QueryMintedResponse)
	}{
		"valid mint request": {
			classId: s.classID,
			valid:   true,
			postTest: func(res *token.QueryMintedResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(4)), res.Amount)
			},
		},
		"invalid class id": {
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryMintedRequest{
				ContractId: tc.classId,
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
		classId  string
		valid    bool
		postTest func(res *token.QueryBurntResponse)
	}{
		"valid burn request": {
			classId: s.classID,
			valid:   true,
			postTest: func(res *token.QueryBurntResponse) {
				s.Require().Equal(s.balance, res.Amount)
			},
		},
		"invalid class id": {
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryBurntRequest{
				ContractId: tc.classId,
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

func (s *KeeperTestSuite) TestQueryTokenClass() {
	// empty request
	_, err := s.queryServer.TokenClass(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		classId  string
		valid    bool
		postTest func(res *token.QueryTokenClassResponse)
	}{
		"valid request": {
			classId: s.classID,
			valid:   true,
			postTest: func(res *token.QueryTokenClassResponse) {
				s.Require().Equal(s.classID, res.Class.ContractId)
			},
		},
		"invalid class id": {
		},
		"no such an id": {
			classId: "00000000",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryTokenClassRequest{
				ContractId: tc.classId,
			}
			res, err := s.queryServer.TokenClass(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryTokenClasses() {
	// empty request
	_, err := s.queryServer.TokenClasses(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		classId  string
		valid    bool
		count    uint64
		postTest func(res *token.QueryTokenClassesResponse)
	}{
		"valid request": {
			classId: s.classID,
			valid:   true,
			postTest: func(res *token.QueryTokenClassesResponse) {
				s.Require().Equal(2, len(res.Classes))
			},
		},
		"valid request with limit": {
			classId: s.classID,
			valid:   true,
			count:   1,
			postTest: func(res *token.QueryTokenClassesResponse) {
				s.Require().Equal(1, len(res.Classes))
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			pageReq := &query.PageRequest{}
			if tc.count != 0 {
				pageReq.Limit = tc.count
			}
			req := &token.QueryTokenClassesRequest{
				Pagination: pageReq,
			}
			res, err := s.queryServer.TokenClasses(s.goCtx, req)
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
		classId  string
		grantee  sdk.AccAddress
		valid    bool
		postTest func(res *token.QueryGranteeGrantsResponse)
	}{
		"valid request": {
			classId: s.classID,
			grantee: s.vendor,
			valid:   true,
			postTest: func(res *token.QueryGranteeGrantsResponse) {
				s.Require().Equal(3, len(res.Grants))
			},
		},
		"invalid class id": {
			grantee: s.vendor,
		},
		"invalid grantee": {
			classId: s.classID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryGranteeGrantsRequest{
				ContractId: tc.classId,
				Grantee: tc.grantee.String(),
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
		classId  string
		proxy    sdk.AccAddress
		approver sdk.AccAddress
		valid    bool
		postTest func(res *token.QueryAuthorizationResponse)
	}{
		"valid request": {
			classId:  s.classID,
			proxy:    s.operator,
			approver: s.customer,
			valid:    true,
			postTest: func(res *token.QueryAuthorizationResponse) {
				s.Require().NotNil(res.Authorization)
			},
		},
		"invalid class id": {
			proxy:    s.operator,
			approver: s.customer,
		},
		"invalid proxy": {
			classId:  s.classID,
			approver: s.customer,
		},
		"invalid approver": {
			classId:  s.classID,
			proxy:    s.operator,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryAuthorizationRequest{
				ContractId:  tc.classId,
				Proxy:    tc.proxy.String(),
				Approver: tc.approver.String(),
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
		classId  string
		proxy    sdk.AccAddress
		valid    bool
		count    uint64
		postTest func(res *token.QueryOperatorAuthorizationsResponse)
	}{
		"valid request": {
			classId: s.classID,
			proxy:   s.operator,
			valid:   true,
			postTest: func(res *token.QueryOperatorAuthorizationsResponse) {
				s.Require().Equal(2, len(res.Authorizations))
			},
		},
		"valid request with limit": {
			classId: s.classID,
			proxy:   s.operator,
			valid:   true,
			count:   1,
			postTest: func(res *token.QueryOperatorAuthorizationsResponse) {
				s.Require().Equal(1, len(res.Authorizations))
			},
		},
		"invalid class id": {
			proxy:   s.operator,
		},
		"invalid proxy": {
			classId: s.classID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			pageReq := &query.PageRequest{}
			if tc.count != 0 {
				pageReq.Limit = tc.count
			}
			req := &token.QueryOperatorAuthorizationsRequest{
				ContractId:    tc.classId,
				Proxy:      tc.proxy.String(),
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
