package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

var _ sdk.Msg = (*MsgIssue)(nil)

type MsgIssue struct {
	Owner    sdk.AccAddress `json:"owner"`
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	TokenURI string         `json:"token_uri"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
}

func NewMsgIssue(owner sdk.AccAddress, name, symbol, tokenURI string, amount sdk.Int, decimal sdk.Int, mintable bool) MsgIssue {
	return MsgIssue{
		Owner:    owner,
		Name:     name,
		Symbol:   symbol,
		TokenURI: tokenURI,
		Amount:   amount,
		Mintable: mintable,
		Decimals: decimal,
	}
}

func (msg MsgIssue) Route() string                { return RouterKey }
func (msg MsgIssue) Type() string                 { return "issue_token" }
func (msg MsgIssue) GetSignBytes() []byte         { return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)) }
func (msg MsgIssue) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }

func (msg MsgIssue) ValidateBasic() sdk.Error {
	if err := linktype.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	if len(msg.Name) == 0 {
		return ErrInvalidTokenName(DefaultCodespace, msg.Name)
	}
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}

	if !ValidateTokenURI(msg.TokenURI) {
		return ErrInvalidTokenURILength(DefaultCodespace, msg.TokenURI)
	}

	if !ValidateName(msg.Name) {
		return ErrInvalidNameLength(DefaultCodespace, msg.Name)
	}

	if msg.Decimals.GT(sdk.NewInt(18)) || msg.Decimals.IsNegative() {
		return ErrInvalidTokenDecimals(DefaultCodespace, msg.Decimals)
	}

	if msg.Amount.IsNegative() {
		return ErrInvalidAmount(DefaultCodespace, msg.Amount.String())
	}

	return nil
}

var _ sdk.Msg = (*MsgMint)(nil)

type MsgMint struct {
	Symbol string         `json:"symbol"`
	From   sdk.AccAddress `json:"from"`
	To     sdk.AccAddress `json:"to"`
	Amount sdk.Int        `json:"amount"`
}

func NewMsgMint(symbol string, from, to sdk.AccAddress, amount sdk.Int) MsgMint {
	return MsgMint{
		Symbol: symbol,
		From:   from,
		To:     to,
		Amount: amount,
	}
}
func (MsgMint) Route() string                    { return RouterKey }
func (MsgMint) Type() string                     { return "mint" }
func (msg MsgMint) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMint) ValidateBasic() sdk.Error {
	if err := linktype.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	if msg.Amount.IsNegative() {
		return ErrInvalidAmount(DefaultCodespace, msg.Amount.String())
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("from address cannot be empty")
	}
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("to address cannot be empty")
	}
	return nil
}

var _ sdk.Msg = (*MsgBurn)(nil)

type MsgBurn struct {
	Symbol string         `json:"symbol"`
	From   sdk.AccAddress `json:"from"`
	Amount sdk.Int        `json:"amount"`
}

func NewMsgBurn(symbol string, from sdk.AccAddress, amount sdk.Int) MsgBurn {
	return MsgBurn{
		Symbol: symbol,
		From:   from,
		Amount: amount,
	}
}
func (MsgBurn) Route() string                    { return RouterKey }
func (MsgBurn) Type() string                     { return "burn" }
func (msg MsgBurn) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurn) ValidateBasic() sdk.Error {
	if err := linktype.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	if msg.Amount.IsNegative() {
		return ErrInvalidAmount(DefaultCodespace, msg.Amount.String())
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("from address cannot be empty")
	}
	return nil
}
