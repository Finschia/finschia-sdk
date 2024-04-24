package config

import (
	"math/big"

	"github.com/Finschia/finschia-sdk/types"
)

type FswapConfig interface {
	OldDenom() string
	NewDenom() string
	SwapCap() types.Int
	SwapMultiple() types.Int
}

func DefaultConfig() *Config {
	oldDenom := "cony"
	newDenom := "peb"
	defaultCap, _ := big.NewInt(0).SetString("1151185567094084523856000000", 10)
	swapCap := types.NewIntFromBigInt(defaultCap)
	num, _ := big.NewInt(0).SetString("148079656000000", 10)
	swapMultiple := types.NewIntFromBigInt(num)

	return &Config{
		oldDenom,
		newDenom,
		swapCap,
		swapMultiple,
	}
}

type Config struct {
	oldDenom     string
	newDenom     string
	swapCap      types.Int
	swapMultiple types.Int
}

func (c Config) OldDenom() string {
	return c.oldDenom
}

func (c Config) NewDenom() string {
	return c.newDenom
}

func (c Config) SwapCap() types.Int {
	return c.swapCap
}

func (c Config) SwapMultiple() types.Int {
	return c.swapMultiple
}
