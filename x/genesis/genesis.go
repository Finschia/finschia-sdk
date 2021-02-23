package genesis

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint:golint
type GenesisState struct {
	GenesisMessage string `json:"genesis_message"`
}

func NewGenesisState(genesisMessage string) GenesisState {
	return GenesisState{GenesisMessage: genesisMessage}
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState("In the beginning God created the heavens and the earth.")
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetGenesisMessage(ctx, data.GenesisMessage)
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(keeper.GetGenesisMessage(ctx))
}

func ValidateGenesis(data GenesisState) error { return nil }
