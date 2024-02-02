package internal

import (
	"time"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (k Keeper) newProposalID(ctx sdk.Context) uint64 {
	id := k.getPreviousProposalID(ctx) + 1
	k.setPreviousProposalID(ctx, id)

	return id
}

func (k Keeper) getPreviousProposalID(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(previousProposalIDKey)
	if err != nil {
		panic(err)
	}
	if len(bz) == 0 {
		panic("previous proposal ID hasn't been set")
	}
	return Uint64FromBytes(bz)
}

func (k Keeper) setPreviousProposalID(ctx sdk.Context, id uint64) {
	store := k.storeService.OpenKVStore(ctx)
	if err := store.Set(previousProposalIDKey, Uint64ToBytes(id)); err != nil {
		panic(err)
	}
}

func (k Keeper) SubmitProposal(ctx sdk.Context, proposers []string, metadata string, msgs []sdk.Msg) (*uint64, error) {
	if err := validateMetadata(metadata, k.config); err != nil {
		return nil, err
	}

	foundationInfo := k.GetFoundationInfo(ctx)
	authority, err := k.addressCodec().StringToBytes(k.GetAuthority())
	if err != nil {
		panic(err)
	}
	if err := ensureMsgAuthz(msgs, authority, k.cdc); err != nil {
		return nil, err
	}

	// Prevent proposal that can not succeed.
	policy := foundationInfo.GetDecisionPolicy()
	if err := policy.Validate(foundationInfo, k.config); err != nil {
		return nil, err
	}

	id := k.newProposalID(ctx)
	proposal := foundation.Proposal{
		Id:                id,
		Metadata:          metadata,
		Proposers:         proposers,
		SubmitTime:        ctx.BlockTime(),
		FoundationVersion: foundationInfo.Version,
		Status:            foundation.PROPOSAL_STATUS_SUBMITTED,
		ExecutorResult:    foundation.PROPOSAL_EXECUTOR_RESULT_NOT_RUN,
		VotingPeriodEnd:   ctx.BlockTime().Add(policy.GetVotingPeriod()),
		FinalTallyResult:  foundation.DefaultTallyResult(),
	}
	if err := proposal.SetMsgs(msgs); err != nil {
		return nil, err
	}

	k.setProposal(ctx, proposal)
	k.addProposalToVPEndQueue(ctx, proposal)

	return &id, nil
}

func (k Keeper) WithdrawProposal(ctx sdk.Context, proposalID uint64) error {
	proposal, err := k.GetProposal(ctx, proposalID)
	if err != nil {
		return err
	}

	// Ensure the proposal can be withdrawn.
	if proposal.Status != foundation.PROPOSAL_STATUS_SUBMITTED {
		return sdkerrors.ErrInvalidRequest.Wrapf("cannot withdraw a proposal with the status of %s", proposal.Status)
	}

	proposal.Status = foundation.PROPOSAL_STATUS_WITHDRAWN
	k.setProposal(ctx, *proposal)

	return nil
}

// pruneProposal deletes a proposal from state.
func (k Keeper) pruneProposal(ctx sdk.Context, proposal foundation.Proposal) {
	k.pruneVotes(ctx, proposal.Id)
	k.removeProposalFromVPEndQueue(ctx, proposal)
	k.deleteProposal(ctx, proposal.Id)
}

// PruneExpiredProposals prunes all proposals which are expired,
// i.e. whose `submit_time + voting_period + max_execution_period` is smaller than (or equal to) now.
func (k Keeper) PruneExpiredProposals(ctx sdk.Context) {
	votingPeriodEnd := ctx.BlockTime().Add(-k.config.MaxExecutionPeriod).Add(time.Nanosecond)

	var proposals []foundation.Proposal
	k.iterateProposalsByVPEnd(ctx, votingPeriodEnd, func(proposal foundation.Proposal) (stop bool) {
		proposals = append(proposals, proposal)
		return false
	})

	for _, proposal := range proposals {
		k.pruneProposal(ctx, proposal)
	}
}

// abortOldProposals aborts all proposals which have lower version than the current foundation's
func (k Keeper) abortOldProposals(ctx sdk.Context) {
	latestVersion := k.GetFoundationInfo(ctx).Version

	k.iterateProposals(ctx, func(proposal foundation.Proposal) (stop bool) {
		if proposal.FoundationVersion == latestVersion {
			return true
		}

		if proposal.Status == foundation.PROPOSAL_STATUS_SUBMITTED {
			k.pruneVotes(ctx, proposal.Id)

			proposal.Status = foundation.PROPOSAL_STATUS_ABORTED
			k.setProposal(ctx, proposal)
		}

		return false
	})
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
	store := k.storeService.OpenKVStore(ctx)
	prefix := proposalKeyPrefix
	adapter := runtime.KVStoreAdapter(store)
	iterator := storetypes.KVStorePrefixIterator(adapter, prefix)
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
	store := k.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	iter := adapter.Iterator(proposalByVPEndKeyPrefix, storetypes.PrefixEndBytes(append(proposalByVPEndKeyPrefix, sdk.FormatTimeBytes(endTime)...)))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		_, id := splitProposalByVPEndKey(iter.Key())

		proposal, err := k.GetProposal(ctx, id)
		if err != nil {
			panic(err)
		}

		if fn(*proposal) {
			break
		}
	}
}

