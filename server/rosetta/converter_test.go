package rosetta_test

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/Finschia/finschia-sdk/testutil/testdata"
	"github.com/Finschia/finschia-sdk/types/tx/signing"
	authtx "github.com/Finschia/finschia-sdk/x/auth/tx"

	rosettatypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/stretchr/testify/suite"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/Finschia/finschia-rdk/server/rosetta"
	crgerrs "github.com/Finschia/finschia-rdk/server/rosetta/lib/errors"
	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/codec"
	codectypes "github.com/Finschia/finschia-sdk/codec/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	authsigning "github.com/Finschia/finschia-sdk/x/auth/signing"
	bank "github.com/Finschia/finschia-sdk/x/bank/types"
)

type ConverterTestSuite struct {
	suite.Suite

	c               rosetta.Converter
	unsignedTxBytes []byte
	unsignedTx      authsigning.Tx

	ir     codectypes.InterfaceRegistry
	cdc    *codec.ProtoCodec
	txConf client.TxConfig
}

// generateMsgSend generate sample unsignedTxHex and pubKeyHex
func generateMsgSend() (unsignedTxHex []byte, pubKeyHex []byte) {
	cdc, _ := rosetta.MakeCodec()
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)

	_, fromPk, fromAddr := testdata.KeyTestPubAddr()
	_, _, toAddr := testdata.KeyTestPubAddr()

	sendMsg := bank.MsgSend{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("stake", 16)),
	}

	txBuilder := txConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(&sendMsg)
	if err != nil {
		return nil, nil
	}
	txBuilder.SetGasLimit(200000)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("stake", 1)))

	sigData := signing.SingleSignatureData{
		SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
		Signature: nil,
	}
	sig := signing.SignatureV2{
		PubKey:   fromPk,
		Data:     &sigData,
		Sequence: 1,
	}
	err = txBuilder.SetSignatures(sig)
	if err != nil {
		return nil, nil
	}

	stdTx := txBuilder.GetTx()
	unsignedTxHex, err = txConfig.TxEncoder()(stdTx)
	if err != nil {
		return nil, nil
	}

	return unsignedTxHex, fromPk.Bytes()
}

func (s *ConverterTestSuite) SetupTest() {
	// create an unsigned tx
	const unsignedTxHex = "0a8a010a87010a1c2f636f736d6f732e62616e6b2e763162657461312e4d736753656e6412670a2b6c696e6b3136773064766c61716e34727779706739757a36713878643778343434687568756d636370756e122b6c696e6b316c33757538657364636c6a3876706a72757737357535666a34773479746475396e6c6b6538721a0b0a057374616b651202313612640a500a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a210377365794209eab396f74316bb32ecb507c0e3788c14edf164f96b25cc4ef624112040a020801180112100a0a0a057374616b6512013110c09a0c1a00"
	unsignedTxBytes, err := hex.DecodeString(unsignedTxHex)
	s.Require().NoError(err)
	s.unsignedTxBytes = unsignedTxBytes
	// instantiate converter
	cdc, ir := rosetta.MakeCodec()
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	s.c = rosetta.NewConverter(cdc, ir, txConfig)
	// add utils
	s.ir = ir
	s.cdc = cdc
	s.txConf = txConfig
	// add authsigning tx
	sdkTx, err := txConfig.TxDecoder()(unsignedTxBytes)
	s.Require().NoError(err)
	builder, err := txConfig.WrapTxBuilder(sdkTx)
	s.Require().NoError(err)

	s.unsignedTx = builder.GetTx()
}

