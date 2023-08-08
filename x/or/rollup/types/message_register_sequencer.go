package types

import (
	sdk "github.com/Finschia/finschia-rdk/types"
	sdkerrors "github.com/Finschia/finschia-rdk/types/errors"

	codectypes "github.com/Finschia/finschia-rdk/codec/types"
	cryptotypes "github.com/Finschia/finschia-rdk/crypto/types"
)

const (
	// TODO: change to parameter
	MinDepositAmount = 10
)

var _ sdk.Msg = (*MsgRegisterSequencer)(nil)

func NewMsgRegisterSequencer(pubkey cryptotypes.PubKey, rollupName, creator string, amount sdk.Coin) (*MsgRegisterSequencer, error) {
	var pkAny *codectypes.Any
	if pubkey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubkey); err != nil {
			return nil, err
		}
	}
	return &MsgRegisterSequencer{
		Creator:    creator,
		Pubkey:     pkAny,
		RollupName: rollupName,
		Value:      amount,
	}, nil
}

// ValidateBasic implements Msg.
func (m MsgRegisterSequencer) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.Creator)
	}

	if err := validateName(m.RollupName); err != nil {
		return err
	}

	if m.Value.Amount.LT(sdk.NewInt(MinDepositAmount)) {
		return ErrIsNotEnoughDepositAmount
	}

	return nil
}

// GetSigners implements Msg
func (m MsgRegisterSequencer) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Creator)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgRegisterSequencer) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgRegisterSequencer) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgRegisterSequencer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
