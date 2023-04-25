package keeper_test

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/query"
	"github.com/Finschia/finschia-sdk/x/token"
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
		"no such a contract id": {
			contractID: "fee1dead",
			valid:      true,
			postTest: func(res *token.QuerySupplyResponse) {
				s.Require().Equal(sdk.ZeroInt(), res.Amount)
			},
		},
		"invalid contract id": {},
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
		"no such a contract id": {
			contractID: "fee1dead",
			valid:      true,
			postTest: func(res *token.QueryMintedResponse) {
				s.Require().Equal(sdk.ZeroInt(), res.Amount)
			},
		},
		"invalid contract id": {},
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
		"no such a contract id": {
			contractID: "fee1dead",
			valid:      true,
			postTest: func(res *token.QueryBurntResponse) {
				s.Require().Equal(sdk.ZeroInt(), res.Amount)
			},
		},
		"invalid contract id": {},
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
		"class not found": {
			contractID: "fee1dead",
			grantee:    s.vendor,
			valid:      true,
			postTest: func(res *token.QueryGranteeGrantsResponse) {
				s.Require().Equal(0, len(res.Grants))
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

func (s *KeeperTestSuite) TestQueryIsOperatorFor() {
	// empty request
	_, err := s.queryServer.IsOperatorFor(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		holder     sdk.AccAddress
		valid      bool
		postTest   func(res *token.QueryIsOperatorForResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			holder:     s.customer,
			valid:      true,
			postTest: func(res *token.QueryIsOperatorForResponse) {
				s.Require().True(res.Authorized)
			},
		},
		"class not found": {
			contractID: "fee1dead",
			operator:   s.operator,
			holder:     s.vendor,
			valid:      true,
			postTest: func(res *token.QueryIsOperatorForResponse) {
				s.Require().False(res.Authorized)
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
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.QueryIsOperatorForRequest{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				Holder:     tc.holder.String(),
			}
			res, err := s.queryServer.IsOperatorFor(s.goCtx, req)
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

func (s *KeeperTestSuite) TestQueryHoldersByOperator() {
	// empty request
	_, err := s.queryServer.HoldersByOperator(s.goCtx, nil)
	s.Require().Error(err)

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		valid      bool
		count      uint64
		postTest   func(res *token.QueryHoldersByOperatorResponse)
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			valid:      true,
			postTest: func(res *token.QueryHoldersByOperatorResponse) {
				s.Require().Equal(2, len(res.Holders))
			},
		},
		"valid request with limit": {
			contractID: s.contractID,
			operator:   s.operator,
			valid:      true,
			count:      1,
			postTest: func(res *token.QueryHoldersByOperatorResponse) {
				s.Require().Equal(1, len(res.Holders))
			},
		},
		"class not found": {
			contractID: "fee1dead",
			operator:   s.operator,
			valid:      true,
			postTest: func(res *token.QueryHoldersByOperatorResponse) {
				s.Require().Equal(0, len(res.Holders))
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
			req := &token.QueryHoldersByOperatorRequest{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				Pagination: pageReq,
			}
			res, err := s.queryServer.HoldersByOperator(s.goCtx, req)
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
