package keeper_test

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/query"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestQueryBalance() {
	// empty request
	_, err := s.queryServer.Balance(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		address    sdk.AccAddress
		code       codes.Code
		postTest   func(res *token.QueryBalanceResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			address:    s.vendor,
			postTest: func(res *token.QueryBalanceResponse) {
				s.Require().Equal(s.balance, res.Amount)
			},
		},
		"invalid contract id": {
			address: s.vendor,
			code:    codes.InvalidArgument,
		},
		"invalid address": {
			contractID: s.contractID,
			code:       codes.InvalidArgument,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryBalanceRequest{
				ContractId: tc.contractID,
				Address:    tc.address.String(),
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

func (s *KeeperTestSuite) TestQuerySupply() {
	// empty request
	_, err := s.queryServer.Supply(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		reqType    string
		code       codes.Code
		postTest   func(res *token.QuerySupplyResponse)
	}{
		"valid supply request": {
			contractID: s.contractID,
			postTest: func(res *token.QuerySupplyResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(3)), res.Amount)
			},
		},
		"invalid contract id": {
			code: codes.InvalidArgument,
		},
		"no such a contract id": {
			contractID: "fee1dead",
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QuerySupplyRequest{
				ContractId: tc.contractID,
			}
			res, err := s.queryServer.Supply(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryMinted() {
	// empty request
	_, err := s.queryServer.Minted(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		code       codes.Code
		postTest   func(res *token.QueryMintedResponse)
	}{
		"valid mint request": {
			contractID: s.contractID,
			postTest: func(res *token.QueryMintedResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(4)), res.Amount)
			},
		},
		"invalid contract id": {
			code: codes.InvalidArgument,
		},
		"no such a contract id": {
			contractID: "fee1dead",
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryMintedRequest{
				ContractId: tc.contractID,
			}
			res, err := s.queryServer.Minted(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryBurnt() {
	// empty request
	_, err := s.queryServer.Burnt(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		code       codes.Code
		postTest   func(res *token.QueryBurntResponse)
	}{
		"valid burn request": {
			contractID: s.contractID,
			postTest: func(res *token.QueryBurntResponse) {
				s.Require().Equal(s.balance, res.Amount)
			},
		},
		"invalid contract id": {
			code: codes.InvalidArgument,
		},
		"no such a contract id": {
			contractID: "fee1dead",
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryBurntRequest{
				ContractId: tc.contractID,
			}
			res, err := s.queryServer.Burnt(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryTokenClass() {
	// empty request
	_, err := s.queryServer.TokenClass(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		code       codes.Code
		postTest   func(res *token.QueryTokenClassResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			postTest: func(res *token.QueryTokenClassResponse) {
				s.Require().Equal(s.contractID, res.Class.ContractId)
			},
		},
		"invalid contract id": {
			code: codes.InvalidArgument,
		},
		"no such an id": {
			contractID: "00000000",
			code:       codes.NotFound,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryTokenClassRequest{
				ContractId: tc.contractID,
			}
			res, err := s.queryServer.TokenClass(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryTokenClasses() {
	// empty request
	_, err := s.queryServer.TokenClasses(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		code       codes.Code
		count      uint64
		postTest   func(res *token.QueryTokenClassesResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			postTest: func(res *token.QueryTokenClassesResponse) {
				s.Require().Equal(2, len(res.Classes))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			count:      1,
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
		postTest   func(res *token.QueryGranteeGrantsResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			grantee:    s.vendor,
			postTest: func(res *token.QueryGranteeGrantsResponse) {
				s.Require().Equal(3, len(res.Grants))
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
			req := &token.QueryGranteeGrantsRequest{
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
		postTest   func(res *token.QueryApprovedResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			address:    s.operator,
			approver:   s.customer,
			postTest: func(res *token.QueryApprovedResponse) {
				s.Require().NotNil(res.Approved)
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
			req := &token.QueryApprovedRequest{
				ContractId: tc.contractID,
				Proxy:      tc.address.String(),
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
		operator   sdk.AccAddress
		code       codes.Code
		count      uint64
		postTest   func(res *token.QueryApproversResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			postTest: func(res *token.QueryApproversResponse) {
				s.Require().Equal(2, len(res.Approvers))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			operator:   s.operator,
			count:      1,
			postTest: func(res *token.QueryApproversResponse) {
				s.Require().Equal(1, len(res.Approvers))
			},
		},
		"invalid contract id": {
			operator: s.operator,
			code:     codes.InvalidArgument,
		},
		"invalid operator": {
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
			req := &token.QueryApproversRequest{
				ContractId: tc.contractID,
				Address:    tc.operator.String(),
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
