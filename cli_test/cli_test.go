// +build cli_test

package clitest

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/line/link/x/bank"
	collectionmodule "github.com/line/link/x/collection"
	"github.com/line/link/x/proxy"
	sbox "github.com/line/link/x/safetybox"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestModifyToken(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()
	defer f.Cleanup()

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
	f.Cleanup()
}

func TestModifyCollection(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()
	defer f.Cleanup()

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
	f.Cleanup()
}

func TestLinkCLIProxy(t *testing.T) {
	t.Skip("SKIP: Proxy module is not in use")

	t.Parallel()
	f := InitFixtures(t)

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()
	defer f.Cleanup()

	denom := DenomLink
	tinaTheProxy := f.KeyAddress(UserTina).String()
	rinahTheOnBehalfOf := f.KeyAddress(UserRinah).String()
	evelynTheReceiver := f.KeyAddress(UserEvelyn).String()

	// rinahTheOnBehalfOf's initial balance
	t.Logf("[Proxy] Initial balance check for the OnBeHalfOf")
	initialBalance := f.QueryAccount(f.KeyAddress(UserRinah)).Coins
	{
		require.Equal(t, initialBalance, defaultCoins)
	}

	// `tinaTheProxy` tries to send 5 link to `evelynTheReceiver` on behalf of `rinahTheOnBehalfOf`
	fiveCoins := sdk.NewInt(5)
	t.Logf("[Proxy] The proxy tries to send %d link to the receiver on behalf of the coin owner - should fail", fiveCoins)
	{
		result, stdout, stderr := f.TxProxySendCoinsFrom(tinaTheProxy, rinahTheOnBehalfOf, evelynTheReceiver, denom, fiveCoins, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		// should fail as it's not approved with ErrProxyNotExist
		require.True(t, result)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(proxy.ErrProxyNotExist, "Proxy: %s, Account: %s", tinaTheProxy, rinahTheOnBehalfOf).Error(), "\""),
			strings.Split(stdout, "\\\\\\\"")[9],
		)
		require.Equal(t, "", stderr)
	}

	// `rinahTheOnBehalfOf` approves 5 link for `tinaTheProxy`
	approved := sdk.NewInt(5)
	t.Logf("[Proxy] The coin owner approves %d link for the proxy", approved)
	{
		result, _, stderr := f.TxProxyApproveCoins(tinaTheProxy, rinahTheOnBehalfOf, denom, approved, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		// should succeed
		require.True(t, result)
		require.Equal(t, "", stderr)
	}

	// check the allowance
	t.Logf("[Proxy] Check the allowance - should be %d", approved)
	{
		allowance := f.QueryProxyAllowance(tinaTheProxy, rinahTheOnBehalfOf, denom)
		require.Equal(t, tinaTheProxy, allowance.Proxy.String())
		require.Equal(t, rinahTheOnBehalfOf, allowance.OnBehalfOf.String())
		require.Equal(t, denom, allowance.Denom)
		require.Equal(t, approved, allowance.Amount)
	}

	// 'tinaTheProxy' tries to send 6 link to `evelynTheReceiver` on behalf of `rinahTheOnBehalfOf`
	sixCoins := sdk.NewInt(6)
	t.Logf("[Proxy] The proxy tries to send %d link to the receiver on behalf of the coin owner - should fail", sixCoins)
	{
		result, stdout, stderr := f.TxProxySendCoinsFrom(tinaTheProxy, rinahTheOnBehalfOf, evelynTheReceiver, denom, sixCoins, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		// should fail as it's more than approved
		require.True(t, result)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(proxy.ErrProxyNotEnoughApprovedCoins, "Approved: %v, Requested: %v", approved, sixCoins).Error(), "\""),
			strings.Split(stdout, "\\\\\\\"")[9],
		)
		require.Equal(t, "", stderr)
	}

	// `tinaTheProxy` sends 2 link to `evelynTheReceiver` on behalf of `rinahTheOnBehalfOf`
	sentAmount1 := sdk.NewInt(2)
	t.Logf("[Proxy] The proxy sends %d link to the receiver on behalf of the coin owner", sentAmount1)
	{
		result, _, stderr := f.TxProxySendCoinsFrom(tinaTheProxy, rinahTheOnBehalfOf, evelynTheReceiver, denom, sentAmount1, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		// should succeed
		require.True(t, result)
		require.Equal(t, "", stderr)
	}

	// check the allowance
	t.Logf("[Proxy] Check the allowance - should be %d", approved.Sub(sentAmount1))
	{
		allowance := f.QueryProxyAllowance(tinaTheProxy, rinahTheOnBehalfOf, denom)
		require.Equal(t, tinaTheProxy, allowance.Proxy.String())
		require.Equal(t, rinahTheOnBehalfOf, allowance.OnBehalfOf.String())
		require.Equal(t, denom, allowance.Denom)
		require.Equal(t, approved.Sub(sentAmount1), allowance.Amount)
	}

	// check balance of `rinahTheOnBehalfOf` and `evelynTheReceiver`
	t.Logf("[Proxy] Check balances to confirm")
	{
		diff := sdk.Coins{sdk.Coin{DenomLink, sentAmount1}}
		rinahBalance := f.QueryAccount(f.KeyAddress(UserRinah)).Coins
		require.Equal(t, rinahBalance, defaultCoins.Sub(diff))
		evelynBalance := f.QueryAccount(f.KeyAddress(UserEvelyn)).Coins
		require.Equal(t, evelynBalance, defaultCoins.Add(diff...))
	}

	// `rinahTheOnBehalfOf` tries to disapprove 4 link from `tinaTheProxy`
	fourCoins := sdk.NewInt(4)
	t.Logf("[Proxy] The coin owner disapproves %d link from the proxy - should fail", fourCoins)
	{
		result, stdout, stderr := f.TxProxyDisapproveCoins(tinaTheProxy, rinahTheOnBehalfOf, denom, fourCoins, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		// should fail as only 3 approved coins are left
		require.True(t, result)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(proxy.ErrProxyNotEnoughApprovedCoins, "Approved: %v, Requested: %v", approved.Sub(sentAmount1), fourCoins).Error(), "\""),
			strings.Split(stdout, "\\\\\\\"")[9],
		)
		require.Equal(t, "", stderr)
	}

	// `rinahTheOnBehalfOf` disapprove 1 link from `tinaTheProxy`
	disapproved := sdk.OneInt()
	t.Logf("[Proxy] The coin owner disapproves %d link from the proxy", disapproved)
	{
		result, _, stderr := f.TxProxyDisapproveCoins(tinaTheProxy, rinahTheOnBehalfOf, denom, disapproved, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		// should succeed
		require.True(t, result)
		require.Equal(t, "", stderr)
	}

	// check the allowance
	t.Logf("[Proxy] Check the allowance - should be %d", approved.Sub(sentAmount1).Sub(disapproved))
	{
		allowance := f.QueryProxyAllowance(tinaTheProxy, rinahTheOnBehalfOf, denom)
		require.Equal(t, tinaTheProxy, allowance.Proxy.String())
		require.Equal(t, rinahTheOnBehalfOf, allowance.OnBehalfOf.String())
		require.Equal(t, denom, allowance.Denom)
		require.Equal(t, approved.Sub(sentAmount1).Sub(disapproved), allowance.Amount)
	}

	// `tinaTheProxy` sends 2 link to `evelynTheReceiver` on behalf of `rinahTheOnBehalfOf`
	sentAmount2 := sdk.NewInt(2)
	t.Logf("[Proxy] The proxy sends %d link to the receiver on behalf of the coin owner", sentAmount2)
	{
		result, _, stderr := f.TxProxySendCoinsFrom(tinaTheProxy, rinahTheOnBehalfOf, evelynTheReceiver, denom, sentAmount2, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		// should succeed
		require.True(t, result)
		require.Equal(t, "", stderr)
	}

	// check balance of `rinahTheOnBehalfOf` and `evelynTheReceiver`
	t.Logf("[Proxy] Check balances to confirm")
	{
		diff := sdk.Coins{sdk.Coin{DenomLink, sentAmount1.Add(sentAmount2)}}
		rinahBalance := f.QueryAccount(f.KeyAddress(UserRinah)).Coins
		require.Equal(t, rinahBalance, defaultCoins.Sub(diff))
		evelynBalance := f.QueryAccount(f.KeyAddress(UserEvelyn)).Coins
		require.Equal(t, evelynBalance, defaultCoins.Add(diff...))
	}

	// 'tinaTheProxy' tries to send 1 link to `evelynTheReceiver` on behalf of `rinahTheOnBehalfOf`
	oneCoin := sdk.NewInt(1)
	t.Logf("[Proxy] The proxy tries to send %d link to the receiver on behalf of the coin owner - should fail", oneCoin)
	{
		result, stdout, stderr := f.TxProxySendCoinsFrom(tinaTheProxy, rinahTheOnBehalfOf, evelynTheReceiver, denom, oneCoin, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		// should fail as there is no coin approved (all sent!)
		require.True(t, result)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(proxy.ErrProxyNotExist, "Proxy: %s, Account: %s", tinaTheProxy, rinahTheOnBehalfOf).Error(), "\""),
			strings.Split(stdout, "\\\\\\\"")[9],
		)
		require.Equal(t, "", stderr)
	}

	// 'onBehalfOf' tries to disapprove 1 link from `proxy`
	t.Logf("[Proxy] The coin owner tries to disapprove %d link from the proxy - should fail", oneCoin)
	{
		result, stdout, stderr := f.TxProxyDisapproveCoins(tinaTheProxy, rinahTheOnBehalfOf, denom, oneCoin, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		// should fail as there is no coin approved (all sent!)
		require.True(t, result)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(proxy.ErrProxyNotExist, "Proxy: %s, Account: %s", tinaTheProxy, rinahTheOnBehalfOf).Error(), "\""),
			strings.Split(stdout, "\\\\\\\"")[9],
		)
		require.Equal(t, "", stderr)
	}
}

