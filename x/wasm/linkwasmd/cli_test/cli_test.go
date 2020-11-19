// +build cli_test

package clitest

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	tokenModule "github.com/line/link-modules/x/token"
	"github.com/line/link-modules/x/wasm/linkwasmd/app"
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

	// instantiate a contract with the code escrow
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

func TestLinkCLIWasmTokenTester(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	fooAddr := f.KeyAddress(keyFoo)

	flagFromFoo := fmt.Sprintf("--from=%s", fooAddr)
	flagGas := "--gas=auto --gas-adjustment=1.2"
	workDir, _ := os.Getwd()
	wasmTokenTester := path.Join(workDir, "contracts", "token-tester", "contract.wasm")
	codeId := uint64(1)
	var contractAddress sdk.AccAddress
	tokenContractId := "9be17165"
	tokenName := "TestToken1"
	tokenSymbol := "TT1"
	tokenMeta := "meta"
	tokenImageURI := "http://example.com/image"
	tokenDecimals := sdk.NewInt(8)

	initAmount := 1
	mintAmount := 99
	transferAmount := 10
	burnAmount := 5
	mintFromFooAmount := 1
	burnFromFooAmount := 3

	modifyKey := "meta"
	modifyValue := "update_token_meta"

	// store the contract token-tester
	{
		f.LogResult(f.TxStoreWasm(wasmTokenTester, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// instantiate
	{
		msgJson := "{}"
		flagLabel := "--label=token-tester"
		f.LogResult(f.TxInstantiateWasm(codeId, msgJson, flagLabel, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate there is only one contract using codeId=1 and get contractAddress
	{
		listContract := f.QueryListContractByCodeWasm(codeId)
		require.Len(t, listContract, 1)
		contractAddress = listContract[0].Address
	}

	// issue token
	{
		msg := map[string]map[string]interface{}{
			"issue": {
				"owner":    contractAddress,
				"to":       contractAddress,
				"name":     tokenName,
				"symbol":   tokenSymbol,
				"meta":     tokenMeta,
				"img_uri":  tokenImageURI,
				"amount":   strconv.Itoa(initAmount),
				"mintable": true,
				"decimals": tokenDecimals,
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate that token is issued
	{
		token := f.QueryToken(tokenContractId)
		require.Equal(t, tokenContractId, token.GetContractID())
		require.Equal(t, tokenName, token.GetName())
		require.Equal(t, tokenSymbol, token.GetSymbol())
		require.Equal(t, tokenMeta, token.GetMeta())
		require.Equal(t, tokenImageURI, token.GetImageURI())
		require.Equal(t, tokenDecimals, token.GetDecimals())
		require.True(t, token.GetMintable())

		perms := f.QueryAccountPermission(contractAddress, tokenContractId)
		require.Len(t, perms, 3)
	}

	// mint token
	{
		msg := map[string]map[string]interface{}{
			"mint": {
				"from":        contractAddress,
				"contract_id": tokenContractId,
				"to":          contractAddress,
				"amount":      strconv.Itoa(mintAmount),
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate total supply, mint, burn
	{
		totalSupply := f.QuerySupplyToken(tokenContractId)
		require.Equal(t, sdk.NewInt(int64(initAmount+mintAmount)), totalSupply)
		totalMint := f.QueryMintToken(tokenContractId)
		require.Equal(t, sdk.NewInt(int64(initAmount+mintAmount)), totalMint)
		totalBurn := f.QueryBurnToken(tokenContractId)
		require.Equal(t, sdk.NewInt(0), totalBurn)
	}

	// transfer token
	{
		msg := map[string]map[string]interface{}{
			"transfer": {
				"from":        contractAddress,
				"contract_id": tokenContractId,
				"to":          fooAddr,
				"amount":      strconv.Itoa(transferAmount),
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate balance
	{
		balanceOfContract := f.QueryBalanceToken(tokenContractId, contractAddress)
		require.Equal(t, sdk.NewInt(int64(initAmount+mintAmount-transferAmount)), balanceOfContract)
		balanceOfFoo := f.QueryBalanceToken(tokenContractId, fooAddr)
		require.Equal(t, sdk.NewInt(int64(transferAmount)), balanceOfFoo)
	}

	// burn token
	{
		msg := map[string]map[string]interface{}{
			"burn": {
				"from":        contractAddress,
				"contract_id": tokenContractId,
				"amount":      strconv.Itoa(burnAmount),
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate total supply, mint, burn
	{
		totalSupply := f.QuerySupplyToken(tokenContractId)
		require.Equal(t, sdk.NewInt(int64(initAmount+mintAmount-burnAmount)), totalSupply)
		totalMint := f.QueryMintToken(tokenContractId)
		require.Equal(t, sdk.NewInt(int64(initAmount+mintAmount)), totalMint)
		totalBurn := f.QueryBurnToken(tokenContractId)
		require.Equal(t, sdk.NewInt(int64(burnAmount)), totalBurn)
	}

	// validate that fooAddr cannot mint token
	{
		perms := f.QueryAccountPermission(fooAddr, tokenContractId)
		require.Len(t, perms, 0)
		_, res, _ := f.TxTokenMint(fooAddr.String(), tokenContractId, fooAddr.String(), strconv.Itoa(mintFromFooAmount), "-y")
		require.True(t, strings.Contains(res, "Permission: mint: failed"))
	}

	// grant permission for fooAddr
	{
		msg := map[string]map[string]interface{}{
			"grant_perm": {
				"from":        contractAddress,
				"contract_id": tokenContractId,
				"to":          fooAddr,
				"permission":  "mint",
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate permission for fooAddr
	{
		perms := f.QueryAccountPermission(fooAddr, tokenContractId)
		require.Len(t, perms, 1)
		require.Equal(t, tokenModule.Permission("mint"), perms[0])
	}

	// mint token from fooAddr
	{
		f.LogResult(f.TxTokenMint(fooAddr.String(), tokenContractId, fooAddr.String(), strconv.Itoa(mintFromFooAmount), "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate balance for fooAddr
	{
		balance := f.QueryBalanceToken(tokenContractId, fooAddr)
		require.Equal(t, sdk.NewInt(int64(transferAmount+mintFromFooAmount)), balance)
	}

	// revoke permission from contractAddress
	{
		msg := map[string]map[string]interface{}{
			"revoke_perm": {
				"from":        contractAddress,
				"contract_id": tokenContractId,
				"permission":  "mint",
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate permission from contractAddress
	{
		perms := f.QueryAccountPermission(contractAddress, tokenContractId)
		require.Len(t, perms, 2)
		require.True(t, perms.HasPermission(tokenModule.NewBurnPermission()))
		require.True(t, perms.HasPermission(tokenModule.NewModifyPermission()))
	}

	// validate that contractAddress cannot mint token
	{
		msg := map[string]map[string]interface{}{
			"mint": {
				"from":        contractAddress,
				"contract_id": tokenContractId,
				"to":          contractAddress,
				"amount":      strconv.Itoa(mintAmount),
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		_, _, errStr := f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y")
		require.True(t, strings.Contains(errStr, "Permission: mint: failed"))
	}

	// modify token
	{
		msg := map[string]map[string]interface{}{
			"modify": {
				"owner":       contractAddress,
				"contract_id": tokenContractId,
				"key":         modifyKey,
				"value":       modifyValue,
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate that updated token
	{
		token := f.QueryToken(tokenContractId)
		require.Equal(t, tokenContractId, token.GetContractID())
		require.Equal(t, tokenName, token.GetName())
		require.Equal(t, tokenSymbol, token.GetSymbol())
		require.Equal(t, modifyValue, token.GetMeta())
		require.Equal(t, tokenImageURI, token.GetImageURI())
		require.Equal(t, tokenDecimals, token.GetDecimals())
		require.True(t, token.GetMintable())
	}

	// validate that fooAddr cannot burn token
	{
		perms := f.QueryAccountPermission(fooAddr, tokenContractId)
		require.Len(t, perms, 1)
		_, res, _ := f.TxTokenBurnFrom(fooAddr.String(), tokenContractId, contractAddress, int64(burnFromFooAmount), "-y")
		require.True(t, strings.Contains(res, "proxy is not approved on the token"))
	}

	// approve and grant burn perm
	{
		msg := map[string]map[string]interface{}{
			"approve": {
				"approver":    contractAddress,
				"contract_id": tokenContractId,
				"proxy":       fooAddr,
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		msg = map[string]map[string]interface{}{
			"grant_perm": {
				"from":        contractAddress,
				"contract_id": tokenContractId,
				"to":          fooAddr,
				"permission":  "burn",
			},
		}
		msgJson, _ = json.Marshal(msg)
		msgString = string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate approved and perm
	{
		res := f.QueryApprovedToken(tokenContractId, fooAddr, contractAddress)
		require.True(t, res)

		perms := f.QueryAccountPermission(fooAddr, tokenContractId)
		require.Len(t, perms, 2)
		require.True(t, perms.HasPermission(tokenModule.NewMintPermission()))
		require.True(t, perms.HasPermission(tokenModule.NewBurnPermission()))
	}

	// burn from fooAddr
	{
		f.LogResult(f.TxTokenBurnFrom(fooAddr.String(), tokenContractId, contractAddress, int64(burnFromFooAmount), "-y"))
		balanceOfContract := f.QueryBalanceToken(tokenContractId, contractAddress)
		require.Equal(t, sdk.NewInt(int64(initAmount+mintAmount-transferAmount-burnAmount-burnFromFooAmount)), balanceOfContract)

	}

	// test for query
	{
		cdc := app.MakeCodec()

		// query tokens
		query := map[string]map[string]interface{}{
			"get_token": {
				"contract_id": tokenContractId,
			},
		}
		queryJson, _ := json.Marshal(query)
		queryString := string(queryJson)
		res := f.QueryContractStateSmartWasm(contractAddress, queryString)
		var token tokenModule.Token
		tokenModule.ModuleCdc.UnmarshalJSON([]byte(res), &token)

		require.Equal(t, tokenContractId, token.GetContractID())
		require.Equal(t, tokenName, token.GetName())
		require.Equal(t, tokenSymbol, token.GetSymbol())
		require.Equal(t, modifyValue, token.GetMeta())
		require.Equal(t, tokenImageURI, token.GetImageURI())
		require.Equal(t, tokenDecimals, token.GetDecimals())
		require.True(t, token.GetMintable())

		// query balances
		query = map[string]map[string]interface{}{
			"get_balance": {
				"contract_id": tokenContractId,
				"address":     contractAddress,
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var balance sdk.Int
		err := cdc.UnmarshalJSON([]byte(res), &balance)
		require.NoError(f.T, err)
		require.Equal(t, sdk.NewInt(int64(initAmount+mintAmount-transferAmount-burnAmount-burnFromFooAmount)), balance)

		// query total
		query = map[string]map[string]interface{}{
			"get_total": {
				"contract_id": tokenContractId,
				"target":      "mint",
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var totalMint sdk.Int
		err = cdc.UnmarshalJSON([]byte(res), &totalMint)
		require.NoError(f.T, err)
		require.Equal(t, sdk.NewInt(int64(initAmount+mintAmount+mintFromFooAmount)), totalMint)

		query = map[string]map[string]interface{}{
			"get_total": {
				"contract_id": tokenContractId,
				"target":      "burn",
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var totalBurn sdk.Int
		err = cdc.UnmarshalJSON([]byte(res), &totalBurn)
		require.NoError(f.T, err)
		require.Equal(t, sdk.NewInt(int64(burnAmount+burnFromFooAmount)), totalBurn)

		query = map[string]map[string]interface{}{
			"get_total": {
				"contract_id": tokenContractId,
				"target":      "supply",
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var totalSupply sdk.Int
		err = cdc.UnmarshalJSON([]byte(res), &totalSupply)
		require.NoError(f.T, err)
		require.Equal(t, sdk.NewInt(int64(initAmount+mintAmount+mintFromFooAmount-burnAmount-burnFromFooAmount)), totalSupply)

		// query perm
		query = map[string]map[string]interface{}{
			"get_perm": {
				"contract_id": tokenContractId,
				"address":     contractAddress,
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var perms tokenModule.Permissions
		err = cdc.UnmarshalJSON([]byte(res), &perms)
		require.NoError(f.T, err)
		require.Equal(t, 2, len(perms))
	}
}

func TestLinkCLIWasmTokenTesterProxy(t *testing.T) {

	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	fooAddr := f.KeyAddress(keyFoo)

	flagFromFoo := fmt.Sprintf("--from=%s", fooAddr)
	flagGas := "--gas=auto --gas-adjustment=1.2"
	workDir, _ := os.Getwd()
	wasmTokenTester := path.Join(workDir, "contracts", "token-tester", "contract.wasm")
	codeId := uint64(1)
	var contractAddress sdk.AccAddress
	tokenContractId := "9be17165"
	tokenName := "TestToken2"
	tokenSymbol := "TT2"
	tokenMeta := "meta"

	// store the contract token-tester
	{
		f.LogResult(f.TxStoreWasm(wasmTokenTester, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// instantiate
	{
		msgJson := "{}"
		flagLabel := "--label=token-tester"
		f.LogResult(f.TxInstantiateWasm(codeId, msgJson, flagLabel, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate there is only one contract using codeId=1 and get contractAddress
	{
		listContract := f.QueryListContractByCodeWasm(codeId)
		require.Len(t, listContract, 1)
		contractAddress = listContract[0].Address
	}

	// set test token and approve for contractAddress
	{
		f.LogResult(f.TxTokenIssue(fooAddr.String(), fooAddr, tokenName, tokenMeta, tokenSymbol, 10000, 6, true, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		f.LogResult(f.TxTokenApprove(fooAddr.String(), tokenContractId, contractAddress, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		f.LogResult(f.TxTokenGrantPerm(fooAddr.String(), contractAddress.String(), tokenContractId, "burn", "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		res := f.QueryApprovedToken(tokenContractId, contractAddress, fooAddr)
		require.True(t, res)

	}

	// burn token from proxy
	{
		msgJson := fmt.Sprintf("{\"burn_from\":{\"proxy\":\"%s\",\"from\":\"%s\",\"contract_id\":\"%s\",\"amount\":\"%s\"}}", contractAddress, fooAddr, tokenContractId, "2")
		f.LogResult(f.TxExecuteWasm(contractAddress, msgJson, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate balance for fooAddr
	{
		balanceOfContract := f.QueryBalanceToken(tokenContractId, fooAddr)
		require.Equal(t, sdk.NewInt(9998), balanceOfContract)
	}

	// transfer token from proxy
	{
		msgJson := fmt.Sprintf("{\"transfer_from\":{\"proxy\":\"%s\",\"from\":\"%s\",\"contract_id\":\"%s\",\"to\":\"%s\",\"amount\":\"%s\"}}", contractAddress, fooAddr, tokenContractId, contractAddress, "3")
		f.LogResult(f.TxExecuteWasm(contractAddress, msgJson, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate balance
	{
		balanceOfFoo := f.QueryBalanceToken(tokenContractId, fooAddr)
		require.Equal(t, sdk.NewInt(9995), balanceOfFoo)
		balanceOfContract := f.QueryBalanceToken(tokenContractId, contractAddress)
		require.Equal(t, sdk.NewInt(3), balanceOfContract)
	}

	// test for query
	{
		cdc := app.MakeCodec()

		// query isApproved
		query := map[string]map[string]interface{}{
			"get_is_approved": {
				"proxy":       contractAddress,
				"contract_id": tokenContractId,
				"approver":    fooAddr,
			},
		}
		queryJson, _ := json.Marshal(query)
		queryString := string(queryJson)
		res := f.QueryContractStateSmartWasm(contractAddress, queryString)
		var isApproved bool
		err := cdc.UnmarshalJSON([]byte(res), &isApproved)
		require.NoError(f.T, err)
		require.True(t, isApproved)

		// query approvers
		query = map[string]map[string]interface{}{
			"get_approvers": {
				"proxy":       contractAddress,
				"contract_id": tokenContractId,
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var approvers []sdk.AccAddress
		err = cdc.UnmarshalJSON([]byte(res), &approvers)
		require.NoError(f.T, err)
		require.Equal(t, fooAddr, approvers[0])
	}
}

func TestLinkCLIWasmCreateCollection(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)
	defer f.Cleanup()

	// start linkd server
	proc := f.LDStart()
	defer func() { require.NoError(t, proc.Stop(false)) }()

	fooAddr := f.KeyAddress(keyFoo)

	flagFromFoo := fmt.Sprintf("--from=%s", fooAddr)
	flagGas := "--gas=auto --gas-adjustment=1.2"
	workDir, _ := os.Getwd()
	wasmCollectionTester := path.Join(workDir, "contracts", "collection-tester", "contract.wasm")
	codeId := uint64(1)
	var contractAddress sdk.AccAddress
	collectionContractId := "9be17165"
	collectionName := "TestCollection1"
	collectionMeta := "meta"
	collectionBaseImageURI := "http://example.com/image"

	// store the contract collection-tester
	{
		f.LogResult(f.TxStoreWasm(wasmCollectionTester, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// instantiate
	{
		msgString := "{}"
		flagLabel := "--label=collection-tester"
		f.LogResult(f.TxInstantiateWasm(codeId, msgString, flagLabel, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate there is only one contract using codeId=1 and get contractAddress
	{
		listContract := f.QueryListContractByCodeWasm(codeId)
		require.Len(t, listContract, 1)
		contractAddress = listContract[0].Address
	}

	// create collection
	{
		msg := map[string]map[string]interface{}{
			"create": {
				"owner":        contractAddress,
				"name":         collectionName,
				"meta":         collectionMeta,
				"base_img_uri": collectionBaseImageURI,
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate that collection is issued
	{
		collection := f.QueryCollection(collectionContractId)
		require.Equal(t, collectionContractId, collection.GetContractID())
		require.Equal(t, collectionName, collection.GetName())
		require.Equal(t, collectionMeta, collection.GetMeta())
		require.Equal(t, collectionBaseImageURI, collection.GetBaseImgURI())
	}
}
