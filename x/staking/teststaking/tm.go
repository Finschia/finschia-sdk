package teststaking

import (
	occrypto "github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"

	cryptocodec "github.com/Finschia/finschia-sdk/crypto/codec"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/staking/types"
)

// GetOcConsPubKey gets the validator's public key as an occrypto.PubKey.
func GetOcConsPubKey(v types.Validator) (occrypto.PubKey, error) {
	pk, err := v.ConsPubKey()
	if err != nil {
		return nil, err
	}

	return cryptocodec.ToTmPubKeyInterface(pk)
}

// ToOcValidator casts an SDK validator to a tendermint type Validator.
func ToOcValidator(v types.Validator, r sdk.Int) (*tmtypes.Validator, error) {
	ocPk, err := GetOcConsPubKey(v)
	if err != nil {
		return nil, err
	}

	return tmtypes.NewValidator(ocPk, v.ConsensusPower(r)), nil
}

// ToOcValidators casts all validators to the corresponding tendermint type.
func ToOcValidators(v types.Validators, r sdk.Int) ([]*tmtypes.Validator, error) {
	validators := make([]*tmtypes.Validator, len(v))
	var err error
	for i, val := range v {
		validators[i], err = ToOcValidator(val, r)
		if err != nil {
			return nil, err
		}
	}

	return validators, nil
}
