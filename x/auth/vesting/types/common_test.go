package types_test

import (
	"github.com/line/lfb-sdk/simapp"
)

var (
	app         = simapp.Setup(false)
	appCodec, _ = simapp.MakeCodecs()
)
