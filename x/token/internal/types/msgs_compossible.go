package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

var _ sdk.Msg = (*MsgAttach)(nil)
var _ sdk.Msg = (*MsgDetach)(nil)
var _ sdk.Msg = (*MsgAttachFrom)(nil)
var _ sdk.Msg = (*MsgDetachFrom)(nil)

type MsgAttach struct {
	From      sdk.AccAddress `json:"from"`
	Symbol    string         `json:"symbol"`
	ToTokenID string         `json:"to_token_id"`
	TokenID   string         `json:"token_id"`
}

type MsgDetach struct {
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}

func NewMsgAttach(from sdk.AccAddress, symbol string, toTokenID string, tokenID string) MsgAttach {
	return MsgAttach{
		From:      from,
		Symbol:    symbol,
		ToTokenID: toTokenID,
		TokenID:   tokenID,
	}
}

func (msg MsgAttach) MarshalJSON() ([]byte, error) {
	type msgAlias MsgAttach
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgAttach) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgAttach
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgAttach) Route() string { return RouterKey }

func (MsgAttach) Type() string { return "attach" }

func (msg MsgAttach) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}

	if err := linktype.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return sdk.ErrInvalidCoins(msg.Symbol)
	}

	if err := linktype.ValidateTokenID(msg.ToTokenID); err != nil {
		return sdk.ErrInvalidCoins(msg.ToTokenID)
	}

	if err := linktype.ValidateTokenID(msg.TokenID); err != nil {
		return sdk.ErrInvalidCoins(msg.TokenID)
	}

	if msg.ToTokenID == msg.TokenID {
		return ErrCannotAttachToItself(DefaultCodespace, msg.Symbol+msg.TokenID)
	}

	return nil
}

func (msg MsgAttach) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgAttach) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func NewMsgDetach(from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) MsgDetach {
	return MsgDetach{
		From:    from,
		To:      to,
		Symbol:  symbol,
		TokenID: tokenID,
	}
}

func (msg MsgDetach) MarshalJSON() ([]byte, error) {
	type msgAlias MsgDetach
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgDetach) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgDetach
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgDetach) Route() string { return RouterKey }

func (MsgDetach) Type() string { return "detach" }

func (msg MsgDetach) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}

	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("To cannot be empty")
	}

	if err := linktype.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return sdk.ErrInvalidCoins(msg.Symbol)
	}

	if err := linktype.ValidateTokenID(msg.TokenID); err != nil {
		return sdk.ErrInvalidCoins(msg.TokenID)
	}

	return nil
}

func (msg MsgDetach) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDetach) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgAttachFrom struct {
	Proxy     sdk.AccAddress `json:"proxy"`
	From      sdk.AccAddress `json:"from"`
	Symbol    string         `json:"symbol"`
	ToTokenID string         `json:"to_token_id"`
	TokenID   string         `json:"token_id"`
}

func NewMsgAttachFrom(proxy sdk.AccAddress, from sdk.AccAddress, symbol string, toTokenID string, tokenID string) MsgAttachFrom {
	return MsgAttachFrom{
		Proxy:     proxy,
		From:      from,
		Symbol:    symbol,
		ToTokenID: toTokenID,
		TokenID:   tokenID,
	}
}

func (msg MsgAttachFrom) MarshalJSON() ([]byte, error) {
	type msgAlias MsgAttachFrom
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgAttachFrom) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgAttachFrom
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgAttachFrom) Route() string { return RouterKey }

func (MsgAttachFrom) Type() string { return "attach_from" }

func (msg MsgAttachFrom) ValidateBasic() sdk.Error {
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}
	if err := linktype.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return sdk.ErrInvalidCoins(msg.Symbol)
	}
	if err := linktype.ValidateTokenID(msg.ToTokenID); err != nil {
		return sdk.ErrInvalidCoins(msg.ToTokenID)
	}
	if err := linktype.ValidateTokenID(msg.TokenID); err != nil {
		return sdk.ErrInvalidCoins(msg.TokenID)
	}

	if msg.ToTokenID == msg.TokenID {
		return ErrCannotAttachToItself(DefaultCodespace, msg.Symbol+msg.TokenID)
	}

	return nil
}

func (msg MsgAttachFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgAttachFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proxy}
}

type MsgDetachFrom struct {
	Proxy   sdk.AccAddress `json:"proxy"`
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}

func NewMsgDetachFrom(proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) MsgDetachFrom {
	return MsgDetachFrom{
		Proxy:   proxy,
		From:    from,
		To:      to,
		Symbol:  symbol,
		TokenID: tokenID,
	}
}

func (msg MsgDetachFrom) MarshalJSON() ([]byte, error) {
	type msgAlias MsgDetachFrom
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgDetachFrom) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgDetachFrom
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgDetachFrom) Route() string { return RouterKey }

func (MsgDetachFrom) Type() string { return "detach_from" }

func (msg MsgDetachFrom) ValidateBasic() sdk.Error {
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("To cannot be empty")
	}
	if err := linktype.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return sdk.ErrInvalidCoins(msg.Symbol)
	}
	if err := linktype.ValidateTokenID(msg.TokenID); err != nil {
		return sdk.ErrInvalidCoins(msg.TokenID)
	}

	return nil
}

func (msg MsgDetachFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgDetachFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proxy}
}
