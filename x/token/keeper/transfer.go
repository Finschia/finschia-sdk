package keeper

import (
	"github.com/gogo/protobuf/proto"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token"
)

func (k Keeper) Transfer(ctx sdk.Context, from, to sdk.AccAddress, amounts []token.FT) error {
	if err := k.transfer(ctx, from, to, amounts); err != nil {
		return err
	}

	events := make([]proto.Message, len(amounts))
	for _, amount := range amounts {
		events = append(events, &token.EventTransfer{
			ClassId: amount.ClassId,
			From:    from.String(),
			To:      to.String(),
			Amount:  amount.Amount,
		})
	}
	return ctx.EventManager().EmitTypedEvents(events...)
}

func (k Keeper) transfer(ctx sdk.Context, from, to sdk.AccAddress, amounts []token.FT) error {
	if err := k.sendTokens(ctx, from, to, amounts); err != nil {
		return err
	}

	return nil
}

func (k Keeper) TransferFrom(ctx sdk.Context, proxy, from, to sdk.AccAddress, amounts []token.FT) error {
	if err := k.transferFrom(ctx, proxy, from, to, amounts); err != nil {
		return err
	}

	events := make([]proto.Message, len(amounts))
	for _, amount := range amounts {
		events = append(events, &token.EventTransfer{
			ClassId: amount.ClassId,
			From:    from.String(),
			To:      to.String(),
			Amount:  amount.Amount,
		})
	}
	return ctx.EventManager().EmitTypedEvents(events...)
}

func (k Keeper) transferFrom(ctx sdk.Context, proxy, from, to sdk.AccAddress, amounts []token.FT) error {
	for _, amount := range amounts {
		if ok := k.GetApprove(ctx, from, proxy, amount.ClassId); !ok {
			return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not authorized to send %s tokens of %s", proxy, amount.ClassId, from)
		}
	}

	return k.transfer(ctx, from, to, amounts)
}

func (k Keeper) Approve(ctx sdk.Context, approver, proxy sdk.AccAddress, classID string) error {
	if err := k.approve(ctx, approver, proxy, classID); err != nil {
		return err
	}

	return ctx.EventManager().EmitTypedEvent(&token.EventApprove{
		ClassId:  classID,
		Approver: approver.String(),
		Proxy:    proxy.String(),
	})
}

func (k Keeper) approve(ctx sdk.Context, approver, proxy sdk.AccAddress, classID string) error {
	if k.GetApprove(ctx, approver, proxy, classID) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Already approved")
	}
	k.setApprove(ctx, approver, proxy, classID, true)
	return nil
}

func (k Keeper) GetApprove(ctx sdk.Context, approver, proxy sdk.AccAddress, classID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(approveKey(classID, proxy, approver))
}

func (k Keeper) setApprove(ctx sdk.Context, approver, proxy sdk.AccAddress, classID string, set bool) {
	store := ctx.KVStore(k.storeKey)
	key := approveKey(classID, proxy, approver)
	if set {
		store.Set(key, []byte{})
	} else {
		store.Delete(key)
	}
}

func (k Keeper) sendTokens(ctx sdk.Context, from, to sdk.AccAddress, amounts []token.FT) error {
	if err := k.subtractTokens(ctx, from, amounts); err != nil {
		return err
	}

	if err := k.addTokens(ctx, to, amounts); err != nil {
		return err
	}

	// TODO: replace it to HasAccount()
	if acc := k.accountKeeper.GetAccount(ctx, to); acc == nil {
		k.accountKeeper.SetAccount(ctx, k.accountKeeper.NewAccountWithAddress(ctx, to))
	}

	return nil
}

func (k Keeper) subtractTokens(ctx sdk.Context, addr sdk.AccAddress, amounts []token.FT) error {
	for _, amount := range amounts {
		if amount.Amount.IsNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, amount.Amount.String())
		}

		balance := k.GetBalance(ctx, addr, amount.ClassId)
		newAmount := balance.Amount.Sub(amount.Amount)
		if newAmount.IsNegative() {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, "%s is smaller than %s", balance.Amount, amount.Amount)
		}

		newBalance := token.FT{
			ClassId: amount.ClassId,
			Amount:  newAmount,
		}

		if err := k.setBalance(ctx, addr, newBalance); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) addTokens(ctx sdk.Context, addr sdk.AccAddress, amounts []token.FT) error {
	for _, amount := range amounts {
		if amount.Amount.IsNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, amount.Amount.String())
		}

		balance := k.GetBalance(ctx, addr, amount.ClassId)
		newAmount := balance.Amount.Add(amount.Amount)

		newBalance := token.FT{
			ClassId: amount.ClassId,
			Amount:  newAmount,
		}

		if err := k.setBalance(ctx, addr, newBalance); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, classID string) token.FT {
	store := ctx.KVStore(k.storeKey)
	amount := sdk.ZeroInt()
	bz := store.Get(balanceKey(addr, classID))
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

// setBalance sets balance.
// The caller must validate `balance`.
func (k Keeper) setBalance(ctx sdk.Context, addr sdk.AccAddress, balance token.FT) error {
	store := ctx.KVStore(k.storeKey)
	key := balanceKey(addr, balance.ClassId)
	if balance.Amount.IsZero() {
		store.Delete(key)
	} else {
		bz, err := balance.Amount.Marshal()
		if err != nil {
			return err
		}
		store.Set(key, bz)
	}

	return nil
}
