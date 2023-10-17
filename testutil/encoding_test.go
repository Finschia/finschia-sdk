package testutil_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/testutil"
	sdk "github.com/Finschia/finschia-sdk/types"
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
	testCases := map[string]struct {
		acceptedType any
	}{
		"string": {
			acceptedType: "test",
		},
		"sdk.AccAddress": {
			acceptedType: sdk.AccAddress("address"),
		},
		"sdk.Coin": {
			acceptedType: sdk.NewInt(1),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			switch v := tc.acceptedType.(type) {
			case string:
				require.Equal(t, []byte(fmt.Sprintf(`"%s"`, v)), testutil.W(v))
			case fmt.Stringer:
				require.Equal(t, []byte(fmt.Sprintf(`"%s"`, v.String())), testutil.W(v))
			default:
				t.Fatalf("not supported types")
			}
		})
	}
}
