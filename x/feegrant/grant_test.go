package feegrant_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/feegrant"
)

func TestGrant(t *testing.T) {
	app := simapp.Setup(false)
	addr, err := sdk.AccAddressFromBech32("link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5")
	require.NoError(t, err)
	addr2, err := sdk.AccAddressFromBech32("link1ghekyjucln7y67ntx7cf27m9dpuxxemnqk82wt")
	require.NoError(t, err)
	atom := sdk.NewCoins(sdk.NewInt64Coin("atom", 555))
	ctx := app.BaseApp.NewContext(false, tmproto.Header{
		Time: time.Now(),
	})
	now := ctx.BlockTime()
	oneYear := now.AddDate(1, 0, 0)

	zeroAtoms := sdk.NewCoins(sdk.NewInt64Coin("atom", 0))
	cdc := app.AppCodec()

	cases := map[string]struct {
		granter sdk.AccAddress
		grantee sdk.AccAddress
		limit   sdk.Coins
		expires time.Time
		valid   bool
	}{
		"good": {
			granter: addr2,
			grantee: addr,
			limit:   atom,
			expires: oneYear,
			valid:   true,
		},
		"no grantee": {
			granter: addr2,
			grantee: nil,
			limit:   atom,
			expires: oneYear,
			valid:   false,
		},
		"no granter": {
			granter: nil,
			grantee: addr,
			limit:   atom,
			expires: oneYear,
			valid:   false,
		},
		"self-grant": {
			granter: addr2,
			grantee: addr2,
			limit:   atom,
			expires: oneYear,
			valid:   false,
		},
		"zero allowance": {
			granter: addr2,
			grantee: addr,
			limit:   zeroAtoms,
			expires: oneYear,
			valid:   false,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			grant, err := feegrant.NewGrant(tc.granter, tc.grantee, &feegrant.BasicAllowance{
				SpendLimit: tc.limit,
				Expiration: &tc.expires,
			})
			require.NoError(t, err)
			err = grant.ValidateBasic()

			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// if it is valid, let's try to serialize, deserialize, and make sure it matches
			bz, err := cdc.Marshal(&grant)
			require.NoError(t, err)
			var loaded feegrant.Grant
			err = cdc.Unmarshal(bz, &loaded)
			require.NoError(t, err)

			err = loaded.ValidateBasic()
			require.NoError(t, err)

			require.Equal(t, grant, loaded)
		})
	}
}
