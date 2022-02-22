package keeper

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	types "github.com/line/lbm-sdk/x/wasm/types"
	wasmvm "github.com/line/wasmvm"
	wasmvmtypes "github.com/line/wasmvm/types"
)

type cosmwasmAPIImpl struct {
	gasMultiplier GasMultiplier
}

const (
	// DefaultDeserializationCostPerByte The formular should be `len(data) * deserializationCostPerByte`
	DefaultDeserializationCostPerByte = 1
)

var (
	costJsonDeserialization = wasmvmtypes.UFraction{
		Numerator:   DefaultDeserializationCostPerByte * types.DefaultGasMultiplier,
		Denominator: 1,
	}
)

func (a cosmwasmAPIImpl) humanAddress(canon []byte) (string, uint64, error) {
	gas := a.gasMultiplier.FromWasmVMGas(5)
	if len(canon) != sdk.BytesAddrLen {
		//nolint:stylecheck
		return "", gas, fmt.Errorf("expected %d byte address", sdk.BytesAddrLen)
	}

	return sdk.BytesToAccAddress(canon).String(), gas, nil
}

func (a cosmwasmAPIImpl) canonicalAddress(human string) ([]byte, uint64, error) {
	bz, err := sdk.AccAddressToBytes(human)
	return bz, a.gasMultiplier.ToWasmVMGas(4), err
}

func (k Keeper) cosmwasmAPI(ctx sdk.Context) wasmvm.GoAPI {
	x := cosmwasmAPIImpl{
		gasMultiplier: k.getGasMultiplier(ctx),
	}
	return wasmvm.GoAPI{
		HumanAddress:     x.humanAddress,
		CanonicalAddress: x.canonicalAddress,
	}
}
