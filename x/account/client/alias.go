package client

import (
	"github.com/line/lbm-sdk/x/account/client/cli"
	"github.com/line/lbm-sdk/x/account/client/rest"
	cosmoscli "github.com/line/lbm-sdk/x/auth/client/cli"
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