func (s *ConverterTestSuite) TestFromRosettaOpsToTxSuccess() {
	addr1 := sdk.AccAddress("address1").String()
	addr2 := sdk.AccAddress("address2").String()

	msg1 := &bank.MsgSend{
		FromAddress: addr1,
		ToAddress:   addr2,
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("test", 10)),
	}

	msg2 := &bank.MsgSend{
		FromAddress: addr2,
		ToAddress:   addr1,
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("utxo", 10)),
	}

	ops, err := s.c.ToRosetta().Ops("", msg1)
	s.Require().NoError(err)

	ops2, err := s.c.ToRosetta().Ops("", msg2)
	s.Require().NoError(err)

	ops = append(ops, ops2...)

	tx, err := s.c.ToSDK().UnsignedTx(ops)
	s.Require().NoError(err)

	getMsgs := tx.GetMsgs()

	s.Require().Equal(2, len(getMsgs))

	s.Require().Equal(getMsgs[0], msg1)
	s.Require().Equal(getMsgs[1], msg2)
}

func (s *ConverterTestSuite) TestFromRosettaOpsToTxErrors() {
	s.Run("unrecognized op", func() {
		op := &rosettatypes.Operation{
			Type: "non-existent",
		}

		_, err := s.c.ToSDK().UnsignedTx([]*rosettatypes.Operation{op})

		s.Require().ErrorIs(err, crgerrs.ErrBadArgument)
	})

	s.Run("codec type but not sdk.Msg", func() {
		op := &rosettatypes.Operation{
			Type: "cosmos.crypto.ed25519.PubKey",
		}

		_, err := s.c.ToSDK().UnsignedTx([]*rosettatypes.Operation{op})

		s.Require().ErrorIs(err, crgerrs.ErrBadArgument)
	})
}

func (s *ConverterTestSuite) TestMsgToMetaMetaToMsg() {
	msg := &bank.MsgSend{
		FromAddress: "addr1",
		ToAddress:   "addr2",
		Amount:      sdk.NewCoins(sdk.NewInt64Coin("test", 10)),
	}
	msg.Route()

	meta, err := s.c.ToRosetta().Meta(msg)
	s.Require().NoError(err)

	copyMsg := new(bank.MsgSend)
	err = s.c.ToSDK().Msg(meta, copyMsg)
	s.Require().NoError(err)
	s.Require().Equal(msg, copyMsg)
}

func (s *ConverterTestSuite) TestSignedTx() {
	s.Run("success", func() {
		const payloadsJSON = `[{"hex_bytes":"82ccce81a3e4a7272249f0e25c3037a316ee2acce76eb0c25db00ef6634a4d57303b2420edfdb4c9a635ad8851fe5c7a9379b7bc2baadc7d74f7e76ac97459b5","public_key":{"curve_type":"secp256k1","hex_bytes":"0377365794209eab396f74316bb32ecb507c0e3788c14edf164f96b25cc4ef6241"},"signature_type":"ecdsa","signing_payload":{"account_identifier":{"address":"link16w0dvlaqn4rwypg9uz6q8xd7x444huhumccpun"},"address":"link16w0dvlaqn4rwypg9uz6q8xd7x444huhumccpun","hex_bytes":"ea43c4019ee3c888a7f99acb57513f708bb8915bc84e914cf4ecbd08ab2d9e51","signature_type":"ecdsa"}}]`
		const expectedSignedTxHex = "0a8a010a87010a1c2f636f736d6f732e62616e6b2e763162657461312e4d736753656e6412670a2b6c696e6b3136773064766c61716e34727779706739757a36713878643778343434687568756d636370756e122b6c696e6b316c33757538657364636c6a3876706a72757737357535666a34773479746475396e6c6b6538721a0b0a057374616b651202313612640a500a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a210377365794209eab396f74316bb32ecb507c0e3788c14edf164f96b25cc4ef624112040a02087f180112100a0a0a057374616b6512013110c09a0c1a4082ccce81a3e4a7272249f0e25c3037a316ee2acce76eb0c25db00ef6634a4d57303b2420edfdb4c9a635ad8851fe5c7a9379b7bc2baadc7d74f7e76ac97459b5"

		var payloads []*rosettatypes.Signature
		s.Require().NoError(json.Unmarshal([]byte(payloadsJSON), &payloads))

		signedTx, err := s.c.ToSDK().SignedTx(s.unsignedTxBytes, payloads)
		s.Require().NoError(err)

		signedTxHex := hex.EncodeToString(signedTx)

		s.Require().Equal(signedTxHex, expectedSignedTxHex)
	})

	s.Run("signers data and signing payloads mismatch", func() {
		_, err := s.c.ToSDK().SignedTx(s.unsignedTxBytes, nil)
		s.Require().ErrorIs(err, crgerrs.ErrInvalidTransaction)
	})
}

