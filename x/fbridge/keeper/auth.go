package keeper

import (
	"encoding/binary"
	"errors"
	"fmt"
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
	"time"
)

func (k Keeper) RegisterRoleProposal(ctx sdk.Context, proposer, target sdk.AccAddress, role types.Role) (types.RoleProposal, error) {
	if k.GetRoleMetadata(ctx).Guardian < 1 && k.authority != proposer.String() {
		return types.RoleProposal{}, sdkerrors.ErrUnauthorized.Wrapf("only %s can submit a role proposal if there are no guardians", k.authority)
	} else if k.GetRole(ctx, proposer) != types.RoleGuardian {
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

func (k Keeper) IsValidRole(role types.Role) error {
	switch role {
	case types.RoleGuardian, types.RoleOperator, types.RoleJudge:
		return nil
	}

	return errors.New("unsupported role")
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

// GetProposals returns all the role proposals from store
func (k Keeper) GetProposals(ctx sdk.Context) (proposals []types.RoleProposal) {
	k.IterateProposals(ctx, func(proposal types.RoleProposal) bool {
		proposals = append(proposals, proposal)
		return false
	})
	return
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

func (k Keeper) setRoleMetadata(ctx sdk.Context, data types.RoleMetadata) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&data)
	store.Set(types.KeyRoleMetadata, bz)
}

func (k Keeper) GetRoleMetadata(ctx sdk.Context) types.RoleMetadata {
	store := ctx.KVStore(k.storeKey)

	data := types.RoleMetadata{}
	bz := store.Get(types.KeyRoleMetadata)
	if bz == nil {
		panic("role metadata must be set at genesis")
	}
	k.cdc.MustUnmarshal(bz, &data)
	return data
}
