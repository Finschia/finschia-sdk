package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
)

var _ sdk.Msg = (*MsgWithdraw)(nil)

func NewMsgWithdraw(rollupName string, sequencerAddress string, amount sdk.Coin) *MsgWithdraw {
	return &MsgWithdraw{
		RollupName:       rollupName,
		SequencerAddress: sequencerAddress,
		Value:            amount,
	}
}

// ValidateBasic implements Msg.
func (m MsgWithdraw) ValidateBasic() error {

	return nil
}

// GetSigners implements Msg
func (m MsgWithdraw) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.SequencerAddress)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgWithdraw) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgWithdraw) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgWithdraw) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
