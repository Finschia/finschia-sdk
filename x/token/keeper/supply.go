package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token"
)

func (k Keeper) Issue(ctx sdk.Context, class token.TokenClass, owner, to sdk.AccAddress, amount sdk.Int) {
	k.issue(ctx, class)

	event := token.EventIssued{
		Creator:    owner.String(),
		ContractId: class.ContractId,
		Name:       class.Name,
		Symbol:     class.Symbol,
		Uri:        class.ImageUri,
		Meta:       class.Meta,
		Decimals:   class.Decimals,
		Mintable:   class.Mintable,
	}
	ctx.EventManager().EmitEvent(token.NewEventIssueToken(event, to, amount)) // deprecated
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	permissions := []token.Permission{
		token.PermissionModify,
	}
	if class.Mintable {
		permissions = append(permissions,
			token.PermissionMint,
			token.PermissionBurn,
		)
	}

	// legacy
	eventGrant := token.EventGranted{
		ContractId: class.ContractId,
		Grantee:    to.String(),
	}
	ctx.EventManager().EmitEvent(token.NewEventGrantPermTokenHead(eventGrant))
	for _, permission := range permissions {
		eventGrant.Permission = permission
		ctx.EventManager().EmitEvent(token.NewEventGrantPermTokenBody(eventGrant))
		k.Grant(ctx, class.ContractId, nil, owner, permission)
	}

	k.mintToken(ctx, class.ContractId, to, amount)

	if err := ctx.EventManager().EmitTypedEvent(&token.EventMinted{
		ContractId: class.ContractId,
		Operator:   owner.String(),
		To:         to.String(),
		Amount:     amount,
	}); err != nil {
		panic(err)
	}
}

func (k Keeper) issue(ctx sdk.Context, class token.TokenClass) {
	if _, err := k.GetClass(ctx, class.ContractId); err == nil {
		panic(sdkerrors.ErrInvalidRequest.Wrapf("ID already exists: %s", class.ContractId))
	}
	k.setClass(ctx, class)
}

func (k Keeper) GetClass(ctx sdk.Context, contractID string) (*token.TokenClass, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(classKey(contractID))
	if bz == nil {
		return nil, token.ErrContractNotFound.Wrap(contractID)
	}

	var class token.TokenClass
	if err := k.cdc.Unmarshal(bz, &class); err != nil {
		panic(err)
	}

	return &class, nil
}

func (k Keeper) setClass(ctx sdk.Context, class token.TokenClass) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&class)
	if err != nil {
		panic(err)
	}

	store.Set(classKey(class.ContractId), bz)
}

func (k Keeper) Mint(ctx sdk.Context, contractID string, grantee, to sdk.AccAddress, amount sdk.Int) error {
	if err := k.mint(ctx, contractID, grantee, to, amount); err != nil {
		return err
	}

	event := token.EventMinted{
		ContractId: contractID,
		Operator:   grantee.String(),
		To:         to.String(),
		Amount:     amount,
	}
	ctx.EventManager().EmitEvent(token.NewEventMintToken(event)) // deprecated
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}
	return nil
}

func (k Keeper) mint(ctx sdk.Context, contractID string, grantee, to sdk.AccAddress, amount sdk.Int) error {
	if _, err := k.GetGrant(ctx, contractID, grantee, token.PermissionMint); err != nil {
		return sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	k.mintToken(ctx, contractID, to, amount)

	return nil
}

func (k Keeper) mintToken(ctx sdk.Context, contractID string, addr sdk.AccAddress, amount sdk.Int) {
	k.addToken(ctx, contractID, addr, amount)

	minted := k.GetMinted(ctx, contractID)
	minted = minted.Add(amount)
	k.setMinted(ctx, contractID, minted)

	supply := k.GetSupply(ctx, contractID)
	supply = supply.Add(amount)
	k.setSupply(ctx, contractID, supply)
}

func (k Keeper) Burn(ctx sdk.Context, contractID string, from sdk.AccAddress, amount sdk.Int) error {
	if err := k.burn(ctx, contractID, from, amount); err != nil {
		return err
	}

	event := token.EventBurned{
		ContractId: contractID,
		Operator:   from.String(),
		From:       from.String(),
		Amount:     amount,
	}
	ctx.EventManager().EmitEvent(token.NewEventBurnToken(event)) // deprecated
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}
	return nil
}

func (k Keeper) burn(ctx sdk.Context, contractID string, from sdk.AccAddress, amount sdk.Int) error {
	if _, err := k.GetGrant(ctx, contractID, from, token.PermissionBurn); err != nil {
		return sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if err := k.burnToken(ctx, contractID, from, amount); err != nil {
		return err
	}

	return nil
}

func (k Keeper) OperatorBurn(ctx sdk.Context, contractID string, operator, from sdk.AccAddress, amount sdk.Int) error {
	if err := k.operatorBurn(ctx, contractID, operator, from, amount); err != nil {
		return err
	}

	event := token.EventBurned{
		ContractId: contractID,
		Operator:   operator.String(),
		From:       from.String(),
		Amount:     amount,
	}
	ctx.EventManager().EmitEvent(token.NewEventBurnTokenFrom(event)) // deprecated
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}
	return nil
}

func (k Keeper) operatorBurn(ctx sdk.Context, contractID string, operator, from sdk.AccAddress, amount sdk.Int) error {
	_, err := k.GetGrant(ctx, contractID, operator, token.PermissionBurn)
	if err != nil {
		return sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}
	if _, err := k.GetAuthorization(ctx, contractID, from, operator); err != nil {
		return sdkerrors.ErrUnauthorized.Wrap(err.Error())
	}

	if err := k.burnToken(ctx, contractID, from, amount); err != nil {
		return err
	}

	return nil
}

