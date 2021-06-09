package types

import (
	"encoding/json"
	"testing"

	sdk "github.com/line/lfb-sdk/types"
	"github.com/stretchr/testify/require"
)

func testHandler(jsonMsg json.RawMessage) ([]sdk.Msg, error) { return []sdk.Msg{}, nil }

func TestRouterSeal(t *testing.T) {
	r := NewRouter()
	r.Seal()
	require.Panics(t, func() { r.AddRoute("test", nil) })
	require.Panics(t, func() { r.Seal() })
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	r.AddRoute("test", testHandler)
	require.True(t, r.HasRoute("test"))
	require.Panics(t, func() { r.AddRoute("test", testHandler) })
	require.Panics(t, func() { r.AddRoute("    ", testHandler) })
}
