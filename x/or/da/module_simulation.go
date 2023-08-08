package da

import (
	"math/rand"

	simappparams "github.com/Finschia/finschia-rdk/simapp/params"
	sdk "github.com/Finschia/finschia-rdk/types"
	"github.com/Finschia/finschia-rdk/types/module"
	simtypes "github.com/Finschia/finschia-rdk/types/simulation"
	dasimulation "github.com/Finschia/finschia-rdk/x/or/da/simulation"
	datest "github.com/Finschia/finschia-rdk/x/or/da/testutil"
	"github.com/Finschia/finschia-rdk/x/or/da/types"
	"github.com/Finschia/finschia-rdk/x/simulation"
)

// avoid unused import issue
var (
	_ = datest.AccAddress
	_ = dasimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
)

func GenBatchMaxBytes(r *rand.Rand) uint64 {
	return uint64(r.Intn(100000*1000000)%1000 + 1)
}

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

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	return operations
}
