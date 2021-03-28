package types

import (
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/contract"
)

var _ contract.Msg = (*MsgIssueFT)(nil)

type MsgIssueFT struct {
	Owner      sdk.AccAddress `json:"owner"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	Name       string         `json:"name"`
	Meta       string         `json:"meta"`
	Amount     sdk.Int        `json:"amount"`
	Mintable   bool           `json:"mintable"`
	Decimals   sdk.Int        `json:"decimals"`
}

func NewMsgIssueFT(owner, to sdk.AccAddress, contractID string, name, meta string, amount sdk.Int, decimal sdk.Int, mintable bool) MsgIssueFT {
	return MsgIssueFT{
		Owner:      owner,
		ContractID: contractID,
		To:         to,
		Name:       name,
		Meta:       meta,
		Amount:     amount,
		Mintable:   mintable,
		Decimals:   decimal,
	}
}

func (msg MsgIssueFT) Route() string                { return RouterKey }
func (msg MsgIssueFT) Type() string                 { return "issue_ft" }
func (msg MsgIssueFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }
func (msg MsgIssueFT) GetContractID() string        { return msg.ContractID }
func (msg MsgIssueFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIssueFT) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if len(msg.Name) == 0 {
		return ErrInvalidTokenName
	}
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner address cannot be empty")
	}
	if msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "to address cannot be empty")
	}

	if msg.Amount.Equal(sdk.NewInt(1)) && msg.Decimals.Equal(sdk.NewInt(0)) && !msg.Mintable {
		return ErrInvalidIssueFT
	}

	if msg.Decimals.GT(sdk.NewInt(18)) || msg.Decimals.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidTokenDecimals, "Decimals: %s", msg.Decimals)
	}

	if !ValidateName(msg.Name) {
		return sdkerrors.Wrapf(ErrInvalidNameLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", msg.Name, MaxTokenNameLength, utf8.RuneCountInString(msg.Name))
	}
	if !ValidateMeta(msg.Meta) {
		return sdkerrors.Wrapf(ErrInvalidMetaLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", msg.Meta, MaxTokenMetaLength, utf8.RuneCountInString(msg.Meta))
	}
	return nil
}

var _ contract.Msg = (*MsgIssueNFT)(nil)

type MsgIssueNFT struct {
	Owner      sdk.AccAddress `json:"owner"`
	ContractID string         `json:"contract_id"`
	Name       string         `json:"name"`
	Meta       string         `json:"meta"`
}

func NewMsgIssueNFT(owner sdk.AccAddress, contractID, name, meta string) MsgIssueNFT {
	return MsgIssueNFT{
		Owner:      owner,
		ContractID: contractID,
		Name:       name,
		Meta:       meta,
	}
}

func (MsgIssueNFT) Route() string                    { return RouterKey }
func (MsgIssueNFT) Type() string                     { return "issue_nft" }
func (msg MsgIssueNFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }
func (msg MsgIssueNFT) GetContractID() string        { return msg.ContractID }
func (msg MsgIssueNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIssueNFT) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}

	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner address cannot be empty")
	}
	if !ValidateName(msg.Name) {
		return sdkerrors.Wrapf(ErrInvalidNameLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", msg.Name, MaxTokenNameLength, utf8.RuneCountInString(msg.Name))
	}
	if !ValidateMeta(msg.Meta) {
		return sdkerrors.Wrapf(ErrInvalidMetaLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", msg.Meta, MaxTokenMetaLength, utf8.RuneCountInString(msg.Meta))
	}
	return nil
}
