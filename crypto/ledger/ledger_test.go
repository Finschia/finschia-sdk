package ledger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-rdk/codec/legacy"
	"github.com/Finschia/finschia-rdk/crypto/hd"
	"github.com/Finschia/finschia-rdk/crypto/types"
	"github.com/Finschia/finschia-rdk/testutil/testdata"
	sdk "github.com/Finschia/finschia-rdk/types"
)

func TestErrorHandling(t *testing.T) {
	// first, try to generate a key, must return an error
	// (no panic)
	path := *hd.NewParams(44, 555, 0, false, 0)
	_, err := NewPrivKeySecp256k1Unsafe(path)
	require.Error(t, err)
}

func TestPublicKeyUnsafe(t *testing.T) {
	path := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	priv, err := NewPrivKeySecp256k1Unsafe(path)
	require.NoError(t, err)
	checkDefaultPubKey(t, priv)
}

func checkDefaultPubKey(t *testing.T, priv types.LedgerPrivKey) {
	require.NotNil(t, priv)
	expectedPkStr := "PubKeySecp256k1{02BCD96CAB102304AC10E90C5A0F29358E3A4A6FB1217B83E5DB657918EA28BEC1}"
	require.Equal(t, "eb5ae9872102bcd96cab102304ac10e90c5a0f29358e3a4a6fb1217b83e5db657918ea28bec1",
		fmt.Sprintf("%x", cdc.Amino.MustMarshalBinaryBare(priv.PubKey())),
		"Is your device using test mnemonic: %s ?", testdata.TestMnemonic)
	require.Equal(t, expectedPkStr, priv.PubKey().String())
	addr := sdk.AccAddress(priv.PubKey().Address()).String()
	require.Equal(t, "link1tdl7n2acgmec0y5nng0q2fahl9khyct3cgsktn",
		addr, "Is your device using test mnemonic: %s ?", testdata.TestMnemonic)
}

func TestPublicKeyUnsafeHDPath(t *testing.T) {
	expectedAnswers := []string{
		"PubKeySecp256k1{02BCD96CAB102304AC10E90C5A0F29358E3A4A6FB1217B83E5DB657918EA28BEC1}",
		"PubKeySecp256k1{025543C96E09392B9B7F416B47BE4520804ED064F9E342062DDE825E0554E24188}",
		"PubKeySecp256k1{02B2DDE2CD78977664A3F3F231EB0B30818AB174D54FFE552205BED3B67DCD852C}",
		"PubKeySecp256k1{020A64E4DF39E6CD90820A29758A80145EE4352B4BC3C1539184B52B58A9313CC1}",
		"PubKeySecp256k1{031E5CA16AFBF9E576B77C176C0CD5B904F16B733668EA37DF16F624B64AD92D2F}",
		"PubKeySecp256k1{026656F9797699F1C7D0217038378D523898257BF38511AB971655F4F775EAB510}",
		"PubKeySecp256k1{027A0D54AC45A3FCB0702754EB08FDCE246336AF1B532B8F5D07FEDA03F4AA5B05}",
		"PubKeySecp256k1{032B1121146F78E24C15710028E3465FCAEC6F0711D2B401923F8D07BDFFFB750D}",
		"PubKeySecp256k1{03ED8BFEBE2B45E4A8F74267E6662B35ACBE1563029A1D8C52217AEF93893BB46F}",
		"PubKeySecp256k1{02D70D97BDC48A1045CA51A5B86BE8208D7BB6431F64BADC424FDFBCFD292B8AD0}",
	}

	const numIters = 10

	privKeys := make([]types.LedgerPrivKey, numIters)

	// Check with device
	for i := uint32(0); i < 10; i++ {
		path := *hd.NewFundraiserParams(0, sdk.CoinType, i)
		t.Logf("Checking keys at %v\n", path)

		priv, err := NewPrivKeySecp256k1Unsafe(path)
		require.NoError(t, err)
		require.NotNil(t, priv)

		// Check other methods
		tmp := priv.(PrivKeyLedgerSecp256k1)
		require.NoError(t, tmp.ValidateKey())
		(&tmp).AssertIsPrivKeyInner()

		// in this test we are chekcking if the generated keys are correct.
		require.Equal(t, expectedAnswers[i], priv.PubKey().String(),
			"Is your device using test mnemonic: %s ?", testdata.TestMnemonic)

		// Store and restore
		serializedPk := priv.Bytes()
		require.NotNil(t, serializedPk)
		require.True(t, len(serializedPk) >= 50)

		privKeys[i] = priv
	}

	// Now check equality
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			require.Equal(t, i == j, privKeys[i].Equals(privKeys[j]))
			require.Equal(t, i == j, privKeys[j].Equals(privKeys[i]))
		}
	}
}

