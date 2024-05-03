package types

import (
	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	sdk "github.com/Finschia/finschia-sdk/types"
)

func DefaultGenesisState() *GenesisState {
	dummyGuardian := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().String())

	return &GenesisState{
		Params: DefaultParams(),
		SendingState: SendingState{
			NextSeq: 1,
		},
		ReceivingState:     ReceivingState{},
		NextRoleProposalId: 1,
		// WARN: you must set your own guardian address in production
		Roles: []RolePair{{Role: RoleGuardian, Address: dummyGuardian.String()}},
		// WARN: you must set your own guardian address in production
		BridgeSwitches: []BridgeSwitch{{Guardian: dummyGuardian.String(), Status: StatusActive}},
	}
}

func ValidateGenesis(data GenesisState) error {
	if data.SendingState.NextSeq < 1 {
		panic("next sequence must be positive")
	}

	if data.NextRoleProposalId < 1 {
		panic("next role proposal ID must be positive")
	}

	return nil
}
