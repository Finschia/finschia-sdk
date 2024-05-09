package keeper

import "github.com/Finschia/finschia-sdk/types"

func CalcSwap(swapRate types.Dec, fromCoinAmount types.Int) types.Int {
	return swapRate.MulTruncate(types.NewDecFromBigInt(fromCoinAmount.BigInt())).TruncateInt()
}
