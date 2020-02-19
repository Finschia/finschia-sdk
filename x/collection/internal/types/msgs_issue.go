package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
)

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
