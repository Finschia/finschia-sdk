package keeper

import "github.com/Finschia/finschia-rdk/x/gov/types"

// UnsafeSetHooks updates the gov keeper's hooks, overriding any potential
// pre-existing hooks.
// WARNING: this function should only be used in tests.
func UnsafeSetHooks(k *Keeper, h types.GovHooks) {
	k.hooks = h
}
