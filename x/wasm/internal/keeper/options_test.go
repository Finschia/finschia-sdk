package keeper

import (
	"github.com/line/lbm-sdk/v2/x/wasm/internal/keeper/wasmtesting"
	"github.com/line/lbm-sdk/v2/x/wasm/internal/types"
	authkeeper "github.com/line/lbm-sdk/v2/x/auth/keeper"
	distributionkeeper "github.com/line/lbm-sdk/v2/x/distribution/keeper"
	paramtypes "github.com/line/lbm-sdk/v2/x/params/types"
	stakingkeeper "github.com/line/lbm-sdk/v2/x/staking/keeper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstructorOptions(t *testing.T) {
	specs := map[string]struct {
		srcOpt Option
		verify func(Keeper)
	}{
		"wasm engine": {
			srcOpt: WithWasmEngine(&wasmtesting.MockWasmer{}),
			verify: func(k Keeper) {
				assert.IsType(t, k.wasmer, &wasmtesting.MockWasmer{})
			},
		},
		"message handler": {
			srcOpt: WithMessageHandler(&wasmtesting.MockMessageHandler{}),
			verify: func(k Keeper) {
				assert.IsType(t, k.messenger, &wasmtesting.MockMessageHandler{})
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
			k := NewKeeper(
				nil,
				nil,
				paramtypes.NewSubspace(nil, nil, nil, nil, ""),
				authkeeper.AccountKeeper{},
				nil,
				stakingkeeper.Keeper{},
				distributionkeeper.Keeper{},
				nil,
				nil,
				nil,
				nil,
				nil,
				"tempDir",
				types.DefaultWasmConfig(),
				SupportedFeatures,
				nil,
				nil,
				spec.srcOpt,
			)
			spec.verify(k)
		})
	}

}
