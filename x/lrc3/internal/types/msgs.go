package types

import (
	"strings"

	nft "github.com/link-chain/link/x/nft"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

/* --------------------------------------------------------------------------- */
// MsgInit
/* --------------------------------------------------------------------------- */
// MsgInit defines a MsgInit message
type MsgInit struct {
	Denom        string
	OwnerAddress sdk.AccAddress
}

// NewMsgInit is a constructor function for MsgInitLRC3
func NewMsgInit(denom string, ownerAddress sdk.AccAddress) MsgInit {
	return MsgInit{
		Denom:        strings.TrimSpace(denom),
		OwnerAddress: ownerAddress,
	}
}

// Route Implements Msg
func (msg MsgInit) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgInit) Type() string { return "init" }

// ValidateBasic Implements Msg.
func (msg MsgInit) ValidateBasic() sdk.Error {
	if strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidDenom(DefaultCodespace)
	}
	if msg.OwnerAddress.Empty() {
		return sdk.ErrInvalidAddress("invalid sender address")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgInit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgInit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.OwnerAddress}
}

/* --------------------------------------------------------------------------- */
// MsgMintNFT
/* --------------------------------------------------------------------------- */

// MsgMintNFT defines a MintNFT message
type MsgMintNFT struct {
	Sender    sdk.AccAddress
	Recipient sdk.AccAddress
	Denom     string
	TokenURI  string
}

// NewMsgMintNFT is a constructor function for MsgMintNFT
func NewMsgMintNFT(sender, recipient sdk.AccAddress, denom, tokenURI string) MsgMintNFT {
	return MsgMintNFT{
		Sender:    sender,
		Recipient: recipient,
		Denom:     denom,
		TokenURI:  strings.TrimSpace(tokenURI),
	}
}

// Route Implements Msg
func (msg MsgMintNFT) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgMintNFT) Type() string { return "mint_nft" }

// ValidateBasic Implements Msg.
func (msg MsgMintNFT) ValidateBasic() sdk.Error {
	if strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidDenom(DefaultCodespace)
	}
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress("invalid froｍ address")
	}
	if msg.Recipient.Empty() {
		return sdk.ErrInvalidAddress("invalid to address")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgMintNFT) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgMintNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

/* --------------------------------------------------------------------------- */
// MsgBurn
/* --------------------------------------------------------------------------- */

// MsgBurn defines a BurnNFT message
type MsgBurn struct {
	MsgBurnNFT nft.MsgBurnNFT
	Executor   sdk.AccAddress
}

// NewMsgBurn is a constructor function for MsgBurnNFT
func NewMsgBurn(sender sdk.AccAddress, denom, tokenId string, executor sdk.AccAddress) MsgBurn {
	msgBurnNFT := nft.NewMsgBurnNFT(sender, tokenId, denom)
	return MsgBurn{
		MsgBurnNFT: msgBurnNFT,
		Executor:   executor,
	}
}

// Route Implements Msg
func (msg MsgBurn) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgBurn) Type() string { return "burn" }

// ValidateBasic Implements Msg.
func (msg MsgBurn) ValidateBasic() sdk.Error {
	if msg.MsgBurnNFT.Sender.Empty() {
		return sdk.ErrInvalidAddress("invalid sender address")
	}
	if strings.TrimSpace(msg.MsgBurnNFT.Denom) == "" {
		return ErrInvalidDenom(DefaultCodespace)
	}
	if strings.TrimSpace(msg.MsgBurnNFT.ID) == "" {
		return ErrInvalidTokenId(DefaultCodespace)
	}
	if msg.Executor.Empty() {
		return sdk.ErrInvalidAddress("invalid executor address")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBurn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgBurn) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Executor}
}

/* --------------------------------------------------------------------------- */
// MsgTransfer
/* --------------------------------------------------------------------------- */
// MsgTransfer defines a TransferFrom message
type MsgTransfer struct {
	MsgTransferNFT nft.MsgTransferNFT
	Executor       sdk.AccAddress
}

// NewMsgTransfer is a constructor function for MsgTransferFrom
func NewMsgTransfer(sender sdk.AccAddress, recipient sdk.AccAddress, denom string, tokenId string, executor sdk.AccAddress) MsgTransfer {
	msgTransferNFT := nft.NewMsgTransferNFT(sender, recipient, denom, tokenId)
	return MsgTransfer{
		MsgTransferNFT: msgTransferNFT,
		Executor:       executor,
	}
}

// Route Implements Msg
func (msg MsgTransfer) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgTransfer) Type() string { return "transfer" }

