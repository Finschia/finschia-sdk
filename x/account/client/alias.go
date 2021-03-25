package client

import (
	cosmoscli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/line/link-modules/x/account/client/cli"
	"github.com/line/link-modules/x/account/client/rest"
)

var (
	GetAccountCmd                    = cosmoscli.GetAccountCmd
	QueryTxsByEventsCmd              = cli.QueryTxsByEventsCmd
	QueryTxCmd                       = cli.QueryTxCmd
	QueryBlockWithTxResponsesCommand = cli.QueryBlockWithTxResponsesCommand
	GetSignCommand                   = cosmoscli.GetSignCommand
	GetMultiSignCommand              = cosmoscli.GetMultiSignCommand
	GetBroadcastCommand              = cosmoscli.GetBroadcastCommand
	GetEncodeCommand                 = cosmoscli.GetEncodeCommand
	RegisterTxRoutes                 = rest.RegisterTxRoutes
)
