package baseapp

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestConvertByteSliceToString(t *testing.T) {
	b := []byte{65, 66, 67, 0, 65, 66, 67}
	s := string(b)
	require.Equal(t, len(b), len(s))
	require.Equal(t, uint8(0), s[3])
}

func TestRegister(t *testing.T) {
	app := setupBaseApp(t)

	privs := newTestPrivKeys(3)
	tx := newTestTx(privs)

	waits, signals := app.checkAccountWGs.Register(tx)

	require.Equal(t, 0, len(waits))
	require.Equal(t, 3, len(signals))

	for _, signal := range signals {
		require.Equal(t, app.checkAccountWGs.wgs[signal.acc], signal.wg)
	}
}

func TestDontPanicWithNil(t *testing.T) {
	app := setupBaseApp(t)

	require.NotPanics(t, func() { app.checkAccountWGs.Waits(nil) })
	require.NotPanics(t, func() { app.checkAccountWGs.Done(nil) })
}

func TestGetUniqSigners(t *testing.T) {
	privs := newTestPrivKeys(3)

	addrs := getAddrs(privs)
	addrs = append(addrs, addrs[1], addrs[0])
	require.Equal(t, 5, len(addrs))

	tx := newTestTx(privs)
	signers := getUniqSigners(tx)

	// length should be reduced because `duplicated` is removed
	require.Less(t, len(signers), len(addrs))

	// check uniqueness
	for i, iv := range signers {
		for j, jv := range signers {
			if i != j {
				require.True(t, iv != jv)
			}
		}
	}
}

type AccountLockTestTx struct {
	Msgs []sdk.Msg
}

var _ sdk.Tx = AccountLockTestTx{}

func (tx AccountLockTestTx) GetMsgs() []sdk.Msg {
	return tx.Msgs
}

func (tx AccountLockTestTx) ValidateBasic() error {
	return nil
}

func newTestPrivKeys(num int) []crypto.PrivKey {
	privs := make([]crypto.PrivKey, 0, num)
	for i := 0; i < num; i++ {
		privs = append(privs, secp256k1.GenPrivKey())
	}
	return privs
}

func getAddrs(privs []crypto.PrivKey) []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 0, len(privs))
	for _, priv := range privs {
		addrs = append(addrs, sdk.AccAddress(priv.PubKey().Address()))
	}
	return addrs
}

func newTestTx(privs []crypto.PrivKey) sdk.Tx {
	addrs := getAddrs(privs)
	msgs := make([]sdk.Msg, len(addrs))
	for i, addr := range addrs {
		msgs[i] = sdk.NewTestMsg(addr)
	}
	return AccountLockTestTx{Msgs: msgs}
}
