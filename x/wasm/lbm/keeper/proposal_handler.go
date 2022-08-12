package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	wasmkeeper "github.com/line/lbm-sdk/x/wasm/keeper"
	lbmwasmtypes "github.com/line/lbm-sdk/x/wasm/lbm/types"
	wasmtypes "github.com/line/lbm-sdk/x/wasm/types"
)

// NewWasmProposalHandler creates a new governance Handler for wasm proposals
func NewWasmProposalHandler(k Keeper, enabledProposalTypes []wasmtypes.ProposalType) govtypes.Handler {
	wasmProposalHandler := wasmkeeper.NewWasmProposalHandler(k, enabledProposalTypes)
	return NewWasmProposalHandlerX(wasmProposalHandler, NewGovPermissionKeeper(k), enabledProposalTypes)
}

// NewWasmProposalHandlerX creates a new governance Handler for wasm proposals
func NewWasmProposalHandlerX(wasmProposalHandler govtypes.Handler, k lbmwasmtypes.ContractOpsKeeper, enabledProposalTypes []wasmtypes.ProposalType) govtypes.Handler {
	enabledTypes := make(map[string]struct{}, len(enabledProposalTypes))
	for i := range enabledProposalTypes {
		enabledTypes[string(enabledProposalTypes[i])] = struct{}{}
	}
	return func(ctx sdk.Context, content govtypes.Content) error {
		if err := wasmProposalHandler(ctx, content); err != nil {
			if content == nil {
				return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "content must not be empty")
			}
			if _, ok := enabledTypes[content.ProposalType()]; !ok {
				return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unsupported wasm proposal content type: %q", content.ProposalType())
			}
			switch c := content.(type) {
			case *lbmwasmtypes.DeactivateContract:
				return handleDeactivateContractProposal(ctx, k, *c)
			case *lbmwasmtypes.ActivateContract:
				return handleActivateContractProposal(ctx, k, *c)
			default:
				return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized wasm proposal content type: %T", c)
			}
		}
		return nil
	}
}

func handleDeactivateContractProposal(ctx sdk.Context, k lbmwasmtypes.ContractOpsKeeper, p lbmwasmtypes.DeactivateContract) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	contractAddress, err := sdk.AccAddressFromBech32(p.Contract)
	if err != nil {
		return sdkerrors.Wrap(err, "contract")
	}

	err = k.DeactivateContract(ctx, contractAddress)
	if err != nil {
		return err
	}

	event := lbmwasmtypes.EventDeactivateContractProposal{
		ContractAddress: contractAddress.String(),
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		return err
	}

	return nil
}

func handleActivateContractProposal(ctx sdk.Context, k lbmwasmtypes.ContractOpsKeeper, p lbmwasmtypes.ActivateContract) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	contractAddress, err := sdk.AccAddressFromBech32(p.Contract)
	if err != nil {
		return sdkerrors.Wrap(err, "contract")
	}

	err = k.ActivateContract(ctx, contractAddress)
	if err != nil {
		return err
	}

	event := lbmwasmtypes.EventActivateContractProposal{ContractAddress: contractAddress.String()}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		return nil
	}

	return nil
}
