package keeper_test

import (
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/testutil/testdata"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestSubmitProposal() {
	testCases := map[string]struct {
		proposers []string
		metadata  string
		msg       sdk.Msg
		valid     bool
	}{
		"valid proposal": {
			proposers: []string{s.members[0].String()},
			msg:       testdata.NewTestMsg(s.operator),
			valid:     true,
		},
		"long metadata": {
			proposers: []string{s.members[0].String()},
			metadata:  string(make([]rune, 256)),
			msg:       testdata.NewTestMsg(s.operator),
		},
		"unauthorized msg": {
			proposers: []string{s.members[0].String()},
			msg:       testdata.NewTestMsg(s.stranger),
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
			id: s.withdrawnProposal,
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

func TestAbortProposal(t *testing.T) {
	checkTx := false
	app := simapp.Setup(checkTx)
	testdata.RegisterInterfaces(app.InterfaceRegistry())

	ctx := app.BaseApp.NewContext(checkTx, ocproto.Header{})
	keeper := app.FoundationKeeper

	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	operator := keeper.GetOperator(ctx)

	members := make([]sdk.AccAddress, 10)
	for i := range members {
		members[i] = createAddress()
	}
	err := keeper.UpdateMembers(ctx, []foundation.MemberRequest{
		{
			Address: members[0].String(),
		},
	})
	require.NoError(t, err)

	// create proposals of different versions and abort them
	for _, newMember := range members[1:] {
		_, err := keeper.SubmitProposal(ctx, []string{members[0].String()}, "", []sdk.Msg{testdata.NewTestMsg(operator)})
		require.NoError(t, err)

		err = keeper.UpdateMembers(ctx, []foundation.MemberRequest{
			{
				Address: newMember.String(),
			},
		})
		require.NoError(t, err)
	}

	for _, proposal := range keeper.GetProposals(ctx) {
		require.Equal(t, foundation.PROPOSAL_STATUS_ABORTED, proposal.Status)
	}
}
