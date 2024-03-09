package internal_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestHooks() {
	testCases := map[string]struct {
		malleate func(ctx sdk.Context)
		grantee  sdk.AccAddress
		valid    bool
	}{
		"valid request": {
			grantee: s.stranger,
			valid:   true,
		},
		"no authorization": {
			grantee: s.members[0],
			valid:   false,
		},
		"not being censored": {
			malleate: func(ctx sdk.Context) {
				err := s.impl.UpdateCensorship(ctx, foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*stakingtypes.MsgCreateValidator)(nil)),
					Authority:  foundation.CensorshipAuthorityUnspecified,
				})
				s.Require().NoError(err)
			},
			grantee: s.members[0],
			valid:   true,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			err := s.keeper.Hooks().AfterValidatorCreated(ctx, sdk.ValAddress(tc.grantee))

			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
