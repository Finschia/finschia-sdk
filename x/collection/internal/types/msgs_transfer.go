package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
)

var _ sdk.Msg = (*MsgTransferFT)(nil)
var _ sdk.Msg = (*MsgTransferNFT)(nil)
var _ sdk.Msg = (*MsgTransferFTFrom)(nil)
var _ sdk.Msg = (*MsgTransferNFTFrom)(nil)

var _ json.Marshaler = (*MsgTransferFT)(nil)
var _ json.Unmarshaler = (*MsgTransferFT)(nil)
var _ json.Marshaler = (*MsgTransferNFT)(nil)
var _ json.Unmarshaler = (*MsgTransferNFT)(nil)
var _ json.Marshaler = (*MsgTransferFTFrom)(nil)
var _ json.Unmarshaler = (*MsgTransferFTFrom)(nil)
var _ json.Marshaler = (*MsgTransferNFTFrom)(nil)
var _ json.Unmarshaler = (*MsgTransferNFTFrom)(nil)

type MsgTransferFT struct {
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
	Amount  sdk.Int        `json:"amount"`
}

func NewMsgTransferFT(from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) MsgTransferFT {
	return MsgTransferFT{
		From:    from,
		To:      to,
		Symbol:  symbol,
		TokenID: tokenID,
		Amount:  amount,
	}
}

func (msg MsgTransferFT) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferFT
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferFT) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferFT
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferFT) Route() string { return RouterKey }

func (MsgTransferFT) Type() string { return "transfer_ft" }

func (msg MsgTransferFT) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}

	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("To cannot be empty")
	}

	if len(msg.Symbol) == 0 {
		return sdk.ErrInvalidCoins("Token symbol is empty")
	}

	if err := types.ValidateSymbolCollectionToken(msg.Symbol + msg.TokenID); err != nil {
		return sdk.ErrInvalidCoins("Only user defined token is possible: " + msg.Symbol + msg.TokenID)
	}

	if !msg.Amount.IsPositive() {
		return sdk.ErrInsufficientCoins("send amount must be positive")
	}
	return nil
}

func (msg MsgTransferFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgTransferNFT struct {
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}

func NewMsgTransferNFT(from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) MsgTransferNFT {
	return MsgTransferNFT{
		From:    from,
		To:      to,
		Symbol:  symbol,
		TokenID: tokenID,
	}
}

func (msg MsgTransferNFT) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferNFT
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferNFT) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferNFT
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferNFT) Route() string { return RouterKey }

func (MsgTransferNFT) Type() string { return "transfer_nft" }

func (msg MsgTransferNFT) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}

	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("To cannot be empty")
	}

	if len(msg.Symbol) == 0 {
		return sdk.ErrInvalidCoins("Token symbol is empty")
	}

	if err := types.ValidateSymbolCollectionToken(msg.Symbol + msg.TokenID); err != nil {
		return sdk.ErrInvalidCoins("Only user defined token is possible: " + msg.Symbol + msg.TokenID)
	}

	return nil
}

func (msg MsgTransferNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgTransferFTFrom struct {
	Proxy   sdk.AccAddress `json:"proxy"`
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
	Amount  sdk.Int        `json:"amount"`
}

func NewMsgTransferFTFrom(proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) MsgTransferFTFrom {
	return MsgTransferFTFrom{
		Proxy:   proxy,
		From:    from,
		To:      to,
		Symbol:  symbol,
		TokenID: tokenID,
		Amount:  amount,
	}
}

func (msg MsgTransferFTFrom) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferFTFrom
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferFTFrom) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferFTFrom
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferFTFrom) Route() string { return RouterKey }

func (MsgTransferFTFrom) Type() string { return "transfer_ft_from" }

func (msg MsgTransferFTFrom) ValidateBasic() sdk.Error {
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("To cannot be empty")
	}
	if msg.From.Equals(msg.Proxy) {
		return ErrApproverProxySame(DefaultCodespace, msg.From.String())
	}
	if err := types.ValidateSymbolCollectionToken(msg.Symbol + msg.TokenID); err != nil {
		return sdk.ErrInvalidCoins("Only user defined token is possible: " + msg.Symbol + msg.TokenID)
	}
	if !msg.Amount.IsPositive() {
		return sdk.ErrInsufficientCoins("send amount must be positive")
	}
	return nil
}

func (msg MsgTransferFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferFTFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proxy}
}

type MsgTransferNFTFrom struct {
	Proxy   sdk.AccAddress `json:"proxy"`
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}

func NewMsgTransferNFTFrom(proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) MsgTransferNFTFrom {
	return MsgTransferNFTFrom{
		Proxy:   proxy,
		From:    from,
		To:      to,
		Symbol:  symbol,
		TokenID: tokenID,
	}
}

func (msg MsgTransferNFTFrom) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferNFTFrom
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferNFTFrom) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferNFTFrom
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferNFTFrom) Route() string { return RouterKey }

func (MsgTransferNFTFrom) Type() string { return "transfer_nft_from" }

func (msg MsgTransferNFTFrom) ValidateBasic() sdk.Error {
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("To cannot be empty")
	}
	if msg.From.Equals(msg.Proxy) {
		return ErrApproverProxySame(DefaultCodespace, msg.From.String())
	}
	if err := types.ValidateSymbolCollectionToken(msg.Symbol + msg.TokenID); err != nil {
		return sdk.ErrInvalidCoins("Only user defined token is possible: " + msg.Symbol + msg.TokenID)
	}
	return nil
}

func (msg MsgTransferNFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferNFTFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proxy}
}
