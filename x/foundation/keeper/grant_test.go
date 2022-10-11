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
