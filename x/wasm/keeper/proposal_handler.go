package keeper

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/line/lfb-sdk/types"
	sdkerrors "github.com/line/lfb-sdk/types/errors"
	govtypes "github.com/line/lfb-sdk/x/gov/types"
	"github.com/line/lfb-sdk/x/wasm/types"
)

// NewWasmProposalHandler creates a new governance Handler for wasm proposals
func NewWasmProposalHandler(k decoratedKeeper, enabledProposalTypes []types.ProposalType) govtypes.Handler {
	return NewWasmProposalHandlerX(NewGovPermissionKeeper(k), enabledProposalTypes)
}

// NewWasmProposalHandlerX creates a new governance Handler for wasm proposals
func NewWasmProposalHandlerX(k types.ContractOpsKeeper, enabledProposalTypes []types.ProposalType) govtypes.Handler {
	enabledTypes := make(map[string]struct{}, len(enabledProposalTypes))
	for i := range enabledProposalTypes {
		enabledTypes[string(enabledProposalTypes[i])] = struct{}{}
	}
	return func(ctx sdk.Context, content govtypes.Content) error {
		if content == nil {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "content must not be empty")
		}
		if _, ok := enabledTypes[content.ProposalType()]; !ok {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unsupported wasm proposal content type: %q", content.ProposalType())
		}
		switch c := content.(type) {
		case *types.StoreCodeProposal:
			return handleStoreCodeProposal(ctx, k, *c)
		case *types.InstantiateContractProposal:
			return handleInstantiateProposal(ctx, k, *c)
		case *types.MigrateContractProposal:
			return handleMigrateProposal(ctx, k, *c)
		case *types.UpdateAdminProposal:
			return handleUpdateAdminProposal(ctx, k, *c)
		case *types.ClearAdminProposal:
			return handleClearAdminProposal(ctx, k, *c)
		case *types.PinCodesProposal:
			return handlePinCodesProposal(ctx, k, *c)
		case *types.UnpinCodesProposal:
			return handleUnpinCodesProposal(ctx, k, *c)
		case *types.UpdateContractStatusProposal:
			return handleUpdateContractStatusProposal(ctx, k, *c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized wasm proposal content type: %T", c)
		}
	}
}

func handleStoreCodeProposal(ctx sdk.Context, k types.ContractOpsKeeper, p types.StoreCodeProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	runAsAddr, err := sdk.AccAddressFromBech32(p.RunAs)
	if err != nil {
		return sdkerrors.Wrap(err, "run as address")
	}
	codeID, err := k.Create(ctx, runAsAddr, p.WASMByteCode, p.Source, p.Builder, p.InstantiatePermission)
	if err != nil {
		return err
	}

	ourEvent := sdk.NewEvent(
		types.EventTypeStoreCode,
		sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", codeID)),
	)
	ctx.EventManager().EmitEvent(ourEvent)
	return nil
}

func handleInstantiateProposal(ctx sdk.Context, k types.ContractOpsKeeper, p types.InstantiateContractProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}
	runAsAddr, err := sdk.AccAddressFromBech32(p.RunAs)
	if err != nil {
		return sdkerrors.Wrap(err, "run as address")
	}
	adminAddr, err := sdk.AccAddressFromBech32(p.Admin)
	if err != nil {
		return sdkerrors.Wrap(err, "admin")
	}

	contractAddr, _, err := k.Instantiate(ctx, p.CodeID, runAsAddr, adminAddr, p.InitMsg, p.Label, p.Funds)
	if err != nil {
		return err
	}

	ourEvent := sdk.NewEvent(
		types.EventTypeInstantiateContract,
		sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", p.CodeID)),
		sdk.NewAttribute(types.AttributeKeyContract, contractAddr.String()),
	)
	ctx.EventManager().EmitEvent(ourEvent)
	return nil
}

