package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
)

func (k Keeper) Vote(ctx sdk.Context, vote foundation.Vote) error {
	if err := validateMetadata(vote.Metadata, k.config); err != nil {
		return err
	}

	// Make sure that a voter hasn't already voted.
	if k.hasVote(ctx, vote.ProposalId, sdk.AccAddress(vote.Voter)) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Already voted: %s", vote.Voter)
	}

	proposal, err := k.GetProposal(ctx, vote.ProposalId)
	if err != nil {
		return err
	}

	// Ensure that we can still accept votes for this proposal.
	if proposal.Status != foundation.PROPOSAL_STATUS_SUBMITTED {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "not possible with proposal status: %s", proposal.Status.String())
	}
	if !ctx.BlockTime().Before(proposal.VotingPeriodEnd) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "voting period has ended already")
	}

	vote.SubmitTime = ctx.BlockTime()

	if err := k.setVote(ctx, vote); err != nil {
		return err
	}

	return ctx.EventManager().EmitTypedEvent(&foundation.EventVote{
		Vote: vote,
	})
}

func (k Keeper) hasVote(ctx sdk.Context, proposalID uint64, voter sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	key := voteKey(proposalID, voter)
	return store.Has(key)
}

func (k Keeper) GetVote(ctx sdk.Context, proposalID uint64, voter sdk.AccAddress) (*foundation.Vote, error) {
	store := ctx.KVStore(k.storeKey)
	key := voteKey(proposalID, voter)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "No vote for proposal %d: %s", proposalID, voter)
	}

	var vote foundation.Vote
	if err := k.cdc.Unmarshal(bz, &vote); err != nil {
		return nil, err
	}
	return &vote, nil
}

func (k Keeper) setVote(ctx sdk.Context, vote foundation.Vote) error {
	bz, err := k.cdc.Marshal(&vote)
	if err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	key := voteKey(vote.ProposalId, sdk.AccAddress(vote.Voter))
	store.Set(key, bz)

	return nil
}

func (k Keeper) iterateVotes(ctx sdk.Context, proposalID uint64, fn func(vote foundation.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefix := append(voteKeyPrefix, Uint64ToBytes(proposalID)...)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var vote foundation.Vote
		k.cdc.MustUnmarshal(iterator.Value(), &vote)
		if stop := fn(vote); stop {
			break
		}
	}
}

func (k Keeper) GetVotes(ctx sdk.Context, proposalID uint64) []foundation.Vote {
	var votes []foundation.Vote
	k.iterateVotes(ctx, proposalID, func(vote foundation.Vote) (stop bool) {
		votes = append(votes, vote)
		return false
	})

	return votes
}

// pruneVotes prunes all votes for a proposal from state.
func (k Keeper) pruneVotes(ctx sdk.Context, proposalID uint64) {
	keys := [][]byte{}
	k.iterateVotes(ctx, proposalID, func(vote foundation.Vote) (stop bool) {
		keys = append(keys, voteKey(proposalID, sdk.AccAddress(vote.Voter)))
		return false
	})

	store := ctx.KVStore(k.storeKey)
	for _, key := range keys {
		store.Delete(key)
	}
}
