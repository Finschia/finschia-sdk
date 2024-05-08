package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSwapKey(t *testing.T) {
	tests := []struct {
		name        string
		fromDenom   string
		toDenom     string
		expectedKey []byte
	}{
		{
			name:        "swapKey",
			fromDenom:   "cony",
			toDenom:     "peb",
			expectedKey: []byte{0x1, 0x4, 0x63, 0x6f, 0x6e, 0x79, 0x3, 0x70, 0x65, 0x62},
			// expectedKey: append(swapPrefix, append(append([]byte{byte(len("cony"))}, []byte("cony")...), append([]byte{byte(len("peb"))}, []byte("peb")...)...)...),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualKey := swapKey(tc.fromDenom, tc.toDenom)
			require.Equal(t, tc.expectedKey, actualKey)
		})
	}
}

func TestSwappedKey(t *testing.T) {
	tests := []struct {
		name        string
		fromDenom   string
		toDenom     string
		expectedKey []byte
	}{
		{
			name:        "swappedKey",
			fromDenom:   "cony",
			toDenom:     "peb",
			expectedKey: []byte{0x3, 0x4, 0x63, 0x6f, 0x6e, 0x79, 0x3, 0x70, 0x65, 0x62},
			// expectedKey: append(swappedKeyPrefix, append(append([]byte{byte(len("cony"))}, []byte("cony")...), append([]byte{byte(len("peb"))}, []byte("peb")...)...)...),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualKey := swappedKey(tc.fromDenom, tc.toDenom)
			require.Equal(t, tc.expectedKey, actualKey)
		})
	}
}
