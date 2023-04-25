package internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/simapp"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/keeper/internal"
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
			msg:       testdata.NewTestMsg(s.authority),
			valid:     true,
		},
		"long metadata": {
			proposers: []string{s.members[0].String()},
			metadata:  string(make([]rune, 256)),
			msg:       testdata.NewTestMsg(s.authority),
		},
		"unauthorized msg": {
			proposers: []string{s.members[0].String()},
			msg:       testdata.NewTestMsg(s.stranger),
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
	checkTx := false
	app := simapp.Setup(checkTx)
	testdata.RegisterInterfaces(app.InterfaceRegistry())

	ctx := app.BaseApp.NewContext(checkTx, tmproto.Header{})
	impl := internal.NewKeeper(
		app.AppCodec(),
		app.GetKey(foundation.ModuleName),
		app.MsgServiceRouter(),
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.FeeCollectorName,
		foundation.DefaultConfig(),
		foundation.DefaultAuthority().String(),
	)

	createAddress := func() sdk.AccAddress {
		return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	authority := sdk.MustAccAddressFromBech32(impl.GetAuthority())

	members := make([]sdk.AccAddress, 10)
	for i := range members {
		members[i] = createAddress()
	}
	impl.SetMember(ctx, foundation.Member{
		Address: members[0].String(),
	})

	info := foundation.DefaultFoundation()
	info.TotalWeight = sdk.NewDec(int64(len(members)))
	err := info.SetDecisionPolicy(workingPolicy())
	require.NoError(t, err)
	impl.SetFoundationInfo(ctx, info)

	// create proposals of different versions and abort them
	for _, newMember := range members[1:] {
		_, err := impl.SubmitProposal(ctx, []string{members[0].String()}, "", []sdk.Msg{testdata.NewTestMsg(authority)})
		require.NoError(t, err)

		err = impl.UpdateMembers(ctx, []foundation.MemberRequest{
			{
				Address: newMember.String(),
			},
		})
		require.NoError(t, err)
	}

	for _, proposal := range impl.GetProposals(ctx) {
		require.Equal(t, foundation.PROPOSAL_STATUS_ABORTED, proposal.Status)
	}
}
