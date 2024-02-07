package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/foundation/keeper/internal/migrations/v2"
	"github.com/Finschia/finschia-sdk/x/foundation/module"
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
	foundationKey := storetypes.NewKVStoreKey(foundation.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, foundationKey, storetypes.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})

	ctx := testCtx.Ctx

	for name, tc := range map[string]struct {
		malleate func(ctx sdk.Context)
		valid    bool
		tax      math.LegacyDec
	}{
		"valid": {
			malleate: func(ctx sdk.Context) {
				// set old keys
				bz := encCfg.Codec.MustMarshal(&foundation.Params{
					FoundationTax: math.LegacyMustNewDecFromStr("0.123456789"),
				})
				store := ctx.KVStore(foundationKey)
				store.Set(v2.ParamsKey, bz)
			},
			valid: true,
			tax:   math.LegacyMustNewDecFromStr("0.123456789"),
		},
		"no params found": {},
		"unmarshal fails": {
			malleate: func(ctx sdk.Context) {
				// invalid contents
				bz := encCfg.Codec.MustMarshal(&foundation.Censorship{
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
			err := v2.MigrateStore(ctx, runtime.NewKVStoreService(foundationKey), encCfg.Codec, subspace)
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
