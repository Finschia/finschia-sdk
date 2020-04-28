package loadgenerator

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestRampUpPacer_Pace(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		rampUpTime   time.Duration
		freq         int
		per          time.Duration
		elapsedTime  time.Duration
		hits         uint64
		expectedWait time.Duration
		expectedStop bool
	}{
		{25 * time.Second, 100, time.Second, 0, 0, 0, false},
		{25 * time.Second, 100, time.Second, 1 * time.Second, 0, 0, false},
		{25 * time.Second, 100, time.Second, 1 * time.Second, 2, 0, false},
		{25 * time.Second, 100, time.Second, 2 * time.Second, 9, 0, false},
		{25 * time.Second, 100, time.Second, 3 * time.Second, 20, 0, false},
		{25 * time.Second, 100, time.Second, 4 * time.Second, 35, 0, false},
		{25 * time.Second, 100, time.Second, 5 * time.Second, 54, 0, false},
		{25 * time.Second, 100, time.Second, 6 * time.Second, 77, 0, false},
		{25 * time.Second, 100, time.Second, 16 * time.Second, 527, 0, false},
		{25 * time.Second, 100, time.Second, 32 * time.Second, 1981, 0, false},
		{25 * time.Second, 100, time.Second, 64 * time.Second, 5213, 0, false},
		{25 * time.Second, 100, time.Second, 128 * time.Second, 11677, 0, false},
		{25 * time.Second, 100, time.Second, 1 * time.Second, 3, 200 * time.Millisecond, false},
		{25 * time.Second, 100, time.Second, 6 * time.Second, 78, 40 * time.Millisecond, false},
		{25 * time.Second, 100, time.Second, 32 * time.Second, 1982, 10 * time.Millisecond, false},
		{25 * time.Second, 100, time.Second, 64 * time.Second, 5214, 10 * time.Millisecond, false},
		{25 * time.Second, 100, time.Second, 128 * time.Second, 11678, 10 * time.Millisecond, false},

		// Zero frequency.
		{25 * time.Second, 0, time.Second, time.Second, 0, 0, false},
		// Zero per.
		{25 * time.Second, 1, 0, time.Second, 0, 0, false},
		// Zero frequency + per.
		{25 * time.Second, 0, 0, time.Second, 0, 0, false},
		// Negative frequency.
		{25 * time.Second, -1, time.Second, time.Second, 0, 0, true},
		// Negative per.
		{25 * time.Second, 1, -time.Second, time.Second, 0, 0, true},
		// Negative frequency + per.
		{25 * time.Second, -1, -time.Second, time.Second, 0, 0, true},
	}

	for i, tt := range tests {
		t.Logf("Test #%d", i)
		{
			p := RampUpPacer{
				Constant: vegeta.Rate{
					Freq: tt.freq,
					Per:  tt.per,
				},
				RampUpTime: tt.rampUpTime,
			}

			wait, stop := p.Pace(tt.elapsedTime, tt.hits)

			require.Equal(t, tt.expectedWait, wait)
			require.Equal(t, tt.expectedStop, stop)
		}
	}
}

func TestRampUpPacer_hits(t *testing.T) {
	var tests = []struct {
		rampUpTime   time.Duration
		freq         int
		elapsedTime  time.Duration
		expectedHits float64
	}{
		{0 * time.Second, 100, 0, 0},
		{0 * time.Second, 100, time.Second / 2, 50.5},
		{0 * time.Second, 100, 1 * time.Second, 101},
		{0 * time.Second, 100, 2 * time.Second, 202},
		{0 * time.Second, 100, 4 * time.Second, 404},
		{0 * time.Second, 100, 8 * time.Second, 808},
		{0 * time.Second, 100, 16 * time.Second, 1616},
		{0 * time.Second, 100, 32 * time.Second, 3232},
		{0 * time.Second, 100, 64 * time.Second, 6464},
		{0 * time.Second, 100, 128 * time.Second, 12928},

		{25 * time.Second, 100, time.Second / 2, 1},
		{25 * time.Second, 100, 1 * time.Second, 3},
		{25 * time.Second, 100, 2 * time.Second, 10},
		{25 * time.Second, 100, 4 * time.Second, 36},
		{25 * time.Second, 100, 8 * time.Second, 136},
		{25 * time.Second, 100, 16 * time.Second, 528},
		{25 * time.Second, 100, 32 * time.Second, 1982},
		{25 * time.Second, 100, 64 * time.Second, 5214},
		{25 * time.Second, 100, 128 * time.Second, 11678},
	}

	for i, tt := range tests {
		t.Logf("Test #%d", i)
		{
			p := RampUpPacer{
				Constant: vegeta.Rate{
					Freq: tt.freq,
					Per:  time.Second,
				},
				RampUpTime: tt.rampUpTime,
			}

			hits := p.hits(tt.elapsedTime)

			floatDelta := 1e-6 * math.Min(math.Abs(tt.expectedHits), math.Abs(hits))
			require.InDelta(t, tt.expectedHits, hits, floatDelta)
		}
	}
}

func TestRampUpPacer_Rate(t *testing.T) {
	var tests = []struct {
		rampUpTime   time.Duration
		freq         int
		elapsedTime  time.Duration
		expectedRate float64
	}{
		{0 * time.Second, 100, 10 * time.Second, 100},
		{0 * time.Second, 300, 10 * time.Second, 300},
		{10 * time.Second, 100, 5 * time.Second, 51},
		{10 * time.Second, 100, 6 * time.Second, 61},
		{10 * time.Second, 100, 7 * time.Second, 71},
		{10 * time.Second, 100, 8 * time.Second, 81},
		{10 * time.Second, 100, 10 * time.Second, 100},
		{10 * time.Second, 100, 15 * time.Second, 100},
	}

	for i, tt := range tests {
		t.Logf("Test #%d", i)
		{
			p := RampUpPacer{
				Constant: vegeta.Rate{
					Freq: tt.freq,
					Per:  time.Second,
				},
				RampUpTime: tt.rampUpTime,
			}

			fmt.Println(p.Constant.Rate(10 * time.Second))

			require.Equal(t, tt.expectedRate, p.Rate(tt.elapsedTime))
		}
	}
}