func (s *ConverterTestSuite) TestOpsAndSigners() {
	s.Run("success", func() {
		addr1 := sdk.AccAddress("address1").String()
		addr2 := sdk.AccAddress("address2").String()

		msg := &bank.MsgSend{
			FromAddress: addr1,
			ToAddress:   addr2,
			Amount:      sdk.NewCoins(sdk.NewInt64Coin("test", 10)),
		}

		builder := s.txConf.NewTxBuilder()
		s.Require().NoError(builder.SetMsgs(msg))

		sdkTx := builder.GetTx()
		txBytes, err := s.txConf.TxEncoder()(sdkTx)
		s.Require().NoError(err)

		ops, signers, err := s.c.ToRosetta().OpsAndSigners(txBytes)
		s.Require().NoError(err)

		s.Require().Equal(len(ops), len(sdkTx.GetMsgs())*len(sdkTx.GetSigners()), "operation number mismatch")

		s.Require().Equal(len(signers), len(sdkTx.GetSigners()), "signers number mismatch")
	})
}

func (s *ConverterTestSuite) TestBeginEndBlockAndHashToTxType() {
	const deliverTxHex = "5229A67AA008B5C5F1A0AEA77D4DEBE146297A30AAEF01777AF10FAD62DD36AB"

	deliverTxBytes, err := hex.DecodeString(deliverTxHex)
	s.Require().NoError(err)

	endBlockTxHex := s.c.ToRosetta().EndBlockTxHash(deliverTxBytes)
	beginBlockTxHex := s.c.ToRosetta().BeginBlockTxHash(deliverTxBytes)

	txType, hash := s.c.ToSDK().HashToTxType(deliverTxBytes)

	s.Require().Equal(rosetta.DeliverTxTx, txType)
	s.Require().Equal(deliverTxBytes, hash, "deliver tx hash should not change")

	endBlockTxBytes, err := hex.DecodeString(endBlockTxHex)
	s.Require().NoError(err)

	txType, hash = s.c.ToSDK().HashToTxType(endBlockTxBytes)

	s.Require().Equal(rosetta.EndBlockTx, txType)
	s.Require().Equal(deliverTxBytes, hash, "end block tx hash should be equal to a block hash")

	beginBlockTxBytes, err := hex.DecodeString(beginBlockTxHex)
	s.Require().NoError(err)

	txType, hash = s.c.ToSDK().HashToTxType(beginBlockTxBytes)

	s.Require().Equal(rosetta.BeginBlockTx, txType)
	s.Require().Equal(deliverTxBytes, hash, "begin block tx hash should be equal to a block hash")

	txType, hash = s.c.ToSDK().HashToTxType([]byte("invalid"))

	s.Require().Equal(rosetta.UnrecognizedTx, txType)
	s.Require().Nil(hash)

	txType, hash = s.c.ToSDK().HashToTxType(append([]byte{0x3}, deliverTxBytes...))
	s.Require().Equal(rosetta.UnrecognizedTx, txType)
	s.Require().Nil(hash)
}

