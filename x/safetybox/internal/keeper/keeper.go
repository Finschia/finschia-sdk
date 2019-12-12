package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cbank "github.com/cosmos/cosmos-sdk/x/bank"
	iam "github.com/line/link/x/iam/exported"
	"github.com/line/link/x/safetybox/internal/types"
	"github.com/tendermint/tendermint/crypto"
)

type Keeper struct {
	cdc           *codec.Codec
	storeKey      sdk.StoreKey
	iamKeeper     iam.IamKeeper
	bankKeeper    cbank.Keeper
	accountKeeper auth.AccountKeeper
}

func NewKeeper(cdc *codec.Codec, iamKeeper iam.IamKeeper, bankKeeper cbank.Keeper, accountKeeper auth.AccountKeeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		iamKeeper:     iamKeeper,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
	}
}

func (k Keeper) NewSafetyBox(ctx sdk.Context, msg types.MsgSafetyBoxCreate) (types.SafetyBox, sdk.Error) {
	// create new safety box account
	newSafetyBoxAccount, err := k.newSafetyBoxAccount(ctx, msg.SafetyBoxId)
	if err != nil {
		return types.SafetyBox{}, err
	}

	// create new safety box
	sb := types.NewSafetyBox(msg.SafetyBoxOwner, msg.SafetyBoxId, newSafetyBoxAccount)

	// reject if the safety box id exists
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.SafetyBoxKey(sb.ID)) {
		return types.SafetyBox{}, types.ErrSafetyBoxIdExist(types.DefaultCodespace)
	}
	store.Set(types.SafetyBoxKey(sb.ID), k.cdc.MustMarshalBinaryBare(sb))

	// grant the owner a permission to whitelist operators
	k.iamKeeper.GrantPermission(ctx, sb.Owner, types.NewWhitelistOperatorsPermission(sb.ID))

	return sb, nil
}

func (k Keeper) newSafetyBoxAccount(ctx sdk.Context, safetyBoxId string) (sdk.AccAddress, sdk.Error) {
	// hash safety box id
	newAddress := sdk.AccAddress(crypto.AddressHash(types.SafetyBoxKey(safetyBoxId)))

	// check if exist
	acc := k.accountKeeper.GetAccount(ctx, newAddress)
	if acc != nil {
		return nil, types.ErrSafetyBoxAccountExist(types.DefaultCodespace)
	}

	// create new account and return its address
	newAccount := k.accountKeeper.NewAccountWithAddress(ctx, newAddress)
	k.accountKeeper.SetAccount(ctx, newAccount)

	return newAccount.GetAddress(), nil
}

func (k Keeper) GetSafetyBox(ctx sdk.Context, safetyBoxId string) (types.SafetyBox, sdk.Error) {
	sb, err := k.get(ctx, safetyBoxId)
	if err != nil {
		return types.SafetyBox{}, err
	}
	return sb, nil
}

func (k Keeper) Allocate(ctx sdk.Context, msg types.MsgSafetyBoxAllocateCoins) sdk.Error {
	sb, err := k.get(ctx, msg.SafetyBoxId)
	if err != nil {
		return err
	}

	// from the allocator, to the safety box
	fromAddress := msg.AllocatorAddress
	toAddress := sb.Address

	// only allocator could allocate
	allocatePermission := types.NewAllocatePermission(msg.SafetyBoxId)
	if !k.iamKeeper.HasPermission(ctx, fromAddress, allocatePermission) {
		return types.ErrSafetyBoxPermissionAllocate(types.DefaultCodespace)
	}

	// allocation
	err = k.sendCoins(ctx, fromAddress, toAddress, msg.Coins)
	if err != nil {
		return err
	}

	// increase the total allocation and cumulative allocation
	sb.TotalAllocation = sb.TotalAllocation.Add(msg.Coins)
	sb.CumulativeAllocation = sb.CumulativeAllocation.Add(msg.Coins)

	return k.set(ctx, msg.SafetyBoxId, sb)
}

