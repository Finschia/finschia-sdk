package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	simappparams "github.com/Finschia/finschia-sdk/simapp/params"
	"github.com/Finschia/finschia-sdk/testutil"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/foundation"

	"github.com/Finschia/finschia-sdk/x/foundation/keeper/internal/migrations/v2"
)

type mockSubspace struct {
	params foundation.Params
}

func (ms mockSubspace) GetParamSet(ctx sdk.Context, ps v2.ParamSet) {
	*ps.(*foundation.Params) = ms.params
}

func (ms *mockSubspace) SetParamSet(ctx sdk.Context, ps v2.ParamSet) {
	ms.params = *ps.(*foundation.Params)
}

func TestMigrateStore(t *testing.T) {
	foundationKey := sdk.NewKVStoreKey(foundation.StoreKey)
	newKey := sdk.NewTransientStoreKey("transient_test")
	encCfg := simappparams.MakeTestEncodingConfig()
	ctx := testutil.DefaultContext(foundationKey, newKey)

	for name, tc := range map[string]struct {
		malleate func(ctx sdk.Context)
		valid    bool
		tax      sdk.Dec
	}{
		"valid": {
			malleate: func(ctx sdk.Context) {
				// set old keys
				bz := encCfg.Marshaler.MustMarshal(&foundation.Params{
					FoundationTax: sdk.MustNewDecFromStr("0.123456789"),
				})
				store := ctx.KVStore(foundationKey)
				store.Set(v2.ParamsKey, bz)
			},
			valid: true,
			tax:   sdk.MustNewDecFromStr("0.123456789"),
		},
		"no params found": {},
		"unmarshal fails": {
			malleate: func(ctx sdk.Context) {
				// invalid contents
				bz := encCfg.Marshaler.MustMarshal(&foundation.Censorship{
					MsgTypeUrl: sdk.MsgTypeURL((*foundation.MsgWithdrawFromTreasury)(nil)),
					Authority:  foundation.CensorshipAuthorityFoundation,
				})
				store := ctx.KVStore(foundationKey)
				store.Set(v2.ParamsKey, bz)
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, _ := ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			// migrate
			subspace := &mockSubspace{}
			err := v2.MigrateStore(ctx, foundationKey, encCfg.Marshaler, subspace)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			store := ctx.KVStore(foundationKey)
			require.Nil(t, store.Get(v2.ParamsKey))

			var params foundation.Params
			subspace.GetParamSet(ctx, &params)
			require.EqualValues(t, tc.tax, params.FoundationTax)
		})
	}
}
