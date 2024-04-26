package types

// NewGenesis creates a new genesis state
func NewGenesisState(fswapInit FswapInit, swapped Swapped) *GenesisState {
	return &GenesisState{
		FswapInit: fswapInit,
		Swapped:   swapped,
	}
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{DefaultFswapInit(), DefaultSwapped()}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
// need confirm: should we validate? Since it may nil
func (gs GenesisState) Validate() error {
	if gs == *DefaultGenesis() {
		return nil
	}
	if err := gs.FswapInit.ValidateBasic(); err != nil {
		return err
	}
	if err := gs.Swapped.Validate(); err != nil {
		return err
	}
	if gs.FswapInit.AmountLimit.LT(gs.Swapped.NewCoinAmount) {
		return ErrExceedSwappable
	}
	return nil
}
