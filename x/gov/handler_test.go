package gov_test

import (
	"strings"
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/line/lbm-sdk/v2/testutil/testdata"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/gov"
	"github.com/line/lbm-sdk/v2/x/gov/keeper"
)

func TestInvalidMsg(t *testing.T) {
	k := keeper.Keeper{}
	h := gov.NewHandler(k)

	res, err := h(sdk.NewContext(nil, tmproto.Header{}, false, nil), testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognized gov message type"))
}
