package foundation

import (
	"time"
)

// Config is a config struct used for intialising the group module to avoid using globals.
type Config struct {
	// MaxExecutionPeriod defines the max duration after a proposal's voting period ends that members can send a MsgExec to execute the proposal.
	MaxExecutionPeriod time.Duration
	// MaxMetadataLen defines the max length of the metadata bytes field for various entities within the foundation module. Defaults to 255 if not explicitly set.
	MaxMetadataLen uint64
}

func DefaultConfig() Config {
	return Config{
		MaxExecutionPeriod: 2 * 7 * 24 * time.Hour, // two weeks
		MaxMetadataLen:     255,
	}
}
