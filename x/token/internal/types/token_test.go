package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalToken(t *testing.T) {
	// Given a token
	token := NewToken(defaultContractID, defaultName, defaultSymbol, defaultImageURI, sdk.NewInt(defaultDecimals), true)
	var token2 Token

	// When marshal and unmarshal the token
	bz, err := ModuleCdc.MarshalJSON(token)
	require.NoError(t, err)
	err = ModuleCdc.UnmarshalJSON(bz, &token2)
	require.NoError(t, err)

	// Then the properties are same
	r := require.New(t)
	r.EqualValues(defaultName, token.GetName(), token2.GetName())
	r.Equal(defaultContractID, token.GetContractID(), token2.GetContractID())
	r.Equal(defaultImageURI, token.GetImageURI(), token2.GetImageURI())
	r.Equal(int64(defaultDecimals), token.GetDecimals().Int64(), token2.GetDecimals().Int64())
	r.Equal(true, token.GetMintable(), token2.GetMintable())
}

func TestSetToken(t *testing.T) {
	// Given a token
	token := NewToken(defaultContractID, defaultName, defaultSymbol, defaultImageURI, sdk.NewInt(defaultDecimals), true)

	// When change name and test uri, Then they are changed
	require.Equal(t, "new_name", token.SetName("new_name").GetName())
	require.Equal(t, "new_token_uri", token.SetImageURI("new_token_uri").GetImageURI())
}

func TestBaseToken_String(t *testing.T) {
	token := NewToken(defaultContractID, defaultName, defaultSymbol, defaultImageURI, sdk.NewInt(defaultDecimals), true)

	require.Equal(t, `{"contract_id":"linktkn","name":"name","symbol":"BTC","image_uri":"image-uri","decimals":"6","mintable":true}`, token.String())
}

func TestTokensString(t *testing.T) {
	tokens := Tokens{
		NewToken(defaultContractID+"1", defaultName, defaultSymbol, defaultImageURI, sdk.NewInt(defaultDecimals), true),
		NewToken(defaultContractID+"2", defaultName, defaultSymbol, defaultImageURI, sdk.NewInt(defaultDecimals), true),
	}

	require.Equal(t, `[{"contract_id":"linktkn1","name":"name","symbol":"BTC","image_uri":"image-uri","decimals":"6","mintable":true},{"contract_id":"linktkn2","name":"name","symbol":"BTC","image_uri":"image-uri","decimals":"6","mintable":true}]`, tokens.String())
}
