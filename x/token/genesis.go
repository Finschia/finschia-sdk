package token

import (
	"math"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

// ValidateGenesis check the given genesis state has no integrity issues
func ValidateGenesis(data GenesisState) error {
	if err := validateClassGenesis(*data.ClassState); err != nil {
		return err
	}

	for _, balance := range data.Balances {
		if err := sdk.ValidateAccAddress(balance.Address); err != nil {
			return err
		}
		for _, amount := range balance.Tokens {
			if err := validateAmount(amount.Amount); err != nil {
				return err
			}
		}
	}

	for _, class := range data.Classes {
		if err := validateName(class.Name); err != nil {
			return err
		}
		if err := validateSymbol(class.Symbol); err != nil {
			return err
		}
		if err := validateImageUri(class.ImageUri); err != nil {
			return err
		}
		if err := validateMeta(class.Meta); err != nil {
			return err
		}
		if err := validateDecimals(class.Decimals); err != nil {
			return err
		}
	}

	for _, grant := range data.Grants {
		if err := sdk.ValidateAccAddress(grant.Grantee); err != nil {
			return err
		}
	}

	for _, approve := range data.Approves {
		if err := sdk.ValidateAccAddress(approve.Approver); err != nil {
			return err
		}
		if err := sdk.ValidateAccAddress(approve.Proxy); err != nil {
			return err
		}
	}

	return nil
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{ClassState: defaultClassGenesisState()}
}

// For Class keeper
func validateClassGenesis(data ClassGenesisState) error {
	if data.Nonce.GT(sdk.NewUint(math.MaxUint64)) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid nonce: %s", data.Nonce)
	}
	
	return nil
}

func defaultClassGenesisState() *ClassGenesisState {
	return &ClassGenesisState{
		Nonce: sdk.ZeroUint(),
	}
}
