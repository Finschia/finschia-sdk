package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
func (k Keeper) ApproveCoins(ctx sdk.Context, msg types.MsgProxyApproveCoins) error {
	requestDenom := types.NewProxyDenom(msg.Proxy, msg.OnBehalfOf, msg.Denom)
	requestAllowance := types.NewProxyAllowance(requestDenom, msg.Amount)

	if !k.hasProxyAllowance(ctx, requestAllowance.ProxyDenom) {
		// if Proxy of `by` does not exist, save new approval proxy
		k.setProxyAllowance(ctx, requestAllowance)
		return nil
	}

	px, err := k.GetProxyAllowance(ctx, requestAllowance.ProxyDenom)
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
func (k Keeper) DisapproveCoins(ctx sdk.Context, msg types.MsgProxyDisapproveCoins) error {
	requestDenom := types.NewProxyDenom(msg.Proxy, msg.OnBehalfOf, msg.Denom)
	requestAllowance := types.NewProxyAllowance(requestDenom, msg.Amount)

	px, err := k.GetProxyAllowance(ctx, requestAllowance.ProxyDenom)
	if err != nil {
		return err
	}

	// check if approved coins are enough
	gte, err := px.GTE(requestAllowance)
	if err != nil {
		return err
	}
	if !gte {
		return sdkerrors.Wrapf(types.ErrProxyNotEnoughApprovedCoins, "Approved: %v, Requested: %v", px.Amount, requestAllowance.Amount)
	}

	// subtract approved coins
	newAllowance, err := px.SubAllowance(requestAllowance)
	if err != nil {
		return nil
	}

	// if succeed, update or delete the proxy
	if newAllowance.Amount.Equal(sdk.ZeroInt()) {
		k.deleteProxyAllowance(ctx, newAllowance.ProxyDenom)
	} else {
		k.setProxyAllowance(ctx, newAllowance)
	}

	return nil
}

func (k Keeper) SendCoinsFrom(ctx sdk.Context, msg types.MsgProxySendCoinsFrom) error {
	// convert the request as a proxy
	requestDenom := types.NewProxyDenom(msg.Proxy, msg.OnBehalfOf, msg.Denom)
	requestAllowance := types.NewProxyAllowance(requestDenom, msg.Amount)

	// get a proxy compatible with the request
	px, err := k.GetProxyAllowance(ctx, requestAllowance.ProxyDenom)
	if err != nil {
		return err
	}

	// check if approved coins are enough
	gte, err := px.GTE(requestAllowance)
	if err != nil {
		return err
	}
	if !gte {
		return sdkerrors.Wrapf(types.ErrProxyNotEnoughApprovedCoins, "Approved: %v, Requested: %v", px.Amount, requestAllowance.Amount)
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
		k.deleteProxyAllowance(ctx, newAllowance.ProxyDenom)
	} else {
		k.setProxyAllowance(ctx, newAllowance)
	}

	return nil
}

func (k Keeper) hasProxyAllowance(ctx sdk.Context, pd types.ProxyDenom) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(k.proxyAllowanceKeyOf(pd))
}

func (k Keeper) setProxyAllowance(ctx sdk.Context, pa types.ProxyAllowance) {
	store := ctx.KVStore(k.storeKey)
	store.Set(k.proxyAllowanceKeyOf(pa.ProxyDenom), k.cdc.MustMarshalBinaryBare(pa))
}

func (k Keeper) GetProxyAllowance(ctx sdk.Context, pd types.ProxyDenom) (types.ProxyAllowance, error) {
	store := ctx.KVStore(k.storeKey)

	// retrieve the proxy
	bz := store.Get(k.proxyAllowanceKeyOf(pd))
	if bz == nil {
		return types.ProxyAllowance{}, sdkerrors.Wrapf(types.ErrProxyNotExist, "Proxy: %s, Account: %s", pd.Proxy.String(), pd.OnBehalfOf.String())
	}
	r := &types.ProxyAllowance{}
	k.cdc.MustUnmarshalBinaryBare(bz, r)

	return *r, nil
}

func (k Keeper) deleteProxyAllowance(ctx sdk.Context, pd types.ProxyDenom) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(k.proxyAllowanceKeyOf(pd))
}

// proxy key pattern: #{proxy_address}:#{denom}:#{on_behalf_of_address}
// to extend to query by each from left as prefix
func (k Keeper) proxyAllowanceKeyOf(pd types.ProxyDenom) []byte {
	prefixed := pd.Proxy.String() + ":" + pd.Denom + ":" + pd.OnBehalfOf.String()
	return types.ProxyKey(prefixed)
}
