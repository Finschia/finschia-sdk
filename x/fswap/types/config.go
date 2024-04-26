package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
)

// Config is a config struct used for intialising the fswap module to avoid using globals.
type Config struct {
	// OldCoinDenom defines the old coin denom.
	OldCoinDenom string
	// NewCoinDenom defines the new coin denom.
	NewCoinDenom string
	// SwapRate defines the swap rate.
	SwapRate sdk.Int
}

// DefaultConfig returns the default config for fswap.
func DefaultConfig() Config {
	return Config{
		OldCoinDenom: "cony",
		NewCoinDenom: "PDT",
		SwapRate:     sdk.NewInt(148079656000000),
	}
}
