package collection_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

func TestFTClass(t *testing.T) {
	nextIDs := collection.DefaultNextClassIDs("deadbeef")
	testCases := map[string]struct {
		id       string
		name     string
		meta     string
		decimals int32
		valid    bool
	}{
		"valid class": {
			valid: true,
		},
		"invalid id": {
			id: "invalid",
		},
		"invalid name": {
			name: string(make([]rune, 21)),
		},
		"invalid meta": {
			meta: string(make([]rune, 1001)),
		},
		"invalid decimals": {
			decimals: 19,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var class collection.TokenClass
			class = &collection.FTClass{
				Id:       tc.id,
				Decimals: tc.decimals,
			}

			if len(tc.id) == 0 {
				class.SetId(&nextIDs)
			}
			class.SetName(tc.name)
			class.SetMeta(tc.meta)

			err := class.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNFTClass(t *testing.T) {
	nextIDs := collection.DefaultNextClassIDs("deadbeef")
	testCases := map[string]struct {
		name  string
		meta  string
		valid bool
	}{
		"valid class": {
			valid: true,
		},
		"invalid name": {
			name: string(make([]rune, 21)),
		},
		"invalid meta": {
			meta: string(make([]rune, 1001)),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var class collection.TokenClass
			class = &collection.NFTClass{}
			class.SetId(&nextIDs)
			class.SetName(tc.name)
			class.SetMeta(tc.meta)

			err := class.ValidateBasic()
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestParseCoin(t *testing.T) {
	testCases := map[string]struct {
		input    string
		valid    bool
		expected collection.Coin
	}{
		"valid coin": {
			input:    "00bab10c00000000:10",
			valid:    true,
			expected: collection.NewFTCoin("00bab10c", sdk.NewInt(10)),
		},
		"invalid expression": {
			input: "oobabloc00000000:10",
		},
		"invalid amount": {
			input: "00bab10c00000000:" + fmt.Sprintf("1%0127d", 0),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			parsed, err := collection.ParseCoin(tc.input)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tc.expected, *parsed)
			require.Equal(t, tc.input, parsed.String())
		})
	}
}

func TestParseCoins(t *testing.T) {
	testCases := map[string]struct {
		input    string
		valid    bool
		expected collection.Coins
	}{
		"valid single coins": {
			input: "00bab10c00000000:10",
			valid: true,
			expected: collection.NewCoins(
				collection.NewFTCoin("00bab10c", sdk.NewInt(10)),
			),
		},
		"valid multiple coins": {
			input: "deadbeef00000001:1,00bab10c00000000:10",
			valid: true,
			expected: collection.NewCoins(
				collection.NewNFTCoin("deadbeef", 1),
				collection.NewFTCoin("00bab10c", sdk.NewInt(10)),
			),
		},
		"empty string": {},
		"invalid coin": {
			input: "oobabloc00000000:10",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			parsed, err := collection.ParseCoins(tc.input)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tc.expected, parsed)
			require.Equal(t, tc.input, parsed.String())
		})
	}
}
