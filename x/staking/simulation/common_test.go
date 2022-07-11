package simulation_test

import (
	"math/big"

	sdk "github.com/line/lbm-sdk/types"
)

func init() {
	sdk.DefaultPowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}
