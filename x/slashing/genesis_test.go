package slashing_test

import (
	"testing"
	"time"

	ostproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/simapp"
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/slashing"
	"github.com/line/lfb-sdk/x/slashing/testslashing"
	"github.com/line/lfb-sdk/x/slashing/types"
)

func TestExportAndInitGenesis(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ostproto.Header{})

	app.SlashingKeeper.SetParams(ctx, testslashing.TestParams())

	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 2, sdk.TokensFromConsensusPower(200))
	info1 := types.NewValidatorSigningInfo(addrDels[0].ToConsAddress(), int64(4), int64(3),
		time.Now().UTC().Add(100000000000), false, int64(10))
	info2 := types.NewValidatorSigningInfo(addrDels[1].ToConsAddress(), int64(5), int64(4),
		time.Now().UTC().Add(10000000000), false, int64(10))

	app.SlashingKeeper.SetValidatorSigningInfo(ctx, addrDels[0].ToConsAddress(), info1)
	app.SlashingKeeper.SetValidatorSigningInfo(ctx, addrDels[1].ToConsAddress(), info2)
	genesisState := slashing.ExportGenesis(ctx, app.SlashingKeeper)

	require.Equal(t, genesisState.Params, testslashing.TestParams())
	require.Len(t, genesisState.SigningInfos, 2)
	require.Equal(t, genesisState.SigningInfos[0].ValidatorSigningInfo, info1)

	// Tombstone validators after genesis shouldn't effect genesis state
	app.SlashingKeeper.Tombstone(ctx, addrDels[0].ToConsAddress())
	app.SlashingKeeper.Tombstone(ctx, addrDels[1].ToConsAddress())

	ok := app.SlashingKeeper.IsTombstoned(ctx, addrDels[0].ToConsAddress())
	require.True(t, ok)

	newInfo1, ok := app.SlashingKeeper.GetValidatorSigningInfo(ctx, addrDels[0].ToConsAddress())
	require.NotEqual(t, info1, newInfo1)
	// Initialise genesis with genesis state before tombstone
	slashing.InitGenesis(ctx, app.SlashingKeeper, app.StakingKeeper, genesisState)

	// Validator isTombstoned should return false as GenesisState is initialised
	ok = app.SlashingKeeper.IsTombstoned(ctx, addrDels[0].ToConsAddress())
	require.False(t, ok)

	newInfo1, ok = app.SlashingKeeper.GetValidatorSigningInfo(ctx, addrDels[0].ToConsAddress())
	newInfo2, ok := app.SlashingKeeper.GetValidatorSigningInfo(ctx, addrDels[1].ToConsAddress())
	require.True(t, ok)
	require.Equal(t, info1, newInfo1)
	require.Equal(t, info2, newInfo2)
}
