package keeper

import (
	"encoding/binary"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, gs *types.GenesisState) error {
	k.setNextSequence(ctx, gs.SendingState.NextSeq)
	for _, info := range gs.SendingState.SeqToBlocknum {
		k.setSeqToBlocknum(ctx, info.Seq, info.Blocknum)
	}

	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		SendingState: types.SendingState{
			NextSeq:       k.GetNextSequence(ctx),
			SeqToBlocknum: k.getAllSeqToBlocknums(ctx),
		},
	}
}

func (k Keeper) getAllSeqToBlocknums(ctx sdk.Context) []types.BlockSeqInfo {
	infos := make([]types.BlockSeqInfo, 0)
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeySeqToBlocknumPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		seq := binary.BigEndian.Uint64(iterator.Key()[1:])
		v := binary.BigEndian.Uint64(iterator.Value())
		info := types.BlockSeqInfo{Seq: seq, Blocknum: v}
		infos = append(infos, info)
	}

	return infos
}
