// DONTCOVER
// nolint
package v036

import (
	sdk "github.com/line/lbm-sdk/v2/types"
)

const ModuleName = "supply"

type (
	GenesisState struct {
		Supply sdk.Coins `json:"supply" yaml:"supply"`
	}
)

func EmptyGenesisState() GenesisState {
	return GenesisState{
		Supply: sdk.NewCoins(), // leave this empty as it's filled on initialization
	}
}