func (s *ConverterTestSuite) TestSigningComponents() {
	s.Run("invalid metadata coins", func() {
		_, _, err := s.c.ToRosetta().SigningComponents(nil, &rosetta.ConstructionMetadata{GasPrice: "invalid"}, nil)
		s.Require().ErrorIs(err, crgerrs.ErrBadArgument)
	})

	s.Run("length signers data does not match signers", func() {
		_, _, err := s.c.ToRosetta().SigningComponents(s.unsignedTx, &rosetta.ConstructionMetadata{GasPrice: "10stake"}, nil)
		s.Require().ErrorIs(err, crgerrs.ErrBadArgument)
	})

	s.Run("length pub keys does not match signers", func() {
		_, _, err := s.c.ToRosetta().SigningComponents(
			s.unsignedTx,
			&rosetta.ConstructionMetadata{GasPrice: "10stake", SignersData: []*rosetta.SignerData{
				{
					AccountNumber: 0,
					Sequence:      0,
				},
			}},
			nil)
		s.Require().ErrorIs(err, crgerrs.ErrBadArgument)
	})

	s.Run("ros pub key is valid but not the one we expect", func() {
		validButUnexpected, err := hex.DecodeString("030da9096a40eb1d6c25f1e26e9cbf8941fc84b8f4dc509c8df5e62a29ab8f2415")
		s.Require().NoError(err)

		_, _, err = s.c.ToRosetta().SigningComponents(
			s.unsignedTx,
			&rosetta.ConstructionMetadata{GasPrice: "10stake", SignersData: []*rosetta.SignerData{
				{
					AccountNumber: 0,
					Sequence:      0,
				},
			}},
			[]*rosettatypes.PublicKey{
				{
					Bytes:     validButUnexpected,
					CurveType: rosettatypes.Secp256k1,
				},
			})
		s.Require().ErrorIs(err, crgerrs.ErrBadArgument)
	})

	s.Run("success", func() {
		expectedPubKey, err := hex.DecodeString("0377365794209eab396f74316bb32ecb507c0e3788c14edf164f96b25cc4ef6241")
		s.Require().NoError(err)

		_, _, err = s.c.ToRosetta().SigningComponents(
			s.unsignedTx,
			&rosetta.ConstructionMetadata{GasPrice: "10stake", SignersData: []*rosetta.SignerData{
				{
					AccountNumber: 0,
					Sequence:      0,
				},
			}},
			[]*rosettatypes.PublicKey{
				{
					Bytes:     expectedPubKey,
					CurveType: rosettatypes.Secp256k1,
				},
			})
		s.Require().NoError(err)
	})
}

func (s *ConverterTestSuite) TestBalanceOps() {
	s.Run("not a balance op", func() {
		notBalanceOp := abci.Event{
			Type: "not-a-balance-op",
		}

		ops := s.c.ToRosetta().BalanceOps("", []abci.Event{notBalanceOp})
		s.Len(ops, 0, "expected no balance ops")
	})

	s.Run("multiple balance ops from 2 multicoins event", func() {
		subBalanceOp := bank.NewCoinSpentEvent(
			sdk.AccAddress("test"),
			sdk.NewCoins(sdk.NewInt64Coin("test", 10), sdk.NewInt64Coin("utxo", 10)),
		)

		addBalanceOp := bank.NewCoinReceivedEvent(
			sdk.AccAddress("test"),
			sdk.NewCoins(sdk.NewInt64Coin("test", 10), sdk.NewInt64Coin("utxo", 10)),
		)

		ops := s.c.ToRosetta().BalanceOps("", []abci.Event{(abci.Event)(subBalanceOp), (abci.Event)(addBalanceOp)})
		s.Len(ops, 4)
	})

	s.Run("spec broken", func() {
		s.Require().Panics(func() {
			specBrokenSub := abci.Event{
				Type: bank.EventTypeCoinSpent,
			}
			_ = s.c.ToRosetta().BalanceOps("", []abci.Event{specBrokenSub})
		})

		s.Require().Panics(func() {
			specBrokenSub := abci.Event{
				Type: bank.EventTypeCoinBurn,
			}
			_ = s.c.ToRosetta().BalanceOps("", []abci.Event{specBrokenSub})
		})

		s.Require().Panics(func() {
			specBrokenSub := abci.Event{
				Type: bank.EventTypeCoinReceived,
			}
			_ = s.c.ToRosetta().BalanceOps("", []abci.Event{specBrokenSub})
		})
	})
}

func TestConverterTestSuite(t *testing.T) {
	suite.Run(t, new(ConverterTestSuite))
}
