// +build cli_test

package clitest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	collectionModule "github.com/line/lbm-sdk/x/collection"
	tokenModule "github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/wasm/linkwasmd/app"
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
	dirContract := path.Join(workDir, "contracts", "escrow")
	hashFile := path.Join(dirContract, "hash.txt")
	wasmEscrow := path.Join(dirContract, "contract.wasm")
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

		//validate the hash is the same
		expectedRow, _ := ioutil.ReadFile(hashFile)
		expected, err := hex.DecodeString(string(expectedRow[:64]))
		require.NoError(t, err)
		actual := listCode[0].GetDataHash().Bytes()
		require.Equal(t, expected, actual)
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
		contractAddress = listContract[0].GetAddress()
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
		contractAddress = listContract[0].GetAddress()
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
		contractAddress = listContract[0].GetAddress()
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

func TestLinkCLIWasmCollectionTester(t *testing.T) {
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

	cdc := app.MakeCodec()

	nftName1 := "TestNFT1"
	nftMeta1 := "nftMeta1"
	nftName2 := "TestNFT2"
	nftMeta2 := "nftMeta2"
	nftName3 := "TestNFT3"
	nftMeta3 := "nftMeta3"

	tokenTypeID1 := "10000001"
	tokenTypeID2 := "10000002"
	tokenTypeID3 := "10000003"

	index1 := "00000001"
	index2 := "00000002"
	nft0ID := tokenTypeID1 + index1
	nft1ID := tokenTypeID1 + index2
	nft2ID := tokenTypeID2 + index1
	nft3ID := tokenTypeID3 + index1
	mintNftName := "nft-0"
	mintNftMeta := ""

	ftName := "TestFT1"
	ftMeta := "ftMeta"
	tokenID := "0000000100000000"

	initFtAmount := 1
	mintFtAmount := 99
	burnFtAmount := 3
	transferFtAmount := 10
	mintFt := strconv.Itoa(mintFtAmount) + ":" + tokenID
	burnFt := strconv.Itoa(burnFtAmount) + ":" + tokenID
	transferFt := strconv.Itoa(transferFtAmount) + ":" + tokenID

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
		contractAddress = listContract[0].GetAddress()
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

	// issue nft
	{
		msg := map[string]map[string]interface{}{
			"issue_nft": {
				"owner":       contractAddress,
				"contract_id": collectionContractId,
				"name":        nftName1,
				"meta":        nftMeta1,
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate issued nft
	{
		tokenType := f.QueryTokenTypeCollection(collectionContractId, tokenTypeID1)
		require.Equal(t, nftName1, tokenType.GetName())
		require.Equal(t, nftMeta1, tokenType.GetMeta())
		require.Equal(t, collectionContractId, tokenType.GetContractID())
		require.Equal(t, tokenTypeID1, tokenType.GetTokenType())
	}

	// issue ft
	{
		msg := map[string]map[string]interface{}{
			"issue_ft": {
				"owner":       contractAddress,
				"contract_id": collectionContractId,
				"to":          contractAddress,
				"name":        ftName,
				"meta":        ftMeta,
				"amount":      strconv.Itoa(initFtAmount),
				"mintable":    true,
				"decimals":    sdk.NewInt(8),
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate issued ft
	{
		token := f.QueryTokenCollection(collectionContractId, tokenID).(collectionModule.FT)
		require.Equal(t, collectionContractId, token.GetContractID())
		require.Equal(t, ftName, token.GetName())
		require.Equal(t, ftMeta, token.GetMeta())
		require.Equal(t, tokenID, token.GetTokenID())
		require.Equal(t, sdk.NewInt(8), token.GetDecimals())
		require.True(t, token.GetMintable())

		balanceOfContract := f.QueryBalanceCollection(collectionContractId, tokenID, contractAddress)
		require.Equal(t, sdk.NewInt(int64(initFtAmount)), balanceOfContract)

		totalMint := f.QueryTotalMintTokenCollection(collectionContractId, tokenID)
		require.Equal(t, sdk.NewInt(int64(initFtAmount)), totalMint)

		totalSupply := f.QueryTotalSupplyTokenCollection(collectionContractId, tokenID)
		require.Equal(t, sdk.NewInt(int64(initFtAmount)), totalSupply)

		totalBurn := f.QueryTotalBurnTokenCollection(collectionContractId, tokenID)
		require.Equal(t, sdk.NewInt(0), totalBurn)
	}

	// mint nft
	{
		msg := map[string]map[string]interface{}{
			"mint_nft": {
				"from":        contractAddress,
				"contract_id": collectionContractId,
				"to":          contractAddress,
				"token_types": []string{tokenTypeID1},
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate minted nft
	{
		nft := f.QueryTokenCollection(collectionContractId, nft0ID).(collectionModule.NFT)
		require.Equal(t, collectionContractId, nft.GetContractID())
		require.Equal(t, nft0ID, nft.GetTokenID())
		require.Equal(t, contractAddress, nft.GetOwner())
		require.Equal(t, mintNftName, nft.GetName())
		require.Equal(t, mintNftMeta, nft.GetMeta())
	}

	// mint ft
	{
		msg := map[string]map[string]interface{}{
			"mint_ft": {
				"from":        contractAddress,
				"contract_id": collectionContractId,
				"to":          contractAddress,
				"tokens":      []string{mintFt},
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate minted ft
	{
		balanceOfContract := f.QueryBalanceCollection(collectionContractId, tokenID, contractAddress)
		require.Equal(t, sdk.NewInt(int64(initFtAmount+mintFtAmount)), balanceOfContract)
	}

	// modify token info
	{
		msg := map[string]map[string]interface{}{
			"modify": {
				"owner":       contractAddress,
				"contract_id": collectionContractId,
				"token_type":  tokenTypeID1,
				"token_index": index1,
				"key":         "meta",
				"value":       "modified_meta",
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate modified nft
	{
		nft := f.QueryTokenCollection(collectionContractId, nft0ID).(collectionModule.NFT)
		require.Equal(t, collectionContractId, nft.GetContractID())
		require.Equal(t, nft0ID, nft.GetTokenID())
		require.Equal(t, contractAddress, nft.GetOwner())
		require.Equal(t, mintNftName, nft.GetName())
		require.Equal(t, "modified_meta", nft.GetMeta())
	}

	// burn nft
	{
		msg := map[string]map[string]interface{}{
			"burn_nft": {
				"from":        contractAddress,
				"contract_id": collectionContractId,
				"token_id":    nft0ID,
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate burn nft
	{
		f.QueryTokenCollectionExpectEmpty(collectionContractId, nft0ID)
	}

	// burn ft
	{
		msg := map[string]map[string]interface{}{
			"burn_ft": {
				"from":        contractAddress,
				"contract_id": collectionContractId,
				"amounts":     []string{burnFt},
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate burn ft
	{
		balanceOfContract := f.QueryBalanceCollection(collectionContractId, tokenID, contractAddress)
		require.Equal(t, sdk.NewInt(int64(initFtAmount+mintFtAmount-burnFtAmount)), balanceOfContract)

		totalMint := f.QueryTotalMintTokenCollection(collectionContractId, tokenID)
		require.Equal(t, sdk.NewInt(int64(initFtAmount+mintFtAmount)), totalMint)

		totalSupply := f.QueryTotalSupplyTokenCollection(collectionContractId, tokenID)
		require.Equal(t, sdk.NewInt(int64(initFtAmount+mintFtAmount-burnFtAmount)), totalSupply)

		totalBurn := f.QueryTotalBurnTokenCollection(collectionContractId, tokenID)
		require.Equal(t, sdk.NewInt(int64(burnFtAmount)), totalBurn)
	}

	// transfer nft
	{
		// prepare nft
		msg := map[string]map[string]interface{}{
			"mint_nft": {
				"from":        contractAddress,
				"contract_id": collectionContractId,
				"to":          contractAddress,
				"token_types": []string{tokenTypeID1},
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// transfer nft to fooAddr
		msg = map[string]map[string]interface{}{
			"transfer_nft": {
				"from":        contractAddress,
				"contract_id": collectionContractId,
				"to":          fooAddr,
				"token_ids":   []string{nft1ID},
			},
		}
		msgJson, _ = json.Marshal(msg)
		msgString = string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate transfered nft
	{
		nft := f.QueryTokenCollection(collectionContractId, nft1ID).(collectionModule.NFT)
		require.Equal(t, collectionContractId, nft.GetContractID())
		require.Equal(t, nft1ID, nft.GetTokenID())
		require.Equal(t, fooAddr, nft.GetOwner())
		require.Equal(t, mintNftName, nft.GetName())
		require.Equal(t, mintNftMeta, nft.GetMeta())
	}

	//  transfer ft
	{
		msg := map[string]map[string]interface{}{
			"transfer_ft": {
				"from":        contractAddress,
				"contract_id": collectionContractId,
				"to":          fooAddr,
				"tokens":      []string{transferFt},
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate transfered ft
	{
		balanceOfContract := f.QueryBalanceCollection(collectionContractId, tokenID, contractAddress)
		require.Equal(t, sdk.NewInt(int64(initFtAmount+mintFtAmount-burnFtAmount-transferFtAmount)), balanceOfContract)
		balanceOfFoo := f.QueryBalanceCollection(collectionContractId, tokenID, fooAddr)
		require.Equal(t, sdk.NewInt(int64(transferFtAmount)), balanceOfFoo)
	}

	// validate that contractAddress has all the permissions
	{
		perms := f.QueryAccountPermissionCollection(contractAddress, collectionContractId)
		require.Equal(t, 4, len(perms))
		require.True(t, perms.HasPermission(collectionModule.NewMintPermission()))
		require.True(t, perms.HasPermission(collectionModule.NewBurnPermission()))
		require.True(t, perms.HasPermission(collectionModule.NewIssuePermission()))
		require.True(t, perms.HasPermission(collectionModule.NewModifyPermission()))
	}

	// validate that fooAddr cannot mint token
	{
		perms := f.QueryAccountPermissionCollection(fooAddr, collectionContractId)
		require.Len(t, perms, 0)
		mintParam := strings.Join([]string{tokenTypeID1, "description", "meta"}, ":")
		_, res, _ := f.TxTokenMintNFTCollection(fooAddr.String(), collectionContractId, fooAddr.String(), mintParam, "-y")
		require.True(t, strings.Contains(res, "Permission: mint: failed"))
	}

	// grant permission
	{
		msg := map[string]map[string]interface{}{
			"grant_perm": {
				"from":        contractAddress,
				"contract_id": collectionContractId,
				"to":          fooAddr,
				"permission":  "mint",
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate perm
	{
		perms := f.QueryAccountPermissionCollection(fooAddr, collectionContractId)
		require.Equal(t, 1, len(perms))
		require.True(t, perms.HasPermission(collectionModule.NewMintPermission()))

		f.TxTokenMintFTCollection(fooAddr.String(), collectionContractId, fooAddr.String(), mintFt, "-y")
		balanceOfFoo := f.QueryBalanceCollection(collectionContractId, tokenID, fooAddr)
		require.Equal(t, sdk.NewInt(int64(transferFtAmount+mintFtAmount)), balanceOfFoo)
	}

	// revoke perm
	{
		msg := map[string]map[string]interface{}{
			"revoke_perm": {
				"from":        contractAddress,
				"contract_id": collectionContractId,
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
		perms := f.QueryAccountPermissionCollection(contractAddress, collectionContractId)
		require.Len(t, perms, 3)
		require.True(t, perms.HasPermission(collectionModule.NewBurnPermission()))
		require.True(t, perms.HasPermission(collectionModule.NewIssuePermission()))
		require.True(t, perms.HasPermission(collectionModule.NewModifyPermission()))
	}

	// attach token
	{
		// prepare token
		msg := map[string]map[string]interface{}{
			"issue_nft": {
				"owner":       contractAddress,
				"contract_id": collectionContractId,
				"name":        nftName2,
				"meta":        nftMeta2,
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		msg = map[string]map[string]interface{}{
			"issue_nft": {
				"owner":       contractAddress,
				"contract_id": collectionContractId,
				"name":        nftName3,
				"meta":        nftMeta3,
			},
		}
		msgJson, _ = json.Marshal(msg)
		msgString = string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		msg = map[string]map[string]interface{}{
			"mint_nft": {
				"from":        contractAddress,
				"contract_id": collectionContractId,
				"to":          contractAddress,
				"token_types": []string{tokenTypeID2, tokenTypeID3},
			},
		}
		msgJson, _ = json.Marshal(msg)
		msgString = string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// attach
		msg = map[string]map[string]interface{}{
			"attach": {
				"from":        contractAddress,
				"contract_id": collectionContractId,
				"to_token_id": nft3ID,
				"token_id":    nft2ID,
			},
		}
		msgJson, _ = json.Marshal(msg)
		msgString = string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate attached token
	{
		parent := f.QueryParentTokenCollection(collectionContractId, nft2ID)
		require.Equal(t, nft3ID, parent.GetTokenID())

		children := f.QueryChildrenTokenCollection(collectionContractId, nft3ID)
		require.Equal(t, 1, len(children))
		require.Equal(t, nft2ID, children[0].GetTokenID())

		// validate for query encoder
		query := map[string]map[string]interface{}{
			"get_root_or_parent_or_children": {
				"contract_id": collectionContractId,
				"token_id":    nft2ID,
				"target":      "parent",
			},
		}
		queryJson, _ := json.Marshal(query)
		queryString := string(queryJson)
		res := f.QueryContractStateSmartWasm(contractAddress, queryString)

		var parentToken collectionModule.Token
		err := cdc.UnmarshalJSON([]byte(res), &parentToken)
		require.NoError(f.T, err)
		require.Equal(t, nft3ID, parentToken.(collectionModule.NFT).GetTokenID())

		query = map[string]map[string]interface{}{
			"get_root_or_parent_or_children": {
				"contract_id": collectionContractId,
				"token_id":    nft2ID,
				"target":      "root",
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)

		var rootToken collectionModule.Token
		err = cdc.UnmarshalJSON([]byte(res), &rootToken)
		require.NoError(f.T, err)
		require.Equal(t, nft3ID, rootToken.(collectionModule.NFT).GetTokenID())

		query = map[string]map[string]interface{}{
			"get_root_or_parent_or_children": {
				"contract_id": collectionContractId,
				"token_id":    nft3ID,
				"target":      "children",
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)

		var childrenTokens []collectionModule.Token
		err = cdc.UnmarshalJSON([]byte(res), &childrenTokens)
		require.NoError(f.T, err)
		require.Equal(t, 1, len(childrenTokens))
		require.Equal(t, nft2ID, childrenTokens[0].(collectionModule.NFT).GetTokenID())
	}

	// detach token
	{
		msg := map[string]map[string]interface{}{
			"detach": {
				"contract_id": collectionContractId,
				"from":        contractAddress,
				"token_id":    nft2ID,
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate attached token
	{
		parentToken := f.QueryParentTokenCollection(collectionContractId, nft2ID)
		require.Nil(t, parentToken)

		childrenTokens := f.QueryChildrenTokenCollection(collectionContractId, nft3ID)
		require.Equal(t, 0, len(childrenTokens))

		rootToken := f.QueryRootTokenCollection(collectionContractId, nft2ID)
		require.Equal(t, nft2ID, rootToken.(collectionModule.NFT).GetTokenID())
	}

	// test for query
	{
		cdc := app.MakeCodec()

		// query collection
		query := map[string]map[string]interface{}{
			"get_collection": {
				"contract_id": collectionContractId,
			},
		}
		queryJson, _ := json.Marshal(query)
		queryString := string(queryJson)
		res := f.QueryContractStateSmartWasm(contractAddress, queryString)
		var collection collectionModule.Collection
		err := cdc.UnmarshalJSON([]byte(res), &collection)
		require.NoError(f.T, err)
		require.Equal(t, collectionContractId, collection.GetContractID())

		// query balance
		query = map[string]map[string]interface{}{
			"get_balance": {
				"contract_id": collectionContractId,
				"token_id":    tokenID,
				"addr":        contractAddress,
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var balanceOfContract sdk.Int
		err = cdc.UnmarshalJSON([]byte(res), &balanceOfContract)
		require.NoError(f.T, err)
		require.Equal(t, sdk.NewInt(int64(initFtAmount+mintFtAmount-burnFtAmount-transferFtAmount)), balanceOfContract)

		// query token type
		query = map[string]map[string]interface{}{
			"get_token_type": {
				"contract_id": collectionContractId,
				"token_id":    tokenTypeID1,
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)

		var tokenType collectionModule.TokenType
		err = cdc.UnmarshalJSON([]byte(res), &tokenType)
		require.NoError(f.T, err)
		require.Equal(t, nftName1, tokenType.GetName())
		require.Equal(t, nftMeta1, tokenType.GetMeta())
		require.Equal(t, collectionContractId, tokenType.GetContractID())
		require.Equal(t, tokenTypeID1, tokenType.GetTokenType())

		// query token types
		query = map[string]map[string]interface{}{
			"get_token_types": {
				"contract_id": collectionContractId,
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)

		var tokenTypes []collectionModule.TokenType
		err = cdc.UnmarshalJSON([]byte(res), &tokenTypes)
		require.NoError(f.T, err)
		require.Equal(t, 3, len(tokenTypes))

		// query token
		query = map[string]map[string]interface{}{
			"get_token": {
				"contract_id": collectionContractId,
				"token_id":    tokenID,
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var token collectionModule.FT
		err = cdc.UnmarshalJSON([]byte(res), &token)
		require.NoError(f.T, err)

		require.Equal(t, collectionContractId, token.GetContractID())
		require.Equal(t, ftName, token.GetName())
		require.Equal(t, ftMeta, token.GetMeta())
		require.Equal(t, tokenID, token.GetTokenID())
		require.Equal(t, sdk.NewInt(8), token.GetDecimals())
		require.True(t, token.GetMintable())

		// query tokens
		query = map[string]map[string]interface{}{
			"get_tokens": {
				"contract_id": collectionContractId,
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var tokens []collectionModule.Token
		err = cdc.UnmarshalJSON([]byte(res), &tokens)
		require.NoError(f.T, err)
		require.Equal(t, 4, len(tokens))

		// query nft total
		query = map[string]map[string]interface{}{
			"get_nft_count": {
				"contract_id": collectionContractId,
				"token_id":    tokenTypeID3,
				"target":      "count",
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var nftCount sdk.Int
		err = cdc.UnmarshalJSON([]byte(res), &nftCount)
		require.NoError(f.T, err)
		require.Equal(t, sdk.NewInt(1), nftCount)

		query = map[string]map[string]interface{}{
			"get_nft_count": {
				"contract_id": collectionContractId,
				"token_id":    tokenTypeID1,
				"target":      "mint",
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var nftMint sdk.Int
		err = cdc.UnmarshalJSON([]byte(res), &nftMint)
		require.NoError(f.T, err)
		require.Equal(t, sdk.NewInt(2), nftMint)

		query = map[string]map[string]interface{}{
			"get_nft_count": {
				"contract_id": collectionContractId,
				"token_id":    tokenTypeID1,
				"target":      "burn",
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var nftBurn sdk.Int
		err = cdc.UnmarshalJSON([]byte(res), &nftBurn)
		require.NoError(f.T, err)
		require.Equal(t, sdk.NewInt(1), nftBurn)

		// query total
		query = map[string]map[string]interface{}{
			"get_total": {
				"contract_id": collectionContractId,
				"token_id":    tokenID,
				"target":      "supply",
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var totalSupply sdk.Int
		err = cdc.UnmarshalJSON([]byte(res), &totalSupply)
		require.NoError(f.T, err)
		require.Equal(t, sdk.NewInt(int64(initFtAmount+2*mintFtAmount-burnFtAmount)), totalSupply)

		query = map[string]map[string]interface{}{
			"get_total": {
				"contract_id": collectionContractId,
				"token_id":    tokenID,
				"target":      "mint",
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var totalMint sdk.Int
		err = cdc.UnmarshalJSON([]byte(res), &totalMint)
		require.NoError(f.T, err)
		require.Equal(t, sdk.NewInt(int64(initFtAmount+2*mintFtAmount)), totalMint)

		query = map[string]map[string]interface{}{
			"get_total": {
				"contract_id": collectionContractId,
				"token_id":    tokenID,
				"target":      "burn",
			},
		}
		queryJson, _ = json.Marshal(query)
		queryString = string(queryJson)
		res = f.QueryContractStateSmartWasm(contractAddress, queryString)
		var totalBurn sdk.Int
		err = cdc.UnmarshalJSON([]byte(res), &totalBurn)
		require.NoError(f.T, err)
		require.Equal(t, sdk.NewInt(int64(burnFtAmount)), totalBurn)

	}
}

func TestLinkCLIWasmCollectionTesterProxy(t *testing.T) {
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

	nftName1 := "TestNFT1"
	nftMeta1 := "nftMeta1"
	nftName2 := "TestNFT2"
	nftMeta2 := "nftMeta2"
	nftName3 := "TestNFT3"
	nftMeta3 := "nftMeta3"

	tokenTypeID1 := "10000001"
	tokenTypeID2 := "10000002"
	tokenTypeID3 := "10000003"

	index1 := "00000001"
	index2 := "00000002"
	nft0ID := tokenTypeID1 + index1
	nft1ID := tokenTypeID1 + index2
	nft2ID := tokenTypeID2 + index1
	nft3ID := tokenTypeID3 + index1

	ftName := "TestFT1"
	ftMeta := "ftMeta"
	tokenID := "0000000100000000"

	burnFtAmount := 3
	transferFtAmount := 10
	burnFt := strconv.Itoa(burnFtAmount) + ":" + tokenID
	transferFt := strconv.Itoa(transferFtAmount) + ":" + tokenID

	// store the contract collection-tester
	{
		f.LogResult(f.TxStoreWasm(wasmCollectionTester, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// instantiate
	{
		msgJson := "{}"
		flagLabel := "--label=collection-tester"
		f.LogResult(f.TxInstantiateWasm(codeId, msgJson, flagLabel, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate there is only one contract using codeId=1 and get contractAddress
	{
		listContract := f.QueryListContractByCodeWasm(codeId)
		require.Len(t, listContract, 1)
		contractAddress = listContract[0].GetAddress()
	}

	// set test token and approve for contractAddress
	{
		// create collection
		f.LogResult(f.TxTokenCreateCollection(fooAddr.String(), collectionName, collectionMeta, collectionBaseImageURI, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// issue ft
		f.LogResult(f.TxTokenIssueFTCollection(fooAddr.String(), collectionContractId, fooAddr, ftName, ftMeta, 10000, 6, true, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// issue nft
		f.LogResult(f.TxTokenIssueNFTCollection(fooAddr.String(), collectionContractId, nftName1, nftMeta1, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// mint nft
		mintParam := strings.Join([]string{tokenTypeID1, "description", "meta"}, ":")
		f.LogResult(f.TxTokenMintNFTCollection(fooAddr.String(), collectionContractId, fooAddr.String(), mintParam, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// approve and grant burn perm
	{
		nft := f.QueryTokenCollection(collectionContractId, nft0ID).(collectionModule.NFT)
		require.Equal(t, collectionContractId, nft.GetContractID())
		require.Equal(t, nft0ID, nft.GetTokenID())
		require.Equal(t, fooAddr, nft.GetOwner())

		f.LogResult(f.TxCollectionApprove(fooAddr.String(), collectionContractId, contractAddress, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		f.LogResult(f.TxCollectionGrantPerm(fooAddr.String(), contractAddress, collectionContractId, "burn", "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		res := f.QueryApprovedTokenCollection(collectionContractId, contractAddress, fooAddr)
		require.True(t, res)
	}

	// query isApproved
	{
		cdc := app.MakeCodec()

		query := map[string]map[string]interface{}{
			"get_approved": {
				"proxy":       contractAddress,
				"contract_id": collectionContractId,
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
				"contract_id": collectionContractId,
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

	// burn ft from proxy
	{
		msg := map[string]map[string]interface{}{
			"burn_ft_from": {
				"proxy":       contractAddress,
				"from":        fooAddr,
				"contract_id": collectionContractId,
				"amounts":     []string{burnFt},
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate balance for fooAddr
	{
		balanceOfFoo := f.QueryBalanceCollection(collectionContractId, tokenID, fooAddr)
		require.Equal(t, sdk.NewInt(9997), balanceOfFoo)
	}

	// burn nft from proxy
	{
		msg := map[string]map[string]interface{}{
			"burn_nft_from": {
				"proxy":       contractAddress,
				"from":        fooAddr,
				"contract_id": collectionContractId,
				"token_ids":   []string{nft0ID},
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate burn nft
	{
		f.QueryTokenCollectionExpectEmpty(collectionContractId, nft0ID)
	}

	// transfer ft from
	{
		msg := map[string]map[string]interface{}{
			"transfer_ft_from": {
				"proxy":       contractAddress,
				"contract_id": collectionContractId,
				"from":        fooAddr,
				"to":          contractAddress,
				"tokens":      []string{transferFt},
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate balance for contractAddress
	{
		balanceOfContract := f.QueryBalanceCollection(collectionContractId, tokenID, contractAddress)
		require.Equal(t, sdk.NewInt(int64(transferFtAmount)), balanceOfContract)
	}

	// transfer nft from
	{
		// prepare nft
		mintParam := strings.Join([]string{tokenTypeID1, "description", "meta"}, ":")
		f.LogResult(f.TxTokenMintNFTCollection(fooAddr.String(), collectionContractId, fooAddr.String(), mintParam, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		msg := map[string]map[string]interface{}{
			"transfer_nft_from": {
				"proxy":       contractAddress,
				"contract_id": collectionContractId,
				"from":        fooAddr,
				"to":          contractAddress,
				"token_ids":   []string{nft1ID},
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate transfered nft
	{
		nft := f.QueryTokenCollection(collectionContractId, nft1ID).(collectionModule.NFT)
		require.Equal(t, collectionContractId, nft.GetContractID())
		require.Equal(t, nft1ID, nft.GetTokenID())
		require.Equal(t, contractAddress, nft.GetOwner())
	}

	// token attach, detach from
	{
		// issue nft
		f.LogResult(f.TxTokenIssueNFTCollection(fooAddr.String(), collectionContractId, nftName2, nftMeta2, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		f.LogResult(f.TxTokenIssueNFTCollection(fooAddr.String(), collectionContractId, nftName3, nftMeta3, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// mint nft
		mintParam := strings.Join([]string{tokenTypeID2, "description", "meta"}, ":")
		f.LogResult(f.TxTokenMintNFTCollection(fooAddr.String(), collectionContractId, fooAddr.String(), mintParam, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		mintParam = strings.Join([]string{tokenTypeID3, "description", "meta"}, ":")
		f.LogResult(f.TxTokenMintNFTCollection(fooAddr.String(), collectionContractId, fooAddr.String(), mintParam, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)

		// attach from
		msg := map[string]map[string]interface{}{
			"attach_from": {
				"proxy":       contractAddress,
				"contract_id": collectionContractId,
				"from":        fooAddr,
				"to_token_id": nft3ID,
				"token_id":    nft2ID,
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate attach from
	{
		parent := f.QueryParentTokenCollection(collectionContractId, nft2ID)
		require.Equal(t, nft3ID, parent.GetTokenID())

		children := f.QueryChildrenTokenCollection(collectionContractId, nft3ID)
		require.Equal(t, 1, len(children))
		require.Equal(t, nft2ID, children[0].GetTokenID())

		root := f.QueryRootTokenCollection(collectionContractId, nft2ID)
		require.Equal(t, nft3ID, root.GetTokenID())
	}

	// detach from
	{
		msg := map[string]map[string]interface{}{
			"detach_from": {
				"proxy":       contractAddress,
				"contract_id": collectionContractId,
				"from":        fooAddr,
				"token_id":    nft2ID,
			},
		}
		msgJson, _ := json.Marshal(msg)
		msgString := string(msgJson)
		f.LogResult(f.TxExecuteWasm(contractAddress, msgString, flagFromFoo, flagGas, "-y"))
		tests.WaitForNextNBlocksTM(1, f.Port)
	}

	// validate detach from
	{
		parentToken := f.QueryParentTokenCollection(collectionContractId, nft2ID)
		require.Nil(t, parentToken)

		childrenTokens := f.QueryChildrenTokenCollection(collectionContractId, nft3ID)
		require.Equal(t, 0, len(childrenTokens))

		rootToken := f.QueryRootTokenCollection(collectionContractId, nft2ID)
		require.Equal(t, nft2ID, rootToken.(collectionModule.NFT).GetTokenID())
	}
}
