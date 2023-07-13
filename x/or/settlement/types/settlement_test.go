package types

import (
	fmt "fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChallengeIsSearching(t *testing.T) {
	testCases := []struct {
		challenge Challenge
		expected  bool
	}{
		{Challenge{L: 15, R: 16}, false},
		{Challenge{L: 14, R: 16}, true},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i+1), func(t *testing.T) {
			require.Equal(t, tc.expected, tc.challenge.IsSearching())
		})
	}
}

// TODO: add proper test cases.
func TestChallengeGetSteps(t *testing.T) {
	testCases := []struct {
		challenge Challenge
		expected  []uint64
	}{
		{Challenge{L: 1, R: 17}, []uint64{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
		{Challenge{L: 1, R: 10}, []uint64{2, 3, 4, 5, 6, 7, 8, 9}},
		{Challenge{L: 2, R: 10}, []uint64{3, 4, 5, 6, 7, 8, 9}},
		{Challenge{L: 2, R: 9}, []uint64{3, 4, 5, 6, 7, 8}},
		{Challenge{L: 3, R: 9}, []uint64{4, 5, 6, 7, 8}},
		{Challenge{L: 3, R: 8}, []uint64{4, 5, 6, 7}},
		{Challenge{L: 4, R: 8}, []uint64{5, 6, 7}},
		{Challenge{L: 4, R: 7}, []uint64{5, 6}},
		{Challenge{L: 5, R: 7}, []uint64{6}},
		{Challenge{L: 1, R: 20}, []uint64{2, 3, 4, 5, 6, 8, 9, 10, 11, 12, 14, 15, 16, 17, 18}},
		{Challenge{L: 1, R: 100}, []uint64{7, 13, 19, 25, 31, 38, 44, 50, 56, 62, 69, 75, 81, 87, 93}},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d", i+1), func(t *testing.T) {
			require.Equal(t, tc.expected, tc.challenge.GetSteps())
		})
	}
}
