package keeper

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func (k Keeper) handleBridgeTransfer(ctx sdk.Context, sender sdk.AccAddress, amount sdk.Int) (uint64, error) {
	token := sdk.Coins{sdk.Coin{Denom: k.targetDenom, Amount: amount}}
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, token); err != nil {
		panic(err)
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, token); err != nil {
		panic(fmt.Errorf("cannot burn coins after a successful send to a module account: %v", err))
	}

	nextSeq := k.GetNextSequence(ctx) + 1
	k.setNextSequence(ctx, nextSeq)

	return nextSeq, nil
}

func (k Keeper) GetNextSequence(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyNextSeqSend)
	if len(bz) == 0 {
		panic(errors.New("sending next sequence should have been set"))
	}

	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) setNextSequence(ctx sdk.Context, seq uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, seq)
	store.Set(types.KeyNextSeqSend, bz)
}

func IsValidEthereumAddress(address string) bool {
	matched, err := regexp.MatchString(`^0x[a-fA-F0-9]{40}$`, address)
	if err != nil || !matched {
		return false
	}

	address = address[2:]
	addressLower := strings.ToLower(address)

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(addressLower))
	addressHash := hex.EncodeToString(hasher.Sum(nil))

	checksumAddress := ""
	for i := 0; i < len(addressLower); i++ {
		c, err := strconv.ParseUint(string(addressHash[i]), 16, 4)
		if err != nil {
			panic(err)
		}
		if c < 8 {
			checksumAddress += string(addressLower[i])
		} else {
			checksumAddress += strings.ToUpper(string(addressLower[i]))
		}
	}

	return address == checksumAddress
}
