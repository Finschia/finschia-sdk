package types_test

import (
	"github.com/Finschia/finschia-sdk/simapp"
)

var (
	app         = simapp.Setup(false)
	ecdc        = simapp.MakeTestEncodingConfig()
	appCodec, _ = ecdc.Marshaler, ecdc.Amino
)
