package keeper_test

import (
	"testing"
	"time"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/simapp"
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/slashing/types"
)

func TestGetSetValidatorSigningInfo(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})
	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(200))

	info, found := app.SlashingKeeper.GetValidatorSigningInfo(ctx, addrDels[0].ToConsAddress())
	require.False(t, found)
	newInfo := types.NewValidatorSigningInfo(
		addrDels[0].ToConsAddress(),
		int64(4),
		int64(3),
		time.Unix(2, 0),
		false,
		int64(10),
	)
	app.SlashingKeeper.SetValidatorSigningInfo(ctx, addrDels[0].ToConsAddress(), newInfo)
	info, found = app.SlashingKeeper.GetValidatorSigningInfo(ctx, addrDels[0].ToConsAddress())
	require.True(t, found)
	require.Equal(t, info.StartHeight, int64(4))
	require.Equal(t, info.IndexOffset, int64(3))
	require.Equal(t, info.JailedUntil, time.Unix(2, 0).UTC())
	require.Equal(t, info.MissedBlocksCounter, int64(10))
}

func TestGetSetValidatorMissedBlockBitArray(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})
	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(200))

	missed := app.SlashingKeeper.GetValidatorMissedBlockBitArray(ctx, addrDels[0].ToConsAddress(), 0)
	require.False(t, missed) // treat empty key as not missed
	app.SlashingKeeper.SetValidatorMissedBlockBitArray(ctx, addrDels[0].ToConsAddress(), 0, true)
	missed = app.SlashingKeeper.GetValidatorMissedBlockBitArray(ctx, addrDels[0].ToConsAddress(), 0)
	require.True(t, missed) // now should be missed
}

func TestTombstoned(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})
	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(200))

	require.Panics(t, func() { app.SlashingKeeper.Tombstone(ctx, addrDels[0].ToConsAddress()) })
	require.False(t, app.SlashingKeeper.IsTombstoned(ctx, addrDels[0].ToConsAddress()))

	newInfo := types.NewValidatorSigningInfo(
		addrDels[0].ToConsAddress(),
		int64(4),
		int64(3),
		time.Unix(2, 0),
		false,
		int64(10),
	)
	app.SlashingKeeper.SetValidatorSigningInfo(ctx, addrDels[0].ToConsAddress(), newInfo)

	require.False(t, app.SlashingKeeper.IsTombstoned(ctx, addrDels[0].ToConsAddress()))
	app.SlashingKeeper.Tombstone(ctx, addrDels[0].ToConsAddress())
	require.True(t, app.SlashingKeeper.IsTombstoned(ctx, addrDels[0].ToConsAddress()))
	require.Panics(t, func() { app.SlashingKeeper.Tombstone(ctx, addrDels[0].ToConsAddress()) })
}

func TestJailUntil(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})
	addrDels := simapp.AddTestAddrsIncremental(app, ctx, 1, sdk.TokensFromConsensusPower(200))

	require.Panics(t, func() { app.SlashingKeeper.JailUntil(ctx, addrDels[0].ToConsAddress(), time.Now()) })

	newInfo := types.NewValidatorSigningInfo(
		addrDels[0].ToConsAddress(),
		int64(4),
		int64(3),
		time.Unix(2, 0),
		false,
		int64(10),
	)
	app.SlashingKeeper.SetValidatorSigningInfo(ctx, addrDels[0].ToConsAddress(), newInfo)
	app.SlashingKeeper.JailUntil(ctx, addrDels[0].ToConsAddress(), time.Unix(253402300799, 0).UTC())

	info, ok := app.SlashingKeeper.GetValidatorSigningInfo(ctx, addrDels[0].ToConsAddress())
	require.True(t, ok)
	require.Equal(t, time.Unix(253402300799, 0).UTC(), info.JailedUntil)
}
