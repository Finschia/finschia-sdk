package token

import (
	"math"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

// ValidateGenesis check the given genesis state has no integrity issues
func ValidateGenesis(data GenesisState) error {
	if data.ClassState != nil {
		if err := ValidateClassGenesis(*data.ClassState); err != nil {
			return err
		}
	}

	for _, contractBalances := range data.Balances {
		if err := ValidateContractID(contractBalances.ContractId); err != nil {
			return err
		}

		if len(contractBalances.Balances) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("balances cannot be empty")
		}
		for _, balance := range contractBalances.Balances {
			if err := sdk.ValidateAccAddress(balance.Address); err != nil {
				return err
			}
			if err := validateAmount(balance.Amount); err != nil {
				return err
			}
		}
	}

	for _, c := range data.Classes {
		if err := ValidateContractID(c.ContractId); err != nil {
			return err
		}
		if err := validateName(c.Name); err != nil {
			return err
		}
		if err := validateSymbol(c.Symbol); err != nil {
			return err
		}
		if err := validateImageURI(c.ImageUri); err != nil {
			return err
		}
		if err := validateMeta(c.Meta); err != nil {
			return err
		}
		if err := validateDecimals(c.Decimals); err != nil {
			return err
		}
	}

	for _, contractGrants := range data.Grants {
		if err := ValidateContractID(contractGrants.ContractId); err != nil {
			return err
		}

		if len(contractGrants.Grants) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("grants cannot be empty")
		}
		for _, grant := range contractGrants.Grants {
			if err := sdk.ValidateAccAddress(grant.Grantee); err != nil {
				return err
			}
			if err := ValidatePermission(grant.Permission); err != nil {
				return err
			}
		}
	}

	for _, contractAuthorizations := range data.Authorizations {
		if err := ValidateContractID(contractAuthorizations.ContractId); err != nil {
			return err
		}

		if len(contractAuthorizations.Authorizations) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("authorizations cannot be empty")
		}
		for _, authorization := range contractAuthorizations.Authorizations {
			if err := sdk.ValidateAccAddress(authorization.Holder); err != nil {
				return err
			}
			if err := sdk.ValidateAccAddress(authorization.Operator); err != nil {
				return err
			}
		}
	}

	return nil
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{ClassState: DefaultClassGenesisState()}
}

// For Class keeper
func ValidateClassGenesis(data ClassGenesisState) error {
	if data.Nonce.GT(sdk.NewUint(math.MaxUint64)) {
		return sdkerrors.ErrInvalidRequest.Wrapf("Invalid nonce: %s", data.Nonce)
	}

	return nil
}

func DefaultClassGenesisState() *ClassGenesisState {
	return &ClassGenesisState{
		Nonce: sdk.ZeroUint(),
	}
}
