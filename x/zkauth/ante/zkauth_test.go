package ante_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	codectypes "github.com/Finschia/finschia-sdk/codec/types"
	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	types2 "github.com/Finschia/finschia-sdk/crypto/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	authante "github.com/Finschia/finschia-sdk/x/auth/ante"
	banktype "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/zkauth/ante"
	"github.com/Finschia/finschia-sdk/x/zkauth/testutil"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
)

func TestNewDecorators(t *testing.T) {
	f := testutil.ZkAuthKeeper(t)
	decorators := []sdk.AnteDecorator{
		ante.NewZKAuthDeductFeeDecorator(f.AccountKeeper, f.BankKeeper, f.FeeGrantKeeper),
		ante.NewZKAuthSigGasConsumeDecorator(f.AccountKeeper, authante.DefaultSigVerificationGasConsumer),
		ante.NewZKAuthMsgDecorator(f.ZKAuthKeeper, f.AccountKeeper, f.SignModeHandler),
	}
	accounts, err := f.CreateTestAccounts(2)
	require.NoError(t, err)
	zkauthAddress, err := f.AddTestAccounts([]string{"link1g7ud63eqllj7zj4q7fkca5h7s223j78tyvr0e2cxuw4qyyaaf3usa64dqc"})
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
				IssBase64:    "aHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29t",
				HeaderBase64: "eyJhbGciOiJSUzI1NiIsImtpZCI6IjU1YzE4OGE4MzU0NmZjMTg4ZTUxNTc2YmE3MjgzNmUwNjAwZThiNzMiLCJ0eXAiOiJKV1QifQ",
				AddressSeed:  "15035161560159971633800983619931498696152633426768016966057770643262022096073",
			},
			MaxBlockHeight: 32754,
		},
	}
	err = f.TxBuilder.SetMsgs(msg)
	require.NoError(t, err)

	f.TxBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin("cony", sdk.NewInt(1000))))
	ephPubKey, ok := new(big.Int).SetString("18948426102457371978524559226152399917062673825697601263047735920285791872240", 10)
	require.True(t, ok)
	pub := secp256k1.PubKey{Key: ephPubKey.Bytes()}
	tx, err := f.CreateTestTx([]types2.PubKey{&pub}, []uint64{uint64(0)})
	require.NoError(t, err)

	for _, decorator := range decorators {
		_, err = decorator.AnteHandle(f.Ctx, tx, false, func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, err error) {
			return ctx, nil
		})
		require.NoError(t, err)
	}

	balance, err := f.BankKeeper.Balance(sdk.WrapSDKContext(f.Ctx), &banktype.QueryBalanceRequest{zkauthAddress[0].GetAddress().String(), "cony"})
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(9999000), balance.Balance.Amount)
}
