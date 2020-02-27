package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
)

var _ sdk.Msg = (*MsgIssueFT)(nil)

type MsgIssueFT struct {
	Owner    sdk.AccAddress `json:"owner"`
	Symbol   string         `json:"symbol"`
	Name     string         `json:"name"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
}

func NewMsgIssueFT(owner sdk.AccAddress, name, symbol string, amount sdk.Int, decimal sdk.Int, mintable bool) MsgIssueFT {
	return MsgIssueFT{
		Owner:    owner,
		Symbol:   symbol,
		Name:     name,
		Amount:   amount,
		Mintable: mintable,
		Decimals: decimal,
	}
}

func (msg MsgIssueFT) Route() string                { return RouterKey }
func (msg MsgIssueFT) Type() string                 { return "issue_ft" }
func (msg MsgIssueFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }
func (msg MsgIssueFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIssueFT) ValidateBasic() sdk.Error {
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

var _ sdk.Msg = (*MsgIssueNFT)(nil)

type MsgIssueNFT struct {
	Owner  sdk.AccAddress `json:"owner"`
	Symbol string         `json:"symbol"`
	Name   string         `json:"name"`
}

func NewMsgIssueNFT(owner sdk.AccAddress, symbol, name string) MsgIssueNFT {
	return MsgIssueNFT{
		Owner:  owner,
		Symbol: symbol,
		Name:   name,
	}
}

func (MsgIssueNFT) Route() string                    { return RouterKey }
func (MsgIssueNFT) Type() string                     { return "issue_nft" }
func (msg MsgIssueNFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }
func (msg MsgIssueNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIssueNFT) ValidateBasic() sdk.Error {
	if err := types.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}
	return nil
}
