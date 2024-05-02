package types

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
)

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		SendingState: SendingState{
			NextSeq: 1,
		},
		ReceivingState:     ReceivingState{},
		NextRoleProposalId: 1,
		RoleMetadata:       RoleMetadata{Guardian: 0, Operator: 0, Judge: 0},
	}
}

func ValidateGenesis(data GenesisState) error {
	if data.SendingState.NextSeq < 1 {
		panic("next sequence must be positive")
	}

	if data.NextRoleProposalId < 1 {
		panic("next role proposal ID must be positive")
	}

	if data.RoleMetadata.Guardian < 0 || data.RoleMetadata.Operator < 0 || data.RoleMetadata.Judge < 0 {
		panic("length of each group must be positive")
	}

	return nil
}

func DefaultAuthority() sdk.AccAddress {
	return authtypes.NewModuleAddress(govtypes.ModuleName)
}
