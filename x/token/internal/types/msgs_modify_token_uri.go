package types

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
)

var _ sdk.Msg = (*MsgModifyTokenURI)(nil)
var _ json.Marshaler = (*MsgModifyTokenURI)(nil)
var _ json.Unmarshaler = (*MsgModifyTokenURI)(nil)

const (
	ModifyActionName = "modify"
)

type MsgModifyTokenURI struct {
	Owner    sdk.AccAddress `json:"owner"`
	Symbol   string         `json:"symbol"`
	TokenURI string         `json:"token_uri"`
	TokenID  string         `json:"token_id"`
}

func (msg MsgModifyTokenURI) MarshalJSON() ([]byte, error) {
	type msgAlias MsgModifyTokenURI

	return json.Marshal((msgAlias)(msg))
}

func (msg *MsgModifyTokenURI) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgModifyTokenURI
	return json.Unmarshal(data, msgAlias(msg))
}

func NewMsgModifyTokenURI(owner sdk.AccAddress, symbol, tokenURI, tokenID string) MsgModifyTokenURI {
	return MsgModifyTokenURI{
		Owner:    owner,
		Symbol:   symbol,
		TokenURI: tokenURI,
		TokenID:  tokenID,
	}
}

func (msg MsgModifyTokenURI) Route() string {
	return RouterKey
}

func (msg MsgModifyTokenURI) Type() string {
	return ModifyActionName
}

func (msg MsgModifyTokenURI) ValidateBasic() sdk.Error {

	if msg.Symbol == "" {
		return sdk.ErrInvalidAddress("symbol cannot be empty")
	}

	if err := types.ValidateSymbol(msg.Symbol); err != nil {
		return sdk.ErrInvalidAddress(fmt.Sprintf("invalid symbol pattern found %s", msg.Symbol))
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}

	if msg.TokenID != "" {
		if err := types.ValidateSymbolTokenID(msg.TokenID); err != nil {
			return sdk.ErrInvalidAddress(fmt.Sprintf("invalid tokenId pattern found %s", msg.TokenID))
		}
	}
	return nil
}

func (msg MsgModifyTokenURI) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgModifyTokenURI) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
