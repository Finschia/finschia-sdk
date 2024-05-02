package types

import "time"

func DefaultParams() Params {
	return Params{
		GuardianTrustLevel: Fraction{Numerator: 2, Denominator: 3},
		OperatorTrustLevel: Fraction{Numerator: 2, Denominator: 3},
		JudgeTrustLevel:    Fraction{Numerator: 2, Denominator: 3},
		ProposalPeriod:     uint64(time.Minute * 60),
	}
}
