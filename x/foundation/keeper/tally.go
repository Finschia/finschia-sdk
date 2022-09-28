package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
)

// doTallyAndUpdate performs a tally, and updates the proposal's
// `FinalTallyResult` field only if the tally is final.
func (k Keeper) doTallyAndUpdate(ctx sdk.Context, p *foundation.Proposal) error {
	pSubmittedAt, err := gogotypes.TimestampProto(p.SubmitTime)
	if err != nil {
		return err
	}
	submittedAt, err := gogotypes.TimestampFromProto(pSubmittedAt)
	if err != nil {
		return err
	}

	tallyResult, err := k.tally(ctx, *p)
	if err != nil {
		return err
	}

	info := k.GetFoundationInfo(ctx)
	policy := info.GetDecisionPolicy()
	decision, err := policy.Allow(tallyResult, info.TotalWeight, ctx.BlockTime().Sub(submittedAt))
	switch {
	case err != nil:
		return err
	case decision.Final:
		k.pruneVotes(ctx, p.Id)
		p.FinalTallyResult = tallyResult
		if decision.Allow {
			p.Result = foundation.PROPOSAL_RESULT_ACCEPTED
			p.Status = foundation.PROPOSAL_STATUS_CLOSED
		} else {
			p.Result = foundation.PROPOSAL_RESULT_REJECTED
			p.Status = foundation.PROPOSAL_STATUS_CLOSED
		}
	}

	return nil
}

// tally is a function that tallies a proposal by iterating through its votes,
// and returns the tally result without modifying the proposal or any state.
func (k Keeper) tally(ctx sdk.Context, p foundation.Proposal) (foundation.TallyResult, error) {
	// If proposal has already been tallied and updated, then its status is
	// closed, in which case we just return the previously stored result.
	if p.Status == foundation.PROPOSAL_STATUS_CLOSED {
		return p.FinalTallyResult, nil
	}

	tallyResult := foundation.DefaultTallyResult()
	var errors []error
	k.iterateVotes(ctx, p.Id, func(vote foundation.Vote) (stop bool) {
		voter, err := sdk.AccAddressFromBech32(vote.Voter)
		if err != nil {
			errors = append(errors, err)
			return true
		}
		_, err = k.GetMember(ctx, voter)
		switch {
		case sdkerrors.ErrNotFound.Is(err):
			// If the member left the group after voting, then we simply skip the
			// vote.
			return false
		case err != nil:
			// For any other errors, we stop and return the error.
			errors = append(errors, err)
			return true
		}

		if err := tallyResult.Add(vote.Option); err != nil {
			errors = append(errors, err)
			return true
		}

		return false
	})

	if len(errors) != 0 {
		return tallyResult, errors[0]
	}

	return tallyResult, nil
}
