package types

import (
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/contract"
)

var _ contract.Msg = (*MsgMintNFT)(nil)

type MintNFTParam struct {
	Name      string `json:"name"`
	Meta      string `json:"meta"`
	TokenType string `json:"token_type"`
}

func NewMintNFTParam(name, meta, tokenType string) MintNFTParam {
	return MintNFTParam{
		Name:      name,
		Meta:      meta,
		TokenType: tokenType,
	}
}

type MsgMintNFT struct {
	From          sdk.AccAddress `json:"from"`
	ContractID    string         `json:"contract_id"`
	To            sdk.AccAddress `json:"to"`
	MintNFTParams []MintNFTParam `json:"params"`
}

func NewMsgMintNFT(from sdk.AccAddress, contractID string, to sdk.AccAddress, mintNFTParams ...MintNFTParam) MsgMintNFT {
	return MsgMintNFT{
		From:          from,
		ContractID:    contractID,
		To:            to,
		MintNFTParams: mintNFTParams,
	}
}

func (msg MsgMintNFT) Route() string                { return RouterKey }
func (msg MsgMintNFT) Type() string                 { return "mint_nft" }
func (msg MsgMintNFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgMintNFT) GetContractID() string        { return msg.ContractID }
func (msg MsgMintNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMintNFT) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "from address cannot be empty")
	}
	if msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "to address cannot be empty")
	}

	if len(msg.MintNFTParams) == 0 {
		return sdkerrors.Wrap(ErrEmptyField, "params cannot be empty")
	}
	for _, mintNFTParam := range msg.MintNFTParams {
		if len(mintNFTParam.Name) == 0 {
			return ErrInvalidTokenName
		}
		if !ValidateName(mintNFTParam.Name) {
			return sdkerrors.Wrapf(ErrInvalidNameLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", mintNFTParam.Name, MaxTokenNameLength, utf8.RuneCountInString(mintNFTParam.Name))
		}
		if !ValidateMeta(mintNFTParam.Meta) {
			return sdkerrors.Wrapf(ErrInvalidMetaLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", mintNFTParam.Meta, MaxTokenMetaLength, utf8.RuneCountInString(mintNFTParam.Meta))
		}
		if err := ValidateTokenTypeNFT(mintNFTParam.TokenType); err != nil {
			return sdkerrors.Wrap(ErrInvalidTokenID, err.Error())
		}
	}

	return nil
}

var _ contract.Msg = (*MsgBurnNFT)(nil)

type MsgBurnNFT struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	TokenIDs   []string       `json:"token_ids"`
}

func NewMsgBurnNFT(from sdk.AccAddress, contractID string, tokenIDs ...string) MsgBurnNFT {
	return MsgBurnNFT{
		From:       from,
		ContractID: contractID,
		TokenIDs:   tokenIDs,
	}
}

func (msg MsgBurnNFT) Route() string                { return RouterKey }
func (msg MsgBurnNFT) Type() string                 { return "burn_nft" }
func (msg MsgBurnNFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgBurnNFT) GetContractID() string        { return msg.ContractID }
func (msg MsgBurnNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnNFT) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner address cannot be empty")
	}

	if len(msg.TokenIDs) == 0 {
		return sdkerrors.Wrap(ErrEmptyField, "token_ids cannot be empty")
	}
	for _, tokenID := range msg.TokenIDs {
		if err := ValidateTokenID(tokenID); err != nil {
			return sdkerrors.Wrap(ErrInvalidTokenID, err.Error())
		}
		if err := ValidateTokenTypeNFT(tokenID[:TokenTypeLength]); err != nil {
			return sdkerrors.Wrap(ErrInvalidTokenID, err.Error())
		}
	}

	return nil
}

var _ contract.Msg = (*MsgBurnNFTFrom)(nil)

type MsgBurnNFTFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	TokenIDs   []string       `json:"token_ids"`
}

func NewMsgBurnNFTFrom(proxy sdk.AccAddress, contractID string, from sdk.AccAddress, tokenIDs ...string) MsgBurnNFTFrom {
	return MsgBurnNFTFrom{
		Proxy:      proxy,
		ContractID: contractID,
		From:       from,
		TokenIDs:   tokenIDs,
	}
}