func TestLinkCLISafetyBox(t *testing.T) {
	t.Skip("SKIP: SafetyBox module is not in use")

	t.Parallel()
	f := InitFixtures(t)

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	var result bool

	id := "some_safety_box"
	denom := DenomLink
	rinahTheOwnerAddress := f.KeyAddress(UserRinah).String()

	// create a safety box w/ user rinah as the owner
	{
		result, _, _ = f.TxSafetyBoxCreate(id, rinahTheOwnerAddress, denom, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)
	}

	// rinah is the owner
	{
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleOwner, rinahTheOwnerAddress)
		require.True(t, sbr.HasRole)
	}

	// rinah is not in any other roles
	{
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleOperator, rinahTheOwnerAddress)
		require.False(t, sbr.HasRole)

		sbr = f.QuerySafetyBoxRole(id, sbox.RoleAllocator, rinahTheOwnerAddress)
		require.False(t, sbr.HasRole)

		sbr = f.QuerySafetyBoxRole(id, sbox.RoleIssuer, rinahTheOwnerAddress)
		require.False(t, sbr.HasRole)

		sbr = f.QuerySafetyBoxRole(id, sbox.RoleReturner, rinahTheOwnerAddress)
		require.False(t, sbr.HasRole)
	}

	// query the safety box
	{
		sb := f.QuerySafetyBox(id)
		require.Equal(t, sb.ID, id)
		require.Equal(t, sb.Owner.String(), rinahTheOwnerAddress)
	}

	// sending coins to the safety box directly should fail
	{
		sb := f.QuerySafetyBox(id)
		result, stdoutSendToBlacklist, _ := f.TxSend(keyFoo, sb.Address, sdk.NewCoin(denom, sdk.OneInt()), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)

		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(bank.ErrCanNotTransferToBlacklisted, "Addr: %s", sb.Address.String()).Error(), "\""),
			strings.Split(stdoutSendToBlacklist, "\\\\\\\"")[9],
		)
	}

	// create a safety box w/ multiple denoms should fail
	{
		tooManyDenoms := []string{DenomLink, DenomStake}
		result, stdoutBoxCreate, _ := f.TxSafetyBoxCreate("new_id", rinahTheOwnerAddress, DenomLink+","+DenomStake, "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)

		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxTooManyCoinDenoms, "Requested: %v", tooManyDenoms).Error(), "\""),
			strings.Split(stdoutBoxCreate, "\\\\\\\"")[9],
		)
	}

	// any coin transfer to the safety box from the owner should fail
	{
		resultAllocation, stdoutAllocation, _ := f.TxSafetyBoxSendCoins(id, sbox.ActionAllocate, DenomLink, int64(1), rinahTheOwnerAddress, "", "-y")
		resultRecall, stdoutRecall, _ := f.TxSafetyBoxSendCoins(id, sbox.ActionRecall, DenomLink, int64(1), rinahTheOwnerAddress, "", "-y")
		resultIssue, stdoutIssue, _ := f.TxSafetyBoxSendCoins(id, sbox.ActionIssue, DenomLink, int64(1), rinahTheOwnerAddress, rinahTheOwnerAddress, "-y")
		resultReturn, stdoutReturn, _ := f.TxSafetyBoxSendCoins(id, sbox.ActionReturn, DenomLink, int64(1), rinahTheOwnerAddress, "", "-y")

		// test all four txs in a single block to reduce the testing time
		// check the error message to get expected errors
		tests.WaitForNextNBlocksTM(1, f.Port)

		require.True(t, resultAllocation)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxPermissionAllocate, "Account: %s", rinahTheOwnerAddress).Error(), "\""),
			strings.Split(stdoutAllocation, "\\\\\\\"")[9],
		)

		require.True(t, resultRecall)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxPermissionRecall, "Account: %s", rinahTheOwnerAddress).Error(), "\""),
			strings.Split(stdoutRecall, "\\\\\\\"")[9],
		)

		require.True(t, resultIssue)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxPermissionIssue, "Account: %s", rinahTheOwnerAddress).Error(), "\""),
			strings.Split(stdoutIssue, "\\\\\\\"")[9],
		)

		require.True(t, resultReturn)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxPermissionReturn, "Account: %s", rinahTheOwnerAddress).Error(), "\""),
			strings.Split(stdoutReturn, "\\\\\\\"")[9],
		)
	}

	// the owner registers an operator
	{
		// register user tina as an operator
		result, _, _ = f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleOperator, rinahTheOwnerAddress, f.KeyAddress(UserTina).String(), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)

		// tina is now an operator
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleOperator, f.KeyAddress(UserTina).String())
		require.True(t, sbr.HasRole)
	}

	// the owner can't register allocator, issuer and returner
	{
		// registering as allocator, issuer and returner should fail
		resultAllocator, _, _ := f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleAllocator, rinahTheOwnerAddress, f.KeyAddress(UserKevin).String(), "-y")
		resultIssuer, _, _ := f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleIssuer, rinahTheOwnerAddress, f.KeyAddress(UserKevin).String(), "-y")
		resultReturner, _, _ := f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleReturner, rinahTheOwnerAddress, f.KeyAddress(UserKevin).String(), "-y")

		// test all four txs in a single block to reduce the testing time
		tests.WaitForNextNBlocksTM(1, f.Port)

		require.True(t, resultAllocator)
		require.True(t, resultIssuer)
		require.True(t, resultReturner)

		// kevin should not have the role
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleAllocator, f.KeyAddress(UserKevin).String())
		require.False(t, sbr.HasRole)

		sbr = f.QuerySafetyBoxRole(id, sbox.RoleIssuer, f.KeyAddress(UserKevin).String())
		require.False(t, sbr.HasRole)

		sbr = f.QuerySafetyBoxRole(id, sbox.RoleReturner, f.KeyAddress(UserKevin).String())
		require.False(t, sbr.HasRole)
	}

	// any coin transfer to the safety box from the operator should fail
	{
		resultAllocate, stdoutAllocate, _ := f.TxSafetyBoxSendCoins(id, sbox.ActionAllocate, DenomLink, int64(1), f.KeyAddress(UserKevin).String(), "", "-y")
		resultRecall, stdoutRecall, _ := f.TxSafetyBoxSendCoins(id, sbox.ActionRecall, DenomLink, int64(1), f.KeyAddress(UserKevin).String(), "", "-y")
		resultIssue, stdoutIssue, _ := f.TxSafetyBoxSendCoins(id, sbox.ActionIssue, DenomLink, int64(1), f.KeyAddress(UserKevin).String(), f.KeyAddress(UserKevin).String(), "-y")
		resultReturn, stdoutReturn, _ := f.TxSafetyBoxSendCoins(id, sbox.ActionReturn, DenomLink, int64(1), f.KeyAddress(UserKevin).String(), "", "-y")

		// test all four txs in a single block to reduce the testing time
		// check the error message to get expected errors
		tests.WaitForNextNBlocksTM(1, f.Port)

		require.True(t, resultAllocate)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxPermissionAllocate, "Account: %s", f.KeyAddress(UserKevin).String()).Error(), "\""),
			strings.Split(stdoutAllocate, "\\\\\\\"")[9],
		)

		require.True(t, resultRecall)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxPermissionRecall, "Account: %s", f.KeyAddress(UserKevin).String()).Error(), "\""),
			strings.Split(stdoutRecall, "\\\\\\\"")[9],
		)

		require.True(t, resultIssue)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxPermissionIssue, "Account: %s", f.KeyAddress(UserKevin).String()).Error(), "\""),
			strings.Split(stdoutIssue, "\\\\\\\"")[9],
		)

		require.True(t, resultReturn)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxPermissionReturn, "Account: %s", f.KeyAddress(UserKevin).String()).Error(), "\""),
			strings.Split(stdoutReturn, "\\\\\\\"")[9],
		)
	}

	// an operator registers an allocator
	{
		// tina, the operator registers kevin as an allocator
		result, _, _ = f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleAllocator, f.KeyAddress(UserTina).String(), f.KeyAddress(UserKevin).String(), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)

		// kevin is now an operator
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleAllocator, f.KeyAddress(UserKevin).String())
		require.True(t, sbr.HasRole)
	}

	// an allocator can't be an issuer or a returner
	{
		// try registering kevin as a returner should fail
		resultReturner, _, _ := f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleReturner, f.KeyAddress(UserTina).String(), f.KeyAddress(UserKevin).String(), "-y")

		// try registering kevin as an issuer should fail
		resultIssuer, _, _ := f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleIssuer, f.KeyAddress(UserTina).String(), f.KeyAddress(UserKevin).String(), "-y")

		tests.WaitForNextNBlocksTM(1, f.Port)

		require.True(t, resultReturner)
		require.True(t, resultIssuer)

		// kevin is not a returner
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleReturner, f.KeyAddress(UserKevin).String())
		require.False(t, sbr.HasRole)

		// kevin is not an issuer
		sbr = f.QuerySafetyBoxRole(id, sbox.RoleIssuer, f.KeyAddress(UserKevin).String())
		require.False(t, sbr.HasRole)
	}

	// an allocator is able to allocate coins to the safety box
	{
		// kevin allocates 1link to the safety box
		result, _, _ = f.TxSafetyBoxSendCoins(id, sbox.ActionAllocate, DenomLink, int64(1), f.KeyAddress(UserKevin).String(), "", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)

		// check the safety box balance
		sb := f.QuerySafetyBox(id)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.OneInt()}}, sb.TotalAllocation)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.OneInt()}}, sb.CumulativeAllocation)
		require.Equal(t, sdk.Coins(nil), sb.TotalIssuance)

		// check the kevin's balance
		kevinAccount := f.QueryAccount(f.KeyAddress(UserKevin))
		require.Equal(t, kevinAccount.Coins,
			sdk.Coins{
				sdk.Coin{DenomLink, sdk.NewInt(999999999)},
				sdk.Coin{DenomStake, sdk.NewInt(100000000000000)},
			},
		)
	}

	// an operator registers an issuer
	{
		// tina, the operator registers brian as an issuer
		result, _, _ = f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleIssuer, f.KeyAddress(UserTina).String(), f.KeyAddress(UserBrian).String(), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)

		// brian is now an issuer
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleIssuer, f.KeyAddress(UserBrian).String())
		require.True(t, sbr.HasRole)
	}

	// an issuer can't be an allocator or a returner
	{
		// try registering brian as an allocator should fail
		resultAllocator, _, _ := f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleAllocator, f.KeyAddress(UserTina).String(), f.KeyAddress(UserBrian).String(), "-y")

		// try registering brian as a returner should fail
		resultReturner, _, _ := f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleReturner, f.KeyAddress(UserTina).String(), f.KeyAddress(UserBrian).String(), "-y")

		tests.WaitForNextNBlocksTM(1, f.Port)

		require.True(t, resultAllocator)
		require.True(t, resultReturner)

		// brian is not an allocator
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleAllocator, f.KeyAddress(UserBrian).String())
		require.False(t, sbr.HasRole)

		// brian is not a returner
		sbr = f.QuerySafetyBoxRole(id, sbox.RoleReturner, f.KeyAddress(UserBrian).String())
		require.False(t, sbr.HasRole)
	}

	// an issuer is able to issue coins from the safety box to itself
	{
		// brian issues 1link from the safety box to himself
		result, _, _ = f.TxSafetyBoxSendCoins(id, sbox.ActionIssue, DenomLink, int64(1), f.KeyAddress(UserBrian).String(), f.KeyAddress(UserBrian).String(), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)

		// check the safety box balance
		sb := f.QuerySafetyBox(id)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.OneInt()}}, sb.TotalAllocation)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.OneInt()}}, sb.CumulativeAllocation)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.OneInt()}}, sb.TotalIssuance)

		// check the brian's balance
		brianAccount := f.QueryAccount(f.KeyAddress(UserBrian))
		require.Equal(
			t,
			sdk.Coins{
				sdk.Coin{DenomLink, sdk.NewInt(1000000001)},
				sdk.Coin{DenomStake, sdk.NewInt(100000000000000)},
			},
			brianAccount.Coins,
		)
	}

	// an issuer is able to issue coins from the safety box to another issuer
	{
		// kevin allocates 1 link to the safety box
		_, _, _ = f.TxSafetyBoxSendCoins(id, sbox.ActionAllocate, DenomLink, int64(1), f.KeyAddress(UserKevin).String(), "", "-y")

		// tina, the operator registers sam as an issuer
		result, _, _ = f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleIssuer, f.KeyAddress(UserTina).String(), f.KeyAddress(UserSam).String(), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)

		// sam is now an issuer
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleIssuer, f.KeyAddress(UserSam).String())
		require.True(t, sbr.HasRole)

		// brian issues 1link from the safety box to Sam
		result, _, _ = f.TxSafetyBoxSendCoins(id, sbox.ActionIssue, DenomLink, int64(1), f.KeyAddress(UserBrian).String(), f.KeyAddress(UserSam).String(), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)

		// check the safety box balance
		sb := f.QuerySafetyBox(id)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.NewInt(2)}}, sb.TotalAllocation)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.NewInt(2)}}, sb.CumulativeAllocation)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.NewInt(2)}}, sb.TotalIssuance)

		// check the sam's balance
		samAccount := f.QueryAccount(f.KeyAddress(UserSam))
		require.Equal(
			t,
			sdk.Coins{
				sdk.Coin{DenomLink, sdk.NewInt(1000000001)},
				sdk.Coin{DenomStake, sdk.NewInt(100000000000000)},
			},
			samAccount.Coins,
		)
	}

	// an issuer try issuing coins from the safety box to non-issuer should fail
	{
		// sam issues 1link from the safety box to non-issuer, evelyn
		result, stdout, _ := f.TxSafetyBoxSendCoins(id, sbox.ActionIssue, DenomLink, int64(1), f.KeyAddress(UserSam).String(), f.KeyAddress(UserEvelyn).String(), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxPermissionIssue, "Account: %s", f.KeyAddress(UserEvelyn).String()).Error(), "\""),
			strings.Split(stdout, "\\\\\\\"")[9],
		)
	}

	// an operator registers a returner
	{
		// tina, the operator registers evelyn as a returner
		result, _, _ = f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleReturner, f.KeyAddress(UserTina).String(), f.KeyAddress(UserEvelyn).String(), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)

		// evelyn is now a returner
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleReturner, f.KeyAddress(UserEvelyn).String())
		require.True(t, sbr.HasRole)
	}

	// a returner can't be an issuer or an allocator
	{
		// try registering evelyn as an allocator should fail
		resultAllocator, _, _ := f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleAllocator, f.KeyAddress(UserTina).String(), f.KeyAddress(UserEvelyn).String(), "-y")

		// try registering evelyn as an issuer should fail
		resultIssuer, _, _ := f.TxSafetyBoxRole(id, sbox.RegisterRole, sbox.RoleIssuer, f.KeyAddress(UserTina).String(), f.KeyAddress(UserEvelyn).String(), "-y")

		tests.WaitForNextNBlocksTM(1, f.Port)

		require.True(t, resultAllocator)
		require.True(t, resultIssuer)

		// evelyn is not an allocator
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleAllocator, f.KeyAddress(UserEvelyn).String())
		require.False(t, sbr.HasRole)

		// evelyn is not an issuer
		sbr = f.QuerySafetyBoxRole(id, sbox.RoleIssuer, f.KeyAddress(UserEvelyn).String())
		require.False(t, sbr.HasRole)
	}

	// a returner is able to return coins to the safety box
	{
		// evelyn returns 1link to the safety box
		result, _, _ = f.TxSafetyBoxSendCoins(id, sbox.ActionReturn, DenomLink, int64(1), f.KeyAddress(UserEvelyn).String(), "", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)

		// check the safety box balance
		sb := f.QuerySafetyBox(id)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.NewInt(2)}}, sb.TotalAllocation)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.NewInt(2)}}, sb.CumulativeAllocation)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.NewInt(1)}}, sb.TotalIssuance)

		// check the evelyn's balance
		evelynAccount := f.QueryAccount(f.KeyAddress(UserEvelyn))
		require.Equal(
			t,
			sdk.Coins{
				sdk.Coin{DenomLink, sdk.NewInt(999999999)},
				sdk.Coin{DenomStake, sdk.NewInt(100000000000000)},
			},
			evelynAccount.Coins,
		)
	}

	// an allocator is able to recall coins from the safety box
	{
		// kevin recalls 1link from the safety box
		result, _, _ = f.TxSafetyBoxSendCoins(id, sbox.ActionRecall, DenomLink, int64(1), f.KeyAddress(UserKevin).String(), "", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)

		// check the safety box balance
		sb := f.QuerySafetyBox(id)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.NewInt(1)}}, sb.TotalAllocation)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.NewInt(2)}}, sb.CumulativeAllocation)
		require.Equal(t, sdk.Coins{sdk.Coin{DenomLink, sdk.NewInt(1)}}, sb.TotalIssuance)

		// check the kevin's balance
		kevinAccount := f.QueryAccount(f.KeyAddress(UserKevin))
		require.Equal(t,
			kevinAccount.Coins,
			sdk.Coins{
				sdk.Coin{DenomLink, sdk.NewInt(999999999)},
				sdk.Coin{DenomStake, sdk.NewInt(100000000000000)},
			},
		)
	}

	// can't allocate, recall, issue nor return coins other than specified denom in the safety box
	{
		// current balances
		sb := f.QuerySafetyBox(id)
		initialBalances := []sdk.Coins{
			sb.TotalAllocation,
			sb.CumulativeAllocation,
			sb.TotalIssuance,
			f.QueryAccount(f.KeyAddress(UserKevin)).Coins,  // kevin, an allocator
			f.QueryAccount(f.KeyAddress(UserBrian)).Coins,  // brian, an issuer
			f.QueryAccount(f.KeyAddress(UserSam)).Coins,    // sam, an issuer
			f.QueryAccount(f.KeyAddress(UserEvelyn)).Coins, // evelyn, a returner
		}

		// kevin allocates 1stake2 to the safety box, should fail
		result, stdout, _ := f.TxSafetyBoxSendCoins(id, sbox.ActionAllocate, DenomStake, int64(1), f.KeyAddress(UserKevin).String(), "", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxIncorrectDenom, "Expected: %s, Requested: %s", DenomLink, DenomStake).Error(), "\""),
			strings.Split(stdout, "\\\\\\\"")[9],
		)

		// brian issues 1stake2 from the safety box to Sam, should fail
		result, _, _ = f.TxSafetyBoxSendCoins(id, sbox.ActionIssue, DenomStake, int64(1), f.KeyAddress(UserBrian).String(), f.KeyAddress(UserSam).String(), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxIncorrectDenom, "Expected: %s, Requested: %s", DenomLink, DenomStake).Error(), "\""),
			strings.Split(stdout, "\\\\\\\"")[9],
		)

		// evelyn returns 1stake2 to the safety box, should fail
		result, _, _ = f.TxSafetyBoxSendCoins(id, sbox.ActionReturn, DenomStake, int64(1), f.KeyAddress(UserEvelyn).String(), "", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxIncorrectDenom, "Expected: %s, Requested: %s", DenomLink, DenomStake).Error(), "\""),
			strings.Split(stdout, "\\\\\\\"")[9],
		)

		// kevin recalls 1stake2 from the safety box, should fail
		result, _, _ = f.TxSafetyBoxSendCoins(id, sbox.ActionRecall, DenomStake, int64(1), f.KeyAddress(UserKevin).String(), "", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)
		require.True(t, result)
		require.Contains(
			t,
			strings.Split(sdkerrors.Wrapf(sbox.ErrSafetyBoxIncorrectDenom, "Expected: %s, Requested: %s", DenomLink, DenomStake).Error(), "\""),
			strings.Split(stdout, "\\\\\\\"")[9],
		)

		// no balance should have changed
		sb = f.QuerySafetyBox(id)
		finalBalances := []sdk.Coins{
			sb.TotalAllocation,
			sb.CumulativeAllocation,
			sb.TotalIssuance,
			f.QueryAccount(f.KeyAddress(UserKevin)).Coins,  // kevin, an allocator
			f.QueryAccount(f.KeyAddress(UserBrian)).Coins,  // brian, an issuer
			f.QueryAccount(f.KeyAddress(UserSam)).Coins,    // sam, an issuer
			f.QueryAccount(f.KeyAddress(UserEvelyn)).Coins, // evelyn, a returner
		}
		require.Equal(t, initialBalances, finalBalances)
	}

	// deregister roles
	{
		// tina, the operator deregisters kevin as an allocator
		result, _, _ = f.TxSafetyBoxRole(id, sbox.DeregisterRole, sbox.RoleAllocator, f.KeyAddress(UserTina).String(), f.KeyAddress(UserKevin).String(), "-y")
		require.True(t, result)

		// tina, the operator deregisters brian as an allocator
		result, _, _ = f.TxSafetyBoxRole(id, sbox.DeregisterRole, sbox.RoleIssuer, f.KeyAddress(UserTina).String(), f.KeyAddress(UserBrian).String(), "-y")
		require.True(t, result)

		// tina, the operator deregisters sam as an allocator
		result, _, _ = f.TxSafetyBoxRole(id, sbox.DeregisterRole, sbox.RoleIssuer, f.KeyAddress(UserTina).String(), f.KeyAddress(UserSam).String(), "-y")
		require.True(t, result)

		// tina, the operator deregisters evelyn as a returner
		result, _, _ = f.TxSafetyBoxRole(id, sbox.DeregisterRole, sbox.RoleReturner, f.KeyAddress(UserTina).String(), f.KeyAddress(UserEvelyn).String(), "-y")
		require.True(t, result)

		tests.WaitForNextNBlocksTM(1, f.Port)

		// kevin is not an allocator anymore
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleAllocator, f.KeyAddress(UserKevin).String())
		require.False(t, sbr.HasRole)

		// brian is not an issuer anymore
		sbr = f.QuerySafetyBoxRole(id, sbox.RoleIssuer, f.KeyAddress(UserBrian).String())
		require.False(t, sbr.HasRole)

		// sam is not an issuer anymore
		sbr = f.QuerySafetyBoxRole(id, sbox.RoleIssuer, f.KeyAddress(UserSam).String())
		require.False(t, sbr.HasRole)

		// evelyn is not a returner anymore
		sbr = f.QuerySafetyBoxRole(id, sbox.RoleReturner, f.KeyAddress(UserEvelyn).String())
		require.False(t, sbr.HasRole)
	}

	// deregister operator
	{
		// rinah, the owner of the safety box deregisters tina as an operator
		result, _, _ = f.TxSafetyBoxRole(id, sbox.DeregisterRole, sbox.RoleOperator, f.KeyAddress(UserRinah).String(), f.KeyAddress(UserTina).String(), "-y")
		require.True(t, result)

		tests.WaitForNextNBlocksTM(1, f.Port)

		// tina is not an operator anymore
		sbr := f.QuerySafetyBoxRole(id, sbox.RoleOperator, f.KeyAddress(UserTina).String())
		require.False(t, sbr.HasRole)
	}
}

