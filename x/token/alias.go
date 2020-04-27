package token

import (
	"github.com/line/link/x/token/internal/keeper"
	"github.com/line/link/x/token/internal/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey
)

type (
	Account     = types.Account
	Token       = types.Token
	Tokens      = types.Tokens
	Permissions = types.Permissions
	Keeper      = keeper.Keeper
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
	NewToken      = types.NewToken
	NewKeeper     = keeper.NewKeeper

	ErrTokenNotExist       = types.ErrTokenNotExist
	ErrInsufficientBalance = types.ErrInsufficientBalance
)
