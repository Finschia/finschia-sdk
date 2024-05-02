package keeper

import (
	"encoding/binary"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, gs *types.GenesisState) error {
	k.SetParams(ctx, gs.Params)
	k.setNextSequence(ctx, gs.SendingState.NextSeq)
	k.setNextProposalID(ctx, gs.NextRoleProposalId)
	k.setRoleMetadata(ctx, gs.RoleMetadata)

	for _, proposal := range gs.RoleProposals {
		k.setRoleProposal(ctx, proposal)
	}
	// TODO: we initialize the appropriate genesis parameters whenever the feature is added

	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
		SendingState: types.SendingState{
			NextSeq: k.GetNextSequence(ctx),
		},
		NextRoleProposalId: k.GetNextProposalID(ctx),
		RoleMetadata:       k.GetRoleMetadata(ctx),
		RoleProposals:      k.GetProposals(ctx),
		Votes:              k.GetAllVotes(ctx),
	}
}

// IterateVotes iterates over the all the votes for role proposals and performs a callback function
func (k Keeper) IterateVotes(ctx sdk.Context, cb func(proposal types.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyProposalPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		id, voter := types.SplitVoterVoteKey(iterator.Key())
		opt := types.VoteOption(binary.BigEndian.Uint32(iterator.Value()))
		v := types.Vote{ProposalId: id, Voter: voter.String(), Option: opt}
		k.cdc.MustUnmarshal(iterator.Value(), &v)
		if cb(v) {
			break
		}
	}
}

// GetAllVotes returns all the votes from the store
func (k Keeper) GetAllVotes(ctx sdk.Context) (votes []types.Vote) {
	k.IterateVotes(ctx, func(vote types.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}
