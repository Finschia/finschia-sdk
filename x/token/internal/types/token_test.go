package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalToken(t *testing.T) {
	// Given a token
	token := NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)
	var token2 Token

	// When marshal and unmarshal the token
	bz, err := ModuleCdc.MarshalJSON(token)
	require.NoError(t, err)
	err = ModuleCdc.UnmarshalJSON(bz, &token2)
	require.NoError(t, err)

	// Then the properties are same
	r := require.New(t)
	r.EqualValues(defaultName, token.GetName(), token2.GetName())
	r.Equal(defaultSymbol, token.GetSymbol(), token2.GetSymbol())
	r.Equal(defaultTokenURI, token.GetTokenURI(), token2.GetTokenURI())
	r.Equal(int64(defaultDecimals), token.GetDecimals().Int64(), token2.GetDecimals().Int64())
	r.Equal(true, token.GetMintable(), token2.GetMintable())
}

func TestSetToken(t *testing.T) {
	// Given a token
	token := NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)

	// When change name and test uri, Then they are changed
	require.Equal(t, "new_name", token.SetName("new_name").GetName())
	require.Equal(t, "new_token_uri", token.SetTokenURI("new_token_uri").GetTokenURI())
}

func TestBaseToken_String(t *testing.T) {
	token := NewToken(defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(defaultDecimals), true)

	require.Equal(t, `{"name":"name","symbol":"linktkn","token_uri":"token-uri","decimals":"6","mintable":true}`, token.String())
}

func TestTokensString(t *testing.T) {
	tokens := Tokens{
		NewToken(defaultName, defaultSymbol+"1", defaultTokenURI, sdk.NewInt(defaultDecimals), true),
		NewToken(defaultName, defaultSymbol+"2", defaultTokenURI, sdk.NewInt(defaultDecimals), true),
	}

	require.Equal(t, `[{"name":"name","symbol":"linktkn1","token_uri":"token-uri","decimals":"6","mintable":true},{"name":"name","symbol":"linktkn2","token_uri":"token-uri","decimals":"6","mintable":true}]`, tokens.String())
}
