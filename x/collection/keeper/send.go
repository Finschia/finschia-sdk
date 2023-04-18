package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection"
)

func (k Keeper) SendCoins(ctx sdk.Context, contractID string, from, to sdk.AccAddress, amount []collection.Coin) error {
	if err := k.subtractCoins(ctx, contractID, from, amount); err != nil {
		return err
	}
	k.addCoins(ctx, contractID, to, amount)

	// legacy
	for _, coin := range amount {
		if err := collection.ValidateNFTID(coin.TokenId); err == nil {
			k.iterateDescendants(ctx, contractID, coin.TokenId, func(descendantID string, _ int) (stop bool) {
				event := collection.EventOwnerChanged{
					ContractId: contractID,
					TokenId:    descendantID,
					From:       from.String(),
					To:         to.String(),
				}
				if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
					panic(err)
				}
				return false
			})
		}
	}

	return nil
}

func (k Keeper) addCoins(ctx sdk.Context, contractID string, address sdk.AccAddress, amount []collection.Coin) {
	for _, coin := range amount {
		balance := k.GetBalance(ctx, contractID, address, coin.TokenId)
		newBalance := balance.Add(coin.Amount)
		k.setBalance(ctx, contractID, address, coin.TokenId, newBalance)

		if err := collection.ValidateNFTID(coin.TokenId); err == nil {
			k.setOwner(ctx, contractID, coin.TokenId, address)
		}
	}
}

func (k Keeper) subtractCoins(ctx sdk.Context, contractID string, address sdk.AccAddress, amount []collection.Coin) error {
	for _, coin := range amount {
		balance := k.GetBalance(ctx, contractID, address, coin.TokenId)
		newBalance := balance.Sub(coin.Amount)
		if newBalance.IsNegative() {
			return collection.ErrInsufficientToken.Wrapf("%s is smaller than %s", balance, coin.Amount)
		}
		k.setBalance(ctx, contractID, address, coin.TokenId, newBalance)

		if err := collection.ValidateNFTID(coin.TokenId); err == nil {
			k.deleteOwner(ctx, contractID, coin.TokenId)
		}
	}

	return nil
}

func (k Keeper) GetBalance(ctx sdk.Context, contractID string, address sdk.AccAddress, tokenID string) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	key := balanceKey(contractID, address, tokenID)
	bz := store.Get(key)
	if bz == nil {
		return sdk.ZeroInt()
	}

	var balance sdk.Int
	if err := balance.Unmarshal(bz); err != nil {
		panic(err)
	}
	return balance
}

func (k Keeper) setBalance(ctx sdk.Context, contractID string, address sdk.AccAddress, tokenID string, balance sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	key := balanceKey(contractID, address, tokenID)

	if balance.IsZero() {
		store.Delete(key)
	} else {
		bz, err := balance.Marshal()
		if err != nil {
			panic(err)
		}
		store.Set(key, bz)
	}
}

func (k Keeper) AuthorizeOperator(ctx sdk.Context, contractID string, holder, operator sdk.AccAddress) error {
	if _, err := k.GetContract(ctx, contractID); err != nil {
		panic(err)
	}

	if _, err := k.GetAuthorization(ctx, contractID, holder, operator); err == nil {
		return collection.ErrCollectionAlreadyApproved.Wrap("Already authorized")
	}

	k.setAuthorization(ctx, contractID, holder, operator)

	return nil
}

func (k Keeper) RevokeOperator(ctx sdk.Context, contractID string, holder, operator sdk.AccAddress) error {
	if _, err := k.GetAuthorization(ctx, contractID, holder, operator); err != nil {
		return err
	}

	k.deleteAuthorization(ctx, contractID, holder, operator)
	return nil
}

func (k Keeper) GetAuthorization(ctx sdk.Context, contractID string, holder, operator sdk.AccAddress) (*collection.Authorization, error) {
	store := ctx.KVStore(k.storeKey)
	if store.Has(authorizationKey(contractID, operator, holder)) {
		return &collection.Authorization{
			Holder:   holder.String(),
			Operator: operator.String(),
		}, nil
	}
	return nil, collection.ErrCollectionNotApproved.Wrapf("no authorization by %s to %s", holder, operator)
}

func (k Keeper) setAuthorization(ctx sdk.Context, contractID string, holder, operator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := authorizationKey(contractID, operator, holder)
	store.Set(key, []byte{})
}

func (k Keeper) deleteAuthorization(ctx sdk.Context, contractID string, holder, operator sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := authorizationKey(contractID, operator, holder)
	store.Delete(key)
}
