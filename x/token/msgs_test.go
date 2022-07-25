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
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		to         sdk.AccAddress
		amount     sdk.Int
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			amount:     sdk.OneInt(),
			valid:      true,
		},
		"empty from": {
			contractID: "deadbeef",
			to:         addrs[1],
			amount:     sdk.OneInt(),
		},
		"invalid contract id": {
			from:   addrs[0],
			to:     addrs[1],
			amount: sdk.OneInt(),
		},
		"invalid to": {
			contractID: "deadbeef",
			from:       addrs[0],
			amount:     sdk.OneInt(),
		},
		"zero amount": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			amount:     sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgSend{
			ContractId: tc.contractID,
			From:       tc.from.String(),
			To:         tc.to.String(),
			Amount:     tc.amount,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
	}
}

func TestMsgOperatorSend(t *testing.T) {
	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		to         sdk.AccAddress
		amount     sdk.Int
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			operator:   addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			amount:     sdk.OneInt(),
			valid:      true,
		},
		"invalid operator": {
			contractID: "deadbeef",
			from:       addrs[1],
			to:         addrs[2],
			amount:     sdk.OneInt(),
		},
		"invalid contract id": {
			operator: addrs[0],
			from:     addrs[1],
			to:       addrs[2],
			amount:   sdk.OneInt(),
		},
		"empty from": {
			contractID: "deadbeef",
			operator:   addrs[0],
			to:         addrs[1],
			amount:     sdk.OneInt(),
		},
		"invalid to": {
			contractID: "deadbeef",
			operator:   addrs[0],
			from:       addrs[1],
			amount:     sdk.OneInt(),
		},
		"zero amount": {
			contractID: "deadbeef",
			operator:   addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			amount:     sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgOperatorSend{
			ContractId: tc.contractID,
			Operator:   tc.operator.String(),
			From:       tc.from.String(),
			To:         tc.to.String(),
			Amount:     tc.amount,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())
	}
}

func TestMsgTransferFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		proxy      sdk.AccAddress
		from       sdk.AccAddress
		to         sdk.AccAddress
		amount     sdk.Int
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			amount:     sdk.OneInt(),
			valid:      true,
		},
		"invalid proxy": {
			contractID: "deadbeef",
			from:       addrs[1],
			to:         addrs[2],
			amount:     sdk.OneInt(),
		},
		"invalid contract id": {
			proxy:  addrs[0],
			from:   addrs[1],
			to:     addrs[2],
			amount: sdk.OneInt(),
		},
		"empty from": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			to:         addrs[1],
			amount:     sdk.OneInt(),
		},
		"invalid to": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			amount:     sdk.OneInt(),
		},
		"zero amount": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			amount:     sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgTransferFrom{
			ContractId: tc.contractID,
			Proxy:      tc.proxy.String(),
			From:       tc.from.String(),
			To:         tc.to.String(),
			Amount:     tc.amount,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.proxy}, msg.GetSigners())
	}
}

func TestMsgAuthorizeOperator(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		holder     sdk.AccAddress
		operator   sdk.AccAddress
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			holder:     addrs[0],
			operator:   addrs[1],
			valid:      true,
		},
		"invalid contract id": {
			holder:   addrs[0],
			operator: addrs[1],
		},
		"invalid holder": {
			contractID: "deadbeef",
			operator:   addrs[1],
		},
		"empty operator": {
			holder:     addrs[0],
			contractID: "deadbeef",
		},
	}

	for name, tc := range testCases {
		msg := token.MsgAuthorizeOperator{
			ContractId: tc.contractID,
			Holder:     tc.holder.String(),
			Operator:   tc.operator.String(),
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.holder}, msg.GetSigners())
	}
}

func TestMsgRevokeOperator(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		holder     sdk.AccAddress
		operator   sdk.AccAddress
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			holder:     addrs[0],
			operator:   addrs[1],
			valid:      true,
		},
		"invalid contract id": {
			holder:   addrs[0],
			operator: addrs[1],
		},
		"invalid holder": {
			contractID: "deadbeef",
			operator:   addrs[1],
		},
		"empty operator": {
			contractID: "deadbeef",
			holder:     addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := token.MsgRevokeOperator{
			ContractId: tc.contractID,
			Holder:     tc.holder.String(),
			Operator:   tc.operator.String(),
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.holder}, msg.GetSigners())
	}
}

