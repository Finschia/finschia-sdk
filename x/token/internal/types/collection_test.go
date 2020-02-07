package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func getPreparedCollection(symbol, name string, t *testing.T) Collection {
	var err error
	collection := NewCollection(symbol, name)

	for idx := 0; idx < 10; idx++ {
		collection, err = collection.AddToken(NewCollectiveFT(collection, defaultName, defaultTokenURI, sdk.NewInt(0), true))
		require.NoError(t, err)
	}
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	collection, err = collection.AddToken(NewCollectiveNFT(collection, defaultName, "link", defaultTokenURI, defaultAddr))
	require.NoError(t, err)
	return collection
}

func TestCollection_UpdateToken(t *testing.T) {
	collection := NewCollection(defaultSymbol, defaultSymbol)
	token := NewCollectiveFT(collection, defaultName, defaultTokenURI, sdk.NewInt(0), true)
	collection, err := collection.AddToken(token)
	require.NoError(t, err)

	token.SetTokenURI("changed")
	collection, err = collection.UpdateToken(token)
	require.NoError(t, err)

	token2, err := collection.GetToken(token.GetTokenID())
	require.NoError(t, err)
	require.Equal(t, "changed", token2.GetTokenURI())
}

func TestCollection_DeleteToken(t *testing.T) {
	collection := NewCollection(defaultSymbol, defaultSymbol)
	token := NewCollectiveFT(collection, defaultName, defaultTokenURI, sdk.NewInt(0), true)
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
	require.Equal(t, "linl", collection.NextTokenTypeNFT())
	require.Equal(t, "000b", collection.NextTokenTypeFT())
	require.Equal(t, "link0009", collection.NextTokenID(""))
	require.Equal(t, "linl0001", collection.NextTokenID(collection.NextTokenTypeNFT()))
}

func TestCollection_GetAllTokens(t *testing.T) {
	collection := getPreparedCollection(defaultSymbol, defaultName, t)
	require.Equal(t, 10, collection.GetFTokens().Len())
	require.Equal(t, 8, collection.GetNFTokens().Len())
	require.Equal(t, 18, collection.GetAllTokens().Len())
	require.Equal(t, collection.GetAllTokens().Len(), collection.GetFTokens().Len()+collection.GetNFTokens().Len())
}

func TestCollection_String(t *testing.T) {
	var collections []Collection
	collections = append(collections, getPreparedCollection(defaultSymbol, defaultName, t))
	collections = append(collections, getPreparedCollection(defaultSymbol+"1", defaultName, t))

	_, err := ModuleCdc.MarshalJSON(collections)
	require.NoError(t, err)
}
