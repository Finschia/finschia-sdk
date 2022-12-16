package keeper

import (
	wasmvm "github.com/line/wasmvm"
	wasmvmtypes "github.com/line/wasmvm/types"

	sdk "github.com/line/lbm-sdk/types"
	types "github.com/line/lbm-sdk/x/wasm/types"
)

type cosmwasmAPIImpl struct {
	gasMultiplier GasMultiplier
	keeper        *Keeper
	ctx           *sdk.Context
}

const (
	// DefaultDeserializationCostPerByte The formular should be `len(data) * deserializationCostPerByte`
	DefaultDeserializationCostPerByte = 1
)

var (
	costJSONDeserialization = wasmvmtypes.UFraction{
		Numerator:   DefaultDeserializationCostPerByte * types.DefaultGasMultiplier,
		Denominator: 1,
	}
)

func (a cosmwasmAPIImpl) humanAddress(canon []byte) (string, uint64, error) {
	gas := a.gasMultiplier.FromWasmVMGas(5)
	if err := sdk.VerifyAddressFormat(canon); err != nil {
		return "", gas, err
	}

	return sdk.AccAddress(canon).String(), gas, nil
}

func (a cosmwasmAPIImpl) canonicalAddress(human string) ([]byte, uint64, error) {
	bz, err := sdk.AccAddressFromBech32(human)
	return bz, a.gasMultiplier.ToWasmVMGas(4), err
}

func (a cosmwasmAPIImpl) GetContractEnv(contractAddrStr string, inputSize uint64) (wasmvm.Env, *wasmvm.Cache, wasmvm.KVStore, wasmvm.Querier, wasmvm.GasMeter, []byte, uint64, uint64, error) {
	contractAddr := sdk.MustAccAddressFromBech32(contractAddrStr)
	contractInfo, codeInfo, prefixStore, err := a.keeper.contractInstance(*a.ctx, contractAddr)
	if err != nil {
		return wasmvm.Env{}, nil, nil, nil, nil, wasmvm.Checksum{}, 0, 0, err
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
	gas := a.gasMultiplier.ToWasmVMGas(11)
	instantiateCost := a.gasMultiplier.ToWasmVMGas(a.keeper.instantiateContractCosts(a.keeper.gasRegister, *a.ctx, a.keeper.IsPinnedCode(*a.ctx, contractInfo.CodeID), int(inputSize)))
	wasmStore := types.NewWasmStore(prefixStore)
	env := types.NewEnv(*a.ctx, contractAddr)

	return env, cache, wasmStore, querier, a.keeper.gasMeter(*a.ctx), codeInfo.CodeHash, instantiateCost, gas, nil
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
