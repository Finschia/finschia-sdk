package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestMsgSend() {
	testCases := map[string]struct {
		contractID string
		amount     sdk.Int
		valid      bool
	}{
		"valid request": {
			contractID: s.contractID,
			amount:     s.balance,
			valid:      true,
		},
		"insufficient funds (no such a class)": {
			contractID: "fee1dead",
			amount:     sdk.OneInt(),
		},
		"insufficient funds (not enough balance)": {
			contractID: s.contractID,
			amount:     s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgSend{
				ContractId: tc.contractID,
				From:       s.vendor.String(),
				To:         s.customer.String(),
				Amount:     tc.amount,
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
		operator sdk.AccAddress
		from     sdk.AccAddress
		amount   sdk.Int
		valid    bool
	}{
		"valid request": {
			operator: s.operator,
			from:     s.customer,
			amount:   s.balance,
			valid:    true,
		},
		"not approved": {
			operator: s.vendor,
			from:     s.customer,
			amount:   s.balance,
		},
		"insufficient funds (not enough balance)": {
			operator: s.operator,
			from:     s.customer,
			amount:   s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgOperatorSend{
				ContractId: s.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				To:         s.vendor.String(),
				Amount:     tc.amount,
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

func (s *KeeperTestSuite) TestMsgTransferFrom() {
	testCases := map[string]struct {
		proxy  sdk.AccAddress
		from   sdk.AccAddress
		amount sdk.Int
		valid  bool
	}{
		"valid request": {
			proxy:  s.operator,
			from:   s.customer,
			amount: s.balance,
			valid:  true,
		},
		"not approved": {
			proxy:  s.vendor,
			from:   s.customer,
			amount: s.balance,
		},
		"insufficient funds (not enough balance)": {
			proxy:  s.operator,
			from:   s.customer,
			amount: s.balance.Add(sdk.OneInt()),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgTransferFrom{
				ContractId: s.contractID,
				Proxy:      tc.proxy.String(),
				From:       tc.from.String(),
				To:         s.vendor.String(),
				Amount:     tc.amount,
			}
			res, err := s.msgServer.TransferFrom(sdk.WrapSDKContext(ctx), req)
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
		holder   sdk.AccAddress
		operator sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			holder:   s.customer,
			operator: s.vendor,
			valid:    true,
		},
		"already approved": {
			holder:   s.customer,
			operator: s.operator,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgAuthorizeOperator{
				ContractId: s.contractID,
				Holder:     tc.holder.String(),
				Operator:   tc.operator.String(),
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
		holder   sdk.AccAddress
		operator sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			holder:   s.customer,
			operator: s.operator,
			valid:    true,
		},
		"no authorization": {
			holder:   s.customer,
			operator: s.vendor,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgRevokeOperator{
				ContractId: s.contractID,
				Holder:     tc.holder.String(),
				Operator:   tc.operator.String(),
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

func (s *KeeperTestSuite) TestMsgApprove() {
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

			req := &token.MsgApprove{
				ContractId: s.contractID,
				Approver:   tc.approver.String(),
				Proxy:      tc.proxy.String(),
			}
			res, err := s.msgServer.Approve(sdk.WrapSDKContext(ctx), req)
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
		granter    sdk.AccAddress
		grantee    sdk.AccAddress
		permission token.Permission
		valid      bool
	}{
		"valid request": {
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.PermissionModify,
			valid:      true,
		},
		"already granted": {
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.PermissionMint,
		},
		"granter has no permission": {
			granter:    s.customer,
			grantee:    s.operator,
			permission: token.PermissionModify,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgGrant{
				ContractId: s.contractID,
				Granter:    tc.granter.String(),
				Grantee:    tc.grantee.String(),
				Permission: tc.permission,
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
		grantee    sdk.AccAddress
		permission token.Permission
		valid      bool
	}{
		"valid request": {
			grantee:    s.operator,
			permission: token.PermissionMint,
			valid:      true,
		},
		"not granted yet": {
			grantee:    s.operator,
			permission: token.PermissionModify,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgAbandon{
				ContractId: s.contractID,
				Grantee:    tc.grantee.String(),
				Permission: tc.permission,
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

func (s *KeeperTestSuite) TestMsgGrantPermission() {
	testCases := map[string]struct {
		granter    sdk.AccAddress
		grantee    sdk.AccAddress
		permission string
		valid      bool
	}{
		"valid request": {
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.LegacyPermissionModify.String(),
			valid:      true,
		},
		"already granted": {
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.LegacyPermissionMint.String(),
		},
		"granter has no permission": {
			granter:    s.customer,
			grantee:    s.operator,
			permission: token.LegacyPermissionModify.String(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgGrantPermission{
				ContractId: s.contractID,
				From:       tc.granter.String(),
				To:         tc.grantee.String(),
				Permission: tc.permission,
			}
			res, err := s.msgServer.GrantPermission(sdk.WrapSDKContext(ctx), req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokePermission() {
	testCases := map[string]struct {
		from       sdk.AccAddress
		permission string
		valid      bool
	}{
		"valid request": {
			from:       s.operator,
			permission: token.LegacyPermissionMint.String(),
			valid:      true,
		},
		"not granted yet": {
			from:       s.operator,
			permission: token.LegacyPermissionModify.String(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgRevokePermission{
				ContractId: s.contractID,
				From:       tc.from.String(),
				Permission: tc.permission,
			}
			res, err := s.msgServer.RevokePermission(sdk.WrapSDKContext(ctx), req)
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
				ContractId: s.contractID,
				From:       tc.grantee.String(),
				To:         s.customer.String(),
				Amount:     sdk.OneInt(),
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
			from: s.customer,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgBurn{
				ContractId: s.contractID,
				From:       tc.from.String(),
				Amount:     s.balance,
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
		operator sdk.AccAddress
		from     sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			operator: s.operator,
			from:     s.customer,
			valid:    true,
		},
		"not approved": {
			operator: s.vendor,
			from:     s.customer,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgOperatorBurn{
				ContractId: s.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				Amount:     s.balance,
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

func (s *KeeperTestSuite) TestMsgBurnFrom() {
	testCases := map[string]struct {
		proxy sdk.AccAddress
		from  sdk.AccAddress
		valid bool
	}{
		"valid request": {
			proxy: s.operator,
			from:  s.customer,
			valid: true,
		},
		"not approved": {
			proxy: s.vendor,
			from:  s.customer,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgBurnFrom{
				ContractId: s.contractID,
				Proxy:      tc.proxy.String(),
				From:       tc.from.String(),
				Amount:     s.balance,
			}
			res, err := s.msgServer.BurnFrom(sdk.WrapSDKContext(ctx), req)
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
				ContractId: s.contractID,
				Owner:      tc.grantee.String(),
				Changes:    []token.Pair{{Field: token.AttributeKeyName.String(), Value: "hello"}},
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
