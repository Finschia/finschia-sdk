package keeper

import (
	"encoding/binary"
	"fmt"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
	"time"
)

func (k Keeper) RegisterRoleProposal(ctx sdk.Context, proposer, target sdk.AccAddress, role types.Role) (types.RoleProposal, error) {
	if k.GetRole(ctx, proposer) != types.RoleGuardian {
		return types.RoleProposal{}, sdkerrors.ErrUnauthorized.Wrap("only guardian can submit a role proposal")
	}

	if k.GetRole(ctx, target) == role {
		return types.RoleProposal{}, sdkerrors.ErrUnauthorized.Wrap("target already has same role")
	}

	proposalID := k.GetNextProposalID(ctx)
	proposal := types.RoleProposal{
		Id:        proposalID,
		Proposer:  proposer.String(),
		Target:    target.String(),
		Role:      role,
		ExpiredAt: ctx.BlockTime().Add(time.Duration(k.GetParams(ctx).ProposalPeriod)).UTC(),
	}

	k.setRoleProposal(ctx, proposal)
	k.setNextProposalID(ctx, proposalID+1)

	return proposal, nil
}

func (k Keeper) addVote(ctx sdk.Context, proposalID uint64, voter sdk.AccAddress, option types.VoteOption) error {
	if k.GetRole(ctx, voter) != types.RoleGuardian {
		return sdkerrors.ErrUnauthorized.Wrap("only guardian can vote for a role proposal")
	}

	_, found := k.GetRoleProposal(ctx, proposalID)
	if !found {
		return types.ErrUnknownProposal.Wrapf("#%d not found", proposalID)
	}

	if err := types.IsValidVoteOption(option); err != nil {
		return err
	}

	k.setVote(ctx, proposalID, voter, option)

	return nil
}

func (k Keeper) UpdateRole(ctx sdk.Context, role types.Role, addr sdk.AccAddress) error {
	previousRole := k.GetRole(ctx, addr)
	if previousRole == role {
		return sdkerrors.ErrInvalidRequest.Wrap("target already has same role")
	}

	roleMeta := k.GetRoleMetadata(ctx)
	bsMeta := k.GetBridgeStatusMetadata(ctx)

	switch previousRole {
	case types.RoleGuardian:
		roleMeta.Guardian--

		sw, err := k.GetBridgeSwitch(ctx, addr)
		if err != nil {
			panic(err)
		}

		if sw.Status == types.StatusActive {
			bsMeta.Active--
		} else {
			bsMeta.Inactive--
		}

		k.deleteBridgeSwitch(ctx, addr)

	case types.RoleOperator:
		roleMeta.Operator--
	case types.RoleJudge:
		roleMeta.Judge--
	}

	if role == types.RoleEmpty {
		k.deleteRole(ctx, addr)
		return nil
	} else {
		k.setRole(ctx, role, addr)
	}

	switch role {
	case types.RoleGuardian:
		roleMeta.Guardian++
		if err := k.setBridgeSwitch(ctx, addr, types.StatusActive); err != nil {
			panic(err)
		}
		bsMeta.Active++
	case types.RoleOperator:
		roleMeta.Operator++
	case types.RoleJudge:
		roleMeta.Judge++
	}

	k.setRoleMetadata(ctx, roleMeta)
	k.setBridgeStatusMetadata(ctx, bsMeta)

	return nil
}

func (k Keeper) setNextProposalID(ctx sdk.Context, seq uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, seq)
	store.Set(types.KeyNextProposalID, bz)
}

func (k Keeper) GetNextProposalID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyNextProposalID)
	if bz == nil {
		panic("next role proposal ID must be set at genesis")
	}

	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) setRoleProposal(ctx sdk.Context, proposal types.RoleProposal) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&proposal)
	store.Set(types.ProposalKey(proposal.Id), bz)
}

func (k Keeper) GetRoleProposal(ctx sdk.Context, id uint64) (proposal types.RoleProposal, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposalKey(id))
	if bz == nil {
		return proposal, false
	}

	k.cdc.MustUnmarshal(bz, &proposal)
	return proposal, true
}

func (k Keeper) DeleteRoleProposal(ctx sdk.Context, id uint64) {
	store := ctx.KVStore(k.storeKey)
	if _, found := k.GetRoleProposal(ctx, id); !found {
		panic(fmt.Sprintf("role proposal #%d not found", id))
	}
	store.Delete(types.ProposalKey(id))
}

