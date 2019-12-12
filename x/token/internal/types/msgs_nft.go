package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/link-chain/link/types"
)

var _ sdk.Msg = (*MsgIssueNFT)(nil)
var _ sdk.Msg = (*MsgIssueNFTCollection)(nil)

type MsgIssueNFT struct {
	MsgIssueCommon
}

func NewMsgIssueNFT(name, symbol, tokenURI string, owner sdk.AccAddress) MsgIssueNFT {
	return MsgIssueNFT{
		MsgIssueCommon: NewMsgIssueCommon(name, symbol, tokenURI, owner),
	}
}

// ATTENTION: to avoid embedded serialization implement this custom JSON marshaler
// see https://github.com/tendermint/go-amino/issues/269
func (msg MsgIssueNFT) MarshalJSON() ([]byte, error) {
	type msgAlias MsgIssueNFT
	return json.Marshal((msgAlias)(msg))
}

func (msg MsgIssueNFT) Route() string { return RouterKey }

func (msg MsgIssueNFT) Type() string { return "issue_nft" }

func (msg MsgIssueNFT) ValidateBasic() sdk.Error {
	if err := msg.MsgIssueCommon.ValidateBasic(); err != nil {
		return err
	}
	return nil
}

func (msg MsgIssueNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIssueNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type MsgIssueNFTCollection struct {
	MsgIssueNFT
	MsgCollection
}

func NewMsgIssueNFTCollection(name, symbol, tokenURI string, owner sdk.AccAddress, tokenID string) MsgIssueNFTCollection {
	return MsgIssueNFTCollection{
		MsgIssueNFT:   NewMsgIssueNFT(name, symbol, tokenURI, owner),
		MsgCollection: NewMsgCollection(tokenID),
	}
}

// ATTENTION: to avoid embedded serialization implement this custom JSON marshaler
// see https://github.com/tendermint/go-amino/issues/269
func (msg MsgIssueNFTCollection) MarshalJSON() ([]byte, error) {
	type msgAlias MsgIssueNFTCollection
	return json.Marshal((msgAlias)(msg))
}

func (msg MsgIssueNFTCollection) Type() string { return "issue_nft_collection" }

func (msg MsgIssueNFTCollection) ValidateBasic() sdk.Error {
	if err := msg.MsgIssueNFT.ValidateBasic(); err != nil {
		return err
	}

	if err := msg.MsgCollection.ValidateBasic(); err != nil {
		return err
	}

	coin := sdk.NewCoin(types.SymbolCollectionToken(msg.Symbol, msg.TokenID), sdk.NewInt(1))
	if !coin.IsValid() {
		return sdk.ErrInvalidCoins(coin.String())
	}

	return nil
}