func (msg MsgBurnNFTFrom) Route() string                { return RouterKey }
func (msg MsgBurnNFTFrom) Type() string                 { return "burn_nft_from" }
func (msg MsgBurnNFTFrom) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Proxy} }
func (msg MsgBurnNFTFrom) GetContractID() string        { return msg.ContractID }
func (msg MsgBurnNFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnNFTFrom) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.Proxy.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty")
	}
	if msg.Proxy.Equals(msg.From) {
		return sdkerrors.Wrapf(ErrApproverProxySame, "Approver: %s", msg.Proxy.String())
	}

	if len(msg.TokenIDs) == 0 {
		return sdkerrors.Wrap(ErrEmptyField, "token_ids cannot be empty")
	}
	for _, tokenID := range msg.TokenIDs {
		if err := ValidateTokenID(tokenID); err != nil {
			return sdkerrors.Wrap(ErrInvalidTokenID, err.Error())
		}
		if err := ValidateTokenTypeNFT(tokenID[:TokenTypeLength]); err != nil {
			return sdkerrors.Wrap(ErrInvalidTokenID, err.Error())
		}
	}
	return nil
}

var _ contract.Msg = (*MsgMintFT)(nil)

type MsgMintFT struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	Amount     Coins          `json:"amount"`
}

func NewMsgMintFT(from sdk.AccAddress, contractID string, to sdk.AccAddress, amount ...Coin) MsgMintFT {
	return MsgMintFT{
		From:       from,
		ContractID: contractID,
		To:         to,
		Amount:     amount,
	}
}
func (MsgMintFT) Route() string                    { return RouterKey }
func (MsgMintFT) Type() string                     { return "mint_ft" }
func (msg MsgMintFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgMintFT) GetContractID() string        { return msg.ContractID }
func (msg MsgMintFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgMintFT) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}

	for _, tokenID := range msg.Amount {
		if err := ValidateDenom(tokenID.Denom); err != nil {
			return sdkerrors.Wrap(ErrInvalidTokenID, "invalid token id")
		}
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAmount, msg.Amount.String())
	}

	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "from address cannot be empty")
	}
	if msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "to address cannot be empty")
	}

	return nil
}

var _ contract.Msg = (*MsgBurnFT)(nil)

type MsgBurnFT struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	Amount     Coins          `json:"amount"`
}

func NewMsgBurnFT(from sdk.AccAddress, contractID string, amount ...Coin) MsgBurnFT {
	return MsgBurnFT{
		From:       from,
		ContractID: contractID,
		Amount:     amount,
	}
}
func (MsgBurnFT) Route() string                    { return RouterKey }
func (MsgBurnFT) Type() string                     { return "burn_ft" }
func (msg MsgBurnFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgBurnFT) GetContractID() string        { return msg.ContractID }
func (msg MsgBurnFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnFT) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}

	for _, tokenID := range msg.Amount {
		if err := ValidateDenom(tokenID.Denom); err != nil {
			return sdkerrors.Wrap(ErrInvalidTokenID, "invalid token id")
		}
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAmount, msg.Amount.String())
	}

	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "from address cannot be empty")
	}
	return nil
}

var _ contract.Msg = (*MsgBurnFTFrom)(nil)

type MsgBurnFTFrom struct {
	Proxy      sdk.AccAddress `json:"proxy"`
	ContractID string         `json:"contract_id"`
	From       sdk.AccAddress `json:"from"`
	Amount     Coins          `json:"amount"`
}

func NewMsgBurnFTFrom(proxy sdk.AccAddress, contractID string, from sdk.AccAddress, amount ...Coin) MsgBurnFTFrom {
	return MsgBurnFTFrom{
		Proxy:      proxy,
		ContractID: contractID,
		From:       from,
		Amount:     amount,
	}
}

func (MsgBurnFTFrom) Route() string                    { return RouterKey }
func (MsgBurnFTFrom) Type() string                     { return "burn_ft_from" }
func (msg MsgBurnFTFrom) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Proxy} }
func (msg MsgBurnFTFrom) GetContractID() string        { return msg.ContractID }
func (msg MsgBurnFTFrom) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgBurnFTFrom) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.Proxy.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty")
	}
	if msg.Proxy.Equals(msg.From) {
		return sdkerrors.Wrapf(ErrApproverProxySame, "Approver: %s", msg.Proxy.String())
	}
	for _, tokenID := range msg.Amount {
		if err := ValidateDenom(tokenID.Denom); err != nil {
			return sdkerrors.Wrap(ErrInvalidTokenID, "invalid token id")
		}
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(ErrInvalidAmount, msg.Amount.String())
	}
	return nil
}
