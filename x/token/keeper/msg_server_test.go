package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/class"
)

func (s *KeeperTestSuite) TestMsgSend() {
	testCases := map[string]struct {
		contractID string
		amount     sdk.Int
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			amount:     s.balance,
		},
		"contract not found": {
			contractID: "fee1dead",
			amount:     sdk.OneInt(),
			err:        class.ErrContractNotExist,
		},
		"insufficient funds": {
			contractID: s.contractID,
			amount:     s.balance.Add(sdk.OneInt()),
			err:        token.ErrInsufficientBalance,
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorSend() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		amount     sdk.Int
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			amount:     s.balance,
		},
		"contract not found": {
			contractID: "fee1dead",
			operator:   s.operator,
			from:       s.customer,
			amount:     s.balance,
			err:        class.ErrContractNotExist,
		},
		"not approved": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			amount:     s.balance,
			err:        token.ErrTokenNotApproved,
		},
		"insufficient funds": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
			amount:     s.balance.Add(sdk.OneInt()),
			err:        token.ErrInsufficientBalance,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgOperatorSend{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				To:         s.vendor.String(),
				Amount:     tc.amount,
			}
			res, err := s.msgServer.OperatorSend(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokeOperator() {
	testCases := map[string]struct {
		contractID string
		holder     sdk.AccAddress
		operator   sdk.AccAddress
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			holder:     s.customer,
			operator:   s.operator,
		},
		"contract not found": {
			contractID: "fee1dead",
			holder:     s.customer,
			operator:   s.operator,
			err:        class.ErrContractNotExist,
		},
		"no authorization": {
			contractID: s.contractID,
			holder:     s.customer,
			operator:   s.vendor,
			err:        token.ErrTokenNotApproved,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgRevokeOperator{
				ContractId: tc.contractID,
				Holder:     tc.holder.String(),
				Operator:   tc.operator.String(),
			}
			res, err := s.msgServer.RevokeOperator(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgAuthorizeOperator() {
	testCases := map[string]struct {
		contractID string
		holder     sdk.AccAddress
		operator   sdk.AccAddress
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			holder:     s.customer,
			operator:   s.vendor,
		},
		"contract not found": {
			contractID: "fee1dead",
			holder:     s.customer,
			operator:   s.vendor,
			err:        class.ErrContractNotExist,
		},
		"already approved": {
			contractID: s.contractID,
			holder:     s.customer,
			operator:   s.operator,
			err:        token.ErrTokenAlreadyApproved,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgAuthorizeOperator{
				ContractId: tc.contractID,
				Holder:     tc.holder.String(),
				Operator:   tc.operator.String(),
			}
			res, err := s.msgServer.AuthorizeOperator(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgIssue() {
	testCases := map[string]struct {
		amount sdk.Int
		err    error
	}{
		"valid request": {
			amount: sdk.OneInt(),
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
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgGrantPermission() {
	testCases := map[string]struct {
		contractID string
		granter    sdk.AccAddress
		grantee    sdk.AccAddress
		permission string
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.LegacyPermissionModify.String(),
		},
		"contract not found": {
			contractID: "fee1dead",
			granter:    s.vendor,
			grantee:    s.operator,
			permission: token.LegacyPermissionModify.String(),
			err:        class.ErrContractNotExist,
		},
		"granter has no permission": {
			contractID: s.contractID,
			granter:    s.customer,
			grantee:    s.operator,
			permission: token.LegacyPermissionModify.String(),
			err:        token.ErrTokenNoPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgGrantPermission{
				ContractId: tc.contractID,
				From:       tc.granter.String(),
				To:         tc.grantee.String(),
				Permission: tc.permission,
			}
			res, err := s.msgServer.GrantPermission(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevokePermission() {
	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		permission string
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.operator,
			permission: token.LegacyPermissionMint.String(),
		},
		"contract not found": {
			contractID: "fee1dead",
			from:       s.operator,
			permission: token.LegacyPermissionMint.String(),
			err:        class.ErrContractNotExist,
		},
		"not granted yet": {
			contractID: s.contractID,
			from:       s.operator,
			permission: token.LegacyPermissionModify.String(),
			err:        token.ErrTokenNoPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgRevokePermission{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				Permission: tc.permission,
			}
			res, err := s.msgServer.RevokePermission(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgMint() {
	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			grantee:    s.operator,
		},
		"contract not found": {
			contractID: "fee1dead",
			grantee:    s.operator,
			err:        class.ErrContractNotExist,
		},
		"not granted": {
			contractID: s.contractID,
			grantee:    s.customer,
			err:        token.ErrTokenNoPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgMint{
				ContractId: tc.contractID,
				From:       tc.grantee.String(),
				To:         s.customer.String(),
				Amount:     sdk.OneInt(),
			}
			res, err := s.msgServer.Mint(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgBurn() {
	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			from:       s.vendor,
		},
		"contract not found": {
			contractID: "fee1dead",
			from:       s.vendor,
			err:        class.ErrContractNotExist,
		},
		"not granted": {
			contractID: s.contractID,
			from:       s.customer,
			err:        token.ErrTokenNoPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgBurn{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				Amount:     s.balance,
			}
			res, err := s.msgServer.Burn(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgOperatorBurn() {
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			operator:   s.operator,
			from:       s.customer,
		},
		"contract not found": {
			contractID: "fee1dead",
			operator:   s.operator,
			from:       s.customer,
			err:        class.ErrContractNotExist,
		},
		"not approved": {
			contractID: s.contractID,
			operator:   s.vendor,
			from:       s.customer,
			err:        token.ErrTokenNotApproved,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgOperatorBurn{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				Amount:     s.balance,
			}
			res, err := s.msgServer.OperatorBurn(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgModify() {
	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		err        error
	}{
		"valid request": {
			contractID: s.contractID,
			grantee:    s.vendor,
		},
		"contract not found": {
			contractID: "fee1dead",
			grantee:    s.vendor,
			err:        class.ErrContractNotExist,
		},
		"not granted": {
			contractID: s.contractID,
			grantee:    s.operator,
			err:        token.ErrTokenNoPermission,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			req := &token.MsgModify{
				ContractId: tc.contractID,
				Owner:      tc.grantee.String(),
				Changes:    []token.Attribute{{Key: token.AttributeKeyImageURI.String(), Value: "uri"}},
			}
			res, err := s.msgServer.Modify(sdk.WrapSDKContext(ctx), req)
			s.Require().ErrorIs(err, tc.err)
			if tc.err != nil {
				return
			}

			s.Require().NotNil(res)
		})
	}
}
