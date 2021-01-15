package baseapp

import (
	"reflect"
	"sort"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestAccountLock(t *testing.T) {
	app := setupBaseApp(t)
	ctx := app.NewContext(true, abci.Header{})

	privs := newTestPrivKeys(3)
	tx := newTestTx(privs)

	accKeys := app.accountLock.Lock(ctx, tx)

	for _, accKey := range accKeys {
		require.True(t, isMutexLock(&app.accountLock.accMtx[accKey]))
	}

	app.accountLock.Unlock(accKeys)

	for _, accKey := range accKeys {
		require.False(t, isMutexLock(&app.accountLock.accMtx[accKey]))
	}
}

func TestUnlockDoNothingWithNil(t *testing.T) {
	app := setupBaseApp(t)
	require.NotPanics(t, func() { app.accountLock.Unlock(nil) })
}

func TestGetSigner(t *testing.T) {
	privs := newTestPrivKeys(3)
	tx := newTestTx(privs)
	signers := getSigners(tx)

	require.Equal(t, getAddrs(privs), signers)
}

func TestGetUniqSortedAddressKey(t *testing.T) {
	privs := newTestPrivKeys(3)

	addrs := getAddrs(privs)
	addrs = append(addrs, addrs[1], addrs[0])
	require.Equal(t, 5, len(addrs))

	accKeys := getUniqSortedAddressKey(addrs)

	// length should be reduced because `duplicated` is removed
	require.Less(t, len(accKeys), len(addrs))

	// check uniqueness
	for i, iv := range accKeys {
		for j, jv := range accKeys {
			if i != j {
				require.True(t, iv != jv)
			}
		}
	}

	// should be sorted
	require.True(t, sort.IsSorted(uint32Slice(accKeys)))
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

// Hack (too slow)
func isMutexLock(mtx *sync.Mutex) bool {
	state := reflect.ValueOf(mtx).Elem().FieldByName("state")
	return state.Int() == 1
}
