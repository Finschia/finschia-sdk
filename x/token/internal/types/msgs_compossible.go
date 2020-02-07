package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

var _ sdk.Msg = (*MsgAttach)(nil)
var _ sdk.Msg = (*MsgDetach)(nil)

type MsgAttach struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	Symbol      string         `json:"symbol"`
	ToTokenID   string         `json:"to_token_id"`
	TokenID     string         `json:"token_id"`
}

type MsgDetach struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Symbol      string         `json:"symbol"`
	TokenID     string         `json:"token_id"`
}

func NewMsgAttach(from sdk.AccAddress, symbol string, toTokenID string, tokenID string) MsgAttach {
	return MsgAttach{
		FromAddress: from,
		Symbol:      symbol,
		ToTokenID:   toTokenID,
		TokenID:     tokenID,
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
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("FromAddress cannot be empty")
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
	return []sdk.AccAddress{msg.FromAddress}
}

func NewMsgDetach(from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) MsgDetach {
	return MsgDetach{
		FromAddress: from,
		ToAddress:   to,
		Symbol:      symbol,
		TokenID:     tokenID,
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
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("FromAddress cannot be empty")
	}

	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("FromAddress cannot be empty")
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
	return []sdk.AccAddress{msg.FromAddress}
}
