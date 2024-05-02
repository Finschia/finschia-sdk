package types

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"time"
)

func DefaultParams() Params {
	return Params{
		GuardianTrustLevel: Fraction{Numerator: 2, Denominator: 3},
		OperatorTrustLevel: Fraction{Numerator: 2, Denominator: 3},
		JudgeTrustLevel:    Fraction{Numerator: 1, Denominator: 1},
		ProposalPeriod:     uint64(time.Minute * 60),
	}
}

func CheckTrustLevelThreshold(total, current uint64, trustLevel Fraction) bool {
	if err := ValidateTrustLevel(trustLevel); err != nil {
		panic(err)
	}

	if total*trustLevel.Numerator <= current*trustLevel.Denominator &&
		total != 0 && current != 0 &&
		current <= total {
		return true
	}

	return false
}

func ValidateTrustLevel(trustLevel Fraction) error {
	if trustLevel.Denominator < 1 || trustLevel.Numerator < 1 {
		return sdkerrors.ErrInvalidRequest.Wrap("trust level must be positive")
	} else if trustLevel.Denominator < trustLevel.Numerator {
		return sdkerrors.ErrInvalidRequest.Wrap("trust level denominator must be greater than or equal to the numerator")
	}

	return nil
}
