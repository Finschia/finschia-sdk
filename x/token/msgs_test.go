package token_test

import (
	"testing"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
	"github.com/stretchr/testify/require"
)

func TestMsgTransfer(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId string
		from sdk.AccAddress
		to sdk.AccAddress
		amount sdk.Int
		valid bool
	}{
		"valid msg": {
			classId: "deadbeef",
			from: addrs[0],
			to: addrs[1],
			amount: sdk.OneInt(),
			valid: true,
		},
		"empty from": {
			classId: "deadbeef",
			from: "",
			to: addrs[1],
			amount: sdk.OneInt(),
			valid: false,
		},
		"invalid to": {
			classId: "deadbeef",
			from: addrs[0],
			to: "invalid",
			amount: sdk.OneInt(),
			valid: false,
		},
		"zero amount": {
			classId: "deadbeef",
			from: addrs[0],
			to: addrs[1],
			amount: sdk.ZeroInt(),
			valid: false,
		},
	}

	for name, tc := range testCases {
		msg := token.MsgTransfer{
			ClassId: tc.classId,
			From: tc.from.String(),
			To: tc.to.String(),
			Amount: tc.amount,
		}

		require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgTransferFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId string
		proxy sdk.AccAddress
		from sdk.AccAddress
		to sdk.AccAddress
		amount sdk.Int
		valid bool
	}{
		"valid msg": {
			classId: "deadbeef",
			proxy: addrs[0],
			from: addrs[1],
			to: addrs[2],
			amount: sdk.OneInt(),
			valid: true,
		},
		"invalid proxy": {
			classId: "deadbeef",
			proxy: "invalid",
			from: addrs[1],
			to: addrs[2],
			amount: sdk.OneInt(),
			valid: false,
		},
		"empty from": {
			classId: "deadbeef",
			proxy: addrs[0],
			from: "",
			to: addrs[1],
			amount: sdk.OneInt(),
			valid: false,
		},
		"invalid to": {
			classId: "deadbeef",
			proxy: addrs[0],
			from: addrs[1],
			to: "invalid",
			amount: sdk.OneInt(),
			valid: false,
		},
		"zero amount": {
			classId: "deadbeef",
			proxy: addrs[0],
			from: addrs[1],
			to: addrs[2],
			amount: sdk.ZeroInt(),
			valid: false,
		},
	}

	for name, tc := range testCases {
		msg := token.MsgTransferFrom{
			ClassId: tc.classId,
			Proxy: tc.proxy.String(),
			From: tc.from.String(),
			To: tc.to.String(),
			Amount: tc.amount,
		}

		require.Equal(t, []sdk.AccAddress{tc.proxy}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgApprove(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId string
		approver sdk.AccAddress
		proxy sdk.AccAddress
		valid bool
	}{
		"valid msg": {
			classId: "deadbeef",
			approver: addrs[0],
			proxy: addrs[1],
			valid: true,
		},
		"invalid approver": {
			classId: "deadbeef",
			approver: "invalid",
			proxy: addrs[1],
			valid: false,
		},
		"empty proxy": {
			classId: "deadbeef",
			approver: addrs[0],
			proxy: "",
			valid: false,
		},
	}

	for name, tc := range testCases {
		msg := token.MsgApprove{
			ClassId: tc.classId,
			Approver: tc.approver.String(),
			Proxy: tc.proxy.String(),
		}

		require.Equal(t, []sdk.AccAddress{tc.approver}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgIssue(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		owner sdk.AccAddress
		to sdk.AccAddress
		name string
		symbol string
		imageUri string
		meta string
		decimals int32
		amount sdk.Int
		valid bool
	}{
		"valid msg": {
			owner: addrs[0],
			to: addrs[1],
			name: "test",
			symbol: "TT",
			imageUri: "some URI",
			meta: "some meta",
			decimals: 8,
			amount: sdk.OneInt(),
			valid: true,
		},
		"invalid owner": {
			owner: "invalid",
			to: addrs[1],
			name: "test",
			symbol: "TT",
			imageUri: "some URI",
			meta: "some meta",
			decimals: 8,
			amount: sdk.OneInt(),
			valid: false,
		},
		"empty to": {
			owner: addrs[0],
			to: "",
			name: "test",
			symbol: "TT",
			imageUri: "some URI",
			meta: "some meta",
			decimals: 8,
			amount: sdk.OneInt(),
			valid: false,
		},
		"empty name": {
			owner: addrs[0],
			to: addrs[1],
			name: "",
			symbol: "TT",
			imageUri: "some URI",
			meta: "some meta",
			decimals: 8,
			amount: sdk.OneInt(),
			valid: false,
		},
		"long name": {
			owner: addrs[0],
			to: addrs[1],
			name: "TOO Looooooooooooooog",
			symbol: "TT",
			imageUri: "some URI",
			meta: "some meta",
			decimals: 8,
			amount: sdk.OneInt(),
			valid: false,
		},
		"invalid symbol": {
			owner: addrs[0],
			to: addrs[1],
			name: "test",
			symbol: "tt",
			imageUri: "some URI",
			meta: "some meta",
			decimals: 8,
			amount: sdk.OneInt(),
			valid: false,
		},
		"invalid image uri": {
			owner: addrs[0],
			to: addrs[1],
			name: "test",
			symbol: "TT",
			imageUri: string(make([]rune, 1001)),
			meta: "some meta",
			decimals: 8,
			amount: sdk.OneInt(),
			valid: false,
		},
		"invalid meta": {
			owner: addrs[0],
			to: addrs[1],
			name: "test",
			symbol: "TT",
			imageUri: "some URI",
			meta: string(make([]rune, 1001)),
			decimals: 8,
			amount: sdk.OneInt(),
			valid: false,
		},
		"invalid decimals": {
			owner: addrs[0],
			to: addrs[1],
			name: "test",
			symbol: "TT",
			imageUri: "some URI",
			meta: "some meta",
			decimals: 19,
			amount: sdk.OneInt(),
			valid: false,
		},
		"valid supply": {
			owner: addrs[0],
			to: addrs[1],
			name: "test",
			symbol: "TT",
			imageUri: "some URI",
			meta: "some meta",
			decimals: 8,
			amount: sdk.ZeroInt(),
			valid: false,
		},
	}

	for name, tc := range testCases {
		msg := token.MsgIssue{
			Owner: tc.owner.String(),
			To: tc.to.String(),
			Name: tc.name,
			Symbol: tc.symbol,
			ImageUri: tc.imageUri,
			Meta: tc.meta,
			Decimals: tc.decimals,
			Amount: tc.amount,
		}

		require.Equal(t, []sdk.AccAddress{tc.owner}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgMint(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId string
		grantee sdk.AccAddress
		to sdk.AccAddress
		amount sdk.Int
		valid bool
	}{
		"valid msg": {
			classId: "deadbeef",
			grantee: addrs[0],
			to: addrs[1],
			amount: sdk.OneInt(),
			valid: true,
		},
		"invalid grantee": {
			classId: "deadbeef",
			grantee: "invalid",
			to: addrs[1],
			amount: sdk.OneInt(),
			valid: false,
		},
		"empty to": {
			classId: "deadbeef",
			grantee: addrs[0],
			to: "",
			amount: sdk.OneInt(),
			valid: false,
		},
		"zero amount": {
			classId: "deadbeef",
			grantee: addrs[0],
			to: addrs[1],
			amount: sdk.ZeroInt(),
			valid: false,
		},
	}

	for name, tc := range testCases {
		msg := token.MsgMint{
			ClassId: tc.classId,
			Grantee: tc.grantee.String(),
			To: tc.to.String(),
			Amount: tc.amount,
		}

		require.Equal(t, []sdk.AccAddress{tc.grantee}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgBurn(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId string
		from sdk.AccAddress
		amount sdk.Int
		valid bool
	}{
		"valid msg": {
			classId: "deadbeef",
			from: addrs[0],
			amount: sdk.OneInt(),
			valid: true,
		},
		"invalid from": {
			classId: "deadbeef",
			from: "invalid",
			amount: sdk.OneInt(),
			valid: false,
		},
		"zero amount": {
			classId: "deadbeef",
			from: addrs[0],
			amount: sdk.ZeroInt(),
			valid: false,
		},
	}

	for name, tc := range testCases {
		msg := token.MsgBurn{
			ClassId: tc.classId,
			From: tc.from.String(),
			Amount: tc.amount,
		}

		require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgBurnFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId string
		grantee sdk.AccAddress
		from sdk.AccAddress
		amount sdk.Int
		valid bool
	}{
		"valid msg": {
			classId: "deadbeef",
			grantee: addrs[0],
			from: addrs[1],
			amount: sdk.OneInt(),
			valid: true,
		},
		"invalid grantee": {
			classId: "deadbeef",
			grantee: "invalid",
			from: addrs[1],
			amount: sdk.OneInt(),
			valid: false,
		},
		"empty from": {
			classId: "deadbeef",
			grantee: addrs[0],
			from: "",
			amount: sdk.OneInt(),
			valid: false,
		},
		"zero amount": {
			classId: "deadbeef",
			grantee: addrs[0],
			from: addrs[1],
			amount: sdk.ZeroInt(),
			valid: false,
		},
	}

	for name, tc := range testCases {
		msg := token.MsgBurnFrom{
			ClassId: tc.classId,
			Grantee: tc.grantee.String(),
			From: tc.from.String(),
			Amount: tc.amount,
		}

		require.Equal(t, []sdk.AccAddress{tc.grantee}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgModify(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	validChange := token.Pair{Key: "name", Value: "New test"}
	testCases := map[string]struct {
		classId string
		grantee sdk.AccAddress
		changes []token.Pair
		valid bool
	}{
		"valid msg": {
			classId: "deadbeef",
			grantee: addrs[0],
			changes: []token.Pair{validChange},
			valid: true,
		},
		"invalid grantee": {
			classId: "deadbeef",
			grantee: "invalid",
			changes: []token.Pair{validChange},
			valid: false,
		},
		"invalid key of change": {
			classId: "deadbeef",
			grantee: addrs[0],
			changes: []token.Pair{{Key: "invalid", Value: "tt"}},
			valid: false,
		},
		"invalid value of change": {
			classId: "deadbeef",
			grantee: addrs[0],
			changes: []token.Pair{{Key: "symbol", Value: "tt"}},
			valid: false,
		},
		"empty changes": {
			classId: "deadbeef",
			grantee: addrs[0],
			changes: []token.Pair{},
			valid: false,
		},
		"duplicated changes": {
			classId: "deadbeef",
			grantee: addrs[0],
			changes: []token.Pair{
				{Key: "name", Value: "hello"},
				{Key: "name", Value: "world"},
			},
			valid: false,
		},
	}

	for name, tc := range testCases {
		msg := token.MsgModify{
			ClassId: tc.classId,
			Grantee: tc.grantee.String(),
			Changes: tc.changes,
		}

		require.Equal(t, []sdk.AccAddress{tc.grantee}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgGrant(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId string
		granter sdk.AccAddress
		grantee sdk.AccAddress
		action string
		valid bool
	}{
		"valid msg": {
			classId: "deadbeef",
			granter: addrs[0],
			grantee: addrs[1],
			action: "mint",
			valid: true,
		},
		"empty granter": {
			classId: "deadbeef",
			granter: "",
			grantee: addrs[1],
			action: "mint",
			valid: false,
		},
		"invalid grantee": {
			classId: "deadbeef",
			granter: addrs[0],
			grantee: "invalid",
			action: "mint",
			valid: false,
		},
		"invalid action": {
			classId: "deadbeef",
			granter: addrs[0],
			grantee: addrs[1],
			action: "invalid",
			valid: false,
		},
	}

	for name, tc := range testCases {
		msg := token.MsgGrant{
			ClassId: tc.classId,
			Granter: tc.granter.String(),
			Grantee: tc.grantee.String(),
			Action: tc.action,
		}

		require.Equal(t, []sdk.AccAddress{tc.granter}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgRevoke(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId string
		grantee sdk.AccAddress
		action string
		valid bool
	}{
		"valid msg": {
			classId: "deadbeef",
			grantee: addrs[0],
			action: "mint",
			valid: true,
		},
		"invalid grantee": {
			classId: "deadbeef",
			grantee: "invalid",
			action: "mint",
			valid: false,
		},
		"invalid action": {
			classId: "deadbeef",
			grantee: addrs[0],
			action: "invalid",
			valid: false,
		},
	}

	for name, tc := range testCases {
		msg := token.MsgRevoke{
			ClassId: tc.classId,
			Grantee: tc.grantee.String(),
			Action: tc.action,
		}

		require.Equal(t, []sdk.AccAddress{tc.grantee}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}
