package types_test

import (
	"github.com/line/lbm-sdk/v2/simapp"
)

var (
	app                   = simapp.Setup(false)
	appCodec, legacyAmino = simapp.MakeCodecs()
)
