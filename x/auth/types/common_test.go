package types_test

import (
	"github.com/Finschia/finschia-rdk/simapp"
)

var (
	app                   = simapp.Setup(false)
	ecdc                  = simapp.MakeTestEncodingConfig()
	appCodec, legacyAmino = ecdc.Marshaler, ecdc.Amino
)
