package keeper

import (
	"time"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
)

// handleUpdateFoundationParamsProposal is a handler for update foundation params proposal
func (k Keeper) handleUpdateFoundationParamsProposal(ctx sdk.Context, p *foundation.UpdateFoundationParamsProposal) error {
	// TODO: validate param changes
	params := p.Params
	k.SetParams(ctx, params)

	if !params.Enabled {
		k.Cleanup(ctx)
	}

	return ctx.EventManager().EmitTypedEvent(&foundation.EventUpdateFoundationParams{
		Params: params,
	})
}

// handleUpdateValidatorAuthsProposal is a handler for update validator auths proposal
func (k Keeper) handleUpdateValidatorAuthsProposal(ctx sdk.Context, p *foundation.UpdateValidatorAuthsProposal) error {
	for _, auth := range p.Auths {
		if err := k.SetValidatorAuth(ctx, auth); err != nil {
			return err
		}
	}

	return ctx.EventManager().EmitTypedEvent(&foundation.EventUpdateValidatorAuths{
		Auths: p.Auths,
	})
}

func (k Keeper) NewProposalId(ctx sdk.Context) uint64 {
	id := k.getPreviousProposalId(ctx) + 1
	k.setPreviousProposalId(ctx, id)

	return id
}

func (k Keeper) getPreviousProposalId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(previousProposalIdKey)
	if len(bz) == 0 {
		panic("previous proposal ID hasn't been set")
	}
	return Uint64FromBytes(bz)
}

func (k Keeper) setPreviousProposalId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(previousProposalIdKey, Uint64ToBytes(id))
}

func (k Keeper) submitProposal(ctx sdk.Context, proposers []string, metadata string, msgs []sdk.Msg) (uint64, error) {
	foundationInfo := k.GetFoundationInfo(ctx)
	if err := ensureMsgAuthz(msgs, sdk.AccAddress(foundationInfo.Operator)); err != nil {
		return 0, err
	}

	id := k.NewProposalId(ctx)
	proposal := foundation.Proposal{
		Id: id,
		Metadata: metadata,
		Proposers: proposers,
		SubmitTime: ctx.BlockTime(),
		FoundationVersion: foundationInfo.Version,
		Result: foundation.PROPOSAL_RESULT_UNFINALIZED,
		Status: foundation.PROPOSAL_STATUS_SUBMITTED,
		ExecutorResult: foundation.PROPOSAL_EXECUTOR_RESULT_NOT_RUN,
		VotingPeriodEnd: ctx.BlockTime().Add(foundationInfo.GetDecisionPolicy().GetVotingPeriod()),
		FinalTallyResult: DefaultTallyResult(),
	}
	if err := proposal.SetMsgs(msgs); err != nil {
		return 0, err
	}

	if err := k.setProposal(ctx, proposal); err != nil {
		return 0, err
	}
	k.addProposalToVPEndQueue(ctx, proposal)

	return id, nil
}

func (k Keeper) withdrawProposal(ctx sdk.Context, proposal foundation.Proposal) error {
	// Ensure the proposal can be withdrawn.
	if proposal.Status != foundation.PROPOSAL_STATUS_SUBMITTED {
		return sdkerrors.ErrInvalidRequest.Wrapf("cannot withdraw a proposal with the status of %s", proposal.Status.String())
	}

	proposal.Result = foundation.PROPOSAL_RESULT_UNFINALIZED
	proposal.Status = foundation.PROPOSAL_STATUS_WITHDRAWN

	k.setProposal(ctx, proposal)

	return nil
}

// pruneProposal deletes a proposal from state.
func (k Keeper) pruneProposal(ctx sdk.Context, proposal foundation.Proposal) {
	k.removeProposalFromVPEndQueue(ctx, proposal)
	k.deleteProposal(ctx, proposal.Id)
	k.pruneVotes(ctx, proposal.Id)
}

// pruneExpiredProposals prunes all proposals:
// 1. which are expired, i.e. whose `submit_time + voting_period + max_execution_period` is greater than now.
func (k Keeper) pruneExpiredProposals(ctx sdk.Context) {
	votingPeriodEnd := ctx.BlockTime().Add(-k.config.MaxExecutionPeriod)

	var proposals []foundation.Proposal
	k.iterateProposalsByVPEnd(ctx, votingPeriodEnd, func(proposal foundation.Proposal) (stop bool) {
		proposals = append(proposals, proposal)
		return false
	})

	for _, proposal := range proposals {
		k.pruneProposal(ctx, proposal)
	}
}

