package simulation

import (
	"fmt"
	"math/rand"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/line/lbm-sdk/v2/x/simulation"

	simtypes "github.com/line/lbm-sdk/v2/types/simulation"
	"github.com/line/lbm-sdk/v2/x/ibc/applications/transfer/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeySendEnabled),
			func(r *rand.Rand) string {
				sendEnabled := RadomEnabled(r)
				return fmt.Sprintf("%s", types.ModuleCdc.MustMarshalJSON(&gogotypes.BoolValue{Value: sendEnabled}))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyReceiveEnabled),
			func(r *rand.Rand) string {
				receiveEnabled := RadomEnabled(r)
				return fmt.Sprintf("%s", types.ModuleCdc.MustMarshalJSON(&gogotypes.BoolValue{Value: receiveEnabled}))
			},
		),
	}
}
