package types

// NewGenesisState creates a new GenesisState object
func NewGenesisState(enabled bool, allowedValidators []string) *GenesisState {
	return &GenesisState{
		Enabled:           enabled,
		AllowedValidators: allowedValidators,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Enabled:           false,
		AllowedValidators: []string{},
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	if !data.Enabled && len(data.AllowedValidators) != 0 {
		return nil // TODO 
	}
	return nil
}
