package tx_test

import (
	gocontext "context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/tx"
	"github.com/Finschia/finschia-sdk/crypto/hd"
	"github.com/Finschia/finschia-sdk/crypto/keyring"
	cryptotypes "github.com/Finschia/finschia-sdk/crypto/types"
	"github.com/Finschia/finschia-sdk/simapp"
	"github.com/Finschia/finschia-sdk/testutil/network"
	sdk "github.com/Finschia/finschia-sdk/types"
	txtypes "github.com/Finschia/finschia-sdk/types/tx"
	signingtypes "github.com/Finschia/finschia-sdk/types/tx/signing"
	"github.com/Finschia/finschia-sdk/x/auth/signing"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
)

func NewTestTxConfig() client.TxConfig {
	cfg := simapp.MakeTestEncodingConfig()
	return cfg.TxConfig
}

// mockContext is a mock client.Context to return abitrary simulation response, used to
// unit test CalculateGas.
type mockContext struct {
	gasUsed uint64
	wantErr bool
}

func (m mockContext) Invoke(grpcCtx gocontext.Context, method string, req, reply interface{}, opts ...grpc.CallOption) (err error) {
	if m.wantErr {
		return fmt.Errorf("mock err")
	}

	*(reply.(*txtypes.SimulateResponse)) = txtypes.SimulateResponse{
		GasInfo: &sdk.GasInfo{GasUsed: m.gasUsed, GasWanted: m.gasUsed},
		Result:  &sdk.Result{Data: []byte("tx data"), Log: "log"},
	}

	return nil
}

func (mockContext) NewStream(gocontext.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error) {
	panic("not implemented")
}

func TestCalculateGas(t *testing.T) {
	type args struct {
		mockGasUsed uint64
		mockWantErr bool
		adjustment  float64
	}

	testCases := []struct {
		name         string
		args         args
		wantEstimate uint64
		wantAdjusted uint64
		expPass      bool
	}{
		{"error", args{0, true, 1.2}, 0, 0, false},
		{"adjusted gas", args{10, false, 1.2}, 10, 12, true},
	}

	for _, tc := range testCases {
		stc := tc
		txCfg := NewTestTxConfig()

		txf := tx.Factory{}.
			WithChainID("test-chain").
			WithTxConfig(txCfg).WithSignMode(txCfg.SignModeHandler().DefaultMode())

		t.Run(stc.name, func(t *testing.T) {
			mockClientCtx := mockContext{
				gasUsed: tc.args.mockGasUsed,
				wantErr: tc.args.mockWantErr,
			}
			simRes, gotAdjusted, err := tx.CalculateGas(mockClientCtx, txf.WithGasAdjustment(stc.args.adjustment))
			if stc.expPass {
				require.NoError(t, err)
				require.Equal(t, simRes.GasInfo.GasUsed, stc.wantEstimate)
				require.Equal(t, gotAdjusted, stc.wantAdjusted)
				require.NotNil(t, simRes.Result)
			} else {
				require.Error(t, err)
				require.Nil(t, simRes)
			}
		})
	}
}

