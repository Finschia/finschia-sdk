package da

import (
	"math/rand"

	dasimulation "github.com/Finschia/finschia-sdk/x/or/da/simulation"
	"github.com/Finschia/finschia-sdk/x/or/da/types"

	"github.com/Finschia/finschia-sdk/baseapp"
	simappparams "github.com/Finschia/finschia-sdk/simapp/params"
	"github.com/Finschia/finschia-sdk/x/simulation"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/module"
	simtypes "github.com/Finschia/finschia-sdk/types/simulation"
	datest "github.com/Finschia/finschia-sdk/x/or/da/testutil"
)

// avoid unused import issue
var (
	_ = datest.AccAddress
	_ = dasimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	daGenesis := types.GenesisState{
		Params: types.DefaultParams(),
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&daGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {
	daParams := types.DefaultParams()
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyPlaceholder), func(r *rand.Rand) string {
			return string(types.Amino.MustMarshalJSON(daParams.Placeholder))
		}),
	}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	return operations
}
