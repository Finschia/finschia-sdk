package consortium

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:         &Params{Enabled: false},
		ValidatorAuths: []*ValidatorAuth{},
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	// validator auths are redundant where consortium is off
	if !data.Params.Enabled && len(data.ValidatorAuths) != 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "redundant validator auths for disabled consortium")
	}

	for _, auth := range data.ValidatorAuths {
		if err := sdk.ValidateValAddress(auth.OperatorAddress); err != nil {
			return err
		}
	}

	return nil
}
