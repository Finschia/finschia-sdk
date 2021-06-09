package gov_test

import (
	"strings"
	"testing"

	ostproto "github.com/line/ostracon/proto/ostracon/types"

	"github.com/line/lfb-sdk/testutil/testdata"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/gov"
	"github.com/line/lfb-sdk/x/gov/keeper"
)

func TestInvalidMsg(t *testing.T) {
	k := keeper.Keeper{}
	h := gov.NewHandler(k)

	res, err := h(sdk.NewContext(nil, ostproto.Header{}, false, nil), testdata.NewTestMsg())
	require.Error(t, err)
	require.Nil(t, res)
	require.True(t, strings.Contains(err.Error(), "unrecognized gov message type"))
}
