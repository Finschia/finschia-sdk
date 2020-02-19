package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

var _ sdk.Msg = (*MsgTransfer)(nil)

type MsgTransfer struct {
	From   sdk.AccAddress `json:"from"`
	To     sdk.AccAddress `json:"to"`
	Symbol string         `json:"symbol"`
	Amount sdk.Int        `json:"amount"`
}

func NewMsgTransfer(from, to sdk.AccAddress, symbol string, amount sdk.Int) MsgTransfer {
	return MsgTransfer{From: from, To: to, Symbol: symbol, Amount: amount}
}

func (msg MsgTransfer) Route() string { return RouterKey }

func (msg MsgTransfer) Type() string { return "transfer_ft" }

func (msg MsgTransfer) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}

	if err := linktype.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return sdk.ErrInvalidCoins("Only user defined token is possible: " + msg.Symbol)
	}

	if !msg.Amount.IsPositive() {
		return sdk.ErrInsufficientCoins("send amount must be positive")
	}
	return nil
}

func (msg MsgTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
