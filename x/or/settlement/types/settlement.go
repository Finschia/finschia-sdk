package types

import (
	paramtypes "github.com/Finschia/finschia-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Challenge)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Challenge{})
}

// ParamSetPairs get the params.ParamSet
func (c *Challenge) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}
