package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
)

var _ sdk.Msg = (*MsgIssue)(nil)

type MsgIssue struct {
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
}

func NewMsgIssue(name, symbol, tokenURI string, owner sdk.AccAddress, amount sdk.Int, decimal sdk.Int, mintable bool) MsgIssue {
	return MsgIssue{
		Name:     name,
		Symbol:   symbol,
		Owner:    owner,
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
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	if len(msg.Name) == 0 {
		return ErrInvalidTokenName(DefaultCodespace, msg.Name)
	}
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}

	if msg.Amount.Equal(sdk.NewInt(1)) && msg.Decimals.Equal(sdk.NewInt(0)) && !msg.Mintable {
		return ErrInvalidIssueFT(DefaultCodespace)
	}

	if msg.Decimals.GT(sdk.NewInt(18)) || msg.Decimals.IsNegative() {
		return ErrInvalidTokenDecimals(DefaultCodespace, msg.Decimals)
	}

	coin := sdk.NewCoin(msg.Symbol, msg.Amount)
	if !coin.IsValid() {
		return sdk.ErrInvalidCoins(coin.String())
	}

	return nil
}

var _ sdk.Msg = (*MsgMint)(nil)

type MsgMint struct {
	From   sdk.AccAddress `json:"from"`
	To     sdk.AccAddress `json:"to"`
	Amount sdk.Coins      `json:"amount"`
}

func NewMsgMint(from, to sdk.AccAddress, amount sdk.Coins) MsgMint {
	return MsgMint{
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
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("amount is not valid")
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
	From   sdk.AccAddress `json:"from"`
	Amount sdk.Coins      `json:"amount"`
}

func NewMsgBurn(from sdk.AccAddress, amount sdk.Coins) MsgBurn {
	return MsgBurn{
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
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("amount is not valid")
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("from address cannot be empty")
	}
	return nil
}

var _ sdk.Msg = (*MsgModifyTokenURI)(nil)

type MsgModifyTokenURI struct {
	Owner    sdk.AccAddress `json:"owner"`
	Symbol   string         `json:"symbol"`
	TokenURI string         `json:"token_uri"`
	TokenID  string         `json:"token_id"`
}

func NewMsgModifyTokenURI(owner sdk.AccAddress, symbol, tokenURI, tokenID string) MsgModifyTokenURI {
	return MsgModifyTokenURI{
		Owner:    owner,
		Symbol:   symbol,
		TokenURI: tokenURI,
		TokenID:  tokenID,
	}
}

func (msg MsgModifyTokenURI) Route() string { return RouterKey }
func (msg MsgModifyTokenURI) Type() string  { return "modify_token" }
func (msg MsgModifyTokenURI) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgModifyTokenURI) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }

func (msg MsgModifyTokenURI) ValidateBasic() sdk.Error {
	if msg.Symbol == "" {
		return sdk.ErrInvalidAddress("symbol cannot be empty")
	}

	if err := types.ValidateSymbol(msg.Symbol); err != nil {
		return sdk.ErrInvalidAddress(fmt.Sprintf("invalid symbol pattern found %s", msg.Symbol))
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}

	if msg.TokenID != "" {
		if err := types.ValidateTokenID(msg.TokenID); err != nil {
			return sdk.ErrInvalidAddress(fmt.Sprintf("invalid tokenId pattern found %s", msg.TokenID))
		}
	}
	return nil
}
