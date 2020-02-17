package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/line/link/x/iam/exported"
)

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*AccountPermissionI)(nil), nil)
	cdc.RegisterConcrete(&AccountPermission{}, "link/AccountPermission", nil)
	cdc.RegisterInterface((*InheritedAccountPermissionI)(nil), nil)
	cdc.RegisterConcrete(&InheritedAccountPermission{}, "link/InheritedAccountPermission", nil)
	cdc.RegisterInterface((*exported.PermissionI)(nil), nil)
}
