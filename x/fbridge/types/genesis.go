package types

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		SendingState: SendingState{
			NextSeq: 1,
		},
		ReceivingState:     ReceivingState{},
		NextRoleProposalId: 1,
		Roles:              []RolePair{{Role: RoleGuardian, Address: "<first guardian address>"}},
		BridgeSwitches:     []BridgeSwitch{{Guardian: "<first guardian address>", Status: StatusActive}},
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
