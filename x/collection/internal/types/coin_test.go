// copied from https://github.com/cosmos/cosmos-sdk/blob/v0.38.1/types/coin_test.go
package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/codec"

	sdk "github.com/line/lbm-sdk/types"
)

var (
	testDenom1 = "0000000100000000"
	testDenom2 = "0000000200000000"
	testDenom3 = "0000000300000000"
	testDenom4 = "0000000400000000"
	testDenomA = "0000000a00000000"
)

// ----------------------------------------------------------------------------
// Coin tests

func TestCoin(t *testing.T) {
	require.Panics(t, func() { NewInt64Coin(testDenom1, -1) })
	require.Panics(t, func() { NewCoin(testDenom1, sdk.NewInt(-1)) })
	require.Panics(t, func() { NewInt64Coin(strings.ToUpper(testDenomA), 10) })
	require.Panics(t, func() { NewCoin(strings.ToUpper(testDenomA), sdk.NewInt(10)) })
	require.Equal(t, sdk.NewInt(5), NewInt64Coin(testDenom1, 5).Amount)
	require.Equal(t, sdk.NewInt(5), NewCoin(testDenom1, sdk.NewInt(5)).Amount)
	require.Equal(t, OneCoins(testDenom1)[0].Amount.Int64(), int64(1))
}

func TestIsEqualCoin(t *testing.T) {
	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
		panics   bool
	}{
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom1, 1), true, false},
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom2, 1), false, true},
		{NewInt64Coin(testDenom3, 1), NewInt64Coin(testDenom3, 10), false, false},
	}

	for tcIndex, tc := range cases {
		tc := tc
		if tc.panics {
			require.Panics(t, func() { tc.inputOne.IsEqual(tc.inputTwo) })
		} else {
			res := tc.inputOne.IsEqual(tc.inputTwo)
			require.Equal(t, tc.expected, res, "coin equality relation is incorrect, tc #%d", tcIndex)
		}
	}
}

func TestCoinIsValid(t *testing.T) {
	cases := []struct {
		coin       Coin
		expectPass bool
	}{
		{Coin{testDenom1, sdk.NewInt(-1)}, false},
		{Coin{testDenom1, sdk.NewInt(0)}, true},
		{Coin{testDenom1, sdk.NewInt(1)}, true},
		{Coin{"a", sdk.NewInt(1)}, false},
		{Coin{"a very long coin denom", sdk.NewInt(1)}, false},
		{Coin{"atOm", sdk.NewInt(1)}, false},
		{Coin{"     ", sdk.NewInt(1)}, false},
	}

	for i, tc := range cases {
		require.Equal(t, tc.expectPass, tc.coin.IsValid(), "unexpected result for IsValid, tc #%d", i)
	}
}

func TestAddCoin(t *testing.T) {
	cases := []struct {
		inputOne    Coin
		inputTwo    Coin
		expected    Coin
		shouldPanic bool
	}{
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom1, 2), false},
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom1, 0), NewInt64Coin(testDenom1, 1), false},
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom2, 1), NewInt64Coin(testDenom1, 1), true},
	}

	for tcIndex, tc := range cases {
		tc := tc
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Add(tc.inputTwo) })
		} else {
			res := tc.inputOne.Add(tc.inputTwo)
			require.Equal(t, tc.expected, res, "sum of coins is incorrect, tc #%d", tcIndex)
		}
	}
}

func TestSubCoin(t *testing.T) {
	cases := []struct {
		inputOne    Coin
		inputTwo    Coin
		expected    Coin
		shouldPanic bool
	}{
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom2, 1), NewInt64Coin(testDenom1, 1), true},
		{NewInt64Coin(testDenom1, 10), NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom1, 9), false},
		{NewInt64Coin(testDenom1, 5), NewInt64Coin(testDenom1, 3), NewInt64Coin(testDenom1, 2), false},
		{NewInt64Coin(testDenom1, 5), NewInt64Coin(testDenom1, 0), NewInt64Coin(testDenom1, 5), false},
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom1, 5), Coin{}, true},
	}

	for tcIndex, tc := range cases {
		tc := tc
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Sub(tc.inputTwo) })
		} else {
			res := tc.inputOne.Sub(tc.inputTwo)
			require.Equal(t, tc.expected, res, "difference of coins is incorrect, tc #%d", tcIndex)
		}
	}

	tc := struct {
		inputOne Coin
		inputTwo Coin
		expected int64
	}{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom1, 1), 0}
	res := tc.inputOne.Sub(tc.inputTwo)
	require.Equal(t, tc.expected, res.Amount.Int64())
}

