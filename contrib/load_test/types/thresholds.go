package types

import "time"

type Thresholds struct {
	Latency time.Duration
	TPS     int
}
