package client

import (
	cosmoscli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/line/link/x/account/client/cli"
	"github.com/line/link/x/account/client/rest"
)

var (
	GetAccountCmd       = cosmoscli.GetAccountCmd
	QueryTxsByEventsCmd = cli.QueryTxsByEventsCmd
	QueryTxCmd          = cli.QueryTxCmd
	GetSignCommand      = cosmoscli.GetSignCommand
	GetMultiSignCommand = cosmoscli.GetMultiSignCommand
	GetBroadcastCommand = cosmoscli.GetBroadcastCommand
	GetEncodeCommand    = cosmoscli.GetEncodeCommand
	RegisterTxRoutes    = rest.RegisterTxRoutes
)
