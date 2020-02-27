package types

/*
func TestTokens_NextTokenID(t *testing.T) {
	ts := Tokens{}
	ts = ts.Append(
		&BaseFT{"link0002", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0003", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0004", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0005", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0006", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0007", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0008", sdk.NewInt(0), true, defaultName, "link0002"},
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
		&BaseFT{"link0001", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0002", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"link0003", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"cony0003", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"conyxxx3", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"conyzzzy", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"conyzzzz", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"line0001", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"line0002", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"line0003", sdk.NewInt(0), true, defaultName, "link0002"},
		&BaseFT{"linezzzz", sdk.NewInt(0), true, defaultName, "link0002"},
	)

	require.Equal(t, 3, ts.GetTokens("link").Len())
	require.Equal(t, 4, ts.GetTokens("cony").Len())
	require.Equal(t, 7, ts.GetTokens("li").Len())
	require.Equal(t, 7, ts.GetTokens("lin").Len())
}

*/
