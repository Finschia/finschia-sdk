package simulation_test

import (
	"math/big"

	sdk "github.com/line/lbm-sdk/v2/types"
)

func init() {
	sdk.PowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}
