package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token"
)

func (k Keeper) Issue(ctx sdk.Context, class token.TokenClass, owner, to sdk.AccAddress, amount sdk.Int) error {
	if err := k.issue(ctx, class); err != nil {
		return err
	}

	permissions := []token.Permission{
		token.Permission_Modify,
	}
	if class.Mintable {
		permissions = append(permissions,
			token.Permission_Mint,
			token.Permission_Burn,
		)
	}
	for _, permission := range permissions {
		if err := k.Grant(ctx, class.ContractId, owner, owner, permission); err != nil {
			return err
		}
	}

	if err := k.mintToken(ctx, class.ContractId, to, amount); err != nil {
		return err
	}
	if err := ctx.EventManager().EmitTypedEvent(&token.EventMinted{
		ContractId: class.ContractId,
		Operator:   owner.String(),
		To:         to.String(),
		Amount:     amount,
	}); err != nil {
		return err
	}

	event := token.EventIssue{
		ContractId: class.ContractId,
		Name:       class.Name,
		Symbol:     class.Symbol,
		Uri:        class.ImageUri,
		Meta:       class.Meta,
		Decimals:   class.Decimals,
		Mintable:   class.Mintable,
	}
	ctx.EventManager().EmitEvent(token.NewEventIssueToken(event, owner, to, amount)) // deprecated
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) issue(ctx sdk.Context, class token.TokenClass) error {
	if _, err := k.GetClass(ctx, class.ContractId); err == nil {
		return sdkerrors.ErrNotFound.Wrapf("ID already exists: %s", class.ContractId)
	}
	if err := k.setClass(ctx, class); err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetClass(ctx sdk.Context, contractID string) (*token.TokenClass, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(classKey(contractID))
	if bz == nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("No information for %s", contractID)
	}

	var class token.TokenClass
	if err := k.cdc.Unmarshal(bz, &class); err != nil {
		return nil, err
	}

	return &class, nil
}

func (k Keeper) setClass(ctx sdk.Context, class token.TokenClass) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&class)
	if err != nil {
		return err
	}

	store.Set(classKey(class.ContractId), bz)

	return nil
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
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) mint(ctx sdk.Context, contractID string, grantee, to sdk.AccAddress, amount sdk.Int) error {
	if k.GetGrant(ctx, contractID, grantee, token.Permission_Mint) == nil {
		return sdkerrors.ErrUnauthorized.Wrapf("%s is not authorized for %s %s tokens", grantee, token.Permission_Mint.String(), contractID)
	}

	if err := k.mintToken(ctx, contractID, to, amount); err != nil {
		return err
	}

	return nil
}

func (k Keeper) mintToken(ctx sdk.Context, contractID string, addr sdk.AccAddress, amount sdk.Int) error {
	if err := k.addToken(ctx, contractID, addr, amount); err != nil {
		return err
	}

	minted := k.GetMinted(ctx, contractID)
	minted = minted.Add(amount)
	if err := k.setMinted(ctx, contractID, minted); err != nil {
		return err
	}

	supply := k.GetSupply(ctx, contractID)
	supply = supply.Add(amount)
	if err := k.setSupply(ctx, contractID, supply); err != nil {
		return err
	}

	return nil
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
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) burn(ctx sdk.Context, contractID string, from sdk.AccAddress, amount sdk.Int) error {
	if k.GetGrant(ctx, contractID, from, token.Permission_Burn) == nil {
		return sdkerrors.ErrUnauthorized.Wrapf("%s is not authorized for %s %s tokens", from, token.Permission_Burn.String(), contractID)
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
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) operatorBurn(ctx sdk.Context, contractID string, operator, from sdk.AccAddress, amount sdk.Int) error {
	grant := k.GetGrant(ctx, contractID, operator, token.Permission_Burn)
	if grant == nil {
		return sdkerrors.ErrUnauthorized.Wrapf("%s has no permission: %s", operator, token.Permission_Burn)
	}
	if _, err := k.GetAuthorization(ctx, contractID, from, operator); err != nil {
		return sdkerrors.ErrUnauthorized.Wrapf(err.Error())
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
	if err := k.setBurnt(ctx, contractID, burnt); err != nil {
		return err
	}

	supply := k.GetSupply(ctx, contractID)
	supply = supply.Sub(amount)
	if err := k.setSupply(ctx, contractID, supply); err != nil {
		return err
	}

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
func (k Keeper) setStatistics(ctx sdk.Context, contractID string, amount sdk.Int, keyPrefix []byte) error {
	store := ctx.KVStore(k.storeKey)
	key := statisticsKey(keyPrefix, contractID)
	if amount.IsZero() {
		store.Delete(key)
	} else {
		bz, err := amount.Marshal()
		if err != nil {
			return err
		}
		store.Set(key, bz)
	}

	return nil
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

func (k Keeper) setSupply(ctx sdk.Context, contractID string, amount sdk.Int) error {
	return k.setStatistics(ctx, contractID, amount, supplyKeyPrefix)
}

func (k Keeper) setMinted(ctx sdk.Context, contractID string, amount sdk.Int) error {
	return k.setStatistics(ctx, contractID, amount, mintKeyPrefix)
}

func (k Keeper) setBurnt(ctx sdk.Context, contractID string, amount sdk.Int) error {
	return k.setStatistics(ctx, contractID, amount, burnKeyPrefix)
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
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) modify(ctx sdk.Context, contractID string, changes []token.Pair) error {
	class, err := k.GetClass(ctx, contractID)
	if err != nil {
		return err
	}

	modifiers := map[string]func(string){
		token.AttributeKey_Name.String(): func(name string) {
			class.Name = name
		},
		token.AttributeKey_ImageURI.String(): func(uri string) {
			class.ImageUri = uri
		},
		token.AttributeKey_Meta.String(): func(meta string) {
			class.Meta = meta
		},
	}
	for _, change := range changes {
		modifiers[change.Field](change.Value)
	}

	k.setClass(ctx, *class)

	return nil
}

func (k Keeper) Grant(ctx sdk.Context, contractID string, granter, grantee sdk.AccAddress, permission token.Permission) error {
	k.grant(ctx, contractID, grantee, permission)

	event := token.EventGrant{
		ContractId: contractID,
		Granter:    granter.String(),
		Grantee:    grantee.String(),
		Permission: permission.String(),
	}
	ctx.EventManager().EmitEvent(token.NewEventGrantPermToken(event)) // deprecated
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) grant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission token.Permission) {
	k.setGrant(ctx, contractID, grantee, permission)

	if !k.accountKeeper.HasAccount(ctx, grantee) {
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, grantee))
	}
}

func (k Keeper) Abandon(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission token.Permission) error {
	k.deleteGrant(ctx, contractID, grantee, permission)

	event := token.EventAbandon{
		ContractId: contractID,
		Grantee:    grantee.String(),
		Permission: permission.String(),
	}
	ctx.EventManager().EmitEvent(token.NewEventRevokePermToken(event)) // deprecated
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) GetGrant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission token.Permission) *token.Grant {
	var grant *token.Grant
	store := ctx.KVStore(k.storeKey)
	if store.Has(grantKey(contractID, grantee, permission)) {
		grant = &token.Grant{
			Grantee:    grantee.String(),
			Permission: permission.String(),
		}
	}

	return grant
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
