package keeper_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/gov/types"
)

func TestGetSetProposal(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	app.GovKeeper.SetProposal(ctx, proposal)

	gotProposal, ok := app.GovKeeper.GetProposal(ctx, proposalID)
	require.True(t, ok)
	require.True(t, proposal.Equal(gotProposal))
}

func TestActivateVotingPeriod(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	tp := TestProposal
	proposal, err := app.GovKeeper.SubmitProposal(ctx, tp)
	require.NoError(t, err)

	require.True(t, proposal.VotingStartTime.Equal(time.Time{}))

	app.GovKeeper.ActivateVotingPeriod(ctx, proposal)

	require.True(t, proposal.VotingStartTime.Equal(ctx.BlockHeader().Time))

	proposal, ok := app.GovKeeper.GetProposal(ctx, proposal.ProposalId)
	require.True(t, ok)

	activeIterator := app.GovKeeper.ActiveProposalQueueIterator(ctx, proposal.VotingEndTime)
	require.True(t, activeIterator.Valid())

	proposalID := types.GetProposalIDFromBytes(activeIterator.Value())
	require.Equal(t, proposalID, proposal.ProposalId)
	activeIterator.Close()
}

type invalidProposalRoute struct{ types.TextProposal }

func (invalidProposalRoute) ProposalRoute() string { return "nonexistingroute" }

func TestSubmitProposal(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	testCases := []struct {
		content     types.Content
		expectedErr error
	}{
		{&types.TextProposal{Title: "title", Description: "description"}, nil},
		// Keeper does not check the validity of title and description, no error
		{&types.TextProposal{Title: "", Description: "description"}, nil},
		{&types.TextProposal{Title: strings.Repeat("1234567890", 100), Description: "description"}, nil},
		{&types.TextProposal{Title: "title", Description: ""}, nil},
		{&types.TextProposal{Title: "title", Description: strings.Repeat("1234567890", 1000)}, nil},
		// error only when invalid route
		{&invalidProposalRoute{}, types.ErrNoProposalHandlerExists},
	}

	for i, tc := range testCases {
		_, err := app.GovKeeper.SubmitProposal(ctx, tc.content)
		require.True(t, errors.Is(tc.expectedErr, err), "tc #%d; got: %v, expected: %v", i, err, tc.expectedErr)
	}
}

func TestGetProposalsFiltered(t *testing.T) {
	proposalID := uint64(1)
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	status := []types.ProposalStatus{types.StatusDepositPeriod, types.StatusVotingPeriod}

	addr1 := sdk.AccAddress("foo_________________")

	for _, s := range status {
		for i := 0; i < 50; i++ {
			p, err := types.NewProposal(TestProposal, proposalID, time.Now(), time.Now())
			require.NoError(t, err)

			p.Status = s

			if i%2 == 0 {
				d := types.NewDeposit(proposalID, addr1, nil)
				v := types.NewVote(proposalID, addr1, types.NewNonSplitVoteOption(types.OptionYes))
				app.GovKeeper.SetDeposit(ctx, d)
				app.GovKeeper.SetVote(ctx, v)
			}

			app.GovKeeper.SetProposal(ctx, p)
			proposalID++
		}
	}

	testCases := []struct {
		params             types.QueryProposalsParams
		expectedNumResults int
	}{
		{types.NewQueryProposalsParams(1, 50, types.StatusNil, "", ""), 50},
		{types.NewQueryProposalsParams(1, 50, types.StatusDepositPeriod, "", ""), 50},
		{types.NewQueryProposalsParams(1, 50, types.StatusVotingPeriod, "", ""), 50},
		{types.NewQueryProposalsParams(1, 25, types.StatusNil, "", ""), 25},
		{types.NewQueryProposalsParams(2, 25, types.StatusNil, "", ""), 25},
		{types.NewQueryProposalsParams(1, 50, types.StatusRejected, "", ""), 0},
		{types.NewQueryProposalsParams(1, 50, types.StatusNil, addr1, ""), 50},
		{types.NewQueryProposalsParams(1, 50, types.StatusNil, "", addr1), 50},
		{types.NewQueryProposalsParams(1, 50, types.StatusNil, addr1, addr1), 50},
		{types.NewQueryProposalsParams(1, 50, types.StatusDepositPeriod, addr1, addr1), 25},
		{types.NewQueryProposalsParams(1, 50, types.StatusDepositPeriod, "", ""), 50},
		{types.NewQueryProposalsParams(1, 50, types.StatusVotingPeriod, "", ""), 50},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test Case %d", i), func(t *testing.T) {
			proposals := app.GovKeeper.GetProposalsFiltered(ctx, tc.params)
			require.Len(t, proposals, tc.expectedNumResults)

			for _, p := range proposals {
				if types.ValidProposalStatus(tc.params.ProposalStatus) {
					require.Equal(t, tc.params.ProposalStatus, p.Status)
				}
			}
		})
	}
}