func (k Keeper) UpdateTallyOfVPEndProposals(ctx sdk.Context) error {
	var proposals []foundation.Proposal
	k.iterateProposalsByVPEnd(ctx, ctx.BlockTime(), func(proposal foundation.Proposal) (stop bool) {
		proposals = append(proposals, proposal)
		return false
	})

	for _, proposal := range proposals {
		proposal := proposal

		if proposal.Status == foundation.PROPOSAL_STATUS_ABORTED || proposal.Status == foundation.PROPOSAL_STATUS_WITHDRAWN {
			k.pruneProposal(ctx, proposal)
			continue
		}

		if err := k.doTallyAndUpdate(ctx, &proposal); err != nil {
			return err
		}
		k.setProposal(ctx, proposal)
	}

	return nil
}

func (k Keeper) GetProposal(ctx sdk.Context, id uint64) (*foundation.Proposal, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := proposalKey(id)
	bz, err := store.Get(key)
	if err != nil {
		return nil, err
	}
	if len(bz) == 0 {
		return nil, sdkerrors.ErrNotFound.Wrapf("No proposal for id: %d", id)
	}

	var proposal foundation.Proposal
	k.cdc.MustUnmarshal(bz, &proposal)

	return &proposal, nil
}

func (k Keeper) setProposal(ctx sdk.Context, proposal foundation.Proposal) {
	store := k.storeService.OpenKVStore(ctx)
	key := proposalKey(proposal.Id)

	bz := k.cdc.MustMarshal(&proposal)
	if err := store.Set(key, bz); err != nil {
		panic(err)
	}
}

func (k Keeper) deleteProposal(ctx sdk.Context, proposalID uint64) {
	store := k.storeService.OpenKVStore(ctx)
	key := proposalKey(proposalID)
	if err := store.Delete(key); err != nil {
		panic(err)
	}
}

func (k Keeper) addProposalToVPEndQueue(ctx sdk.Context, proposal foundation.Proposal) {
	store := k.storeService.OpenKVStore(ctx)
	key := proposalByVPEndKey(proposal.VotingPeriodEnd, proposal.Id)
	if err := store.Set(key, []byte{}); err != nil {
		panic(err)
	}
}

func (k Keeper) removeProposalFromVPEndQueue(ctx sdk.Context, proposal foundation.Proposal) {
	store := k.storeService.OpenKVStore(ctx)
	key := proposalByVPEndKey(proposal.VotingPeriodEnd, proposal.Id)
	if err := store.Delete(key); err != nil {
		panic(err)
	}
}

func validateActorForProposal(address string, proposal foundation.Proposal) error {
	for _, proposer := range proposal.Proposers {
		if address == proposer {
			return nil
		}
	}

	return sdkerrors.ErrUnauthorized.Wrapf("not a proposer: %s", address)
}
