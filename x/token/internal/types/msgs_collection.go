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

var _ sdk.Msg = (*MsgIssueCollection)(nil)

type MsgIssueCollection struct {
	Symbol   string         `json:"symbol"`
	Name     string         `json:"name"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
	TokenID  string         `json:"token_id"`
}

func NewMsgIssueCollection(name, symbol, tokenURI string, owner sdk.AccAddress, amount sdk.Int, decimal sdk.Int, mintable bool, tokenID string) MsgIssueCollection {
	return MsgIssueCollection{
		Symbol:   symbol,
		Name:     name,
		Owner:    owner,
		TokenURI: tokenURI,
		Amount:   amount,
		Mintable: mintable,
		Decimals: decimal,
		TokenID:  tokenID,
	}
}

func (msg MsgIssueCollection) Route() string                { return RouterKey }
func (msg MsgIssueCollection) Type() string                 { return "issue_token_collection" }
func (msg MsgIssueCollection) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }
func (msg MsgIssueCollection) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIssueCollection) ValidateBasic() sdk.Error {
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

	if msg.TokenID != "" {
		if err := types.ValidateSymbolTokenID(msg.TokenID); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
		if msg.TokenID[0] != '0' {
			return ErrInvalidTokenID(DefaultCodespace, "token-id for FT should be started with 0")
		}
	}

	return nil
}

var _ sdk.Msg = (*MsgIssueNFTCollection)(nil)

type MsgIssueNFTCollection struct {
	Symbol   string         `json:"symbol"`
	Name     string         `json:"name"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
	TokenID  string         `json:"token_id"`
}

func NewMsgIssueNFTCollection(name, symbol, tokenURI string, owner sdk.AccAddress, tokenID string) MsgIssueNFTCollection {
	return MsgIssueNFTCollection{
		Symbol:   symbol,
		Name:     name,
		Owner:    owner,
		TokenURI: tokenURI,
		TokenID:  tokenID,
	}
}

func (msg MsgIssueNFTCollection) Route() string                { return RouterKey }
func (msg MsgIssueNFTCollection) Type() string                 { return "issue_nft_collection" }
func (msg MsgIssueNFTCollection) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }
func (msg MsgIssueNFTCollection) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIssueNFTCollection) ValidateBasic() sdk.Error {
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}
	if len(msg.Name) == 0 {
		return ErrInvalidTokenName(DefaultCodespace, msg.Name)
	}
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}

	if msg.TokenID != "" {
		if err := types.ValidateSymbolTokenID(msg.TokenID); err != nil {
			if err := types.ValidateSymbolBaseTokenID(msg.TokenID); err != nil {
				return ErrInvalidTokenID(DefaultCodespace, err.Error())
			}
		}

		if msg.TokenID[0] == '0' {
			return ErrInvalidTokenID(DefaultCodespace, "token-id for NFT should not be started with 0")
		}
	}

	return nil
}

var _ sdk.Msg = (*MsgMintCollection)(nil)

type MsgMintCollection struct {
	From   sdk.AccAddress            `json:"from"`
	To     sdk.AccAddress            `json:"to"`
	Amount linktype.CoinWithTokenIDs `json:"amount"`
}

func NewMsgCollectionTokenMint(from, to sdk.AccAddress, amount linktype.CoinWithTokenIDs) MsgMintCollection {
	return MsgMintCollection{
		From:   from,
		To:     to,
		Amount: amount,
	}
}
func (MsgMintCollection) Route() string                    { return RouterKey }
func (MsgMintCollection) Type() string                     { return "mint_token" }
func (msg MsgMintCollection) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgMintCollection) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMintCollection) ValidateBasic() sdk.Error {
	for _, coin := range msg.Amount {
		if err := linktype.ValidateSymbolUserDefined(coin.Symbol); err != nil {
			return ErrInvalidTokenSymbol(DefaultCodespace, fmt.Sprintf("[%s] is invalid token symbol", coin.Symbol))
		}
		if coin.TokenID == "" {
			return ErrInvalidTokenID(DefaultCodespace, "token id cannot be empty")
		}
		if err := linktype.ValidateSymbolTokenID(coin.TokenID); err != nil {
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

var _ sdk.Msg = (*MsgBurnCollection)(nil)

type MsgBurnCollection struct {
	From   sdk.AccAddress            `json:"from"`
	Amount linktype.CoinWithTokenIDs `json:"amount"`
}

func NewMsgCollectionTokenBurn(from sdk.AccAddress, amount linktype.CoinWithTokenIDs) MsgBurnCollection {
	return MsgBurnCollection{
		From:   from,
		Amount: amount,
	}
}
func (MsgBurnCollection) Route() string                    { return RouterKey }
func (MsgBurnCollection) Type() string                     { return "burn_token" }
func (msg MsgBurnCollection) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgBurnCollection) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnCollection) ValidateBasic() sdk.Error {
	for _, coin := range msg.Amount {
		if err := linktype.ValidateSymbolUserDefined(coin.Symbol); err != nil {
			return ErrInvalidTokenSymbol(DefaultCodespace, fmt.Sprintf("[%s] is invalid token symbol", coin.Symbol))
		}
		if coin.TokenID == "" {
			return ErrInvalidTokenID(DefaultCodespace, "token id cannot be empty")
		}
		if err := linktype.ValidateSymbolTokenID(coin.TokenID); err != nil {
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
