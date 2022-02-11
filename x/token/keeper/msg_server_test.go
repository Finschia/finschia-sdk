package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func (s *KeeperTestSuite) TestMsgTransfer() {
	testCases := map[string]struct{
		classId string
		amount sdk.Int
		valid bool
	}{
		"valid request": {
			classId: s.mintableClass,
			amount: s.balance,
			valid: true,
		},
		"insufficient funds (amount)": {
			classId: s.mintableClass,
			amount: s.balance.Add(sdk.OneInt()),
			valid: false,
		},
		"insufficient funds (no such a class)": {
			classId: "typo",
			amount: s.balance,
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgTransfer{
				ClassId: tc.classId,
				From: s.vendor.String(),
				To: s.customer.String(),
				Amount: tc.amount,
			}
			res, err := s.msgServer.Transfer(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgTransferFrom() {
	testCases := map[string]struct{
		classId string
		proxy sdk.AccAddress
		from sdk.AccAddress
		valid bool
	}{
		"valid request": {
			classId: s.mintableClass,
			proxy: s.operator,
			from: s.customer,
			valid: true,
		},
		"not approved": {
			classId: s.mintableClass,
			proxy: s.vendor,
			from: s.customer,
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgTransferFrom{
				ClassId: tc.classId,
				Proxy: tc.proxy.String(),
				From: tc.from.String(),
				To: s.vendor.String(),
				Amount: s.balance,
			}
			res, err := s.msgServer.TransferFrom(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NotNil(res)
		})
	}
}

func (s *KeeperTestSuite) TestMsgApprove() {
	testCases := map[string]struct{
		classId string
		approver sdk.AccAddress
		proxy sdk.AccAddress
		valid bool
	}{
		"valid request": {
			classId: s.mintableClass,
			approver: s.customer,
			proxy: s.vendor,
			valid: true,
		},
		"already approved": {
			classId: s.mintableClass,
			approver: s.customer,
			proxy: s.operator,
			valid: false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			req := &token.MsgApprove{
				ClassId: tc.classId,
				Approver: tc.approver.String(),
				Proxy: tc.proxy.String(),
			}
			res, err := s.msgServer.Approve(s.goCtx, req)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NotNil(res)
		})
	}
}
