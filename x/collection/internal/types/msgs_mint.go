package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
	linktype "github.com/line/link/types"
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
	From    sdk.AccAddress `json:"from"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}

func NewMsgBurnNFT(from sdk.AccAddress, symbol, tokenID string) MsgBurnNFT {
	return MsgBurnNFT{
		From:    from,
		Symbol:  symbol,
		TokenID: tokenID,
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

	if err := types.ValidateTokenID(msg.TokenID); err != nil {
		return ErrInvalidTokenID(DefaultCodespace, err.Error())
	}

	if err := types.ValidateTokenTypeNFT(msg.TokenID[:TokenTypeLength]); err != nil {
		return ErrInvalidTokenID(DefaultCodespace, err.Error())
	}

	return nil
}

var _ sdk.Msg = (*MsgBurnNFTFrom)(nil)

type MsgBurnNFTFrom struct {
	Proxy   sdk.AccAddress `json:"proxy"`
	From    sdk.AccAddress `json:"from"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}

func NewMsgBurnNFTFrom(proxy sdk.AccAddress, from sdk.AccAddress, symbol, tokenID string) MsgBurnNFTFrom {
	return MsgBurnNFTFrom{
		Proxy:   proxy,
		From:    from,
		Symbol:  symbol,
		TokenID: tokenID,
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

	if err := types.ValidateTokenID(msg.TokenID); err != nil {
		return ErrInvalidTokenID(DefaultCodespace, err.Error())
	}

	if err := types.ValidateTokenTypeNFT(msg.TokenID[:TokenTypeLength]); err != nil {
		return ErrInvalidTokenID(DefaultCodespace, err.Error())
	}

	return nil
}

var _ sdk.Msg = (*MsgMintFT)(nil)

type MsgMintFT struct {
	From   sdk.AccAddress            `json:"from"`
	To     sdk.AccAddress            `json:"to"`
	Amount linktype.CoinWithTokenIDs `json:"amount"`
}

func NewMsgMintFT(from, to sdk.AccAddress, amount linktype.CoinWithTokenIDs) MsgMintFT {
	return MsgMintFT{
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
	for _, coin := range msg.Amount {
		if err := linktype.ValidateSymbolUserDefined(coin.Symbol); err != nil {
			return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
		}
		if err := linktype.ValidateTokenID(coin.TokenID); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
		if err := linktype.ValidateTokenTypeFT(coin.TokenID[:TokenTypeLength]); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
		if !coin.IsPositive() {
			return sdk.ErrInvalidCoins("amount is not valid")
		}
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
	From   sdk.AccAddress            `json:"from"`
	Amount linktype.CoinWithTokenIDs `json:"amount"`
}

func NewMsgBurnFT(from sdk.AccAddress, amount linktype.CoinWithTokenIDs) MsgBurnFT {
	return MsgBurnFT{
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
	for _, coin := range msg.Amount {
		if err := linktype.ValidateSymbolUserDefined(coin.Symbol); err != nil {
			return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
		}
		if err := linktype.ValidateTokenID(coin.TokenID); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
		if err := linktype.ValidateTokenTypeFT(coin.TokenID[:TokenTypeLength]); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
		if !coin.IsPositive() {
			return sdk.ErrInvalidCoins("Amount is not valid")
		}
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From address cannot be empty")
	}
	return nil
}

var _ sdk.Msg = (*MsgBurnFTFrom)(nil)

type MsgBurnFTFrom struct {
	Proxy  sdk.AccAddress            `json:"proxy"`
	From   sdk.AccAddress            `json:"from"`
	Amount linktype.CoinWithTokenIDs `json:"amount"`
}

func NewMsgBurnFTFrom(proxy sdk.AccAddress, from sdk.AccAddress, amount linktype.CoinWithTokenIDs) MsgBurnFTFrom {
	return MsgBurnFTFrom{
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
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}
	if msg.Proxy.Equals(msg.From) {
		return ErrApproverProxySame(DefaultCodespace, msg.Proxy.String())
	}
	for _, coin := range msg.Amount {
		if err := linktype.ValidateSymbolUserDefined(coin.Symbol); err != nil {
			return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
		}
		if err := linktype.ValidateTokenID(coin.TokenID); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
		if err := linktype.ValidateTokenTypeFT(coin.TokenID[:TokenTypeLength]); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
		if !coin.IsPositive() {
			return sdk.ErrInvalidCoins("Amount is not valid")
		}
	}
	return nil
}
