package slashing_test

import (
	"testing"
	"time"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/slashing"
	"github.com/line/lbm-sdk/x/slashing/testslashing"
	"github.com/line/lbm-sdk/x/slashing/types"
)

func TestExportAndInitGenesis(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	app.SlashingKeeper.SetParams(ctx, testslashing.TestParams())

	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 2, app.StakingKeeper.TokensFromConsensusPower(ctx, 200))
	info1 := types.NewValidatorSigningInfo(sdk.ConsAddress(addrDels[0]),
		time.Now().UTC().Add(100000000000), false, int64(10), int64(3))
	info2 := types.NewValidatorSigningInfo(sdk.ConsAddress(addrDels[1]),
		time.Now().UTC().Add(10000000000), false, int64(10), int64(4))

	app.SlashingKeeper.SetValidatorSigningInfo(ctx, sdk.ConsAddress(addrDels[0]), info1)
	app.SlashingKeeper.SetValidatorSigningInfo(ctx, sdk.ConsAddress(addrDels[1]), info2)
	genesisState := slashing.ExportGenesis(ctx, app.SlashingKeeper)

	require.Equal(t, genesisState.Params, testslashing.TestParams())
	require.Len(t, genesisState.SigningInfos, 2)
	require.Equal(t, genesisState.SigningInfos[0].ValidatorSigningInfo, info1)

	// Tombstone validators after genesis shouldn't effect genesis state
	app.SlashingKeeper.Tombstone(ctx, sdk.ConsAddress(addrDels[0]))
	app.SlashingKeeper.Tombstone(ctx, sdk.ConsAddress(addrDels[1]))

	ok := app.SlashingKeeper.IsTombstoned(ctx, sdk.ConsAddress(addrDels[0]))
	require.True(t, ok)

	newInfo1, ok := app.SlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(addrDels[0]))
	require.NotEqual(t, info1, newInfo1)
	// Initialise genesis with genesis state before tombstone
	slashing.InitGenesis(ctx, app.SlashingKeeper, app.StakingKeeper, genesisState)

	// Validator isTombstoned should return false as GenesisState is initialised
	ok = app.SlashingKeeper.IsTombstoned(ctx, sdk.ConsAddress(addrDels[0]))
	require.False(t, ok)

	newInfo1, ok = app.SlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(addrDels[0]))
	newInfo2, ok := app.SlashingKeeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(addrDels[1]))
	require.True(t, ok)
	require.Equal(t, info1, newInfo1)
	require.Equal(t, info2, newInfo2)
}
