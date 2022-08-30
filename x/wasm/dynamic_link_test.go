package wasm

import (
	"fmt"
	"testing"

	"github.com/line/lbm-sdk/x/wasm/keeper"
	abci "github.com/line/ostracon/abci/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// These come from https://github.com/line/cosmwasm/tree/main/contracts.
	// Hashes of them are in testdata directory.
	calleeContract     = mustLoad("./testdata/dynamic_callee_contract.wasm")
	callerContract     = mustLoad("./testdata/dynamic_caller_contract.wasm")
	numberContract     = mustLoad("./testdata/number.wasm")
	callNumberContract = mustLoad("./testdata/call_number.wasm")
)

// This tests dynamic calls using callee_contract's pong
func TestDynamicPingPongWorks(t *testing.T) {
	// setup
	data := setupTest(t)

	h := data.module.Route().Handler()

	// store dynamic callee code
	storeCalleeMsg := &MsgStoreCode{
		Sender:       addr1,
		WASMByteCode: calleeContract,
	}
	res, err := h(data.ctx, storeCalleeMsg)
	require.NoError(t, err)

	calleeCodeId := uint64(1)
	assertStoreCodeResponse(t, res.Data, calleeCodeId)

	// store dynamic caller code
	storeCallerMsg := &MsgStoreCode{
		Sender:       addr1,
		WASMByteCode: callerContract,
	}
	res, err = h(data.ctx, storeCallerMsg)
	require.NoError(t, err)

	callerCodeId := uint64(2)
	assertStoreCodeResponse(t, res.Data, callerCodeId)

	// instantiate callee contract
	instantiateCalleeMsg := &MsgInstantiateContract{
		Sender:  addr1,
		CodeID:  calleeCodeId,
		Label:   "callee",
		InitMsg: []byte(`{}`),
		Funds:   nil,
	}
	res, err = h(data.ctx, instantiateCalleeMsg)
	require.NoError(t, err)

	calleeContractAddress := parseInitResponse(t, res.Data)

	// instantiate caller contract
	cosmwasmInstantiateCallerMsg := fmt.Sprintf(`{"callee_addr":"%s"}`, calleeContractAddress)
	instantiateCallerMsg := &MsgInstantiateContract{
		Sender:  addr1,
		CodeID:  callerCodeId,
		Label:   "caller",
		InitMsg: []byte(cosmwasmInstantiateCallerMsg),
		Funds:   nil,
	}
	res, err = h(data.ctx, instantiateCallerMsg)
	require.NoError(t, err)

	callerContractAddress := parseInitResponse(t, res.Data)

	// execute ping
	cosmwasmExecuteMsg := `{"ping":{"ping_num":"100"}}`
	executeMsg := MsgExecuteContract{
		Sender:   addr1,
		Contract: callerContractAddress,
		Msg:      []byte(cosmwasmExecuteMsg),
		Funds:    nil,
	}
	res, err = h(data.ctx, &executeMsg)
	require.NoError(t, err)

	assert.Equal(t, len(res.Events), 3)
	assert.Equal(t, "wasm", res.Events[0].Type)
	assert.Equal(t, len(res.Events[0].Attributes), 6)
	assertAttribute(t, "returned_pong", "101", res.Events[0].Attributes[1])
	assertAttribute(t, "returned_pong_with_struct", "hello world 101", res.Events[0].Attributes[2])
	assertAttribute(t, "returned_pong_with_tuple", "(hello world, 42)", res.Events[0].Attributes[3])
	assertAttribute(t, "returned_pong_with_tuple_takes_2_args", "(hello world, 42)", res.Events[0].Attributes[4])
}

// This tests re-entrancy in dynamic call fails
func TestDynamicReEntrancyFails(t *testing.T) {
	// setup
	data := setupTest(t)

	h := data.module.Route().Handler()

	// store dynamic callee code
	storeCalleeMsg := &MsgStoreCode{
		Sender:       addr1,
		WASMByteCode: calleeContract,
	}
	res, err := h(data.ctx, storeCalleeMsg)
	require.NoError(t, err)

	calleeCodeId := uint64(1)
	assertStoreCodeResponse(t, res.Data, calleeCodeId)

	// store dynamic caller code
	storeCallerMsg := &MsgStoreCode{
		Sender:       addr1,
		WASMByteCode: callerContract,
	}
	res, err = h(data.ctx, storeCallerMsg)
	require.NoError(t, err)

	callerCodeId := uint64(2)
	assertStoreCodeResponse(t, res.Data, callerCodeId)

	// instantiate callee contract
	instantiateCalleeMsg := &MsgInstantiateContract{
		Sender:  addr1,
		CodeID:  calleeCodeId,
		Label:   "callee",
		InitMsg: []byte(`{}`),
		Funds:   nil,
	}
	res, err = h(data.ctx, instantiateCalleeMsg)
	require.NoError(t, err)

	calleeContractAddress := parseInitResponse(t, res.Data)

	// instantiate caller contract
	cosmwasmInstantiateCallerMsg := fmt.Sprintf(`{"callee_addr":"%s"}`, calleeContractAddress)
	instantiateCallerMsg := &MsgInstantiateContract{
		Sender:  addr1,
		CodeID:  callerCodeId,
		Label:   "caller",
		InitMsg: []byte(cosmwasmInstantiateCallerMsg),
		Funds:   nil,
	}
	res, err = h(data.ctx, instantiateCallerMsg)
	require.NoError(t, err)

	callerContractAddress := parseInitResponse(t, res.Data)

	// execute ping
	cosmwasmExecuteMsg := `{"try_re_entrancy":{}}`
	executeMsg := MsgExecuteContract{
		Sender:   addr1,
		Contract: callerContractAddress,
		Msg:      []byte(cosmwasmExecuteMsg),
		Funds:    nil,
	}
	res, err = h(data.ctx, &executeMsg)
	assert.ErrorContains(t, err, "A contract can only be called once per one call stack.")
}

