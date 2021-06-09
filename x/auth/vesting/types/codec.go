package types

import (
	"github.com/line/lfb-sdk/codec"
	"github.com/line/lfb-sdk/codec/types"
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/types/msgservice"
	authtypes "github.com/line/lfb-sdk/x/auth/types"
	"github.com/line/lfb-sdk/x/auth/vesting/exported"
)

// RegisterLegacyAminoCodec registers the vesting interfaces and concrete types on the
// provided LegacyAmino codec. These types are used for Amino JSON serialization
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*exported.VestingAccount)(nil), nil)
	cdc.RegisterConcrete(&BaseVestingAccount{}, "lfb-sdk/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&ContinuousVestingAccount{}, "lfb-sdk/ContinuousVestingAccount", nil)
	cdc.RegisterConcrete(&DelayedVestingAccount{}, "lfb-sdk/DelayedVestingAccount", nil)
	cdc.RegisterConcrete(&PeriodicVestingAccount{}, "lfb-sdk/PeriodicVestingAccount", nil)
}

// RegisterInterface associates protoName with AccountI and VestingAccount
// Interfaces and creates a registry of it's concrete implementations
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface(
		"lfb.vesting.v1beta1.VestingAccount",
		(*exported.VestingAccount)(nil),
		&ContinuousVestingAccount{},
		&DelayedVestingAccount{},
		&PeriodicVestingAccount{},
	)

	registry.RegisterImplementations(
		(*authtypes.AccountI)(nil),
		&BaseVestingAccount{},
		&DelayedVestingAccount{},
		&ContinuousVestingAccount{},
		&PeriodicVestingAccount{},
	)

	registry.RegisterImplementations(
		(*authtypes.GenesisAccount)(nil),
		&BaseVestingAccount{},
		&DelayedVestingAccount{},
		&ContinuousVestingAccount{},
		&PeriodicVestingAccount{},
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateVestingAccount{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var amino = codec.NewLegacyAmino()

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
