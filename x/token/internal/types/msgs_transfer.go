package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
)

var _ sdk.Msg = (*MsgTransferFT)(nil)
var _ sdk.Msg = (*MsgTransferCFT)(nil)
var _ sdk.Msg = (*MsgTransferCNFT)(nil)

var _ json.Marshaler = (*MsgTransferCFT)(nil)
var _ json.Unmarshaler = (*MsgTransferCFT)(nil)
var _ json.Marshaler = (*MsgTransferCNFT)(nil)
var _ json.Unmarshaler = (*MsgTransferCNFT)(nil)

type MsgTransferFT struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Symbol      string         `json:"symbol"`
	Amount      sdk.Int        `json:"amount"`
}

func NewMsgTransferFT(fromAddr, toAddr sdk.AccAddress, symbol string, amount sdk.Int) MsgTransferFT {
	return MsgTransferFT{FromAddress: fromAddr, ToAddress: toAddr, Symbol: symbol, Amount: amount}
}

func (msg MsgTransferFT) Route() string { return RouterKey }

func (msg MsgTransferFT) Type() string { return "transfer-ft" }

func (msg MsgTransferFT) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}

	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return sdk.ErrInvalidCoins("Only user defined token is possible: " + msg.Symbol)
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
	return []sdk.AccAddress{msg.FromAddress}
}

type MsgTransferCFT struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Symbol      string         `json:"symbol"`
	TokenID     string         `json:"token_id"`
	Amount      sdk.Int        `json:"amount"`
}

func NewMsgTransferCFT(from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string, amount sdk.Int) MsgTransferCFT {
	return MsgTransferCFT{
		FromAddress: from,
		ToAddress:   to,
		Symbol:      symbol,
		TokenID:     tokenID,
		Amount:      amount,
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

func (MsgTransferCFT) Type() string { return "transfer-cft" }

func (msg MsgTransferCFT) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("FromAddress cannot be empty")
	}

	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("ToAddress cannot be empty")
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
	return []sdk.AccAddress{msg.FromAddress}
}

type MsgTransferCNFT struct {
	FromAddress sdk.AccAddress `json:"from_address"`
	ToAddress   sdk.AccAddress `json:"to_address"`
	Symbol      string         `json:"symbol"`
	TokenID     string         `json:"token_id"`
}

func NewMsgTransferCNFT(from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenID string) MsgTransferCNFT {
	return MsgTransferCNFT{
		FromAddress: from,
		ToAddress:   to,
		Symbol:      symbol,
		TokenID:     tokenID,
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

func (MsgTransferCNFT) Type() string { return "transfer-cnft" }

func (msg MsgTransferCNFT) ValidateBasic() sdk.Error {
	if msg.FromAddress.Empty() {
		return sdk.ErrInvalidAddress("FromAddress cannot be empty")
	}

	if msg.ToAddress.Empty() {
		return sdk.ErrInvalidAddress("ToAddress cannot be empty")
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
	return []sdk.AccAddress{msg.FromAddress}
}
