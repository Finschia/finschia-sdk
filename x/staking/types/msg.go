package types

import (
	"bytes"

	codectypes "github.com/line/lfb-sdk/codec/types"
	cryptotypes "github.com/line/lfb-sdk/crypto/types"
	sdk "github.com/line/lfb-sdk/types"
	sdkerrors "github.com/line/lfb-sdk/types/errors"
)

// staking message types
const (
	TypeMsgUndelegate      = "begin_unbonding"
	TypeMsgEditValidator   = "edit_validator"
	TypeMsgCreateValidator = "create_validator"
	TypeMsgDelegate        = "delegate"
	TypeMsgBeginRedelegate = "begin_redelegate"
)

var (
	_ sdk.Msg                            = &MsgCreateValidator{}
	_ codectypes.UnpackInterfacesMessage = (*MsgCreateValidator)(nil)
	_ sdk.Msg                            = &MsgCreateValidator{}
	_ sdk.Msg                            = &MsgEditValidator{}
	_ sdk.Msg                            = &MsgDelegate{}
	_ sdk.Msg                            = &MsgUndelegate{}
	_ sdk.Msg                            = &MsgBeginRedelegate{}
)

// NewMsgCreateValidator creates a new MsgCreateValidator instance.
// Delegator address and validator address are the same.
func NewMsgCreateValidator(
	valAddr sdk.ValAddress, pubKey cryptotypes.PubKey, //nolint:interfacer
	selfDelegation sdk.Coin, description Description, commission CommissionRates, minSelfDelegation sdk.Int,
) (*MsgCreateValidator, error) {
	var pkAny *codectypes.Any
	if pubKey != nil {
		var err error
		if pkAny, err = codectypes.NewAnyWithValue(pubKey); err != nil {
			return nil, err
		}
	}
	return &MsgCreateValidator{
		Description:       description,
		DelegatorAddress:  valAddr.ToAccAddress().String(),
		ValidatorAddress:  valAddr.String(),
		Pubkey:            pkAny,
		Value:             selfDelegation,
		Commission:        commission,
		MinSelfDelegation: minSelfDelegation,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgCreateValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCreateValidator) Type() string { return TypeMsgCreateValidator }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the validator address is not same as delegator's, then the validator must
// sign the msg as well.
func (msg MsgCreateValidator) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	addrs := []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
	valAddr := sdk.ValAddress(msg.ValidatorAddress).ToAccAddress()
	if !addrs[0].Equals(valAddr) {
		addrs = append(addrs, valAddr)
	}

	return addrs
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateValidator) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures either empty or valid
	delAddrBytes, err := sdk.AccAddressToBytes(msg.DelegatorAddress)
	if err != nil {
		return err
	}
	if len(delAddrBytes) == 0 {
		return ErrEmptyDelegatorAddr
	}

	if msg.ValidatorAddress == "" {
		return ErrEmptyValidatorAddr
	}

	valAddrBytes, err := sdk.ValAddressToBytes(msg.ValidatorAddress)
	if err != nil {
		return err
	}
	if !bytes.Equal(delAddrBytes, valAddrBytes) {
		return ErrBadValidatorAddr
	}

	if msg.Pubkey == nil {
		return ErrEmptyValidatorPubKey
	}

	if !msg.Value.IsValid() || !msg.Value.Amount.IsPositive() {
		return ErrBadDelegationAmount
	}

	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.Commission == (CommissionRates{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty commission")
	}

	if err := msg.Commission.Validate(); err != nil {
		return err
	}

	if !msg.MinSelfDelegation.IsPositive() {
		return ErrMinSelfDelegationInvalid
	}

	if msg.Value.Amount.LT(msg.MinSelfDelegation) {
		return ErrSelfDelegationBelowMinimum
	}

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateValidator) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pubKey cryptotypes.PubKey
	return unpacker.UnpackAny(msg.Pubkey, &pubKey)
}

// NewMsgEditValidator creates a new MsgEditValidator instance
//nolint:interfacer
func NewMsgEditValidator(valAddr sdk.ValAddress, description Description, newRate *sdk.Dec, newMinSelfDelegation *sdk.Int) *MsgEditValidator {
	return &MsgEditValidator{
		Description:       description,
		CommissionRate:    newRate,
		ValidatorAddress:  valAddr.String(),
		MinSelfDelegation: newMinSelfDelegation,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgEditValidator) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgEditValidator) Type() string { return TypeMsgEditValidator }

