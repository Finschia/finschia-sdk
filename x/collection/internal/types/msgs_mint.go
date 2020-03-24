package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
	"github.com/line/link/x/contract"
)

var _ contract.Msg = (*MsgMintNFT)(nil)

type MsgMintNFT struct {
	From       sdk.AccAddress `json:"from"`
	ContractID string         `json:"contract_id"`
	To         sdk.AccAddress `json:"to"`
	Name       string         `json:"name"`
	Meta       string         `json:"meta"`
	TokenType  string         `json:"token_type"`
}

func NewMsgMintNFT(from sdk.AccAddress, contractID string, to sdk.AccAddress, name, meta, tokenType string) MsgMintNFT {
	return MsgMintNFT{
		From:       from,
		ContractID: contractID,
		To:         to,
		Name:       name,
		Meta:       meta,
		TokenType:  tokenType,
	}
}

func (msg MsgMintNFT) Route() string                { return RouterKey }
func (msg MsgMintNFT) Type() string                 { return "mint_nft" }
func (msg MsgMintNFT) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MsgMintNFT) GetContractID() string        { return msg.ContractID }
func (msg MsgMintNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgMintNFT) ValidateBasic() sdk.Error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
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
	if !ValidateName(msg.Name) {
		return ErrInvalidNameLength(DefaultCodespace, msg.Name)
	}
	if !ValidateMeta(msg.Meta) {
		return ErrInvalidMetaLength(DefaultCodespace, msg.Meta)
	}

	if err := types.ValidateTokenTypeNFT(msg.TokenType); err != nil {
		return ErrInvalidTokenID(DefaultCodespace, err.Error())
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

func (msg MsgBurnNFT) ValidateBasic() sdk.Error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}

	for _, tokenID := range msg.TokenIDs {
		if err := types.ValidateTokenID(tokenID); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
		if err := types.ValidateTokenTypeNFT(tokenID[:TokenTypeLength]); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
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

func (msg MsgBurnNFTFrom) ValidateBasic() sdk.Error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}
	if msg.Proxy.Equals(msg.From) {
		return ErrApproverProxySame(DefaultCodespace, msg.Proxy.String())
	}

	for _, tokenID := range msg.TokenIDs {
		if err := types.ValidateTokenID(tokenID); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
		}
		if err := types.ValidateTokenTypeNFT(tokenID[:TokenTypeLength]); err != nil {
			return ErrInvalidTokenID(DefaultCodespace, err.Error())
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

func (msg MsgMintFT) ValidateBasic() sdk.Error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}

	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("from address cannot be empty")
	}
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress("to address cannot be empty")
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

func (msg MsgBurnFT) ValidateBasic() sdk.Error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}

	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}

	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From address cannot be empty")
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

func (msg MsgBurnFTFrom) ValidateBasic() sdk.Error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}
	if msg.Proxy.Empty() {
		return sdk.ErrInvalidAddress("Proxy cannot be empty")
	}
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("From cannot be empty")
	}
	if msg.Proxy.Equals(msg.From) {
		return ErrApproverProxySame(DefaultCodespace, msg.Proxy.String())
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}
	return nil
}
