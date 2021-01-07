package types_test

import (
	"github.com/line/lbm-sdk/simapp"
)

var (
	app      = simapp.Setup(false)
	appCodec = simapp.MakeTestEncodingConfig().Marshaler
)
