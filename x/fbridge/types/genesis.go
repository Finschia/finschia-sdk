package types

import (
	"errors"
	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
)

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		SendingState: SendingState{
			NextSeq: 1,
		},
		ReceivingState:     ReceivingState{},
		NextRoleProposalId: 1,
	}
}

func DefaultAuthority() sdk.AccAddress {
	return authtypes.NewModuleAddress(govtypes.ModuleName)
}

func ValidateGenesis(data GenesisState) error {
	if err := ValidateParams(data.Params); err != nil {
		return err
	}

	if data.SendingState.NextSeq < 1 {
		return errors.New("next sequence must be positive")
	}

	if data.NextRoleProposalId < 1 {
		return errors.New("next role proposal ID must be positive")
	}

	return nil
}
