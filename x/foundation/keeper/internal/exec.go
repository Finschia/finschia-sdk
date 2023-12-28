package internal

import (
	"bytes"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

// ensureMsgAuthz checks that if a message requires signers that all of them are equal to the given account address of the authority.
func ensureMsgAuthz(msgs []sdk.Msg, authority sdk.AccAddress, cdc codec.Codec) error {
	for _, msg := range msgs {
		// In practice, GetSigners() should return a non-empty array without
		// duplicates, so the code below is equivalent to:
		// `msgs[i].GetSigners()[0] == authority`
		// but we prefer to loop through all GetSigners just to be sure.
		signers, _, err := cdc.GetMsgV1Signers(msg)
		if err != nil {
			return err
		}

		for _, signer := range signers {
			if !bytes.Equal(signer, authority) {
				return sdkerrors.ErrUnauthorized.Wrapf("bad signer; expected %s, got %s", authority, signer)
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
		proposal.Status != foundation.PROPOSAL_STATUS_ACCEPTED {
		return sdkerrors.ErrInvalidRequest.Wrapf("not possible with proposal status: %s", proposal.Status)
	}

	if proposal.Status == foundation.PROPOSAL_STATUS_SUBMITTED {
		if err := k.doTallyAndUpdate(ctx, proposal); err != nil {
			return err
		}
	}

	// Execute proposal payload.
	var logs string
	if proposal.Status == foundation.PROPOSAL_STATUS_ACCEPTED &&
		proposal.ExecutorResult != foundation.PROPOSAL_EXECUTOR_RESULT_SUCCESS {
		logger := ctx.Logger().With("module", fmt.Sprintf("x/%s", foundation.ModuleName))
		// Caching context so that we don't update the store in case of failure.
		cacheCtx, flush := ctx.CacheContext()

		if results, err := k.doExecuteMsgs(cacheCtx, *proposal); err != nil {
			proposal.ExecutorResult = foundation.PROPOSAL_EXECUTOR_RESULT_FAILURE
			logs = fmt.Sprintf("proposal execution failed on proposal %d, because of error %s", proposalID, err.Error())
			logger.Info("proposal execution failed", "cause", err, "proposalID", proposal.Id)
		} else {
			proposal.ExecutorResult = foundation.PROPOSAL_EXECUTOR_RESULT_SUCCESS
			flush()

			for _, res := range results {
				// NOTE: The sdk msg handler creates a new EventManager, so events must be correctly propagated back to the current context
				ctx.EventManager().EmitEvents(res.GetEvents())
			}
		}
	}

	// If proposal has successfully run, delete it from state.
	if proposal.ExecutorResult == foundation.PROPOSAL_EXECUTOR_RESULT_SUCCESS {
		k.pruneProposal(ctx, *proposal)
	} else {
		k.setProposal(ctx, *proposal)
	}

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventExec{
		ProposalId: proposal.Id,
		Logs:       logs,
		Result:     proposal.ExecutorResult,
	}); err != nil {
		panic(err)
	}

	return nil
}

// doExecuteMsgs routes the messages to the registered handlers.
func (k Keeper) doExecuteMsgs(ctx sdk.Context, proposal foundation.Proposal) ([]sdk.Result, error) {
	msgs := proposal.GetMsgs()
	results := make([]sdk.Result, len(msgs))

	authority, err := k.addressCodec.StringToBytes(k.GetAuthority())
	if err != nil {
		panic(err)
	}

	if err := ensureMsgAuthz(msgs, authority, k.cdc); err != nil {
		return nil, err
	}

	for i, msg := range msgs {
		handler := k.router.Handler(msg)
		if handler == nil {
			return nil, sdkerrors.ErrUnknownRequest.Wrapf("no message handler found for %q", sdk.MsgTypeURL(msg))
		}
		r, err := handler(ctx, msg)
		if err != nil {
			return nil, errorsmod.Wrapf(err, "message %q at position %d", msg, i)
		}
		if r != nil {
			results[i] = *r
		}
	}
	return results, nil
}
