package keeper

import (
	bankpluskeeper "github.com/line/lbm-sdk/x/bankplus/keeper"
	"testing"

	authkeeper "github.com/line/lbm-sdk/x/auth/keeper"
	distributionkeeper "github.com/line/lbm-sdk/x/distribution/keeper"
	paramtypes "github.com/line/lbm-sdk/x/params/types"
	stakingkeeper "github.com/line/lbm-sdk/x/staking/keeper"
	"github.com/line/lbm-sdk/x/wasm/keeper/wasmtesting"
	"github.com/line/lbm-sdk/x/wasm/types"
	"github.com/stretchr/testify/assert"
)

func TestConstructorOptions(t *testing.T) {
	specs := map[string]struct {
		srcOpt Option
		verify func(Keeper)
	}{
		"wasm engine": {
			srcOpt: WithWasmEngine(&wasmtesting.MockWasmer{}),
			verify: func(k Keeper) {
				assert.IsType(t, k.wasmVM, &wasmtesting.MockWasmer{})
			},
		},
		"message handler": {
			srcOpt: WithMessageHandler(&wasmtesting.MockMessageHandler{}),
			verify: func(k Keeper) {
				assert.IsType(t, k.messenger, &wasmtesting.MockMessageHandler{})
			},
		},
		"query plugins": {
			srcOpt: WithQueryHandler(&wasmtesting.MockQueryHandler{}),
			verify: func(k Keeper) {
				assert.IsType(t, k.wasmVMQueryHandler, &wasmtesting.MockQueryHandler{})
			},
		},
		"coin transferrer": {
			srcOpt: WithCoinTransferrer(&wasmtesting.MockCoinTransferrer{}),
			verify: func(k Keeper) {
				assert.IsType(t, k.bank, &wasmtesting.MockCoinTransferrer{})
			},
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			k := NewKeeper(nil, nil, paramtypes.NewSubspace(nil, nil, nil, ""), authkeeper.AccountKeeper{}, bankpluskeeper.BaseKeeper{}, stakingkeeper.Keeper{}, distributionkeeper.Keeper{}, nil, nil, nil, nil, nil, nil, "tempDir", types.DefaultWasmConfig(), SupportedFeatures, nil, nil, spec.srcOpt)
			spec.verify(k)
		})
	}

}
