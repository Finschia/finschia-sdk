package coin

import (
	"github.com/line/link/x/coin/client/cli"
	"github.com/line/link/x/coin/internal/keeper"
	"github.com/line/link/x/coin/internal/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
)

type (
	Keeper = keeper.Keeper

	MsgMultiSend = types.MsgMultiSend

	Input  = types.Input
	Output = types.Output
)

var (
	SendTxCmd                      = cli.SendTxCmd
	NewMsgSend                     = types.NewMsgSend
	NewKeeper                      = keeper.NewKeeper
	ActionTransferTo               = types.ActionTransferTo
	ErrCanNotTransferToBlacklisted = types.ErrCanNotTransferToBlacklisted
)
