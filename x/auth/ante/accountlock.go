package ante

import (
	"encoding/binary"
	"sort"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AccountLockDecorator struct {
	addrMtx [1 << 24]sync.Mutex
}

func NewAccountLockDecorator() *AccountLockDecorator {
	return &AccountLockDecorator{}
}

func (ald *AccountLockDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if !ctx.IsCheckTx() || ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	stdTx, ok := tx.(types.StdTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a StdTx")
	}

	signers := stdTx.GetSigners()
	addrKeys := getUniqSortedAddressKey(signers)

	ald.lock(addrKeys)
	defer ald.unlock(addrKeys)

	return next(ctx, tx, simulate)
}

func (ald *AccountLockDecorator) lock(addrKeys []uint32) {
	for _, key := range addrKeys {
		ald.addrMtx[key].Lock()
	}
}

func (ald *AccountLockDecorator) unlock(addrKeys []uint32) {
	// NOTE reverse order
	for i, length := 0, len(addrKeys); i < length; i++ {
		key := addrKeys[length-1-i]
		ald.addrMtx[key].Unlock()
	}
}

func getUniqSortedAddressKey(addrs []sdk.AccAddress) []uint32 {
	addrKeys := make([]uint32, 0, len(addrs))
	for _, addr := range addrs {
		tail3 := addr[len(addr)-3:]
		tail := append([]byte{0}, tail3...)

		addrKey := binary.BigEndian.Uint32(tail)
		addrKeys = append(addrKeys, addrKey)
	}

	addrKeys = uniq(addrKeys)
	sort.Slice(addrKeys, func(i, j int) bool {
		return i < j
	})

	return addrKeys
}

func uniq(u []uint32) []uint32 {
	encountered := map[uint32]bool{}
	var ret []uint32

	for _, v := range u {
		if !encountered[v] {
			encountered[v] = true
			ret = append(ret, v)
		}
	}

	return ret
}
