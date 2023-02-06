package token_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/auth/legacy/legacytx"
	"github.com/line/lbm-sdk/x/token"
	"github.com/line/lbm-sdk/x/token/class"
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
		err        error
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			amount:     sdk.OneInt(),
		},
		"invalid from": {
			contractID: "deadbeef",
			to:         addrs[1],
			amount:     sdk.OneInt(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid contract id": {
			from:   addrs[0],
			to:     addrs[1],
			amount: sdk.OneInt(),
			err:    class.ErrInvalidContractID,
		},
		"invalid to": {
			contractID: "deadbeef",
			from:       addrs[0],
			amount:     sdk.OneInt(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid amount": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			amount:     sdk.ZeroInt(),
			err:        token.ErrInvalidAmount,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := token.MsgSend{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				To:         tc.to.String(),
				Amount:     tc.amount,
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
		})
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
		err        error
	}{
		"valid msg": {
			contractID: "deadbeef",
			operator:   addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			amount:     sdk.OneInt(),
		},
		"invalid operator": {
			contractID: "deadbeef",
			from:       addrs[1],
			to:         addrs[2],
			amount:     sdk.OneInt(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid contract id": {
			operator: addrs[0],
			from:     addrs[1],
			to:       addrs[2],
			amount:   sdk.OneInt(),
			err:      class.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: "deadbeef",
			operator:   addrs[0],
			to:         addrs[1],
			amount:     sdk.OneInt(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			contractID: "deadbeef",
			operator:   addrs[0],
			from:       addrs[1],
			amount:     sdk.OneInt(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid amount": {
			contractID: "deadbeef",
			operator:   addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			amount:     sdk.ZeroInt(),
			err:        token.ErrInvalidAmount,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := token.MsgOperatorSend{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				To:         tc.to.String(),
				Amount:     tc.amount,
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())
		})
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
		err        error
	}{
		"valid msg": {
			contractID: "deadbeef",
			holder:     addrs[0],
			operator:   addrs[1],
		},
		"invalid contract id": {
			holder:   addrs[0],
			operator: addrs[1],
			err:      class.ErrInvalidContractID,
		},
		"invalid holder": {
			contractID: "deadbeef",
			operator:   addrs[1],
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid operator": {
			contractID: "deadbeef",
			holder:     addrs[0],
			err:        sdkerrors.ErrInvalidAddress,
		},
		"operator and holder should be different": {
			contractID: "deadbeef",
			holder:     addrs[0],
			operator:   addrs[0],
			err:        token.ErrApproverProxySame,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := token.MsgRevokeOperator{
				ContractId: tc.contractID,
				Holder:     tc.holder.String(),
				Operator:   tc.operator.String(),
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.holder}, msg.GetSigners())
		})
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
		err        error
	}{
		"valid msg": {
			contractID: "deadbeef",
			holder:     addrs[0],
			operator:   addrs[1],
		},
		"invalid contract id": {
			holder:   addrs[0],
			operator: addrs[1],
			err:      class.ErrInvalidContractID,
		},
		"invalid holder": {
			contractID: "deadbeef",
			operator:   addrs[1],
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty operator": {
			contractID: "deadbeef",
			holder:     addrs[0],
			err:        sdkerrors.ErrInvalidAddress,
		},
		"proxy and approver should be different": {
			contractID: "deadbeef",
			holder:     addrs[0],
			operator:   addrs[0],
			err:        token.ErrApproverProxySame,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := token.MsgAuthorizeOperator{
				ContractId: tc.contractID,
				Holder:     tc.holder.String(),
				Operator:   tc.operator.String(),
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.holder}, msg.GetSigners())
		})
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
		err      error
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
		},
		"invalid owner": {
			to:       addrs[1],
			name:     "test",
			symbol:   "TT",
			imageUri: "some URI",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
			err:      sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			owner:    addrs[0],
			name:     "test",
			symbol:   "TT",
			imageUri: "some URI",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
			err:      sdkerrors.ErrInvalidAddress,
		},
		"empty name": {
			owner:    addrs[0],
			to:       addrs[1],
			symbol:   "TT",
			imageUri: "some URI",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
			err:      token.ErrInvalidTokenName,
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
			err:      token.ErrInvalidNameLength,
		},
		"invalid symbol": {
			owner:    addrs[0],
			to:       addrs[1],
			name:     "test",
			imageUri: "some URI",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
			err:      token.ErrInvalidTokenSymbol,
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
			err:      token.ErrInvalidImageURILength,
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
			err:      token.ErrInvalidMetaLength,
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
			err:      token.ErrInvalidTokenDecimals,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := token.MsgIssue{
				Owner:    tc.owner.String(),
				To:       tc.to.String(),
				Name:     tc.name,
				Symbol:   tc.symbol,
				Uri:      tc.imageUri,
				Meta:     tc.meta,
				Decimals: tc.decimals,
				Amount:   tc.amount,
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.owner}, msg.GetSigners())
		})
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
		err        error
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			to:         addrs[1],
			amount:     sdk.OneInt(),
		},
		"invalid contract id": {
			grantee: addrs[0],
			to:      addrs[1],
			amount:  sdk.OneInt(),
			err:     class.ErrInvalidContractID,
		},
		"invalid grantee": {
			contractID: "deadbeef",
			to:         addrs[1],
			amount:     sdk.OneInt(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			amount:     sdk.OneInt(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid amount": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			to:         addrs[1],
			amount:     sdk.ZeroInt(),
			err:        token.ErrInvalidAmount,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := token.MsgMint{
				ContractId: tc.contractID,
				From:       tc.grantee.String(),
				To:         tc.to.String(),
				Amount:     tc.amount,
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.grantee}, msg.GetSigners())
		})
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
		err        error
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:       addrs[0],
			amount:     sdk.OneInt(),
		},
		"invalid contract id": {
			from:   addrs[0],
			amount: sdk.OneInt(),
			err:    class.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: "deadbeef",
			amount:     sdk.OneInt(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid amount": {
			contractID: "deadbeef",
			from:       addrs[0],
			amount:     sdk.ZeroInt(),
			err:        token.ErrInvalidAmount,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := token.MsgBurn{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				Amount:     tc.amount,
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
		})
	}
}

