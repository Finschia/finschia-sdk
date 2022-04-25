package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
)

func (k Keeper) vote(ctx sdk.Context, id uint64, voter string, option foundation.VoteOption, metadata string) error {
	// Make sure that a voter hasn't already voted.
	if k.hasVote(ctx, id, sdk.AccAddress(voter)) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Already voted: %s", voter)
	}

	proposal, err := k.GetProposal(ctx, id)
	if err != nil {
		return err
	}

	// Ensure that we can still accept votes for this proposal.
	if proposal.Status != foundation.PROPOSAL_STATUS_SUBMITTED {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "proposal not open for voting")
	}
	if ctx.BlockTime().After(proposal.VotingPeriodEnd) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "voting period has ended already")
	}

	newVote := foundation.Vote{
		ProposalId: id,
		Voter:      voter,
		Option:     option,
		Metadata:   metadata,
		SubmitTime: ctx.BlockTime(),
	}

	return k.setVote(ctx, newVote)
}

func (k Keeper) hasVote(ctx sdk.Context, proposalId uint64, voter sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	key := voteKey(proposalId, voter)
	return store.Has(key)
}

func (k Keeper) GetVote(ctx sdk.Context, proposalId uint64, voter sdk.AccAddress) (*foundation.Vote, error) {
	store := ctx.KVStore(k.storeKey)
	key := voteKey(proposalId, voter)
	bz := store.Get(key)
	if len(bz) == 0 {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrNotFound, "No vote for proposal %d: %s", proposalId, voter)
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

func (k Keeper) iterateVotes(ctx sdk.Context, proposalId uint64, fn func(vote foundation.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefix := append(voteKeyPrefix, Uint64ToBytes(proposalId)...)
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

func (k Keeper) GetVotes(ctx sdk.Context, proposalId uint64) []foundation.Vote {
	var votes []foundation.Vote
	k.iterateVotes(ctx, proposalId, func(vote foundation.Vote) (stop bool) {
		votes = append(votes, vote)
		return false
	})

	return votes
}

// pruneVotes prunes all votes for a proposal from state.
func (k Keeper) pruneVotes(ctx sdk.Context, proposalId uint64) {
	keys := [][]byte{}
	k.iterateVotes(ctx, proposalId, func(vote foundation.Vote) (stop bool) {
		keys = append(keys, voteKey(proposalId, sdk.AccAddress(vote.Voter)))
		return false
	})

	store := ctx.KVStore(k.storeKey)
	for _, key := range keys {
		store.Delete(key)
	}
}
