package iam

import (
	"github.com/line/link/x/iam/internal/keeper"
	"github.com/line/link/x/iam/internal/types"
)

const (
	ModuleName = types.ModuleName

	StoreKey = types.StoreKey
)

type (
	Keeper = keeper.Keeper
)

var (
	RegisterCodec = types.RegisterCodec
	NewKeeper     = keeper.NewKeeper
)
