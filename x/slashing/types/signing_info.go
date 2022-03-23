package types

import (
	"fmt"
	"time"

	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
)

// NewValidatorSigningInfo creates a new ValidatorSigningInfo instance
//nolint:interfacer
func NewValidatorSigningInfo(
	condAddr sdk.ConsAddress, jailedUntil time.Time, tombstoned bool,
	missedBlocksCounter, voterSetCounter int64,
) ValidatorSigningInfo {

	return ValidatorSigningInfo{
		Address:             condAddr.String(),
		JailedUntil:         jailedUntil,
		Tombstoned:          tombstoned,
		MissedBlocksCounter: missedBlocksCounter,
		VoterSetCounter:     voterSetCounter,
	}
}

// String implements the stringer interface for ValidatorSigningInfo
func (i ValidatorSigningInfo) String() string {
	return fmt.Sprintf(`Validator Signing Info:
  Address:               %s
  Jailed Until:          %v
  Tombstoned:            %t
  Missed Blocks Counter: %d
  Voter Set Counter:     %d`,
		i.Address, i.JailedUntil, i.Tombstoned,
		i.MissedBlocksCounter, i.VoterSetCounter)
}

// unmarshal a validator signing info from a store value
func UnmarshalValSigningInfo(cdc codec.Codec, value []byte) (signingInfo ValidatorSigningInfo, err error) {
	err = cdc.Unmarshal(value, &signingInfo)
	return signingInfo, err
}
