package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = (*MsgPublishToken)(nil)
var _ sdk.Msg = (*MsgMint)(nil)
var _ sdk.Msg = (*MsgBurn)(nil)

type MsgPublishToken struct {
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Amount   sdk.Int        `json:"amount"`
	Owner    sdk.AccAddress `json:"owner"`
	Mintable bool           `json:"mintable"`
}

func NewMsgPublishToken(name, symbol string, amount sdk.Int, owner sdk.AccAddress, mintable bool) MsgPublishToken {
	return MsgPublishToken{
		Name:     name,
		Symbol:   symbol,
		Amount:   amount,
		Owner:    owner,
		Mintable: mintable,
	}
}

func (msg MsgPublishToken) Route() string { return RouterKey }

func (msg MsgPublishToken) Type() string { return "publish_token" }

func (msg MsgPublishToken) ValidateBasic() sdk.Error {
	if len(msg.Name) == 0 || len(msg.Symbol) == 0 {
		return sdk.ErrUnknownRequest("name and symbol cannot be empty")
	}

	coin := sdk.NewCoin(msg.Name, msg.Amount)
	if !coin.IsValid() {
		return sdk.ErrInvalidCoins(coin.String())
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}

	if len(msg.Symbol) <= 5 {
		return ErrTokenSymbolLength(DefaultCodespace)
	}
	return nil
}

func (msg MsgPublishToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgPublishToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func NewMsgMint(to sdk.AccAddress, amount sdk.Coins) MsgMint {
	return MsgMint{
		To:     to,
		Amount: amount,
	}
}

type MsgMint struct {
	To     sdk.AccAddress `json:"to"`
	Amount sdk.Coins      `json:"amount"`
}

func (MsgMint) Route() string { return RouterKey }

func (MsgMint) Type() string { return "mint_token" }

func (msg MsgMint) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("amount is not valid")
	}
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("To address cannot be empty")
	}
	return nil
}

func (msg MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.To}
}

func NewMsgBurn(from sdk.AccAddress, amount sdk.Coins) MsgBurn {
	return MsgBurn{
		From:   from,
		Amount: amount,
	}
}

type MsgBurn struct {
	From   sdk.AccAddress `json:"from"`
	Amount sdk.Coins      `json:"amount"`
}

func (MsgBurn) Route() string { return RouterKey }

func (MsgBurn) Type() string { return "burn_token" }

func (msg MsgBurn) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("amount is not valid")
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From address cannot be empty")
	}
	return nil
}

func (msg MsgBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func NewMsgGrantPermission(from, to sdk.AccAddress, perm Permission) MsgGrantPermission {
	return MsgGrantPermission{
		From:       from,
		To:         to,
		Permission: perm,
	}
}

type MsgGrantPermission struct {
	From       sdk.AccAddress `json:"from"`
	To         sdk.AccAddress `json:"to"`
	Permission Permission     `json:"permission"`
}

func (MsgGrantPermission) Route() string { return RouterKey }

func (MsgGrantPermission) Type() string { return "grant_permission" }

func (msg MsgGrantPermission) ValidateBasic() sdk.Error {
	if msg.From.Empty() || msg.To.Empty() {
		return sdk.ErrInvalidAddress("addresses cannot be empty")
	}

	if msg.From.Equals(msg.To) {
		return sdk.ErrInvalidAddress("from, to address can not be the same")
	}

	if len(msg.Permission.GetAction()) <= 0 || len(msg.Permission.GetResource()) <= 0 {
		return ErrTokenPermission(DefaultCodespace)
	}
	return nil
}

func (msg MsgGrantPermission) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgGrantPermission) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

func NewMsgRevokePermission(from sdk.AccAddress, perm Permission) MsgRevokePermission {
	return MsgRevokePermission{
		From:       from,
		Permission: perm,
	}
}

type MsgRevokePermission struct {
	From       sdk.AccAddress `json:"from"`
	Permission Permission     `json:"permission"`
}

func (MsgRevokePermission) Route() string { return RouterKey }

func (MsgRevokePermission) Type() string { return "revoke_permission" }

func (msg MsgRevokePermission) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("addresses cannot be empty")
	}

	if len(msg.Permission.GetAction()) <= 0 || len(msg.Permission.GetResource()) <= 0 {
		return ErrTokenPermission(DefaultCodespace)
	}
	return nil
}

func (msg MsgRevokePermission) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgRevokePermission) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
