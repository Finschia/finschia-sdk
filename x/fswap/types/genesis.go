package types

// NewGenesis creates a new genesis state
func NewGenesisState(config Config) *GenesisState {
	return &GenesisState{
		Swapped: NewSwapped(config),
	}
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesisState(DefaultConfig())
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Swapped.Validate(); err != nil {
		return err
	}

	return nil
}
