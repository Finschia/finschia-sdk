package types

import (
	"time"

	sdktypes "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

func DefaultParams() Params {
	return Params{
		GuardianTrustLevel: Fraction{Numerator: 2, Denominator: 3},
		OperatorTrustLevel: Fraction{Numerator: 2, Denominator: 3},
		JudgeTrustLevel:    Fraction{Numerator: 1, Denominator: 1},
		ProposalPeriod:     uint64(time.Minute * 60),
		TimelockPeriod:     uint64(time.Hour * 24),
		TargetDenom:        sdktypes.DefaultBondDenom,
	}
}

func (p Params) ValidateParams() error {
	if err := ValidateTrustLevel(p.GuardianTrustLevel); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("guardian trust level: " + err.Error())
	}

	if err := ValidateTrustLevel(p.OperatorTrustLevel); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("operator trust level: " + err.Error())
	}

	if err := ValidateTrustLevel(p.JudgeTrustLevel); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrap("judge trust level: " + err.Error())
	}

	if p.ProposalPeriod == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("proposal period cannot be 0")
	}

	if p.TimelockPeriod == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("timelock period cannot be 0")
	}

	if err := sdktypes.ValidateDenom(p.TargetDenom); err != nil {
		return err
	}

	return nil
}

func CheckTrustLevelThreshold(total, current uint64, trustLevel Fraction) bool {
	if err := ValidateTrustLevel(trustLevel); err != nil {
		panic(err)
	}

	if total*trustLevel.Numerator <= current*trustLevel.Denominator &&
		total > 0 &&
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
