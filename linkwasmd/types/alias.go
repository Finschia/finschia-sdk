package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	AddrLen = sdk.AddrLen // 20

	PrefixPublic    = sdk.PrefixPublic    // "pub"
	PrefixValidator = sdk.PrefixValidator // "val"
	PrefixOperator  = sdk.PrefixOperator  // "oper"
	PrefixConsensus = sdk.PrefixConsensus // "cons"
	PrefixAddress   = sdk.PrefixAddress   // "addr"
)
