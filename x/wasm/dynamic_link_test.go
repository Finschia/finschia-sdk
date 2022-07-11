package wasm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	calleeContract     = mustLoad("./testdata/dynamic_callee_contract.wasm")
	callerContract     = mustLoad("./testdata/dynamic_caller_contract.wasm")
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
		Sender: addr1,
		CodeID: calleeCodeId,
		Label: "callee",
		InitMsg: []byte(`{}`),
		Funds: nil,
	}
	res, err = h(data.ctx, instantiateCalleeMsg)
	require.NoError(t, err)

	calleeContractAddress := parseInitResponse(t, res.Data)

	// instantiate caller contract
	cosmwasmInstantiateCallerMsg := fmt.Sprintf(`{"callee_addr":"%s"}`, calleeContractAddress)
	instantiateCallerMsg := &MsgInstantiateContract{
		Sender: addr1,
		CodeID: callerCodeId,
		Label: "caller",
		InitMsg: []byte(cosmwasmInstantiateCallerMsg),
		Funds: nil,
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
	assertAttribute(t, "contract_address", "link10pyejy66429refv3g35g2t7am0was7yaducgya", res.Events[0].Attributes[0])
	assertAttribute(t, "returned_pong", "101", res.Events[0].Attributes[1])
	assertAttribute(t, "returned_pong_with_struct", "hello world 101", res.Events[0].Attributes[2])
	assertAttribute(t, "returned_pong_with_tuple", "(hello world, 42)", res.Events[0].Attributes[3])
	assertAttribute(t, "returned_pong_with_tuple_takes_2_args", "(hello world, 42)", res.Events[0].Attributes[4])
	assertAttribute(t, "returned_contract_address", "link18vd8fpwxzck93qlwghaj6arh4p7c5n89fvcmzu", res.Events[0].Attributes[5])
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
		Sender: addr1,
		CodeID: calleeCodeId,
		Label: "callee",
		InitMsg: []byte(`{}`),
		Funds: nil,
	}
	res, err = h(data.ctx, instantiateCalleeMsg)
	require.NoError(t, err)

	calleeContractAddress := parseInitResponse(t, res.Data)

	// instantiate caller contract
	cosmwasmInstantiateCallerMsg := fmt.Sprintf(`{"callee_addr":"%s"}`, calleeContractAddress)
	instantiateCallerMsg := &MsgInstantiateContract{
		Sender: addr1,
		CodeID: callerCodeId,
		Label: "caller",
		InitMsg: []byte(cosmwasmInstantiateCallerMsg),
		Funds: nil,
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
