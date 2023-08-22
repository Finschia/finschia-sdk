package feegrant_test

import (
	"testing"
	"time"

	proto "github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/Finschia/finschia-rdk/simapp"
	sdk "github.com/Finschia/finschia-rdk/types"
	banktypes "github.com/Finschia/finschia-rdk/x/bank/types"
	"github.com/Finschia/finschia-rdk/x/feegrant"
)

func TestFilteredFeeValidAllow(t *testing.T) {
	app := simapp.Setup(false)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{
		Time: time.Now(),
	})
	eth := sdk.NewCoins(sdk.NewInt64Coin("eth", 10))
	atom := sdk.NewCoins(sdk.NewInt64Coin("atom", 555))
	smallAtom := sdk.NewCoins(sdk.NewInt64Coin("atom", 43))
	bigAtom := sdk.NewCoins(sdk.NewInt64Coin("atom", 1000))
	leftAtom := sdk.NewCoins(sdk.NewInt64Coin("atom", 512))
	now := ctx.BlockTime()
	oneHour := now.Add(1 * time.Hour)

	// msg we will call in the all cases
	call := banktypes.MsgSend{}
	cases := map[string]struct {
		allowance *feegrant.BasicAllowance
		msgs      []string
		fee       sdk.Coins
		blockTime time.Time
		accept    bool
		remove    bool
		remains   sdk.Coins
	}{
		"msg contained": {
			allowance: &feegrant.BasicAllowance{},
			msgs:      []string{sdk.MsgTypeURL(&call)},
			accept:    true,
		},
		"msg not contained": {
			allowance: &feegrant.BasicAllowance{},
			msgs:      []string{"/cosmos.gov.v1beta1.MsgVote"},
			accept:    false,
		},
		"small fee without expire": {
			allowance: &feegrant.BasicAllowance{
				SpendLimit: atom,
			},
			msgs:    []string{sdk.MsgTypeURL(&call)},
			fee:     smallAtom,
			accept:  true,
			remove:  false,
			remains: leftAtom,
		},
		"all fee without expire": {
			allowance: &feegrant.BasicAllowance{
				SpendLimit: smallAtom,
			},
			msgs:   []string{sdk.MsgTypeURL(&call)},
			fee:    smallAtom,
			accept: true,
			remove: true,
		},
		"wrong fee": {
			allowance: &feegrant.BasicAllowance{
				SpendLimit: smallAtom,
			},
			msgs:   []string{sdk.MsgTypeURL(&call)},
			fee:    eth,
			accept: false,
		},
		"non-expired": {
			allowance: &feegrant.BasicAllowance{
				SpendLimit: atom,
				Expiration: &oneHour,
			},
			msgs:      []string{sdk.MsgTypeURL(&call)},
			fee:       smallAtom,
			blockTime: now,
			accept:    true,
			remove:    false,
			remains:   leftAtom,
		},
		"expired": {
			allowance: &feegrant.BasicAllowance{
				SpendLimit: atom,
				Expiration: &now,
			},
			msgs:      []string{sdk.MsgTypeURL(&call)},
			fee:       smallAtom,
			blockTime: oneHour,
			accept:    false,
			remove:    true,
		},
		"fee more than allowed": {
			allowance: &feegrant.BasicAllowance{
				SpendLimit: atom,
				Expiration: &oneHour,
			},
			msgs:      []string{sdk.MsgTypeURL(&call)},
			fee:       bigAtom,
			blockTime: now,
			accept:    false,
		},
		"with out spend limit": {
			allowance: &feegrant.BasicAllowance{
				Expiration: &oneHour,
			},
			msgs:      []string{sdk.MsgTypeURL(&call)},
			fee:       bigAtom,
			blockTime: now,
			accept:    true,
		},
		"expired no spend limit": {
			allowance: &feegrant.BasicAllowance{
				Expiration: &now,
			},
			msgs:      []string{sdk.MsgTypeURL(&call)},
			fee:       bigAtom,
			blockTime: oneHour,
			accept:    false,
		},
	}

	for name, stc := range cases {
		tc := stc // to make scopelint happy
		t.Run(name, func(t *testing.T) {
			err := tc.allowance.ValidateBasic()
			require.NoError(t, err)

			ctx := app.BaseApp.NewContext(false, tmproto.Header{}).WithBlockTime(tc.blockTime)

			// create grant
			var granter, grantee sdk.AccAddress
			allowance, err := feegrant.NewAllowedMsgAllowance(tc.allowance, tc.msgs)
			require.NoError(t, err)
			grant, err := feegrant.NewGrant(granter, grantee, allowance)
			require.NoError(t, err)

			// now try to deduct
			removed, err := allowance.Accept(ctx, tc.fee, []sdk.Msg{&call})
			if !tc.accept {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tc.remove, removed)
			if !removed {
				// mimic save & load process (#10564)
				// the cached allowance was correct even before the fix,
				// however, the saved value was not.
				// so we need this to catch the bug.

				newGranter, _ := sdk.AccAddressFromBech32(grant.Granter)
				newGrantee, _ := sdk.AccAddressFromBech32(grant.Grantee)
				// create a new updated grant
				newGrant, err := feegrant.NewGrant(
					newGranter,
					newGrantee,
					allowance)
				require.NoError(t, err)

				// save the grant
				cdc := simapp.MakeTestEncodingConfig().Marshaler
				bz, err := cdc.Marshal(&newGrant)
				require.NoError(t, err)

				// load the grant
				var loadedGrant feegrant.Grant
				err = cdc.Unmarshal(bz, &loadedGrant)
				require.NoError(t, err)

				newAllowance, err := loadedGrant.GetGrant()
				require.NoError(t, err)
				feeAllowance, err := newAllowance.(*feegrant.AllowedMsgAllowance).GetAllowance()
				require.NoError(t, err)
				assert.Equal(t, tc.remains, feeAllowance.(*feegrant.BasicAllowance).SpendLimit)
			}
		})
	}
}

