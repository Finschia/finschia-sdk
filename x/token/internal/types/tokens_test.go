package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestNextTokenID(t *testing.T) {
	require.Equal(t, "b", NextTokenID("a", ""))
	require.Equal(t, "0001", NextTokenID("0000", ""))
	require.Equal(t, "000a", NextTokenID("0009", ""))
	require.Equal(t, "0010", NextTokenID("000z", ""))
	require.Equal(t, "0000", NextTokenID("zzzz", ""))
	require.Equal(t, "00000000", NextTokenID("zzzzzzzz", ""))
	require.Equal(t, "abce0000", NextTokenID("abcdzzzz", ""))
	require.Equal(t, "abcdabc1", NextTokenID("abcdabc0", ""))

	require.Equal(t, "", NextTokenID("", ""))
	require.Equal(t, "", NextTokenID("", "zzzzz"))
	require.Equal(t, "z0", NextTokenID("zz", "z"))
	require.Equal(t, "item0001", NextTokenID("item0000", "item"))
	require.Equal(t, "item0010", NextTokenID("item000z", "item"))
	require.Equal(t, "itemyyz0", NextTokenID("itemyyyz", "item"))
	require.Equal(t, "itemyz00", NextTokenID("itemyyzz", "item"))
	require.Equal(t, "item999a", NextTokenID("item9999", "item"))
	require.Equal(t, "item99a0", NextTokenID("item999z", "item"))
	require.Equal(t, "z0000000", NextTokenID("zzzzzzzz", "z"))
	require.Equal(t, "zz000000", NextTokenID("zzzzzzzz", "zz"))
	require.Equal(t, "zzzzzzz0", NextTokenID("zzzzzzzz", "zzzzzzz"))
	require.Equal(t, "zzzzzzzz", NextTokenID("zzzzzzzz", "zzzzzzzz"))
	require.Equal(t, "item0000", NextTokenID("itemzzzz", "item"))
	require.Equal(t, "itemz000", NextTokenID("itemyzzz", "item"))
	require.Equal(t, "item0000", NextTokenID("itezzzzz", "item"))

	nextID := "0000"
	for idx := 0; idx < 36*36*36*36; idx++ {
		nextID = NextTokenID(nextID, "")
	}
	require.Equal(t, "0000", nextID)
}

func TestTokens_NextTokenID(t *testing.T) {
	ts := Tokens{}
	collection := NewCollection("link", "linktest")
	ts = ts.Append(
		NewCollectiveFT(collection, defaultName, "link0001", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "link0002", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "link0003", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "link0004", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "link0005", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "link0006", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "link0007", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "link0008", defaultTokenURI, sdk.NewInt(0), true),
	)
	require.Equal(t, "link0009", ts.NextTokenID(""))
	require.Equal(t, "link0009", ts.NextTokenID("link"))

	require.Equal(t, "", ts.NextTokenID("1234567890"))
	require.Equal(t, "linl", ts.NextBaseID())
	require.Equal(t, NextTokenID("link", ""), ts.NextBaseID())

}

func TestTokens_Iterate(t *testing.T) {

	ts := Tokens{}
	require.Equal(t, 0, ts.GetTokens("link").Len())
	require.Equal(t, 0, ts.GetTokens("").Len())
	collection := NewCollection("link", "linktest")
	ts = ts.Append(
		NewCollectiveFT(collection, defaultName, "link0001", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "link0002", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "link0003", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "cony0003", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "conyxxx3", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "conyzzzy", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "conyzzzz", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "line0001", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "line0002", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "line0003", defaultTokenURI, sdk.NewInt(0), true),
		NewCollectiveFT(collection, defaultName, "linezzzz", defaultTokenURI, sdk.NewInt(0), true),
	)

	require.Equal(t, 3, ts.GetTokens("link").Len())
	require.Equal(t, 4, ts.GetTokens("cony").Len())
	require.Equal(t, 7, ts.GetTokens("li").Len())
	require.Equal(t, 7, ts.GetTokens("lin").Len())
}
