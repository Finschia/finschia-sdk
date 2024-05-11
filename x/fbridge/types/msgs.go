package types

import sdk "github.com/Finschia/finschia-sdk/types"

var (
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgTransfer{}
	_ sdk.Msg = &MsgProvision{}
	_ sdk.Msg = &MsgHoldTransfer{}
	_ sdk.Msg = &MsgReleaseTransfer{}
	_ sdk.Msg = &MsgRemoveProvision{}
	_ sdk.Msg = &MsgClaimBatch{}
	_ sdk.Msg = &MsgClaim{}
	_ sdk.Msg = &MsgSuggestRole{}
	_ sdk.Msg = &MsgAddVoteForRole{}
	_ sdk.Msg = &MsgSetBridgeStatus{}
)

func (m MsgUpdateParams) ValidateBasic() error { return nil }

func (m MsgUpdateParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

func (m MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgTransfer) ValidateBasic() error { return nil }

func (m MsgTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Sender)}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// Type implements the LegacyMsg.Type method.
func (m MsgTransfer) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgTransfer) Route() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgProvision) ValidateBasic() error { return nil }

func (m MsgProvision) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.From)}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgProvision) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// Type implements the LegacyMsg.Type method.
func (m MsgProvision) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgProvision) Route() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgHoldTransfer) ValidateBasic() error { return nil }

func (m MsgHoldTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.From)}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgHoldTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// Type implements the LegacyMsg.Type method.
func (m MsgHoldTransfer) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgHoldTransfer) Route() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgReleaseTransfer) ValidateBasic() error { return nil }

func (m MsgReleaseTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.From)}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgReleaseTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// Type implements the LegacyMsg.Type method.
func (m MsgReleaseTransfer) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgReleaseTransfer) Route() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgRemoveProvision) ValidateBasic() error { return nil }

func (m MsgRemoveProvision) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.From)}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgRemoveProvision) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// Type implements the LegacyMsg.Type method.
func (m MsgRemoveProvision) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgRemoveProvision) Route() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgClaimBatch) ValidateBasic() error { return nil }

func (m MsgClaimBatch) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.From)}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgClaimBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// Type implements the LegacyMsg.Type method.
func (m MsgClaimBatch) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgClaimBatch) Route() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgClaim) ValidateBasic() error { return nil }

func (m MsgClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.From)}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// Type implements the LegacyMsg.Type method.
func (m MsgClaim) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgClaim) Route() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgSuggestRole) ValidateBasic() error { return nil }

func (m MsgSuggestRole) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.From)}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgSuggestRole) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// Type implements the LegacyMsg.Type method.
func (m MsgSuggestRole) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgSuggestRole) Route() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgAddVoteForRole) ValidateBasic() error { return nil }

func (m MsgAddVoteForRole) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.From)}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgAddVoteForRole) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// Type implements the LegacyMsg.Type method.
func (m MsgAddVoteForRole) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgAddVoteForRole) Route() string {
	return sdk.MsgTypeURL(&m)
}

func (m MsgSetBridgeStatus) ValidateBasic() error { return nil }

func (m MsgSetBridgeStatus) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Guardian)}
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgSetBridgeStatus) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// Type implements the LegacyMsg.Type method.
func (m MsgSetBridgeStatus) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgSetBridgeStatus) Route() string {
	return sdk.MsgTypeURL(&m)
}
