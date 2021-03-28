package account

import (
	"github.com/line/lbm-sdk/v2/x/account/client/cli"
	"github.com/line/lbm-sdk/v2/x/account/internal/types"
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
