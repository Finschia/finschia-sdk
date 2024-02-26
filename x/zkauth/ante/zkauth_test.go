package ante_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/client"
	codectypes "github.com/Finschia/finschia-sdk/codec/types"
	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	types2 "github.com/Finschia/finschia-sdk/crypto/types"
	"github.com/Finschia/finschia-sdk/simapp"
	"github.com/Finschia/finschia-sdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/tx/signing"
	xauthsigning "github.com/Finschia/finschia-sdk/x/auth/signing"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	banktype "github.com/Finschia/finschia-sdk/x/bank/types"
	minttypes "github.com/Finschia/finschia-sdk/x/mint/types"
	"github.com/Finschia/finschia-sdk/x/zkauth/ante"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
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
		acc := f.app.AccountKeeper.NewAccountWithAddress(f.ctx, addr)
		if err := acc.SetAccountNumber(uint64(i)); err != nil {
			return nil, err
		}

		f.app.AccountKeeper.SetAccount(f.ctx, acc)
		someCoins := sdk.Coins{sdk.NewInt64Coin("cony", 10000000)}
		if err := f.app.BankKeeper.MintCoins(f.ctx, minttypes.ModuleName, someCoins); err != nil {
			return nil, err
		}

		if err := f.app.BankKeeper.SendCoinsFromModuleToAccount(f.ctx, minttypes.ModuleName, addr, someCoins); err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (f *fixture) CreateTestTx(pubs []types2.PubKey, accSeqs []uint64) (xauthsigning.Tx, error) {
	var sigsV2 []signing.SignatureV2
	for i, pub := range pubs {
		sigV2 := signing.SignatureV2{
			PubKey: pub,
			Data: &signing.SingleSignatureData{
				SignMode:  f.clientCtx.TxConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := f.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	return f.txBuilder.GetTx(), nil
}

func TestNewZKAuthMsgDecorator(t *testing.T) {
	f := initFixture(t)
	decorator := ante.NewZKAuthMsgDecorator(f.app.ZKAuthKeeper)
	accounts, err := f.CreateTestAccounts(2)
	require.NoError(t, err)

	// bank msg
	subMsg := &banktype.MsgSend{
		FromAddress: accounts[0].String(),
		ToAddress:   accounts[1].String(),
		Amount:      sdk.Coins{sdk.NewInt64Coin("cony", 1000000)},
	}
	any, err := codectypes.NewAnyWithValue(subMsg)
	const proofStr = "{\n \"pi_a\": [\n  \"7575287679446209007446416020137456670042570578978230730578011103770415897062\",\n  \"20469978368515629364541212704109752583692706286549284712208570249653184893207\",\n  \"1\"\n ],\n \"pi_b\": [\n  [\n   \"4001119070037193619600086014535210556571209449080681376392853276923728808564\",\n   \"18475391841797083641468254159150812922259839776046448499150732610021959794558\"\n  ],\n  [\n   \"19781252109528278034156073207688818205850783935629584279449144780221040670063\",\n   \"5873714313814830719712095806732872482213125567325442209795797618441438990229\"\n  ],\n  [\n   \"1\",\n   \"0\"\n  ]\n ],\n \"pi_c\": [\n  \"18920522434978516095250248740518039198650690968720755259416280639852277665022\",\n  \"1945774583580804632084048753815901730674007769630810705050114062476636502591\",\n  \"1\"\n ],\n \"protocol\": \"groth16\",\n \"curve\": \"bn128\"\n}"
	msg := &types.MsgExecution{
		Msgs: []*codectypes.Any{any},
		ZkAuthSignature: types.ZKAuthSignature{
			ZkAuthInputs: &types.ZKAuthInputs{
				ProofPoints:  []byte(proofStr),
				IssF:         "aHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29t",
				HeaderBase64: "eyJhbGciOiJSUzI1NiIsImtpZCI6IjU1YzE4OGE4MzU0NmZjMTg4ZTUxNTc2YmE3MjgzNmUwNjAwZThiNzMiLCJ0eXAiOiJKV1QifQ",
				AddressSeed:  "15035161560159971633800983619931498696152633426768016966057770643262022096073",
			},
			MaxBlockHeight: 32754,
		},
	}
	err = f.txBuilder.SetMsgs(msg)
	require.NoError(t, err)

	ephPubKey, ok := new(big.Int).SetString("18948426102457371978524559226152399917062673825697601263047735920285791872240", 10)
	require.True(t, ok)
	pub := secp256k1.PubKey{Key: ephPubKey.Bytes()}
	tx, err := f.CreateTestTx([]types2.PubKey{&pub}, []uint64{uint64(0)})
	require.NoError(t, err)

	_, err = decorator.AnteHandle(f.ctx, tx, false, func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
		return ctx, nil
	})
	require.NoError(t, err)
}
