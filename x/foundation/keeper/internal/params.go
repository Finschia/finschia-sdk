package internal

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (k Keeper) GetParams(ctx sdk.Context) foundation.Params {
	var params foundation.Params
	k.paramSpace.GetParamSet(ctx, &params)

	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params foundation.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// aliases
func (k Keeper) GetFoundationTax(ctx sdk.Context) sdk.Dec {
	params := k.GetParams(ctx)

	return params.FoundationTax
}

func (k Keeper) IsCensoredMessage(ctx sdk.Context, msgTypeURL string) bool {
	_, err := k.GetCensorship(ctx, msgTypeURL)
	return err == nil
}
