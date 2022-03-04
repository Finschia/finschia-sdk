package types_test

import (
	"testing"

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
	res := types.AddressFromBalancesStore(key)
	require.Equal(t, res, addr)
}
