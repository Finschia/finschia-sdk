package lbmtypes

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	wasmtypes "github.com/line/lbm-sdk/x/wasm/types"
)

func (msg MsgStoreCodeAndInstantiateContract) Route() string {
	return wasmtypes.RouterKey
}

func (msg MsgStoreCodeAndInstantiateContract) Type() string {
	return "store-code-and-instantiate"
}

func (msg MsgStoreCodeAndInstantiateContract) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return err
	}

	if err := validateWasmCode(msg.WASMByteCode); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "code bytes %s", err.Error())
	}

	if msg.InstantiatePermission != nil {
		if err := msg.InstantiatePermission.ValidateBasic(); err != nil {
			return sdkerrors.Wrap(err, "instantiate permission")
		}
	}

	if err := validateLabel(msg.Label); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "label is required")
	}

	if !msg.Funds.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}

	if len(msg.Admin) != 0 {
		if _, err := sdk.AccAddressFromBech32(msg.Admin); err != nil {
			return sdkerrors.Wrap(err, "admin")
		}
	}

	if err := msg.Msg.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "payload msg")
	}
	return nil
}

func (msg MsgStoreCodeAndInstantiateContract) GetSignBytes() []byte {
	return sdk.MustSortJSON(wasmtypes.ModuleCdc.MustMarshalJSON(&msg))
}

func (msg MsgStoreCodeAndInstantiateContract) GetSigners() []sdk.AccAddress {
	senderAddr := sdk.MustAccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{senderAddr}
}
