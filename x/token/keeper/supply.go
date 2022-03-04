package keeper

import (
	"github.com/gogo/protobuf/proto"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token"
)

func (k Keeper) Issue(ctx sdk.Context, class token.Token, owner, to sdk.AccAddress, amount sdk.Int) error {
	if err := k.issue(ctx, class, owner, to, amount); err != nil {
		return err
	}

	return ctx.EventManager().EmitTypedEvent(&token.EventIssue{
		ClassId: class.Id,
	})
}

func (k Keeper) issue(ctx sdk.Context, class token.Token, owner, to sdk.AccAddress, amount sdk.Int) error {
	if _, err := k.GetClass(ctx, class.Id); err == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ID already exists: %s", class.Id)
	}
	if err := k.setClass(ctx, class); err != nil {
		return err
	}

	if class.Mintable {
		mintActions := []string{
			token.ActionMint,
			token.ActionBurn,
		}
		for _, action := range mintActions {
			k.setGrant(ctx, owner, class.Id, action, true)
		}
	}
	k.setGrant(ctx, owner, class.Id, token.ActionModify, true)

	// zero check?
	amounts := []token.FT{
		{
			ClassId: class.Id,
			Amount:  amount,
		},
	}
	if err := k.mintTokens(ctx, to, amounts); err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetClass(ctx sdk.Context, classID string) (*token.Token, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(classKey(classID))
	if bz == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "No information for %s", classID)
	}

	var class token.Token
	if err := k.cdc.UnmarshalBinaryBare(bz, &class); err != nil {
		return nil, err
	}

	return &class, nil
}

func (k Keeper) setClass(ctx sdk.Context, class token.Token) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryBare(&class)
	if err != nil {
		return err
	}

	store.Set(classKey(class.Id), bz)

	return nil
}

