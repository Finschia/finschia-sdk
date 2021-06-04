package simulation_test

import (
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	types2 "github.com/line/lbm-sdk/v2/store/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/v2/simapp"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/gov/simulation"
	"github.com/line/lbm-sdk/v2/x/gov/types"
)

var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.AccAddress(delPk1.Address())
)

func TestDecodeStore(t *testing.T) {
	cdc, _ := simapp.MakeCodecs()
	dec := simulation.NewDecodeStore(cdc)

	endTime := time.Now().UTC()

	content := types.ContentFromProposalType("test", "test", types.ProposalTypeText)
	proposal, err := types.NewProposal(content, 1, endTime, endTime.Add(24*time.Hour))
	require.NoError(t, err)

	proposalIDBz := make([]byte, 8)
	binary.LittleEndian.PutUint64(proposalIDBz, 1)
	deposit := types.NewDeposit(1, delAddr1, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.OneInt())))
	vote := types.NewVote(1, delAddr1, types.OptionYes)

	kvPairs := []types2.KOPair{
			{Key: types.ProposalKey(1), Value: &proposal},
			{Key: types.InactiveProposalQueueKey(1, endTime), Value: proposalIDBz},
			{Key: types.DepositKey(1, delAddr1), Value: &deposit},
			{Key: types.VoteKey(1, delAddr1), Value: &vote},
			{Key: []byte{0x99}, Value: []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"proposals", fmt.Sprintf("%v\n%v", proposal, proposal)},
		{"proposal IDs", "proposalIDA: 1\nProposalIDB: 1"},
		{"deposits", fmt.Sprintf("%v\n%v", deposit, deposit)},
		{"votes", fmt.Sprintf("%v\n%v", vote, vote)},
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
