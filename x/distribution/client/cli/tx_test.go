package cli

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	"github.com/Finschia/finschia-sdk/simapp/params"
	"github.com/Finschia/finschia-sdk/testutil"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
)

func Test_splitAndCall_NoMessages(t *testing.T) {
	clientCtx := client.Context{}

	// empty tx will always trigger genOrBroadcastFn
	maxMsgsCases := []int{0, 1, 10}
	calledCounter := 0
	for _, maxMsgs := range maxMsgsCases {
		err := newSplitAndApply(
			func(clientCtx client.Context, fs *pflag.FlagSet, msgs ...sdk.Msg) error {
				// dummy genOrBroadcastFn called once for each case
				calledCounter++
				return nil
			}, clientCtx, nil, nil, maxMsgs)
		assert.NoError(t, err, "")
	}
	assert.Equal(t, calledCounter, len(maxMsgsCases))
}

func Test_splitAndCall_Splitting(t *testing.T) {
	clientCtx := client.Context{}

	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	// Add five messages
	msgs := []sdk.Msg{
		testdata.NewTestMsg(addr),
		testdata.NewTestMsg(addr),
		testdata.NewTestMsg(addr),
		testdata.NewTestMsg(addr),
		testdata.NewTestMsg(addr),
	}

	// Keep track of number of calls
	const chunkSize = 2

	callCount := 0
	err := newSplitAndApply(
		func(clientCtx client.Context, fs *pflag.FlagSet, msgs ...sdk.Msg) error {
			callCount++

			assert.NotNil(t, clientCtx)
			assert.NotNil(t, msgs)

			if callCount < 3 {
				assert.Equal(t, len(msgs), 2)
			} else {
				assert.Equal(t, len(msgs), 1)
			}

			return nil
		},
		clientCtx, nil, msgs, chunkSize)

	assert.NoError(t, err, "")
	assert.Equal(t, 3, callCount)
}

func TestParseProposal(t *testing.T) {
	encodingConfig := params.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Community Pool Spend",
  "description": "Pay me some Atoms!",
  "recipient": "link1s5afhd6gxevu37mkqcvvsj8qeylhn0rz46zdlq",
  "amount": "1000stake",
  "deposit": "1000stake"
}
`)

	proposal, err := ParseCommunityPoolSpendProposalWithDeposit(encodingConfig.Marshaler, okJSON.Name())
	require.NoError(t, err)

	require.Equal(t, "Community Pool Spend", proposal.Title)
	require.Equal(t, "Pay me some Atoms!", proposal.Description)
	require.Equal(t, "link1s5afhd6gxevu37mkqcvvsj8qeylhn0rz46zdlq", proposal.Recipient)
	require.Equal(t, "1000stake", proposal.Deposit)
	require.Equal(t, "1000stake", proposal.Amount)
}
