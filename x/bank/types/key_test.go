package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/bank/types"
)

func cloneAppend(bz []byte, tail []byte) (res []byte) {
	res = make([]byte, len(bz)+len(tail))
	copy(res, bz)
	copy(res[len(bz):], tail)
	return
}

func TestAddressFromBalancesStore(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("link19tzp7e489drh9qfs9m84k2qe5a5yyknzen48tz")
	require.NoError(t, err)

	key := cloneAppend(addr.Bytes(), []byte("stake"))
	res := types.AddressFromBalancesStore(key)
	require.Equal(t, res, addr)
}
