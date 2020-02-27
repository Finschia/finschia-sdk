package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

var _ sdk.Msg = (*MsgModify)(nil)

type MsgModify struct {
	Owner      sdk.AccAddress   `json:"owner"`
	Symbol     string           `json:"symbol"`
	TokenType  string           `json:"token_type"`
	TokenIndex string           `json:"token_index"`
	Changes    linktype.Changes `json:"changes"`
}

func NewMsgModify(owner sdk.AccAddress, symbol, tokenType, tokenIndex string, changes linktype.Changes) MsgModify {
	return MsgModify{
		Owner:      owner,
		Symbol:     symbol,
		TokenType:  tokenType,
		TokenIndex: tokenIndex,
		Changes:    changes,
	}
}

func (msg MsgModify) Route() string { return RouterKey }
func (msg MsgModify) Type() string  { return "modify_token" }
func (msg MsgModify) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgModify) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }

func (msg MsgModify) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}

	if err := linktype.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, msg.Symbol)
	}
	if msg.TokenType != "" && linktype.ValidateTokenTypeNFT(msg.TokenType) != nil {
		return ErrInvalidTokenType(DefaultCodespace, msg.TokenType)
	}
	if msg.TokenIndex != "" && linktype.ValidateTokenIndex(msg.TokenIndex) != nil {
		return ErrInvalidTokenIndex(DefaultCodespace, msg.TokenIndex)
	}

	validator := NewChangesValidator()
	if err := validator.SetMode(msg.TokenType, msg.TokenIndex); err != nil {
		return err
	}
	if err := validator.Validate(msg.Changes); err != nil {
		return err
	}

	return nil
}