func TestIsGTECoin(t *testing.T) {
	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
		panics   bool
	}{
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom1, 1), true, false},
		{NewInt64Coin(testDenom1, 2), NewInt64Coin(testDenom1, 1), true, false},
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom2, 1), false, true},
	}

	for tcIndex, tc := range cases {
		tc := tc
		if tc.panics {
			require.Panics(t, func() { tc.inputOne.IsGTE(tc.inputTwo) })
		} else {
			res := tc.inputOne.IsGTE(tc.inputTwo)
			require.Equal(t, tc.expected, res, "coin GTE relation is incorrect, tc #%d", tcIndex)
		}
	}
}

func TestIsLTCoin(t *testing.T) {
	cases := []struct {
		inputOne Coin
		inputTwo Coin
		expected bool
		panics   bool
	}{
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom1, 1), false, false},
		{NewInt64Coin(testDenom1, 2), NewInt64Coin(testDenom1, 1), false, false},
		{NewInt64Coin(testDenom1, 0), NewInt64Coin(testDenom2, 1), false, true},
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom2, 1), false, true},
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom1, 1), false, false},
		{NewInt64Coin(testDenom1, 1), NewInt64Coin(testDenom1, 2), true, false},
	}

	for tcIndex, tc := range cases {
		tc := tc
		if tc.panics {
			require.Panics(t, func() { tc.inputOne.IsLT(tc.inputTwo) })
		} else {
			res := tc.inputOne.IsLT(tc.inputTwo)
			require.Equal(t, tc.expected, res, "coin LT relation is incorrect, tc #%d", tcIndex)
		}
	}
}

func TestCoinIsZero(t *testing.T) {
	coin := NewInt64Coin(testDenom1, 0)
	res := coin.IsZero()
	require.True(t, res)

	coin = NewInt64Coin(testDenom1, 1)
	res = coin.IsZero()
	require.False(t, res)
}

// ----------------------------------------------------------------------------
// Coins tests

func TestIsZeroCoins(t *testing.T) {
	cases := []struct {
		inputOne Coins
		expected bool
	}{
		{Coins{}, true},
		{Coins{NewInt64Coin(testDenom1, 0)}, true},
		{Coins{NewInt64Coin(testDenom1, 0), NewInt64Coin(testDenom2, 0)}, true},
		{Coins{NewInt64Coin(testDenom1, 1)}, false},
		{Coins{NewInt64Coin(testDenom1, 0), NewInt64Coin(testDenom2, 1)}, false},
	}

	for _, tc := range cases {
		res := tc.inputOne.IsZero()
		require.Equal(t, tc.expected, res)
	}
}

func TestEqualCoins(t *testing.T) {
	cases := []struct {
		inputOne Coins
		inputTwo Coins
		expected bool
		panics   bool
	}{
		{Coins{}, Coins{}, true, false},
		{Coins{NewInt64Coin(testDenom1, 0)}, Coins{NewInt64Coin(testDenom1, 0)}, true, false},
		{Coins{NewInt64Coin(testDenom1, 0), NewInt64Coin(testDenom2, 1)}, Coins{NewInt64Coin(testDenom1, 0), NewInt64Coin(testDenom2, 1)}, true, false},
		{Coins{NewInt64Coin(testDenom1, 0)}, Coins{NewInt64Coin(testDenom2, 0)}, false, true},
		{Coins{NewInt64Coin(testDenom1, 0)}, Coins{NewInt64Coin(testDenom1, 1)}, false, false},
		{Coins{NewInt64Coin(testDenom1, 0)}, Coins{NewInt64Coin(testDenom1, 0), NewInt64Coin(testDenom2, 1)}, false, false},
		{Coins{NewInt64Coin(testDenom1, 0), NewInt64Coin(testDenom2, 1)}, Coins{NewInt64Coin(testDenom1, 0), NewInt64Coin(testDenom2, 1)}, true, false},
	}

	for tcnum, tc := range cases {
		tc := tc
		if tc.panics {
			require.Panics(t, func() { tc.inputOne.IsEqual(tc.inputTwo) })
		} else {
			res := tc.inputOne.IsEqual(tc.inputTwo)
			require.Equal(t, tc.expected, res, "Equality is differed from exported. tc #%d, expected %b, actual %b.", tcnum, tc.expected, res)
		}
	}
}

