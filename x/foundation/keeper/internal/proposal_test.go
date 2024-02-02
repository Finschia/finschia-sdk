package internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (s *KeeperTestSuite) TestSubmitProposal() {
	testCases := map[string]struct {
		proposers []string
		metadata  string
		msg       sdk.Msg
		valid     bool
	}{
		"valid proposal": {
			proposers: []string{s.bytesToString(s.members[0])},
			msg:       s.newTestMsg(s.authority),
			valid:     true,
		},
		"long metadata": {
			proposers: []string{s.bytesToString(s.members[0])},
			metadata:  string(make([]rune, 256)),
			msg:       s.newTestMsg(s.authority),
		},
		"unauthorized msg": {
			proposers: []string{s.bytesToString(s.members[0])},
			msg:       s.newTestMsg(s.stranger),
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			_, err := s.impl.SubmitProposal(ctx, tc.proposers, tc.metadata, []sdk.Msg{tc.msg})
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
		"no such a proposal": {
			id: s.nextProposal,
		},
		"not active": {
			id: s.withdrawnProposal,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			ctx, _ := s.ctx.CacheContext()

			err := s.impl.WithdrawProposal(ctx, tc.id)
			if tc.valid {
				s.Require().NoError(err)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func TestAbortProposal(t *testing.T) {
	impl, _, _, _, _, addressCodec, ctx := setupFoundationKeeper(t, nil, nil)

	bytesToString := func(addr sdk.AccAddress) string {
		str, err := addressCodec.BytesToString(addr)
		require.NoError(t, err)
		return str
	}

	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	gs := &foundation.GenesisState{
		Params: foundation.DefaultParams(),
	}

	members := make([]sdk.AccAddress, 10)
	for i := range members {
		members[i] = createAddress()
	}
	gs.Members = []foundation.Member{{
		Address: bytesToString(members[0]),
	}}

	info := foundation.DefaultFoundation()
	info.TotalWeight = math.LegacyNewDec(1)
	err := info.SetDecisionPolicy(workingPolicy())
	require.NoError(t, err)
	gs.Foundation = info

	impl.InitGenesis(ctx, gs)

	// create proposals of different versions and abort them
	for _, newMember := range members[1:] {
		_, err := impl.SubmitProposal(ctx, []string{bytesToString(members[0])}, "", []sdk.Msg{&testdata.TestMsg{
			Signers: []string{impl.GetAuthority()},
		}})
		require.NoError(t, err)

		err = impl.UpdateMembers(ctx, []foundation.MemberRequest{
			{
				Address: bytesToString(newMember),
			},
		})
		require.NoError(t, err)
	}

	for _, proposal := range impl.GetProposals(ctx) {
		require.Equal(t, foundation.PROPOSAL_STATUS_ABORTED, proposal.Status)
	}
}