// GetSigners implements the sdk.Msg interface.
func (msg MsgEditValidator) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.ValAddress(msg.ValidatorAddress).ToAccAddress()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgEditValidator) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgEditValidator) ValidateBasic() error {
	if msg.ValidatorAddress == "" {
		return ErrEmptyValidatorAddr
	}
	if err := sdk.ValidateValAddress(msg.ValidatorAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid validator address (%s)", err)
	}

	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.MinSelfDelegation != nil && !msg.MinSelfDelegation.IsPositive() {
		return ErrMinSelfDelegationInvalid
	}

	if msg.CommissionRate != nil {
		if msg.CommissionRate.GT(sdk.OneDec()) || msg.CommissionRate.IsNegative() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "commission rate must be between 0 and 1 (inclusive)")
		}
	}

	return nil
}

// NewMsgDelegate creates a new MsgDelegate instance.
//nolint:interfacer
func NewMsgDelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) *MsgDelegate {
	return &MsgDelegate{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
		Amount:           amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgDelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgDelegate) Type() string { return TypeMsgDelegate }

// GetSigners implements the sdk.Msg interface.
func (msg MsgDelegate) GetSigners() []sdk.AccAddress {
	err := sdk.ValidateAccAddress(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgDelegate) ValidateBasic() error {
	if msg.DelegatorAddress == "" {
		return ErrEmptyDelegatorAddr
	}
	if err := sdk.ValidateAccAddress(msg.DelegatorAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid delegator address (%s)", err)
	}

	if msg.ValidatorAddress == "" {
		return ErrEmptyValidatorAddr
	}
	if err := sdk.ValidateValAddress(msg.ValidatorAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid validator address (%s)", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return ErrBadDelegationAmount
	}

	return nil
}

// NewMsgBeginRedelegate creates a new MsgBeginRedelegate instance.
//nolint:interfacer
func NewMsgBeginRedelegate(
	delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, amount sdk.Coin,
) *MsgBeginRedelegate {
	return &MsgBeginRedelegate{
		DelegatorAddress:    delAddr.String(),
		ValidatorSrcAddress: valSrcAddr.String(),
		ValidatorDstAddress: valDstAddr.String(),
		Amount:              amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgBeginRedelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface
func (msg MsgBeginRedelegate) Type() string { return TypeMsgBeginRedelegate }

// GetSigners implements the sdk.Msg interface
func (msg MsgBeginRedelegate) GetSigners() []sdk.AccAddress {
	err := sdk.ValidateAccAddress(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgBeginRedelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgBeginRedelegate) ValidateBasic() error {
	if msg.DelegatorAddress == "" {
		return ErrEmptyDelegatorAddr
	}
	if err := sdk.ValidateAccAddress(msg.DelegatorAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid delegator address (%s)", err)
	}

	if msg.ValidatorSrcAddress == "" {
		return ErrEmptyValidatorAddr
	}
	if err := sdk.ValidateValAddress(msg.ValidatorSrcAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid source validator address (%s)", err)
	}

	if msg.ValidatorDstAddress == "" {
		return ErrEmptyValidatorAddr
	}
	if err := sdk.ValidateValAddress(msg.ValidatorDstAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid destination validator address (%s)", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return ErrBadSharesAmount
	}

	return nil
}

// NewMsgUndelegate creates a new MsgUndelegate instance.
//nolint:interfacer
func NewMsgUndelegate(delAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Coin) *MsgUndelegate {
	return &MsgUndelegate{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
		Amount:           amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgUndelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgUndelegate) Type() string { return TypeMsgUndelegate }

// GetSigners implements the sdk.Msg interface.
func (msg MsgUndelegate) GetSigners() []sdk.AccAddress {
	err := sdk.ValidateAccAddress(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(msg.DelegatorAddress)}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgUndelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUndelegate) ValidateBasic() error {
	if msg.DelegatorAddress == "" {
		return ErrEmptyDelegatorAddr
	}
	if err := sdk.ValidateAccAddress(msg.DelegatorAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid delegator address (%s)", err)
	}

	if msg.ValidatorAddress == "" {
		return ErrEmptyValidatorAddr
	}
	if err := sdk.ValidateValAddress(msg.ValidatorAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid validator address (%s)", err)
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return ErrBadSharesAmount
	}

	return nil
}
