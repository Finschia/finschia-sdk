package types

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/contract"
)

var _ contract.Msg = (*MsgModify)(nil)

type MsgModify struct {
	Owner      sdk.AccAddress `json:"owner"`
	ContractID string         `json:"contract_id"`
	TokenType  string         `json:"token_type"`
	TokenIndex string         `json:"token_index"`
	Changes    Changes        `json:"changes"`
}

func NewMsgModify(owner sdk.AccAddress, contractID, tokenType, tokenIndex string, changes Changes) MsgModify {
	return MsgModify{
		Owner:      owner,
		ContractID: contractID,
		TokenType:  tokenType,
		TokenIndex: tokenIndex,
		Changes:    changes,
	}
}

func (msg MsgModify) Route() string         { return RouterKey }
func (msg MsgModify) Type() string          { return "modify_token" }
func (msg MsgModify) GetContractID() string { return msg.ContractID }
func (msg MsgModify) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
func (msg MsgModify) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }

func (msg MsgModify) ValidateBasic() error {
	if err := contract.ValidateContractIDBasic(msg); err != nil {
		return err
	}

	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner address cannot be empty")
	}

	if msg.TokenType != "" {
		if err := ValidateTokenType(msg.TokenType); err != nil {
			return sdkerrors.Wrap(ErrInvalidTokenType, msg.TokenType)
		}
		if ValidateTokenTypeFT(msg.TokenType) == nil && msg.TokenIndex == "" {
			return sdkerrors.Wrap(ErrTokenTypeFTWithoutIndex, msg.TokenType)
		}
	}
	if msg.TokenIndex != "" && ValidateTokenIndex(msg.TokenIndex) != nil {
		return sdkerrors.Wrap(ErrInvalidTokenIndex, msg.TokenIndex)
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
