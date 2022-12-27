package keeper

import (
	"fmt"
	"io/ioutil"
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	wasmtype "github.com/line/lbm-sdk/x/wasm/types"
	wasmvm "github.com/line/wasmvm"
	wasmvmapi "github.com/line/wasmvm/api"
	wasmvmtypes "github.com/line/wasmvm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newAPI(t *testing.T) wasmvm.GoAPI {
	ctx, keepers := CreateTestInput(t, false, SupportedFeatures, nil, nil)
	return keepers.WasmKeeper.cosmwasmAPI(ctx)
}

func TestAPIHumanAddress(t *testing.T) {
	// prepare API
	api := newAPI(t)

	t.Run("valid address", func(t *testing.T) {
		// address for alice in testnet
		addr := "link1twsfmuj28ndph54k4nw8crwu8h9c8mh3rtx705"
		bz, err := sdk.AccAddressFromBech32(addr)
		require.NoError(t, err)
		result, gas, err := api.HumanAddress(bz)
		require.NoError(t, err)
		assert.Equal(t, addr, result)
		assert.Equal(t, wasmtype.DefaultGasMultiplier * 5, gas)
	})

	t.Run("invalid address", func(t *testing.T) {
		_, gas, err := api.HumanAddress([]byte("invalid_address"))
		require.Error(t, err)
		assert.Equal(t, wasmtype.DefaultGasMultiplier * 5, gas)
	})
}

func TestAPICanonicalAddress(t *testing.T) {
	// prepare API
	api := newAPI(t)

	t.Run("valid address", func(t *testing.T) {
		addr := "link1twsfmuj28ndph54k4nw8crwu8h9c8mh3rtx705"
		expected, err := sdk.AccAddressFromBech32(addr)
		require.NoError(t, err)
		result, gas, err := api.CanonicalAddress(addr)
		require.NoError(t, err)
		assert.Equal(t, expected.Bytes(), result)
		assert.Equal(t, wasmtype.DefaultGasMultiplier * 4, gas)
	})

	t.Run("invalid address", func(t * testing.T) {
		_, gas, err := api.CanonicalAddress("invalid_address")
		assert.Error(t, err)
		assert.Equal(t, wasmtype.DefaultGasMultiplier * 4, gas)
	})
}

func TestAPIGetContractEnv(t *testing.T) {
	// prepare ctx and keeper
	ctx, keepers := CreateTestInput(t, false, SupportedFeatures, nil, nil)

	// instantiate a number contract
	numberWasm, err := ioutil.ReadFile("../testdata/number.wasm")
	require.NoError(t, err)
	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	creator := keepers.Faucet.NewFundedAccount(ctx, deposit...)
	em := sdk.NewEventManager()
	codeID, err := keepers.ContractKeeper.Create(ctx.WithEventManager(em), creator, numberWasm, nil)
	require.NoError(t, err)
	value := 42
	initMsg := []byte(fmt.Sprintf(`{"value":%d}`, value))
	contractAddr, _, err := keepers.ContractKeeper.Instantiate(ctx.WithEventManager(em), codeID, creator, nil, initMsg, "number", nil)
	require.NoError(t, err)
	msgLen := 101010

	// prepare API
	api := keepers.WasmKeeper.cosmwasmAPI(ctx)

	t.Run("succeed", func(t *testing.T) {
		// omitted value is MultipliedGasMeter. It is not tested here.
		env, cache, store, querier, _, hash, instantiateCost, gas, err := api.GetContractEnv(contractAddr.String(), uint64(msgLen))

		require.NoError(t, err)

		assert.Equal(t, uint64(ctx.BlockHeight()), env.Block.Height)
		assert.Equal(t, uint64(ctx.BlockTime().UnixNano()), env.Block.Time)
		assert.Equal(t, ctx.ChainID(), env.Block.ChainID)
		assert.Equal(t, contractAddr.String(), env.Contract.Address)

		code, err := wasmvmapi.GetCode(*cache, hash)
		assert.Equal(t, numberWasm, code)

		// "number" comes from https://github.com/line/cosmwasm/blob/d08b5a59115cc3d28f375b7425b902bfd1dac6a4/contracts/number/src/contract.rs#L9
		assert.Equal(t, []byte{uint8(0), uint8(0), uint8(0), uint8(value)}, store.Get([]byte("number")))

		queryMsg := []byte(`{"number":{}}`)
		query := wasmvmtypes.QueryRequest {
			Wasm: &wasmvmtypes.WasmQuery {
				Smart: &wasmvmtypes.SmartQuery {
					ContractAddr: contractAddr.String(),
					Msg: queryMsg,
				},
			},
		}
		queryResult, err := querier.Query(query, 10_000_000_000_000)
		require.NoError(t, err)
		assert.Equal(t, []byte(`{"value":42}`), queryResult)

		expectedInstantiateCost := keepers.WasmKeeper.instantiateContractCosts(keepers.WasmKeeper.gasRegister, ctx, false, msgLen)
		assert.Equal(t, wasmtype.DefaultGasMultiplier * expectedInstantiateCost, instantiateCost)

		assert.Equal(t, wasmtype.DefaultGasMultiplier * 11, gas)
	})

	t.Run("non-existed contract", func(t *testing.T){
		nonExistedContractAddr := "link1nc5tatafv6eyq7llkr2gv50ff9e22mnf70qgjlv737ktmt4eswrqyu0w3p"
		require.NotEqual(t, nonExistedContractAddr, contractAddr)
		_, _, _, _, _, _, _, _, err := api.GetContractEnv(nonExistedContractAddr, uint64(msgLen))
		require.Error(t, err)
	})
}
