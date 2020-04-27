package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

const (
	DefaultParamspace = ModuleName

	DefaultMaxComposableDepth uint64 = 20
	DefaultMaxComposableWidth uint64 = 20
)

var (
	KeyMaxComposableDepth = []byte("MaxComposableDepth")
	KeyMaxComposableWidth = []byte("MaxComposableWidth")
)

var _ subspace.ParamSet = &Params{}

type Params struct {
	MaxComposableDepth uint64 `json:"max_composable_depth" yaml:"max_composable_depth"`
	MaxComposableWidth uint64 `json:"max_composable_width" yaml:"max_composable_width"`
}

func NewParams(maxComposableDepth, maxComposableWidth uint64) Params {
	return Params{
		MaxComposableDepth: maxComposableDepth,
		MaxComposableWidth: maxComposableWidth,
	}
}

func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		params.NewParamSetPair(KeyMaxComposableDepth, &p.MaxComposableDepth, validateMaxComposableDepth),
		params.NewParamSetPair(KeyMaxComposableWidth, &p.MaxComposableWidth, validateMaxComposableWidth),
	}
}

func (p Params) Validate() error {
	if err := validateMaxComposableDepth(p.MaxComposableDepth); err != nil {
		return err
	}

	if err := validateMaxComposableWidth(p.MaxComposableWidth); err != nil {
		return err
	}

	return nil
}

func validateMaxComposableDepth(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid max composable depth: %d", v)
	}

	return nil
}

func validateMaxComposableWidth(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("invalid max composable width: %d", v)
	}

	return nil
}

func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() Params {
	return NewParams(DefaultMaxComposableDepth, DefaultMaxComposableWidth)
}
