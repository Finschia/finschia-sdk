package types

// For testing purposes, you must not use the DummyGuardian address in production
const DummyGuardian = "link1zmm9v8wucqecl75q22hddz0qypdgyvdpgg9a6d"

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		SendingState: SendingState{
			NextSeq: 1,
		},
		ReceivingState:     ReceivingState{},
		NextRoleProposalId: 1,
		// WARN: you must set your own guardian address in production
		Roles: []RolePair{{Role: RoleGuardian, Address: DummyGuardian}},
		// WARN: you must set your own guardian address in production
		BridgeSwitches: []BridgeSwitch{{Guardian: DummyGuardian, Status: StatusActive}},
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
