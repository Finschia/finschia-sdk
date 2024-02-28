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
	_, err := k.AddTestAccounts([]string{"link19qvkxwmln5kf0z59ecw744ue2gzsndlwcuz9uq2kav4evuerjlysxduwzj"})
	require.NoError(t, err)

	const sampleTxBase64 = "CtcICtQICiUvZmluc2NoaWEuemthdXRoLnYxYmV0YTEuTXNnRXhlY3V0aW9uEqoICo4BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm4KP2xpbmsxZzd1ZDYzZXFsbGo3emo0cTdma2NhNWg3czIyM2o3OHR5dnIwZTJjeHV3NHF5eWFhZjN1c2E2NGRxYxIrbGluazEwMDh3ZW5ncjI4ejVxdWF0MmR6cnBydDloOGV1YXY0aGVyZnl1bRKWBwqPBwqvBXsicGlfYSI6WyIzNjE1OTE2NzQ3NTQxMzkyMzMzOTA1MjYyNjkyNjQ3NTE1MjE2MDc2MTQzMTkwNzQwNDUyOTEyNjMzMzExMDAyMDc4NjYzNjEwNTU5IiwiMjAxMzg5MzA5NDk5MjEyMTA2Njk0NDYyMDA4ODAwNzI2MDg2MjUwOTM0MzQ5OTc3NDc2NDkyMTU2MDA5NDU0MTg4ODA1ODM2NzE5MDciLCIxIl0sInBpX2IiOltbIjE2MjY0ODM5OTg1NTE5MTY1MzM1MTQzMjk5MTczMzI0NjU1NDgwNjI1MjE2NDUyNDkzNzQzOTE0NDI3MTA0MDUzMDU0OTA2NzYwMTk1IiwiMjAyNTY3NzY4MTY5MjEzMjUzOTkwMjI4NjYyMTkyNzAzNTc5Mzc4Mzc4ODk2NjA2MDQ0OTU5ODk3Njg3MjczMjUzNzY1ODIyMzM0NTgiXSxbIjY2MzU2NzY4NzMxOTU4OTkyMDUxMDM3MjYwNDE4ODIyNTkxNjc3MDgwNDQwNzc4MzQ2MzM5NzA5ODIxMzc4NDIxNjkyNjg0OTU1OTYiLCI0ODc4ODc2OTQ4MTE1NTQwMDQ4NTE5MzE2MjI2NzE0OTQ3NDMxMDQ5ODc1Njg1MTcxNTg1Mzk3OTE4NDMzNjEwNzI5NTExMjcyNzM1Il0sWyIxIiwiMCJdXSwicGlfYyI6WyIxMTAzOTE3OTgzNjgwNTQyNTU5NzQzNDYzOTYxMzM0MzExOTI3MzAyNjk4NTc3ODQyMDIxMzEwNDEzMzMzNTU5NTEzOTI2Mzg2MjAzNyIsIjEwMjM4MzE5NjY1MjY0MDc0ODQwNTY3MTc5Njg1MjYzOTA1NDY5NTE4NDgzMTk4ODkzMzUwNDMxMDMwNjE3NTMwODY3NjA1NDczNjEzIiwiMSJdfRIkYUhSMGNITTZMeTloWTJOdmRXNTBjeTVuYjI5bmJHVXVZMjl0GmZleUpoYkdjaU9pSlNVekkxTmlJc0ltdHBaQ0k2SWpabU9UYzNOMkUyT0RVNU1EYzNPVGhsWmpjNU5EQTJNbU13TUdJMk5XUTJObU15TkRCaU1XSWlMQ0owZVhBaU9pSktWMVFpZlEiTTE1MDM1MTYxNTYwMTU5OTcxNjMzODAwOTgzNjE5OTMxNDk4Njk2MTUyNjMzNDI2NzY4MDE2OTY2MDU3NzcwNjQzMjYyMDIyMDk2MDczEISIBxJkCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEC620v5dX6qHW1pxNtJNnq/I/T5a7hOOaJx/I9MtIFFcwSBAoCCAESEgoMCgRjb255EgQyMDAwEMCaDBpA6pSglqOdkufji1fUl1NBOBCimsduiA4GD/bd3OXFlBoi35CE2qNfdclGgF8ZO7WchVsnx1PuYnJalNq/RJw8RA=="
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
