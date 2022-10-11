package types_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/x/ibc/core/02-client/types"
)

// tests ParseClientIdentifier and IsValidClientID
func TestParseClientIdentifier(t *testing.T) {
	testCases := []struct {
		name       string
		clientID   string
		clientType string
		expSeq     uint64
		expPass    bool
	}{
		{"valid 0", "ostracon-0", "ostracon", 0, true},
		{"valid 1", "ostracon-1", "ostracon", 1, true},
		{"valid solemachine", "solomachine-v1-1", "solomachine-v1", 1, true},
		{"valid large sequence", types.FormatClientIdentifier("ostracon", math.MaxUint64), "ostracon", math.MaxUint64, true},
		{"valid short client type", "t-0", "t", 0, true},
		// one above uint64 max
		{"invalid uint64", "ostracon-18446744073709551616", "ostracon", 0, false},
		// uint64 == 20 characters
		{"invalid large sequence", "ostracon-2345682193567182931243", "ostracon", 0, false},
		{"invalid newline in clientID", "ostraco\nn-1", "ostraco\nn", 0, false},
		{"invalid newline character before dash", "ostracon\n-1", "ostracon", 0, false},
		{"missing dash", "ostracon0", "ostracon", 0, false},
		{"blank id", "               ", "    ", 0, false},
		{"empty id", "", "", 0, false},
		{"negative sequence", "ostracon--1", "ostracon", 0, false},
		{"invalid format", "ostracon-tm", "ostracon", 0, false},
		{"empty clientype", " -100", "ostracon", 0, false},
	}

	for _, tc := range testCases {

		clientType, seq, err := types.ParseClientIdentifier(tc.clientID)
		valid := types.IsValidClientID(tc.clientID)
		require.Equal(t, tc.expSeq, seq, tc.clientID)

		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.True(t, valid)
			require.Equal(t, tc.clientType, clientType)
		} else {
			require.Error(t, err, tc.name, tc.clientID)
			require.False(t, valid)
			require.Equal(t, "", clientType)
		}
	}
}
