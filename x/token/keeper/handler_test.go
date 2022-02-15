package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/keeper"
	banktypes "github.com/line/lbm-sdk/x/bank/types"
)

func (s *KeeperTestSuite) TestNewHandler() {
	amount := sdk.OneInt()
	testCases := map[string]struct{
		msg sdk.Msg
		valid bool
	}{
		"MsgTransfer": {
			&token.MsgTransfer{
				ClassId: s.classID,
				From: s.vendor.String(),
				To: s.customer.String(),
				Amount: amount,
			},
			true,
		},
		"MsgTransferFrom": {
			&token.MsgTransferFrom{
				ClassId: s.classID,
				Proxy: s.operator.String(),
				From: s.customer.String(),
				To: s.vendor.String(),
				Amount: amount,
			},
			true,
		},
		"MsgApprove": {
			&token.MsgApprove{
				ClassId: s.classID,
				Approver: s.customer.String(),
				Proxy: s.vendor.String(),
			},
			true,
		},
		"MsgIssue": {
			&token.MsgIssue{
				Owner: s.vendor.String(),
				To: s.vendor.String(),
				Name: "for test",
				Symbol: "TT",
				Amount: sdk.OneInt(),
			},
			true,
		},
		"MsgGrant": {
			&token.MsgGrant{
				ClassId: s.classID,
				Granter: s.vendor.String(),
				Grantee: s.operator.String(),
				Action: token.ActionModify,
			},
			true,
		},
		"MsgRevoke": {
			&token.MsgRevoke{
				ClassId: s.classID,
				Grantee: s.operator.String(),
				Action: token.ActionMint,
			},
			true,
		},
		"MsgMint": {
			&token.MsgMint{
				ClassId: s.classID,
				Grantee: s.vendor.String(),
				To: s.customer.String(),
				Amount: amount,
			},
			true,
		},
		"MsgBurn": {
			&token.MsgBurn{
				ClassId: s.classID,
				From: s.operator.String(),
				Amount: amount,
			},
			true,
		},
		"MsgBurnFrom": {
			&token.MsgBurnFrom{
				ClassId: s.classID,
				Grantee: s.operator.String(),
				From: s.customer.String(),
				Amount: amount,
			},
			true,
		},
		"MsgModify": {
			&token.MsgModify{
				ClassId: s.classID,
				Grantee: s.vendor.String(),
				Changes: []token.Pair{{Key: "name", Value: "new name"}},
			},
			true,
		},
		"MsgSend": {
			&banktypes.MsgSend{},
			false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			handler := keeper.NewHandler(s.keeper)
			_, err := handler(s.ctx, tc.msg)
			if !tc.valid {
				s.Require().Error(err)
				return
			}
			s.Require().NoError(err)
		})
	}
}
