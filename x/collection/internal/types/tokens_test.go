package types

/*
func TestTokens_NextTokenID(t *testing.T) {
	ts := Tokens{}
	ts = ts.Append(
		&BaseFT{"link0002", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0003", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0004", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0005", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0006", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0007", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0008", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
	)
	require.Equal(t, "link0009", ts.NextTokenID(""))
	require.Equal(t, "link0009", ts.NextTokenID("link"))

	require.Equal(t, "", ts.NextTokenID("1234567890"))
	require.Equal(t, "linl", ts.NextTokenTypeForNFT())
	require.Equal(t, NextID("link", ""), ts.NextTokenTypeForNFT())
}
func TestTokens_Iterate(t *testing.T) {
	ts := Tokens{}
	ts = ts.Append(
		&BaseFT{"link0001", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0002", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0003", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"cony0003", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"conyxxx3", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"conyzzzy", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"conyzzzz", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"line0001", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"line0002", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"line0003", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"linezzzz", defaultTokenURI, sdk.NewInt(0), true, defaultName, "link0002"},
	)

	require.Equal(t, 3, ts.GetTokens("link").Len())
	require.Equal(t, 4, ts.GetTokens("cony").Len())
	require.Equal(t, 7, ts.GetTokens("li").Len())
	require.Equal(t, 7, ts.GetTokens("lin").Len())
}

*/
