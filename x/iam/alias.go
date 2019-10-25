package iam

import (
	"github.com/link-chain/link/x/iam/internal/keeper"
	"github.com/link-chain/link/x/iam/internal/types"
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