func (k Keeper) Recall(ctx sdk.Context, msg types.MsgSafetyBoxRecallCoins) sdk.Error {
	sb, err := k.get(ctx, msg.SafetyBoxId)
	if err != nil {
		return err
	}

	// from the safety box, to the allocator
	fromAddress := sb.Address
	toAddress := msg.AllocatorAddress

	// only allocator could recall
	recallPermission := types.NewRecallPermission(msg.SafetyBoxId)
	if !k.iamKeeper.HasPermission(ctx, toAddress, recallPermission) {
		return types.ErrSafetyBoxPermissionRecall(types.DefaultCodespace)
	}

	// check not to recall more than allocated
	if msg.Coins.IsAnyGT(sb.TotalAllocation) {
		return types.ErrSafetyBoxRecallMoreThanAllocated(types.DefaultCodespace)
	}

	// recall
	err = k.sendCoins(ctx, fromAddress, toAddress, msg.Coins)
	if err != nil {
		return err
	}

	// decrease the total allocation
	sb.TotalAllocation = sb.TotalAllocation.Sub(msg.Coins)

	return k.set(ctx, msg.SafetyBoxId, sb)
}

func (k Keeper) Issue(ctx sdk.Context, msg types.MsgSafetyBoxIssueCoins) sdk.Error {
	sb, err := k.get(ctx, msg.SafetyBoxId)
	if err != nil {
		return err
	}

	// both issuer and issuee must be issuers
	issuerAddress := msg.FromAddress
	toAddress := msg.ToAddress

	issuePermission := types.NewIssuePermission(msg.SafetyBoxId)
	if !k.iamKeeper.HasPermission(ctx, issuerAddress, issuePermission) {
		return types.ErrSafetyBoxPermissionIssue(types.DefaultCodespace)
	}
	if !k.iamKeeper.HasPermission(ctx, toAddress, issuePermission) {
		return types.ErrSafetyBoxPermissionIssue(types.DefaultCodespace)
	}

	// issue from the safety box to an issuer
	fromAddress := sb.Address
	err = k.sendCoins(ctx, fromAddress, toAddress, msg.Coins)
	if err != nil {
		return err
	}

	// increase the total issuance
	sb.TotalIssuance = sb.TotalIssuance.Add(msg.Coins)

	return k.set(ctx, msg.SafetyBoxId, sb)
}

func (k Keeper) Return(ctx sdk.Context, msg types.MsgSafetyBoxReturnCoins) sdk.Error {
	sb, err := k.get(ctx, msg.SafetyBoxId)
	if err != nil {
		return err
	}

	// from the returner, to the safety box
	fromAddress := msg.ReturnerAddress
	toAddress := sb.Address

	// only returner could return
	returnPermission := types.NewReturnPermission(msg.SafetyBoxId)
	if !k.iamKeeper.HasPermission(ctx, fromAddress, returnPermission) {
		return types.ErrSafetyBoxPermissionReturn(types.DefaultCodespace)
	}

	// check not to return more than issued
	if msg.Coins.IsAnyGT(sb.TotalIssuance) {
		return types.ErrSafetyBoxReturnMoreThanIssued(types.DefaultCodespace)
	}

	// return
	err = k.sendCoins(ctx, fromAddress, toAddress, msg.Coins)
	if err != nil {
		return err
	}

	// decrease the total issuance
	sb.TotalIssuance = sb.TotalIssuance.Sub(msg.Coins)

	return k.set(ctx, msg.SafetyBoxId, sb)
}

