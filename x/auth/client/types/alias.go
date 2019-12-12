package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	ABCIMessageLogs = sdk.ABCIMessageLogs
	StringEvents    = sdk.StringEvents
	Tx              = sdk.Tx
)

var (
	ParseABCILogs   = sdk.ParseABCILogs
	StringifyEvents = sdk.StringifyEvents
)
