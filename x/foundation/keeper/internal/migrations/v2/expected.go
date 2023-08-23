package v2

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	paramtypes "github.com/Finschia/finschia-sdk/x/params/types"
)

type (
	ParamSet = paramtypes.ParamSet

	Subspace interface {
		SetParamSet(ctx sdk.Context, ps ParamSet)
	}
)
