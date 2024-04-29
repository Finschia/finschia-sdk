package types

func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	return nil
}
