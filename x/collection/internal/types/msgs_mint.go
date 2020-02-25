package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
)

var _ sdk.Msg = (*MsgMintNFT)(nil)

type MsgMintNFT struct {
	From      sdk.AccAddress `json:"from"`
	To        sdk.AccAddress `json:"to"`
	Symbol    string         `json:"symbol"`
	Name      string         `json:"name"`
	TokenURI  string         `json:"token_uri"`
	TokenType string         `json:"token_type"`
}

func NewMsgMintNFT(from, to sdk.AccAddress, name, symbol, tokenURI string, tokenType string) MsgMintNFT {
	return MsgMintNFT{
		From:      from,
		To:        to,
		Symbol:    symbol,
		Name:      name,
		TokenURI:  tokenURI,
		TokenType: tokenType,
	}
}

func (msg MsgMintNFT) Route() string                { return RouterKey }
func (msg MsgMintNFT) Type() string                 { return "mint_nft" }
func (msg MsgMintNFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgMintNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMintNFT) ValidateBasic() sdk.Error {
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	if len(msg.Name) == 0 {
		return ErrInvalidTokenName(DefaultCodespace, msg.Name)
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("from address cannot be empty")
	}
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("to address cannot be empty")
	}

	if err := types.ValidateTokenTypeNFT(msg.TokenType); err != nil {
		return ErrInvalidTokenID(DefaultCodespace, err.Error())
	}

	if !ValidTokenURI(msg.TokenURI) {
		return ErrInvalidTokenURILength(DefaultCodespace, msg.TokenURI)
	}

	return nil
}

var _ sdk.Msg = (*MsgBurnNFT)(nil)

type MsgBurnNFT struct {
	From     sdk.AccAddress `json:"from"`
	Symbol   string         `json:"symbol"`
	TokenIDs []string       `json:"token_ids"`
}

func NewMsgBurnNFT(from sdk.AccAddress, symbol string, tokenIDs ...string) MsgBurnNFT {
	return MsgBurnNFT{
		From:     from,
		Symbol:   symbol,
		TokenIDs: tokenIDs,
	}
}

func (msg MsgBurnNFT) Route() string                { return RouterKey }
func (msg MsgBurnNFT) Type() string                 { return "burn_nft" }
func (msg MsgBurnNFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgBurnNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnNFT) ValidateBasic() sdk.Error {
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}

	for _, tokenID := range msg.TokenIDs {
		if err := types.ValidateTokenID(tokenID); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
		if err := types.ValidateTokenTypeNFT(tokenID[:TokenTypeLength]); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
	}

	return nil
}

var _ sdk.Msg = (*MsgBurnNFTFrom)(nil)

type MsgBurnNFTFrom struct {
	Proxy    sdk.AccAddress `json:"proxy"`
	From     sdk.AccAddress `json:"from"`
	Symbol   string         `json:"symbol"`
	TokenIDs []string       `json:"token_ids"`
}

func NewMsgBurnNFTFrom(proxy sdk.AccAddress, from sdk.AccAddress, symbol string, tokenIDs ...string) MsgBurnNFTFrom {
	return MsgBurnNFTFrom{
		Proxy:    proxy,
		From:     from,
		Symbol:   symbol,
		TokenIDs: tokenIDs,
	}
}

func (msg MsgBurnNFTFrom) Route() string                { return RouterKey }
func (msg MsgBurnNFTFrom) Type() string                 { return "burn_nft_from" }
func (msg MsgBurnNFTFrom) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Proxy} }
func (msg MsgBurnNFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnNFTFrom) ValidateBasic() sdk.Error {
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}
	if msg.Proxy.Equals(msg.From) {
		return ErrApproverProxySame(DefaultCodespace, msg.Proxy.String())
	}
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}

	for _, tokenID := range msg.TokenIDs {
		if err := types.ValidateTokenID(tokenID); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
		if err := types.ValidateTokenTypeNFT(tokenID[:TokenTypeLength]); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
	}
	return nil
}

var _ sdk.Msg = (*MsgMintFT)(nil)

type MsgMintFT struct {
	Symbol string         `json:"symbol"`
	From   sdk.AccAddress `json:"from"`
	To     sdk.AccAddress `json:"to"`
	Amount Coins          `json:"amount"`
}

func NewMsgMintFT(symbol string, from, to sdk.AccAddress, amount ...Coin) MsgMintFT {
	return MsgMintFT{
		Symbol: symbol,
		From:   from,
		To:     to,
		Amount: amount,
	}
}
func (MsgMintFT) Route() string                    { return RouterKey }
func (MsgMintFT) Type() string                     { return "mint_ft" }
func (msg MsgMintFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgMintFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMintFT) ValidateBasic() sdk.Error {
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}

	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("from address cannot be empty")
	}
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("to address cannot be empty")
	}

	return nil
}

var _ sdk.Msg = (*MsgBurnFT)(nil)

type MsgBurnFT struct {
	Symbol string         `json:"symbol"`
	From   sdk.AccAddress `json:"from"`
	Amount Coins          `json:"amount"`
}

func NewMsgBurnFT(symbol string, from sdk.AccAddress, amount ...Coin) MsgBurnFT {
	return MsgBurnFT{
		Symbol: symbol,
		From:   from,
		Amount: amount,
	}
}
func (MsgBurnFT) Route() string                    { return RouterKey }
func (MsgBurnFT) Type() string                     { return "burn_ft" }
func (msg MsgBurnFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgBurnFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnFT) ValidateBasic() sdk.Error {
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}

	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From address cannot be empty")
	}
	return nil
}

var _ sdk.Msg = (*MsgBurnFTFrom)(nil)

type MsgBurnFTFrom struct {
	Symbol string         `json:"symbol"`
	Proxy  sdk.AccAddress `json:"proxy"`
	From   sdk.AccAddress `json:"from"`
	Amount Coins          `json:"amount"`
}

func NewMsgBurnFTFrom(symbol string, proxy sdk.AccAddress, from sdk.AccAddress, amount ...Coin) MsgBurnFTFrom {
	return MsgBurnFTFrom{
		Symbol: symbol,
		Proxy:  proxy,
		From:   from,
		Amount: amount,
	}
}

func (MsgBurnFTFrom) Route() string                    { return RouterKey }
func (MsgBurnFTFrom) Type() string                     { return "burn_ft_from" }
func (msg MsgBurnFTFrom) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Proxy} }
func (msg MsgBurnFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnFTFrom) ValidateBasic() sdk.Error {
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}
	if msg.Proxy.Equals(msg.From) {
		return ErrApproverProxySame(DefaultCodespace, msg.Proxy.String())
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}
	return nil
}