func TestAddCoins(t *testing.T) {
	zero := sdk.NewInt(0)
	one := sdk.NewInt(1)
	two := sdk.NewInt(2)

	cases := []struct {
		inputOne Coins
		inputTwo Coins
		expected Coins
	}{
		{Coins{{testDenom1, one}, {testDenom2, one}}, Coins{{testDenom1, one}, {testDenom2, one}}, Coins{{testDenom1, two}, {testDenom2, two}}},
		{Coins{{testDenom1, zero}, {testDenom2, one}}, Coins{{testDenom1, zero}, {testDenom2, zero}}, Coins{{testDenom2, one}}},
		{Coins{{testDenom1, two}}, Coins{{testDenom2, zero}}, Coins{{testDenom1, two}}},
		{Coins{{testDenom1, one}}, Coins{{testDenom1, one}, {testDenom2, two}}, Coins{{testDenom1, two}, {testDenom2, two}}},
		{Coins{{testDenom1, zero}, {testDenom2, zero}}, Coins{{testDenom1, zero}, {testDenom2, zero}}, Coins(nil)},
	}

	for tcIndex, tc := range cases {
		res := tc.inputOne.Add(tc.inputTwo...)
		assert.True(t, res.IsValid())
		require.Equal(t, tc.expected, res, "sum of coins is incorrect, tc #%d", tcIndex)
	}
}

func TestSubCoins(t *testing.T) {
	zero := sdk.NewInt(0)
	one := sdk.NewInt(1)
	two := sdk.NewInt(2)

	testCases := []struct {
		inputOne    Coins
		inputTwo    Coins
		expected    Coins
		shouldPanic bool
	}{
		{Coins{{testDenom1, two}}, Coins{{testDenom1, one}, {testDenom2, two}}, Coins{{testDenom1, one}, {testDenom2, two}}, true},
		{Coins{{testDenom1, two}}, Coins{{testDenom2, zero}}, Coins{{testDenom1, two}}, false},
		{Coins{{testDenom1, one}}, Coins{{testDenom2, zero}}, Coins{{testDenom1, one}}, false},
		{Coins{{testDenom1, one}, {testDenom2, one}}, Coins{{testDenom1, one}}, Coins{{testDenom2, one}}, false},
		{Coins{{testDenom1, one}, {testDenom2, one}}, Coins{{testDenom1, two}}, Coins{}, true},
	}

	for i, tc := range testCases {
		tc := tc
		if tc.shouldPanic {
			require.Panics(t, func() { tc.inputOne.Sub(tc.inputTwo) })
		} else {
			res := tc.inputOne.Sub(tc.inputTwo)
			assert.True(t, res.IsValid())
			require.Equal(t, tc.expected, res, "sum of coins is incorrect, tc #%d", i)
		}
	}
}

