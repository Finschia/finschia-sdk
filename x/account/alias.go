package account

import (
	"github.com/line/link/x/account/client/cli"
	"github.com/line/link/x/account/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
)

var (
	CreateAccountTxCmd = cli.CreateAccountCmd
	EmptyTxCmd         = cli.EmptyCmd
	RegisterCodec      = types.RegisterCodec
	ModuleCdc          = types.ModuleCdc
)