func (k Keeper) burnToken(ctx sdk.Context, contractID string, addr sdk.AccAddress, amount sdk.Int) error {
	if err := k.subtractToken(ctx, contractID, addr, amount); err != nil {
		return err
	}

	burnt := k.GetBurnt(ctx, contractID)
	burnt = burnt.Add(amount)
	k.setBurnt(ctx, contractID, burnt)

	supply := k.GetSupply(ctx, contractID)
	supply = supply.Sub(amount)
	k.setSupply(ctx, contractID, supply)

	return nil
}

func (k Keeper) getStatistics(ctx sdk.Context, contractID string, keyPrefix []byte) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	amount := sdk.ZeroInt()
	bz := store.Get(statisticsKey(keyPrefix, contractID))
	if bz != nil {
		if err := amount.Unmarshal(bz); err != nil {
			panic(err)
		}
	}

	return amount
}

// The caller must validate `amount`.
func (k Keeper) setStatistics(ctx sdk.Context, contractID string, amount sdk.Int, keyPrefix []byte) {
	store := ctx.KVStore(k.storeKey)
	key := statisticsKey(keyPrefix, contractID)
	if amount.IsZero() {
		store.Delete(key)
	} else {
		bz, err := amount.Marshal()
		if err != nil {
			panic(err)
		}
		store.Set(key, bz)
	}
}

func (k Keeper) GetSupply(ctx sdk.Context, contractID string) sdk.Int {
	return k.getStatistics(ctx, contractID, supplyKeyPrefix)
}

func (k Keeper) GetMinted(ctx sdk.Context, contractID string) sdk.Int {
	return k.getStatistics(ctx, contractID, mintKeyPrefix)
}

func (k Keeper) GetBurnt(ctx sdk.Context, contractID string) sdk.Int {
	return k.getStatistics(ctx, contractID, burnKeyPrefix)
}

func (k Keeper) setSupply(ctx sdk.Context, contractID string, amount sdk.Int) {
	k.setStatistics(ctx, contractID, amount, supplyKeyPrefix)
}

func (k Keeper) setMinted(ctx sdk.Context, contractID string, amount sdk.Int) {
	k.setStatistics(ctx, contractID, amount, mintKeyPrefix)
}

func (k Keeper) setBurnt(ctx sdk.Context, contractID string, amount sdk.Int) {
	k.setStatistics(ctx, contractID, amount, burnKeyPrefix)
}

func (k Keeper) Modify(ctx sdk.Context, contractID string, grantee sdk.AccAddress, changes []token.Pair) error {
	if err := k.modify(ctx, contractID, changes); err != nil {
		return err
	}

	event := token.EventModified{
		ContractId: contractID,
		Operator:   grantee.String(),
		Changes:    changes,
	}
	ctx.EventManager().EmitEvents(token.NewEventModifyToken(event)) // deprecated
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}
	return nil
}

func (k Keeper) modify(ctx sdk.Context, contractID string, changes []token.Pair) error {
	class, err := k.GetClass(ctx, contractID)
	if err != nil {
		panic(err)
	}

	modifiers := map[token.AttributeKey]func(string){
		token.AttributeKeyName: func(name string) {
			class.Name = name
		},
		token.AttributeKeyImageURI: func(uri string) {
			class.ImageUri = uri
		},
		token.AttributeKeyMeta: func(meta string) {
			class.Meta = meta
		},
	}
	for _, change := range changes {
		key := token.AttributeKeyFromString(change.Field)
		modifiers[key](change.Value)
	}

	k.setClass(ctx, *class)

	return nil
}

func (k Keeper) Grant(ctx sdk.Context, contractID string, granter, grantee sdk.AccAddress, permission token.Permission) {
	k.grant(ctx, contractID, grantee, permission)

	event := token.EventGranted{
		ContractId: contractID,
		Granter:    granter.String(),
		Grantee:    grantee.String(),
		Permission: permission,
	}
	ctx.EventManager().EmitEvent(token.NewEventGrantPermToken(event)) // deprecated
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}
}

func (k Keeper) grant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission token.Permission) {
	k.setGrant(ctx, contractID, grantee, permission)

	if !k.accountKeeper.HasAccount(ctx, grantee) {
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, grantee))
	}
}

func (k Keeper) Abandon(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission token.Permission) {
	k.deleteGrant(ctx, contractID, grantee, permission)

	event := token.EventRenounced{
		ContractId: contractID,
		Grantee:    grantee.String(),
		Permission: permission,
	}
	ctx.EventManager().EmitEvent(token.NewEventRevokePermToken(event)) // deprecated
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}
}

func (k Keeper) GetGrant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission token.Permission) (*token.Grant, error) {
	var grant *token.Grant
	store := ctx.KVStore(k.storeKey)
	if store.Has(grantKey(contractID, grantee, permission)) {
		grant = &token.Grant{
			Grantee:    grantee.String(),
			Permission: permission,
		}
		return grant, nil
	}

	return nil, token.ErrGrantNotFound.Wrapf("%s has no %s", grantee, permission)
}

func (k Keeper) setGrant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission token.Permission) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(contractID, grantee, permission)
	store.Set(key, []byte{})
}

func (k Keeper) deleteGrant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission token.Permission) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(contractID, grantee, permission)
	store.Delete(key)
}
