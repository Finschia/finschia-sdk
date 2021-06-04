package simulation_test

import (
	"fmt"
	"testing"
	"time"

	types2 "github.com/line/lbm-sdk/v2/store/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/v2/simapp"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/evidence/simulation"
	"github.com/line/lbm-sdk/v2/x/evidence/types"
)

func TestDecodeStore(t *testing.T) {
	cdc, _ := simapp.MakeCodecs()
	dec := simulation.NewDecodeStore(cdc)

	delPk1 := ed25519.GenPrivKey().PubKey()

	ev := &types.Equivocation{
		Height:           10,
		Time:             time.Now().UTC(),
		Power:            1000,
		ConsensusAddress: sdk.ConsAddress(delPk1.Address()).String(),
	}

	kvPairs := []types2.KOPair{
			{
				Key:   types.KeyPrefixEvidence,
				Value: ev,
			},
			{
				Key:   []byte{0x99},
				Value: []byte{0x99},
			},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Evidence", fmt.Sprintf("%v\n%v", ev, ev)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec.LogPair(kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec.LogPair(kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
