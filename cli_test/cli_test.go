// +build cli_test

package clitest

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	collectionmodule "github.com/line/link/x/collection"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestModifyToken(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	const (
		contractID = "9be17165"
		firstName  = "itisbrown"
		name       = "description"
		meta       = "{}"
		symbol     = "BTC"
		amount     = 10000
		decimals   = 6
	)

	barAddr := f.KeyAddress(keyBar)
	fooAddr := f.KeyAddress(keyFoo)
	// Given user
	sendTokens := sdk.TokensFromConsensusPower(1)
	f.LogResult(f.TxSend(keyFoo, barAddr, sdk.NewCoin(denom, sendTokens), "-y"))
	// And token
	f.TxTokenIssue(keyFoo, fooAddr, name, meta, symbol, amount, decimals, true, "-y")
	tests.WaitForNextNBlocksTM(1, f.Port)
	firstResult := f.QueryToken(contractID)
	require.Equal(t, name, firstResult.GetName())

	t.Log("Modify token")
	{
		// When modify token name
		modifiedName := firstName + "modified"
		f.LogResult(f.TxTokenModify(keyFoo, contractID, "name", modifiedName, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// Then the name is modified
		require.Equal(t, modifiedName, f.QueryToken(contractID).GetName())
	}
}

func TestModifyCollection(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	const (
		contractID   = "9be17165"
		tokenType    = "10000001"
		tokenIndex   = "00000001"
		tokenID      = tokenType + tokenIndex
		tokenTypeFT  = "00000001"
		tokenIndexFT = "00000000"
		tokenIDFT    = tokenTypeFT + tokenIndexFT
		firstName    = "itisbrown"
		firstMeta    = "{}"
		firstBaseURI = "uri:itisbrown"
		amount       = 10000
		decimals     = 6
	)

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)
	// Given user
	sendTokens := sdk.TokensFromConsensusPower(1)
	f.LogResult(f.TxSend(keyFoo, barAddr, sdk.NewCoin(denom, sendTokens), "-y"))
	// And collection
	f.TxTokenCreateCollection(keyFoo, firstName, firstMeta, firstBaseURI, "-y")
	tests.WaitForNextNBlocksTM(1, f.Port)
	// And FT
	f.LogResult(f.TxTokenIssueFTCollection(keyFoo, contractID, fooAddr, firstName, firstMeta, amount, decimals, true, "-y"))
	tests.WaitForNextNBlocksTM(1, f.Port)
	// And NFT
	f.LogResult(f.TxTokenIssueNFTCollection(keyFoo, contractID, firstName, firstMeta, "-y"))
	tests.WaitForNextNBlocksTM(1, f.Port)
	mintParam := strings.Join([]string{tokenType, firstName, firstMeta}, ":")
	f.LogResult(f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam, "-y"))
	tests.WaitForNextNBlocksTM(1, f.Port)
	firstResult := f.QueryTokenCollection(contractID, tokenID).(collectionmodule.NFT)
	require.Equal(t, tokenID, firstResult.GetTokenID())

	t.Log("Modify collection")
	{
		// When modify collection uri
		modifiedURI := firstBaseURI + "modified"
		f.LogResult(f.TxCollectionModify(keyFoo, contractID, "", "", "base_img_uri", modifiedURI, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// Then the uri is modified
		require.Equal(t, modifiedURI, f.QueryCollection(contractID).GetBaseImgURI())
	}
	t.Log("Modify meta")
	{
		// When modify meta
		modifiedMeta := firstMeta + "modified"
		f.LogResult(f.TxCollectionModify(keyFoo, contractID, "", "", "meta", modifiedMeta, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// Then the meta is modified
		require.Equal(t, modifiedMeta, f.QueryCollection(contractID).GetMeta())
	}
	t.Log("Modify token type")
	{
		// When modify token name
		modifiedName := firstName + "modified"
		f.LogResult(f.TxCollectionModify(keyFoo, contractID, tokenType, "", "name", modifiedName, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// Then the name is modified
		require.Equal(t, modifiedName, f.QueryTokenTypeCollection(contractID, tokenType).GetName())

		// When modify meta
		modifiedMeta := firstMeta + "modified"
		f.LogResult(f.TxCollectionModify(keyFoo, contractID, tokenType, "", "meta", modifiedMeta, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// Then the meta is modified
		require.Equal(t, modifiedMeta, f.QueryTokenTypeCollection(contractID, tokenType).GetMeta())
	}
	t.Log("Modify nft token")
	{
		// When modify token name
		modifiedName := firstName + "modified"
		f.LogResult(f.TxCollectionModify(keyFoo, contractID, tokenType, tokenIndex, "name", modifiedName, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// Then the name is modified
		require.Equal(t, modifiedName, f.QueryTokenCollection(contractID, tokenID).(collectionmodule.NFT).GetName())

		// When modify meta
		modifiedMeta := firstMeta + "modified"
		f.LogResult(f.TxCollectionModify(keyFoo, contractID, tokenType, tokenIndex, "meta", modifiedMeta, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// Then the meta is modified
		require.Equal(t, modifiedMeta, f.QueryTokenCollection(contractID, tokenID).(collectionmodule.NFT).GetMeta())
	}
	t.Log("Modify ft token")
	{
		// When modify token name
		modifiedName := firstName + "modified"
		f.LogResult(f.TxCollectionModify(keyFoo, contractID, tokenTypeFT, tokenIndexFT, "name", modifiedName, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// Then the name is modified
		require.Equal(t, modifiedName, f.QueryTokenCollection(contractID, tokenIDFT).(collectionmodule.FT).GetName())

		// When modify meta
		modifiedMeta := firstMeta + "modified"
		f.LogResult(f.TxCollectionModify(keyFoo, contractID, tokenTypeFT, tokenIndexFT, "meta", modifiedMeta, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// Then the meta is modified
		require.Equal(t, modifiedMeta, f.QueryTokenCollection(contractID, tokenIDFT).(collectionmodule.FT).GetMeta())
	}
}

func TestLinkCLIMempool(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	// Save key addresses for later use
	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	fooAcc := f.QueryAccount(fooAddr)
	startTokens := sdk.TokensFromConsensusPower(50)
	require.Equal(t, startTokens, fooAcc.GetCoins().AmountOf(denom))

	// Send some tokens from one account to the other
	sendTokens := sdk.TokensFromConsensusPower(10)
	success, stdout, stderr := f.TxSend(keyFoo, barAddr, sdk.NewCoin(denom, sendTokens), "-y", "-b sync")
	require.True(t, success)
	require.NotEmpty(t, stdout)
	require.Empty(t, stderr)

	// check mempool txs
	{
		result := f.MempoolUnconfirmedTxHashes()
		require.Equal(t, 1, result.Count)
		require.Equal(t, 1, result.Total)

		txHash := UnmarshalTxResponse(t, stdout).TxHash
		require.Equal(t, txHash, result.Txs[0])
	}

	// Ensure account balances match expected
	tests.WaitForNextNBlocksTM(2, f.Port)

	// check mempool empty
	{
		result := f.MempoolNumUnconfirmedTxs()
		require.Equal(t, 0, result.Count)
		require.Equal(t, 0, result.Total)
		require.Empty(t, result.Txs)
	}

	barAcc := f.QueryAccount(barAddr)
	require.Equal(t, sendTokens, barAcc.GetCoins().AmountOf(denom))
	fooAcc = f.QueryAccount(fooAddr)
	require.Equal(t, startTokens.Sub(sendTokens), fooAcc.GetCoins().AmountOf(denom))
}

func TestLinkCLITokenIssue(t *testing.T) {

	const (
		contractID1 = "9be17165"
		contractID2 = "678c146a"
		contractID3 = "3336b76f"
		description = "description"
		symbol      = "BTC"
		meta        = "{}"
		amount      = 10000
		decimals    = 6
	)

	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	// Issue Token.
	{
		f.TxTokenIssue(keyFoo, fooAddr, description, meta, symbol, amount, decimals, false, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		token := f.QueryToken(contractID1)
		require.Equal(t, description, token.GetName())
		require.Equal(t, contractID1, token.GetContractID())
		require.Equal(t, int64(decimals), token.GetDecimals().Int64())
		require.Equal(t, false, token.GetMintable())

		require.Equal(t, sdk.NewInt(amount), f.QueryBalanceToken(contractID1, fooAddr))
	}

	// Issue Token.
	{
		f.TxTokenIssue(keyFoo, fooAddr, description, meta, symbol, amount, decimals, true, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		token := f.QueryToken(contractID2)
		require.Equal(t, description, token.GetName())
		require.Equal(t, contractID2, token.GetContractID())
		require.Equal(t, int64(decimals), token.GetDecimals().Int64())
		require.Equal(t, true, token.GetMintable())

		require.Equal(t, sdk.NewInt(amount), f.QueryBalanceToken(contractID2, fooAddr))
	}

	// Query permissions for foo account
	{
		pms := f.QueryAccountPermission(f.KeyAddress(keyFoo), contractID1)
		require.Equal(t, 1, len(pms))
		require.Equal(t, "modify", pms[0].String())
	}
	{
		pms := f.QueryAccountPermission(f.KeyAddress(keyFoo), contractID2)
		require.Equal(t, 3, len(pms))
		require.Equal(t, "modify", pms[0].String())
		require.Equal(t, "mint", pms[1].String())
		require.Equal(t, "burn", pms[2].String())
	}

	// Query permissions for bar account
	{
		pms := f.QueryAccountPermission(f.KeyAddress(keyBar), contractID1)
		require.Equal(t, 0, len(pms))
		pms = f.QueryAccountPermission(f.KeyAddress(keyBar), contractID2)
		require.Equal(t, 0, len(pms))
	}

	// Transfer Token
	{
		f.TxTokenTransfer(keyFoo, barAddr, contractID1, amount/2, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		balance := f.QueryBalanceToken(contractID1, barAddr)
		require.Equal(t, int64(amount/2), balance.Int64())

		const keyMarshall = "marshall"
		f.KeysDelete(keyMarshall)
		f.KeysAdd(keyMarshall)
		marshallAddr := f.KeyAddress(keyMarshall)
		f.TxTokenTransfer(keyBar, marshallAddr, contractID1, amount/4, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		tinaAccount := f.QueryAccount(marshallAddr)
		require.Equal(t, marshallAddr.String(), tinaAccount.Address.String())
	}
}

func TestLinkCLITokenMintBurn(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	const (
		contractID = "9be17165"

		initAmount    = 2000000
		initAmountStr = "2000000"
		mintAmount    = 200
		mintAmountStr = "200"
		burnAmount    = 100
		burnAmountStr = "100"
		description   = "decription"
		symbol        = "BTC"
		meta          = "meta"
	)

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	// Create Account bar
	{
		sendTokens := sdk.TokensFromConsensusPower(1)
		f.TxSend(keyFoo, barAddr, sdk.NewCoin(denom, sendTokens), "-y")
	}
	// Issue a Token and check the amount
	{
		f.TxTokenIssue(keyFoo, fooAddr, description, meta, symbol, initAmount, 6, true, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		token := f.QueryToken(contractID)
		require.Equal(t, description, token.GetName())
		require.Equal(t, contractID, token.GetContractID())
		require.Equal(t, int64(6), token.GetDecimals().Int64())
		require.Equal(t, true, token.GetMintable())

		require.Equal(t, int64(initAmount), f.QuerySupplyToken(contractID).Int64())
		require.Equal(t, int64(initAmount), f.QueryMintToken(contractID).Int64())
		require.Equal(t, int64(0), f.QueryBurnToken(contractID).Int64())
		require.Equal(t, int64(initAmount), f.QueryBalanceToken(contractID, fooAddr).Int64())
	}
	// Mint/Burn by token owner
	{
		f.TxTokenMint(keyFoo, contractID, fooAddr.String(), mintAmountStr, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.Equal(t, int64(initAmount+mintAmount), f.QuerySupplyToken(contractID).Int64())
		require.Equal(t, int64(initAmount+mintAmount), f.QueryMintToken(contractID).Int64())
		require.Equal(t, int64(0), f.QueryBurnToken(contractID).Int64())
		require.Equal(t, int64(initAmount+mintAmount), f.QueryBalanceToken(contractID, fooAddr).Int64())

		f.TxTokenBurn(keyFoo, contractID, burnAmountStr, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.Equal(t, int64(initAmount+mintAmount-burnAmount), f.QuerySupplyToken(contractID).Int64())
		require.Equal(t, int64(initAmount+mintAmount), f.QueryMintToken(contractID).Int64())
		require.Equal(t, int64(burnAmount), f.QueryBurnToken(contractID).Int64())
		require.Equal(t, int64(initAmount+mintAmount-burnAmount), f.QueryBalanceToken(contractID, fooAddr).Int64())
	}

	// Mint/Burn Fail
	{
		// Foo try to burn but insufficient
		_, stdOut, _ := f.TxTokenBurn(keyFoo, contractID, initAmountStr+initAmountStr, "-y")
		require.Contains(t, stdOut, "not enough coins")
		// bar try to mint but has no permission
		_, stdOut, _ = f.TxTokenMint(keyBar, contractID, barAddr.String(), mintAmountStr, "-y")
		require.Contains(t, stdOut, "account does not have the permission")

		// Amount not changed
		require.Equal(t, int64(initAmount+mintAmount-burnAmount), f.QuerySupplyToken(contractID).Int64())
		require.Equal(t, int64(initAmount+mintAmount), f.QueryMintToken(contractID).Int64())
		require.Equal(t, int64(burnAmount), f.QueryBurnToken(contractID).Int64())
		require.Equal(t, int64(initAmount+mintAmount-burnAmount), f.QueryBalanceToken(contractID, fooAddr).Int64())
	}

	// Grant Permission to bar
	{
		f.TxTokenGrantPerm(keyFoo, barAddr.String(), contractID, "mint", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		f.TxTokenMint(keyBar, contractID, barAddr.String(), mintAmountStr, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.Equal(t, int64(initAmount+mintAmount-burnAmount+mintAmount), f.QuerySupplyToken(contractID).Int64())
		require.Equal(t, int64(initAmount+mintAmount+mintAmount), f.QueryMintToken(contractID).Int64())
		require.Equal(t, int64(burnAmount), f.QueryBurnToken(contractID).Int64())
		require.Equal(t, int64(mintAmount), f.QueryBalanceToken(contractID, barAddr).Int64())
	}

	// Revoke permission from foo
	{
		f.TxTokenRevokePerm(keyFoo, contractID, "mint", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		_, stdOut, _ := f.TxTokenMint(keyFoo, contractID, fooAddr.String(), mintAmountStr, "-y")
		require.Contains(t, stdOut, "account does not have the permission")

		// Amount not changed
		require.Equal(t, int64(initAmount+mintAmount-burnAmount+mintAmount), f.QuerySupplyToken(contractID).Int64())
		require.Equal(t, int64(initAmount+mintAmount+mintAmount), f.QueryMintToken(contractID).Int64())
		require.Equal(t, int64(burnAmount), f.QueryBurnToken(contractID).Int64())
		require.Equal(t, int64(initAmount+mintAmount-burnAmount), f.QueryBalanceToken(contractID, fooAddr).Int64())
		require.Equal(t, int64(mintAmount), f.QueryBalanceToken(contractID, barAddr).Int64())
	}

	// Burn from bar without permissions; burn failure
	{
		f.TxTokenBurn(keyBar, contractID, burnAmountStr, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.Equal(t, int64(initAmount+mintAmount-burnAmount+mintAmount), f.QuerySupplyToken(contractID).Int64())
		require.Equal(t, int64(initAmount+mintAmount+mintAmount), f.QueryMintToken(contractID).Int64())
		require.Equal(t, int64(burnAmount), f.QueryBurnToken(contractID).Int64())
		require.Equal(t, int64(mintAmount), f.QueryBalanceToken(contractID, barAddr).Int64())
	}
}

func TestLinkCLITokenCollection(t *testing.T) {

	const (
		contractID1 = "9be17165"
		tokenID01   = "0000000100000000"
		tokenID02   = "0000000200000000"
		tokenID03   = "0000000300000000"
		tokenID04   = "0000000400000000"
		description = "description"
		meta        = "meta"
		tokenuri    = "uri:itisbrown"
	)

	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)
	// barSuffix := types.AccAddrSuffix(barAddr)

	// Create Account bar
	{
		sendTokens := sdk.TokensFromConsensusPower(1)
		f.TxSend(keyFoo, barAddr, sdk.NewCoin(denom, sendTokens), "-y")
	}
	// Create Collection
	{
		f.TxTokenCreateCollection(keyFoo, description, meta, tokenuri, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
	}
	// Issue collective token brown with token id
	{
		f.TxTokenIssueFTCollection(keyFoo, contractID1, fooAddr, description, meta, 10000, 6, false, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		f.QueryTokenExpectEmpty(contractID1)
		token := f.QueryTokenCollection(contractID1, tokenID01)
		require.Equal(t, description, token.GetName())
		require.Equal(t, contractID1, token.GetContractID())
		require.Equal(t, tokenID01, token.GetTokenID())
		require.Equal(t, int64(6), token.(collectionmodule.FT).GetDecimals().Int64())
		require.Equal(t, false, token.(collectionmodule.FT).GetMintable())
		require.Equal(t, sdk.NewInt(10000), f.QueryBalanceCollection(contractID1, tokenID01, fooAddr))
		require.Equal(t, sdk.NewInt(10000), f.QueryTotalSupplyTokenCollection(contractID1, tokenID01))
		require.Equal(t, sdk.NewInt(10000), f.QueryTotalMintTokenCollection(contractID1, tokenID01))
		require.Equal(t, sdk.ZeroInt(), f.QueryTotalBurnTokenCollection(contractID1, tokenID01))
	}
	{
		f.TxTokenIssueFTCollection(keyFoo, contractID1, fooAddr, description, meta, 20000, 6, true, "-y")
		f.TxTokenIssueFTCollection(keyFoo, contractID1, fooAddr, description, meta, 30000, 6, true, "-y")

		token := f.QueryTokenCollection(contractID1, tokenID01)
		require.Equal(t, contractID1, token.GetContractID())
		require.Equal(t, tokenID01, token.GetTokenID())
		require.Equal(t, sdk.NewInt(10000), f.QueryTotalSupplyTokenCollection(contractID1, tokenID01))
		require.Equal(t, sdk.NewInt(10000), f.QueryTotalMintTokenCollection(contractID1, tokenID01))
		require.Equal(t, sdk.ZeroInt(), f.QueryTotalBurnTokenCollection(contractID1, tokenID01))

		token = f.QueryTokenCollection(contractID1, tokenID02)
		require.Equal(t, contractID1, token.GetContractID())
		require.Equal(t, tokenID02, token.GetTokenID())
		require.Equal(t, sdk.NewInt(20000), f.QueryTotalSupplyTokenCollection(contractID1, tokenID02))
		require.Equal(t, sdk.NewInt(20000), f.QueryTotalMintTokenCollection(contractID1, tokenID02))
		require.Equal(t, sdk.ZeroInt(), f.QueryTotalBurnTokenCollection(contractID1, tokenID02))

		token = f.QueryTokenCollection(contractID1, tokenID03)
		require.Equal(t, contractID1, token.GetContractID())
		require.Equal(t, tokenID03, token.GetTokenID())
		require.Equal(t, sdk.NewInt(30000), f.QueryTotalSupplyTokenCollection(contractID1, tokenID03))
		require.Equal(t, sdk.NewInt(30000), f.QueryTotalMintTokenCollection(contractID1, tokenID03))
		require.Equal(t, sdk.ZeroInt(), f.QueryTotalBurnTokenCollection(contractID1, tokenID03))

		tokensByOption := f.QueryTokensByTokenTypeCollection(contractID1, tokenID01[:8])
		require.Equal(t, 1, len(tokensByOption))
		require.Equal(t, contractID1, tokensByOption[0].GetContractID())
		require.Equal(t, tokenID01, tokensByOption[0].GetTokenID())
		require.Equal(t, tokenID01[:8], tokensByOption[0].GetTokenType())
		require.Equal(t, meta, tokensByOption[0].GetMeta())
		require.Equal(t, description, tokensByOption[0].GetName())

		tokensByNoOption := f.QueryTokensCollection(contractID1)
		require.Equal(t, 3, len(tokensByNoOption))
		require.Equal(t, contractID1, tokensByNoOption[0].GetContractID())
		require.Equal(t, tokenID01, tokensByNoOption[0].GetTokenID())
		require.Equal(t, tokenID01[:8], tokensByNoOption[0].GetTokenType())
		require.Equal(t, meta, tokensByNoOption[0].GetMeta())
		require.Equal(t, description, tokensByNoOption[0].GetName())

		require.Equal(t, contractID1, tokensByNoOption[1].GetContractID())
		require.Equal(t, tokenID02, tokensByNoOption[1].GetTokenID())
		require.Equal(t, tokenID02[:8], tokensByNoOption[1].GetTokenType())
		require.Equal(t, meta, tokensByNoOption[1].GetMeta())
		require.Equal(t, description, tokensByNoOption[1].GetName())

		require.Equal(t, contractID1, tokensByNoOption[2].GetContractID())
		require.Equal(t, tokenID03, tokensByNoOption[2].GetTokenID())
		require.Equal(t, tokenID03[:8], tokensByNoOption[2].GetTokenType())
		require.Equal(t, meta, tokensByNoOption[2].GetMeta())
		require.Equal(t, description, tokensByNoOption[2].GetName())

		toknenIDNoExist := "00000009"
		tokensEmpty := f.QueryTokensByTokenTypeCollection(contractID1, toknenIDNoExist)
		require.Empty(t, tokensEmpty)
		toknenIDNoExist2 := "0000000a"
		tokensEmpty2 := f.QueryTokensByTokenTypeCollection(contractID1, toknenIDNoExist2)
		require.Empty(t, tokensEmpty2)
	}

	// Bar cannot issue with the collection
	{
		f.TxTokenIssueFTCollection(keyBar, contractID1, barAddr, description, meta, 10000, 6, false, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		f.QueryTokenCollectionExpectEmpty(contractID1, tokenID04)
	}

	// Bar can issue collective token when granted the issue permission
	{
		f.TxCollectionGrantPerm(keyFoo, barAddr, contractID1, "issue", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		f.TxCollectionGrantPerm(keyFoo, barAddr, contractID1, "mint", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		f.TxCollectionGrantPerm(keyFoo, barAddr, contractID1, "burn", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		f.TxTokenIssueFTCollection(keyBar, contractID1, barAddr, description, meta, 40000, 6, true, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		token := f.QueryTokenCollection(contractID1, tokenID04)
		require.Equal(t, contractID1, token.GetContractID())
		require.Equal(t, tokenID04, token.GetTokenID())
		require.Equal(t, sdk.NewInt(40000), f.QueryTotalSupplyTokenCollection(contractID1, tokenID04))
		require.Equal(t, sdk.NewInt(40000), f.QueryTotalMintTokenCollection(contractID1, tokenID04))
		require.Equal(t, sdk.ZeroInt(), f.QueryTotalBurnTokenCollection(contractID1, tokenID04))
	}

	// Mint and Burn FTs in the collection
	{
		amount := fmt.Sprintf("%d:%s", 1000, tokenID04)
		f.TxTokenMintFTCollection(keyBar, contractID1, barAddr.String(), amount, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		require.Equal(t, sdk.NewInt(41000), f.QueryTotalSupplyTokenCollection(contractID1, tokenID04))
		require.Equal(t, sdk.NewInt(41000), f.QueryTotalMintTokenCollection(contractID1, tokenID04))
		require.Equal(t, sdk.ZeroInt(), f.QueryTotalBurnTokenCollection(contractID1, tokenID04))

		f.TxTokenBurnFTCollection(keyBar, contractID1, tokenID04, int64(2000), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		require.Equal(t, sdk.NewInt(39000), f.QueryTotalSupplyTokenCollection(contractID1, tokenID04))
		require.Equal(t, sdk.NewInt(41000), f.QueryTotalMintTokenCollection(contractID1, tokenID04))
		require.Equal(t, sdk.NewInt(2000), f.QueryTotalBurnTokenCollection(contractID1, tokenID04))
	}

	// Multi-transfer and multi-mint FTs
	{
		amount := fmt.Sprintf("%d:%s,%d:%s", 10000, tokenID02, 20000, tokenID03)
		f.TxTokenTransferFTCollection(keyFoo, contractID1, barAddr.String(), amount, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		require.Equal(t, sdk.NewInt(10000), f.QueryBalanceCollection(contractID1, tokenID02, fooAddr))
		require.Equal(t, sdk.NewInt(10000), f.QueryBalanceCollection(contractID1, tokenID02, barAddr))
		require.Equal(t, sdk.NewInt(10000), f.QueryBalanceCollection(contractID1, tokenID03, fooAddr))
		require.Equal(t, sdk.NewInt(20000), f.QueryBalanceCollection(contractID1, tokenID03, barAddr))

		f.TxTokenMintFTCollection(keyFoo, contractID1, fooAddr.String(), amount, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		require.Equal(t, sdk.NewInt(20000), f.QueryBalanceCollection(contractID1, tokenID02, fooAddr))
		require.Equal(t, sdk.NewInt(30000), f.QueryTotalSupplyTokenCollection(contractID1, tokenID02))
		require.Equal(t, sdk.NewInt(30000), f.QueryTotalMintTokenCollection(contractID1, tokenID02))
		require.Equal(t, sdk.ZeroInt(), f.QueryTotalBurnTokenCollection(contractID1, tokenID02))

		require.Equal(t, sdk.NewInt(30000), f.QueryBalanceCollection(contractID1, tokenID03, fooAddr))
		require.Equal(t, sdk.NewInt(50000), f.QueryTotalSupplyTokenCollection(contractID1, tokenID03))
		require.Equal(t, sdk.NewInt(50000), f.QueryTotalMintTokenCollection(contractID1, tokenID03))
		require.Equal(t, sdk.ZeroInt(), f.QueryTotalBurnTokenCollection(contractID1, tokenID03))
	}
}

func TestLinkCliTokenApprove(t *testing.T) {
	const (
		contractID  = "9be17165"
		description = "description"
		meta        = "meta"
		tokenuri    = "uri:itisbrown"
	)

	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()
	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	fooAddr := f.KeyAddress(keyFoo)
	tinaAddr := f.KeyAddress(UserTina)
	kelvinAddr := f.KeyAddress(UserKevin)

	// Create Collection
	{
		f.LogResult(f.TxTokenCreateCollection(keyFoo, description, meta, tokenuri, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
		f.LogResult(f.TxTokenCreateCollection(UserTina, description, meta, tokenuri, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}
	// tx Collection Approve
	{
		f.LogResult(f.TxCollectionApprove(keyFoo, contractID, kelvinAddr, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
		f.LogResult(f.TxCollectionApprove(UserTina, contractID, kelvinAddr, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// query Collection Approved
	{
		isKelApprovedByFoo := f.QueryApprovedTokenCollection(contractID, kelvinAddr, fooAddr)
		require.True(t, isKelApprovedByFoo)
		isKelApprovedByTina := f.QueryApprovedTokenCollection(contractID, kelvinAddr, tinaAddr)
		require.True(t, isKelApprovedByTina)
		isTinaApprovedByFoo := f.QueryApprovedTokenCollection(contractID, tinaAddr, fooAddr)
		require.False(t, isTinaApprovedByFoo)
	}

	// query CollectionApprovers
	{
		approverAddresses := f.QueryApproversTokenCollection(contractID, kelvinAddr)
		if bytes.Compare(tinaAddr, fooAddr) < 0 {
			require.Equal(t, tinaAddr.String(), approverAddresses[0].String())
			require.Equal(t, fooAddr.String(), approverAddresses[1].String())
		} else {
			require.Equal(t, tinaAddr.String(), approverAddresses[1].String())
			require.Equal(t, fooAddr.String(), approverAddresses[0].String())
		}
		addsEmpty := f.QueryApproversTokenCollection(contractID, fooAddr)
		require.Empty(t, addsEmpty)
	}
}

func TestLinkCLITokenNFT(t *testing.T) {

	const (
		contractID  = "9be17165"
		tokenType   = "10000001"
		tokenType2  = "10000002"
		tokenType3  = "10000003"
		tokenID01   = "1000000100000001"
		tokenID02   = "1000000100000002"
		tokenID03   = "1000000100000003"
		tokenID04   = "1000000100000004"
		description = "description"
		meta        = "meta"
		tokenuri    = "uri:itisbrown"
	)

	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	// Create Account bar
	{
		sendTokens := sdk.TokensFromConsensusPower(1)
		f.TxSend(keyFoo, barAddr, sdk.NewCoin(denom, sendTokens), "-y")
	}
	// Create Collection
	{
		f.TxTokenCreateCollection(keyFoo, description, meta, tokenuri, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
	}
	// Issue Collective NFT for the collection
	{
		f.LogResult(f.TxTokenIssueNFTCollection(keyFoo, contractID, description, meta, "-y"))
		f.LogResult(f.TxTokenIssueNFTCollection(keyFoo, contractID, description, meta, "-y"))
		f.LogResult(f.TxTokenIssueNFTCollection(keyFoo, contractID, description, meta, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
		mintParam := strings.Join([]string{tokenType, description, meta}, ":")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam, "-y")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam, "-y")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam, "-y")
		mintParam2 := strings.Join([]string{tokenType2, description, meta}, ":")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam2, "-y")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam2, "-y")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam2, "-y")
		mintParam3 := strings.Join([]string{tokenType3, description, meta}, ":")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam3, "-y")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam3, "-y")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam3, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		token := f.QueryTokenCollection(contractID, tokenID01)
		require.Equal(t, contractID, token.GetContractID())
		require.Equal(t, tokenID01, token.GetTokenID())
		token = f.QueryTokenCollection(contractID, tokenID02)
		require.Equal(t, contractID, token.GetContractID())
		require.Equal(t, tokenID02, token.GetTokenID())
		token = f.QueryTokenCollection(contractID, tokenID03)
		require.Equal(t, contractID, token.GetContractID())
		require.Equal(t, tokenID03, token.GetTokenID())

		tokensByOption := f.QueryTokensByTokenTypeCollection(contractID, tokenType)
		require.Equal(t, 3, len(tokensByOption))
		require.Equal(t, contractID, tokensByOption[0].GetContractID())
		require.Equal(t, tokenID01, tokensByOption[0].GetTokenID())
		require.Equal(t, tokenID01[:8], tokensByOption[0].GetTokenType())
		require.Equal(t, meta, tokensByOption[0].GetMeta())
		require.Equal(t, description, tokensByOption[0].GetName())
		require.Equal(t, contractID, tokensByOption[1].GetContractID())
		require.Equal(t, tokenID02, tokensByOption[1].GetTokenID())
		require.Equal(t, tokenID02[:8], tokensByOption[1].GetTokenType())
		require.Equal(t, meta, tokensByOption[1].GetMeta())
		require.Equal(t, description, tokensByOption[1].GetName())

		tokensByNoOption := f.QueryTokensCollection(contractID)
		require.Equal(t, 9, len(tokensByNoOption))
		require.Equal(t, contractID, tokensByNoOption[0].GetContractID())
		require.Equal(t, tokenID01, tokensByNoOption[0].GetTokenID())
		require.Equal(t, tokenID01[:8], tokensByNoOption[0].GetTokenType())
		require.Equal(t, meta, tokensByNoOption[0].GetMeta())
		require.Equal(t, description, tokensByNoOption[0].GetName())
		require.Equal(t, contractID, tokensByNoOption[1].GetContractID())
		require.Equal(t, tokenID02, tokensByNoOption[1].GetTokenID())
		require.Equal(t, tokenID02[:8], tokensByNoOption[1].GetTokenType())
		require.Equal(t, meta, tokensByNoOption[1].GetMeta())
		require.Equal(t, description, tokensByNoOption[1].GetName())

		tokenTypeNoExist := "10000009"
		tokensEmpty := f.QueryTokensByTokenTypeCollection(contractID, tokenTypeNoExist)
		require.Empty(t, tokensEmpty)
	}

	// Multi-transfer NFTs
	{
		f.TxTokenTransferNFTCollection(keyFoo, contractID, barAddr.String(), tokenID01, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		require.Equal(t, sdk.ZeroInt(), f.QueryBalanceCollection(contractID, tokenID01, fooAddr))
		require.Equal(t, sdk.NewInt(1), f.QueryBalanceCollection(contractID, tokenID01, barAddr))

		tokenIDs := tokenID02 + "," + tokenID03
		f.TxTokenTransferNFTCollection(keyFoo, contractID, barAddr.String(), tokenIDs, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		require.Equal(t, sdk.ZeroInt(), f.QueryBalanceCollection(contractID, tokenID02, fooAddr))
		require.Equal(t, sdk.ZeroInt(), f.QueryBalanceCollection(contractID, tokenID03, fooAddr))
		require.Equal(t, sdk.NewInt(1), f.QueryBalanceCollection(contractID, tokenID02, barAddr))
		require.Equal(t, sdk.NewInt(1), f.QueryBalanceCollection(contractID, tokenID03, barAddr))
	}

	// Multi-mint NFTs
	{
		myTokenType := "10000002"
		myTokenID01 := "1000000200000001"

		f.LogResult(f.TxTokenIssueNFTCollection(keyFoo, contractID, description, meta, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		mint1 := strings.Join([]string{tokenType, description, meta}, ":")
		mint2 := strings.Join([]string{myTokenType, description, meta}, ":")
		mintParam := strings.Join([]string{mint1, mint2}, ",")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		token := f.QueryTokenCollection(contractID, tokenID04)
		require.Equal(t, contractID, token.GetContractID())
		require.Equal(t, tokenID04, token.GetTokenID())
		token = f.QueryTokenCollection(contractID, myTokenID01)
		require.Equal(t, contractID, token.GetContractID())
		require.Equal(t, myTokenID01, token.GetTokenID())

		// check minting order
		_, stdout, _ := f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam, "-y")
		require.NotEmpty(t, stdout)
		require.Less(t, strings.Index(stdout, tokenType), strings.Index(stdout, myTokenType))

		mintParam2 := strings.Join([]string{mint2, mint1}, ",")
		_, stdout2, _ := f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam2, "-y")
		require.NotEmpty(t, stdout2)
		require.Greater(t, strings.Index(stdout2, tokenType), strings.Index(stdout2, myTokenType))
	}
}

func TestLinkCLISendGenerateSignAndBroadcastWithToken(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	contractID := "9be17165"

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	fooAddr := f.KeyAddress(keyFoo)

	success, stdout, stderr := f.TxTokenIssue(fooAddr.String(), fooAddr, "test", "{}", "BTC", 10000, 6, true, "--generate-only")
	require.True(t, success)
	require.Empty(t, stderr)
	msg := UnmarshalStdTx(t, stdout)
	require.True(t, msg.Fee.Gas > 0)
	require.Equal(t, len(msg.Msgs), 1)

	// Write the output to disk
	unsignedTxFile := WriteToNewTempFile(t, stdout)
	defer os.Remove(unsignedTxFile.Name())

	// Test sign --validate-signatures
	success, stdout, _ = f.TxSign(keyFoo, unsignedTxFile.Name(), "--validate-signatures")
	require.False(t, success)
	require.Equal(t, fmt.Sprintf("Signers:\n  0: %v\n\nSignatures:\n\n", fooAddr.String()), stdout)

	// Test sign
	success, stdout, _ = f.TxSign(keyFoo, unsignedTxFile.Name())
	require.True(t, success)
	msg = UnmarshalStdTx(t, stdout)
	require.Equal(t, len(msg.Msgs), 1)
	require.Equal(t, 1, len(msg.GetSignatures()))
	require.Equal(t, fooAddr.String(), msg.GetSigners()[0].String())

	// Write the output to disk
	signedTxFile := WriteToNewTempFile(t, stdout)
	defer os.Remove(signedTxFile.Name())

	// Test sign --validate-signatures
	success, stdout, _ = f.TxSign(keyFoo, signedTxFile.Name(), "--validate-signatures")
	require.True(t, success)
	require.Equal(t, fmt.Sprintf("Signers:\n  0: %v\n\nSignatures:\n  0: %v\t\t\t[OK]\n\n", fooAddr.String(),
		fooAddr.String()), stdout)

	f.QueryTokenExpectEmpty(contractID)

	// Test broadcast
	success, stdout, _ = f.TxBroadcast(signedTxFile.Name())
	require.True(t, success)
	tests.WaitForNextNBlocksTM(1, f.Port)

	token := f.QueryToken(contractID)
	require.Equal(t, "test", token.GetName())
	require.Equal(t, int64(6), token.GetDecimals().Int64())
}

func TestLinkCLIEmpty(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	brianAddr := f.KeyAddress(UserBrian).String()

	t.Logf("[Account] Do nothing with empty message")
	success, stdout, stderr := f.TxEmpty(brianAddr, "-y")
	{
		require.True(t, success)
		require.NotEmpty(t, stdout)
		require.Empty(t, stderr)
	}
}

/*
from cosmos-sdk v0.38.0, multiple txs from an account for a block is not allowed
Modify the IncrementSequenceDecorator to avoid the restriction
*/
func TestLinkCLIIncrementSequenceDecorator(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	sendTokens := sdk.TokensFromConsensusPower(1)

	fooAcc := f.QueryAccount(fooAddr)

	// Prepare signed Tx
	var signedTxFiles []*os.File
	for idx := 0; idx < 3; idx++ {
		// Test generate sendTx, estimate gas
		success, stdout, _ := f.TxSend(fooAddr.String(), barAddr, sdk.NewCoin(denom, sendTokens), "--generate-only")
		require.True(t, success)

		// Write the output to disk
		unsignedTxFile := WriteToNewTempFile(t, stdout)
		defer os.Remove(unsignedTxFile.Name())

		// Test sign
		success, stdout, _ = f.TxSign(keyFoo, unsignedTxFile.Name(), "--offline", "--account-number", strconv.Itoa(int(fooAcc.AccountNumber)), "--sequence", strconv.Itoa(int(fooAcc.Sequence)+idx))
		require.True(t, success)

		// Write the output to disk
		signedTxFile := WriteToNewTempFile(t, stdout)
		signedTxFiles = append(signedTxFiles, signedTxFile)
		defer os.Remove(signedTxFile.Name())
	}
	// Wait for a new block
	tests.WaitForNextNBlocksTM(1, f.Port)

	var txHashes []string
	// Broadcast the signed Txs
	for _, signedTxFile := range signedTxFiles {
		// Test broadcast
		success, stdout, _ := f.TxBroadcast(signedTxFile.Name(), "--broadcast-mode", "sync")
		require.True(t, success)
		sendResp := UnmarshalTxResponse(t, stdout)
		txHashes = append(txHashes, sendResp.TxHash)
	}

	// Wait for a new block
	tests.WaitForNextNBlocksTM(2, f.Port)

	// All Txs are in one block
	height := f.QueryTx(txHashes[0]).Height
	for _, txHash := range txHashes {
		require.Equal(t, height, f.QueryTx(txHash).Height)
	}
}

func TestLinkCliTokenProxy(t *testing.T) {
	const (
		contractID = "9be17165"
		firstName  = "itisbrown"
		name       = "description"
		meta       = "{}"
		symbol     = "BTC"
		amount     = 10000
		decimals   = 6
		sendAmount = 10
		burnAmount = 5
	)

	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)
	kevinAddr := f.KeyAddress(UserKevin)

	f.TxTokenIssue(keyFoo, fooAddr, name, meta, symbol, amount, decimals, true, "-y")
	tests.WaitForNextNBlocksTM(1, f.Port)
	firstResult := f.QueryToken(contractID)
	require.Equal(t, name, firstResult.GetName())

	// query token approved and get false
	{
		isApprovedByFoo := f.QueryApprovedToken(contractID, kevinAddr, fooAddr)
		require.False(t, isApprovedByFoo)
	}

	// tx token Approve
	{
		f.LogResult(f.TxTokenApprove(keyFoo, contractID, kevinAddr, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// query token approved and get true
	{
		isApprovedByFoo := f.QueryApprovedToken(contractID, kevinAddr, fooAddr)
		require.True(t, isApprovedByFoo)
	}

	// tx token transfer-from
	{
		f.LogResult(f.TxTokenTransferFrom(UserKevin, contractID, fooAddr, barAddr, sendAmount, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// check amount if transfer succeeds
	{
		amountOfBar := f.QueryBalanceToken(contractID, barAddr, "-y").Int64()
		require.Equal(t, int(amountOfBar), sendAmount)
	}

	// Grant permission for burn
	{
		f.TxTokenGrantPerm(keyFoo, kevinAddr.String(), contractID, "burn", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// tx token burn-from
	{
		f.LogResult(f.TxTokenBurnFrom(UserKevin, contractID, fooAddr, burnAmount, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// check amount if transfer succeeds
	{
		amountOfBar := f.QueryBalanceToken(contractID, fooAddr, "-y").Int64()
		require.Equal(t, int(amountOfBar), 9985) // 10000 - 10(transfer) - 5(burn) = 9985
	}
}
