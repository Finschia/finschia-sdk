package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
	linktype "github.com/line/link/types"
)

var _ sdk.Msg = (*MsgCreateCollection)(nil)

type MsgCreateCollection struct {
	Owner  sdk.AccAddress `json:"owner"`
	Name   string         `json:"name"`
	Symbol string         `json:"symbol"`
}

func NewMsgCreateCollection(owner sdk.AccAddress, name, symbol string) MsgCreateCollection {
	return MsgCreateCollection{
		Owner:  owner,
		Name:   name,
		Symbol: symbol,
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
	Owner    sdk.AccAddress `json:"owner"`
	Symbol   string         `json:"symbol"`
	Name     string         `json:"name"`
	TokenURI string         `json:"token_uri"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
}

func NewMsgIssueCFT(owner sdk.AccAddress, name, symbol, tokenURI string, amount sdk.Int, decimal sdk.Int, mintable bool) MsgIssueCFT {
	return MsgIssueCFT{
		Owner:    owner,
		Symbol:   symbol,
		Name:     name,
		TokenURI: tokenURI,
		Amount:   amount,
		Mintable: mintable,
		Decimals: decimal,
	}
}

func (msg MsgIssueCFT) Route() string                { return RouterKey }
func (msg MsgIssueCFT) Type() string                 { return "issue_cft" }
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
	Owner  sdk.AccAddress `json:"owner"`
	Symbol string         `json:"symbol"`
}

func NewMsgIssueCNFT(owner sdk.AccAddress, symbol string) MsgIssueCNFT {
	return MsgIssueCNFT{
		Owner:  owner,
		Symbol: symbol,
	}
}

func (MsgIssueCNFT) Route() string                    { return RouterKey }
func (MsgIssueCNFT) Type() string                     { return "issue_cnft" }
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

var _ sdk.Msg = (*MsgApproveCollection)(nil)

type MsgApproveCollection struct {
	Approver sdk.AccAddress `json:"approver"`
	Proxy    sdk.AccAddress `json:"proxy"`
	Symbol   string         `json:"symbol"`
}

func NewMsgApproveCollection(approver sdk.AccAddress, proxy sdk.AccAddress, symbol string) MsgApproveCollection {
	return MsgApproveCollection{
		Approver: approver,
		Proxy:    proxy,
		Symbol:   symbol,
	}
}

func (msg MsgApproveCollection) ValidateBasic() sdk.Error {
	if msg.Approver.Empty() {
		return sdk.ErrInvalidAddress("Approver cannot be empty")
	}
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.Approver.Equals(msg.Proxy) {
		return ErrApproverProxySame(DefaultCodespace, msg.Approver.String())
	}
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	return nil
}

func (MsgApproveCollection) Route() string { return RouterKey }
func (MsgApproveCollection) Type() string  { return "approve_collection" }
func (msg MsgApproveCollection) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Approver}
}
func (msg MsgApproveCollection) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

var _ sdk.Msg = (*MsgDisapproveCollection)(nil)

type MsgDisapproveCollection struct {
	Approver sdk.AccAddress `json:"approver"`
	Proxy    sdk.AccAddress `json:"proxy"`
	Symbol   string         `json:"symbol"`
}

func NewMsgDisapproveCollection(approver sdk.AccAddress, proxy sdk.AccAddress, symbol string) MsgDisapproveCollection {
	return MsgDisapproveCollection{
		Approver: approver,
		Proxy:    proxy,
		Symbol:   symbol,
	}
}

func (msg MsgDisapproveCollection) ValidateBasic() sdk.Error {
	if msg.Approver.Empty() {
		return sdk.ErrInvalidAddress("Approver cannot be empty")
	}
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.Approver.Equals(msg.Proxy) {
		return ErrApproverProxySame(DefaultCodespace, msg.Approver.String())
	}
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	return nil
}

func (MsgDisapproveCollection) Route() string { return RouterKey }
func (MsgDisapproveCollection) Type() string  { return "disapprove_collection" }
func (msg MsgDisapproveCollection) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Approver}
}
func (msg MsgDisapproveCollection) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
