package testutil_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/testutil"
)

func TestMustJSONMarshal(t *testing.T) {
	type tc struct {
		Name  string `json:"myName"`
		Order string `json:"myOrder"`
	}

	a := tc{
		Name:  "test",
		Order: "first",
	}
	b := new(tc)

	marshaled := testutil.MustJSONMarshal(a)
	err := json.Unmarshal(marshaled, b)
	require.NoError(t, err)
	require.Equal(t, a, *b)
	require.Panics(t, func() { testutil.MustJSONMarshal(make(chan int)) })
}

func TestW(t *testing.T) {
	require.Equal(t, []byte(`"test"`), testutil.W("test"))
}