func TestLinkCLIMempool(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

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

	// check mempool size
	{
		result := f.MempoolNumUnconfirmedTxs()
		require.Equal(t, 1, result.Count)
		require.Equal(t, 1, result.Total)
		require.Empty(t, result.Txs)
	}

	// check mempool txs
	{
		result := f.MempoolUnconfirmedTxHashes()
		require.Equal(t, 1, result.Count)
		require.Equal(t, 1, result.Total)

		txHash := UnmarshalTxResponse(t, stdout).TxHash
		require.Equal(t, txHash, result.Txs[0])
	}

	// Ensure account balances match expected
	tests.WaitForNextNBlocksTM(1, f.Port)

	barAcc := f.QueryAccount(barAddr)
	require.Equal(t, sendTokens, barAcc.GetCoins().AmountOf(denom))
	fooAcc = f.QueryAccount(fooAddr)
	require.Equal(t, startTokens.Sub(sendTokens), fooAcc.GetCoins().AmountOf(denom))

	f.Cleanup()
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

	// Query for all tokens
	{
		allTokens := f.QueryTokens()
		require.Equal(t, 2, len(allTokens))

		require.Equal(t, description, allTokens[0].GetName())
		require.Equal(t, contractID2, allTokens[0].GetContractID())
		require.Equal(t, int64(decimals), allTokens[0].GetDecimals().Int64())
		require.Equal(t, true, allTokens[0].GetMintable())

		require.Equal(t, description, allTokens[1].GetName())
		require.Equal(t, contractID1, allTokens[1].GetContractID())
		require.Equal(t, int64(decimals), allTokens[1].GetDecimals().Int64())
		require.Equal(t, false, allTokens[1].GetMintable())

	}

	// Query permissions for foo account
	{
		pms := f.QueryAccountPermission(f.KeyAddress(keyFoo), contractID1)
		require.Equal(t, 1, len(pms))
		require.Equal(t, contractID1, pms[0].GetResource())
		require.Equal(t, "modify", pms[0].GetAction())
	}
	{
		pms := f.QueryAccountPermission(f.KeyAddress(keyFoo), contractID2)
		require.Equal(t, 3, len(pms))
		require.Equal(t, contractID2, pms[0].GetResource())
		require.Equal(t, "modify", pms[0].GetAction())
		require.Equal(t, contractID2, pms[1].GetResource())
		require.Equal(t, "mint", pms[1].GetAction())
		require.Equal(t, contractID2, pms[2].GetResource())
		require.Equal(t, "burn", pms[2].GetAction())
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

	f.Cleanup()
}
func TestLinkCLITokenMintBurn(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

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
		//Foo try to burn but insufficient
		_, stdOut, _ := f.TxTokenBurn(keyFoo, contractID, initAmountStr+initAmountStr, "-y")
		require.Contains(t, stdOut, "not enough coins")
		//bar try to mint but has no permission
		_, stdOut, _ = f.TxTokenMint(keyBar, contractID, barAddr.String(), mintAmountStr, "-y")
		require.Contains(t, stdOut, "account does not have the permission")

		//Amount not changed
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

	f.Cleanup()
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

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)
	//barSuffix := types.AccAddrSuffix(barAddr)

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

	f.Cleanup()
}

func TestLinkCLITokenNFT(t *testing.T) {

	const (
		contractID  = "9be17165"
		tokenType   = "10000001"
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
		tests.WaitForNextNBlocksTM(1, f.Port)
		mintParam := strings.Join([]string{tokenType, description, meta}, ":")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam, "-y")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam, "-y")
		f.TxTokenMintNFTCollection(keyFoo, contractID, fooAddr.String(), mintParam, "-y")
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
	}

	f.Cleanup()
}

func TestLinkCLISendGenerateSignAndBroadcastWithToken(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

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

	f.Cleanup()
}

func TestLinkCLIEmpty(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()
	defer f.Cleanup()

	brianAddr := f.KeyAddress(UserBrian).String()

	t.Logf("[Account] Do nothing with empty message")
	success, stdout, stderr := f.TxEmpty(brianAddr, "-y")
	{
		require.True(t, success)
		require.NotEmpty(t, stdout)
		require.Empty(t, stderr)
	}
}