// invalidInterfaceAllowance does not implement proto.Message
type invalidInterfaceAllowance struct {
}

// compilation time interface implementation check
var _ feegrant.FeeAllowanceI = (*invalidInterfaceAllowance)(nil)

func (i invalidInterfaceAllowance) Accept(ctx sdk.Context, fee sdk.Coins, msgs []sdk.Msg) (remove bool, err error) {
	return false, nil
}

func (i invalidInterfaceAllowance) ValidateBasic() error {
	return nil
}

// invalidProtoAllowance can not run proto.Marshal
type invalidProtoAllowance struct {
	invalidInterfaceAllowance
}

// compilation time interface implementation check
var _ feegrant.FeeAllowanceI = (*invalidProtoAllowance)(nil)
var _ proto.Message = (*invalidProtoAllowance)(nil)

func (i invalidProtoAllowance) Reset() {
}

func (i invalidProtoAllowance) String() string {
	return ""
}

func (i invalidProtoAllowance) ProtoMessage() {
}

func TestSetAllowance(t *testing.T) {
	cases := map[string]struct {
		allowance feegrant.FeeAllowanceI
		valid     bool
	}{
		"valid allowance": {
			allowance: &feegrant.BasicAllowance{},
			valid:     true,
		},
		"invalid interface allowance": {
			allowance: &invalidInterfaceAllowance{},
			valid:     false,
		},
		"empty allowance": {
			allowance: (*invalidProtoAllowance)(nil),
			valid:     false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			allowance := &feegrant.BasicAllowance{}
			msgs := []string{sdk.MsgTypeURL(&banktypes.MsgSend{})}
			allowed, err := feegrant.NewAllowedMsgAllowance(allowance, msgs)
			require.NoError(t, err)
			require.NotNil(t, allowed)
			err = allowed.SetAllowance(tc.allowance)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