func TestCoins(t *testing.T) {
	good := Coins{
		{testDenom1, sdk.NewInt(1)},
		{testDenom2, sdk.NewInt(1)},
		{testDenom3, sdk.NewInt(1)},
	}
	mixedCase1 := Coins{
		{"000A0000", sdk.NewInt(1)},
		{"A0000000", sdk.NewInt(1)},
		{"ABCD0000", sdk.NewInt(1)},
	}
	mixedCase2 := Coins{
		{"000A0000", sdk.NewInt(1)},
		{testDenom3, sdk.NewInt(1)},
	}
	mixedCase3 := Coins{
		{"000A0000", sdk.NewInt(1)},
	}
	empty := NewCoins()
	badSort1 := Coins{
		{testDenom3, sdk.NewInt(1)},
		{testDenom1, sdk.NewInt(1)},
		{testDenom2, sdk.NewInt(1)},
	}

	// both are after the first one, but the second and third are in the wrong order
	badSort2 := Coins{
		{testDenom1, sdk.NewInt(1)},
		{testDenom3, sdk.NewInt(1)},
		{testDenom2, sdk.NewInt(1)},
	}
	badAmt := Coins{
		{testDenom1, sdk.NewInt(1)},
		{testDenom3, sdk.NewInt(0)},
		{testDenom2, sdk.NewInt(1)},
	}
	dup := Coins{
		{testDenom1, sdk.NewInt(1)},
		{testDenom1, sdk.NewInt(1)},
		{testDenom2, sdk.NewInt(1)},
	}
	neg := Coins{
		{testDenom1, sdk.NewInt(-1)},
		{testDenom2, sdk.NewInt(1)},
	}

	assert.True(t, good.IsValid(), "Coins are valid")
	assert.False(t, mixedCase1.IsValid(), "Coins denoms contain upper case characters")
	assert.False(t, mixedCase2.IsValid(), "First Coins denoms contain upper case characters")
	assert.False(t, mixedCase3.IsValid(), "Single denom in Coins contains upper case characters")
	assert.True(t, good.IsAllPositive(), "Expected coins to be positive: %v", good)
	assert.False(t, empty.IsAllPositive(), "Expected coins to not be positive: %v", empty)
	assert.True(t, good.IsAllGTE(empty), "Expected %v to be >= %v", good, empty)
	assert.False(t, good.IsAllLT(empty), "Expected %v to be < %v", good, empty)
	assert.True(t, empty.IsAllLT(good), "Expected %v to be < %v", empty, good)
	assert.False(t, badSort1.IsValid(), "Coins are not sorted")
	assert.False(t, badSort2.IsValid(), "Coins are not sorted")
	assert.False(t, badAmt.IsValid(), "Coins cannot include 0 amounts")
	assert.False(t, dup.IsValid(), "Duplicate coin")
	assert.False(t, neg.IsValid(), "Negative first-denom coin")
}

