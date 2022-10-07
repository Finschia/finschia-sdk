package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
)

// doTallyAndUpdate performs a tally, and, if the tally result is final, then:
// - updates the proposal's `Status` and `FinalTallyResult` fields,
// - prune all the votes.
func (k Keeper) doTallyAndUpdate(ctx sdk.Context, p *foundation.Proposal) error {
	tallyResult, err := k.tally(ctx, *p)
	if err != nil {
		return err
	}

	info := k.GetFoundationInfo(ctx)
	policy := info.GetDecisionPolicy()
	sinceSubmission := ctx.BlockTime().Sub(p.SubmitTime) // duration passed since proposal submission.
	result, err := policy.Allow(tallyResult, info.TotalWeight, sinceSubmission)
	if err != nil {
		return err
	}

	// If the result was final (i.e. enough votes to pass) or if the voting
	// period ended, then we consider the proposal as final.
	if isFinal := result.Final || ctx.BlockTime().After(p.VotingPeriodEnd); isFinal {
		k.pruneVotes(ctx, p.Id)
		p.FinalTallyResult = tallyResult
		if result.Allow {
			p.Status = foundation.PROPOSAL_STATUS_ACCEPTED
		} else {
			p.Status = foundation.PROPOSAL_STATUS_REJECTED
		}
	}

	return nil
}

// tally is a function that tallies a proposal by iterating through its votes,
// and returns the tally result without modifying the proposal or any state.
func (k Keeper) tally(ctx sdk.Context, p foundation.Proposal) (foundation.TallyResult, error) {
	// If proposal has already been tallied and updated, then its status is
	// accepted/rejected, in which case we just return the previously stored result.
	//
	// In all other cases (including withdrawn, aborted...) we do the tally
	// again.
	if p.Status == foundation.PROPOSAL_STATUS_ACCEPTED || p.Status == foundation.PROPOSAL_STATUS_REJECTED {
		return p.FinalTallyResult, nil
	}

	tallyResult := foundation.DefaultTallyResult()
	var errIter error
	k.iterateVotes(ctx, p.Id, func(vote foundation.Vote) (stop bool) {
		voter := sdk.MustAccAddressFromBech32(vote.Voter)

		_, err := k.GetMember(ctx, voter)
		switch {
		case sdkerrors.ErrNotFound.Is(err):
			// If the member left the foundation after voting, then we simply skip the
			// vote.
			return false
		case err != nil:
			// For any other errors, we stop and return the error.
			errIter = err
			return true
		}

		if err := tallyResult.Add(vote.Option); err != nil {
			panic(err)
		}

		return false
	})

	if errIter != nil {
		return tallyResult, errIter
	}

	return tallyResult, nil
}
