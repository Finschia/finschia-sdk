package bank

import (
	"github.com/line/link/x/bank/client/cli"
	"github.com/line/link/x/bank/internal/keeper"
	"github.com/line/link/x/bank/internal/types"
)

const (
	ModuleName = types.ModuleName
)

type (
	Keeper = keeper.Keeper

	MsgMultiSend = types.MsgMultiSend

	Input  = types.Input
	Output = types.Output
)

var (
	SendTxCmd = cli.SendTxCmd
	NewKeeper = keeper.NewKeeper
)
