package keeper_test

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

func newParams(enabled bool) *foundation.Params {
	params := foundation.DefaultParams()
	params.Enabled = enabled
	return params
}

func newUpdateFoundationParamsProposal(params *foundation.Params) govtypes.Content {
	return foundation.NewUpdateFoundationParamsProposal("Test", "description", params)
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