func TestPublicKeySafe(t *testing.T) {
	path := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	priv, addr, err := NewPrivKeySecp256k1(path, "link")

	require.NoError(t, err)
	require.NotNil(t, priv)
	require.Nil(t, ShowAddress(path, priv.PubKey(), sdk.GetConfig().GetBech32AccountAddrPrefix()))
	checkDefaultPubKey(t, priv)

	addr2 := sdk.AccAddress(priv.PubKey().Address()).String()
	require.Equal(t, addr, addr2)
}

func TestPublicKeyHDPath(t *testing.T) {
	expectedPubKeys := []string{
		"PubKeySecp256k1{02BCD96CAB102304AC10E90C5A0F29358E3A4A6FB1217B83E5DB657918EA28BEC1}",
		"PubKeySecp256k1{025543C96E09392B9B7F416B47BE4520804ED064F9E342062DDE825E0554E24188}",
		"PubKeySecp256k1{02B2DDE2CD78977664A3F3F231EB0B30818AB174D54FFE552205BED3B67DCD852C}",
		"PubKeySecp256k1{020A64E4DF39E6CD90820A29758A80145EE4352B4BC3C1539184B52B58A9313CC1}",
		"PubKeySecp256k1{031E5CA16AFBF9E576B77C176C0CD5B904F16B733668EA37DF16F624B64AD92D2F}",
		"PubKeySecp256k1{026656F9797699F1C7D0217038378D523898257BF38511AB971655F4F775EAB510}",
		"PubKeySecp256k1{027A0D54AC45A3FCB0702754EB08FDCE246336AF1B532B8F5D07FEDA03F4AA5B05}",
		"PubKeySecp256k1{032B1121146F78E24C15710028E3465FCAEC6F0711D2B401923F8D07BDFFFB750D}",
		"PubKeySecp256k1{03ED8BFEBE2B45E4A8F74267E6662B35ACBE1563029A1D8C52217AEF93893BB46F}",
		"PubKeySecp256k1{02D70D97BDC48A1045CA51A5B86BE8208D7BB6431F64BADC424FDFBCFD292B8AD0}",
	}

	expectedAddrs := []string{
		"link1tdl7n2acgmec0y5nng0q2fahl9khyct3cgsktn",
		"link1lzmehungm97jh0nme768v9wv28l8jr2dkkv357",
		"link1p0yx9c9q4xsnedlcn24gqfry5dcu6e9xkhv9aj",
		"link1cjjev0yzp90dxhsyxlxkwvrgl6vdw5p7qpxgfu",
		"link1gf58l6wlscadkw2c9dpyd8se4hw37gu0zxx2g6",
		"link1n93v9h6slh2e56lfu6vcy302f6ttxv5tt676ej",
		"link124w2ltcv7wdhc85g07kta7yc4pwqk9qqsfq3c0",
		"link17evhyfgwm70xjy8s3la64ek2x068aqezs24lth",
		"link1clgcsmd9gcu4v9ec0qzqpr932aetwrqkxwswdh",
		"link1hcttwju93d5m39467gjcq63p5kc4fdcn30dgd8",
	}

	const numIters = 10

	privKeys := make([]types.LedgerPrivKey, numIters)

	// Check with device
	for i := 0; i < len(expectedAddrs); i++ {
		path := *hd.NewFundraiserParams(0, sdk.CoinType, uint32(i))
		t.Logf("Checking keys at %s\n", path)

		priv, addr, err := NewPrivKeySecp256k1(path, "link")
		require.NoError(t, err)
		require.NotNil(t, addr)
		require.NotNil(t, priv)

		addr2 := sdk.AccAddress(priv.PubKey().Address()).String()
		require.Equal(t, addr2, addr)
		require.Equal(t,
			expectedAddrs[i], addr,
			"Is your device using test mnemonic: %s ?", testdata.TestMnemonic)

		// Check other methods
		tmp := priv.(PrivKeyLedgerSecp256k1)
		require.NoError(t, tmp.ValidateKey())
		(&tmp).AssertIsPrivKeyInner()

		// in this test we are chekcking if the generated keys are correct and stored in a right path.
		require.Equal(t,
			expectedPubKeys[i], priv.PubKey().String(),
			"Is your device using test mnemonic: %s ?", testdata.TestMnemonic)

		// Store and restore
		serializedPk := priv.Bytes()
		require.NotNil(t, serializedPk)
		require.True(t, len(serializedPk) >= 50)

		privKeys[i] = priv
	}

	// Now check equality
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			require.Equal(t, i == j, privKeys[i].Equals(privKeys[j]))
			require.Equal(t, i == j, privKeys[j].Equals(privKeys[i]))
		}
	}
}

