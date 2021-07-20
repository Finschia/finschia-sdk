package wasmtesting

import (
	sdk "github.com/line/lfb-sdk/types"
	wasmvmtypes "github.com/line/wasmvm/types"
)

type MockMsgDispatcher struct {
	DispatchMessagesFn    func(ctx sdk.Context, contractAddr sdk.AccAddress, ibcPort string, msgs []wasmvmtypes.CosmosMsg) error
	DispatchSubmessagesFn func(ctx sdk.Context, contractAddr sdk.AccAddress, ibcPort string, msgs []wasmvmtypes.SubMsg) ([]byte, error)
}

func (m MockMsgDispatcher) DispatchMessages(ctx sdk.Context, contractAddr sdk.AccAddress, ibcPort string, msgs []wasmvmtypes.CosmosMsg) error {
	if m.DispatchMessagesFn == nil {
		panic("not expected to be called")
	}
	return m.DispatchMessagesFn(ctx, contractAddr, ibcPort, msgs)
}

func (m MockMsgDispatcher) DispatchSubmessages(ctx sdk.Context, contractAddr sdk.AccAddress, ibcPort string, msgs []wasmvmtypes.SubMsg) ([]byte, error) {
	if m.DispatchSubmessagesFn == nil {
		panic("not expected to be called")
	}
	return m.DispatchSubmessagesFn(ctx, contractAddr, ibcPort, msgs)
}