func handleMigrateProposal(ctx sdk.Context, k types.ContractOpsKeeper, p types.MigrateContractProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	contractAddr, err := sdk.AccAddressFromBech32(p.Contract)
	if err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	runAsAddr, err := sdk.AccAddressFromBech32(p.RunAs)
	if err != nil {
		return sdkerrors.Wrap(err, "run as address")
	}
	res, err := k.Migrate(ctx, contractAddr, runAsAddr, p.CodeID, p.MigrateMsg)
	if err != nil {
		return err
	}
	ourEvent := sdk.NewEvent(
		types.EventTypeMigrateContract,
		sdk.NewAttribute(types.AttributeKeyContract, p.Contract),
		sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", p.CodeID)),
	)
	ctx.EventManager().EmitEvent(ourEvent)

	for _, e := range res.Events {
		attr := make([]sdk.Attribute, len(e.Attributes))
		for i, a := range e.Attributes {
			attr[i] = sdk.NewAttribute(string(a.Key), string(a.Value))
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(e.Type, attr...))
	}
	return nil
}

func handleUpdateAdminProposal(ctx sdk.Context, k types.ContractOpsKeeper, p types.UpdateAdminProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}
	contractAddr, err := sdk.AccAddressFromBech32(p.Contract)
	if err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	newAdminAddr, err := sdk.AccAddressFromBech32(p.NewAdmin)
	if err != nil {
		return sdkerrors.Wrap(err, "run as address")
	}

	if err := k.UpdateContractAdmin(ctx, contractAddr, nil, newAdminAddr); err != nil {
		return err
	}

	ourEvent := sdk.NewEvent(
		types.EventTypeUpdateAdmin,
		sdk.NewAttribute(types.AttributeKeyContract, p.Contract),
	)
	ctx.EventManager().EmitEvent(ourEvent)
	return nil
}

func handleClearAdminProposal(ctx sdk.Context, k types.ContractOpsKeeper, p types.ClearAdminProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	contractAddr, err := sdk.AccAddressFromBech32(p.Contract)
	if err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	if err := k.ClearContractAdmin(ctx, contractAddr, nil); err != nil {
		return err
	}

	ourEvent := sdk.NewEvent(
		types.EventTypeClearAdmin,
		sdk.NewAttribute(types.AttributeKeyContract, p.Contract),
	)
	ctx.EventManager().EmitEvent(ourEvent)
	return nil
}

func handlePinCodesProposal(ctx sdk.Context, k types.ContractOpsKeeper, p types.PinCodesProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}
	for _, v := range p.CodeIDs {
		if err := k.PinCode(ctx, v); err != nil {
			return sdkerrors.Wrapf(err, "code id: %d", v)
		}
	}
	s := make([]string, len(p.CodeIDs))
	for _, v := range p.CodeIDs {
		ourEvent := sdk.NewEvent(
			types.EventTypePinCode,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyCodeID, strconv.FormatUint(v, 10)),
		)
		ctx.EventManager().EmitEvent(ourEvent)
	}

	ourEvent := sdk.NewEvent(
		types.EventTypePinCode,
		sdk.NewAttribute(types.AttributeKeyCodeIDs, strings.Join(s, ",")),
	)
	ctx.EventManager().EmitEvent(ourEvent)

	return nil
}

func handleUnpinCodesProposal(ctx sdk.Context, k types.ContractOpsKeeper, p types.UnpinCodesProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}
	for _, v := range p.CodeIDs {
		if err := k.UnpinCode(ctx, v); err != nil {
			return sdkerrors.Wrapf(err, "code id: %d", v)
		}
	}
	s := make([]string, len(p.CodeIDs))
	for _, v := range p.CodeIDs {
		ourEvent := sdk.NewEvent(
			types.EventTypeUnpinCode,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyCodeID, strconv.FormatUint(v, 10)),
		)
		ctx.EventManager().EmitEvent(ourEvent)
	}

	ourEvent := sdk.NewEvent(
		types.EventTypeUnpinCode,
		sdk.NewAttribute(types.AttributeKeyCodeIDs, strings.Join(s, ",")),
	)
	ctx.EventManager().EmitEvent(ourEvent)

	return nil
}

func handleUpdateContractStatusProposal(ctx sdk.Context, k types.ContractOpsKeeper, p types.UpdateContractStatusProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}
	contractAddr, err := sdk.AccAddressFromBech32(p.Contract)
	if err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	if err = k.UpdateContractStatus(ctx, contractAddr, nil, p.Status); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUpdateContractStatus,
		sdk.NewAttribute(types.AttributeKeyContract, p.Contract),
		sdk.NewAttribute(types.AttributeKeyContractStatus, p.Status.String()),
	))
	return nil
}
