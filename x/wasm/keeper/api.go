package keeper

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/wasm/types"
	wasmvm "github.com/line/wasmvm"
)

type cosmwasmAPIImpl struct {
	gasMultiplier uint64
	keeper        *Keeper
	ctx           *sdk.Context
}

func (a cosmwasmAPIImpl) humanAddress(canon []byte) (string, uint64, error) {
	gas := 5 * a.gasMultiplier
	if len(canon) != sdk.BytesAddrLen {
		//nolint:stylecheck
		return "", gas, fmt.Errorf("expected %d byte address", sdk.BytesAddrLen)
	}

	return sdk.BytesToAccAddress(canon).String(), gas, nil
}

func (a cosmwasmAPIImpl) canonicalAddress(human string) ([]byte, uint64, error) {
	bz, err := sdk.AccAddressToBytes(human)
	return bz, 4 * a.gasMultiplier, err
}

func (a cosmwasmAPIImpl) GetContractEnv(contractAddrStr string) (wasmvm.Env, *wasmvm.Cache, wasmvm.KVStore, wasmvm.Querier, wasmvm.GasMeter, []byte, uint64, error) {
	contractAddr := sdk.AccAddress(contractAddrStr)
	_, codeInfo, prefixStore, err := a.keeper.contractInstance(*a.ctx, contractAddr)
	if err != nil {
		return wasmvm.Env{}, nil, nil, nil, nil, wasmvm.Checksum{}, 0, err
	}

	cache := a.keeper.wasmVM.GetCache()
	if cache == nil {
		panic("cannot found instance cache")
	}

	// prepare querier
	querier := NewQueryHandler(*a.ctx, a.keeper.wasmVMQueryHandler, contractAddr, a.gasMultiplier)

	// this gas cost is temporal value defined by
	// https://github.com/line/lbm-sdk/runs/8150140720?check_suite_focus=true#step:5:483
	// Before release, it is adjusted by benchmark taken in environment similar to the nodes.
	gas := 11 * a.gasMultiplier
	wasmStore := types.NewWasmStore(prefixStore)
	env := types.NewEnv(*a.ctx, contractAddr)

	return env, cache, wasmStore, querier, a.keeper.gasMeter(*a.ctx), codeInfo.CodeHash, gas, nil
}

func (k Keeper) cosmwasmAPI(ctx sdk.Context) wasmvm.GoAPI {
	x := cosmwasmAPIImpl{
		gasMultiplier: k.getGasMultiplier(ctx),
		keeper:        &k,
		ctx:           &ctx,
	}
	return wasmvm.GoAPI{
		HumanAddress:     x.humanAddress,
		CanonicalAddress: x.canonicalAddress,
		GetContractEnv:   x.GetContractEnv,
	}
}
