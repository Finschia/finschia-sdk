package keeper

import (
	"time"

	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
)

// DefaultTallyResult returns a TallyResult with all counts set to 0.
func DefaultTallyResult() foundation.TallyResult {
	return foundation.TallyResult{
		YesCount:        sdk.ZeroDec(),
		NoCount:         sdk.ZeroDec(),
		NoWithVetoCount: sdk.ZeroDec(),
		AbstainCount:    sdk.ZeroDec(),
	}
}

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
	allow, final, err := makeDecision(ctx, info, tallyResult, ctx.BlockTime().Sub(submittedAt))
	switch {
	case err != nil:
		return err
	case final:
		if err := k.pruneVotes(ctx, p.Id); err != nil {
			return err
		}
		p.FinalTallyResult = tallyResult
		if allow {
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

	tallyResult := DefaultTallyResult()
	var errors []error
	k.iterateVotes(ctx, p.Id, func(vote foundation.Vote) (stop bool) {
		voter, err := k.GetMember(ctx, sdk.AccAddress(vote.Voter))
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

		if err := tallyResult.Add(vote.Option, voter.Weight); err != nil {
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

func makeDecision(ctx sdk.Context, info foundation.FoundationInfo, result foundation.TallyResult, sinceSubmission time.Duration) (allow bool, final bool, err error) {
	policy := info.DecisionPolicy
	if sinceSubmission < policy.Windows.MinExecutionPeriod {
		err = sdkerrors.ErrUnauthorized.Wrapf("must wait %s after submission before execution, currently at %s", policy.Windows.MinExecutionPeriod, sinceSubmission)
		return
	}

	decisions := []func(sdk.Context, foundation.FoundationInfo, foundation.TallyResult)(allow bool, final bool){
		makeDecisionThreshold,
		makeDecisionPercentage,
	}
	allow, final = false, true
	for _, decision := range decisions {
		a, f := decision(ctx, info, result)
		if a {
			allow, final = true, true
			return
		}
		if !f {
			final = false
		}
	}

	return 
}

func makeDecisionThreshold(ctx sdk.Context, info foundation.FoundationInfo, result foundation.TallyResult) (allow bool, final bool) {
	// the real threshold of the policy is `min(threshold,total_weight)`. If
	// the foundation member weights changes (member leaving, member weight update)
	// and the threshold doesn't, we can end up with threshold > total_weight.
	// In this case, as long as everyone votes yes (in which case
	// `yesCount`==`realThreshold`), then the proposal still passes.
	realThreshold := sdk.MinDec(info.DecisionPolicy.Threshold, info.TotalWeight)
	if result.YesCount.GTE(realThreshold) {
		allow, final = true, true
		return
	}

	totalCounts := result.TotalCounts()
	undecided := info.TotalWeight.Sub(totalCounts)

	// maxYesCount is the max potential number of yes count, i.e the current yes count
	// plus all undecided count (supposing they all vote yes).
	maxYesCount := result.YesCount.Add(undecided)
	if maxYesCount.LT(realThreshold) {
		allow, final = false, true
		return
	}

	return
}

func makeDecisionPercentage(ctx sdk.Context, info foundation.FoundationInfo, result foundation.TallyResult) (allow bool, final bool) {
	yesPercentage := result.YesCount.Quo(info.TotalWeight)
	if yesPercentage.GTE(info.DecisionPolicy.Percentage) {
		allow, final = true, true
		return
	}

	totalCounts := result.TotalCounts()
	undecided := info.TotalWeight.Sub(totalCounts)

	maxYesCount := result.YesCount.Add(undecided)
	maxYesPercentage := maxYesCount.Quo(info.TotalWeight)
	if maxYesPercentage.LT(info.DecisionPolicy.Percentage) {
		allow, final = false, true
		return
	}

	return
}