func TestMsgOperatorBurn(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		from       sdk.AccAddress
		amount     sdk.Int
		err        error
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			from:       addrs[1],
			amount:     sdk.OneInt(),
		},
		"invalid contract id": {
			grantee: addrs[0],
			from:    addrs[1],
			amount:  sdk.OneInt(),
			err:     class.ErrInvalidContractID,
		},
		"invalid grantee": {
			contractID: "deadbeef",
			from:       addrs[1],
			amount:     sdk.OneInt(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid from": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			amount:     sdk.OneInt(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid amount": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			from:       addrs[1],
			amount:     sdk.ZeroInt(),
			err:        token.ErrInvalidAmount,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := token.MsgOperatorBurn{
				ContractId: tc.contractID,
				Operator:   tc.grantee.String(),
				From:       tc.from.String(),
				Amount:     tc.amount,
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.grantee}, msg.GetSigners())
		})
	}
}

func TestMsgModify(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	validChange := token.Attribute{Key: token.AttributeKeyName.String(), Value: "New test"}
	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		changes    []token.Attribute
		err        error
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			changes:    []token.Attribute{validChange},
		},
		"invalid contract id": {
			grantee: addrs[0],
			changes: []token.Attribute{validChange},
			err:     class.ErrInvalidContractID,
		},
		"invalid grantee": {
			contractID: "deadbeef",
			changes:    []token.Attribute{validChange},
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid key of change": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			changes:    []token.Attribute{{Key: strings.ToUpper(token.AttributeKeyName.String()), Value: "tt"}},
			err:        token.ErrInvalidChangesField,
		},
		"invalid value of change": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			changes:    []token.Attribute{{Key: token.AttributeKeyName.String(), Value: string(make([]rune, 21))}},
			err:        token.ErrInvalidNameLength,
		},
		"empty changes": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			err:        token.ErrEmptyChanges,
		},
		"duplicated changes": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			changes: []token.Attribute{
				{Key: token.AttributeKeyName.String(), Value: "hello"},
				{Key: token.AttributeKeyName.String(), Value: "world"},
			},
			err: token.ErrDuplicateChangesField,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := token.MsgModify{
				ContractId: tc.contractID,
				Owner:      tc.grantee.String(),
				Changes:    tc.changes,
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.grantee}, msg.GetSigners())
		})
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
		err        error
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			permission: token.LegacyPermissionMint.String(),
		},
		"invalid contract id": {
			from:       addrs[0],
			to:         addrs[1],
			permission: token.LegacyPermissionMint.String(),
			err:        class.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: "deadbeef",
			to:         addrs[1],
			permission: token.LegacyPermissionMint.String(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			contractID: "deadbeef",
			from:       addrs[0],
			permission: token.LegacyPermissionMint.String(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid permission": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			err:        sdkerrors.ErrInvalidPermission,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := token.MsgGrantPermission{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				To:         tc.to.String(),
				Permission: tc.permission,
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
		})
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
		err        error
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:       addrs[0],
			permission: token.LegacyPermissionMint.String(),
		},
		"invalid contract id": {
			from:       addrs[0],
			permission: token.LegacyPermissionMint.String(),
			err:        class.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: "deadbeef",
			permission: token.LegacyPermissionMint.String(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid permission": {
			contractID: "deadbeef",
			from:       addrs[0],
			err:        sdkerrors.ErrInvalidPermission,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := token.MsgRevokePermission{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				Permission: tc.permission,
			}

			err := msg.ValidateBasic()
			require.ErrorIs(t, err, tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
		})
	}
}

