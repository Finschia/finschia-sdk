package collection_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	"github.com/Finschia/finschia-sdk/x/collection"
)

const contractID = "deadbeef"

func TestNFTClass(t *testing.T) {
	nextIDs := collection.DefaultNextClassIDs(contractID)
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
			var class collection.TokenClass = &collection.NFTClass{}
			class.SetID(&nextIDs)
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
			input:    "00bab10c00000001:1",
			valid:    true,
			expected: collection.NewNFTCoin("00bab10c", 1),
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
			input: "00bab10c00000001:1",
			valid: true,
			expected: collection.NewCoins(
				collection.NewNFTCoin("00bab10c", 1),
			),
		},
		"valid multiple coins": {
			input: "deadbeef00000001:1,deadbeef0000000a:1",
			valid: true,
			expected: collection.NewCoins(
				collection.NewNFTCoin(contractID, 1),
				collection.NewNFTCoin(contractID, 10),
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

func TestDefaultNextClassIDs(t *testing.T) {
	require.Equal(t, collection.NextClassIDs{
		ContractId:  contractID,
		NonFungible: math.NewUint(1 << 28).Incr(), // "10000000 + 1"
	},
		collection.DefaultNextClassIDs(contractID),
	)
}
