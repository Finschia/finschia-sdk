package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
)

var _ sdk.Msg = (*MsgDeposit)(nil)

func NewMsgDeposit(rollupName string, sequencerAddress string, amount sdk.Coin) *MsgDeposit {
	return &MsgDeposit{
		RollupName:       rollupName,
		SequencerAddress: sequencerAddress,
		Value:            amount,
	}
}

// ValidateBasic implements Msg.
func (m MsgDeposit) ValidateBasic() error {

	return nil
}

// GetSigners implements Msg
func (m MsgDeposit) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.SequencerAddress)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgDeposit) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgDeposit) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgDeposit) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
