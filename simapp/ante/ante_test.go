package ante_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"

	appante "github.com/Finschia/finschia-sdk/simapp/ante"
	"github.com/Finschia/finschia-sdk/x/auth/ante"
	"github.com/Finschia/finschia-sdk/x/zkauth/testutil"
)

func TestAnteHandler(t *testing.T) {
	k := testutil.ZkAuthKeeper(t)
	_, err := k.AddTestAccounts([]string{"link1z972lulhxgjyql3em8lj0hdx576mc99a5cpcvyrfsffyuqxe25lsxudr8l"})
	require.NoError(t, err)

	const sampleTxBase64 = "CrcCCrQCCiUvZmluc2NoaWEuemthdXRoLnYxYmV0YTEuTXNnRXhlY3V0aW9uEooCCiAKHC9jb3Ntb3MuYmFuay52MWJldGExLk1zZ1NlbmQSABLlAQreAQoAEiRhSFIwY0hNNkx5OWhZMk52ZFc1MGN5NW5iMjluYkdVdVkyOXQaZmV5SmhiR2NpT2lKU1V6STFOaUlzSW10cFpDSTZJalUxWXpFNE9HRTRNelUwTm1aak1UZzRaVFV4TlRjMlltRTNNamd6Tm1Vd05qQXdaVGhpTnpNaUxDSjBlWEFpT2lKS1YxUWlmUSJMNjQyNzU4NzQ4MjYxNDE4NjcwMDA5Mzc1NTYxNjE1MjEwMDcxMjIxMTQyNTQ4NjAzNjM4MTc1ODMyMTgxMDUxNDc5MzA3NTUwMDQ2MBCz/AESZApOCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAwQi109b7c6R9/HIIyz5OkrFQPV0Egis+NY87xM1MgKMEgQKAggBEhIKDAoEY29ueRIEMjAwMBDAmgwaQJwrov/j228lteH+sqxUT3HZWDx8HSAJq4+6dVPibNSOC8PZiuHSdk2noTGxDQ1BVjz3QviklWdy+6A/kZtSGL8="
	sampleTxBytes, err := base64.StdEncoding.DecodeString(sampleTxBase64)
	require.NoError(t, err)

	tx, err := k.ClientCtx.TxConfig.TxDecoder()(sampleTxBytes)
	require.NoError(t, err)

	anteHandler, err := appante.NewAnteHandler(
		appante.HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   k.Simapp.AccountKeeper,
				BankKeeper:      k.Simapp.BankKeeper,
				FeegrantKeeper:  k.Simapp.FeeGrantKeeper,
				SignModeHandler: k.ClientCtx.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			ZKAuthKeeper: k.ZKAuthKeeper,
		},
	)

	_, err = anteHandler(k.Ctx, tx, false)
	require.NoError(t, err)
}