func TestAminoJSON(t *testing.T) {
	tx := legacytx.StdTx{}
	var contractId = "deadbeef"

	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		msg          legacytx.LegacyMsg
		expectedType string
		expected     string
	}{
		"MsgSend": {
			&token.MsgSend{
				ContractId: contractId,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				Amount:     sdk.OneInt(),
			},
			"/lbm.token.v1.MsgSend",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSend\",\"value\":{\"amount\":\"1\",\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgOperatorSend": {
			&token.MsgOperatorSend{
				ContractId: contractId,
				Operator:   addrs[0].String(),
				From:       addrs[1].String(),
				To:         addrs[2].String(),
				Amount:     sdk.OneInt(),
			},
			"/lbm.token.v1.MsgOperatorSend",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgOperatorSend\",\"value\":{\"amount\":\"1\",\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"operator\":\"%s\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String(), addrs[2].String()),
		},
		"MsgRevokeOperator": {
			&token.MsgRevokeOperator{
				ContractId: contractId,
				Holder:     addrs[0].String(),
				Operator:   addrs[1].String(),
			},
			"/lbm.token.v1.MsgRevokeOperator",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/token/MsgRevokeOperator\",\"value\":{\"contract_id\":\"deadbeef\",\"holder\":\"%s\",\"operator\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgAuthorizeOperator": {
			&token.MsgAuthorizeOperator{
				ContractId: contractId,
				Holder:     addrs[0].String(),
				Operator:   addrs[1].String(),
			},
			"/lbm.token.v1.MsgAuthorizeOperator",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/token/MsgAuthorizeOperator\",\"value\":{\"contract_id\":\"deadbeef\",\"holder\":\"%s\",\"operator\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgIssue": {
			&token.MsgIssue{
				Name:     "Test Name",
				Symbol:   "LN",
				Uri:      "http://image.url",
				Meta:     "This is test",
				Decimals: 6,
				Mintable: false,
				Owner:    addrs[0].String(),
				To:       addrs[1].String(),
				Amount:   sdk.NewInt(1000000),
			},
			"/lbm.token.v1.MsgIssue",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgIssue\",\"value\":{\"amount\":\"1000000\",\"decimals\":6,\"meta\":\"This is test\",\"name\":\"Test Name\",\"owner\":\"%s\",\"symbol\":\"LN\",\"to\":\"%s\",\"uri\":\"http://image.url\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgGrantPermission": {
			&token.MsgGrantPermission{
				ContractId: contractId,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				Permission: token.LegacyPermissionMint.String(),
			},
			"/lbm.token.v1.MsgGrantPermission",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/token/MsgGrantPermission\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"permission\":\"mint\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgRevokePermission": {
			&token.MsgRevokePermission{
				ContractId: contractId,
				From:       addrs[0].String(),
				Permission: token.LegacyPermissionMint.String(),
			},
			"/lbm.token.v1.MsgRevokePermission",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/token/MsgRevokePermission\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"permission\":\"mint\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgMint": {
			&token.MsgMint{
				ContractId: contractId,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				Amount:     sdk.NewInt(1000000),
			},
			"/lbm.token.v1.MsgMint",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgMint\",\"value\":{\"amount\":\"1000000\",\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgBurn": {
			&token.MsgBurn{
				ContractId: contractId,
				From:       addrs[0].String(),
				Amount:     sdk.Int{},
			},
			"/lbm.token.v1.MsgBurn",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgBurn\",\"value\":{\"amount\":\"0\",\"contract_id\":\"deadbeef\",\"from\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgOperatorBurn": {
			&token.MsgOperatorBurn{
				ContractId: contractId,
				Operator:   addrs[0].String(),
				From:       addrs[1].String(),
				Amount:     sdk.NewInt(1000000),
			},
			"/lbm.token.v1.MsgOperatorBurn",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgOperatorBurn\",\"value\":{\"amount\":\"1000000\",\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"operator\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String()),
		},
		"MsgModify": {
			&token.MsgModify{
				ContractId: contractId,
				Owner:      addrs[0].String(),
				Changes:    []token.Attribute{{Key: token.AttributeKeyName.String(), Value: "New test"}},
			},
			"/lbm.token.v1.MsgModify",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/token/MsgModify\",\"value\":{\"changes\":[{\"key\":\"name\",\"value\":\"New test\"}],\"contract_id\":\"deadbeef\",\"owner\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			tx.Msgs = []sdk.Msg{tc.msg}
			require.Equal(t, token.RouterKey, tc.msg.Route())
			require.Equal(t, tc.expectedType, tc.msg.Type())
			require.Equal(t, tc.expected, string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{tc.msg}, "memo")))
		})
	}
}
