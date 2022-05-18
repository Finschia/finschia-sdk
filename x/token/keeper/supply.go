package keeper

import (
	"github.com/gogo/protobuf/proto"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token"
)

func (k Keeper) Issue(ctx sdk.Context, class token.TokenClass, owner, to sdk.AccAddress, amount sdk.Int) error {
	if err := k.issue(ctx, class, owner, to, amount); err != nil {
		return err
	}

	return ctx.EventManager().EmitTypedEvent(&token.EventIssue{
		ClassId: class.ContractId,
	})
}

func (k Keeper) issue(ctx sdk.Context, class token.TokenClass, owner, to sdk.AccAddress, amount sdk.Int) error {
	if _, err := k.GetClass(ctx, class.ContractId); err == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ID already exists: %s", class.Id)
	}
	if err := k.setClass(ctx, class); err != nil {
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
		k.setGrant(ctx, owner, class.ContractId, token.Permission_name[int32(permission)], true)
	}

	if err := k.mintToken(ctx, class.ContractId, to, amount); err != nil {
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
	if err := k.mint(ctx, grantee, to, amounts); err != nil {
		return err
	}

	events := make([]proto.Message, 0, len(amounts))
	for _, amount := range amounts {
		events = append(events, &token.EventMint{
			ClassId: amount.ClassId,
			To:      to.String(),
			Amount:  amount.Amount,
		})
	}
	return ctx.EventManager().EmitTypedEvents(events...)
}

func (k Keeper) mint(ctx sdk.Context, classID string, grantee, to sdk.AccAddress, amount sdk.Int) error {
	if ok := k.GetGrant(ctx, grantee, amount.ClassId, token.ActionMint); !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized for %s %s tokens", grantee, token.ActionMint, amount.ClassId)
	}

	if err := k.mintTokens(ctx, to, amounts); err != nil {
		return err
	}

	return nil
}

