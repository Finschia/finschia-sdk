package types

import (
	"errors"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		SwapInit: []SwapInit{},
		Swapped:  []Swapped{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
// need confirm: should we validate? Since it may nil
func (gs GenesisState) Validate() error {
	if len(gs.GetSwapInit()) == 0 {
		return nil
	}
	if len(gs.GetSwapped()) == 0 {
		return nil
	}
	if len(gs.GetSwapInit()) > 1 {
		return errors.New("cannot have more than one swapInit") // TODO(bjs) to sentinel
	}
	if len(gs.GetSwapped()) > 1 {
		return errors.New("cannot have more than one swapped") // TODO(bjs) to sentinel
	}
	swapInit := gs.GetSwapInit()[0]
	if err := swapInit.ValidateBasic(); err != nil {
		return err
	}
	swapped := gs.GetSwapped()[0]
	if err := swapped.Validate(); err != nil {
		return err
	}

	if swapInit.AmountCapForToDenom.LT(swapped.GetToCoinAmount().Amount) {
		return ErrExceedSwappable
	}
	return nil
}
