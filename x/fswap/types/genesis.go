package types

import sdk "github.com/Finschia/finschia-sdk/types"

var (
	DefaultTotalSupply = sdk.NewCoin(DefaultConfig().NewCoinDenom, sdk.NewInt(100000))
)

// NewGenesis creates a new genesis state
func NewGenesisState(config Config, swappableNewCoinAmount sdk.Coin) *GenesisState {
	return &GenesisState{
		Swapped:                NewSwapped(config),
		SwappableNewCoinAmount: swappableNewCoinAmount,
	}
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesisState(DefaultConfig(), DefaultTotalSupply)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Swapped.Validate(); err != nil {
		return err
	}
	if err := ValidateCoin(gs.SwappableNewCoinAmount); err != nil {
		return err
	}
	return nil
}
