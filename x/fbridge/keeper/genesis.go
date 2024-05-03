package keeper

import (
	"encoding/binary"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, gs *types.GenesisState) error {

	k.SetParams(ctx, gs.Params)
	k.setNextSequence(ctx, gs.SendingState.NextSeq)
	for _, info := range gs.SendingState.SeqToBlocknum {
		k.setSeqToBlocknum(ctx, info.Seq, info.Blocknum)
	}

	for _, pair := range gs.Roles {
		k.setRole(ctx, pair.Role, sdk.MustAccAddressFromBech32(pair.Address))
	}

	for _, sw := range gs.BridgeSwitches {
		if err := k.setBridgeSwitch(ctx, sdk.MustAccAddressFromBech32(sw.Guardian), sw.Status); err != nil {
			panic(err)
		}
	}

	k.setNextProposalID(ctx, gs.NextRoleProposalId)
	for _, proposal := range gs.RoleProposals {
		k.setRoleProposal(ctx, proposal)
	}

	for _, vote := range gs.Votes {
		k.setVote(ctx, vote.ProposalId, sdk.MustAccAddressFromBech32(vote.Voter), vote.Option)
	}

	// TODO: we initialize the appropriate genesis parameters whenever the feature is added

	k.InitMemStore(ctx)

	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
		SendingState: types.SendingState{
			NextSeq:       k.GetNextSequence(ctx),
			SeqToBlocknum: k.getAllSeqToBlocknums(ctx),
		},
		NextRoleProposalId: k.GetNextProposalID(ctx),
		RoleProposals:      k.GetRoleProposals(ctx),
		Votes:              k.GetAllVotes(ctx),
		Roles:              k.GetRolePairs(ctx),
	}
}

func (k Keeper) getAllSeqToBlocknums(ctx sdk.Context) []types.BlockSeqInfo {
	infos := make([]types.BlockSeqInfo, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeySeqToBlocknumPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		seq := binary.BigEndian.Uint64(iterator.Key()[1:])
		v := binary.BigEndian.Uint64(iterator.Value())
		info := types.BlockSeqInfo{Seq: seq, Blocknum: v}
		infos = append(infos, info)
	}

	return infos
}

// IterateVotes iterates over the all the votes for role proposals and performs a callback function
func (k Keeper) IterateVotes(ctx sdk.Context, cb func(proposal types.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyProposalVotePrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		id, voter := types.SplitVoterVoteKey(iterator.Key())
		opt := types.VoteOption(binary.BigEndian.Uint32(iterator.Value()))
		v := types.Vote{ProposalId: id, Voter: voter.String(), Option: opt}
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
