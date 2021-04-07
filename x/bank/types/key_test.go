package types_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/x/bank/keeper"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/address"
	"github.com/line/lbm-sdk/x/bank/types"
)

func cloneAppend(bz []byte, tail []byte) (res []byte) {
	res = make([]byte, len(bz)+len(tail))
	copy(res, bz)
	copy(res[len(bz):], tail)
	return
}

func TestAddressFromBalancesStore(t *testing.T) {
	addr := sdk.AccAddress("link19tzp7e489drh9qfs9m84k2qe5a5yyknzen48tz")
	err := sdk.ValidateAccAddress(addr.String())
	require.NoError(t, err)

	key := cloneAppend(address.MustLengthPrefix(addr.Bytes()), []byte("stake"))
	res, err := types.AddressFromBalancesStore(key)
	require.NoError(t, err)
	require.Equal(t, res, addr)
}

func TestInvalidAddressFromBalancesStore(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
	}{
		{"empty", []byte("")},
		{"invalid", []byte("3AA")},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := types.AddressFromBalancesStore(tc.key)
			assert.Error(t, err)
			assert.True(t, errors.Is(types.ErrInvalidKey, err))
		})
	}
}
