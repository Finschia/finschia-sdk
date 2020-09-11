package genesis

import (
	"github.com/line/link-modules/x/genesis/internal/keeper"
	"github.com/line/link-modules/x/genesis/internal/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey
)

type (
	Keeper = keeper.Keeper
)

var (
	ModuleCdc = types.ModuleCdc
	NewKeeper = keeper.NewKeeper
)
