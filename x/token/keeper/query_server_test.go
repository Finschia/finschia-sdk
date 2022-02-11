package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestQueryBalance() {
	testCases := map[string]struct{
		classId string
		address sdk.AccAddress
		valid bool
		postTest func(res *token.QueryBalanceResponse)
	}{
		"valid request": {
			classId: s.mintableClass,
			address: s.vendor,
			valid: true,
			postTest: func(res *token.QueryBalanceResponse) {
				s.Require().Equal(s.balance, res.Amount)
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryBalanceRequest{
				ClassId: tc.classId,
				Address: tc.address.String(),
			}
			res, err := s.queryServer.Balance(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQuerySupply() {
	testCases := map[string]struct{
		classId string
		reqType string
		valid bool
		postTest func(res *token.QuerySupplyResponse)
	}{
		"valid supply request": {
			classId: s.mintableClass,
			reqType: "supply",
			valid: true,
			postTest: func(res *token.QuerySupplyResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(3)), res.Amount)
			},
		},
		"valid mint request": {
			classId: s.mintableClass,
			reqType: "mint",
			valid: true,
			postTest: func(res *token.QuerySupplyResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(3)), res.Amount)
			},
		},
		"valid burn request": {
			classId: s.mintableClass,
			reqType: "burn",
			valid: true,
			postTest: func(res *token.QuerySupplyResponse) {
				s.Require().Equal(sdk.ZeroInt(), res.Amount)
			},
		},
		"invalid request": {
			classId: s.mintableClass,
			reqType: "invalid",
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QuerySupplyRequest{
				ClassId: tc.classId,
				Type: tc.reqType,
			}
			res, err := s.queryServer.Supply(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryToken() {
	testCases := map[string]struct{
		classId string
		valid bool
		postTest func(res *token.QueryTokenResponse)
	}{
		"valid request": {
			classId: s.mintableClass,
			valid: true,
			postTest: func(res *token.QueryTokenResponse) {
				s.Require().Equal(s.mintableClass, res.Token.Id)
			},
		},
		"invalid request": {
			classId: "invalid",
			valid: false,
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
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

func (s *KeeperTestSuite) TestQueryTokens() {
	testCases := map[string]struct{
		classId string
		count uint64
		valid bool
		postTest func(res *token.QueryTokensResponse)
	}{
		"valid request": {
			classId: s.mintableClass,
			valid: true,
			postTest: func(res *token.QueryTokensResponse) {
				s.Require().Equal(2, len(res.Tokens))
			},
		},
		"valid request with limit": {
			classId: s.mintableClass,
			count: 1,
			valid: true,
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
			s.Require().NotNil(res)
			tc.postTest(res)
		})
	}
}

