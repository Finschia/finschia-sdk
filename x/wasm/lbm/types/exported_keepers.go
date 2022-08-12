package types

import (
	sdk "github.com/line/lbm-sdk/types"
)

// ViewKeeper provides read only operations
type ViewKeeper interface {
	IterateInactiveContracts(ctx sdk.Context, fn func(contractAddress sdk.AccAddress) bool)
	IsInactiveContract(ctx sdk.Context, contractAddress sdk.AccAddress) bool
}

// ContractOpsKeeper contains mutable operations on a contract.
type ContractOpsKeeper interface {
	// DeactivateContract add the contract address to inactive contract list.
	DeactivateContract(ctx sdk.Context, contractAddress sdk.AccAddress) error

	// ActivateContract remove the contract address from inactive contract list.
	ActivateContract(ctx sdk.Context, contractAddress sdk.AccAddress) error
}
