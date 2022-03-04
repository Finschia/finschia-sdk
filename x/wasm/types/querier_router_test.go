package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
)

func testQuerierHandler(ctx sdk.Context, jsonQuerier json.RawMessage) ([]byte, error) {
	return nil, nil
}

func TestQuerierRouterSeal(t *testing.T) {
	r := NewQuerierRouter()
	r.Seal()
	require.Panics(t, func() { r.AddRoute("test", nil) })
	require.Panics(t, func() { r.Seal() })
}

func TestQuerierRouter(t *testing.T) {
	r := NewQuerierRouter()
	r.AddRoute("test", testQuerierHandler)
	require.True(t, r.HasRoute("test"))
	require.Panics(t, func() { r.AddRoute("test", testQuerierHandler) })
	require.Panics(t, func() { r.AddRoute("    ", testQuerierHandler) })
}
