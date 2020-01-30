// +build cli_multi_node_test

package clitest

import (
	"errors"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/privval"
)

func TestMultiValidatorAndSendTokens(t *testing.T) {
	t.Parallel()

	const (
		subnet = "192.168.0.0"
	)

	fg := InitFixturesGroup(t, subnet)

	fg.LDStopContainers()
	fg.LDStartContainers()

	defer fg.Cleanup()

	f := fg.Fixture(0)

	var (
		keyFoo = f.Moniker
	)

	fooAddr := f.KeyAddress(keyFoo)
	f.KeysDelete(keyBaz)
	f.KeysAdd(keyBaz)
	bazAddr := f.KeyAddress(keyBaz)

	fg.AddFullNode()
	{

		fooAcc := f.QueryAccount(fooAddr)
		startTokens := sdk.TokensFromConsensusPower(50)
		require.Equal(t, startTokens, fooAcc.GetCoins().AmountOf(denom))

		// Send some tokens from one account to the other
		sendTokens := sdk.TokensFromConsensusPower(10)
		f.TxSend(keyFoo, bazAddr, sdk.NewCoin(denom, sendTokens), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		// Ensure account balances match expected
		barAcc := f.QueryAccount(bazAddr)
		require.Equal(t, sendTokens, barAcc.GetCoins().AmountOf(denom))
		fooAcc = f.QueryAccount(fooAddr)
		require.Equal(t, startTokens.Sub(sendTokens), fooAcc.GetCoins().AmountOf(denom))

		// Test --dry-run
		success, _, _ := f.TxSend(keyFoo, bazAddr, sdk.NewCoin(denom, sendTokens), "--dry-run")
		require.True(t, success)

		// Test --generate-only
		success, stdout, stderr := f.TxSend(
			fooAddr.String(), bazAddr, sdk.NewCoin(denom, sendTokens), "--generate-only=true",
		)
		require.Empty(t, stderr)
		require.True(t, success)
		msg := UnmarshalStdTx(f.T, stdout)
		require.NotZero(t, msg.Fee.Gas)
		require.Len(t, msg.Msgs, 1)
		require.Len(t, msg.GetSignatures(), 0)

		// Check state didn't change
		fooAcc = f.QueryAccount(fooAddr)
		require.Equal(t, startTokens.Sub(sendTokens), fooAcc.GetCoins().AmountOf(denom))

		// test autosequencing
		f.TxSend(keyFoo, bazAddr, sdk.NewCoin(denom, sendTokens), "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		// Ensure account balances match expected
		barAcc = f.QueryAccount(bazAddr)
		require.Equal(t, sendTokens.MulRaw(2), barAcc.GetCoins().AmountOf(denom))
		fooAcc = f.QueryAccount(fooAddr)
		require.Equal(t, startTokens.Sub(sendTokens.MulRaw(2)), fooAcc.GetCoins().AmountOf(denom))

		// test memo
		f.TxSend(keyFoo, bazAddr, sdk.NewCoin(denom, sendTokens), "--memo='testmemo'", "-y")
		tests.WaitForNextNBlocksTM(1, f.Port)

		// Ensure account balances match expected
		barAcc = f.QueryAccount(bazAddr)
		require.Equal(t, sendTokens.MulRaw(3), barAcc.GetCoins().AmountOf(denom))
		fooAcc = f.QueryAccount(fooAddr)
		require.Equal(t, startTokens.Sub(sendTokens.MulRaw(3)), fooAcc.GetCoins().AmountOf(denom))
	}

}

func TestMultiValidatorAddNodeAndPromoteValidator(t *testing.T) {
	t.Parallel()

	const (
		subnet = "192.168.1.0"
	)

	fg := InitFixturesGroup(t, subnet)

	fg.LDStopContainers()
	fg.LDStartContainers()

	defer fg.Cleanup()

	f1 := fg.Fixture(0)

	f2 := fg.AddFullNode()

	{
		f2.KeysDelete(keyBar)
		f2.KeysAdd(keyBar)
	}

	barAddr := f2.KeyAddress(keyBar)
	barVal := sdk.ValAddress(barAddr)

	sendTokens := sdk.TokensFromConsensusPower(10)
	{
		f1.TxSend(f1.Moniker, barAddr, sdk.NewCoin(denom, sendTokens), "-y")
		tests.WaitForNextNBlocksTM(1, f1.Port)

		barAcc := f2.QueryAccount(barAddr)
		require.Equal(t, sendTokens, barAcc.GetCoins().AmountOf(denom))
	}

	newValTokens := sdk.TokensFromConsensusPower(2)
	{
		privVal := privval.LoadFilePVEmptyState(f2.PrivValidatorKeyFile(), "")
		consPubKey := sdk.MustBech32ifyConsPub(privVal.GetPubKey())

		f2.TxStakingCreateValidator(keyBar, consPubKey, sdk.NewCoin(denom, newValTokens), "-y")
		tests.WaitForNextNBlocksTM(1, f2.Port)
	}
	{
		// Ensure funds were deducted properly
		barAcc := f2.QueryAccount(barAddr)
		require.Equal(t, sendTokens.Sub(newValTokens), barAcc.GetCoins().AmountOf(denom))

		// Ensure that validator state is as expected
		validator := f2.QueryStakingValidator(barVal)
		require.Equal(t, validator.OperatorAddress, barVal)
		require.True(sdk.IntEq(t, newValTokens, validator.Tokens))

		// Query delegations to the validator
		validatorDelegations := f2.QueryStakingDelegationsTo(barVal)
		require.Len(t, validatorDelegations, 1)
		require.NotZero(t, validatorDelegations[0].Shares)
	}
}

func TestMultiValidatorAddNodeAndFailedTransactions(t *testing.T) {
	t.Parallel()

	const (
		subnet = "192.168.2.0"
	)

	fg := InitFixturesGroup(t, subnet)

	fg.LDStopContainers()
	fg.LDStartContainers()

	defer fg.Cleanup()

	f1 := fg.Fixture(0)

	var (
		keyFoo = f1.Moniker
	)

	fooAddr := f1.KeyAddress(keyFoo)

	{
		f1.KeysDelete(keyBaz)
		f1.KeysAdd(keyBaz)
	}

	bazAddr := f1.KeyAddress(keyBaz)

	f2 := fg.AddFullNode("--mempool.broadcast=false")

	{
		f2.KeysDelete(keyBar)
		f2.KeysAdd(keyBar)
	}

	barAddr := f2.KeyAddress(keyBar)

	// Send init token
	startTokens := sdk.TokensFromConsensusPower(10)
	{
		success, stdout, stderr := f1.TxSend(keyFoo, bazAddr, sdk.NewCoin(denom, startTokens), "-y")
		require.True(t, success)
		require.NotEmpty(t, stdout)
		require.Empty(t, stderr)
		tests.WaitForNextNBlocksTM(1, f1.Port)

		bazAcc := f2.QueryAccount(bazAddr)
		require.Equal(t, startTokens, bazAcc.GetCoins().AmountOf(denom))
	}
	{
		success, stdout, stderr := f1.TxSend(keyFoo, barAddr, sdk.NewCoin(denom, startTokens), "-y")
		require.True(t, success)
		require.NotEmpty(t, stdout)
		require.Empty(t, stderr)
		tests.WaitForNextNBlocksTM(1, f1.Port)

		barAcc := f2.QueryAccount(barAddr)
		require.Equal(t, startTokens, barAcc.GetCoins().AmountOf(denom))
	}

	// CheckTx failed
	{
		success, stdout, stderr := f1.TxSend(keyBaz, fooAddr, sdk.NewCoin(denom, startTokens.MulRaw(2)), "-y")
		require.True(t, success)
		require.Contains(t, stdout, "insufficient account funds")
		require.Empty(t, stderr)
		tests.WaitForNextNBlocksTM(1, f1.Port)

		bazAcc := f2.QueryAccount(bazAddr)
		require.Equal(t, startTokens, bazAcc.GetCoins().AmountOf(denom))
	}

	// CheckTx passed, but removed in mempool by Recheck
	{
		sendTokens := sdk.TokensFromConsensusPower(6)

		success, stdout1, stderr := f2.TxSend(keyBar, bazAddr, sdk.NewCoin(denom, sendTokens), "-y", "-b sync")
		require.True(t, success)
		require.NotEmpty(t, stdout1)
		require.Empty(t, stderr)

		// CheckTx passed
		sendResp1 := UnmarshalTxResponse(t, stdout1)
		require.Equal(t, sendResp1.Logs[0].Success, true)

		nodeOption := fmt.Sprintf("--node=%s", f1.RPCAddr)
		success, stdout2, stderr := f2.TxSend(keyBar, fooAddr, sdk.NewCoin(denom, sendTokens), "-y", nodeOption)
		require.True(t, success)
		require.NotEmpty(t, stdout2)
		require.Empty(t, stderr)

		// Commit new block
		sendResp2 := UnmarshalTxResponse(t, stdout2)
		require.Equal(t, sendResp2.Logs[0].Success, true)

		tests.WaitForNextNBlocksTM(1, f1.Port)

		barAcc := f2.QueryAccount(barAddr)
		require.Equal(t, startTokens.Sub(sendTokens), barAcc.GetCoins().AmountOf(denom))

		// Recheck and remove invalid Tx in mempool
		err := fmt.Sprintf(
			"ERROR: Tx: response error: RPC error -32603 - Internal error: Tx (%s) not found",
			sendResp1.TxHash,
		)
		f1.QueryTxInvalid(errors.New(err), sendResp1.TxHash)
	}

	// CheckTx passed, but DeliverTx failed
	{
		sendTokens := sdk.TokensFromConsensusPower(6)

		success, stdout1, stderr := f1.TxSend(keyBaz, fooAddr, sdk.NewCoin(denom, sendTokens), "-y", "-b sync", "-s 1")
		require.True(t, success)
		require.NotEmpty(t, stdout1)
		require.Empty(t, stderr)

		success, stdout2, stderr := f1.TxSend(keyBaz, barAddr, sdk.NewCoin(denom, sendTokens), "-y", "-b sync", "-s 2")
		require.True(t, success)
		require.NotEmpty(t, stdout2)
		require.Empty(t, stderr)

		// CheckTx results
		sendResp1 := UnmarshalTxResponse(t, stdout1)
		require.Equal(t, sendResp1.Logs[0].Success, true)

		sendResp2 := UnmarshalTxResponse(t, stdout2)
		require.Equal(t, sendResp2.Logs[0].Success, true)

		tests.WaitForNextNBlocksTM(1, f1.Port)

		// DeliverTx results
		txResult1 := f1.QueryTx(sendResp1.TxHash)
		require.Equal(t, txResult1.Logs[0].Success, true)

		txResult2 := f1.QueryTx(sendResp2.TxHash)
		require.Equal(t, txResult2.Logs[0].Success, false)

		bazAcc := f2.QueryAccount(bazAddr)
		require.Equal(t, startTokens.Sub(sendTokens), bazAcc.GetCoins().AmountOf(denom))
	}
}
