package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	ParamSet = paramtypes.ParamSet

	Subspace interface {
		SetParamSet(ctx sdk.Context, ps ParamSet)
	}
)
