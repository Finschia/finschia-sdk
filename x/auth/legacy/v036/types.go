// DONTCOVER
// nolint
package v036

import v034auth "github.com/line/lbm-sdk/v2/x/auth/legacy/v034"

const (
	ModuleName = "auth"
)

type (
	GenesisState struct {
		Params v034auth.Params `json:"params"`
	}
)

func NewGenesisState(params v034auth.Params) GenesisState {
	return GenesisState{params}
}
