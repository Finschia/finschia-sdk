package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	lbmwasmtypes "github.com/line/lbm-sdk/x/wasm/lbm/types"
)

var _ lbmwasmtypes.ContractOpsKeeper = PermissionedKeeper{}

type decoratedKeeper interface {
	activateContract(ctx sdk.Context, contractAddress sdk.AccAddress) error
	deactivateContract(ctx sdk.Context, contractAddress sdk.AccAddress) error
}

type PermissionedKeeper struct {
	nested decoratedKeeper
}

func NewPermissionedKeeper(nested decoratedKeeper) *PermissionedKeeper {
	return &PermissionedKeeper{nested: nested}
}

func NewGovPermissionKeeper(nested decoratedKeeper) *PermissionedKeeper {
	return NewPermissionedKeeper(nested)
}

func NewDefaultPermissionKeeper(nested decoratedKeeper) *PermissionedKeeper {
	return NewPermissionedKeeper(nested)
}

func (p PermissionedKeeper) DeactivateContract(ctx sdk.Context, contractAddress sdk.AccAddress) error {
	return p.nested.deactivateContract(ctx, contractAddress)
}

func (p PermissionedKeeper) ActivateContract(ctx sdk.Context, contractAddress sdk.AccAddress) error {
	return p.nested.activateContract(ctx, contractAddress)
}
