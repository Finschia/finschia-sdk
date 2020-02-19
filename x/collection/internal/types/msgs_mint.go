package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
	linktype "github.com/line/link/types"
)

var _ sdk.Msg = (*MsgMintCNFT)(nil)

type MsgMintCNFT struct {
	From      sdk.AccAddress `json:"from"`
	To        sdk.AccAddress `json:"to"`
	Symbol    string         `json:"symbol"`
	Name      string         `json:"name"`
	TokenURI  string         `json:"token_uri"`
	TokenType string         `json:"token_type"`
}

func NewMsgMintCNFT(from, to sdk.AccAddress, name, symbol, tokenURI string, tokenType string) MsgMintCNFT {
	return MsgMintCNFT{
		From:      from,
		To:        to,
		Symbol:    symbol,
		Name:      name,
		TokenURI:  tokenURI,
		TokenType: tokenType,
	}
}

func (msg MsgMintCNFT) Route() string                { return RouterKey }
func (msg MsgMintCNFT) Type() string                 { return "mint_cnft" }
func (msg MsgMintCNFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgMintCNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMintCNFT) ValidateBasic() sdk.Error {
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

	return nil
}

var _ sdk.Msg = (*MsgBurnCNFT)(nil)

type MsgBurnCNFT struct {
	From    sdk.AccAddress `json:"from"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}

func NewMsgBurnCNFT(from sdk.AccAddress, symbol, tokenID string) MsgBurnCNFT {
	return MsgBurnCNFT{
		From:    from,
		Symbol:  symbol,
		TokenID: tokenID,
	}
}

func (msg MsgBurnCNFT) Route() string                { return RouterKey }
func (msg MsgBurnCNFT) Type() string                 { return "burn_cnft" }
func (msg MsgBurnCNFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgBurnCNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnCNFT) ValidateBasic() sdk.Error {
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

var _ sdk.Msg = (*MsgBurnCNFTFrom)(nil)

type MsgBurnCNFTFrom struct {
	Proxy   sdk.AccAddress `json:"proxy"`
	From    sdk.AccAddress `json:"from"`
	Symbol  string         `json:"symbol"`
	TokenID string         `json:"token_id"`
}

func NewMsgBurnCNFTFrom(proxy sdk.AccAddress, from sdk.AccAddress, symbol, tokenID string) MsgBurnCNFTFrom {
	return MsgBurnCNFTFrom{
		Proxy:   proxy,
		From:    from,
		Symbol:  symbol,
		TokenID: tokenID,
	}
}

func (msg MsgBurnCNFTFrom) Route() string                { return RouterKey }
func (msg MsgBurnCNFTFrom) Type() string                 { return "burn_cnft_from" }
func (msg MsgBurnCNFTFrom) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Proxy} }
func (msg MsgBurnCNFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnCNFTFrom) ValidateBasic() sdk.Error {
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

var _ sdk.Msg = (*MsgMintCFT)(nil)

type MsgMintCFT struct {
	From   sdk.AccAddress            `json:"from"`
	To     sdk.AccAddress            `json:"to"`
	Amount linktype.CoinWithTokenIDs `json:"amount"`
}

func NewMsgMintCFT(from, to sdk.AccAddress, amount linktype.CoinWithTokenIDs) MsgMintCFT {
	return MsgMintCFT{
		From:   from,
		To:     to,
		Amount: amount,
	}
}
func (MsgMintCFT) Route() string                    { return RouterKey }
func (MsgMintCFT) Type() string                     { return "mint_cft" }
func (msg MsgMintCFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgMintCFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMintCFT) ValidateBasic() sdk.Error {
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

var _ sdk.Msg = (*MsgBurnCFT)(nil)

type MsgBurnCFT struct {
	From   sdk.AccAddress            `json:"from"`
	Amount linktype.CoinWithTokenIDs `json:"amount"`
}

func NewMsgBurnCFT(from sdk.AccAddress, amount linktype.CoinWithTokenIDs) MsgBurnCFT {
	return MsgBurnCFT{
		From:   from,
		Amount: amount,
	}
}
func (MsgBurnCFT) Route() string                    { return RouterKey }
func (MsgBurnCFT) Type() string                     { return "burn_cft" }
func (msg MsgBurnCFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgBurnCFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnCFT) ValidateBasic() sdk.Error {
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

var _ sdk.Msg = (*MsgBurnCFTFrom)(nil)

type MsgBurnCFTFrom struct {
	Proxy  sdk.AccAddress            `json:"proxy"`
	From   sdk.AccAddress            `json:"from"`
	Amount linktype.CoinWithTokenIDs `json:"amount"`
}

func NewMsgBurnCFTFrom(proxy sdk.AccAddress, from sdk.AccAddress, amount linktype.CoinWithTokenIDs) MsgBurnCFTFrom {
	return MsgBurnCFTFrom{
		Proxy:  proxy,
		From:   from,
		Amount: amount,
	}
}

func (MsgBurnCFTFrom) Route() string                    { return RouterKey }
func (MsgBurnCFTFrom) Type() string                     { return "burn_cft_from" }
func (msg MsgBurnCFTFrom) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Proxy} }
func (msg MsgBurnCFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnCFTFrom) ValidateBasic() sdk.Error {
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
