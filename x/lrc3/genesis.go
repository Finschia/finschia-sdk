package lrc3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	nft "github.com/link-chain/link/x/nft"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	LRC3s []nft.Collections `json:"lrc3_records"`
}

func NewGenesisState(lrc3Records []nft.Collections) GenesisState {
	return GenesisState{LRC3s: nil}
}

func ValidateGenesis(data GenesisState) error {
	// TODO: validate
	return nil
}

func DefaultGenesisState() GenesisState {
	return NewGenesisState(nil)
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// TODO:
	return NewGenesisState(nil)
}
