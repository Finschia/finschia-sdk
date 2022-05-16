package keeper

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"

	"github.com/line/lbm-sdk/x/foundation"
)

// ensureMsgAuthz checks that if a message requires signers that all of them are equal to the given account address of the operator.
func ensureMsgAuthz(msgs []sdk.Msg, operator sdk.AccAddress) error {
	for _, msg := range msgs {
		for _, signer := range msg.GetSigners() {
			if !operator.Equals(signer) {
				return sdkerrors.ErrUnauthorized.Wrapf("signer of msg is not the operator: %s, %s", signer, operator)
			}
		}
	}
	return nil
}

func (k Keeper) Exec(ctx sdk.Context, proposalID uint64) error {
	proposal, err := k.GetProposal(ctx, proposalID)
	if err != nil {
		return err
	}

	if proposal.Status != foundation.PROPOSAL_STATUS_SUBMITTED &&
		proposal.Status != foundation.PROPOSAL_STATUS_CLOSED {
		return sdkerrors.ErrInvalidRequest.Wrapf("not possible with proposal status: %s", proposal.Status.String())
	}

	if proposal.Status == foundation.PROPOSAL_STATUS_SUBMITTED {
		if err := k.doTallyAndUpdate(ctx, proposal); err != nil {
			return err
		}
	}

	// Execute proposal payload.
	if proposal.Status == foundation.PROPOSAL_STATUS_CLOSED &&
		proposal.Result == foundation.PROPOSAL_RESULT_ACCEPTED &&
		proposal.ExecutorResult != foundation.PROPOSAL_EXECUTOR_RESULT_SUCCESS {
		logger := ctx.Logger().With("module", fmt.Sprintf("x/%s", foundation.ModuleName))
		// Caching context so that we don't update the store in case of failure.
		ctx, flush := ctx.CacheContext()

		if _, err = k.doExecuteMsgs(ctx, *proposal); err != nil {
			proposal.ExecutorResult = foundation.PROPOSAL_EXECUTOR_RESULT_FAILURE
			logger.Info("proposal execution failed", "cause", err, "proposalID", proposal.Id)
		} else {
			proposal.ExecutorResult = foundation.PROPOSAL_EXECUTOR_RESULT_SUCCESS
			flush()
		}
	}

	// If proposal has successfully run, delete it from state.
	if proposal.ExecutorResult == foundation.PROPOSAL_EXECUTOR_RESULT_SUCCESS {
		k.pruneProposal(ctx, *proposal)
	} else {
		if err := k.setProposal(ctx, *proposal); err != nil {
			return err
		}
	}

	return ctx.EventManager().EmitTypedEvent(&foundation.EventExec{
		ProposalId: proposal.Id,
		Result:     proposal.ExecutorResult,
	})
}

// doExecuteMsgs routes the messages to the registered handlers.
func (k Keeper) doExecuteMsgs(ctx sdk.Context, proposal foundation.Proposal) ([]sdk.Result, error) {
	msgs := proposal.GetMsgs()
	results := make([]sdk.Result, len(msgs))
	if err := ensureMsgAuthz(msgs, k.GetOperator(ctx)); err != nil {
		return nil, err
	}
	for i, msg := range msgs {
		handler := k.router.Handler(msg)
		if handler == nil {
			return nil, sdkerrors.ErrUnknownRequest.Wrapf("no message handler found for %q", sdk.MsgTypeURL(msg))
		}
		r, err := handler(ctx, msg)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "message %q at position %d", msg, i)
		}
		if r != nil {
			results[i] = *r
		}
	}
	return results, nil
}
