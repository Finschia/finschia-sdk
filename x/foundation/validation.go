package foundation

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const (
	maxMetadataLen = 100
)
	
func validateMetadata(metadata string) error {
	if len(metadata) > maxMetadataLen {
		return sdkerrors.ErrInvalidRequest.Wrap("metadata is too large")
	}

	return nil
}

func validateMembers(members []Member) error {
	addrs := map[string]bool{}
	for _, member := range members {
		if err := member.ValidateBasic(); err != nil {
			return err
		}
		if addrs[member.Address] {
			return sdkerrors.ErrInvalidRequest.Wrapf("duplicated address: %s", member.Address)
		}
		addrs[member.Address] = true
	}

	return nil
}

func (m Member) ValidateBasic() error {
	if err := sdk.ValidateAccAddress(m.Address); err != nil {
		return err
	}

	if m.Weight.IsNegative() {
		return sdkerrors.ErrInvalidRequest.Wrapf("expected a non-negative decimal, got %s", m.Weight)
	}

	return nil
}
