package types

import (
	"unicode/utf8"

	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

const (
	nameLengthLimit = 30
)

var _ sdk.Msg = (*MsgCreateRollup)(nil)

func NewMsgCreateRollup(rollupName, creator string, l1ToL2GasRatio uint64, permissionedAddresses *Sequencers) *MsgCreateRollup {
	return &MsgCreateRollup{
		RollupName:            rollupName,
		Creator:               creator,
		L1ToL2GasRatio:        l1ToL2GasRatio,
		PermissionedAddresses: *permissionedAddresses,
	}
}

// ValidateBasic implements Msg.
func (m MsgCreateRollup) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Creator); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", m.Creator)
	}

	if err := validateName(m.RollupName); err != nil {
		return err
	}

	permissionedAddresses := m.GetPermissionedAddresses()
	if permissionedAddresses.Size() > 0 {
		duplicateAddresses := make(map[string]bool)
		for _, item := range permissionedAddresses.GetAddresses() {
			// check if the item/element exist in the duplicateAddresses map
			_, exist := duplicateAddresses[item]
			if exist {
				return sdkerrors.Wrapf(ErrPermissionedAddressesDuplicate, "address: %s", item)
			}
			// check Bech32 format
			if _, err := sdk.AccAddressFromBech32(item); err != nil {
				return sdkerrors.Wrapf(ErrInvalidPermissionedAddress, "invalid permissioned address: %s", err)
			}
			// mark as exist
			duplicateAddresses[item] = true
		}
	}

	return nil
}

// GetSigners implements Msg
func (m MsgCreateRollup) GetSigners() []sdk.AccAddress {
	signer, _ := sdk.AccAddressFromBech32(m.Creator)
	return []sdk.AccAddress{signer}
}

// Type implements the LegacyMsg.Type method.
func (m MsgCreateRollup) Type() string {
	return sdk.MsgTypeURL(&m)
}

// Route implements the LegacyMsg.Route method.
func (m MsgCreateRollup) Route() string {
	return RouterKey
}

// GetSignBytes implements the LegacyMsg.GetSignBytes method.
func (m MsgCreateRollup) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func validateName(name string) error {
	if err := validateStringSize(name, nameLengthLimit, "name"); err != nil {
		return ErrInvalidNameLength.Wrap(err.Error())
	}

	return nil
}

func validateStringSize(str string, limit int, name string) error {
	if length := utf8.RuneCountInString(str); length > limit {
		return sdkerrors.ErrInvalidRequest.Wrapf("%s cannot exceed %d in length: current %d", name, limit, length)
	}
	return nil
}
