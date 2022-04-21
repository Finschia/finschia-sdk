package keeper

import (
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
	id := k.getProposalId(ctx)
	k.setProposalId(ctx, id + 1)

	return id
}

func (k Keeper) getProposalId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(nextProposalIdKey)
	if len(bz) == 0 {
		panic("next proposal ID hasn't been set")
	}
	return Uint64FromBytes(bz)
}

func (k Keeper) setProposalId(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(nextProposalIdKey, Uint64ToBytes(id))
}

func (k Keeper) submitProposal(ctx sdk.Context, proposers []string, metadata string, msgs []sdk.Msg) (uint64, error) {
	for _, proposer := range proposers {
		if _, err := k.GetMember(ctx, sdk.AccAddress(proposer)); err != nil {
			return 0, err
		}
	}

	foundationInfo := k.GetFoundationInfo(ctx)
	if err := ensureMsgAuthz(msgs, sdk.AccAddress(foundationInfo.Operator)); err != nil {
		return 0, err
	}

	// TODO: validate policy

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
func (k Keeper) pruneProposal(ctx sdk.Context, proposalId uint64) {
	k.deleteProposal(ctx, proposalId)
	k.pruneVotes(ctx, proposalId)
}

// pruneProposal deletes a proposal from state.
func (k Keeper) pruneProposals(ctx sdk.Context) {
	var ids []uint64
	k.iterateProposals(ctx, func(proposal foundation.Proposal) (stop bool) {
		ids = append(ids, proposal.Id)
		return false
	})

	for _, id := range ids {
		k.pruneProposal(ctx, id)
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
