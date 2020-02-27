package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

var _ sdk.Msg = (*MsgModify)(nil)

type MsgModify struct {
	Owner   sdk.AccAddress   `json:"owner"`
	Symbol  string           `json:"symbol"`
	Changes linktype.Changes `json:"changes"`
}

func NewMsgModify(owner sdk.AccAddress, symbol string, changes linktype.Changes) MsgModify {
	return MsgModify{
		Owner:   owner,
		Symbol:  symbol,
		Changes: changes,
	}
}

func (msg MsgModify) Route() string { return RouterKey }
func (msg MsgModify) Type() string  { return "modify_token" }
func (msg MsgModify) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgModify) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }

func (msg MsgModify) ValidateBasic() sdk.Error {
	if err := linktype.ValidateSymbolUserDefined(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, msg.Symbol)
	}

	validator := NewChangesValidator()
	if err := validator.Validate(msg.Changes); err != nil {
		return err
	}

	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress("owner address cannot be empty")
	}

	return nil
}
