package v040_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/client"
	"github.com/line/lfb-sdk/simapp"
	sdk "github.com/line/lfb-sdk/types"
	v038evidence "github.com/line/lfb-sdk/x/evidence/legacy/v038"
	v040evidence "github.com/line/lfb-sdk/x/evidence/legacy/v040"
)

func TestMigrate(t *testing.T) {
	encodingConfig := simapp.MakeTestEncodingConfig()
	clientCtx := client.Context{}.
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithJSONMarshaler(encodingConfig.Marshaler)

	addr1, _ := sdk.AccAddressFromBech32("link19d880rl0u36uc4yhdjrsp5yculgd4ak7hdt660")

	evidenceGenState := v038evidence.GenesisState{
		Params: v038evidence.Params{MaxEvidenceAge: v038evidence.DefaultMaxEvidenceAge},
		Evidence: []v038evidence.Evidence{v038evidence.Equivocation{
			Height:           20,
			Power:            100,
			ConsensusAddress: addr1.Bytes(),
		}},
	}

	migrated := v040evidence.Migrate(evidenceGenState)
	expected := `{"evidence":[{"@type":"/lfb.evidence.v1beta1.Equivocation","height":"20","time":"0001-01-01T00:00:00Z","power":"100","consensus_address":"linkvalcons19d880rl0u36uc4yhdjrsp5yculgd4ak7326mca"}]}`

	bz, err := clientCtx.JSONMarshaler.MarshalJSON(migrated)
	require.NoError(t, err)
	require.Equal(t, expected, string(bz))
}
