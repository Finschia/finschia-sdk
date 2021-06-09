package teststaking

import (
	ostcrypto "github.com/line/ostracon/crypto"
	osttypes "github.com/line/ostracon/types"

	cryptocodec "github.com/line/lfb-sdk/crypto/codec"
	"github.com/line/lfb-sdk/x/staking/types"
)

// GetTmConsPubKey gets the validator's public key as a ostcrypto.PubKey.
func GetTmConsPubKey(v types.Validator) (ostcrypto.PubKey, error) {
	pk, err := v.ConsPubKey()
	if err != nil {
		return nil, err
	}

	return cryptocodec.ToTmPubKeyInterface(pk)
}

// ToTmValidator casts an SDK validator to a tendermint type Validator.
func ToTmValidator(v types.Validator) (*osttypes.Validator, error) {
	tmPk, err := GetTmConsPubKey(v)
	if err != nil {
		return nil, err
	}

	return osttypes.NewValidator(tmPk, v.ConsensusPower()), nil
}

// ToTmValidators casts all validators to the corresponding tendermint type.
func ToTmValidators(v types.Validators) ([]*osttypes.Validator, error) {
	validators := make([]*osttypes.Validator, len(v))
	var err error
	for i, val := range v {
		validators[i], err = ToTmValidator(val)
		if err != nil {
			return nil, err
		}
	}

	return validators, nil
}
