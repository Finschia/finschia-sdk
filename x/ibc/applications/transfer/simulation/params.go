package simulation

import (
	"fmt"
	"math/rand"

	simtypes "github.com/line/lbm-sdk/types/simulation"
	"github.com/line/lbm-sdk/x/simulation"
	gogotypes "github.com/gogo/protobuf/types"

	"github.com/line/lbm-sdk/x/ibc/applications/transfer/types"
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
