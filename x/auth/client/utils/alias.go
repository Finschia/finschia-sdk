package utils

import (
	cosmosutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

var (
	// functions aliases
	GenerateOrBroadcastMsgs    = cosmosutils.GenerateOrBroadcastMsgs
	GetTxEncoder               = cosmosutils.GetTxEncoder
	WriteGenerateStdTxResponse = cosmosutils.WriteGenerateStdTxResponse
)
