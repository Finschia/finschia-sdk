package collection

import (
	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/collection/internal/types"
)

type GenesisState struct {
	Params types.Params `json:"params"`
	Tokens []Token      `json:"tokens"`
	// TODO: approvals
}

func NewGenesisState(params types.Params, tokens []Token) GenesisState {
	return GenesisState{
		Params: params,
		Tokens: tokens,
	}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(types.DefaultParams(), nil)
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// TODO: fill it with permission
	keeper.SetParams(ctx, data.Params)
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	params := keeper.GetParams(ctx)
	return NewGenesisState(params, nil)
}

func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return nil
}
