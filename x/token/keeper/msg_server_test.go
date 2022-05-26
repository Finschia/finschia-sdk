package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestMsgSend() {
	testCases := map[string]struct {
		classId string
		amount  sdk.Int
		valid   bool
	}{
		"valid request": {
			classId: s.classID,
			amount:  s.balance,
			valid:   true,
		},
		"insufficient funds (no such a class)": {
			classId: "fee1dead",
			amount:  sdk.OneInt(),
		},
		"insufficient funds (not enough balance)": {
			classId: s.classID,
			amount:  s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgSend{
				ContractId: tc.classId,
				From:    s.vendor.String(),
				To:      s.customer.String(),
				Amount:  tc.amount,
			}
			res, err := s.msgServer.Send(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorSend() {
	testCases := map[string]struct {
		proxy sdk.AccAddress
		from  sdk.AccAddress
		amount sdk.Int
		valid bool
	}{
		"valid request": {
			proxy: s.operator,
			from:  s.customer,
			amount: s.balance,
			valid: true,
		},
		"not approved": {
			proxy: s.vendor,
			from:  s.customer,
			amount: s.balance,
		},
		"insufficient funds (not enough balance)": {
			proxy: s.operator,
			from:  s.customer,
			amount:  s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgOperatorSend{
				ContractId: s.classID,
				Proxy:   tc.proxy.String(),
				From:    tc.from.String(),
				To:      s.vendor.String(),
				Amount:  tc.amount,
			}
			res, err := s.msgServer.OperatorSend(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgAuthorizeOperator() {
	testCases := map[string]struct {
		approver sdk.AccAddress
		proxy    sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			approver: s.customer,
			proxy:    s.vendor,
			valid:    true,
		},
		"already approved": {
			approver: s.customer,
			proxy:    s.operator,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgAuthorizeOperator{
				ContractId:  s.classID,
				Approver: tc.approver.String(),
				Proxy:    tc.proxy.String(),
			}
			res, err := s.msgServer.AuthorizeOperator(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokeOperator() {
	testCases := map[string]struct {
		approver sdk.AccAddress
		proxy    sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			approver: s.customer,
			proxy:    s.operator,
			valid:    true,
		},
		"no authorization": {
			approver: s.customer,
			proxy:    s.vendor,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgRevokeOperator{
				ContractId:  s.classID,
				Approver: tc.approver.String(),
				Proxy:    tc.proxy.String(),
			}
			res, err := s.msgServer.RevokeOperator(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgIssue() {
	testCases := map[string]struct {
		amount sdk.Int
		valid  bool
	}{
		"valid request": {
			amount: sdk.OneInt(),
			valid:  true,
		},
		"invalid amount (not possible to reach here but for cov)": {
			amount: sdk.NewInt(-1),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgIssue{
				Owner:  s.vendor.String(),
				To:     s.vendor.String(),
				Name:   "test",
				Symbol: "TT",
				Amount: tc.amount,
			}
			res, err := s.msgServer.Issue(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgGrant() {
	testCases := map[string]struct {
		granter sdk.AccAddress
		grantee sdk.AccAddress
		permission  string
		valid   bool
	}{
		"valid request": {
			granter: s.vendor,
			grantee: s.operator,
			permission:  token.Permission_Modify.String(),
			valid:   true,
		},
		"already granted": {
			granter: s.vendor,
			grantee: s.operator,
			permission:  token.Permission_Mint.String(),
		},
		"granter has no permission": {
			granter: s.customer,
			grantee: s.operator,
			permission:  token.Permission_Modify.String(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgGrant{
				ContractId: s.classID,
				From: tc.granter.String(),
				To: tc.grantee.String(),
				Permission:  tc.permission,
			}
			res, err := s.msgServer.Grant(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgAbandon() {
	testCases := map[string]struct {
		grantee sdk.AccAddress
		permission  string
		valid   bool
	}{
		"valid request": {
			grantee: s.operator,
			permission:  token.Permission_Mint.String(),
			valid:   true,
		},
		"not granted yet": {
			grantee: s.operator,
			permission:  token.Permission_Modify.String(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgAbandon{
				ContractId: s.classID,
				Grantee: tc.grantee.String(),
				Permission:  tc.permission,
			}
			res, err := s.msgServer.Abandon(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgMint() {
	testCases := map[string]struct {
		grantee sdk.AccAddress
		valid   bool
	}{
		"valid request": {
			grantee: s.operator,
			valid:   true,
		},
		"not granted": {
			grantee: s.customer,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgMint{
				ContractId: s.classID,
				From: tc.grantee.String(),
				To:      s.customer.String(),
				Amount:  sdk.OneInt(),
			}
			res, err := s.msgServer.Mint(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurn() {
	testCases := map[string]struct {
		from  sdk.AccAddress
		valid bool
	}{
		"valid request": {
			from:  s.vendor,
			valid: true,
		},
		"not granted": {
			from:  s.customer,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgBurn{
				ContractId: s.classID,
				From:    tc.from.String(),
				Amount:  s.balance,
			}
			res, err := s.msgServer.Burn(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorBurn() {
	testCases := map[string]struct {
		grantee sdk.AccAddress
		from    sdk.AccAddress
		valid   bool
	}{
		"valid request": {
			grantee: s.operator,
			from:    s.customer,
			valid:   true,
		},
		"not approved": {
			grantee: s.vendor,
			from:    s.customer,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgOperatorBurn{
				ContractId: s.classID,
				Proxy: tc.grantee.String(),
				From:    tc.from.String(),
				Amount:  s.balance,
			}
			res, err := s.msgServer.OperatorBurn(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgModify() {
	testCases := map[string]struct {
		grantee sdk.AccAddress
		valid   bool
	}{
		"valid request": {
			grantee: s.vendor,
			valid:   true,
		},
		"not granted": {
			grantee: s.operator,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgModify{
				ContractId: s.classID,
				Owner: tc.grantee.String(),
				Changes: []token.Pair{{Field: token.AttributeKey_Name.String(), Value: "hello"}},
			}
			res, err := s.msgServer.Modify(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}
