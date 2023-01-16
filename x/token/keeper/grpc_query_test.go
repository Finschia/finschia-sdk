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
		contractID string
		address    sdk.AccAddress
		valid      bool
		postTest   func(res *token.QueryBalanceResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			address:    s.vendor,
			valid:      true,
			postTest: func(res *token.QueryBalanceResponse) {
				s.Require().Equal(s.balance, res.Amount)
			},
		},
		"invalid contract id": {
			address: s.vendor,
		},
		"invalid address": {
			contractID: s.contractID,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryBalanceRequest{
				ContractId: tc.contractID,
				Address:    tc.address.String(),
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
		contractID string
		reqType    string
		valid      bool
		postTest   func(res *token.QuerySupplyResponse)
	}{
		"valid supply request": {
			contractID: s.contractID,
			valid:      true,
			postTest: func(res *token.QuerySupplyResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(3)), res.Amount)
			},
		},
		"invalid contract id": {},
		"no such a contract id": {
			contractID: "fee1dead",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QuerySupplyRequest{
				ContractId: tc.contractID,
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
		valid      bool
		postTest   func(res *token.QueryMintedResponse)
	}{
		"valid mint request": {
			contractID: s.contractID,
			valid:      true,
			postTest: func(res *token.QueryMintedResponse) {
				s.Require().Equal(s.balance.Mul(sdk.NewInt(4)), res.Amount)
			},
		},
		"invalid contract id": {},
		"no such a contract id": {
			contractID: "fee1dead",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryMintedRequest{
				ContractId: tc.contractID,
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
		valid      bool
		postTest   func(res *token.QueryBurntResponse)
	}{
		"valid burn request": {
			contractID: s.contractID,
			valid:      true,
			postTest: func(res *token.QueryBurntResponse) {
				s.Require().Equal(s.balance, res.Amount)
			},
		},
		"invalid contract id": {},
		"no such a contract id": {
			contractID: "fee1dead",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryBurntRequest{
				ContractId: tc.contractID,
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

func (s *KeeperTestSuite) TestQueryContract() {
	// empty request
	_, err := s.queryServer.Contract(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		valid      bool
		postTest   func(res *token.QueryContractResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			valid:      true,
			postTest: func(res *token.QueryContractResponse) {
				s.Require().Equal(s.contractID, res.Contract.Id)
			},
		},
		"invalid contract id": {},
		"no such an id": {
			contractID: "00000000",
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryContractRequest{
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
		contractID string
		valid      bool
		count      uint64
		postTest   func(res *token.QueryContractsResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			valid:      true,
			postTest: func(res *token.QueryContractsResponse) {
				s.Require().Equal(2, len(res.Contracts))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			valid:      true,
			count:      1,
			postTest: func(res *token.QueryContractsResponse) {
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
			req := &token.QueryContractsRequest{
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

func (s *KeeperTestSuite) TestQueryGranteeGrants() {
	// empty request
	_, err := s.queryServer.GranteeGrants(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		valid      bool
		postTest   func(res *token.QueryGranteeGrantsResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			grantee:    s.vendor,
			valid:      true,
			postTest: func(res *token.QueryGranteeGrantsResponse) {
				s.Require().Equal(3, len(res.Grants))
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
			req := &token.QueryGranteeGrantsRequest{
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
		postTest   func(res *token.QueryApprovedResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			address:    s.operator,
			approver:   s.customer,
			valid:      true,
			postTest: func(res *token.QueryApprovedResponse) {
				s.Require().NotNil(res.Approved)
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
			req := &token.QueryApprovedRequest{
				ContractId: tc.contractID,
				Proxy:      tc.address.String(),
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
		operator   sdk.AccAddress
		valid      bool
		count      uint64
		postTest   func(res *token.QueryApproversResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			valid:      true,
			postTest: func(res *token.QueryApproversResponse) {
				s.Require().Equal(2, len(res.Approvers))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			operator:   s.operator,
			valid:      true,
			count:      1,
			postTest: func(res *token.QueryApproversResponse) {
				s.Require().Equal(1, len(res.Approvers))
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
			req := &token.QueryApproversRequest{
				ContractId: tc.contractID,
				Address:    tc.operator.String(),
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
