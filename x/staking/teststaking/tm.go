package teststaking

import (
	occrypto "github.com/Finschia/ostracon/crypto"
	octypes "github.com/Finschia/ostracon/types"

	cryptocodec "github.com/Finschia/finschia-rdk/crypto/codec"
	sdk "github.com/Finschia/finschia-rdk/types"
	"github.com/Finschia/finschia-rdk/x/staking/types"
)

// GetOcConsPubKey gets the validator's public key as an occrypto.PubKey.
func GetOcConsPubKey(v types.Validator) (occrypto.PubKey, error) {
	pk, err := v.ConsPubKey()
	if err != nil {
		return nil, err
	}

	return cryptocodec.ToOcPubKeyInterface(pk)
}

// ToOcValidator casts an SDK validator to a tendermint type Validator.
func ToOcValidator(v types.Validator, r sdk.Int) (*octypes.Validator, error) {
	ocPk, err := GetOcConsPubKey(v)
	if err != nil {
		return nil, err
	}

	return octypes.NewValidator(ocPk, v.ConsensusPower(r)), nil
}

// ToOcValidators casts all validators to the corresponding tendermint type.
func ToOcValidators(v types.Validators, r sdk.Int) ([]*octypes.Validator, error) {
	validators := make([]*octypes.Validator, len(v))
	var err error
	for i, val := range v {
		validators[i], err = ToOcValidator(val, r)
		if err != nil {
			return nil, err
		}
	}

	return validators, nil
}