func (k Keeper) mintToken(ctx sdk.Context, classID string, addr sdk.AccAddress, amount sdk.Int) error {
	if err := k.addToken(ctx, classID, addr, amount); err != nil {
		return err
	}

	mint := k.GetMint(ctx, classID)
	mint = mint.Add(amount)
	if err := k.setMint(ctx, classID, mint); err != nil {
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
	if err := k.burn(ctx, from, amounts); err != nil {
		return err
	}

	events := make([]proto.Message, 0, len(amounts))
	for _, amount := range amounts {
		events = append(events, &token.EventBurn{
			ClassId: amount.ClassId,
			From:    from.String(),
			Amount:  amount.Amount,
		})
	}
	return ctx.EventManager().EmitTypedEvents(events...)
}

func (k Keeper) burn(ctx sdk.Context, contractID string, from sdk.AccAddress, amount sdk.Int) error {
	if ok := k.GetGrant(ctx, contractID, from, token.Permission_Burn); !ok {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized for %s %s tokens", from, token.ActionBurn, amount.ClassId)
	}

	if err := k.burnTokens(ctx, from, amounts); err != nil {
		return err
	}

	return nil
}

func (k Keeper) BurnFrom(ctx sdk.Context, contractID string, proxy, from sdk.AccAddress, amount sdk.Int) error {
	if err := k.burnFrom(ctx, proxy, from, amounts); err != nil {
		return err
	}

	events := make([]proto.Message, 0, len(amounts))
	for _, amount := range amounts {
		events = append(events, &token.EventBurn{
			ClassId: amount.ClassId,
			From:    from.String(),
			Amount:  amount.Amount,
		})
	}
	return ctx.EventManager().EmitTypedEvents(events...)
}

func (k Keeper) burnFrom(ctx sdk.Context, contractID string, proxy, from sdk.AccAddress, amount sdk.Int) error {
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

	burn := k.GetBurn(ctx, contractID)
	burn = burn.Add(amount)
	if err := k.setBurn(ctx, contractID, burn); err != nil {
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

func (k Keeper) GetMint(ctx sdk.Context, classID string) sdk.Int {
	return k.getStatistics(ctx, classID, mintKeyPrefix)
}

func (k Keeper) GetBurn(ctx sdk.Context, classID string) sdk.Int {
	return k.getStatistics(ctx, classID, burnKeyPrefix)
}

func (k Keeper) setSupply(ctx sdk.Context, classID string, amount sdk.Int) error {
	return k.setStatistics(ctx, classID, amount, supplyKeyPrefix)
}

func (k Keeper) setMint(ctx sdk.Context, classID string, amount sdk.Int) error {
	return k.setStatistics(ctx, classID, amount, mintKeyPrefix)
}

func (k Keeper) setBurn(ctx sdk.Context, classID string, amount sdk.Int) error {
	return k.setStatistics(ctx, classID, amount, burnKeyPrefix)
}

func (k Keeper) Modify(ctx sdk.Context, classID string, grantee sdk.AccAddress, changes []token.Pair) error {
	if err := k.modify(ctx, classID, grantee, changes); err != nil {
		return err
	}

	events := make([]proto.Message, 0, len(changes))
	for _, change := range changes {
		events = append(events, &token.EventModify{
			ClassId: classID,
			Key:     change.Key,
			Value:   change.Value,
		})
	}
	return ctx.EventManager().EmitTypedEvents(events...)
}

func (k Keeper) modify(ctx sdk.Context, classID string, grantee sdk.AccAddress, changes []token.Pair) error {
	if !k.GetGrant(ctx, grantee, classID, token.ActionModify) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized for %s", grantee, token.ActionModify)
	}

	class, err := k.GetClass(ctx, classID)
	if err != nil {
		return err
	}

	modifiers := map[string]func(string){
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
		modifiers[change.Key](change.Value)
	}

	k.setClass(ctx, *class)

	return nil
}

func (k Keeper) Grant(ctx sdk.Context, contractID string, granter, grantee sdk.AccAddress, action string) error {
	if err := k.grant(ctx, contractID, granter, grantee, action); err != nil {
		return err
	}

	return ctx.EventManager().EmitTypedEvent(&token.EventGrant{
		ClassId: classID,
		Grantee: grantee.String(),
		Action:  action,
	})
}

func (k Keeper) grant(ctx sdk.Context, granter, grantee sdk.AccAddress, classID, action string) error {
	if !k.GetGrant(ctx, granter, classID, action) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized for %s", granter, action)
	}
	if k.GetGrant(ctx, grantee, classID, action) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "%s is already granted for %s", grantee, action)
	}

	k.setGrant(ctx, grantee, classID, action, true)

	// TODO: replace it to HasAccount()
	if acc := k.accountKeeper.GetAccount(ctx, grantee); acc == nil {
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, grantee))
	}

	return nil
}

func (k Keeper) Abandon(ctx sdk.Context, grantee sdk.AccAddress, classID, action string) error {
	if err := k.abandon(ctx, grantee, classID, action); err != nil {
		return err
	}

	return ctx.EventManager().EmitTypedEvent(&token.EventAbandon{
		ClassId: classID,
		Grantee: grantee.String(),
		Action:  action,
	})
}

func (k Keeper) abandon(ctx sdk.Context, grantee sdk.AccAddress, classID, permission string) error {
	if k.GetGrant(ctx, classID, grantee, permission) == nil {
		return sdkerrors.ErrNotFound.Wrapf("%s is not authorized for %s", grantee, permission)
	}

	k.deleteGrant(ctx, classID, grantee, permission)

	return nil
}

func (k Keeper) GetGrant(ctx sdk.Context, classID string, grantee sdk.AccAddress, permission token.Permission) *token.Grant {
	var grant *token.Grant
	store := ctx.KVStore(k.storeKey)
	if store.Has(grantKey(classID, grantee, permission)) {
		grant = &token.Grant{
			Grantee: grantee.String(),
			Permission: token.Permission_name[int32(permission)],
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

func (k Keeper) setGrant(ctx sdk.Context, classID string, grantee sdk.AccAddress, permission string) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(classID, grantee, permission)
	store.Set(key, []byte{})
}

func (k Keeper) deleteGrant(ctx sdk.Context, classID string, grantee sdk.AccAddress, permission string) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(classID, grantee, permission)
	store.Delete(key)
}
