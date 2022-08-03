package simulation

import (
	"fmt"
	"math/rand"

	"github.com/line/lbm-sdk/codec"
	simtypes "github.com/line/lbm-sdk/types/simulation"
	"github.com/line/lbm-sdk/x/simulation"
	lbmwasmtypes "github.com/line/lbm-sdk/x/wasm/lbm/types"
	"github.com/line/lbm-sdk/x/wasm/types"
)

func ParamChanges(r *rand.Rand, cdc codec.Codec) []simtypes.ParamChange {
	params := RandomParams(r)
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyUploadAccess),
			func(r *rand.Rand) string {
				jsonBz, err := cdc.MarshalJSON(&params.CodeUploadAccess)
				if err != nil {
					panic(err)
				}
				return string(jsonBz)
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.ParamStoreKeyInstantiateAccess),
			func(r *rand.Rand) string {
				return fmt.Sprintf("%q", params.CodeUploadAccess.Permission.String())
			},
		),
	}
}

func RandomParams(r *rand.Rand) lbmwasmtypes.Params {
	permissionType := types.AccessType(simtypes.RandIntBetween(r, 1, 3))
	account, _ := simtypes.RandomAcc(r, simtypes.RandomAccounts(r, 10))
	accessConfig := permissionType.With(account.Address)
	return lbmwasmtypes.Params{
		CodeUploadAccess:             accessConfig,
		InstantiateDefaultPermission: accessConfig.Permission,
		GasMultiplier:                types.DefaultGasMultiplier,
		InstanceCost:                 types.DefaultInstanceCost,
		CompileCost:                  types.DefaultCompileCost,
	}
}
