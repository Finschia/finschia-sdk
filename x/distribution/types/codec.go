package types

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/codec/legacy"
	"github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/msgservice"
	govtypes "github.com/line/lbm-sdk/x/gov/types"
)

// RegisterLegacyAminoCodec registers the necessary x/distribution interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgWithdrawDelegatorReward{}, "cosmos-sdk/MsgWithdrawDelegationReward", nil)
	cdc.RegisterConcrete(&MsgWithdrawValidatorCommission{}, "cosmos-sdk/MsgWithdrawValidatorCommission", nil)
	cdc.RegisterConcrete(&MsgSetWithdrawAddress{}, "cosmos-sdk/MsgModifyWithdrawAddress", nil)
	cdc.RegisterConcrete(&MsgFundCommunityPool{}, "cosmos-sdk/MsgFundCommunityPool", nil)
	cdc.RegisterConcrete(&CommunityPoolSpendProposal{}, "cosmos-sdk/CommunityPoolSpendProposal", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgWithdrawDelegatorReward{},
		&MsgWithdrawValidatorCommission{},
		&MsgSetWithdrawAddress{},
		&MsgFundCommunityPool{},
	)
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&CommunityPoolSpendProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	RegisterLegacyAminoCodec(legacy.Cdc)
}
