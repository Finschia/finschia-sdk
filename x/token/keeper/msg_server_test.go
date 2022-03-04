package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestMsgTransfer() {
	testCases := map[string]struct {
		classId string
		amount  sdk.Int
		valid   bool
	}{
		"valid request": {
			classId: s.classID,
			amount:  sdk.OneInt(),
			valid:   true,
		},
		"insufficient funds (no such a class)": {
			classId: "00000000",
			amount:  sdk.OneInt(),
			valid:   false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgTransfer{
				ClassId: tc.classId,
				From:    s.vendor.String(),
				To:      s.customer.String(),
				Amount:  tc.amount,
			}
			res, err := s.msgServer.Transfer(s.goCtx, req)
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
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgTransferFrom{
				ClassId: s.classID,
				Proxy:   tc.proxy.String(),
				From:    tc.from.String(),
				To:      s.vendor.String(),
				Amount:  sdk.OneInt(),
			}
			res, err := s.msgServer.TransferFrom(s.goCtx, req)
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
			valid:    false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgApprove{
				ClassId:  s.classID,
				Approver: tc.approver.String(),
				Proxy:    tc.proxy.String(),
			}
			res, err := s.msgServer.Approve(s.goCtx, req)
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
			valid:  false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgIssue{
				Owner:  s.vendor.String(),
				To:     s.vendor.String(),
				Name:   "test",
				Symbol: "TT",
				Amount: tc.amount,
			}
			res, err := s.msgServer.Issue(s.goCtx, req)
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
		action  string
		valid   bool
	}{
		"valid request": {
			granter: s.vendor,
			grantee: s.operator,
			action:  token.ActionModify,
			valid:   true,
		},
		"already granted": {
			granter: s.vendor,
			grantee: s.operator,
			action:  token.ActionMint,
			valid:   false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgGrant{
				ClassId: s.classID,
				Granter: tc.granter.String(),
				Grantee: tc.grantee.String(),
				Action:  tc.action,
			}
			res, err := s.msgServer.Grant(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgRevoke() {
	testCases := map[string]struct {
		grantee sdk.AccAddress
		action  string
		valid   bool
	}{
		"valid request": {
			grantee: s.operator,
			action:  token.ActionMint,
			valid:   true,
		},
		"not granted yet": {
			grantee: s.operator,
			action:  token.ActionModify,
			valid:   false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgRevoke{
				ClassId: s.classID,
				Grantee: tc.grantee.String(),
				Action:  tc.action,
			}
			res, err := s.msgServer.Revoke(s.goCtx, req)
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
			valid:   false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgMint{
				ClassId: s.classID,
				Grantee: tc.grantee.String(),
				To:      s.customer.String(),
				Amount:  sdk.OneInt(),
			}
			res, err := s.msgServer.Mint(s.goCtx, req)
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
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgBurn{
				ClassId: s.classID,
				From:    tc.from.String(),
				Amount:  sdk.OneInt(),
			}
			res, err := s.msgServer.Burn(s.goCtx, req)
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
			valid:   false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgBurnFrom{
				ClassId: s.classID,
				Grantee: tc.grantee.String(),
				From:    tc.from.String(),
				Amount:  sdk.OneInt(),
			}
			res, err := s.msgServer.BurnFrom(s.goCtx, req)
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
			valid:   false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgModify{
				ClassId: s.classID,
				Grantee: tc.grantee.String(),
				Changes: []token.Pair{{Key: token.AttributeKeyName, Value: "hello"}},
			}
			res, err := s.msgServer.Modify(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
			s.Require().NotNil(res)
		})
	}
}
