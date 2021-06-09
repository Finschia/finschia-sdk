package simulation

import (
	"github.com/line/lfb-sdk/types/module"
	"github.com/line/lfb-sdk/x/wasm/internal/types"
)

// RandomizeGenState generates a random GenesisState for wasm
func RandomizedGenState(simstate *module.SimulationState) {
	params := RandomParams(simstate.Rand)
	wasmGenesis := types.GenesisState{
		Params:    params,
		Codes:     nil,
		Contracts: nil,
		Sequences: nil,
		GenMsgs:   nil,
	}

	_, err := simstate.Cdc.MarshalJSON(&wasmGenesis)
	if err != nil {
		panic(err)
	}

	simstate.GenState[types.ModuleName] = simstate.Cdc.MustMarshalJSON(&wasmGenesis)
}