func TestBuildSimTx(t *testing.T) {
	txCfg := NewTestTxConfig()

	msg := banktypes.NewMsgSend(sdk.AccAddress("from"), sdk.AccAddress("to"), nil)
	kb, err := keyring.New(t.Name(), "test", t.TempDir(), nil)
	require.NoError(t, err)

	txf := tx.Factory{}.
		WithTxConfig(txCfg).
		WithAccountNumber(50).
		WithSequence(23).
		WithFees("50stake").
		WithMemo("memo").
		WithChainID("test-chain").
		WithSignMode(txCfg.SignModeHandler().DefaultMode()).
		WithKeybase(kb)

	// empty keybase list
	bz, err := tx.BuildSimTx(txf, msg)
	require.Error(t, err)
	require.Nil(t, bz)

	// with keybase list
	path := hd.CreateHDPath(118, 0, 0).String()
	_, _, err = kb.NewMnemonic("test_key1", keyring.English, path, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	require.NoError(t, err)
	txf = txf.WithKeybase(kb)
	bz, err = tx.BuildSimTx(txf, msg)
	require.NoError(t, err)
	require.NotNil(t, bz)

	// no ChainID
	txf = txf.WithChainID("")
	bz, err = tx.BuildSimTx(txf, msg)
	require.Error(t, err)
	require.Nil(t, bz)
}

func TestBuildUnsignedTx(t *testing.T) {
	kb, err := keyring.New(t.Name(), "test", t.TempDir(), nil)
	require.NoError(t, err)

	path := hd.CreateHDPath(118, 0, 0).String()

	_, _, err = kb.NewMnemonic("test_key1", keyring.English, path, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	require.NoError(t, err)

	txf := tx.Factory{}.
		WithTxConfig(NewTestTxConfig()).
		WithAccountNumber(50).
		WithSequence(23).
		WithGasPrices("50stake,50cony").
		WithMemo("memo").
		WithChainID("test-chain")

	msg := banktypes.NewMsgSend(sdk.AccAddress("from"), sdk.AccAddress("to"), nil)
	txBuiler, err := tx.BuildUnsignedTx(txf, msg)
	require.NoError(t, err)
	require.NotNil(t, txBuiler)

	sigs, err := txBuiler.GetTx().(signing.SigVerifiableTx).GetSignaturesV2()
	require.NoError(t, err)
	require.Empty(t, sigs)

	// no ChainID
	txf = txf.WithChainID("")
	txBuiler, err = tx.BuildUnsignedTx(txf, msg)
	require.Nil(t, txBuiler)
	require.Error(t, err)

	// both fees and gas prices
	txf = txf.
		WithChainID("test-chain").
		WithFees("50stake")
	txBuiler, err = tx.BuildUnsignedTx(txf, msg)
	require.Nil(t, txBuiler)
	require.Error(t, err)
}

func TestPrintUnsignedTx(t *testing.T) {
	txConfig := NewTestTxConfig()
	txf := tx.Factory{}.
		WithTxConfig(txConfig).
		WithAccountNumber(50).
		WithSequence(23).
		WithFees("50stake").
		WithMemo("memo").
		WithChainID("test-chain")

	msg := banktypes.NewMsgSend(sdk.AccAddress("from"), sdk.AccAddress("to"), nil)
	clientCtx := client.Context{}.
		WithTxConfig(txConfig)
	err := txf.PrintUnsignedTx(clientCtx, msg)
	require.NoError(t, err)

	// no ChainID
	txf = txf.WithChainID("")
	err = txf.PrintUnsignedTx(clientCtx, msg)
	require.Error(t, err)

	// SimulateAndExecute
	// failed at CaculateGas
	txf = txf.WithSimulateAndExecute(true)
	err = txf.PrintUnsignedTx(clientCtx, msg)
	require.Error(t, err)

	// Offline
	clientCtx = clientCtx.WithOffline(true)
	err = txf.PrintUnsignedTx(clientCtx, msg)
	require.Error(t, err)
}

func TestSign(t *testing.T) {
	requireT := require.New(t)
	path := hd.CreateHDPath(118, 0, 0).String()
	kr, err := keyring.New(t.Name(), "test", t.TempDir(), nil)
	requireT.NoError(err)

	from1 := "test_key1"
	from2 := "test_key2"

	// create a new key using a mnemonic generator and test if we can reuse seed to recreate that account
	_, seed, err := kr.NewMnemonic(from1, keyring.English, path, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	requireT.NoError(err)
	requireT.NoError(kr.Delete(from1))
	info1, _, err := kr.NewMnemonic(from1, keyring.English, path, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	requireT.NoError(err)

	info2, err := kr.NewAccount(from2, seed, "", path, hd.Secp256k1)
	requireT.NoError(err)

	pubKey1 := info1.GetPubKey()
	pubKey2 := info2.GetPubKey()
	requireT.NotEqual(pubKey1.Bytes(), pubKey2.Bytes())
	t.Log("Pub keys:", pubKey1, pubKey2)

	txfNoKeybase := tx.Factory{}.
		WithTxConfig(NewTestTxConfig()).
		WithAccountNumber(50).
		WithSequence(23).
		WithFees("50stake").
		WithMemo("memo").
		WithChainID("test-chain")
	txfDirect := txfNoKeybase.
		WithKeybase(kr).
		WithSignMode(signingtypes.SignMode_SIGN_MODE_DIRECT)
	txfAmino := txfDirect.
		WithSignMode(signingtypes.SignMode_SIGN_MODE_LEGACY_AMINO_JSON)
	msg1 := banktypes.NewMsgSend(info1.GetAddress(), sdk.AccAddress("to"), nil)
	msg2 := banktypes.NewMsgSend(info2.GetAddress(), sdk.AccAddress("to"), nil)

	txb, err := tx.BuildUnsignedTx(txfNoKeybase, msg1, msg2)
	requireT.NoError(err)
	txb2, err := tx.BuildUnsignedTx(txfNoKeybase, msg1, msg2)
	requireT.NoError(err)
	txbSimple, err := tx.BuildUnsignedTx(txfNoKeybase, msg2)
	requireT.NoError(err)

	testCases := []struct {
		name         string
		txf          tx.Factory
		txb          client.TxBuilder
		from         string
		overwrite    bool
		expectedPKs  []cryptotypes.PubKey
		matchingSigs []int // if not nil, check matching signature against old ones.
	}{
		{
			"should fail if txf without keyring",
			txfNoKeybase, txb, from1, true, nil, nil,
		},
		{
			"should fail for non existing key",
			txfAmino, txb, "unknown", true, nil, nil,
		},
		{
			"amino: should succeed with keyring",
			txfAmino, txbSimple, from1, true,
			[]cryptotypes.PubKey{pubKey1},
			nil,
		},
		{
			"direct: should succeed with keyring",
			txfDirect, txbSimple, from1, true,
			[]cryptotypes.PubKey{pubKey1},
			nil,
		},

		/**** test double sign Amino mode ****/
		{
			"amino: should sign multi-signers tx",
			txfAmino, txb, from1, true,
			[]cryptotypes.PubKey{pubKey1},
			nil,
		},
		{
			"amino: should append a second signature and not overwrite",
			txfAmino, txb, from2, false,
			[]cryptotypes.PubKey{pubKey1, pubKey2},
			[]int{0, 0},
		},
		{
			"amino: should overwrite a signature",
			txfAmino, txb, from2, true,
			[]cryptotypes.PubKey{pubKey2},
			[]int{1, 0},
		},

		/**** test double sign Direct mode
		  signing transaction with 2 or more DIRECT signers should fail in DIRECT mode ****/
		{
			"direct: should append a DIRECT signature with existing AMINO",
			// txb already has 1 AMINO signature
			txfDirect, txb, from1, false,
			[]cryptotypes.PubKey{pubKey2, pubKey1},
			nil,
		},
		{
			"direct: should add single DIRECT sig in multi-signers tx",
			txfDirect, txb2, from1, false,
			[]cryptotypes.PubKey{pubKey1},
			nil,
		},
		{
			"direct: should fail to append 2nd DIRECT sig in multi-signers tx",
			txfDirect, txb2, from2, false,
			[]cryptotypes.PubKey{},
			nil,
		},
		{
			"amino: should append 2nd AMINO sig in multi-signers tx with 1 DIRECT sig",
			// txb2 already has 1 DIRECT signature
			txfAmino, txb2, from2, false,
			[]cryptotypes.PubKey{},
			nil,
		},
		{
			"direct: should overwrite multi-signers tx with DIRECT sig",
			txfDirect, txb2, from1, true,
			[]cryptotypes.PubKey{pubKey1},
			nil,
		},
	}
	var prevSigs []signingtypes.SignatureV2
	for _, tc := range testCases {
		t.Run(tc.name, func(_ *testing.T) {
			err = tx.Sign(tc.txf, tc.from, tc.txb, tc.overwrite)
			if len(tc.expectedPKs) == 0 {
				requireT.Error(err)
			} else {
				requireT.NoError(err)
				sigs := testSigners(requireT, tc.txb.GetTx(), tc.expectedPKs...)
				if tc.matchingSigs != nil {
					requireT.Equal(prevSigs[tc.matchingSigs[0]], sigs[tc.matchingSigs[1]])
				}
				prevSigs = sigs
			}
		})
	}
}

func testSigners(require *require.Assertions, tr signing.Tx, pks ...cryptotypes.PubKey) []signingtypes.SignatureV2 {
	sigs, err := tr.GetSignaturesV2()
	require.Len(sigs, len(pks))
	require.NoError(err)
	require.Len(sigs, len(pks))
	for i := range pks {
		require.True(sigs[i].PubKey.Equals(pks[i]), "Signature is signed with a wrong pubkey. Got: %s, expected: %s", sigs[i].PubKey, pks[i])
	}
	return sigs
}

func TestPrepare(t *testing.T) {
	txf := tx.Factory{}.
		WithAccountRetriever(authtypes.AccountRetriever{})

	cfg := network.DefaultConfig()
	cfg.NumValidators = 1

	network := network.New(t, cfg)
	defer network.Cleanup()

	_, err := network.WaitForHeight(3)
	require.NoError(t, err)

	val := network.Validators[0]
	clientCtx := val.ClientCtx.
		WithHeight(2).
		WithFromAddress(val.Address)

	_, err = txf.Prepare(clientCtx)
	require.NoError(t, err)
}
