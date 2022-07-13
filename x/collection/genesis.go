package collection

import (
	codectypes "github.com/line/lbm-sdk/codec/types"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"

	"github.com/line/lbm-sdk/x/token/class"
)

const (
	DefaultDepthLimit = 3
	DefaultWidthLimit = 8
)

// ValidateGenesis check the given genesis state has no integrity issues
func ValidateGenesis(data GenesisState) error {

	// the legacy module did not validate the data.
	// if LegacyMode {
	// 	return nil
	// }

	for _, contract := range data.Contracts {
		if err := class.ValidateID(contract.ContractId); err != nil {
			return err
		}

		if err := validateName(contract.Name); err != nil {
			return err
		}
		if err := validateBaseImgURI(contract.BaseImgUri); err != nil {
			return err
		}
		if err := validateMeta(contract.Meta); err != nil {
			return err
		}
	}

	for _, contractClasses := range data.Classes {
		if err := class.ValidateID(contractClasses.ContractId); err != nil {
			return err
		}

		if len(contractClasses.Classes) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("classes cannot be empty")
		}
		for i := range contractClasses.Classes {
			any := &contractClasses.Classes[i]
			class := TokenClassFromAny(any)
			if err := class.ValidateBasic(); err != nil {
				return err
			}
		}
	}

	for _, contractBalances := range data.Balances {
		if err := class.ValidateID(contractBalances.ContractId); err != nil {
			return err
		}

		if len(contractBalances.Balances) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("balances cannot be empty")
		}
		for _, balance := range contractBalances.Balances {
			if err := sdk.ValidateAccAddress(balance.Address); err != nil {
				return err
			}
			if err := balance.Amount.ValidateBasic(); err != nil {
				return err
			}
		}
	}

	for _, contractParents := range data.Parents {
		if err := class.ValidateID(contractParents.ContractId); err != nil {
			return err
		}

		if len(contractParents.Relations) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("parents cannot be empty")
		}
		for _, relation := range contractParents.Relations {
			if err := ValidateTokenID(relation.Self); err != nil {
				return err
			}
			if err := ValidateTokenID(relation.Other); err != nil {
				return err
			}
		}
	}

	for _, contractAuthorizations := range data.Authorizations {
		if err := class.ValidateID(contractAuthorizations.ContractId); err != nil {
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

	for _, contractGrants := range data.Grants {
		if err := class.ValidateID(contractGrants.ContractId); err != nil {
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

	return nil
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: Params{
			DepthLimit: DefaultDepthLimit,
			WidthLimit: DefaultWidthLimit,
		},
	}
}

func (data GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, contractClasses := range data.Classes {
		for i := range contractClasses.Classes {
			any := &contractClasses.Classes[i]
			if err := TokenClassUnpackInterfaces(any, unpacker); err != nil {
				return err
			}
		}
	}

	return nil
}