func (k Keeper) sendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, coins sdk.Coins) sdk.Error {
	err := k.bankKeeper.SendCoins(ctx, fromAddr, toAddr, coins)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) GrantPermission(ctx sdk.Context, safetyBoxId string, by sdk.AccAddress, acc sdk.AccAddress, action string) sdk.Error {
	// reject self-grant
	if by.Equals(acc) {
		return types.ErrSafetyBoxSelfPermission(types.DefaultCodespace)
	}

	// grant operator
	if action == types.RoleOperator {
		// check whitelist permission
		whitelistOperatorsPermission := types.NewWhitelistOperatorsPermission(safetyBoxId)
		if !k.iamKeeper.HasPermission(ctx, by, whitelistOperatorsPermission) {
			return types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace)
		}

		if k.IsOperator(ctx, safetyBoxId, acc) {
			return types.ErrSafetyBoxHasPermissionAlready(types.DefaultCodespace)
		}
		if k.IsAllocator(ctx, safetyBoxId, acc) || k.IsIssuer(ctx, safetyBoxId, acc) || k.IsReturner(ctx, safetyBoxId, acc) {
			return types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace)
		}
		k.iamKeeper.GrantPermission(ctx, acc, types.NewWhitelistOtherRolesPermission(safetyBoxId))
	}

	// grant issuer
	if action == types.RoleIssuer {
		// check whitelist permission
		whitelistOtherRolesPermission := types.NewWhitelistOtherRolesPermission(safetyBoxId)
		if !k.iamKeeper.HasPermission(ctx, by, whitelistOtherRolesPermission) {
			return types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace)
		}

		if k.IsIssuer(ctx, safetyBoxId, acc) {
			return types.ErrSafetyBoxHasPermissionAlready(types.DefaultCodespace)
		}
		if k.IsOperator(ctx, safetyBoxId, acc) || k.IsReturner(ctx, safetyBoxId, acc) || k.IsAllocator(ctx, safetyBoxId, acc) {
			return types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace)
		}
		k.iamKeeper.GrantPermission(ctx, acc, types.NewIssuePermission(safetyBoxId))
	}

	// grant returner
	if action == types.RoleReturner {
		// check whitelist permission
		whitelistOtherRolesPermission := types.NewWhitelistOtherRolesPermission(safetyBoxId)
		if !k.iamKeeper.HasPermission(ctx, by, whitelistOtherRolesPermission) {
			return types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace)
		}

		if k.IsReturner(ctx, safetyBoxId, acc) {
			return types.ErrSafetyBoxHasPermissionAlready(types.DefaultCodespace)
		}
		if k.IsOperator(ctx, safetyBoxId, acc) || k.IsIssuer(ctx, safetyBoxId, acc) || k.IsAllocator(ctx, safetyBoxId, acc) {
			return types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace)
		}
		k.iamKeeper.GrantPermission(ctx, acc, types.NewReturnPermission(safetyBoxId))
	}

	// grant allocator
	if action == types.RoleAllocator {
		// check whitelist permission
		whitelistOtherRolesPermission := types.NewWhitelistOtherRolesPermission(safetyBoxId)
		if !k.iamKeeper.HasPermission(ctx, by, whitelistOtherRolesPermission) {
			return types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace)
		}

		if k.IsAllocator(ctx, safetyBoxId, acc) {
			return types.ErrSafetyBoxHasPermissionAlready(types.DefaultCodespace)
		}
		if k.IsOperator(ctx, safetyBoxId, acc) || k.IsIssuer(ctx, safetyBoxId, acc) || k.IsReturner(ctx, safetyBoxId, acc) {
			return types.ErrSafetyBoxHasOtherPermission(types.DefaultCodespace)
		}
		k.iamKeeper.GrantPermission(ctx, acc, types.NewAllocatePermission(safetyBoxId))
		k.iamKeeper.GrantPermission(ctx, acc, types.NewRecallPermission(safetyBoxId))
	}

	return nil
}

func (k Keeper) RevokePermission(ctx sdk.Context, safetyBoxId string, by sdk.AccAddress, acc sdk.AccAddress, action string) sdk.Error {
	// reject self-revoke
	if by.Equals(acc) {
		return types.ErrSafetyBoxSelfPermission(types.DefaultCodespace)
	}

	// revoke operator
	if action == types.RoleOperator {
		// check whitelist permission
		whitelistOperatorsPermission := types.NewWhitelistOperatorsPermission(safetyBoxId)
		if !k.iamKeeper.HasPermission(ctx, by, whitelistOperatorsPermission) {
			return types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace)
		}

		if !k.IsOperator(ctx, safetyBoxId, acc) {
			return types.ErrSafetyBoxDoesNotHavePermission(types.DefaultCodespace)
		}
		k.iamKeeper.RevokePermission(ctx, acc, types.NewWhitelistOtherRolesPermission(safetyBoxId))
	}

	// revoke issuer
	if action == types.RoleIssuer {
		// check whitelist permission
		whitelistOtherRolesPermission := types.NewWhitelistOtherRolesPermission(safetyBoxId)
		if !k.iamKeeper.HasPermission(ctx, by, whitelistOtherRolesPermission) {
			return types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace)
		}

		if !k.IsIssuer(ctx, safetyBoxId, acc) {
			return types.ErrSafetyBoxDoesNotHavePermission(types.DefaultCodespace)
		}
		k.iamKeeper.RevokePermission(ctx, acc, types.NewIssuePermission(safetyBoxId))
	}

	// revoke returner
	if action == types.RoleReturner {
		// check whitelist permission
		whitelistOtherRolesPermission := types.NewWhitelistOtherRolesPermission(safetyBoxId)
		if !k.iamKeeper.HasPermission(ctx, by, whitelistOtherRolesPermission) {
			return types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace)
		}

		if !k.IsReturner(ctx, safetyBoxId, acc) {
			return types.ErrSafetyBoxDoesNotHavePermission(types.DefaultCodespace)
		}
		k.iamKeeper.RevokePermission(ctx, acc, types.NewReturnPermission(safetyBoxId))
	}

	// revoke allocator
	if action == types.RoleAllocator {
		// check whitelist permission
		whitelistOtherRolesPermission := types.NewWhitelistOtherRolesPermission(safetyBoxId)
		if !k.iamKeeper.HasPermission(ctx, by, whitelistOtherRolesPermission) {
			return types.ErrSafetyBoxPermissionWhitelist(types.DefaultCodespace)
		}

		if !k.IsAllocator(ctx, safetyBoxId, acc) {
			return types.ErrSafetyBoxDoesNotHavePermission(types.DefaultCodespace)
		}
		k.iamKeeper.RevokePermission(ctx, acc, types.NewAllocatePermission(safetyBoxId))
		k.iamKeeper.RevokePermission(ctx, acc, types.NewRecallPermission(safetyBoxId))
	}

	return nil
}

