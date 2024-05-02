package types

import (
	"errors"

	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
)

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		SendingState: SendingState{
			NextSeq: 1,
		},
		ReceivingState: ReceivingState{},
	}
}

func ValidateGenesis(data GenesisState) error {
	if data.SendingState.NextSeq < 1 {
		panic(errors.New("next sequence must be positive"))
	}

	return nil
}

func DefaultAuthority() sdk.AccAddress {
	return authtypes.NewModuleAddress(govtypes.ModuleName)
}
