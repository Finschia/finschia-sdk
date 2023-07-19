package types

import paramtypes "github.com/Finschia/finschia-sdk/x/params/types"

// NewParams creates a new parameter configuration for the bank module
func NewParams() Params {
	return Params{}
}

// DefaultParams is the default parameter configuration for the bank module
func DefaultParams() Params {
	return Params{}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}
