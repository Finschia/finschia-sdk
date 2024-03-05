package v3_test

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
	"github.com/Finschia/finschia-sdk/x/foundation/keeper/internal/migrations/v3"
	"github.com/Finschia/finschia-sdk/x/foundation/module"
)

type mockSubspace struct {
	params foundation.Params
}

func (ms mockSubspace) GetParamSet(ctx sdk.Context, ps v3.ParamSet) {
	*ps.(*foundation.Params) = ms.params
}

func (ms *mockSubspace) SetParamSet(ctx sdk.Context, ps v3.ParamSet) {
	ms.params = *ps.(*foundation.Params)
}

func TestMigrateStore(t *testing.T) {
	foundationKey := storetypes.NewKVStoreKey(foundation.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, foundationKey, storetypes.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModuleBasic{})

	ctx := testCtx.Ctx

	for name, tc := range map[string]struct {
		malleate func(sdk.Context, *mockSubspace)
		valid    bool
		tax      math.LegacyDec
	}{
		"valid": {
			malleate: func(ctx sdk.Context, subspace *mockSubspace) {
				// set old keys
				params := &foundation.Params{
					FoundationTax: math.LegacyMustNewDecFromStr("0.123456789"),
				}
				subspace.SetParamSet(ctx, params)
			},
			valid: true,
			tax:   math.LegacyMustNewDecFromStr("0.123456789"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, _ := ctx.CacheContext()

			subspace := &mockSubspace{}
			if tc.malleate != nil {
				tc.malleate(ctx, subspace)
			}

			// migrate
			err := v3.MigrateStore(ctx, runtime.NewKVStoreService(foundationKey), encCfg.Codec, subspace)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			store := ctx.KVStore(foundationKey)
			bz := store.Get(v3.ParamsKey)
			require.NotNil(t, bz)

			var params foundation.Params
			err = encCfg.Codec.Unmarshal(bz, &params)
			require.NoError(t, err)
			require.EqualValues(t, tc.tax, params.FoundationTax)
		})
	}
}
