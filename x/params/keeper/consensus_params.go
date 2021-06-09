package keeper

import (
	abci "github.com/line/ostracon/abci/types"
	ostproto "github.com/line/ostracon/proto/ostracon/types"

	"github.com/line/lfb-sdk/baseapp"
	"github.com/line/lfb-sdk/x/params/types"
)

// ConsensusParamsKeyTable returns an x/params module keyTable to be used in
// the BaseApp's ParamStore. The KeyTable registers the types along with the
// standard validation functions. Applications can choose to adopt this KeyTable
// or provider their own when the existing validation functions do not suite their
// needs.
func ConsensusParamsKeyTable() types.KeyTable {
	return types.NewKeyTable(
		types.NewParamSetPair(
			baseapp.ParamStoreKeyBlockParams, abci.BlockParams{}, baseapp.ValidateBlockParams,
		),
		types.NewParamSetPair(
			baseapp.ParamStoreKeyEvidenceParams, ostproto.EvidenceParams{}, baseapp.ValidateEvidenceParams,
		),
		types.NewParamSetPair(
			baseapp.ParamStoreKeyValidatorParams, ostproto.ValidatorParams{}, baseapp.ValidateValidatorParams,
		),
	)
}
