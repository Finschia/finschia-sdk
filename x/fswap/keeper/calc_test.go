package keeper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/Finschia/finschia-sdk/types"
)

func TestCalcSwap(t *testing.T) {
	rateDot5, err := sdk.NewDecFromStr("0.5")
	require.NoError(t, err)
	rateDot3maxPrecision, err := sdk.NewDecFromStr("0.333333333333333333")
	require.NoError(t, err)

	finschiaSwapRate, err := sdk.NewDecFromStr("148.079656")
	require.NoError(t, err)
	conySwapRate := finschiaSwapRate.Mul(sdk.NewDec(1000000))
	pebSwapRateForCony, err := sdk.NewDecFromStr("148079656000000")
	require.NoError(t, err)
	testCases := map[string]struct {
		fromAmount     sdk.Int
		expectedAmount sdk.Int
		swapRate       sdk.Dec
	}{
		"swapRate 0.5": {
			fromAmount:     sdk.ZeroInt(),
			swapRate:       rateDot5,
			expectedAmount: sdk.ZeroInt(),
		},
		"swapRate 0.333333333333333333": {
			fromAmount:     sdk.NewInt(3),
			swapRate:       rateDot3maxPrecision,
			expectedAmount: sdk.ZeroInt(),
		},
		"swapRate conySwapRate(148.079656 * 10^6) fromAmount(1)": {
			fromAmount:     sdk.NewInt(1),
			swapRate:       conySwapRate,
			expectedAmount: sdk.NewInt(148079656),
		},
		"swapRate conySwapRate(148.079656 * 10^6) fromAmount(3)": {
			fromAmount:     sdk.NewInt(3),
			swapRate:       conySwapRate,
			expectedAmount: sdk.NewInt(444238968),
		},
		"pebSwapRateForCony pebSwapRateForCony(148.079656 * 10^12) fromAmount(1)": {
			fromAmount:     sdk.NewInt(1),
			swapRate:       pebSwapRateForCony,
			expectedAmount: sdk.NewInt(148079656000000),
		},
		"pebSwapRateForCony pebSwapRateForCony(148.079656 * 10^12) fromAmount(3)": {
			fromAmount:     sdk.NewInt(3),
			swapRate:       pebSwapRateForCony,
			expectedAmount: sdk.NewInt(444238968000000),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actualAmount := CalcSwap(tc.swapRate, tc.fromAmount)
			require.Equal(t, tc.expectedAmount, actualAmount, fmt.Sprintf("tc.expectedAmount = %v, actualAmount = %v", tc.expectedAmount, actualAmount))
		})
	}
}
