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
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func (k Keeper) handleBridgeTransfer(ctx sdk.Context, sender sdk.AccAddress, amount sdk.Int) (uint64, error) {
	token := sdk.Coins{sdk.Coin{Denom: k.GetParams(ctx).TargetDenom, Amount: amount}}
	if err := k.bankKeeper.IsSendEnabledCoins(ctx, token...); err != nil {
		return 0, err
	}

	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, token); err != nil {
		return 0, err
	}

	if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, token); err != nil {
		panic(fmt.Errorf("cannot burn coins after a successful send to a module account: %v", err))
	}

	seq := k.GetNextSequence(ctx)
	k.setNextSequence(ctx, seq+1)
	k.setSeqToBlocknum(ctx, seq, uint64(ctx.BlockHeight()))

	return seq, nil
}

func (k Keeper) setSeqToBlocknum(ctx sdk.Context, seq, height uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, height)
	store.Set(types.SeqToBlocknumKey(seq), bz)
}

func (k Keeper) GetSeqToBlocknum(ctx sdk.Context, seq uint64) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SeqToBlocknumKey(seq))
	if len(bz) == 0 {
		return 0, sdkerrors.ErrNotFound.Wrapf("sequence %d not found", seq)
	}

	return binary.BigEndian.Uint64(bz), nil
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

func IsValidEthereumAddress(address string) error {
	matched, err := regexp.MatchString(`^0x[a-fA-F0-9]{40}$`, address)
	if err != nil || !matched {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid eth address: %s", address)
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
			return err
		}
		if c < 8 {
			checksumAddress += string(addressLower[i])
		} else {
			checksumAddress += strings.ToUpper(string(addressLower[i]))
		}
	}

	if address != checksumAddress {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid checksum for eth address: %s", address)
	}

	return nil
}
