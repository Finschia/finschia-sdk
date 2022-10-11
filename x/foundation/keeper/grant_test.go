package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestGrant() {
	testCases := map[string]struct {
		grantee sdk.AccAddress
		auth    foundation.Authorization
		valid   bool
	}{
		"valid authz": {
			grantee: s.members[0],
			auth:    &foundation.ReceiveFromTreasuryAuthorization{},
			valid:   true,
		},
		"override attempt": {
			grantee: s.stranger,
			auth:    &foundation.ReceiveFromTreasuryAuthorization{},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.Grant(ctx, tc.grantee, tc.auth)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestRevoke() {
	testCases := map[string]struct {
		grantee sdk.AccAddress
		url     string
		valid   bool
	}{
		"valid url": {
			grantee: s.stranger,
			url:     foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
			valid:   true,
		},
		"grant not found": {
			grantee: s.members[0],
			url:     foundation.ReceiveFromTreasuryAuthorization{}.MsgTypeURL(),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.Revoke(ctx, tc.grantee, tc.url)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestAccept() {
	testCases := map[string]struct {
		grantee sdk.AccAddress
		msg     sdk.Msg
		valid   bool
	}{
		"valid request": {
			grantee: s.stranger,
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To:       s.stranger.String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			},
			valid: true,
		},
		"no authorization": {
			grantee: s.members[0],
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To:       s.members[0].String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.Accept(ctx, tc.grantee, tc.msg)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
