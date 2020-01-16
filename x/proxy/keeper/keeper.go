package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cbank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/line/link/x/proxy/types"
)

type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	bankKeeper cbank.Keeper
}

func NewKeeper(cdc *codec.Codec, bankKeeper cbank.Keeper, accountKeeper auth.AccountKeeper, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		bankKeeper: bankKeeper,
	}
}

// approve coins for `by` to be transferred on behalf of `onBehalfOf`
func (k Keeper) ApproveCoins(ctx sdk.Context, msg types.MsgProxyApproveCoins) sdk.Error {
	// convert the request as a proxy
	requestAllowance := types.NewProxyAllowance(msg.Proxy, msg.OnBehalfOf, msg.Denom, msg.Amount)

	if !k.hasProxyAllowance(ctx, requestAllowance) {
		// if Proxy of `by` does not exist, save new approval proxy
		k.setProxyAllowance(ctx, requestAllowance)
		return nil
	}

	px, err := k.getProxyAllowance(ctx, requestAllowance)
	if err != nil {
		return err
	}

	// add approved coins
	newProxyAllowance, err := px.AddAllowance(requestAllowance)
	if err != nil {
		return err
	}
	k.setProxyAllowance(ctx, newProxyAllowance)
	return nil
}

// disapprove coins for `by` to be transferred on behalf of `onBehalfOf`
func (k Keeper) DisapproveCoins(ctx sdk.Context, msg types.MsgProxyDisapproveCoins) sdk.Error {
	requestAllowance := types.NewProxyAllowance(msg.Proxy, msg.OnBehalfOf, msg.Denom, msg.Amount)

	px, err := k.getProxyAllowance(ctx, requestAllowance)
	if err != nil {
		return err
	}

	// check if approved coins are enough
	gte, err := px.GTE(requestAllowance)
	if err != nil {
		return err
	}
	if !gte {
		return types.ErrProxyNotEnoughApprovedCoins(types.DefaultCodespace, px.Amount, requestAllowance.Amount)
	}

	// subtract approved coins
	newAllowance, err := px.SubAllowance(requestAllowance)
	if err != nil {
		return nil
	}
	k.setProxyAllowance(ctx, newAllowance)
	return nil
}

func (k Keeper) SendCoinsFrom(ctx sdk.Context, msg types.MsgProxySendCoinsFrom) sdk.Error {
	// convert the request as a proxy
	requestAllowance := types.NewProxyAllowance(msg.Proxy, msg.OnBehalfOf, msg.Denom, msg.Amount)

	// get a proxy compatible with the request
	px, err := k.getProxyAllowance(ctx, requestAllowance)
	if err != nil {
		return err
	}

	// check if approved coins are enough
	gte, err := px.GTE(requestAllowance)
	if err != nil {
		return err
	}
	if !gte {
		return types.ErrProxyNotEnoughApprovedCoins(types.DefaultCodespace, px.Amount, requestAllowance.Amount)
	}

	// send coins from `OnBehalfOf` to `ToAddress`
	coins := sdk.NewCoins(sdk.NewCoin(msg.Denom, msg.Amount))
	err = k.bankKeeper.SendCoins(ctx, px.OnBehalfOf, msg.ToAddress, coins)
	if err != nil {
		return err
	}

	// if succeed, update or delete the proxy
	newAllowance, err := px.SubAllowance(requestAllowance)
	if err != nil {
		return err
	}
	if newAllowance.Amount.Equal(sdk.ZeroInt()) {
		k.deleteProxyAllowance(ctx, newAllowance)
	} else {
		k.setProxyAllowance(ctx, newAllowance)
	}

	return nil
}

func (k Keeper) hasProxyAllowance(ctx sdk.Context, pa types.ProxyAllowance) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(k.proxyAllowanceKeyOf(pa))
}

func (k Keeper) setProxyAllowance(ctx sdk.Context, pa types.ProxyAllowance) {
	store := ctx.KVStore(k.storeKey)
	store.Set(k.proxyAllowanceKeyOf(pa), k.cdc.MustMarshalBinaryBare(pa))
}

func (k Keeper) getProxyAllowance(ctx sdk.Context, pa types.ProxyAllowance) (types.ProxyAllowance, sdk.Error) {
	store := ctx.KVStore(k.storeKey)

	// retrieve the proxy
	bz := store.Get(k.proxyAllowanceKeyOf(pa))
	if bz == nil {
		return types.ProxyAllowance{}, types.ErrProxyNotExist(types.DefaultCodespace, pa.Proxy.String(), pa.OnBehalfOf.String())
	}
	r := &types.ProxyAllowance{}
	k.cdc.MustUnmarshalBinaryBare(bz, r)

	return *r, nil
}

func (k Keeper) deleteProxyAllowance(ctx sdk.Context, pa types.ProxyAllowance) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(k.proxyAllowanceKeyOf(pa))
}

// proxy key pattern: #{proxy_address}:#{denom}:#{on_behalf_of_address}
// to extend to query by each from left as prefix
func (k Keeper) proxyAllowanceKeyOf(pa types.ProxyAllowance) []byte {
	prefixed := pa.Proxy.String() + ":" + pa.Denom + ":" + pa.OnBehalfOf.String()
	return types.ProxyKey(prefixed)
}
