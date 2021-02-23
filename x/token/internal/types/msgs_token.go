package types

import (
	"unicode/utf8"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/contract"
)

var _ sdk.Msg = (*MsgIssue)(nil)

type MsgIssue struct {
	Owner    sdk.AccAddress `json:"owner"`
	To       sdk.AccAddress `json:"to"`
	Name     string         `json:"name"`
	Symbol   string         `json:"symbol"`
	ImageURI string         `json:"img_uri"`
	Meta     string         `json:"meta"`
	Amount   sdk.Int        `json:"amount"`
	Mintable bool           `json:"mintable"`
	Decimals sdk.Int        `json:"decimals"`
}

func NewMsgIssue(owner, to sdk.AccAddress, name, symbol, meta string, imageURI string, amount sdk.Int, decimal sdk.Int, mintable bool) MsgIssue {
	return MsgIssue{
		Owner:    owner,
		To:       to,
		Name:     name,
		Symbol:   symbol,
		Meta:     meta,
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

func (msg MsgIssue) ValidateBasic() error {
	if len(msg.Name) == 0 {
		return ErrInvalidTokenName
	}
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner cannot be empty")
	}

	if msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "to cannot be empty")
	}

	if !ValidateName(msg.Name) {
		return sdkerrors.Wrapf(ErrInvalidNameLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", msg.Name, MaxTokenNameLength, utf8.RuneCountInString(msg.Name))
	}
	if !ValidateMeta(msg.Meta) {
		return sdkerrors.Wrapf(ErrInvalidMetaLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", msg.Meta, MaxTokenMetaLength, utf8.RuneCountInString(msg.Meta))
	}

	if err := ValidateTokenSymbol(msg.Symbol); err != nil {
		return sdkerrors.Wrapf(ErrInvalidTokenSymbol, "Symbol: %s", msg.Symbol)
	}

	if !ValidateImageURI(msg.ImageURI) {
		return sdkerrors.Wrapf(ErrInvalidImageURILength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", msg.ImageURI, MaxImageURILength, utf8.RuneCountInString(msg.ImageURI))
	}

	if msg.Decimals.GT(sdk.NewInt(18)) || msg.Decimals.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidTokenDecimals, "Decimals: %s", msg.Decimals)
	}

	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidAmount, msg.Amount.String())
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

func (msg MsgMint) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.Amount.IsNegative() {
		return sdkerrors.Wrap(ErrInvalidAmount, msg.Amount.String())
	}
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "from cannot be empty")
	}
	if msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "to cannot be empty")
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

func (msg MsgBurn) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidAmount, msg.Amount.String())
	}
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner cannot be empty")
	}
	if msg.Amount.IsNegative() {
		return sdkerrors.Wrap(ErrInvalidAmount, msg.Amount.String())
	}
	return nil
}

func (msg MsgBurn) GetContractID() string {
	return msg.ContractID
}

var _ contract.Msg = (*MsgBurnFrom)(nil)

type MsgBurnFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	Amount     sdk.Int        `json:"amount"`
}

func NewMsgBurnFrom(proxy sdk.AccAddress, contractID string, from sdk.AccAddress, amount sdk.Int) MsgBurnFrom {
	return MsgBurnFrom{
		Proxy:      proxy,
		ContractID: contractID,
		From:       from,
		Amount:     amount,
	}
}
func (MsgBurnFrom) Route() string                    { return RouterKey }
func (MsgBurnFrom) Type() string                     { return "burn_from" }
func (msg MsgBurnFrom) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Proxy} }
func (msg MsgBurnFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnFrom) ValidateBasic() error {
	if msg.Proxy.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Proxy cannot be empty")
	}
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if !msg.Amount.IsPositive() {
		return sdkerrors.Wrap(ErrInvalidAmount, msg.Amount.String())
	}
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty")
	}
	if msg.Proxy.Equals(msg.From) {
		return sdkerrors.Wrapf(ErrApproverProxySame, "Approver: %s", msg.Proxy.String())
	}

	return nil
}

func (msg MsgBurnFrom) GetContractID() string {
	return msg.ContractID
}
