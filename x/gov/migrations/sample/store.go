package sample

import (
	"github.com/line/lbm-sdk/codec"
	types2 "github.com/line/lbm-sdk/store/types"
	"github.com/line/lbm-sdk/types"
	v1 "github.com/line/lbm-sdk/x/gov/migrations/v1"
	types3 "github.com/line/lbm-sdk/x/gov/types"
)

// This migrates all ProposalV1s to Proposals
func MigrateStore(ctx types.Context, storeKey types2.StoreKey, cdc codec.BinaryMarshaler) error {
	store := ctx.KVStore(storeKey)

	iterator := types.KVStorePrefixIterator(store, types3.ProposalsKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var proposalv1 v1.ProposalV1
		cdc.MustUnmarshalBinaryBare(iterator.Value(), &proposalv1)
		proposal := v1.ProposalV1ToProposal(proposalv1)
		bz := cdc.MustMarshalBinaryBare(&proposal)

		store.Delete(iterator.Key())
		store.Set(iterator.Key(), bz)
	}
	return nil
}
