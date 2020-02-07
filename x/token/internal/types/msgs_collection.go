package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
	linktype "github.com/line/link/types"
)

var _ sdk.Msg = (*MsgCreateCollection)(nil)

type MsgCreateCollection struct {
	Name   string         `json:"name"`
	Symbol string         `json:"symbol"`
	Owner  sdk.AccAddress `json:"owner"`
}

func NewMsgCreateCollection(name, symbol string, owner sdk.AccAddress) MsgCreateCollection {
	return MsgCreateCollection{
		Name:   name,
		Symbol: symbol,
		Owner:  owner,
	}
}

func (msg MsgCreateCollection) ValidateBasic() sdk.Error {
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	if len(msg.Name) == 0 {
		return ErrInvalidTokenName(DefaultCodespace, msg.Name)
	}
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}
	return nil
}

func (MsgCreateCollection) Route() string { return RouterKey }
func (MsgCreateCollection) Type() string  { return "create_collection" }
func (msg MsgCreateCollection) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgCreateCollection) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

var _ sdk.Msg = (*MsgIssueCFT)(nil)

type MsgIssueCFT struct {
	Symbol   string         `json:"symbol"`
	Name     string         `json:"name"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
}

func NewMsgIssueCFT(name, symbol, tokenURI string, owner sdk.AccAddress, amount sdk.Int, decimal sdk.Int, mintable bool) MsgIssueCFT {
	return MsgIssueCFT{
		Symbol:   symbol,
		Name:     name,
		Owner:    owner,
		TokenURI: tokenURI,
		Amount:   amount,
		Mintable: mintable,
		Decimals: decimal,
	}
}

func (msg MsgIssueCFT) Route() string                { return RouterKey }
func (msg MsgIssueCFT) Type() string                 { return "issue_ft_collection" }
func (msg MsgIssueCFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }
func (msg MsgIssueCFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIssueCFT) ValidateBasic() sdk.Error {
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

	return nil
}

var _ sdk.Msg = (*MsgIssueCNFT)(nil)

type MsgIssueCNFT struct {
	Symbol string         `json:"symbol"`
	Owner  sdk.AccAddress `json:"owner"`
}

func NewMsgIssueCNFT(symbol string, owner sdk.AccAddress) MsgIssueCNFT {
	return MsgIssueCNFT{
		Symbol: symbol,
		Owner:  owner,
	}
}

func (MsgIssueCNFT) Route() string                    { return RouterKey }
func (MsgIssueCNFT) Type() string                     { return "issue_nft_collection" }
func (msg MsgIssueCNFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }
func (msg MsgIssueCNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIssueCNFT) ValidateBasic() sdk.Error {
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}
	return nil
}

var _ sdk.Msg = (*MsgMintCNFT)(nil)

type MsgMintCNFT struct {
	Symbol    string         `json:"symbol"`
	Name      string         `json:"name"`
	From      sdk.AccAddress `json:"from"`
	To        sdk.AccAddress `json:"to"`
	TokenURI  string         `json:"token_uri"`
	TokenType string         `json:"token_type"`
}

func NewMsgMintCNFT(name, symbol, tokenURI string, tokenType string, from, to sdk.AccAddress) MsgMintCNFT {
	return MsgMintCNFT{
		Symbol:    symbol,
		Name:      name,
		From:      from,
		To:        to,
		TokenURI:  tokenURI,
		TokenType: tokenType,
	}
}

func (msg MsgMintCNFT) Route() string                { return RouterKey }
func (msg MsgMintCNFT) Type() string                 { return "mint_nft_collection" }
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

/*
XXX: not yet supported.
TODO: when the requirements are fixed. Enable it.
var _ sdk.Msg = (*MsgBurnCNFT)(nil)
type MsgBurnCNFT struct {
	Symbol  string         `json:"symbol"`
	Owner   sdk.AccAddress `json:"owner"`
	TokenID string         `json:"token_id"`
}

func NewMsgBurnCNFT(symbol, tokenID string) MsgBurnCNFT {
	return MsgBurnCNFT{
		Symbol:  symbol,
		TokenID: tokenID,
	}
}

func (msg MsgBurnCNFT) Route() string                { return RouterKey }
func (msg MsgBurnCNFT) Type() string                 { return "burn_nft_collection" }
func (msg MsgBurnCNFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }
func (msg MsgBurnCNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnCNFT) ValidateBasic() sdk.Error {
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}

	if err := types.ValidateTokenID(msg.TokenID); err != nil {
		return ErrInvalidTokenID(DefaultCodespace, err.Error())
	}

	return nil
}
*/
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
			return ErrInvalidTokenSymbol(DefaultCodespace, fmt.Sprintf("[%s] is invalid token symbol", coin.Symbol))
		}
		if coin.TokenID == "" {
			return ErrInvalidTokenID(DefaultCodespace, "token id cannot be empty")
		}
		if err := linktype.ValidateTokenID(coin.TokenID); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, fmt.Sprintf("[%s] is invalid token id", coin.TokenID))
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
			return ErrInvalidTokenSymbol(DefaultCodespace, fmt.Sprintf("[%s] is invalid token symbol", coin.Symbol))
		}
		if coin.TokenID == "" {
			return ErrInvalidTokenID(DefaultCodespace, "token id cannot be empty")
		}
		if err := linktype.ValidateTokenID(coin.TokenID); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, fmt.Sprintf("[%s] is invalid token id", coin.TokenID))
		}

		if !coin.IsPositive() {
			return sdk.ErrInvalidCoins("amount is not valid")
		}
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From address cannot be empty")
	}
	return nil
}
