package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
	"github.com/line/link/x/contract"
)

var _ sdk.Msg = (*MsgIssue)(nil)

type MsgIssue struct {
	Owner    sdk.AccAddress `json:"owner"`
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	ImageURI string         `json:"image_uri"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
}

func NewMsgIssue(owner sdk.AccAddress, name, symbol, imageURI string, amount sdk.Int, decimal sdk.Int, mintable bool) MsgIssue {
	return MsgIssue{
		Owner:    owner,
		Name:     name,
		Symbol:   symbol,
		ImageURI: imageURI,
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
	if len(msg.Name) == 0 {
		return ErrInvalidTokenName(DefaultCodespace, msg.Name)
	}
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner cannot be empty")
	}

	if !ValidateName(msg.Name) {
		return ErrInvalidNameLength(DefaultCodespace, msg.Name)
	}

	if err := types.ValidateTokenSymbol(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, msg.Symbol)
	}

	if !ValidateImageURI(msg.ImageURI) {
		return ErrInvalidImageURILength(DefaultCodespace, msg.ImageURI)
	}

	if msg.Decimals.GT(sdk.NewInt(18)) || msg.Decimals.IsNegative() {
		return ErrInvalidTokenDecimals(DefaultCodespace, msg.Decimals)
	}

	if msg.Amount.IsNegative() {
		return ErrInvalidAmount(DefaultCodespace, msg.Amount.String())
	}

	return nil
}

var _ contract.Msg = (*MsgMint)(nil)

type MsgMint struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	Amount     sdk.Int        `json:"amount"`
}

func NewMsgMint(from sdk.AccAddress, contractID string, to sdk.AccAddress, amount sdk.Int) MsgMint {
	return MsgMint{
		From:       from,
		ContractID: contractID,
		To:         to,
		Amount:     amount,
	}
}
func (MsgMint) Route() string                    { return RouterKey }
func (MsgMint) Type() string                     { return "mint" }
func (msg MsgMint) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMint) ValidateBasic() sdk.Error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.Amount.IsNegative() {
		return ErrInvalidAmount(DefaultCodespace, msg.Amount.String())
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("from cannot be empty")
	}
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("to cannot be empty")
	}
	return nil
}
func (msg MsgMint) GetContractID() string {
	return msg.ContractID
}

var _ contract.Msg = (*MsgBurn)(nil)

type MsgBurn struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	Amount     sdk.Int        `json:"amount"`
}

func NewMsgBurn(from sdk.AccAddress, contractID string, amount sdk.Int) MsgBurn {
	return MsgBurn{
		From:       from,
		ContractID: contractID,
		Amount:     amount,
	}
}
func (MsgBurn) Route() string                    { return RouterKey }
func (MsgBurn) Type() string                     { return "burn" }
func (msg MsgBurn) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgBurn) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurn) ValidateBasic() sdk.Error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("owner cannot be empty")
	}
	if msg.Amount.IsNegative() {
		return ErrInvalidAmount(DefaultCodespace, msg.Amount.String())
	}
	return nil
}

func (msg MsgBurn) GetContractID() string {
	return msg.ContractID
}
