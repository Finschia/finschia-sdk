package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/auth/legacy/legacytx"
)

var (
	_ sdk.Msg = (*MsgUpdateParams)(nil)
	_ sdk.Msg = (*MsgAppendCCBatch)(nil)
	_ sdk.Msg = (*MsgEnqueue)(nil)
	_ sdk.Msg = (*MsgAppendSCCBatch)(nil)
	_ sdk.Msg = (*MsgRemoveSCCBatch)(nil)

	_ legacytx.LegacyMsg = (*MsgUpdateParams)(nil)
	_ legacytx.LegacyMsg = (*MsgAppendCCBatch)(nil)
	_ legacytx.LegacyMsg = (*MsgEnqueue)(nil)
	_ legacytx.LegacyMsg = (*MsgAppendSCCBatch)(nil)
	_ legacytx.LegacyMsg = (*MsgRemoveSCCBatch)(nil)
)

func NewMsgUpdateParams(params Params, authority sdk.AccAddress) *MsgUpdateParams {
	return &MsgUpdateParams{
		Params:    params,
		Authority: authority.String(),
	}
}

func (msg MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{authority}
}

func (msg MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgAppendCCBatch) GetSigners() []sdk.AccAddress {
	sequencer, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{sequencer}
}

func (msg MsgAppendCCBatch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgEnqueue) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{from}
}

func (msg MsgEnqueue) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgAppendSCCBatch) GetSigners() []sdk.AccAddress {
	from, _ := sdk.AccAddressFromBech32(msg.FromAddress)
	return []sdk.AccAddress{from}
}

func (msg MsgAppendSCCBatch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRemoveSCCBatch) GetSigners() []sdk.AccAddress {
	authority, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{authority}
}

func (msg MsgRemoveSCCBatch) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// Deprecated: those methods will be removed after merge the upstream.
// move ValidateBasic to msg server, please refer to https://github.com/cosmos/cosmos-sdk/issues/15648
func (msg MsgUpdateParams) ValidateBasic() error   { return nil }
func (msg MsgUpdateParams) Route() string          { return RouterKey }
func (msg MsgUpdateParams) Type() string           { return sdk.MsgTypeURL(&msg) }
func (msg MsgAppendCCBatch) ValidateBasic() error  { return nil }
func (msg MsgAppendCCBatch) Route() string         { return RouterKey }
func (msg MsgAppendCCBatch) Type() string          { return sdk.MsgTypeURL(&msg) }
func (msg MsgEnqueue) ValidateBasic() error        { return nil }
func (msg MsgEnqueue) Route() string               { return RouterKey }
func (msg MsgEnqueue) Type() string                { return sdk.MsgTypeURL(&msg) }
func (msg MsgAppendSCCBatch) ValidateBasic() error { return nil }
func (msg MsgAppendSCCBatch) Route() string        { return RouterKey }
func (msg MsgAppendSCCBatch) Type() string         { return sdk.MsgTypeURL(&msg) }
func (msg MsgRemoveSCCBatch) ValidateBasic() error { return nil }
func (msg MsgRemoveSCCBatch) Route() string        { return RouterKey }
func (msg MsgRemoveSCCBatch) Type() string         { return sdk.MsgTypeURL(&msg) }
