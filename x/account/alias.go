package account

import (
	"github.com/line/link-modules/x/account/client/cli"
	"github.com/line/link-modules/x/account/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
)

var (
	CreateAccountTxCmd = cli.CreateAccountCmd
	EmptyTxCmd         = cli.EmptyCmd
	NewMsgEmpty        = types.NewMsgEmpty
	RegisterCodec      = types.RegisterCodec
	ModuleCdc          = types.ModuleCdc
)