func (k Keeper) get(ctx sdk.Context, safetyBoxId string) (types.SafetyBox, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	// retrieve the safety box
	bz := store.Get(types.SafetyBoxKey(safetyBoxId))
	if bz == nil {
		return types.SafetyBox{}, types.ErrSafetyBoxNotExist(types.DefaultCodespace)
	}
	r := &types.SafetyBox{}
	k.cdc.MustUnmarshalBinaryBare(bz, r)
	return *r, nil
}

func (k Keeper) set(ctx sdk.Context, safetyBoxId string, sb types.SafetyBox) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.SafetyBoxKey(safetyBoxId), k.cdc.MustMarshalBinaryBare(sb))
	return nil
}

func (k Keeper) GetPermissions(ctx sdk.Context, safetyBoxId, role string, acc sdk.AccAddress) (types.MsgSafetyBoxRoleResponse, sdk.Error) {
	var hasRole bool
	switch role {
	case types.RoleOwner:
		hasRole = k.IsOwner(ctx, safetyBoxId, acc)
	case types.RoleOperator:
		hasRole = k.IsOperator(ctx, safetyBoxId, acc)
	case types.RoleAllocator:
		hasRole = k.IsAllocator(ctx, safetyBoxId, acc)
	case types.RoleIssuer:
		hasRole = k.IsIssuer(ctx, safetyBoxId, acc)
	case types.RoleReturner:
		hasRole = k.IsReturner(ctx, safetyBoxId, acc)
	default:
		return types.MsgSafetyBoxRoleResponse{HasRole: false}, types.ErrSafetyBoxInvalidRole(types.DefaultCodespace)
	}
	return types.MsgSafetyBoxRoleResponse{HasRole: hasRole}, nil
}

func (k Keeper) IsOwner(ctx sdk.Context, safetyBoxId string, acc sdk.AccAddress) bool {
	sb, err := k.get(ctx, safetyBoxId)
	if err != nil {
		return false
	}

	return sb.Owner.Equals(acc)
}

func (k Keeper) IsOperator(ctx sdk.Context, safetyBoxId string, acc sdk.AccAddress) bool {
	return k.iamKeeper.HasPermission(ctx, acc, types.NewWhitelistOtherRolesPermission(safetyBoxId))
}

func (k Keeper) IsAllocator(ctx sdk.Context, safetyBoxId string, acc sdk.AccAddress) bool {
	canAllocate := k.iamKeeper.HasPermission(ctx, acc, types.NewAllocatePermission(safetyBoxId))
	canRecall := k.iamKeeper.HasPermission(ctx, acc, types.NewRecallPermission(safetyBoxId))
	return canAllocate && canRecall
}

func (k Keeper) IsIssuer(ctx sdk.Context, safetyBoxId string, acc sdk.AccAddress) bool {
	return k.iamKeeper.HasPermission(ctx, acc, types.NewIssuePermission(safetyBoxId))
}

func (k Keeper) IsReturner(ctx sdk.Context, safetyBoxId string, acc sdk.AccAddress) bool {
	return k.iamKeeper.HasPermission(ctx, acc, types.NewReturnPermission(safetyBoxId))
}
