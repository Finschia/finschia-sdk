package internal

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (k Keeper) Vote(ctx sdk.Context, vote foundation.Vote) error {
	if err := validateMetadata(vote.Metadata, k.config); err != nil {
		return err
	}

	// Make sure that a voter hasn't already voted.
	voter, err := k.addressCodec().StringToBytes(vote.Voter)
	if err != nil {
		panic(err)
	}

	if k.hasVote(ctx, vote.ProposalId, voter) {
		return sdkerrors.ErrInvalidRequest.Wrapf("Already voted: %s", vote.Voter)
	}

	proposal, err := k.GetProposal(ctx, vote.ProposalId)
	if err != nil {
		return err
	}

	// Ensure that we can still accept votes for this proposal.
	if proposal.Status != foundation.PROPOSAL_STATUS_SUBMITTED {
		return sdkerrors.ErrInvalidRequest.Wrapf("not possible with proposal status: %s", proposal.Status)
	}
	if !ctx.BlockTime().Before(proposal.VotingPeriodEnd) {
		return sdkerrors.ErrInvalidRequest.Wrap("voting period has ended already")
	}

	vote.SubmitTime = ctx.BlockTime()
	k.setVote(ctx, vote)

	if err := ctx.EventManager().EmitTypedEvent(&foundation.EventVote{
		Vote: vote,
	}); err != nil {
		panic(err)
	}

	return nil
}

func (k Keeper) hasVote(ctx sdk.Context, proposalID uint64, voter sdk.AccAddress) bool {
	store := k.storeService.OpenKVStore(ctx)
	key := voteKey(proposalID, voter)
	has, err := store.Has(key)
	return (err == nil) && has
}

func (k Keeper) GetVote(ctx sdk.Context, proposalID uint64, voter sdk.AccAddress) (*foundation.Vote, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := voteKey(proposalID, voter)
	bz, err := store.Get(key)
	if err != nil {
		return nil, err
	}
	if len(bz) == 0 {
		return nil, sdkerrors.ErrNotFound.Wrapf("No vote for proposal %d: %s", proposalID, voter)
	}

	var vote foundation.Vote
	k.cdc.MustUnmarshal(bz, &vote)

	return &vote, nil
}

func (k Keeper) setVote(ctx sdk.Context, vote foundation.Vote) {
	store := k.storeService.OpenKVStore(ctx)
	voter, err := k.addressCodec().StringToBytes(vote.Voter)
	if err != nil {
		panic(err)
	}

	key := voteKey(vote.ProposalId, voter)
	bz := k.cdc.MustMarshal(&vote)
	store.Set(key, bz)
}

func (k Keeper) iterateVotes(ctx sdk.Context, proposalID uint64, fn func(vote foundation.Vote) (stop bool)) {
	store := k.storeService.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(store)
	prefix := append(voteKeyPrefix, Uint64ToBytes(proposalID)...)
	iterator := storetypes.KVStorePrefixIterator(adapter, prefix)
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
		voter, err := k.addressCodec().StringToBytes(vote.Voter)
		if err != nil {
			panic(err)
		}

		keys = append(keys, voteKey(proposalID, voter))
		return false
	})

	store := k.storeService.OpenKVStore(ctx)
	for _, key := range keys {
		store.Delete(key)
	}
}
