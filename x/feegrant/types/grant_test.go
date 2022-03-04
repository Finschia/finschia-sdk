package types_test

import (
	"testing"
	"time"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/x/feegrant/types"

	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
)

func TestGrant(t *testing.T) {
	app := simapp.Setup(false)
	addr := sdk.AccAddress("link1qyqszqgpqyqszqgpqyqszqgpqyqszqgp8apuk5")
	addr2 := sdk.AccAddress("link1ghekyjucln7y67ntx7cf27m9dpuxxemnqk82wt")
	atom := sdk.NewCoins(sdk.NewInt64Coin("atom", 555))
	ctx := app.BaseApp.NewContext(false, ocproto.Header{
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
			grantee: "",
			limit:   atom,
			expires: oneYear,
			valid:   false,
		},
		"no granter": {
			granter: "",
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
			grant, err := types.NewGrant(tc.granter, tc.grantee, &types.BasicAllowance{
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
			bz, err := cdc.MarshalBinaryBare(&grant)
			require.NoError(t, err)
			var loaded types.Grant
			err = cdc.UnmarshalBinaryBare(bz, &loaded)
			require.NoError(t, err)

			err = loaded.ValidateBasic()
			require.NoError(t, err)

			require.Equal(t, grant, loaded)
		})
	}
}
