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
	From   sdk.AccAddress `json:"from"`
	To     sdk.AccAddress `json:"to"`
	Symbol string         `json:"symbol"`
	Amount Coins          `json:"amount"`
}

func NewMsgTransferFT(from sdk.AccAddress, to sdk.AccAddress, symbol string, amount ...Coin) MsgTransferFT {
	return MsgTransferFT{
		From:   from,
		To:     to,
		Symbol: symbol,
		Amount: amount,
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

	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}

	if !msg.Amount.IsValid() {
		return ErrInvalidAmount(DefaultCodespace, "invalid amount")
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
	From     sdk.AccAddress `json:"from"`
	To       sdk.AccAddress `json:"to"`
	Symbol   string         `json:"symbol"`
	TokenIDs []string       `json:"token_ids"`
}

func NewMsgTransferNFT(from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenIDs ...string) MsgTransferNFT {
	return MsgTransferNFT{
		From:     from,
		To:       to,
		Symbol:   symbol,
		TokenIDs: tokenIDs,
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

	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	for _, tokenID := range msg.TokenIDs {
		if err := types.ValidateTokenID(tokenID); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
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
	Proxy  sdk.AccAddress `json:"proxy"`
	From   sdk.AccAddress `json:"from"`
	To     sdk.AccAddress `json:"to"`
	Symbol string         `json:"symbol"`
	Amount Coins          `json:"amount"`
}

func NewMsgTransferFTFrom(proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, amount ...Coin) MsgTransferFTFrom {
	return MsgTransferFTFrom{
		Proxy:  proxy,
		From:   from,
		To:     to,
		Symbol: symbol,
		Amount: amount,
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
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	if !msg.Amount.IsValid() {
		return ErrInvalidAmount(DefaultCodespace, "invalid amount")
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
	Proxy    sdk.AccAddress `json:"proxy"`
	From     sdk.AccAddress `json:"from"`
	To       sdk.AccAddress `json:"to"`
	Symbol   string         `json:"symbol"`
	TokenIDs []string       `json:"token_ids"`
}

func NewMsgTransferNFTFrom(proxy sdk.AccAddress, from sdk.AccAddress, to sdk.AccAddress, symbol string, tokenIDs ...string) MsgTransferNFTFrom {
	return MsgTransferNFTFrom{
		Proxy:    proxy,
		From:     from,
		To:       to,
		Symbol:   symbol,
		TokenIDs: tokenIDs,
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
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	for _, tokenID := range msg.TokenIDs {
		if err := types.ValidateTokenID(tokenID); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
	}
	return nil
}

func (msg MsgTransferNFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransferNFTFrom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proxy}
}