// pruneOldProposals prunes all proposals:
// 2. which have lower version than the current foundation's
func (k Keeper) pruneOldProposals(ctx sdk.Context) {
	latestVersion := k.GetFoundationInfo(ctx).Version

	var proposals []foundation.Proposal
	k.iterateProposals(ctx, func(proposal foundation.Proposal) (stop bool) {
		if proposal.FoundationVersion == latestVersion {
			return true
		}

		// one may still execute the finalized proposals
		if proposal.Result == foundation.PROPOSAL_RESULT_UNFINALIZED {
			proposals = append(proposals, proposal)
		}

		return false
	})

	for _, proposal := range proposals {
		k.pruneProposal(ctx, proposal)
	}
}

func (k Keeper) GetProposals(ctx sdk.Context) []foundation.Proposal {
	var proposals []foundation.Proposal
	k.iterateProposals(ctx, func(proposal foundation.Proposal) (stop bool) {
		proposals = append(proposals, proposal)
		return false
	})

	return proposals
}

func (k Keeper) iterateProposals(ctx sdk.Context, fn func(proposal foundation.Proposal) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefix := append(proposalKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var proposal foundation.Proposal
		k.cdc.MustUnmarshal(iterator.Value(), &proposal)
		if stop := fn(proposal); stop {
			break
		}
	}
}

func (k Keeper) iterateProposalsByVPEnd(ctx sdk.Context, endTime time.Time, fn func(proposal foundation.Proposal) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := store.Iterator(proposalByVPEndKeyPrefix, sdk.PrefixEndBytes(append(proposalByVPEndKeyPrefix, sdk.FormatTimeBytes(endTime)...)))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var proposal foundation.Proposal
		k.cdc.MustUnmarshal(iter.Value(), &proposal)

		if fn(proposal) {
			break
		}
	}
}

func (k Keeper) UpdateTallyOfVPEndProposals(ctx sdk.Context) {
	k.iterateProposalsByVPEnd(ctx, ctx.BlockTime(), func(proposal foundation.Proposal) (stop bool) {
		if err := k.doTallyAndUpdate(ctx, &proposal); err != nil {
			panic(err)
		}

		if err := k.setProposal(ctx, proposal); err != nil {
			panic(err)
		}

		return false
	})
}

func (k Keeper) GetProposal(ctx sdk.Context, id uint64) (*foundation.Proposal, error) {
	store := ctx.KVStore(k.storeKey)
	key := proposalKey(id)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "No such proposal for id: %d", id)
	}

	var proposal foundation.Proposal
	if err := k.cdc.Unmarshal(bz, &proposal); err != nil {
		return nil, err
	}
	return &proposal, nil
}

func (k Keeper) setProposal(ctx sdk.Context, proposal foundation.Proposal) error {
	bz, err := k.cdc.Marshal(&proposal)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	key := proposalKey(proposal.Id)
	store.Set(key, bz)

	return nil
}

func (k Keeper) deleteProposal(ctx sdk.Context, proposalId uint64) {
	store := ctx.KVStore(k.storeKey)
	key := proposalKey(proposalId)
	store.Delete(key)
}

func (k Keeper) addProposalToVPEndQueue(ctx sdk.Context, proposal foundation.Proposal) {
	store := ctx.KVStore(k.storeKey)
	key := proposalByVPEndKey(proposal.Id, proposal.VotingPeriodEnd)
	store.Set(key, []byte{})
}

func (k Keeper) removeProposalFromVPEndQueue(ctx sdk.Context, proposal foundation.Proposal) {
	store := ctx.KVStore(k.storeKey)
	key := proposalByVPEndKey(proposal.Id, proposal.VotingPeriodEnd)
	store.Delete(key)
}

func validateActorForProposal(address string, proposal foundation.Proposal) error {
	proposers := map[string]bool{}
	for _, proposer := range proposal.Proposers {
		proposers[proposer] = true
	}
	if !proposers[address] {
		return sdkerrors.ErrUnauthorized.Wrapf("given address is not in proposers: %s", address)
	}

	return nil
}
