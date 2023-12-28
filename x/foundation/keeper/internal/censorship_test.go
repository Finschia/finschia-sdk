package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestUpdateCensorship() {
	ctx, _ := s.ctx.CacheContext()

	// add a dummy url to censorship
	msgTypeURL := sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil))
	dummyURL := sdk.MsgTypeURL((*foundation.MsgFundTreasury)(nil))
	for _, url := range []string{
		msgTypeURL,
		dummyURL,
	} {
		s.impl.SetCensorship(ctx, foundation.Censorship{
			MsgTypeUrl: url,
			Authority:  foundation.CensorshipAuthorityFoundation,
		})
	}

	// check preconditions
	s.Require().True(s.impl.IsCensoredMessage(ctx, msgTypeURL))
	_, err := s.impl.GetAuthorization(ctx, s.stranger, msgTypeURL)
	s.Require().NoError(err)

	// test update censorship
	removingCensorship := foundation.Censorship{
		MsgTypeUrl: msgTypeURL,
		Authority:  foundation.CensorshipAuthorityUnspecified,
	}
	s.Require().NoError(removingCensorship.ValidateBasic())
	err = s.impl.UpdateCensorship(ctx, removingCensorship)
	s.Require().NoError(err)

	// check censorship
	_, err = s.impl.GetCensorship(ctx, msgTypeURL)
	s.Require().Error(err)
	s.Require().False(s.impl.IsCensoredMessage(ctx, msgTypeURL))

	// check authorizations
	_, err = s.impl.GetAuthorization(ctx, s.stranger, msgTypeURL)
	s.Require().Error(err)

	// 2. re-enable the removed censorship, which must fail
	newCensorship := foundation.Censorship{
		MsgTypeUrl: msgTypeURL,
		Authority:  foundation.CensorshipAuthorityGovernance,
	}
	s.Require().NoError(newCensorship.ValidateBasic())
	err = s.impl.UpdateCensorship(ctx, newCensorship)
	s.Require().Error(err)
}

func (s *KeeperTestSuite) TestGrant() {
	testCases := map[string]struct {
		malleate func(ctx sdk.Context)
		grantee  sdk.AccAddress
		auth     foundation.Authorization
		valid    bool
	}{
		"valid authz": {
			grantee: s.members[0],
			auth:    &foundation.ReceiveFromTreasuryAuthorization{},
			valid:   true,
		},
		"not being censored": {
			malleate: func(ctx sdk.Context) {
				err := s.impl.UpdateCensorship(ctx, foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityUnspecified,
				})
				s.Require().NoError(err)
			},
			grantee: s.members[0],
			auth:    &foundation.ReceiveFromTreasuryAuthorization{},
		},
		"override attempt": {
			grantee: s.stranger,
			auth:    &foundation.ReceiveFromTreasuryAuthorization{},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			err := s.impl.Grant(ctx, tc.grantee, tc.auth)
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

			err := s.impl.Revoke(ctx, tc.grantee, tc.url)
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
		malleate func(ctx sdk.Context)
		grantee  sdk.AccAddress
		msg      sdk.Msg
		valid    bool
	}{
		"valid request": {
			grantee: s.stranger,
			msg: &foundation.MsgWithdrawFromTreasury{
				Authority: s.bytesToString(s.authority),
				To:        s.bytesToString(s.stranger),
				Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt())),
			},
			valid: true,
		},
		"not being censored": {
			malleate: func(ctx sdk.Context) {
				err := s.impl.UpdateCensorship(ctx, foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityUnspecified,
				})
				s.Require().NoError(err)
			},
			grantee: s.members[0],
			msg: &foundation.MsgWithdrawFromTreasury{
				Authority: s.bytesToString(s.authority),
				To:        s.bytesToString(s.members[0]),
				Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt())),
			},
			valid: true,
		},
		"no authorization": {
			grantee: s.members[0],
			msg: &foundation.MsgWithdrawFromTreasury{
				Authority: s.bytesToString(s.authority),
				To:        s.bytesToString(s.members[0]),
				Amount:    sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.OneInt())),
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			err := s.impl.Accept(ctx, tc.grantee, tc.msg)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
