package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
)

var _ sdk.Msg = (*MsgTransferCFT)(nil)
var _ sdk.Msg = (*MsgTransferCNFT)(nil)
var _ sdk.Msg = (*MsgTransferCFTFrom)(nil)
var _ sdk.Msg = (*MsgTransferCNFTFrom)(nil)

var _ json.Marshaler = (*MsgTransferCFT)(nil)
var _ json.Unmarshaler = (*MsgTransferCFT)(nil)
var _ json.Marshaler = (*MsgTransferCNFT)(nil)
var _ json.Unmarshaler = (*MsgTransferCNFT)(nil)
var _ json.Marshaler = (*MsgTransferCFTFrom)(nil)
var _ json.Unmarshaler = (*MsgTransferCFTFrom)(nil)
var _ json.Marshaler = (*MsgTransferCNFTFrom)(nil)
var _ json.Unmarshaler = (*MsgTransferCNFTFrom)(nil)

type MsgTransferCFT struct {
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
	Amount  sdk.Int        `json:"amount"`
}

func NewMsgTransferCFT(from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) MsgTransferCFT {
	return MsgTransferCFT{
		From:    from,
		To:      to,
		Symbol:  symbol,
		TokenID: tokenID,
		Amount:  amount,
	}
}

func (msg MsgTransferCFT) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferCFT
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferCFT) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferCFT
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferCFT) Route() string { return RouterKey }

func (MsgTransferCFT) Type() string { return "transfer_cft" }

func (msg MsgTransferCFT) ValidateBasic() sdk.Error {
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

func (msg MsgTransferCFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferCFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgTransferCNFT struct {
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}

func NewMsgTransferCNFT(from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) MsgTransferCNFT {
	return MsgTransferCNFT{
		From:    from,
		To:      to,
		Symbol:  symbol,
		TokenID: tokenID,
	}
}

func (msg MsgTransferCNFT) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferCNFT
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferCNFT) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferCNFT
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferCNFT) Route() string { return RouterKey }

func (MsgTransferCNFT) Type() string { return "transfer_cnft" }

func (msg MsgTransferCNFT) ValidateBasic() sdk.Error {
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

func (msg MsgTransferCNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferCNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type MsgTransferCFTFrom struct {
	Proxy   sdk.AccAddress `json:"proxy"`
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
	Amount  sdk.Int        `json:"amount"`
}

func NewMsgTransferCFTFrom(proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) MsgTransferCFTFrom {
	return MsgTransferCFTFrom{
		Proxy:   proxy,
		From:    from,
		To:      to,
		Symbol:  symbol,
		TokenID: tokenID,
		Amount:  amount,
	}
}

func (msg MsgTransferCFTFrom) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferCFTFrom
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferCFTFrom) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferCFTFrom
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferCFTFrom) Route() string { return RouterKey }

func (MsgTransferCFTFrom) Type() string { return "transfer_cft_from" }

func (msg MsgTransferCFTFrom) ValidateBasic() sdk.Error {
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

func (msg MsgTransferCFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferCFTFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proxy}
}

type MsgTransferCNFTFrom struct {
	Proxy   sdk.AccAddress `json:"proxy"`
	From    sdk.AccAddress `json:"from"`
	To      sdk.AccAddress `json:"to"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}

func NewMsgTransferCNFTFrom(proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) MsgTransferCNFTFrom {
	return MsgTransferCNFTFrom{
		Proxy:   proxy,
		From:    from,
		To:      to,
		Symbol:  symbol,
		TokenID: tokenID,
	}
}

func (msg MsgTransferCNFTFrom) MarshalJSON() ([]byte, error) {
	type msgAlias MsgTransferCNFTFrom
	return json.Marshal(msgAlias(msg))
}

func (msg *MsgTransferCNFTFrom) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgTransferCNFTFrom
	return json.Unmarshal(data, msgAlias(msg))
}

func (MsgTransferCNFTFrom) Route() string { return RouterKey }

func (MsgTransferCNFTFrom) Type() string { return "transfer_cnft_from" }

func (msg MsgTransferCNFTFrom) ValidateBasic() sdk.Error {
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

func (msg MsgTransferCNFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferCNFTFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proxy}
}
