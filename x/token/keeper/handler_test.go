package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	banktypes "github.com/line/lbm-sdk/x/bank/types"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/keeper"
)

func (s *KeeperTestSuite) TestNewHandler() {
	testCases := map[string]struct {
		msg   sdk.Msg
		valid bool
	}{
		"MsgTransfer": {
			&token.MsgTransfer{
				ClassId: s.classID,
				From:    s.vendor.String(),
				To:      s.customer.String(),
				Amount:  sdk.OneInt(),
			},
			true,
		},
		"MsgTransferFrom": {
			&token.MsgTransferFrom{
				ClassId: s.classID,
				Proxy:   s.operator.String(),
				From:    s.customer.String(),
				To:      s.vendor.String(),
				Amount:  sdk.OneInt(),
			},
			true,
		},
		"MsgApprove": {
			&token.MsgApprove{
				ClassId:  s.classID,
				Approver: s.customer.String(),
				Proxy:    s.vendor.String(),
			},
			true,
		},
		"MsgIssue": {
			&token.MsgIssue{
				Owner:  s.vendor.String(),
				To:     s.vendor.String(),
				Name:   "for test",
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
				Action:  token.ActionModify,
			},
			true,
		},
		"MsgRevoke": {
			&token.MsgRevoke{
				ClassId: s.classID,
				Grantee: s.operator.String(),
				Action:  token.ActionMint,
			},
			true,
		},
		"MsgMint": {
			&token.MsgMint{
				ClassId: s.classID,
				Grantee: s.vendor.String(),
				To:      s.customer.String(),
				Amount:  sdk.OneInt(),
			},
			true,
		},
		"MsgBurn": {
			&token.MsgBurn{
				ClassId: s.classID,
				From:    s.operator.String(),
				Amount:  sdk.OneInt(),
			},
			true,
		},
		"MsgBurnFrom": {
			&token.MsgBurnFrom{
				ClassId: s.classID,
				Grantee: s.operator.String(),
				From:    s.customer.String(),
				Amount:  sdk.OneInt(),
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
