package types

import (
	"errors"

	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	foundationtypes "github.com/Finschia/finschia-sdk/x/foundation"
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

func AuthorityCandiates() []sdk.AccAddress {
	return []sdk.AccAddress{
		authtypes.NewModuleAddress(govtypes.ModuleName),
		authtypes.NewModuleAddress(foundationtypes.ModuleName),
	}
}

func ValidateGenesis(data GenesisState) error {
	if err := data.Params.ValidateParams(); err != nil {
		return err
	}

	if err := validateSendingState(data.SendingState); err != nil {
		return err
	}

	if data.NextRoleProposalId < 1 {
		return errors.New("next role proposal ID must be positive")
	}

	for _, v := range data.RoleProposals {
		if v.Id < 1 {
			return errors.New("role proposal ID must be positive")
		}
		sdk.MustAccAddressFromBech32(v.Proposer)
		sdk.MustAccAddressFromBech32(v.Target)
		if err := IsValidRole(v.Role); err != nil {
			return err
		}
	}

	for _, v := range data.Votes {
		if v.ProposalId < 1 {
			return errors.New("role proposal ID must be positive")
		}
		sdk.MustAccAddressFromBech32(v.Voter)
		if err := IsValidVoteOption(v.Option); err != nil {
			return err
		}
	}

	for _, v := range data.Roles {
		sdk.MustAccAddressFromBech32(v.Address)
		if err := IsValidRole(v.Role); err != nil {
			return err
		}
	}

	for _, v := range data.BridgeSwitches {
		sdk.MustAccAddressFromBech32(v.Guardian)
		if err := IsValidBridgeStatus(v.Status); err != nil {
			return err
		}
	}

	return nil
}

func validateSendingState(state SendingState) error {
	if state.NextSeq < 1 {
		return errors.New("next sequence must be positive")
	}

	if state.NextSeq-1 != uint64(len(state.SeqToBlocknum)) {
		return errors.New("sequence to blocknum mapping is invalid")
	}

	chkSeq := make(map[uint64]struct{})
	for _, v := range state.SeqToBlocknum {
		if v.Blocknum == 0 || v.Seq == 0 {
			return errors.New("blocknum and seq must be positive")
		}

		if _, ok := chkSeq[v.Seq]; ok {
			return errors.New("duplicate sequence")
		}

		chkSeq[v.Seq] = struct{}{}
	}

	return nil
}