// This tests both of dynamic calls and traditional queries can be used
// in a contract call
func TestDynamicCallAndTraditionalQueryWork(t *testing.T) {
	// setup
	data := setupTest(t)

	h := data.module.Route().Handler()
	q := data.module.LegacyQuerierHandler(nil)

	// store callee code (number)
	storeCalleeMsg := &MsgStoreCode{
		Sender:       addr1,
		WASMByteCode: numberContract,
	}
	res, err := h(data.ctx, storeCalleeMsg)
	require.NoError(t, err)

	calleeCodeId := uint64(1)
	assertStoreCodeResponse(t, res.Data, calleeCodeId)

	// store caller code (call-number)
	storeCallerMsg := &MsgStoreCode{
		Sender:       addr1,
		WASMByteCode: callNumberContract,
	}
	res, err = h(data.ctx, storeCallerMsg)
	require.NoError(t, err)

	callerCodeId := uint64(2)
	assertStoreCodeResponse(t, res.Data, callerCodeId)

	// instantiate callee contract
	instantiateCalleeMsg := &MsgInstantiateContract{
		Sender:  addr1,
		CodeID:  calleeCodeId,
		Label:   "number",
		InitMsg: []byte(`{"value":21}`),
		Funds:   nil,
	}
	res, err = h(data.ctx, instantiateCalleeMsg)
	require.NoError(t, err)

	calleeContractAddress := parseInitResponse(t, res.Data)

	// instantiate caller contract
	cosmwasmInstantiateCallerMsg := fmt.Sprintf(`{"callee_addr":"%s"}`, calleeContractAddress)
	instantiateCallerMsg := &MsgInstantiateContract{
		Sender:  addr1,
		CodeID:  callerCodeId,
		Label:   "call-number",
		InitMsg: []byte(cosmwasmInstantiateCallerMsg),
		Funds:   nil,
	}
	res, err = h(data.ctx, instantiateCallerMsg)
	require.NoError(t, err)

	callerContractAddress := parseInitResponse(t, res.Data)

	// traditional queries from caller
	queryPath := []string{
		QueryGetContractState,
		callerContractAddress,
		keeper.QueryMethodContractStateSmart,
	}
	queryReq := abci.RequestQuery{Data: []byte(`{"number":{}}`)}
	qRes, qErr := q(data.ctx, queryPath, queryReq)
	assert.NoError(t, qErr)
	assert.Equal(t, []byte(`{"value":21}`), qRes)

	// query via dynamic call from caller
	dynQueryReq := abci.RequestQuery{Data: []byte(`{"number_dyn":{}}`)}
	qRes, qErr = q(data.ctx, queryPath, dynQueryReq)
	assert.NoError(t, qErr)
	assert.Equal(t, []byte(`{"value":21}`), qRes)

	// execute mul
	cosmwasmExecuteMsg := `{"mul":{"value":2}}`
	executeMsg := MsgExecuteContract{
		Sender:   addr1,
		Contract: callerContractAddress,
		Msg:      []byte(cosmwasmExecuteMsg),
		Funds:    nil,
	}
	res, err = h(data.ctx, &executeMsg)
	require.NoError(t, err)
	assert.Equal(t, len(res.Events), 3)
	assert.Equal(t, "wasm", res.Events[0].Type)
	assert.Equal(t, len(res.Events[0].Attributes), 3)
	assertAttribute(t, "value_by_dynamic", "42", res.Events[0].Attributes[1])
	assertAttribute(t, "value_by_query", "42", res.Events[0].Attributes[2])

	// queries
	qRes, qErr = q(data.ctx, queryPath, queryReq)
	assert.NoError(t, qErr)
	assert.Equal(t, []byte(`{"value":42}`), qRes)
	qRes, qErr = q(data.ctx, queryPath, dynQueryReq)
	assert.NoError(t, qErr)
	assert.Equal(t, []byte(`{"value":42}`), qRes)
}