func TestCoinsGT(t *testing.T) {
	one := sdk.NewInt(1)
	two := sdk.NewInt(2)

	assert.False(t, Coins{}.IsAllGT(Coins{}))
	assert.True(t, Coins{{testDenom1, one}}.IsAllGT(Coins{}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllGT(Coins{{testDenom1, one}}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllGT(Coins{{testDenom2, one}}))
	assert.True(t, Coins{{testDenom1, one}, {testDenom2, two}}.IsAllGT(Coins{{testDenom2, one}}))
	assert.False(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllGT(Coins{{testDenom2, two}}))
}

// nolint:dupl
func TestCoinsLT(t *testing.T) {
	one := sdk.NewInt(1)
	two := sdk.NewInt(2)

	assert.False(t, Coins{}.IsAllLT(Coins{}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllLT(Coins{}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllLT(Coins{{testDenom1, one}}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllLT(Coins{{testDenom2, one}}))
	assert.False(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllLT(Coins{{testDenom2, one}}))
	assert.False(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllLT(Coins{{testDenom2, two}}))
	assert.False(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllLT(Coins{{testDenom1, one}, {testDenom2, one}}))
	assert.True(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllLT(Coins{{testDenom1, two}, {testDenom2, two}}))
	assert.True(t, Coins{}.IsAllLT(Coins{{testDenom1, one}}))
}

// nolint:dupl
func TestCoinsLTE(t *testing.T) {
	one := sdk.NewInt(1)
	two := sdk.NewInt(2)

	assert.True(t, Coins{}.IsAllLTE(Coins{}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllLTE(Coins{}))
	assert.True(t, Coins{{testDenom1, one}}.IsAllLTE(Coins{{testDenom1, one}}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllLTE(Coins{{testDenom2, one}}))
	assert.False(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllLTE(Coins{{testDenom2, one}}))
	assert.False(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllLTE(Coins{{testDenom2, two}}))
	assert.True(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllLTE(Coins{{testDenom1, one}, {testDenom2, one}}))
	assert.True(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllLTE(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.True(t, Coins{}.IsAllLTE(Coins{{testDenom1, one}}))
}

func TestSortCoins(t *testing.T) {
	good := Coins{
		NewInt64Coin(testDenom1, 1),
		NewInt64Coin(testDenom2, 1),
		NewInt64Coin(testDenom3, 1),
	}
	empty := Coins{
		NewInt64Coin(testDenom4, 0),
	}
	badSort1 := Coins{
		NewInt64Coin(testDenom3, 1),
		NewInt64Coin(testDenom1, 1),
		NewInt64Coin(testDenom2, 1),
	}
	badSort2 := Coins{ // both are after the first one, but the second and third are in the wrong order
		NewInt64Coin(testDenom1, 1),
		NewInt64Coin(testDenom3, 1),
		NewInt64Coin(testDenom2, 1),
	}
	badAmt := Coins{
		NewInt64Coin(testDenom1, 1),
		NewInt64Coin(testDenom3, 0),
		NewInt64Coin(testDenom2, 1),
	}
	dup := Coins{
		NewInt64Coin(testDenom1, 1),
		NewInt64Coin(testDenom1, 1),
		NewInt64Coin(testDenom2, 1),
	}

	cases := []struct {
		coins         Coins
		before, after bool // valid before/after sort
	}{
		{good, true, true},
		{empty, false, false},
		{badSort1, false, true},
		{badSort2, false, true},
		{badAmt, false, false},
		{dup, false, false},
	}

	for tcIndex, tc := range cases {
		require.Equal(t, tc.before, tc.coins.IsValid(), "coin validity is incorrect before sorting, tc #%d", tcIndex)
		tc.coins.Sort()
		require.Equal(t, tc.after, tc.coins.IsValid(), "coin validity is incorrect after sorting, tc #%d", tcIndex)
	}
}

func TestAmountOf(t *testing.T) {
	case0 := Coins{}
	case1 := Coins{
		NewInt64Coin(testDenom4, 0),
	}
	case2 := Coins{
		NewInt64Coin(testDenom1, 1),
		NewInt64Coin(testDenom2, 1),
		NewInt64Coin(testDenom3, 1),
	}
	case3 := Coins{
		NewInt64Coin(testDenom2, 1),
		NewInt64Coin(testDenom3, 1),
	}
	case4 := Coins{
		NewInt64Coin(testDenom1, 8),
	}

	cases := []struct {
		coins           Coins
		amountOf        int64
		amountOfSpace   int64
		amountOfGAS     int64
		amountOfMINERAL int64
		amountOfTREE    int64
	}{
		{case0, 0, 0, 0, 0, 0},
		{case1, 0, 0, 0, 0, 0},
		{case2, 0, 0, 1, 1, 1},
		{case3, 0, 0, 0, 1, 1},
		{case4, 0, 0, 8, 0, 0},
	}

	for _, tc := range cases {
		assert.Equal(t, sdk.NewInt(tc.amountOfGAS), tc.coins.AmountOf(testDenom1))
		assert.Equal(t, sdk.NewInt(tc.amountOfMINERAL), tc.coins.AmountOf(testDenom2))
		assert.Equal(t, sdk.NewInt(tc.amountOfTREE), tc.coins.AmountOf(testDenom3))
	}

	assert.Panics(t, func() { cases[0].coins.AmountOf("Invalid") })
}

// nolint:dupl
func TestCoinsIsAnyGTE(t *testing.T) {
	one := sdk.NewInt(1)
	two := sdk.NewInt(2)

	assert.False(t, Coins{}.IsAnyGTE(Coins{}))
	assert.False(t, Coins{{testDenom1, one}}.IsAnyGTE(Coins{}))
	assert.False(t, Coins{}.IsAnyGTE(Coins{{testDenom1, one}}))
	assert.False(t, Coins{{testDenom1, one}}.IsAnyGTE(Coins{{testDenom1, two}}))
	assert.False(t, Coins{{testDenom1, one}}.IsAnyGTE(Coins{{testDenom2, one}}))
	assert.True(t, Coins{{testDenom1, one}, {testDenom2, two}}.IsAnyGTE(Coins{{testDenom1, two}, {testDenom2, one}}))
	assert.True(t, Coins{{testDenom1, one}}.IsAnyGTE(Coins{{testDenom1, one}}))
	assert.True(t, Coins{{testDenom1, two}}.IsAnyGTE(Coins{{testDenom1, one}}))
	assert.True(t, Coins{{testDenom1, one}}.IsAnyGTE(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.True(t, Coins{{testDenom2, two}}.IsAnyGTE(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.False(t, Coins{{testDenom2, one}}.IsAnyGTE(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.True(t, Coins{{testDenom1, one}, {testDenom2, two}}.IsAnyGTE(Coins{{testDenom1, one}, {testDenom2, one}}))
	assert.True(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAnyGTE(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.True(t, Coins{{"00000aaa00000000", one}, {"00000bbb00000000", one}}.IsAnyGTE(Coins{{testDenom2, one}, {"00000ccc00000000", one}, {"00000bbb00000000", one}, {"00000ddd00000000", one}}))
}

// nolint:dupl
func TestCoinsIsAllGT(t *testing.T) {
	one := sdk.NewInt(1)
	two := sdk.NewInt(2)

	assert.False(t, Coins{}.IsAllGT(Coins{}))
	assert.True(t, Coins{{testDenom1, one}}.IsAllGT(Coins{}))
	assert.False(t, Coins{}.IsAllGT(Coins{{testDenom1, one}}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllGT(Coins{{testDenom1, two}}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllGT(Coins{{testDenom2, one}}))
	assert.False(t, Coins{{testDenom1, one}, {testDenom2, two}}.IsAllGT(Coins{{testDenom1, two}, {testDenom2, one}}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllGT(Coins{{testDenom1, one}}))
	assert.True(t, Coins{{testDenom1, two}}.IsAllGT(Coins{{testDenom1, one}}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllGT(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.False(t, Coins{{testDenom2, two}}.IsAllGT(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.False(t, Coins{{testDenom2, one}}.IsAllGT(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.False(t, Coins{{testDenom1, one}, {testDenom2, two}}.IsAllGT(Coins{{testDenom1, one}, {testDenom2, one}}))
	assert.False(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllGT(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.False(t, Coins{{"xxx", one}, {"yyy", one}}.IsAllGT(Coins{{testDenom2, one}, {"ccc", one}, {"yyy", one}, {"zzz", one}}))
}

func TestCoinsIsAllGTE(t *testing.T) {
	one := sdk.NewInt(1)
	two := sdk.NewInt(2)

	assert.True(t, Coins{}.IsAllGTE(Coins{}))
	assert.True(t, Coins{{testDenom1, one}}.IsAllGTE(Coins{}))
	assert.True(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllGTE(Coins{{testDenom2, one}}))
	assert.False(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllGTE(Coins{{testDenom2, two}}))
	assert.False(t, Coins{}.IsAllGTE(Coins{{testDenom1, one}}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllGTE(Coins{{testDenom1, two}}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllGTE(Coins{{testDenom2, one}}))
	assert.False(t, Coins{{testDenom1, one}, {testDenom2, two}}.IsAllGTE(Coins{{testDenom1, two}, {testDenom2, one}}))
	assert.True(t, Coins{{testDenom1, one}}.IsAllGTE(Coins{{testDenom1, one}}))
	assert.True(t, Coins{{testDenom1, two}}.IsAllGTE(Coins{{testDenom1, one}}))
	assert.False(t, Coins{{testDenom1, one}}.IsAllGTE(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.False(t, Coins{{testDenom2, two}}.IsAllGTE(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.False(t, Coins{{testDenom2, one}}.IsAllGTE(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.True(t, Coins{{testDenom1, one}, {testDenom2, two}}.IsAllGTE(Coins{{testDenom1, one}, {testDenom2, one}}))
	assert.False(t, Coins{{testDenom1, one}, {testDenom2, one}}.IsAllGTE(Coins{{testDenom1, one}, {testDenom2, two}}))
	assert.False(t, Coins{{"xxx", one}, {"yyy", one}}.IsAllGTE(Coins{{testDenom2, one}, {"ccc", one}, {"yyy", one}, {"zzz", one}}))
}

func TestNewCoins(t *testing.T) {
	tenatom := NewInt64Coin(testDenom1, 10)
	tenbtc := NewInt64Coin(testDenom2, 10)
	zeroeth := NewInt64Coin(testDenom3, 0)
	tests := []struct {
		name      string
		coins     Coins
		want      Coins
		wantPanic bool
	}{
		{"empty args", []Coin{}, Coins{}, false},
		{"one coin", []Coin{tenatom}, Coins{tenatom}, false},
		{"sort after create", []Coin{tenbtc, tenatom}, Coins{tenatom, tenbtc}, false},
		{"sort and remove zeroes", []Coin{zeroeth, tenbtc, tenatom}, Coins{tenatom, tenbtc}, false},
		{"panic on dups", []Coin{tenatom, tenatom}, Coins{}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				require.Panics(t, func() { NewCoins(tt.coins...) })
				return
			}
			got := NewCoins(tt.coins...)
			require.True(t, got.IsEqual(tt.want))
		})
	}
}

func TestCoinsIsAnyGT(t *testing.T) {
	twoAtom := NewInt64Coin(testDenom1, 2)
	fiveAtom := NewInt64Coin(testDenom1, 5)
	threeEth := NewInt64Coin(testDenom2, 3)
	sixEth := NewInt64Coin(testDenom2, 6)
	twoBtc := NewInt64Coin(testDenom3, 2)

	require.False(t, Coins{}.IsAnyGT(Coins{}))

	require.False(t, Coins{fiveAtom}.IsAnyGT(Coins{}))
	require.False(t, Coins{}.IsAnyGT(Coins{fiveAtom}))
	require.True(t, Coins{fiveAtom}.IsAnyGT(Coins{twoAtom}))
	require.False(t, Coins{twoAtom}.IsAnyGT(Coins{fiveAtom}))

	require.True(t, Coins{twoAtom, sixEth}.IsAnyGT(Coins{twoBtc, fiveAtom, threeEth}))
	require.False(t, Coins{twoBtc, twoAtom, threeEth}.IsAnyGT(Coins{fiveAtom, sixEth}))
	require.False(t, Coins{twoAtom, sixEth}.IsAnyGT(Coins{twoBtc, fiveAtom}))
}

func TestFindDup(t *testing.T) {
	abc := NewInt64Coin(testDenom1, 10)
	def := NewInt64Coin(testDenom2, 10)
	ghi := NewInt64Coin(testDenom3, 10)

	type args struct {
		coins Coins
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"empty", args{NewCoins()}, -1},
		{"one coin", args{NewCoins(NewInt64Coin("00000abc00000000", 10))}, -1},
		{"no dups", args{Coins{abc, def, ghi}}, -1},
		{"dup at first position", args{Coins{abc, abc, def}}, 1},
		{"dup after first position", args{Coins{abc, def, def}}, 2},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := findDup(tt.args.coins); got != tt.want {
				t.Errorf("findDup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshalJSONCoins(t *testing.T) {
	cdc := codec.New()
	RegisterCodec(cdc)

	testCases := []struct {
		name      string
		input     Coins
		strOutput string
	}{
		{"nil coins", nil, `[]`},
		{"empty coins", Coins{}, `[]`},
		{"non-empty coins", NewCoins(NewInt64Coin("00000fee00000000", 50)), `[{"token_id":"00000fee00000000","amount":"50"}]`},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			bz, err := cdc.MarshalJSON(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.strOutput, string(bz))

			var newCoins Coins
			require.NoError(t, cdc.UnmarshalJSON(bz, &newCoins))

			if tc.input.Empty() {
				require.Nil(t, newCoins)
			} else {
				require.Equal(t, tc.input, newCoins)
			}
		})
	}
}

func TestParseCoin(t *testing.T) {
	cases := []struct {
		coinStr    string
		expectPass bool
	}{
		{"1:0000000100000000", true},
		{"1000:0000000100000000", true},
		{"21302131270312:0000000100000000", true},
		{"1:000000010000000", false},
		{"10000000100000000", false},
		{"0001:0000000100000000", true},
		{"1 : 0000000100000000", false},
		{"1:0000000a00000000", true},
		{"1:0000000A00000000", false},
	}

	for i, tc := range cases {
		_, err := ParseCoin(tc.coinStr)
		if tc.expectPass {
			require.NoError(t, err, "unexpected result for IsValid, tc #%d", i)
		} else {
			require.Error(t, err, "unexpected result for IsValid, tc #%d", i)
		}
	}
}

func TestParseCoins(t *testing.T) {
	cases := []struct {
		coinStr    string
		expectPass bool
	}{
		{"1:0000000100000000", true},
		{"1:0000000100000000,2:0000000200000000", true},
		{"1:0000000100000000,2:0000000200000000,3:0000000300000000", true},
		{"1:0000000100000000 , 2:0000000200000000     , 3:0000000300000000", true},
	}

	for i, tc := range cases {
		_, err := ParseCoins(tc.coinStr)
		if tc.expectPass {
			require.NoError(t, err, "unexpected result for IsValid, tc #%d", i)
		} else {
			require.Error(t, err, "unexpected result for IsValid, tc #%d", i)
		}
	}
}
