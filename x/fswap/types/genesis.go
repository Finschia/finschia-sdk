package types

import sdk "github.com/Finschia/finschia-sdk/types"

// NewGenesis creates a new genesis state
func NewGenesisState(config Config, swappable sdk.Coin) *GenesisState {
	return &GenesisState{
		Params:  NewParams(swappable),
		Swapped: NewSwapped(config),
	}
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesisState(DefaultConfig(), DefaultSwappableNewCoinAmount)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	if err := gs.Swapped.Validate(); err != nil {
		return err
	}
	return nil
}
