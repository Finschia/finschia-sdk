// +build cli_test

package clitest

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestLinkCLIWasmEscrow(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	flagFromFoo := fmt.Sprintf("--from=%s", fooAddr)
	flagFromBar := fmt.Sprintf("--from=%s", barAddr)
	flagGas := "--gas=auto --gas-adjustment=1.2"
	workDir, _ := os.Getwd()
	tmpDir := path.Join(workDir, "tmp-dir-for-test-escrow")
	wasmEscrow := path.Join(workDir, "contracts", "escrow-v6", "contract.wasm")
	codeId := uint64(1)
	amountSend := uint64(10)
	denomSend := fooDenom
	approvalMsgJson := fmt.Sprintf("{\"approve\":{\"quantity\":[{\"amount\":\"%d\",\"denom\":\"%s\"}]}}", amountSend, denomSend)
	var contractAddress sdk.AccAddress

	// make tmpDir
	os.Mkdir(tmpDir, os.ModePerm)

	// get init amount
	initAmountOfFoo := f.QueryAccount(fooAddr).Coins.AmountOf(denomSend).Uint64()
	initAmountOfBar := uint64(0)

	// validate that there are no code in the chain
	{
		listCode := f.QueryListCodeWasm()
		require.Empty(t, listCode)
	}

	// store the contract escrow
	{
		f.LogResult(f.TxStoreWasm(wasmEscrow, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate the code is stored
	{
		listCode := f.QueryListCodeWasm()
		require.Len(t, listCode, 1)
		// TODO after #23: validate the hash is same with the wasm file
	}

	// validate getCode get the exact same wasm
	{
		outputPath := path.Join(tmpDir, "escrow-tmp.wasm")
		f.QueryCodeWasm(codeId, outputPath)
		fLocal, _ := os.Open(wasmEscrow)
		fChain, _ := os.Open(outputPath)

		// 2000000 is enough length
		dataLocal := make([]byte, 2000000)
		dataChain := make([]byte, 2000000)
		fLocal.Read(dataLocal)
		fChain.Read(dataChain)
		require.Equal(t, dataLocal, dataChain)
	}

	// validate that there are no contract using the code (id=1)
	{
		listContract := f.QueryListContractByCodeWasm(codeId)
		require.Empty(t, listContract)
	}

	// instanciate a contract with the code escrow
	{
		msgJson := fmt.Sprintf("{\"arbiter\":\"%s\",\"recipient\":\"%s\"}", fooAddr, barAddr)
		flagLabel := "--label=escrow-test"
		flagAmount := fmt.Sprintf("--amount=%d%s", amountSend, denomSend)
		f.LogResult(f.TxInstantiateWasm(codeId, msgJson, flagLabel, flagAmount, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate foo's amount decreased
	{
		amount := f.QueryAccount(fooAddr).Coins.AmountOf(denomSend).Uint64()
		require.Equal(t, initAmountOfFoo-amountSend, amount)
	}

	// validate there is only one contract using codeId=1 and get contractAddress
	{
		listContract := f.QueryListContractByCodeWasm(codeId)
		require.Len(t, listContract, 1)
		contractAddress = listContract[0].Address
	}

	// check arbiter with query
	{
		res := f.QueryContractStateSmartWasm(contractAddress, "{\"arbiter\":{}}")
		require.Equal(t, fmt.Sprintf("{\"arbiter\":\"%s\"}", fooAddr), res)
	}

	// validate executing approve is failed by invalid account
	{
		succeeded, _, _ := f.TxExecuteWasm(contractAddress, approvalMsgJson, flagFromBar, flagGas, "-y")
		require.False(t, succeeded)
	}

	// execute approve
	{
		f.LogResult(f.TxExecuteWasm(contractAddress, approvalMsgJson, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate the coin Foot is transfered
	{
		amount := f.QueryAccount(barAddr).Coins.AmountOf(denomSend).Uint64()
		require.Equal(t, initAmountOfBar+amountSend, amount)
	}

	// validate approve over amount does not succeed
	{
		succeeded, _, _ := f.TxExecuteWasm(contractAddress, approvalMsgJson, flagFromFoo, flagGas, "-y")
		require.False(t, succeeded)
	}

	// remove tmp dir
	os.RemoveAll(tmpDir)
}
