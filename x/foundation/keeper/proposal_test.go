package keeper_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/keeper"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	"github.com/line/lbm-sdk/x/stakingplus"
)

func newParams(enabled bool) *foundation.Params {
	params := foundation.DefaultParams()
	params.Enabled = enabled
	return params
}

func newUpdateFoundationParamsProposal(params *foundation.Params) govtypes.Content {
	return foundation.NewUpdateFoundationParamsProposal("Test", "description", params)
}

func newValidatorAuths(addrs []sdk.ValAddress, allow bool) []foundation.ValidatorAuth {
	auths := []foundation.ValidatorAuth{}
	for _, addr := range addrs {
		auth := foundation.ValidatorAuth{
			OperatorAddress: addr.String(),
			CreationAllowed: allow,
		}
		auths = append(auths, auth)
	}

	return auths
}

func newUpdateValidatorAuthsProposal(auths []foundation.ValidatorAuth) govtypes.Content {
	return foundation.NewUpdateValidatorAuthsProposal("Test", "description", auths)
}

func TestProposalHandler(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	// turn on the module
	k := app.FoundationKeeper
	params_on := newParams(true)
	k.SetParams(ctx, params_on)

	handler := keeper.NewProposalHandler(k)

	msgTypeURL := stakingplus.CreateValidatorAuthorization{}.MsgTypeURL()
	// test adding creation allowed validators
	adding := newValidatorAuths([]sdk.ValAddress{valAddr}, true)
	ap := newUpdateValidatorAuthsProposal(adding)
	require.NoError(t, ap.ValidateBasic())
	require.NoError(t, handler(ctx, ap))
	for i := range adding {
		valAddr, err := sdk.ValAddressFromBech32(adding[i].OperatorAddress)
		grantee := sdk.AccAddress(valAddr)
		require.NoError(t, err)
		_, err = k.GetAuthorization(ctx, govtypes.ModuleName, grantee, msgTypeURL)
		require.NoError(t, err)
	}

	// test deleting creation allowed validators
	deleting := newValidatorAuths([]sdk.ValAddress{valAddr}, false)
	dp := newUpdateValidatorAuthsProposal(deleting)
	require.NoError(t, dp.ValidateBasic())
	require.NoError(t, handler(ctx, dp))
	for i := range deleting {
		valAddr, err := sdk.ValAddressFromBech32(adding[i].OperatorAddress)
		grantee := sdk.AccAddress(valAddr)
		_, err = k.GetAuthorization(ctx, govtypes.ModuleName, grantee, msgTypeURL)
		require.Error(t, err)
	}

	// disable foundation
	params_off := newParams(false)
	pp := newUpdateFoundationParamsProposal(params_off)
	require.NoError(t, pp.ValidateBasic())
	require.NoError(t, handler(ctx, pp))
	require.Empty(t, k.GetGrants(ctx))
	require.Equal(t, params_off, k.GetParams(ctx))

	// attempt to enable foundation, which fails
	pp = newUpdateFoundationParamsProposal(params_on)
	require.Error(t, pp.ValidateBasic())
}

func (s *KeeperTestSuite) TestSubmitProposal() {
	testCases := map[string]struct {
		proposers []string
		metadata  string
		msg       sdk.Msg
		valid     bool
	}{
		"valid proposal": {
			proposers: []string{s.members[0].String()},
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To:       s.stranger.String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			},
			valid: true,
		},
		"long metadata": {
			proposers: []string{s.members[0].String()},
			metadata:  string(make([]rune, 256)),
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.operator.String(),
				To:       s.stranger.String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			},
		},
		"unauthorized msg": {
			proposers: []string{s.members[0].String()},
			msg: &foundation.MsgWithdrawFromTreasury{
				Operator: s.stranger.String(),
				To:       s.stranger.String(),
				Amount:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())),
			},
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			_, err := s.keeper.SubmitProposal(ctx, tc.proposers, tc.metadata, []sdk.Msg{tc.msg})
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestWithdrawProposal() {
	testCases := map[string]struct {
		id    uint64
		valid bool
	}{
		"valid proposal": {
			id:    s.activeProposal,
			valid: true,
		},
		"not active": {
			id: s.abortedProposal,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.keeper.WithdrawProposal(ctx, tc.id)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}
