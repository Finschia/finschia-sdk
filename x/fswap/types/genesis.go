package types

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
