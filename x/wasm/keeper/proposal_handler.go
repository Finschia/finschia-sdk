package keeper

import (
	"encoding/hex"
	"strconv"
	"strings"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
	"github.com/line/lbm-sdk/x/wasm/types"
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

	err := sdk.ValidateAccAddress(p.RunAs)
	if err != nil {
		return sdkerrors.Wrap(err, "run as address")
	}
	_, err = k.Create(ctx, sdk.AccAddress(p.RunAs), p.WASMByteCode, p.InstantiatePermission)

	return err
}

func handleInstantiateProposal(ctx sdk.Context, k types.ContractOpsKeeper, p types.InstantiateContractProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}
	err := sdk.ValidateAccAddress(p.RunAs)
	if err != nil {
		return sdkerrors.Wrap(err, "run as address")
	}
	err = sdk.ValidateAccAddress(p.Admin)
	if err != nil {
		return sdkerrors.Wrap(err, "admin")
	}

	_, data, err := k.Instantiate(ctx, p.CodeID, sdk.AccAddress(p.RunAs), sdk.AccAddress(p.Admin), p.Msg, p.Label, p.Funds)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGovContractResult,
		sdk.NewAttribute(types.AttributeKeyResultDataHex, hex.EncodeToString(data)),
	))
	return nil
}

func handleMigrateProposal(ctx sdk.Context, k types.ContractOpsKeeper, p types.MigrateContractProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	err := sdk.ValidateAccAddress(p.Contract)
	if err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	err = sdk.ValidateAccAddress(p.RunAs)
	if err != nil {
		return sdkerrors.Wrap(err, "run as address")
	}
	data, err := k.Migrate(ctx, sdk.AccAddress(p.Contract), sdk.AccAddress(p.RunAs), p.CodeID, p.Msg)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeGovContractResult,
		sdk.NewAttribute(types.AttributeKeyResultDataHex, hex.EncodeToString(data)),
	))
	return nil
}

func handleUpdateAdminProposal(ctx sdk.Context, k types.ContractOpsKeeper, p types.UpdateAdminProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}
	err := sdk.ValidateAccAddress(p.Contract)
	if err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	err = sdk.ValidateAccAddress(p.NewAdmin)
	if err != nil {
		return sdkerrors.Wrap(err, "run as address")
	}

	if err := k.UpdateContractAdmin(ctx, sdk.AccAddress(p.Contract), "", sdk.AccAddress(p.NewAdmin)); err != nil {
		return err
	}

	ourEvent := sdk.NewEvent(
		types.EventTypeUpdateAdmin,
		sdk.NewAttribute(types.AttributeKeyContractAddr, p.Contract),
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
	)
	ctx.EventManager().EmitEvent(ourEvent)
	return nil
}

func handleClearAdminProposal(ctx sdk.Context, k types.ContractOpsKeeper, p types.ClearAdminProposal) error {
	if err := p.ValidateBasic(); err != nil {
		return err
	}

	err := sdk.ValidateAccAddress(p.Contract)
	if err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	if err := k.ClearContractAdmin(ctx, sdk.AccAddress(p.Contract), ""); err != nil {
		return err
	}

	ourEvent := sdk.NewEvent(
		types.EventTypeClearAdmin,
		sdk.NewAttribute(types.AttributeKeyContractAddr, p.Contract),
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
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
	err := sdk.ValidateAccAddress(p.Contract)
	if err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	if err = k.UpdateContractStatus(ctx, sdk.AccAddress(p.Contract), "", p.Status); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUpdateContractStatus,
		sdk.NewAttribute(types.AttributeKeyContractAddr, p.Contract),
		sdk.NewAttribute(types.AttributeKeyContractStatus, p.Status.String()),
	))
	return nil
}
