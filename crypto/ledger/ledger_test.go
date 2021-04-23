package ledger

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/codec/legacy"
	"github.com/line/lbm-sdk/v2/crypto/hd"
	"github.com/line/lbm-sdk/v2/crypto/types"
	"github.com/line/lbm-sdk/v2/testutil"
	sdk "github.com/line/lbm-sdk/v2/types"
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
	require.NotNil(t, priv)

	require.Equal(t, "eb5ae9872102bcd96cab102304ac10e90c5a0f29358e3a4a6fb1217b83e5db657918ea28bec1",
		fmt.Sprintf("%x", cdc.Amino.MustMarshalBinaryBare(priv.PubKey())),
		"Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

	pubKeyAddr, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, priv.PubKey())
	require.NoError(t, err)
	require.Equal(t, "linkpub1addwnpepq27djm9tzq3sftqsayx95refxk8r5jn0kyshhql9mdjhjx829zlvzszgelc",
		pubKeyAddr, "Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

	addr := sdk.AccAddress(priv.PubKey().Address()).String()
	require.Equal(t, "link1tdl7n2acgmec0y5nng0q2fahl9khyct3cgsktn",
		addr, "Is your device using test mnemonic: %s ?", testutil.TestMnemonic)
}

func TestPublicKeyUnsafeHDPath(t *testing.T) {
	expectedAnswers := []string{
		"linkpub1addwnpepq27djm9tzq3sftqsayx95refxk8r5jn0kyshhql9mdjhjx829zlvzszgelc",
		"linkpub1addwnpepqf258jtwpyujhxmlg94500j9yzqya5ryl835yp3dm6p9up25ufqcsjs2r6s",
		"linkpub1addwnpepq2edmckd0zthve9r70err6ctxzqc4vt5648lu4fzqkld8dnaekzjcl9j5ft",
		"linkpub1addwnpepqg9xfexl88nvmyyzpg5htz5qz30wgdftf0puz5u3sj6jkk9fxy7vzg7q5xn",
		"linkpub1addwnpepqv09egt2l0u72a4h0stkcrx4hyz0z6mnxe5w5d7lzmmzfdj2mykj755jfdt",
		"linkpub1addwnpepqfn9d7tew6vlr37sy9crsdud2gufsftm7wz3r2uhze2lfam4a263qsj6a8r",
		"linkpub1addwnpepqfaq649vgk3levrsya2wkz8aecjxxd40rdfjhr6aqlld5ql54fds2yg7xpt",
		"linkpub1addwnpepqv43zgg5dauwynq4wyqz3c6xtl9wcmc8z8ftgqvj87xs000lld6s63hlzt9",
		"linkpub1addwnpepq0kchl479dz7f28hgfn7ve3txkktu9trq2dpmrzjy9awlyuf8w6x7nugtw7",
		"linkpub1addwnpepqttsm9aacj9pq3w22xjms6lgyzxhhdjrrajt4hzzfl0melff9w9dq9e2ky7",
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

		pubKeyAddr, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, priv.PubKey())
		require.NoError(t, err)
		require.Equal(t,
			expectedAnswers[i], pubKeyAddr,
			"Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

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
	priv, addr, err := NewPrivKeySecp256k1(path, "cosmos")

	require.NoError(t, err)
	require.NotNil(t, priv)

	require.Nil(t, ShowAddress(path, priv.PubKey(), sdk.GetConfig().GetBech32AccountAddrPrefix()))

	require.Equal(t, "eb5ae9872102bcd96cab102304ac10e90c5a0f29358e3a4a6fb1217b83e5db657918ea28bec1",
		fmt.Sprintf("%x", cdc.Amino.MustMarshalBinaryBare(priv.PubKey())),
		"Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

	pubKeyAddr, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, priv.PubKey())
	require.NoError(t, err)
	require.Equal(t, "linkpub1addwnpepq27djm9tzq3sftqsayx95refxk8r5jn0kyshhql9mdjhjx829zlvzszgelc",
		pubKeyAddr, "Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

	require.Equal(t, "link1tdl7n2acgmec0y5nng0q2fahl9khyct3cgsktn",
		addr, "Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

	addr2 := sdk.AccAddress(priv.PubKey().Address()).String()
	require.Equal(t, addr, addr2)
}

func TestPublicKeyHDPath(t *testing.T) {
	expectedPubKeys := []string{
		"linkpub1addwnpepq27djm9tzq3sftqsayx95refxk8r5jn0kyshhql9mdjhjx829zlvzszgelc",
		"linkpub1addwnpepqf258jtwpyujhxmlg94500j9yzqya5ryl835yp3dm6p9up25ufqcsjs2r6s",
		"linkpub1addwnpepq2edmckd0zthve9r70err6ctxzqc4vt5648lu4fzqkld8dnaekzjcl9j5ft",
		"linkpub1addwnpepqg9xfexl88nvmyyzpg5htz5qz30wgdftf0puz5u3sj6jkk9fxy7vzg7q5xn",
		"linkpub1addwnpepqv09egt2l0u72a4h0stkcrx4hyz0z6mnxe5w5d7lzmmzfdj2mykj755jfdt",
		"linkpub1addwnpepqfn9d7tew6vlr37sy9crsdud2gufsftm7wz3r2uhze2lfam4a263qsj6a8r",
		"linkpub1addwnpepqfaq649vgk3levrsya2wkz8aecjxxd40rdfjhr6aqlld5ql54fds2yg7xpt",
		"linkpub1addwnpepqv43zgg5dauwynq4wyqz3c6xtl9wcmc8z8ftgqvj87xs000lld6s63hlzt9",
		"linkpub1addwnpepq0kchl479dz7f28hgfn7ve3txkktu9trq2dpmrzjy9awlyuf8w6x7nugtw7",
		"linkpub1addwnpepqttsm9aacj9pq3w22xjms6lgyzxhhdjrrajt4hzzfl0melff9w9dq9e2ky7",
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
	for i := uint32(0); i < 10; i++ {
		path := *hd.NewFundraiserParams(0, sdk.CoinType, i)
		t.Logf("Checking keys at %s\n", path)

		priv, addr, err := NewPrivKeySecp256k1(path, "cosmos")
		require.NoError(t, err)
		require.NotNil(t, addr)
		require.NotNil(t, priv)

		addr2 := sdk.AccAddress(priv.PubKey().Address()).String()
		require.Equal(t, addr2, addr)
		require.Equal(t,
			expectedAddrs[i], addr,
			"Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

		// Check other methods
		tmp := priv.(PrivKeyLedgerSecp256k1)
		require.NoError(t, tmp.ValidateKey())
		(&tmp).AssertIsPrivKeyInner()

		pubKeyAddr, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, priv.PubKey())
		require.NoError(t, err)
		require.Equal(t,
			expectedPubKeys[i], pubKeyAddr,
			"Is your device using test mnemonic: %s ?", testutil.TestMnemonic)

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
		require.True(t, valid, "Is your device using test mnemonic: %s ?", testutil.TestMnemonic)
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
	bs = legacy.Cdc.MustMarshalBinaryBare(pub)
	bpub, err := legacy.PubKeyFromBytes(bs)
	require.NoError(t, err)
	require.Equal(t, pub, bpub)
}
