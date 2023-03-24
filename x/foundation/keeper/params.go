package keeper

import (
	"sort"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"

	"github.com/line/lbm-sdk/x/foundation"
)

func (k Keeper) GetParams(ctx sdk.Context) foundation.Params {
	store := ctx.KVStore(k.storeKey)
	key := paramsKey
	bz := store.Get(key)

	var params foundation.Params
	k.cdc.MustUnmarshal(bz, &params)

	return params
}

func (k Keeper) UpdateParams(ctx sdk.Context, params foundation.Params) error {
	// not allowed to set the tax, if it has been already disabled
	if k.GetFoundationTax(ctx).IsZero() && !params.FoundationTax.IsZero() {
		return sdkerrors.ErrInvalidRequest.Wrap("foundation tax has been already disabled")
	}

	// for the cleaning up
	urlRemoved := map[string]bool{}
	for _, url := range k.GetParams(ctx).CensoredMsgTypeUrls {
		urlRemoved[url] = true
	}

	// not allowed to add additional censored messages
	for _, url := range params.CensoredMsgTypeUrls {
		if !urlRemoved[url] {
			return sdkerrors.ErrInvalidRequest.Wrapf("adding a new msg type url of %s", url)
		}
		urlRemoved[url] = false
	}

	// clean up relevant authorizations
	// sort it for the determinism (#)
	urlRemovedSorted := make([]string, 0, len(urlRemoved))
	for url, removed := range urlRemoved {
		if !removed {
			continue
		}
		urlRemovedSorted = append(urlRemovedSorted, url)
	}
	sort.Strings(urlRemovedSorted)

	for _, url := range urlRemovedSorted {
		var grantees []sdk.AccAddress
		k.iterateAuthorizations(ctx, func(grantee sdk.AccAddress, authorization foundation.Authorization) (stop bool) {
			if authorization.MsgTypeURL() == url {
				grantees = append(grantees, grantee)
			}
			return false
		})

		for _, grantee := range grantees {
			k.deleteAuthorization(ctx, grantee, url)
		}
	}

	k.SetParams(ctx, params)
	return nil
}

func (k Keeper) SetParams(ctx sdk.Context, params foundation.Params) {
	bz := k.cdc.MustMarshal(&params)

	store := ctx.KVStore(k.storeKey)
	key := paramsKey
	store.Set(key, bz)
}

// aliases
func (k Keeper) GetFoundationTax(ctx sdk.Context) sdk.Dec {
	params := k.GetParams(ctx)

	return params.FoundationTax
}

func (k Keeper) IsCensoredMessage(ctx sdk.Context, msgTypeURL string) bool {
	params := k.GetParams(ctx)
	for _, censoredURL := range params.CensoredMsgTypeUrls {
		if msgTypeURL == censoredURL {
			return true
		}
	}
	return false
}
