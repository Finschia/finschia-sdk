package types

import (
	ocmath "github.com/line/ostracon/libs/math"
	"github.com/line/ostracon/light"
)

// DefaultTrustLevel is the ostracon light client default trust level
var DefaultTrustLevel = NewFractionFromOc(light.DefaultTrustLevel)

// NewFractionFromOc returns a new Fraction instance from a ocmath.Fraction
func NewFractionFromOc(f ocmath.Fraction) Fraction {
	return Fraction{
		Numerator:   f.Numerator,
		Denominator: f.Denominator,
	}
}

// ToOstracon converts Fraction to ocmath.Fraction
func (f Fraction) ToOstracon() ocmath.Fraction {
	return ocmath.Fraction{
		Numerator:   f.Numerator,
		Denominator: f.Denominator,
	}
}
