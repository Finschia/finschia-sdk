package coin

import (
	"github.com/line/link-modules/x/coin/client/cli"
	"github.com/line/link-modules/x/coin/internal/keeper"
	"github.com/line/link-modules/x/coin/internal/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey
)

type (
	Keeper = keeper.Keeper

	MsgSend      = types.MsgSend
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
