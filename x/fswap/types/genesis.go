package types

import (
	"errors"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		FswapInit: []FswapInit{},
		Swapped:   []Swapped{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
// need confirm: should we validate? Since it may nil
func (gs GenesisState) Validate() error {
	if len(gs.GetFswapInit()) == 0 {
		return nil
	}
	if len(gs.GetSwapped()) == 0 {
		return nil
	}
	if len(gs.GetFswapInit()) > 1 {
		return errors.New("cannot have more than one fswap") // TODO(bjs) to sentinel
	}
	if len(gs.GetSwapped()) > 1 {
		return errors.New("cannot have more than one swapped") // TODO(bjs) to sentinel
	}
	fswap := gs.GetFswapInit()[0]
	if err := fswap.ValidateBasic(); err != nil {
		return err
	}
	swapped := gs.GetSwapped()[0]
	if err := swapped.Validate(); err != nil {
		return err
	}

	if fswap.AmountCapForToDenom.LT(swapped.GetNewCoinAmount().Amount) {
		return ErrExceedSwappable
	}
	return nil
}
