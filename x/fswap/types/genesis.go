package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
)

const (
	DefaultOldCoins string = "cony"
)

var DefaultSwapRate = sdk.NewDecWithPrec(148079656, 6)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:  DefaultParams(),
		Swapped: DefaultSwapped(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	if err := gs.Swapped.Validate(); err != nil {
		return err
	}

	return nil
}
