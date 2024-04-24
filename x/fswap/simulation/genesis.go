package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/module"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// Simulation parameter constants
const (
	OldCoinAmount          = "old_coin_amount"
	NewCoinAmount          = "new_coin_amount"
	SwappableNewCoinAmount = "swappable_new_coin_amount"
)

// GenOldCoinAmount randomized oldCoinAmount
func GenOldCoinAmount(r *rand.Rand) sdk.Int {
	return sdk.NewInt(int64(r.Intn(1001) + 1000))
}

// GenNewCoinAmount randomized oldCoinAmount
func GenNewCoinAmount(r *rand.Rand) sdk.Int {
	return sdk.NewInt(int64(r.Intn(100001) + 100000))
}

// GenSwappableNewCoinAmount randomized swappableNewCoinAmount
func GenSwappableNewCoinAmount(r *rand.Rand) sdk.Int {
	return sdk.NewInt(int64(r.Intn(100001) + 200000))
}

// RandomizedGenState generates a random GenesisState for fswap
func RandomizedGenState(simState *module.SimulationState) {

	var oldCoinAmount sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, OldCoinAmount, &oldCoinAmount, simState.Rand,
		func(r *rand.Rand) { oldCoinAmount = GenOldCoinAmount(r) },
	)

	var newCoinAmount sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, NewCoinAmount, &newCoinAmount, simState.Rand,
		func(r *rand.Rand) { newCoinAmount = GenNewCoinAmount(r) },
	)

	var swappableNewCoinAmount sdk.Int
	simState.AppParams.GetOrGenerate(
		simState.Cdc, SwappableNewCoinAmount, &swappableNewCoinAmount, simState.Rand,
		func(r *rand.Rand) { swappableNewCoinAmount = GenSwappableNewCoinAmount(r) },
	)

	fswapParams := types.NewParams(swappableNewCoinAmount)
	fswapSwapped := types.NewSwapped(oldCoinAmount, newCoinAmount)
	fswapGenesis := types.NewGenesisState(fswapParams, fswapSwapped)

	bz, err := json.MarshalIndent(&fswapGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated fswap parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(fswapGenesis)
}
