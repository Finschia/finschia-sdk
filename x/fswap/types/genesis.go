package types

// NewGenesis creates a new genesis state
func NewGenesisState(params Params, swapped Swapped) *GenesisState {
	return &GenesisState{
		Params:  params,
		Swapped: swapped,
	}
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesisState(DefaultParams(), DefaultSwapped())
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
	if gs.Params.SwappableNewCoinAmount.LT(gs.Swapped.NewCoinAmount) {
		return ErrExceedSwappable
	}
	return nil
}