func TestMsgApprove(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		approver   sdk.AccAddress
		proxy      sdk.AccAddress
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			approver:   addrs[0],
			proxy:      addrs[1],
			valid:      true,
		},
		"invalid contract id": {
			approver: addrs[0],
			proxy:    addrs[1],
		},
		"invalid approver": {
			contractID: "deadbeef",
			proxy:      addrs[1],
		},
		"empty proxy": {
			contractID: "deadbeef",
			approver:   addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := token.MsgApprove{
			ContractId: tc.contractID,
			Approver:   tc.approver.String(),
			Proxy:      tc.proxy.String(),
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.approver}, msg.GetSigners())
	}
}

func TestMsgIssue(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
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

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.owner}, msg.GetSigners())
	}
}

func TestMsgMint(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		to         sdk.AccAddress
		amount     sdk.Int
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			to:         addrs[1],
			amount:     sdk.OneInt(),
			valid:      true,
		},
		"invalid contract id": {
			grantee: addrs[0],
			to:      addrs[1],
			amount:  sdk.OneInt(),
		},
		"invalid grantee": {
			contractID: "deadbeef",
			to:         addrs[1],
			amount:     sdk.OneInt(),
		},
		"empty to": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			amount:     sdk.OneInt(),
		},
		"zero amount": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			to:         addrs[1],
			amount:     sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgMint{
			ContractId: tc.contractID,
			From:       tc.grantee.String(),
			To:         tc.to.String(),
			Amount:     tc.amount,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.grantee}, msg.GetSigners())
	}
}

func TestMsgBurn(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		amount     sdk.Int
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:       addrs[0],
			amount:     sdk.OneInt(),
			valid:      true,
		},
		"invalid contract id": {
			from:   addrs[0],
			amount: sdk.OneInt(),
		},
		"invalid from": {
			contractID: "deadbeef",
			amount:     sdk.OneInt(),
		},
		"zero amount": {
			contractID: "deadbeef",
			from:       addrs[0],
			amount:     sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgBurn{
			ContractId: tc.contractID,
			From:       tc.from.String(),
			Amount:     tc.amount,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
	}
}

func TestMsgOperatorBurn(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		amount     sdk.Int
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			operator:   addrs[0],
			from:       addrs[1],
			amount:     sdk.OneInt(),
			valid:      true,
		},
		"invalid contract id": {
			operator: addrs[0],
			from:     addrs[1],
			amount:   sdk.OneInt(),
		},
		"invalid operator": {
			contractID: "deadbeef",
			from:       addrs[1],
			amount:     sdk.OneInt(),
		},
		"empty from": {
			contractID: "deadbeef",
			operator:   addrs[0],
			amount:     sdk.OneInt(),
		},
		"zero amount": {
			contractID: "deadbeef",
			operator:   addrs[0],
			from:       addrs[1],
			amount:     sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgOperatorBurn{
			ContractId: tc.contractID,
			Operator:   tc.operator.String(),
			From:       tc.from.String(),
			Amount:     tc.amount,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())
	}
}

func TestMsgBurnFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		from       sdk.AccAddress
		amount     sdk.Int
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			from:       addrs[1],
			amount:     sdk.OneInt(),
			valid:      true,
		},
		"invalid contract id": {
			grantee: addrs[0],
			from:    addrs[1],
			amount:  sdk.OneInt(),
		},
		"invalid grantee": {
			contractID: "deadbeef",
			from:       addrs[1],
			amount:     sdk.OneInt(),
		},
		"empty from": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			amount:     sdk.OneInt(),
		},
		"zero amount": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			from:       addrs[1],
			amount:     sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := token.MsgBurnFrom{
			ContractId: tc.contractID,
			Proxy:      tc.grantee.String(),
			From:       tc.from.String(),
			Amount:     tc.amount,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.grantee}, msg.GetSigners())
	}
}

