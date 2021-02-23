package types

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const RouterKey = ModuleName

// MsgSend - high level transaction of the coin module
type MsgSend struct {
	From   sdk.AccAddress `json:"from"`
	To     sdk.AccAddress `json:"to"`
	Amount sdk.Coins      `json:"amount"`
}

var _ sdk.Msg = MsgSend{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSend(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) MsgSend {
	return MsgSend{From: fromAddr, To: toAddr, Amount: amount}
}

// Route Implements Msg.
func (msg MsgSend) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSend) Type() string { return "send" }

// ValidateBasic Implements Msg.
func (msg MsgSend) ValidateBasic() error {
	if msg.From.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	if msg.To.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
	}
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	if !msg.Amount.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	if err := validateDenomination(msg.Amount); err != nil {
		return err
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// MsgMultiSend - high level transaction of the coin module
type MsgMultiSend struct {
	Inputs  []Input  `json:"inputs" yaml:"inputs"`
	Outputs []Output `json:"outputs" yaml:"outputs"`
}

var _ sdk.Msg = MsgMultiSend{}

// NewMsgMultiSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgMultiSend(in []Input, out []Output) MsgMultiSend {
	return MsgMultiSend{Inputs: in, Outputs: out}
}

// Route Implements Msg
func (msg MsgMultiSend) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgMultiSend) Type() string { return "multisend" }

// ValidateBasic Implements Msg.
func (msg MsgMultiSend) ValidateBasic() error {
	// this just makes sure all the inputs and outputs are properly formatted,
	// not that they actually have the money inside
	if len(msg.Inputs) == 0 {
		return ErrNoInputs
	}
	if len(msg.Outputs) == 0 {
		return ErrNoOutputs
	}

	return ValidateInputsOutputs(msg.Inputs, msg.Outputs)
}

// GetSignBytes Implements Msg.
func (msg MsgMultiSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgMultiSend) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, len(msg.Inputs))
	for i, in := range msg.Inputs {
		addrs[i] = in.Address
	}
	return addrs
}

// Input models transaction input
type Input struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Coins   sdk.Coins      `json:"coins" yaml:"coins"`
}

// ValidateBasic - validate transaction input
func (in Input) ValidateBasic() error {
	if len(in.Address) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "input address missing")
	}
	if !in.Coins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, in.Coins.String())
	}
	if !in.Coins.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, in.Coins.String())
	}
	if err := validateDenomination(in.Coins); err != nil {
		return err
	}
	return nil
}

// NewInput - create a transaction input, used with MsgMultiSend
func NewInput(addr sdk.AccAddress, coins sdk.Coins) Input {
	return Input{
		Address: addr,
		Coins:   coins,
	}
}

// Output models transaction outputs
type Output struct {
	Address sdk.AccAddress `json:"address" yaml:"address"`
	Coins   sdk.Coins      `json:"coins" yaml:"coins"`
}

// ValidateBasic - validate transaction output
func (out Output) ValidateBasic() error {
	if len(out.Address) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "output address missing")
	}
	if !out.Coins.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, out.Coins.String())
	}
	if !out.Coins.IsAllPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, out.Coins.String())
	}
	if err := validateDenomination(out.Coins); err != nil {
		return err
	}
	return nil
}

// NewOutput - create a transaction output, used with MsgMultiSend
func NewOutput(addr sdk.AccAddress, coins sdk.Coins) Output {
	return Output{
		Address: addr,
		Coins:   coins,
	}
}

// ValidateInputsOutputs validates that each respective input and output is
// valid and that the sum of inputs is equal to the sum of outputs.
func ValidateInputsOutputs(inputs []Input, outputs []Output) error {
	var totalIn, totalOut sdk.Coins

	for _, in := range inputs {
		if err := in.ValidateBasic(); err != nil {
			return err
		}
		totalIn = totalIn.Add(in.Coins...)
	}

	for _, out := range outputs {
		if err := out.ValidateBasic(); err != nil {
			return err
		}
		totalOut = totalOut.Add(out.Coins...)
	}

	// make sure inputs and outputs match
	if !totalIn.IsEqual(totalOut) {
		return ErrInputOutputMismatch
	}

	return nil
}
func validateDenomination(coins sdk.Coins) error {
	for _, coin := range coins {
		if err := ValidateSymbolReserved(coin.Denom); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidDenom, "invalid denom [%s] send message supports 3~5 length denom only", coin.Denom)
		}
	}
	return nil
}