func getFakeTx(accountNumber uint32) []byte {
	tmp := fmt.Sprintf(
		`{"account_number":"%d","chain_id":"1234","fee":{"amount":[{"amount":"150","denom":"atom"}],"gas":"5000"},"memo":"memo","msgs":[[""]],"sequence":"6"}`,
		accountNumber)

	return []byte(tmp)
}

func TestSignaturesHD(t *testing.T) {
	for account := uint32(0); account < 100; account += 30 {
		msg := getFakeTx(account)

		path := *hd.NewFundraiserParams(account, sdk.CoinType, account/5)
		t.Logf("Checking signature at %v    ---   PLEASE REVIEW AND ACCEPT IN THE DEVICE\n", path)

		priv, err := NewPrivKeySecp256k1Unsafe(path)
		require.NoError(t, err)

		pub := priv.PubKey()
		sig, err := priv.Sign(msg)
		require.NoError(t, err)

		valid := pub.VerifySignature(msg, sig)
		require.True(t, valid, "Is your device using test mnemonic: %s ?", testdata.TestMnemonic)
	}
}

func TestRealDeviceSecp256k1(t *testing.T) {
	msg := getFakeTx(50)
	path := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	priv, err := NewPrivKeySecp256k1Unsafe(path)
	require.NoError(t, err)

	pub := priv.PubKey()
	sig, err := priv.Sign(msg)
	require.NoError(t, err)

	valid := pub.VerifySignature(msg, sig)
	require.True(t, valid)

	// now, let's serialize the public key and make sure it still works
	bs := cdc.Amino.MustMarshalBinaryBare(priv.PubKey())
	pub2, err := legacy.PubKeyFromBytes(bs)
	require.Nil(t, err, "%+v", err)

	// make sure we get the same pubkey when we load from disk
	require.Equal(t, pub, pub2)

	// signing with the loaded key should match the original pubkey
	sig, err = priv.Sign(msg)
	require.NoError(t, err)
	valid = pub.VerifySignature(msg, sig)
	require.True(t, valid)

	// make sure pubkeys serialize properly as well
	bs = legacy.Cdc.MustMarshal(pub)
	bpub, err := legacy.PubKeyFromBytes(bs)
	require.NoError(t, err)
	require.Equal(t, pub, bpub)
}
