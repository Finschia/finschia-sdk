package ante_test

//import (
//	"encoding/base64"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//
//	appante "github.com/Finschia/finschia-sdk/simapp/ante"
//	"github.com/Finschia/finschia-sdk/x/auth/ante"
//	"github.com/Finschia/finschia-sdk/x/zkauth/testutil"
//)
//
//func TestAnteHandler(t *testing.T) {
//	k := testutil.ZkAuthKeeper(t)
//	_, err := k.AddTestAccounts([]string{"link19qvkxwmln5kf0z59ecw744ue2gzsndlwcuz9uq2kav4evuerjlysxduwzj"})
//	require.NoError(t, err)
//
//	const sampleTxBase64 = "CtEICs4ICiUvZmluc2NoaWEuemthdXRoLnYxYmV0YTEuTXNnRXhlY3V0aW9uEqQICo4BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm4KP2xpbmsxZzd1ZDYzZXFsbGo3emo0cTdma2NhNWg3czIyM2o3OHR5dnIwZTJjeHV3NHF5eWFhZjN1c2E2NGRxYxIrbGluazEwMDh3ZW5ncjI4ejVxdWF0MmR6cnBydDloOGV1YXY0aGVyZnl1bRKQBwqJBwqpBXsicGlfYSI6WyIyMjQ0NTE0OTk0ODU0NDU0Mjk1MDA5ODUwODA5ODc2ODk5MzMwODMwMDkzNTc4NDE4NDgyOTU3OTM1NjgxMTc2OTkyMTE0NjA1NTk5IiwiMTAwMTg2MTAwODQ3OTAyODQ2MjIwODE1MDk0NDA4NDk4MDI4ODcxNDY0MTgwOTU3NjY5NDQ1ODAyNDQ1NDk2MDU3MDAzMzczOTQ0MjYiLCIxIl0sInBpX2IiOltbIjQ0MjcwMjAyMzU0MzY1NTcxMjY4NDg4MzM1MDI3NDAwMzI0NTI4ODc1ODE5MDcxMDc1NzcxMDE5NDg0MDQ2MDE3NjE3OTQ3MzMyNzAiLCIxNDEwMzE4ODM4OTI4NTMxMjAzMTY1OTE1MTQwOTYyNjEwNzc1OTE0OTcyNDYwNzUxMDk1ODAwMzQ0MTE3NjcyMjQ1NTcxODIxNjI5Il0sWyIzNjQ0ODkwMzQ4MzQzODM1NjQwNTU1NzUwNDE2NjA0OTE5MzQyNDA4ODQzNjA1NTkzMjY4MzQ0MjkxMDU2NzczNTYxMDMxMDA3NTEzIiwiMjU0MDM5OTcxNzQ2MTE1OTI4NDM2Njg2MDY3NjgzNTQwOTMwNTMwOTIxMDYyMTU3OTYxNTI3MDUwMTk4Mzc1MTY1NzU1NTA1NTY2OSJdLFsiMSIsIjAiXV0sInBpX2MiOlsiNjg4MDM4MTAxNDMzOTQyMDc5NTI4MDM2ODEwNjUyNDk3NDMwNzg1MDg1NjA4NDkyNDQxOTAxNDM0MDY1NDY1Mzg3MzU4NTQyNDA4NiIsIjIwMTU3MTk1MDk4MTA0ODI4ODIxNDM0MTU5Mzg3NzU5MjM2NzMwNzYyMzM0NTA4OTMxNDQ1MTkzNzM4MjMzMjc3MjczNjQ5MDQ3IiwiMSJdfRIkYUhSMGNITTZMeTloWTJOdmRXNTBjeTVuYjI5bmJHVXVZMjl0GmZleUpoYkdjaU9pSlNVekkxTmlJc0ltdHBaQ0k2SWpabU9UYzNOMkUyT0RVNU1EYzNPVGhsWmpjNU5EQTJNbU13TUdJMk5XUTJObU15TkRCaU1XSWlMQ0owZVhBaU9pSktWMVFpZlEiTTE1MDM1MTYxNTYwMTU5OTcxNjMzODAwOTgzNjE5OTMxNDk4Njk2MTUyNjMzNDI2NzY4MDE2OTY2MDU3NzcwNjQzMjYyMDIyMDk2MDczEMLSBhJkCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECslbQJ539OJlcsB+PKQFXul4iGBvBQAFGFB4+K2KWPKsSBAoCCAESEgoMCgRjb255EgQyMDAwEMCaDBpA33tI8LSxX7Em24By+y9IKLjUpfAYAksIaODWmCUKhQEyxo8n9X15kxpEi0wDFxsqYSJ/49Lbk9BO/N95RcKAQQ=="
//	sampleTxBytes, err := base64.StdEncoding.DecodeString(sampleTxBase64)
//	require.NoError(t, err)
//
//	tx, err := k.ClientCtx.TxConfig.TxDecoder()(sampleTxBytes)
//	require.NoError(t, err)
//
//	anteHandler, err := appante.NewAnteHandler(
//		appante.HandlerOptions{
//			HandlerOptions: ante.HandlerOptions{
//				AccountKeeper:   k.Simapp.AccountKeeper,
//				BankKeeper:      k.Simapp.BankKeeper,
//				FeegrantKeeper:  k.Simapp.FeeGrantKeeper,
//				SignModeHandler: k.ClientCtx.TxConfig.SignModeHandler(),
//				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
//			},
//			ZKAuthKeeper: k.ZKAuthKeeper,
//		},
//	)
//
//	_, err = anteHandler(k.Ctx, tx, false)
//	require.NoError(t, err)
//}
