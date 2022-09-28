package v040

import (
	"github.com/golang/protobuf/proto" // nolint: staticcheck

	codectypes "github.com/line/lbm-sdk/codec/types"
	"github.com/line/lbm-sdk/x/bank/types"
)

// SupplyI defines an inflationary supply interface for modules that handle
// token supply.
// It is copy-pasted from:
// https://github.com/cosmos/cosmos-sdk/blob/v042.3/x/bank/exported/exported.go
// where we stripped off the unnecessary methods.
//
// It is used in the migration script, because we save this interface as an Any
// in the supply state.
//
// Deprecated.
type SupplyI interface {
	proto.Message
}

// RegisterInterfaces registers interfaces required for the v0.40 migrations.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"cosmos.bank.v1beta1.SupplyI",
		(*SupplyI)(nil), // nolint: staticcheck
		&types.Supply{}, // nolint: staticcheck
	)
}
