package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

type (
	ParamSet = paramstypes.ParamSet

	Subspace interface {
		GetParamSet(ctx sdk.Context, ps ParamSet)
	}
)
