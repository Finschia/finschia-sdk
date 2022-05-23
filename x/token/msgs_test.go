package token_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func TestMsgSend(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId string
		from    sdk.AccAddress
		to      sdk.AccAddress
		amount  sdk.Int
		valid   bool
	}{
		"valid msg": {
			classId: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
			amount:  sdk.OneInt(),
			valid:   true,
		},
		"empty from": {
			classId: "deadbeef",
			from:    "",
			to:      addrs[1],
			amount:  sdk.OneInt(),
		},
		"invalid class id": {
			from:    addrs[0],
			to:      addrs[1],
			amount:  sdk.OneInt(),
		},
		"invalid to": {
			classId: "deadbeef",
			from:    addrs[0],
			amount:  sdk.OneInt(),
		},
		"zero amount": {
			classId: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
			amount:  sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgSend{
			ContractId: tc.classId,
			From:    tc.from.String(),
			To:      tc.to.String(),
			Amount:  tc.amount,
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

func TestMsgOperatorSend(t *testing.T) {
	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId string
		proxy   sdk.AccAddress
		from    sdk.AccAddress
		to      sdk.AccAddress
		amount  sdk.Int
		valid   bool
	}{
		"valid msg": {
			classId: "deadbeef",
			proxy:   addrs[0],
			from:    addrs[1],
			to:      addrs[2],
			amount:  sdk.OneInt(),
			valid:   true,
		},
		"invalid proxy": {
			classId: "deadbeef",
			from:    addrs[1],
			to:      addrs[2],
			amount:  sdk.OneInt(),
		},
		"invalid class id": {
			proxy:   addrs[0],
			from:    addrs[1],
			to:      addrs[2],
			amount:  sdk.OneInt(),
		},
		"empty from": {
			classId: "deadbeef",
			proxy:   addrs[0],
			to:      addrs[1],
			amount:  sdk.OneInt(),
		},
		"invalid to": {
			classId: "deadbeef",
			proxy:   addrs[0],
			from:    addrs[1],
			amount:  sdk.OneInt(),
		},
		"zero amount": {
			classId: "deadbeef",
			proxy:   addrs[0],
			from:    addrs[1],
			to:      addrs[2],
			amount:  sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgOperatorSend{
			ContractId: tc.classId,
			Proxy:   tc.proxy.String(),
			From:    tc.from.String(),
			To:      tc.to.String(),
			Amount:  tc.amount,
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

func TestMsgAuthorizeOperator(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId  string
		approver sdk.AccAddress
		proxy    sdk.AccAddress
		valid    bool
	}{
		"valid msg": {
			classId:  "deadbeef",
			approver: addrs[0],
			proxy:    addrs[1],
			valid:    true,
		},
		"invalid class id": {
			approver: addrs[0],
			proxy:    addrs[1],
		},
		"invalid approver": {
			classId:  "deadbeef",
			proxy:    addrs[1],
		},
		"empty proxy": {
			classId:  "deadbeef",
			approver: addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := token.MsgAuthorizeOperator{
			ContractId:  tc.classId,
			Approver: tc.approver.String(),
			Proxy:    tc.proxy.String(),
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

func TestMsgRevokeOperator(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId  string
		approver sdk.AccAddress
		proxy    sdk.AccAddress
		valid    bool
	}{
		"valid msg": {
			classId:  "deadbeef",
			approver: addrs[0],
			proxy:    addrs[1],
			valid:    true,
		},
		"invalid class id": {
			approver: addrs[0],
			proxy:    addrs[1],
		},
		"invalid approver": {
			classId:  "deadbeef",
			proxy:    addrs[1],
		},
		"empty proxy": {
			classId:  "deadbeef",
			approver: addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := token.MsgRevokeOperator{
			ContractId:  tc.classId,
			Approver: tc.approver.String(),
			Proxy:    tc.proxy.String(),
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
		owner    sdk.AccAddress
		to       sdk.AccAddress
		name     string
		symbol   string
		imageUri string
		meta     string
		decimals int32
		amount   sdk.Int
		valid    bool
	}{
		"valid msg": {
			owner:    addrs[0],
			to:       addrs[1],
			name:     "test",
			symbol:   "TT",
			imageUri: "some URI",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
			valid:    true,
		},
		"invalid owner": {
			to:       addrs[1],
			name:     "test",
			symbol:   "TT",
			imageUri: "some URI",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
		},
		"empty to": {
			owner:    addrs[0],
			name:     "test",
			symbol:   "TT",
			imageUri: "some URI",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
		},
		"empty name": {
			owner:    addrs[0],
			to:       addrs[1],
			symbol:   "TT",
			imageUri: "some URI",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
		},
		"long name": {
			owner:    addrs[0],
			to:       addrs[1],
			name:     "TOO Looooooooooooooog",
			symbol:   "TT",
			imageUri: "some URI",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
			valid:    false,
		},
		"invalid symbol": {
			owner:    addrs[0],
			to:       addrs[1],
			name:     "test",
			imageUri: "some URI",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
		},
		"invalid image uri": {
			owner:    addrs[0],
			to:       addrs[1],
			name:     "test",
			symbol:   "TT",
			imageUri: string(make([]rune, 1001)),
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
		},
		"invalid meta": {
			owner:    addrs[0],
			to:       addrs[1],
			name:     "test",
			symbol:   "TT",
			imageUri: "some URI",
			meta:     string(make([]rune, 1001)),
			decimals: 8,
			amount:   sdk.OneInt(),
		},
		"invalid decimals": {
			owner:    addrs[0],
			to:       addrs[1],
			name:     "test",
			symbol:   "TT",
			imageUri: "some URI",
			meta:     "some meta",
			decimals: 19,
			amount:   sdk.OneInt(),
		},
		"valid supply": {
			owner:    addrs[0],
			to:       addrs[1],
			name:     "test",
			symbol:   "TT",
			imageUri: "some URI",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgIssue{
			Owner:    tc.owner.String(),
			To:       tc.to.String(),
			Name:     tc.name,
			Symbol:   tc.symbol,
			ImageUri: tc.imageUri,
			Meta:     tc.meta,
			Decimals: tc.decimals,
			Amount:   tc.amount,
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
		to      sdk.AccAddress
		amount  sdk.Int
		valid   bool
	}{
		"valid msg": {
			classId: "deadbeef",
			grantee: addrs[0],
			to:      addrs[1],
			amount:  sdk.OneInt(),
			valid:   true,
		},
		"invalid class id": {
			grantee: addrs[0],
			to:      addrs[1],
			amount:  sdk.OneInt(),
		},
		"invalid grantee": {
			classId: "deadbeef",
			to:      addrs[1],
			amount:  sdk.OneInt(),
		},
		"empty to": {
			classId: "deadbeef",
			grantee: addrs[0],
			amount:  sdk.OneInt(),
		},
		"zero amount": {
			classId: "deadbeef",
			grantee: addrs[0],
			to:      addrs[1],
			amount:  sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgMint{
			ContractId: tc.classId,
			From: tc.grantee.String(),
			To:      tc.to.String(),
			Amount:  tc.amount,
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
		from    sdk.AccAddress
		amount  sdk.Int
		valid   bool
	}{
		"valid msg": {
			classId: "deadbeef",
			from:    addrs[0],
			amount:  sdk.OneInt(),
			valid:   true,
		},
		"invalid class id": {
			from:    addrs[0],
			amount:  sdk.OneInt(),
		},
		"invalid from": {
			classId: "deadbeef",
			amount:  sdk.OneInt(),
		},
		"zero amount": {
			classId: "deadbeef",
			from:    addrs[0],
			amount:  sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgBurn{
			ContractId: tc.classId,
			From:    tc.from.String(),
			Amount:  tc.amount,
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

func TestMsgOperatorBurn(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId string
		grantee sdk.AccAddress
		from    sdk.AccAddress
		amount  sdk.Int
		valid   bool
	}{
		"valid msg": {
			classId: "deadbeef",
			grantee: addrs[0],
			from:    addrs[1],
			amount:  sdk.OneInt(),
			valid:   true,
		},
		"invalid class id": {
			grantee: addrs[0],
			from:    addrs[1],
			amount:  sdk.OneInt(),
		},
		"invalid grantee": {
			classId: "deadbeef",
			from:    addrs[1],
			amount:  sdk.OneInt(),
		},
		"empty from": {
			classId: "deadbeef",
			grantee: addrs[0],
			amount:  sdk.OneInt(),
		},
		"zero amount": {
			classId: "deadbeef",
			grantee: addrs[0],
			from:    addrs[1],
			amount:  sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgOperatorBurn{
			ContractId: tc.classId,
			Proxy: tc.grantee.String(),
			From:    tc.from.String(),
			Amount:  tc.amount,
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

	validChange := token.Pair{Field: "name", Value: "New test"}
	testCases := map[string]struct {
		classId string
		grantee sdk.AccAddress
		changes []token.Pair
		valid   bool
	}{
		"valid msg": {
			classId: "deadbeef",
			grantee: addrs[0],
			changes: []token.Pair{validChange},
			valid:   true,
		},
		"invalid class id": {
			grantee: addrs[0],
			changes: []token.Pair{validChange},
		},
		"invalid grantee": {
			classId: "deadbeef",
			changes: []token.Pair{validChange},
		},
		"invalid key of change": {
			classId: "deadbeef",
			grantee: addrs[0],
			changes: []token.Pair{{Value: "tt"}},
		},
		"invalid value of change": {
			classId: "deadbeef",
			grantee: addrs[0],
			changes: []token.Pair{{Field: "symbol"}},
		},
		"empty changes": {
			classId: "deadbeef",
			grantee: addrs[0],
		},
		"duplicated changes": {
			classId: "deadbeef",
			grantee: addrs[0],
			changes: []token.Pair{
				{Field: "name", Value: "hello"},
				{Field: "name", Value: "world"},
			},
		},
	}

	for name, tc := range testCases {
		msg := token.MsgModify{
			ContractId: tc.classId,
			Owner: tc.grantee.String(),
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
		permission  string
		valid   bool
	}{
		"valid msg": {
			classId: "deadbeef",
			granter: addrs[0],
			grantee: addrs[1],
			permission:  token.Permission_Mint.String(),
			valid:   true,
		},
		"invalid class id": {
			granter: addrs[0],
			grantee: addrs[1],
			permission:  token.Permission_Mint.String(),
		},
		"empty granter": {
			classId: "deadbeef",
			grantee: addrs[1],
			permission:  token.Permission_Mint.String(),
		},
		"invalid grantee": {
			classId: "deadbeef",
			granter: addrs[0],
			permission:  token.Permission_Mint.String(),
		},
		"invalid permission": {
			classId: "deadbeef",
			granter: addrs[0],
			grantee: addrs[1],
			permission:  "invalid",
		},
	}

	for name, tc := range testCases {
		msg := token.MsgGrant{
			ContractId: tc.classId,
			From: tc.granter.String(),
			To: tc.grantee.String(),
			Permission:  tc.permission,
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

func TestMsgAbandon(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		classId string
		grantee sdk.AccAddress
		permission  string
		valid   bool
	}{
		"valid msg": {
			classId: "deadbeef",
			grantee: addrs[0],
			permission:  token.Permission_Mint.String(),
			valid:   true,
		},
		"invalid class id": {
			grantee: addrs[0],
			permission:  token.Permission_Mint.String(),
		},
		"invalid grantee": {
			classId: "deadbeef",
			permission:  token.Permission_Mint.String(),
		},
		"invalid permission": {
			classId: "deadbeef",
			grantee: addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := token.MsgAbandon{
			ContractId: tc.classId,
			Grantee: tc.grantee.String(),
			Permission:  tc.permission,
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
