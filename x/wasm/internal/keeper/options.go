package keeper

import "github.com/line/lbm-sdk/v2/x/wasm/internal/types"

type optsFn func(*Keeper)

func (f optsFn) apply(keeper *Keeper) {
	f(keeper)
}

// WithMessageHandler is an optional constructor parameter to replace the default wasm vm engine with the
// given one.
func WithWasmEngine(x types.WasmerEngine) Option {
	return optsFn(func(k *Keeper) {
		k.wasmer = x
	})
}

// WithMessageHandler is an optional constructor parameter to set a custom message handler.
func WithMessageHandler(x messenger) Option {
	return optsFn(func(k *Keeper) {
		k.messenger = x
	})
}

// WithCoinTransferrer is an optional constructor parameter to set a custom coin transferrer
func WithCoinTransferrer(x coinTransferrer) Option {
	return optsFn(func(k *Keeper) {
		k.bank = x
	})
}
