package keeper_test

import (
	"github.com/Finschia/finschia-sdk/x/foundation"
	govtypes "github.com/Finschia/finschia-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/golang/mock/gomock"

	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/keeper"
	"github.com/Finschia/finschia-sdk/x/fbridge/testutil"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func TestNewKeeper(t *testing.T) {
	key, memKey, _, encCfg, _, bankKeeper, _ := testutil.PrepareFbridgeTest(t, 0)
	authKeeper := testutil.NewMockAccountKeeper(gomock.NewController(t))

	tcs := map[string]struct {
		malleate func()
		isPanic  bool
	}{
		"fbrdige module account has not been set": {
			malleate: func() {
				authKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(nil).Times(1)
				keeper.NewKeeper(encCfg.Codec, key, memKey, authKeeper, bankKeeper, "stake", types.DefaultAuthority().String())
			},
			isPanic: true,
		},
		"x/bridge authority must be the gov or foundation module account": {
			malleate: func() {
				authKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(authtypes.NewModuleAddress(types.ModuleName)).Times(1)
				keeper.NewKeeper(encCfg.Codec, key, memKey, authKeeper, bankKeeper, "stake", authtypes.NewModuleAddress("invalid").String())
			},
			isPanic: true,
		},
		"success - gov authority": {
			malleate: func() {
				authKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(authtypes.NewModuleAddress(types.ModuleName)).Times(1)
				keeper.NewKeeper(encCfg.Codec, key, memKey, authKeeper, bankKeeper, "stake", authtypes.NewModuleAddress(govtypes.ModuleName).String())
			},
			isPanic: false,
		},
		"success - foundation authority": {
			malleate: func() {
				authKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(authtypes.NewModuleAddress(types.ModuleName)).Times(1)
				keeper.NewKeeper(encCfg.Codec, key, memKey, authKeeper, bankKeeper, "stake", authtypes.NewModuleAddress(foundation.ModuleName).String())
			},
			isPanic: false,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			if tc.isPanic {
				require.Panics(t, tc.malleate)
			} else {
				tc.malleate()
			}
		})
	}
}
