package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token"
)

func (k Keeper) Issue(ctx sdk.Context, class token.TokenClass, owner, to sdk.AccAddress, amount sdk.Int) error {
	if err := k.issue(ctx, class, owner, to, amount); err != nil {
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

func (k Keeper) issue(ctx sdk.Context, class token.TokenClass, owner, to sdk.AccAddress, amount sdk.Int) error {
	if _, err := k.GetClass(ctx, class.ContractId); err == nil {
		return sdkerrors.ErrNotFound.Wrapf("ID already exists: %s")
	}
	if err := k.setClass(ctx, class); err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetClass(ctx sdk.Context, classID string) (*token.TokenClass, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(classKey(classID))
	if bz == nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("No information for %s", classID)
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

func (k Keeper) Mint(ctx sdk.Context, classID string, grantee, to sdk.AccAddress, amount sdk.Int) error {
	if err := k.mint(ctx, classID, grantee, to, amount); err != nil {
		return err
	}

	event := token.EventMinted{
		ContractId: classID,
		Operator:   grantee.String(),
		To:         to.String(),
		Amount:     amount,
	}
	ctx.EventManager().EmitEvent(token.NewEventMintToken(event)) // deprecated
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) mint(ctx sdk.Context, classID string, grantee, to sdk.AccAddress, amount sdk.Int) error {
	if k.GetGrant(ctx, classID, grantee, token.Permission_Mint) == nil {
		return sdkerrors.ErrUnauthorized.Wrapf("%s is not authorized for %s %s tokens", grantee, token.Permission_Mint.String(), classID)
	}

	if err := k.mintToken(ctx, classID, to, amount); err != nil {
		return err
	}

	return nil
}

func (k Keeper) mintToken(ctx sdk.Context, classID string, addr sdk.AccAddress, amount sdk.Int) error {
	if err := k.addToken(ctx, classID, addr, amount); err != nil {
		return err
	}

	minted := k.GetMinted(ctx, classID)
	minted = minted.Add(amount)
	if err := k.setMinted(ctx, classID, minted); err != nil {
		return err
	}

	supply := k.GetSupply(ctx, classID)
	supply = supply.Add(amount)
	if err := k.setSupply(ctx, classID, supply); err != nil {
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

func (k Keeper) OperatorBurn(ctx sdk.Context, contractID string, proxy, from sdk.AccAddress, amount sdk.Int) error {
	if err := k.operatorBurn(ctx, contractID, proxy, from, amount); err != nil {
		return err
	}

	event := token.EventBurned{
		ContractId: contractID,
		Operator:   proxy.String(),
		From:       from.String(),
		Amount:     amount,
	}
	ctx.EventManager().EmitEvent(token.NewEventBurnTokenFrom(event)) // deprecated
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) operatorBurn(ctx sdk.Context, contractID string, proxy, from sdk.AccAddress, amount sdk.Int) error {
	grant := k.GetGrant(ctx, contractID, proxy, token.Permission_Burn)
	authorization := k.GetAuthorization(ctx, contractID, from, proxy)
	if grant == nil || authorization == nil || true {
		return sdkerrors.ErrUnauthorized.Wrapf("%s is not authorized for %s tokens from %s", proxy, token.Permission_Burn.String(), from)
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

func (k Keeper) getStatistics(ctx sdk.Context, classID string, keyPrefix []byte) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	amount := sdk.ZeroInt()
	bz := store.Get(statisticsKey(keyPrefix, classID))
	if bz != nil {
		if err := amount.Unmarshal(bz); err != nil {
			panic(err)
		}
	}

	return amount
}

// The caller must validate `amount`.
func (k Keeper) setStatistics(ctx sdk.Context, classID string, amount sdk.Int, keyPrefix []byte) error {
	store := ctx.KVStore(k.storeKey)
	key := statisticsKey(keyPrefix, classID)
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

func (k Keeper) GetSupply(ctx sdk.Context, classID string) sdk.Int {
	return k.getStatistics(ctx, classID, supplyKeyPrefix)
}

func (k Keeper) GetMinted(ctx sdk.Context, classID string) sdk.Int {
	return k.getStatistics(ctx, classID, mintKeyPrefix)
}

func (k Keeper) GetBurnt(ctx sdk.Context, classID string) sdk.Int {
	return k.getStatistics(ctx, classID, burnKeyPrefix)
}

func (k Keeper) setSupply(ctx sdk.Context, classID string, amount sdk.Int) error {
	return k.setStatistics(ctx, classID, amount, supplyKeyPrefix)
}

func (k Keeper) setMinted(ctx sdk.Context, classID string, amount sdk.Int) error {
	return k.setStatistics(ctx, classID, amount, mintKeyPrefix)
}

func (k Keeper) setBurnt(ctx sdk.Context, classID string, amount sdk.Int) error {
	return k.setStatistics(ctx, classID, amount, burnKeyPrefix)
}

func (k Keeper) Modify(ctx sdk.Context, classID string, grantee sdk.AccAddress, changes []token.Pair) error {
	if err := k.modify(ctx, classID, grantee, changes); err != nil {
		return err
	}

	event := token.EventModified{
		ContractId: classID,
		Changes:    changes,
	}
	ctx.EventManager().EmitEvents(token.NewEventModifyToken(event)) // deprecated
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) modify(ctx sdk.Context, classID string, grantee sdk.AccAddress, changes []token.Pair) error {
	if k.GetGrant(ctx, classID, grantee, token.Permission_Modify) == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized for %s", grantee, token.Permission_Modify.String())
	}

	class, err := k.GetClass(ctx, classID)
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
	if err := k.grant(ctx, contractID, grantee, permission); err != nil {
		return err
	}

	event := token.EventGrant{
		ContractId: contractID,
		Grantee:    grantee.String(),
		Permission: permission.String(),
	}
	ctx.EventManager().EmitEvent(token.NewEventGrantPermToken(event)) // deprecated
	return ctx.EventManager().EmitTypedEvent(&event)
}

func (k Keeper) grant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission token.Permission) error {
	k.setGrant(ctx, contractID, grantee, permission)

	if !k.accountKeeper.HasAccount(ctx, grantee) {
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, grantee))
	}

	return nil
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

func (k Keeper) GetGrant(ctx sdk.Context, classID string, grantee sdk.AccAddress, permission token.Permission) *token.Grant {
	var grant *token.Grant
	store := ctx.KVStore(k.storeKey)
	if store.Has(grantKey(classID, grantee, permission)) {
		grant = &token.Grant{
			Grantee:    grantee.String(),
			Permission: permission.String(),
		}
	}

	return grant
}

func (k Keeper) GetGrants(ctx sdk.Context) []token.Grant {
	var grants []token.Grant
	k.iterateGrants(ctx, func(classID string, grant token.Grant) (stop bool) {
		grants = append(grants, grant)
		return false
	})

	return grants
}

func (k Keeper) setGrant(ctx sdk.Context, classID string, grantee sdk.AccAddress, permission token.Permission) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(classID, grantee, permission)
	store.Set(key, []byte{})
}

func (k Keeper) deleteGrant(ctx sdk.Context, classID string, grantee sdk.AccAddress, permission token.Permission) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(classID, grantee, permission)
	store.Delete(key)
}
