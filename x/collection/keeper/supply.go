package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/collection"
)

func (k Keeper) CreateContract(ctx sdk.Context, creator sdk.AccAddress, contract collection.Contract) (*string, error) {
	contractID, err := k.createContract(ctx, contract)
	if err != nil {
		return nil, err
	}

	for permission := range collection.Permission_name {
		p := collection.Permission(permission)
		if p == collection.Permission_Unspecified {
			continue
		}
		// TODO: revisit grant
		k.setGrant(ctx, *contractID, creator, collection.Permission(permission))
	}

	return contractID, nil
}

func (k Keeper) createContract(ctx sdk.Context, contract collection.Contract) (*string, error) {
	contractID := k.classKeeper.NewID(ctx)
	contract.ContractId = contractID
	k.setContract(ctx, contract)

	// set the next class ids
	nextIDs := collection.DefaultNextClassIDs(contractID)
	k.setNextClassIDs(ctx, nextIDs)

	// TODO: emit event
	return &contractID, nil
}

func (k Keeper) GetContract(ctx sdk.Context, contractID string) (*collection.Contract, error) {
	store := ctx.KVStore(k.storeKey)
	key := contractKey(contractID)
	bz := store.Get(key)
	if bz == nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("No such a contract: %s", contractID)
	}

	var contract collection.Contract
	if err := contract.Unmarshal(bz); err != nil {
		panic(err)
	}
	return &contract, nil
}

func (k Keeper) setContract(ctx sdk.Context, contract collection.Contract) {
	store := ctx.KVStore(k.storeKey)
	key := contractKey(contract.ContractId)

	bz, err := contract.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}

func (k Keeper) CreateTokenClass(ctx sdk.Context, class collection.TokenClass) (*string, error) {
	contractID := class.GetContractId()
	if _, err := k.GetContract(ctx, contractID); err != nil {
		return nil, err
	}

	nextClassIDs := k.getNextClassIDs(ctx, contractID)
	class.SetId(&nextClassIDs)
	k.setNextClassIDs(ctx, nextClassIDs)

	if err := class.ValidateBasic(); err != nil {
		return nil, err
	}
	k.setTokenClass(ctx, class)

	// TODO: emit event
	id := class.GetId()
	return &id, nil
}

func (k Keeper) GetTokenClass(ctx sdk.Context, contractID, classID string) (collection.TokenClass, error) {
	store := ctx.KVStore(k.storeKey)
	key := classKey(contractID, classID)
	bz := store.Get(key)
	if bz == nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("No such a class in contract %s: %s", contractID, classID)
	}

	var class collection.TokenClass
	if err := k.cdc.UnmarshalInterface(bz, &class); err != nil {
		panic(err)
	}
	return class, nil
}

func (k Keeper) setTokenClass(ctx sdk.Context, class collection.TokenClass) {
	store := ctx.KVStore(k.storeKey)
	key := classKey(class.GetContractId(), class.GetId())

	bz, err := k.cdc.MarshalInterface(class)
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}

func (k Keeper) getNextClassIDs(ctx sdk.Context, contractID string) collection.NextClassIDs {
	store := ctx.KVStore(k.storeKey)
	key := nextClassIDKey(contractID)
	bz := store.Get(key)
	if bz == nil {
		panic(sdkerrors.ErrNotFound.Wrapf("No next ids of contract %s", contractID))
	}

	var class collection.NextClassIDs
	if err := class.Unmarshal(bz); err != nil {
		panic(err)
	}
	return class
}

func (k Keeper) setNextClassIDs(ctx sdk.Context, ids collection.NextClassIDs) {
	store := ctx.KVStore(k.storeKey)
	key := nextClassIDKey(ids.ContractId)

	bz, err := ids.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}

func (k Keeper) Grant(ctx sdk.Context, contractID string, granter, grantee sdk.AccAddress, permission collection.Permission) error {
	k.grant(ctx, contractID, grantee, permission)

	event := collection.EventGrant{
		ContractId: contractID,
		Granter:    granter.String(),
		Grantee:    grantee.String(),
		Permission: permission.String(),
	}
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) grant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission collection.Permission) {
	k.setGrant(ctx, contractID, grantee, permission)

	if !k.accountKeeper.HasAccount(ctx, grantee) {
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, grantee))
	}
}

func (k Keeper) Abandon(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission collection.Permission) error {
	k.deleteGrant(ctx, contractID, grantee, permission)

	event := collection.EventAbandon{
		ContractId: contractID,
		Grantee:    grantee.String(),
		Permission: permission.String(),
	}
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) GetGrant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission collection.Permission) (*collection.Grant, error) {
	store := ctx.KVStore(k.storeKey)
	if store.Has(grantKey(contractID, grantee, permission)) {
		return &collection.Grant{
			Grantee:    grantee.String(),
			Permission: permission.String(),
		}, nil
	}
	return nil, sdkerrors.ErrNotFound.Wrapf("no %s permission granted on %s", permission, grantee)
}

func (k Keeper) setGrant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission collection.Permission) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(contractID, grantee, permission)
	store.Set(key, []byte{})
}

func (k Keeper) deleteGrant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission collection.Permission) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(contractID, grantee, permission)
	store.Delete(key)
}