func (k Keeper) Mint(ctx sdk.Context, grantee, to sdk.AccAddress, amounts []token.FT) error {
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

func (k Keeper) mint(ctx sdk.Context, grantee, to sdk.AccAddress, amounts []token.FT) error {
	for _, amount := range amounts {
		if ok := k.GetGrant(ctx, grantee, amount.ClassId, token.ActionMint); !ok {
			return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized for %s %s tokens", grantee, token.ActionMint, amount.ClassId)
		}
	}

	if err := k.mintTokens(ctx, to, amounts); err != nil {
		return err
	}

	return nil
}

func (k Keeper) mintTokens(ctx sdk.Context, addr sdk.AccAddress, amounts []token.FT) error {
	if err := k.addTokens(ctx, addr, amounts); err != nil {
		return err
	}

	for _, amount := range amounts {
		mint := k.GetMint(ctx, amount.ClassId)
		mint.Amount = mint.Amount.Add(amount.Amount)
		if err := k.setMint(ctx, mint); err != nil {
			return err
		}

		supply := k.GetSupply(ctx, amount.ClassId)
		supply.Amount = supply.Amount.Add(amount.Amount)
		if err := k.setSupply(ctx, supply); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) Burn(ctx sdk.Context, from sdk.AccAddress, amounts []token.FT) error {
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

func (k Keeper) burn(ctx sdk.Context, from sdk.AccAddress, amounts []token.FT) error {
	for _, amount := range amounts {
		if ok := k.GetGrant(ctx, from, amount.ClassId, token.ActionBurn); !ok {
			return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized for %s %s tokens", from, token.ActionBurn, amount.ClassId)
		}
	}

	if err := k.burnTokens(ctx, from, amounts); err != nil {
		return err
	}

	return nil
}

func (k Keeper) BurnFrom(ctx sdk.Context, proxy, from sdk.AccAddress, amounts []token.FT) error {
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

func (k Keeper) burnFrom(ctx sdk.Context, proxy, from sdk.AccAddress, amounts []token.FT) error {
	for _, amount := range amounts {
		granted := k.GetGrant(ctx, proxy, amount.ClassId, token.ActionBurn)
		approved := k.GetApprove(ctx, from, proxy, amount.ClassId)
		if !granted || !approved {
			return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized for %s %s tokens of %s", proxy, token.ActionBurn, amount.ClassId, from)
		}
	}

	if err := k.burnTokens(ctx, from, amounts); err != nil {
		return err
	}

	return nil
}

func (k Keeper) burnTokens(ctx sdk.Context, addr sdk.AccAddress, amounts []token.FT) error {
	if err := k.subtractTokens(ctx, addr, amounts); err != nil {
		return err
	}

	for _, amount := range amounts {
		burn := k.GetBurn(ctx, amount.ClassId)
		burn.Amount = burn.Amount.Add(amount.Amount)
		if err := k.setBurn(ctx, burn); err != nil {
			return err
		}

		supply := k.GetSupply(ctx, amount.ClassId)
		supply.Amount = supply.Amount.Sub(amount.Amount)
		if err := k.setSupply(ctx, supply); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) getStatistics(ctx sdk.Context, classID string, keyPrefix []byte) token.FT {
	store := ctx.KVStore(k.storeKey)
	amount := sdk.ZeroInt()
	bz := store.Get(statisticsKey(keyPrefix, classID))
	if bz != nil {
		if err := amount.Unmarshal(bz); err != nil {
			panic(err)
		}
	}

	return token.FT{
		ClassId: classID,
		Amount:  amount,
	}
}

// The caller must validate `amount`.
func (k Keeper) setStatistics(ctx sdk.Context, amount token.FT, keyPrefix []byte) error {
	store := ctx.KVStore(k.storeKey)
	key := statisticsKey(keyPrefix, amount.ClassId)
	if amount.Amount.IsZero() {
		store.Delete(key)
	} else {
		bz, err := amount.Amount.Marshal()
		if err != nil {
			return err
		}
		store.Set(key, bz)
	}

	return nil
}

func (k Keeper) GetSupply(ctx sdk.Context, classID string) token.FT {
	return k.getStatistics(ctx, classID, supplyKeyPrefix)
}

func (k Keeper) GetMint(ctx sdk.Context, classID string) token.FT {
	return k.getStatistics(ctx, classID, mintKeyPrefix)
}

func (k Keeper) GetBurn(ctx sdk.Context, classID string) token.FT {
	return k.getStatistics(ctx, classID, burnKeyPrefix)
}

func (k Keeper) setSupply(ctx sdk.Context, amount token.FT) error {
	return k.setStatistics(ctx, amount, supplyKeyPrefix)
}

func (k Keeper) setMint(ctx sdk.Context, amount token.FT) error {
	return k.setStatistics(ctx, amount, mintKeyPrefix)
}

func (k Keeper) setBurn(ctx sdk.Context, amount token.FT) error {
	return k.setStatistics(ctx, amount, burnKeyPrefix)
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

func (k Keeper) Grant(ctx sdk.Context, granter, grantee sdk.AccAddress, classID, action string) error {
	if err := k.grant(ctx, granter, grantee, classID, action); err != nil {
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

func (k Keeper) Revoke(ctx sdk.Context, grantee sdk.AccAddress, classID, action string) error {
	if err := k.revoke(ctx, grantee, classID, action); err != nil {
		return err
	}

	return ctx.EventManager().EmitTypedEvent(&token.EventRevoke{
		ClassId: classID,
		Grantee: grantee.String(),
		Action:  action,
	})
}

func (k Keeper) revoke(ctx sdk.Context, grantee sdk.AccAddress, classID, action string) error {
	if !k.GetGrant(ctx, grantee, classID, action) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized for %s", grantee, action)
	}

	k.setGrant(ctx, grantee, classID, action, false)

	return nil
}

func (k Keeper) GetGrant(ctx sdk.Context, grantee sdk.AccAddress, classID, action string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(grantKey(grantee, classID, action))
}

func (k Keeper) setGrant(ctx sdk.Context, grantee sdk.AccAddress, classID, action string, set bool) {
	store := ctx.KVStore(k.storeKey)
	key := grantKey(grantee, classID, action)
	if set {
		store.Set(key, []byte{})
	} else {
		store.Delete(key)
	}
}
