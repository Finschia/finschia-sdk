package collection

import (
	"math"

	cmath "cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DefaultDepthLimit = 1
	DefaultWidthLimit = 4
)

// ValidateGenesis check the given genesis state has no integrity issues
func ValidateGenesis(data GenesisState) error {
	if err := validateParams(data.Params); err != nil {
		return err
	}

	for _, contract := range data.Contracts {
		if err := ValidateContractID(contract.Id); err != nil {
			return err
		}

		if err := validateName(contract.Name); err != nil {
			return err
		}
		if err := validateURI(contract.Uri); err != nil {
			return err
		}
		if err := validateMeta(contract.Meta); err != nil {
			return err
		}
	}

	for _, nextClassID := range data.NextClassIds {
		if err := ValidateContractID(nextClassID.ContractId); err != nil {
			return err
		}
	}

	for _, contractClasses := range data.Classes {
		if err := ValidateContractID(contractClasses.ContractId); err != nil {
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

	for _, contractNextTokenIDs := range data.NextTokenIds {
		if err := ValidateContractID(contractNextTokenIDs.ContractId); err != nil {
			return err
		}

		if len(contractNextTokenIDs.TokenIds) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("next token ids cannot be empty")
		}
		for _, nextTokenIDs := range contractNextTokenIDs.TokenIds {
			if err := ValidateClassID(nextTokenIDs.ClassId); err != nil {
				return err
			}
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
			if _, err := sdk.AccAddressFromBech32(balance.Address); err != nil {
				return err
			}
			if err := balance.Amount.ValidateBasic(); err != nil {
				return err
			}
		}
	}

	for _, contractNFTs := range data.Nfts {
		if err := ValidateContractID(contractNFTs.ContractId); err != nil {
			return err
		}

		if len(contractNFTs.Nfts) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("nfts cannot be empty")
		}
		for _, token := range contractNFTs.Nfts {
			if err := ValidateTokenID(token.TokenId); err != nil {
				return err
			}
			if err := validateName(token.Name); err != nil {
				return err
			}
			if err := validateMeta(token.Meta); err != nil {
				return err
			}
		}
	}

	for _, contractParents := range data.Parents {
		if err := ValidateContractID(contractParents.ContractId); err != nil {
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
		if err := ValidateContractID(contractAuthorizations.ContractId); err != nil {
			return err
		}

		if len(contractAuthorizations.Authorizations) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("authorizations cannot be empty")
		}
		for _, authorization := range contractAuthorizations.Authorizations {
			if _, err := sdk.AccAddressFromBech32(authorization.Holder); err != nil {
				return err
			}
			if _, err := sdk.AccAddressFromBech32(authorization.Operator); err != nil {
				return err
			}
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
			if _, err := sdk.AccAddressFromBech32(grant.Grantee); err != nil {
				return err
			}
			if err := ValidatePermission(grant.Permission); err != nil {
				return err
			}
		}
	}

	for _, contractSupplies := range data.Supplies {
		if err := ValidateContractID(contractSupplies.ContractId); err != nil {
			return err
		}

		if len(contractSupplies.Statistics) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("supplies cannot be empty")
		}
		for _, supply := range contractSupplies.Statistics {
			if err := ValidateClassID(supply.ClassId); err != nil {
				return err
			}
			if !supply.Amount.IsPositive() {
				return sdkerrors.ErrInvalidRequest.Wrap("supply must be positive")
			}
		}
	}

	for _, contractBurnts := range data.Burnts {
		if err := ValidateContractID(contractBurnts.ContractId); err != nil {
			return err
		}

		if len(contractBurnts.Statistics) == 0 {
			return sdkerrors.ErrInvalidRequest.Wrap("burnts cannot be empty")
		}
		for _, burnt := range contractBurnts.Statistics {
			if err := ValidateClassID(burnt.ClassId); err != nil {
				return err
			}
			if !burnt.Amount.IsPositive() {
				return sdkerrors.ErrInvalidRequest.Wrap("burnt must be positive")
			}
		}
	}

	if data.ClassState != nil {
		if data.ClassState.Nonce.GT(cmath.NewUint(math.MaxUint64)) {
			return sdkerrors.ErrInvalidRequest.Wrapf("Invalid nonce: %s", data.ClassState.Nonce)
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
