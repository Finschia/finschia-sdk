package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestQueryTokenBalance() {
	// empty request
	_, err := s.queryServer.TokenBalance(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		classId  string
		address  sdk.AccAddress
		valid    bool
		postTest func(res *token.QueryTokenBalanceResponse)
	}{
		"valid request": {
			classId: s.classID,
			address: s.vendor,
			valid:   true,
			postTest: func(res *token.QueryTokenBalanceResponse) {
				s.Require().Equal(s.balance, res.Amount)
			},
		},
		"invalid class id": {
			classId: "invalid",
			address: s.vendor,
			valid:   false,
		},
		"invalid address": {
			classId: s.classID,
			address: "invalid",
			valid:   false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryTokenBalanceRequest{
				ClassId: tc.classId,
				Address: tc.address.String(),
			}
			res, err := s.queryServer.TokenBalance(s.goCtx, req)
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
			reqType: "supply",
			valid:   true,
			postTest: func(res *token.QuerySupplyResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(3)), res.Amount)
			},
		},
		"valid mint request": {
			classId: s.classID,
			reqType: "mint",
			valid:   true,
			postTest: func(res *token.QuerySupplyResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(4)), res.Amount)
			},
		},
		"valid burn request": {
			classId: s.classID,
			reqType: "burn",
			valid:   true,
			postTest: func(res *token.QuerySupplyResponse) {
				s.Require().Equal(s.balance, res.Amount)
			},
		},
		"invalid class id": {
			classId: "invalid",
			reqType: "burn",
			valid:   false,
		},
		"invalid request": {
			classId: s.classID,
			reqType: "invalid",
			valid:   false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QuerySupplyRequest{
				ClassId: tc.classId,
				Type:    tc.reqType,
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

func (s *KeeperTestSuite) TestQueryToken() {
	// empty request
	_, err := s.queryServer.Token(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		classId  string
		valid    bool
		postTest func(res *token.QueryTokenResponse)
	}{
		"valid request": {
			classId: s.classID,
			valid:   true,
			postTest: func(res *token.QueryTokenResponse) {
				s.Require().Equal(s.classID, res.Token.Id)
			},
		},
		"invalid class id": {
			classId: "invalid",
			valid:   false,
		},
		"no such an id": {
			classId: "00000000",
			valid:   false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryTokenRequest{
				ClassId: tc.classId,
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
		classId  string
		valid    bool
		count    uint64
		postTest func(res *token.QueryTokensResponse)
	}{
		"valid request": {
			classId: s.classID,
			valid:   true,
			postTest: func(res *token.QueryTokensResponse) {
				s.Require().Equal(2, len(res.Tokens))
			},
		},
		"valid request with limit": {
			classId: s.classID,
			valid:   true,
			count:   1,
			postTest: func(res *token.QueryTokensResponse) {
				s.Require().Equal(1, len(res.Tokens))
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			pageReq := &query.PageRequest{}
			if tc.count != 0 {
				pageReq.Limit = tc.count
			}
			req := &token.QueryTokensRequest{
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

func (s *KeeperTestSuite) TestQueryGrants() {
	// empty request
	_, err := s.queryServer.TokenGrants(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		classId  string
		grantee  sdk.AccAddress
		valid    bool
		postTest func(res *token.QueryGrantsResponse)
	}{
		"valid request": {
			classId: s.classID,
			grantee: s.vendor,
			valid:   true,
			postTest: func(res *token.QueryGrantsResponse) {
				s.Require().Equal(3, len(res.Grants))
			},
		},
		"invalid class id": {
			classId: "invalid",
			grantee: s.vendor,
			valid:   false,
		},
		"invalid grantee": {
			classId: s.classID,
			grantee: "invalid",
			valid:   false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryGrantsRequest{
				ClassId: tc.classId,
				Grantee: tc.grantee.String(),
			}
			res, err := s.queryServer.TokenGrants(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryApprove() {
	// empty request
	_, err := s.queryServer.Approve(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		classId  string
		proxy    sdk.AccAddress
		approver sdk.AccAddress
		valid    bool
		postTest func(res *token.QueryApproveResponse)
	}{
		"valid request": {
			classId:  s.classID,
			proxy:    s.operator,
			approver: s.customer,
			valid:    true,
			postTest: func(res *token.QueryApproveResponse) {
				s.Require().NotNil(res.Approve)
			},
		},
		"invalid class id": {
			classId:  "invalid",
			proxy:    s.operator,
			approver: s.customer,
			valid:    false,
		},
		"invalid proxy": {
			classId:  s.classID,
			proxy:    "invalid",
			approver: s.customer,
			valid:    false,
		},
		"invalid approver": {
			classId:  s.classID,
			proxy:    s.operator,
			approver: "invalid",
			valid:    false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryApproveRequest{
				ClassId:  tc.classId,
				Proxy:    tc.proxy.String(),
				Approver: tc.approver.String(),
			}
			res, err := s.queryServer.Approve(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryApproves() {
	// empty request
	_, err := s.queryServer.Approves(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		classId  string
		proxy    sdk.AccAddress
		valid    bool
		count    uint64
		postTest func(res *token.QueryApprovesResponse)
	}{
		"valid request": {
			classId: s.classID,
			proxy:   s.operator,
			valid:   true,
			postTest: func(res *token.QueryApprovesResponse) {
				s.Require().Equal(2, len(res.Approves))
			},
		},
		"valid request with limit": {
			classId: s.classID,
			proxy:   s.operator,
			valid:   true,
			count:   1,
			postTest: func(res *token.QueryApprovesResponse) {
				s.Require().Equal(1, len(res.Approves))
			},
		},
		"invalid class id": {
			classId: "invalid",
			proxy:   s.operator,
			valid:   false,
		},
		"invalid proxy": {
			classId: s.classID,
			proxy:   "invalid",
			valid:   false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			pageReq := &query.PageRequest{}
			if tc.count != 0 {
				pageReq.Limit = tc.count
			}
			req := &token.QueryApprovesRequest{
				ClassId:    tc.classId,
				Proxy:      tc.proxy.String(),
				Pagination: pageReq,
			}
			res, err := s.queryServer.Approves(s.goCtx, req)
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
