// package fswap

// import (
// 	"math/rand"

// 	"github.com/Finschia/finschia-sdk/testutil/sample"

// 	fswapsimulation "github.com/Finschia/finschia-sdk/x/fswap/simulation"
// 	"github.com/Finschia/finschia-sdk/x/fswap/types"

// 	"github.com/Finschia/finschia-sdk/baseapp"
// 	sdk "github.com/Finschia/finschia-sdk/types"
// 	"github.com/Finschia/finschia-sdk/types/module"
// 	simtypes "github.com/Finschia/finschia-sdk/types/simulation"
// 	"github.com/Finschia/finschia-sdk/x/simulation"
// )

// // avoid unused import issue
// var (
// 	_ = sample.AccAddress
// 	_ = fswapsimulation.FindAccount
// 	_ = simulation.MsgEntryKind
// 	_ = baseapp.Paramspace
// 	_ = rand.Rand{}
// )

// const (
// // this line is used by starport scaffolding # simapp/module/const
// )

// // GenerateGenesisState creates a randomized GenState of the module.
// func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
// 	accs := make([]string, len(simState.Accounts))
// 	for i, acc := range simState.Accounts {
// 		accs[i] = acc.Address.String()
// 	}
// 	fswapGenesis := types.GenesisState{
// 		Params: types.DefaultParams(),
// 		// this line is used by starport scaffolding # simapp/module/genesisState
// 	}
// 	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&fswapGenesis)
// }

// // RegisterStoreDecoder registers a decoder.
// func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// // ProposalContents doesn't return any content functions for governance proposals.
// func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
// 	return nil
// }

// // WeightedOperations returns the all the gov module operations with their respective weights.
// func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
// 	operations := make([]simtypes.WeightedOperation, 0)

// 	// this line is used by starport scaffolding # simapp/module/operation

// 	return operations
// }

// // ProposalMsgs returns msgs used for governance proposals for simulations.
// func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
// 	return []simtypes.WeightedProposalMsg{
// 		// this line is used by starport scaffolding # simapp/module/OpMsg
// 	}
// }
