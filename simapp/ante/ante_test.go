package ante_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/simapp"
	appante "github.com/Finschia/finschia-sdk/simapp/ante"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/auth/ante"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	minttypes "github.com/Finschia/finschia-sdk/x/mint/types"
)

type fixture struct {
	app       *simapp.SimApp
	ctx       sdk.Context
	clientCtx client.Context
	txBuilder client.TxBuilder
}

func initFixture(t *testing.T) *fixture {
	t.Helper()
	const isCheckTx = false
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	// Set up TxConfig
	encodingConfig := simapp.MakeTestEncodingConfig()
	encodingConfig.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	testdata.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	clientCtx := client.Context{}.WithTxConfig(encodingConfig.TxConfig)

	return &fixture{
		app:       app,
		ctx:       ctx,
		clientCtx: clientCtx,
		txBuilder: clientCtx.TxConfig.NewTxBuilder(),
	}
}

func (f *fixture) CreateTestAccounts(numAcc int) ([]authtypes.AccountI, error) {
	var accounts []authtypes.AccountI

	for i := 0; i < numAcc; i++ {
		_, _, addr := testdata.KeyTestPubAddr()
		acc, err := f.addAccount(addr, i)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (f *fixture) addAccount(accAddr sdk.AccAddress, accNum int) (authtypes.AccountI, error) {
	acc := f.app.AccountKeeper.NewAccountWithAddress(f.ctx, accAddr)
	if err := acc.SetAccountNumber(uint64(accNum)); err != nil {
		return nil, err
	}

	f.app.AccountKeeper.SetAccount(f.ctx, acc)
	someCoins := sdk.Coins{sdk.NewInt64Coin("cony", 10000000)}
	if err := f.app.BankKeeper.MintCoins(f.ctx, minttypes.ModuleName, someCoins); err != nil {
		return nil, err
	}

	if err := f.app.BankKeeper.SendCoinsFromModuleToAccount(f.ctx, minttypes.ModuleName, accAddr, someCoins); err != nil {
		return nil, err
	}

	return acc, nil
}

func (f *fixture) AddTestAccounts(addrs []string) ([]authtypes.AccountI, error) {
	var accounts []authtypes.AccountI

	for i, addrStr := range addrs {
		addr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			return nil, err
		}

		acc, err := f.addAccount(addr, i)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func TestAnteHandler(t *testing.T) {
	f := initFixture(t)
	_, err := f.AddTestAccounts([]string{"link1scmaqcayll6gqwuhl24vtvxe6484t35hqw2qn3ewfqsy704986hqcf5hdw"})
	require.NoError(t, err)

	const sampleTxBase64 = "CqcDCqQDCiUvZmluc2NoaWEuemthdXRoLnYxYmV0YTEuTXNnRXhlY3V0aW9uEvoCCo4BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm4KP2xpbmsxZzd1ZDYzZXFsbGo3emo0cTdma2NhNWg3czIyM2o3OHR5dnIwZTJjeHV3NHF5eWFhZjN1c2E2NGRxYxIrbGluazEwMDh3ZW5ncjI4ejVxdWF0MmR6cnBydDloOGV1YXY0aGVyZnl1bRLmAQrfAQoAEiRhSFIwY0hNNkx5OWhZMk52ZFc1MGN5NW5iMjluYkdVdVkyOXQaZmV5SmhiR2NpT2lKU1V6STFOaUlzSW10cFpDSTZJalptT1RjM04yRTJPRFU1TURjM09UaGxaamM1TkRBMk1tTXdNR0kyTldRMk5tTXlOREJpTVdJaUxDSjBlWEFpT2lKS1YxUWlmUSJNMTUwMzUxNjE1NjAxNTk5NzE2MzM4MDA5ODM2MTk5MzE0OTg2OTYxNTI2MzM0MjY3NjgwMTY5NjYwNTc3NzA2NDMyNjIwMjIwOTYwNzMQsvUDEmQKTgpGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQKyVtAnnf04mVywH48pAVe6XiIYG8FAAUYUHj4rYpY8qxIECgIIARISCgwKBGNvbnkSBDIwMDAQwJoMGkBz80Odj/WNAR7enMKapTpn+uHqak/ZLD+zpC7CylCJjGd9ThHLAZteIB3W85ZuZfl5S3c37De7j7Z99p7M7sdv"
	sampleTxBytes, err := base64.StdEncoding.DecodeString(sampleTxBase64)
	require.NoError(t, err)

	tx, err := f.clientCtx.TxConfig.TxDecoder()(sampleTxBytes)
	require.NoError(t, err)

	anteHandler, err := appante.NewAnteHandler(
		appante.HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   f.app.AccountKeeper,
				BankKeeper:      f.app.BankKeeper,
				FeegrantKeeper:  f.app.FeeGrantKeeper,
				SignModeHandler: f.clientCtx.TxConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			ZKAuthKeeper: f.app.ZKAuthKeeper,
		},
	)

	_, err = anteHandler(f.ctx, tx, false)
	require.NoError(t, err)
}