func TestMsgModify(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	validChange := token.Pair{Field: token.AttributeKeyName.String(), Value: "New test"}
	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		changes    []token.Pair
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			changes:    []token.Pair{validChange},
			valid:      true,
		},
		"invalid contract id": {
			grantee: addrs[0],
			changes: []token.Pair{validChange},
		},
		"invalid grantee": {
			contractID: "deadbeef",
			changes:    []token.Pair{validChange},
		},
		"invalid key of change": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			changes:    []token.Pair{{Value: "tt"}},
		},
		"invalid value of change": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			changes:    []token.Pair{{Field: "symbol"}},
		},
		"empty changes": {
			contractID: "deadbeef",
			grantee:    addrs[0],
		},
		"duplicated changes": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			changes: []token.Pair{
				{Field: token.AttributeKeyName.String(), Value: "hello"},
				{Field: token.AttributeKeyName.String(), Value: "world"},
			},
		},
	}

	for name, tc := range testCases {
		msg := token.MsgModify{
			ContractId: tc.contractID,
			Owner:      tc.grantee.String(),
			Changes:    tc.changes,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.grantee}, msg.GetSigners())
	}
}

func TestMsgGrant(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		granter    sdk.AccAddress
		grantee    sdk.AccAddress
		permission token.Permission
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			granter:    addrs[0],
			grantee:    addrs[1],
			permission: token.PermissionMint,
			valid:      true,
		},
		"invalid contract id": {
			granter:    addrs[0],
			grantee:    addrs[1],
			permission: token.PermissionMint,
		},
		"empty granter": {
			contractID: "deadbeef",
			grantee:    addrs[1],
			permission: token.PermissionMint,
		},
		"invalid grantee": {
			contractID: "deadbeef",
			granter:    addrs[0],
			permission: token.PermissionMint,
		},
		"invalid permission": {
			contractID: "deadbeef",
			granter:    addrs[0],
			grantee:    addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := token.MsgGrant{
			ContractId: tc.contractID,
			Granter:    tc.granter.String(),
			Grantee:    tc.grantee.String(),
			Permission: tc.permission,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.granter}, msg.GetSigners())
	}
}

func TestMsgAbandon(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		permission token.Permission
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			permission: token.PermissionMint,
			valid:      true,
		},
		"invalid contract id": {
			grantee:    addrs[0],
			permission: token.PermissionMint,
		},
		"invalid grantee": {
			contractID: "deadbeef",
			permission: token.PermissionMint,
		},
		"invalid permission": {
			contractID: "deadbeef",
			grantee:    addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := token.MsgAbandon{
			ContractId: tc.contractID,
			Grantee:    tc.grantee.String(),
			Permission: tc.permission,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.grantee}, msg.GetSigners())
	}
}

func TestMsgGrantPermission(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		to         sdk.AccAddress
		permission string
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			permission: token.LegacyPermissionMint.String(),
			valid:      true,
		},
		"invalid contract id": {
			from:       addrs[0],
			to:         addrs[1],
			permission: token.LegacyPermissionMint.String(),
		},
		"empty from": {
			contractID: "deadbeef",
			to:         addrs[1],
			permission: token.LegacyPermissionMint.String(),
		},
		"invalid to": {
			contractID: "deadbeef",
			from:       addrs[0],
			permission: token.LegacyPermissionMint.String(),
		},
		"invalid permission": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := token.MsgGrantPermission{
			ContractId: tc.contractID,
			From:       tc.from.String(),
			To:         tc.to.String(),
			Permission: tc.permission,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
	}
}

func TestMsgRevokePermission(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		permission string
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:       addrs[0],
			permission: token.LegacyPermissionMint.String(),
			valid:      true,
		},
		"invalid contract id": {
			from:       addrs[0],
			permission: token.LegacyPermissionMint.String(),
		},
		"invalid from": {
			contractID: "deadbeef",
			permission: token.LegacyPermissionMint.String(),
		},
		"invalid permission": {
			contractID: "deadbeef",
			from:       addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := token.MsgRevokePermission{
			ContractId: tc.contractID,
			From:       tc.from.String(),
			Permission: tc.permission,
		}

		err := msg.ValidateBasic()
		if !tc.valid {
			require.Error(t, err, name)
			return
		}
		require.NoError(t, err, name)

		require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
	}
}
