package types

import (
	ostmath "github.com/line/ostracon/libs/math"
	"github.com/line/ostracon/light"
)

// DefaultTrustLevel is the tendermint light client default trust level
var DefaultTrustLevel = NewFractionFromTm(light.DefaultTrustLevel)

// NewFractionFromTm returns a new Fraction instance from a ostmath.Fraction
func NewFractionFromTm(f ostmath.Fraction) Fraction {
	return Fraction{
		Numerator:   f.Numerator,
		Denominator: f.Denominator,
	}
}

// ToTendermint converts Fraction to ostmath.Fraction
func (f Fraction) ToTendermint() ostmath.Fraction {
	return ostmath.Fraction{
		Numerator:   f.Numerator,
		Denominator: f.Denominator,
	}
}
