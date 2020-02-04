package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func getPreparedCollection(symbol, name string, t *testing.T) Collection {
	collection := NewCollection(symbol, name)

	collection, err := collection.AddToken(NewCollectiveFT(collection, defaultName, "0ink0001", defaultTokenURI, sdk.NewInt(0), true))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveFT(collection, defaultName, "0ink0003", defaultTokenURI, sdk.NewInt(0), true))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveFT(collection, defaultName, "0ink0004", defaultTokenURI, sdk.NewInt(0), true))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveFT(collection, defaultName, "0ink0005", defaultTokenURI, sdk.NewInt(0), true))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveFT(collection, defaultName, "0ink0006", defaultTokenURI, sdk.NewInt(0), true))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveFT(collection, defaultName, "0ink0007", defaultTokenURI, sdk.NewInt(0), true))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveFT(collection, defaultName, "0ink0008", defaultTokenURI, sdk.NewInt(0), true))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveFT(collection, defaultName, "0lnk0009", defaultTokenURI, sdk.NewInt(0), true))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link0000", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link0001", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link0002", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link0003", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link0004", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link0005", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link0006", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link0007", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link0008", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	return collection
}

func TestCollection_AddToken(t *testing.T) {
	collection := NewCollection(defaultSymbol, defaultSymbol)
	collection, err := collection.AddToken(NewCollectiveFT(collection, defaultName, "00000001", defaultTokenURI, sdk.NewInt(0), true))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveFT(collection, defaultName, "00000001", defaultTokenURI, sdk.NewInt(0), true))
	require.Error(t, err)
}

func TestCollection_UpdateToken(t *testing.T) {
	collection := NewCollection(defaultSymbol, defaultSymbol)
	token := NewCollectiveFT(collection, defaultName, "00000001", defaultTokenURI, sdk.NewInt(0), true)
	collection, err := collection.AddToken(token)
	require.NoError(t, err)

	token.SetTokenURI("changed")
	collection, err = collection.UpdateToken(token)
	require.NoError(t, err)

	token2, err := collection.GetToken("00000001")
	require.NoError(t, err)
	require.Equal(t, "changed", token2.GetTokenURI())

	token = NewCollectiveFT(collection, defaultName, "00000002", defaultTokenURI, sdk.NewInt(0), true)
	collection, err = collection.UpdateToken(token)
	require.Error(t, err)
}

func TestCollection_DeleteToken(t *testing.T) {
	collection := NewCollection(defaultSymbol, defaultSymbol)
	token := NewCollectiveFT(collection, defaultName, "00000001", defaultTokenURI, sdk.NewInt(0), true)
	collection, err := collection.AddToken(token)
	require.NoError(t, err)
	collection, err = collection.DeleteToken(token)
	require.NoError(t, err)
	token2, err := collection.GetToken(token.GetTokenID())
	require.Error(t, err)
	require.Nil(t, token2)
	require.False(t, collection.HasToken(token.GetTokenID()))
}

func TestCollection_NextTokenID(t *testing.T) {
	collection := getPreparedCollection(defaultSymbol, defaultName, t)
	require.Equal(t, "linl", collection.NextBaseID())
	require.Equal(t, "0lnk000a", collection.NexTokenIDForFT())
	require.Equal(t, "link0009", collection.NextTokenID(""))
	require.Equal(t, "linl0000", collection.NextTokenID(collection.NextBaseID()))
}

func TestCollection_GetAllTokens(t *testing.T) {
	collection := getPreparedCollection(defaultSymbol, defaultName, t)
	require.Equal(t, 8, collection.GetFTokens().Len())
	require.Equal(t, 9, collection.GetNFTokens().Len())
	require.Equal(t, 17, collection.GetAllTokens().Len())
	require.Equal(t, collection.GetAllTokens().Len(), collection.GetFTokens().Len()+collection.GetNFTokens().Len())
}

func TestCollection_String(t *testing.T) {
	var collections []Collection
	collections = append(collections, getPreparedCollection(defaultSymbol, defaultName, t))
	collections = append(collections, getPreparedCollection(defaultSymbol+"1", defaultName, t))

	bz, err := ModuleCdc.MarshalJSON(collections)
	require.NoError(t, err)
	t.Log(string(bz))
}
