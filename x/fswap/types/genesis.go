package types

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Swaps:     []Swap{},
		SwapStats: SwapStats{},
		Swappeds:  []Swapped{},
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs *GenesisState) Validate() error {
	for _, swap := range gs.GetSwaps() {
		if err := swap.ValidateBasic(); err != nil {
			return err
		}
	}

	if err := gs.SwapStats.ValidateBasic(); err != nil {
		return err
	}

	for _, swapped := range gs.GetSwappeds() {
		if err := swapped.ValidateBasic(); err != nil {
			return err
		}
	}

	if len(gs.GetSwaps()) != len(gs.GetSwappeds()) {
		return ErrInvalidState.Wrap("number of swaps does not match number of Swappeds")
	}

	if len(gs.GetSwaps()) != int(gs.GetSwapStats().SwapCount) {
		return ErrInvalidState.Wrap("number of swaps does not match swap count in SwapStats")
	}

	return nil
}