// This tests dynamic call with writing something to storage fails
// if it is called by a query
func TestDynamicCallWithWriteFailsByQuery(t *testing.T) {
	// setup
	data := setupTest(t)

	h := data.module.Route().Handler()
	q := data.module.LegacyQuerierHandler(nil)

	// store callee code (number)
	storeCalleeMsg := &MsgStoreCode{
		Sender:       addr1,
		WASMByteCode: numberContract,
	}
	res, err := h(data.ctx, storeCalleeMsg)
	require.NoError(t, err)

	calleeCodeId := uint64(1)
	assertStoreCodeResponse(t, res.Data, calleeCodeId)

	// store caller code (call-number)
	storeCallerMsg := &MsgStoreCode{
		Sender:       addr1,
		WASMByteCode: callNumberContract,
	}
	res, err = h(data.ctx, storeCallerMsg)
	require.NoError(t, err)

	callerCodeId := uint64(2)
	assertStoreCodeResponse(t, res.Data, callerCodeId)

	// instantiate callee contract
	instantiateCalleeMsg := &MsgInstantiateContract{
		Sender:  addr1,
		CodeID:  calleeCodeId,
		Label:   "number",
		InitMsg: []byte(`{"value":21}`),
		Funds:   nil,
	}
	res, err = h(data.ctx, instantiateCalleeMsg)
	require.NoError(t, err)

	calleeContractAddress := parseInitResponse(t, res.Data)

	// instantiate caller contract
	cosmwasmInstantiateCallerMsg := fmt.Sprintf(`{"callee_addr":"%s"}`, calleeContractAddress)
	instantiateCallerMsg := &MsgInstantiateContract{
		Sender:  addr1,
		CodeID:  callerCodeId,
		Label:   "call-number",
		InitMsg: []byte(cosmwasmInstantiateCallerMsg),
		Funds:   nil,
	}
	res, err = h(data.ctx, instantiateCallerMsg)
	require.NoError(t, err)

	callerContractAddress := parseInitResponse(t, res.Data)

	// query which tries to write value to storage
	queryPath := []string{
		QueryGetContractState,
		callerContractAddress,
		keeper.QueryMethodContractStateSmart,
	}
	queryReq := abci.RequestQuery{Data: []byte(`{"mul":{"value":2}}`)}
	_, qErr := q(data.ctx, queryPath, queryReq)
	assert.ErrorContains(t, qErr, "Must not call a writing storage function in this context.")
}

// This tests callee_panic in dynamic call fails
func TestDynamicCallCalleeFails(t *testing.T) {
	// setup
	data := setupTest(t)

	h := data.module.Route().Handler()

	// store dynamic callee code
	storeCalleeMsg := &MsgStoreCode{
		Sender:       addr1,
		WASMByteCode: calleeContract,
	}
	res, err := h(data.ctx, storeCalleeMsg)
	require.NoError(t, err)

	calleeCodeId := uint64(1)
	assertStoreCodeResponse(t, res.Data, calleeCodeId)

	// store dynamic caller code
	storeCallerMsg := &MsgStoreCode{
		Sender:       addr1,
		WASMByteCode: callerContract,
	}
	res, err = h(data.ctx, storeCallerMsg)
	require.NoError(t, err)

	callerCodeId := uint64(2)
	assertStoreCodeResponse(t, res.Data, callerCodeId)

	// instantiate callee contract
	instantiateCalleeMsg := &MsgInstantiateContract{
		Sender:  addr1,
		CodeID:  calleeCodeId,
		Label:   "callee",
		InitMsg: []byte(`{}`),
		Funds:   nil,
	}
	res, err = h(data.ctx, instantiateCalleeMsg)
	require.NoError(t, err)

	calleeContractAddress := parseInitResponse(t, res.Data)

	// instantiate caller contract
	cosmwasmInstantiateCallerMsg := fmt.Sprintf(`{"callee_addr":"%s"}`, calleeContractAddress)
	instantiateCallerMsg := &MsgInstantiateContract{
		Sender:  addr1,
		CodeID:  callerCodeId,
		Label:   "caller",
		InitMsg: []byte(cosmwasmInstantiateCallerMsg),
		Funds:   nil,
	}
	res, err = h(data.ctx, instantiateCallerMsg)
	require.NoError(t, err)

	callerContractAddress := parseInitResponse(t, res.Data)

	// execute pong_panic
	cosmwasmExecuteMsg := `{"callee_panic":{}}`
	executeMsg := MsgExecuteContract{
		Sender:   addr1,
		Contract: callerContractAddress,
		Msg:      []byte(cosmwasmExecuteMsg),
		Funds:    nil,
	}
	res, err = h(data.ctx, &executeMsg)
	assert.ErrorContains(t, err, "Error in dynamic link")
}
