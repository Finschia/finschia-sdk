package loadgenerator

import (
	"math"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

// RampUpPacer paces an attack by starting at a 1 RPS and increasing linearly during RampUpTime.
// After RampUpTime, it keeps constant rate.
type RampUpPacer struct {
	Constant   vegeta.Rate
	RampUpTime time.Duration
}

func (p RampUpPacer) Pace(elapsed time.Duration, hits uint64) (time.Duration, bool) {
	switch {
	case p.Constant.Per == 0 || p.Constant.Freq == 0:
		return 0, false // Zero value = infinite rate
	case p.Constant.Per < 0 || p.Constant.Freq < 0:
		return 0, true
	}

	expectedHits := p.hits(elapsed)
	if hits == 0 || hits < uint64(expectedHits) {
		// Running behind, send next hit immediately.
		return 0, false
	}

	rate := p.Rate(elapsed)
	interval := math.Round(1e9 / rate)

	if n := uint64(interval); n != 0 && math.MaxInt64/n < hits {
		// We would overflow wait if we continued, so stop the attack.
		return 0, true
	}

	delta := float64(hits+1) - expectedHits
	wait := time.Duration(interval * delta)

	return wait, false
}

func (p RampUpPacer) Rate(elapsed time.Duration) float64 {
	if elapsed < p.RampUpTime {
		a := float64(p.Constant.Freq) / float64(p.RampUpTime)
		x := float64(elapsed)
		return a*x + 1
	}
	return p.Constant.Rate(elapsed)
}

func (p RampUpPacer) hits(t time.Duration) float64 {
	if t < 0 {
		return 0
	}
	x := t.Seconds()
	slope := float64(p.Constant.Freq) / p.RampUpTime.Seconds()

	if t < p.RampUpTime {
		return slope*math.Pow(x, 2)/2 + x
	}
	return (x-p.RampUpTime.Seconds()/2)*float64(p.Constant.Freq) + x
}
