package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Tokens []Token `json:"tokens"`
	// TODO: approvals
}

func NewGenesisState(tokens []Token) GenesisState {
	return GenesisState{Tokens: tokens}
}

func DefaultGenesisState() GenesisState { return NewGenesisState(nil) }

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// TODO: fill it with permission
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(keeper.GetAllTokens(ctx))
}

func ValidateGenesis(data GenesisState) error { return nil }

// TODO: validate
