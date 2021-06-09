package baseapp_test

import (
	"testing"

	abci "github.com/line/ostracon/abci/types"
	ostproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/baseapp"
)

func TestValidateBlockParams(t *testing.T) {
	testCases := []struct {
		arg       interface{}
		expectErr bool
	}{
		{nil, true},
		{&abci.BlockParams{}, true},
		{abci.BlockParams{}, true},
		{abci.BlockParams{MaxBytes: -1, MaxGas: -1}, true},
		{abci.BlockParams{MaxBytes: 2000000, MaxGas: -5}, true},
		{abci.BlockParams{MaxBytes: 2000000, MaxGas: 300000}, false},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expectErr, baseapp.ValidateBlockParams(tc.arg) != nil)
	}
}

func TestValidateEvidenceParams(t *testing.T) {
	testCases := []struct {
		arg       interface{}
		expectErr bool
	}{
		{nil, true},
		{&ostproto.EvidenceParams{}, true},
		{ostproto.EvidenceParams{}, true},
		{ostproto.EvidenceParams{MaxAgeNumBlocks: -1, MaxAgeDuration: 18004000, MaxBytes: 5000000}, true},
		{ostproto.EvidenceParams{MaxAgeNumBlocks: 360000, MaxAgeDuration: -1, MaxBytes: 5000000}, true},
		{ostproto.EvidenceParams{MaxAgeNumBlocks: 360000, MaxAgeDuration: 18004000, MaxBytes: -1}, true},
		{ostproto.EvidenceParams{MaxAgeNumBlocks: 360000, MaxAgeDuration: 18004000, MaxBytes: 5000000}, false},
		{ostproto.EvidenceParams{MaxAgeNumBlocks: 360000, MaxAgeDuration: 18004000, MaxBytes: 0}, false},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expectErr, baseapp.ValidateEvidenceParams(tc.arg) != nil)
	}
}

func TestValidateValidatorParams(t *testing.T) {
	testCases := []struct {
		arg       interface{}
		expectErr bool
	}{
		{nil, true},
		{&ostproto.ValidatorParams{}, true},
		{ostproto.ValidatorParams{}, true},
		{ostproto.ValidatorParams{PubKeyTypes: []string{}}, true},
		{ostproto.ValidatorParams{PubKeyTypes: []string{"secp256k1"}}, false},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expectErr, baseapp.ValidateValidatorParams(tc.arg) != nil)
	}
}