// IterateProposals iterates over the all the role proposals and performs a callback function
func (k Keeper) IterateProposals(ctx sdk.Context, cb func(proposal types.RoleProposal) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyProposalPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var proposal types.RoleProposal
		k.cdc.MustUnmarshal(iterator.Value(), &proposal)
		if cb(proposal) {
			break
		}
	}
}

// GetRoleProposals returns all the role proposals from store
func (k Keeper) GetRoleProposals(ctx sdk.Context) (proposals []types.RoleProposal) {
	k.IterateProposals(ctx, func(proposal types.RoleProposal) bool {
		proposals = append(proposals, proposal)
		return false
	})
	return
}

func (k Keeper) setVote(ctx sdk.Context, proposalID uint64, voter sdk.AccAddress, option types.VoteOption) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 4)
	binary.BigEndian.PutUint32(bz, uint32(option))
	store.Set(types.VoterVoteKey(proposalID, voter), bz)
}

func (k Keeper) GetVote(ctx sdk.Context, proposalID uint64, voter sdk.AccAddress) (types.VoteOption, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VoterVoteKey(proposalID, voter))
	if bz == nil {
		return types.OptionEmpty, types.ErrUnknownVote
	}

	return types.VoteOption(binary.BigEndian.Uint32(bz)), nil
}

func (k Keeper) GetProposalVotes(ctx sdk.Context, proposalID uint64) []types.Vote {
	store := ctx.KVStore(k.storeKey)

	votes := make([]types.Vote, 0)
	iterator := sdk.KVStorePrefixIterator(store, types.VotesKey(proposalID))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		_, voter := types.SplitVoterVoteKey(iterator.Key())
		v := types.Vote{
			ProposalId: proposalID,
			Voter:      voter.String(),
			Option:     types.VoteOption(binary.BigEndian.Uint32(iterator.Value())),
		}
		votes = append(votes, v)
	}

	return votes
}

func (k Keeper) setRole(ctx sdk.Context, role types.Role, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 4)
	binary.BigEndian.PutUint32(bz, uint32(role))
	store.Set(types.RoleKey(addr), bz)
}

func (k Keeper) GetRole(ctx sdk.Context, addr sdk.AccAddress) types.Role {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RoleKey(addr))
	if bz == nil {
		return types.RoleEmpty
	}

	return types.Role(binary.BigEndian.Uint32(bz))
}

func (k Keeper) GetRolePairs(ctx sdk.Context) []types.RolePair {
	store := ctx.KVStore(k.storeKey)
	pairs := make([]types.RolePair, 0)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyRolePrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		assignee := types.SplitRoleKey(iterator.Key())
		pairs = append(pairs, types.RolePair{Address: assignee.String(), Role: types.Role(binary.BigEndian.Uint32(iterator.Value()))})
	}

	return pairs
}

func (k Keeper) deleteRole(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.RoleKey(addr))
}

func (k Keeper) setBridgeSwitch(ctx sdk.Context, guardian sdk.AccAddress, status types.BridgeStatus) error {
	if k.GetRole(ctx, guardian) != types.RoleGuardian {
		return sdkerrors.ErrUnauthorized.Wrap("only guardian can set bridge switch")
	}

	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 4)
	binary.BigEndian.PutUint32(bz, uint32(status))
	store.Set(types.BridgeSwitchKey(guardian), bz)

	return nil
}

func (k Keeper) deleteBridgeSwitch(ctx sdk.Context, guardian sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.BridgeSwitchKey(guardian))
}

func (k Keeper) GetBridgeSwitch(ctx sdk.Context, guardian sdk.AccAddress) (types.BridgeSwitch, error) {
	if k.GetRole(ctx, guardian) != types.RoleGuardian {
		return types.BridgeSwitch{}, sdkerrors.ErrUnauthorized.Wrap("only guardian can set bridge switch")
	}

	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.BridgeSwitchKey(guardian))
	if bz == nil {
		panic("bridge switch must be set at genesis")
	}

	return types.BridgeSwitch{Guardian: guardian.String(), Status: types.BridgeStatus(binary.BigEndian.Uint32(bz))}, nil
}

func (k Keeper) GetBridgeSwitches(ctx sdk.Context) []types.BridgeSwitch {
	store := ctx.KVStore(k.storeKey)

	bws := make([]types.BridgeSwitch, 0)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyBridgeSwitchPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		addr := types.SplitBridgeSwitchKey(iterator.Key())
		bws = append(bws, types.BridgeSwitch{Guardian: addr.String(), Status: types.BridgeStatus(binary.BigEndian.Uint32(iterator.Value()))})
	}

	return bws
}
