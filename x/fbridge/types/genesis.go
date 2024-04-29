package types

import "errors"

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
