package token

import (
	"math"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token/class"
)

// ValidateGenesis check the given genesis state has no integrity issues
func ValidateGenesis(data GenesisState) error {
	if data.ClassState != nil {
		if err := ValidateClassGenesis(*data.ClassState); err != nil {
			return err
		}
	}

	for _, balance := range data.Balances {
		if err := sdk.ValidateAccAddress(balance.Address); err != nil {
			return err
		}
		if len(balance.Tokens) == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "tokens cannot be empty")
		}
		for _, amount := range balance.Tokens {
			if err := validateAmount(amount.Amount); err != nil {
				return err
			}
		}
	}

	for _, c := range data.Classes {
		if err := class.ValidateID(c.Id); err != nil {
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

	for _, grant := range data.Grants {
		if err := sdk.ValidateAccAddress(grant.Grantee); err != nil {
			return err
		}
		if err := validateAction(grant.Action); err != nil {
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
	return &GenesisState{ClassState: DefaultClassGenesisState()}
}

// For Class keeper
func ValidateClassGenesis(data ClassGenesisState) error {
	if data.Nonce.GT(sdk.NewUint(math.MaxUint64)) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "Invalid nonce: %s", data.Nonce)
	}

	return nil
}

func DefaultClassGenesisState() *ClassGenesisState {
	return &ClassGenesisState{
		Nonce: sdk.ZeroUint(),
	}
}