// ValidateBasic Implements Msg.
func (msg MsgTransfer) ValidateBasic() sdk.Error {
	if msg.MsgTransferNFT.Recipient.Empty() {
		return sdk.ErrInvalidAddress("invalid recipient address")
	}
	if msg.MsgTransferNFT.Sender.Empty() {
		return sdk.ErrInvalidAddress("invalid sender address")
	}
	if strings.TrimSpace(msg.MsgTransferNFT.Denom) == "" {
		return ErrInvalidDenom(DefaultCodespace)
	}
	if strings.TrimSpace(msg.MsgTransferNFT.ID) == "" {
		return ErrInvalidTokenId(DefaultCodespace)
	}
	if msg.Executor.Empty() {
		return sdk.ErrInvalidAddress("invalid executor address")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgTransfer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Executor}
}

/* --------------------------------------------------------------------------- */
// MsgEditMetadata
/* --------------------------------------------------------------------------- */

// MsgEditMetadata edits an NFT's metadata
type MsgEditMetadata struct {
	MsgEditNFTMetadata nft.MsgEditNFTMetadata
	Executor           sdk.AccAddress
}

// NewMsgEditMetadata is a constructor function for MsgSetName
func NewMsgEditMetadata(sender sdk.AccAddress, tokenId, denom, tokenURI string, executor sdk.AccAddress) MsgEditMetadata {
	msgEditNFTMetadata := nft.NewMsgEditNFTMetadata(sender, tokenId, denom, tokenURI)
	return MsgEditMetadata{
		MsgEditNFTMetadata: msgEditNFTMetadata,
		Executor:           executor,
	}
}

// Route Implements Msg
func (msg MsgEditMetadata) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgEditMetadata) Type() string { return "edit_metadata" }

// ValidateBasic Implements Msg.
func (msg MsgEditMetadata) ValidateBasic() sdk.Error {
	if msg.MsgEditNFTMetadata.Sender.Empty() {
		return sdk.ErrInvalidAddress("invalid sender address")
	}
	if strings.TrimSpace(msg.MsgEditNFTMetadata.ID) == "" {
		return ErrNotExistNFT(DefaultCodespace)
	}
	if strings.TrimSpace(msg.MsgEditNFTMetadata.Denom) == "" {
		return ErrNotExistLRC3(DefaultCodespace)
	}
	if msg.Executor.Empty() {
		return sdk.ErrInvalidAddress("invalid executor address")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgEditMetadata) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgEditMetadata) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Executor}
}

/* --------------------------------------------------------------------------- */
// MsgApprove
/* --------------------------------------------------------------------------- */
type MsgApprove struct {
	Sender    sdk.AccAddress
	Denom     string
	TokenID   string
	Recipient sdk.AccAddress
}

// NewMsgApprove is a constructor function for MsgApprove
func NewMsgApprove(sender sdk.AccAddress, denom string, tokenId string, recipient sdk.AccAddress) MsgApprove {
	return MsgApprove{
		Sender:    sender,
		Denom:     strings.TrimSpace(denom),
		TokenID:   tokenId,
		Recipient: recipient,
	}
}

// Route Implements Msg
func (msg MsgApprove) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgApprove) Type() string { return "approve" }

// ValidateBasic Implements Msg.
func (msg MsgApprove) ValidateBasic() sdk.Error {
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress("invalid from address")
	}
	if msg.Recipient.Empty() {
		return sdk.ErrInvalidAddress("invalid to address")
	}
	if strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidDenom(DefaultCodespace)
	}
	if strings.TrimSpace(msg.TokenID) == "" {
		return ErrInvalidTokenId(DefaultCodespace)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgApprove) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgApprove) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

/* --------------------------------------------------------------------------- */
// MsgSetApprovalForAll
/* --------------------------------------------------------------------------- */
type MsgSetApprovalForAll struct {
	Denom    string
	Owner    sdk.AccAddress
	Operator sdk.AccAddress
	Approved bool
}

// NewMsgSetApprovalForAll is a constructor function for MsgSetApprovalForAll
func NewMsgSetApprovalForAll(denom string, owner sdk.AccAddress, operator sdk.AccAddress, approved bool) MsgSetApprovalForAll {
	return MsgSetApprovalForAll{
		Denom:    strings.TrimSpace(denom),
		Owner:    owner,
		Operator: operator,
		Approved: approved,
	}
}

// Route Implements Msg
func (msg MsgSetApprovalForAll) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgSetApprovalForAll) Type() string { return "set_approval_for_all" }

// ValidateBasic Implements Msg.
func (msg MsgSetApprovalForAll) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("invalid owner address")
	}
	if msg.Operator.Empty() {
		return sdk.ErrInvalidAddress("invalid operator address")
	}
	if strings.TrimSpace(msg.Denom) == "" {
		return ErrInvalidDenom(DefaultCodespace)
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetApprovalForAll) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgSetApprovalForAll) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
