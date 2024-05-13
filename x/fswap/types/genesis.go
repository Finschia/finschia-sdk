package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Swaps:     []Swap{},
		SwapStats: SwapStats{},
		Swappeds:  []Swapped{},
	}
}

func DefaultAuthority() sdk.AccAddress {
	return authtypes.NewModuleAddress(foundation.ModuleName)
}

func AuthorityCandidates() []sdk.AccAddress {
	return []sdk.AccAddress{
		authtypes.NewModuleAddress(govtypes.ModuleName),
		authtypes.NewModuleAddress(foundation.ModuleName),
	}
}

// Validate performs basic genesis state validation returning an error upon any failure.
func (gs *GenesisState) Validate() error {
	if err := gs.SwapStats.ValidateBasic(); err != nil {
		return err
	}

	if len(gs.GetSwaps()) != len(gs.GetSwappeds()) {
		return ErrInvalidState.Wrap("number of swaps does not match number of Swappeds")
	}

	if len(gs.GetSwaps()) != int(gs.GetSwapStats().SwapCount) {
		return ErrInvalidState.Wrap("number of swaps does not match swap count in SwapStats")
	}

	swappeds := gs.GetSwappeds()
	for i, swap := range gs.GetSwaps() {
		swapped := swappeds[i]
		if err := swap.ValidateBasic(); err != nil {
			return err
		}
		if err := swapped.ValidateBasic(); err != nil {
			return err
		}
		if swap.FromDenom != swapped.FromCoinAmount.Denom {
			return ErrInvalidState.Wrap("FromCoin in swap and swapped do not correspond")
		}
		if swap.ToDenom != swapped.ToCoinAmount.Denom {
			return ErrInvalidState.Wrap("ToCoin in swap and swapped do not correspond")
		}
		if swap.AmountCapForToDenom.LT(swapped.ToCoinAmount.Amount) {
			return ErrInvalidState.Wrap("AmountCapForToDenom cannot be exceeded")
		}
	}
	return nil
}
