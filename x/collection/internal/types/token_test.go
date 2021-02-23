package types

import (
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalFT(t *testing.T) {
	// Given a FT
	token := NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(defaultDecimals), true)
	var token2 BaseFT

	// When marshal and unmarshal the FT
	bz, err := ModuleCdc.MarshalJSON(token)
	require.NoError(t, err)
	err = ModuleCdc.UnmarshalJSON(bz, &token2)
	require.NoError(t, err)

	// Then the properties are same
	r := require.New(t)
	r.EqualValues(defaultName, token.GetName(), token2.GetName())
	r.Equal(defaultContractID, token.GetContractID(), token2.GetContractID())
	r.Equal(defaultTokenIDFT, token.GetTokenID(), token2.GetTokenID())
	r.Equal(defaultTokenIDFT[:TokenTypeLength], token.GetTokenType(), token2.GetTokenType())
	r.Equal(defaultTokenIDFT[TokenTypeLength:], token.GetTokenIndex(), token2.GetTokenIndex())
	r.Equal(int64(defaultDecimals), token.GetDecimals().Int64(), token2.GetDecimals().Int64())
	r.Equal(true, token.GetMintable(), token2.GetMintable())

	r.Equal(`{"contract_id":"abcdef01","token_id":"0000000100000000","decimals":"6","mintable":true,"name":"name","meta":"{}"}`, token.String())
}

func TestUnmarshalNFT(t *testing.T) {
	// Given a NFT
	token := NewNFT(defaultContractID, defaultTokenID1, defaultName, defaultMeta, addr1)
	var token2 BaseNFT

	// When marshal and unmarshal the FT
	bz, err := ModuleCdc.MarshalJSON(token)
	require.NoError(t, err)
	err = ModuleCdc.UnmarshalJSON(bz, &token2)
	require.NoError(t, err)

	// Then the properties are same
	r := require.New(t)
	r.Equal(defaultName, token.GetName(), token2.GetName())
	r.Equal(defaultContractID, token.GetContractID(), token2.GetContractID())
	r.Equal(defaultTokenID1, token.GetTokenID(), token2.GetTokenID())
	r.Equal(defaultTokenID1[:TokenTypeLength], token.GetTokenType(), token2.GetTokenType())
	r.Equal(defaultTokenID1[TokenTypeLength:], token.GetTokenIndex(), token2.GetTokenIndex())
	r.Equal(addr1, token.GetOwner(), token2.GetOwner())
}

func TestSetName(t *testing.T) {
	// Given a FT, NFT
	tokenFT := NewFT(defaultContractID, defaultTokenIDFT, defaultName, defaultMeta, sdk.NewInt(defaultDecimals), true)
	tokenNFT := NewNFT(defaultContractID, defaultTokenID1, defaultName, defaultMeta, addr1)

	tokenFT.SetName("new_name")
	tokenNFT.SetName("new_name")
	tokenFT.SetMeta("new_meta")
	tokenNFT.SetMeta("new_meta")

	// When change name, Then they are changed
	require.Equal(t, "new_name", tokenFT.GetName())
	require.Equal(t, "new_name", tokenNFT.GetName())
	require.Equal(t, "new_meta", tokenFT.GetMeta())
	require.Equal(t, "new_meta", tokenNFT.GetMeta())

	// Set empty name
	tokenFT.SetName("")
	tokenNFT.SetName("")

	require.Equal(t, "", tokenFT.GetName())
	require.Equal(t, "", tokenNFT.GetName())
}
