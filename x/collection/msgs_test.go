package collection_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

func TestMsgSend(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	amount := collection.NewCoins(
		collection.NewFTCoin("00bab10c", sdk.OneInt()),
	)

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		to         sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			amount:     amount,
			valid:      true,
		},
		"empty from": {
			contractID: contractID,
			to:         addrs[1],
			amount:     amount,
		},
		"invalid contract id": {
			from:   addrs[0],
			to:     addrs[1],
			amount: amount,
		},
		"invalid to": {
			contractID: contractID,
			from:       addrs[0],
			amount:     amount,
		},
		"empty amount": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
		},
		"invalid token id": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			amount: []collection.Coin{{
				Amount: sdk.OneInt(),
			}},
		},
		"duplicate token ids": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			amount:     []collection.Coin{amount[0], amount[0]},
		},
		"invalid amount": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			amount: []collection.Coin{{
				TokenId: amount[0].TokenId,
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgSend{
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

	contractID := "deadbeef"
	amount := collection.NewCoins(
		collection.NewFTCoin("00bab10c", sdk.OneInt()),
	)

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		to         sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			amount:     amount,
			valid:      true,
		},
		"invalid operator": {
			contractID: contractID,
			from:       addrs[1],
			to:         addrs[2],
			amount:     amount,
		},
		"invalid contract id": {
			operator: addrs[0],
			from:     addrs[1],
			to:       addrs[2],
			amount:   amount,
		},
		"empty from": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			amount:     amount,
		},
		"invalid to": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			amount:     amount,
		},
		"empty amount": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			to:         addrs[2],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgOperatorSend{
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

func TestMsgTransferFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(
		collection.NewFTCoin("00bab10c", sdk.OneInt()),
	)

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		to         sdk.AccAddress
		amount     []collection.Coin
		valid      bool
		panic      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			amount:     amount,
			valid:      true,
		},
		"empty from": {
			contractID: "deadbeef",
			to:         addrs[1],
			amount:     amount,
		},
		"invalid contract id": {
			from:   addrs[0],
			to:     addrs[1],
			amount: amount,
		},
		"invalid to": {
			contractID: "deadbeef",
			from:       addrs[0],
			amount:     amount,
		},
		"nil amount": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			amount: []collection.Coin{{
				TokenId: collection.NewFTID("00bab10c"),
			}},
			panic: true,
		},
		"zero amount": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			amount: []collection.Coin{{
				TokenId: collection.NewFTID("00bab10c"),
				Amount:  sdk.ZeroInt(),
			}},
		},
		"invalid token id": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			amount: []collection.Coin{{
				Amount: sdk.OneInt(),
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgTransferFT{
			ContractId: tc.contractID,
			From:       tc.from.String(),
			To:         tc.to.String(),
			Amount:     tc.amount,
		}

		if tc.panic {
			require.Panics(t, func() { msg.ValidateBasic() }, name)
			continue
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

func TestMsgTransferFTFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(
		collection.NewFTCoin("00bab10c", sdk.OneInt()),
	)

	testCases := map[string]struct {
		contractID string
		proxy      sdk.AccAddress
		from       sdk.AccAddress
		to         sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			amount:     amount,
			valid:      true,
		},
		"invalid proxy": {
			contractID: "deadbeef",
			from:       addrs[1],
			to:         addrs[2],
			amount:     amount,
		},
		"invalid contract id": {
			proxy:  addrs[0],
			from:   addrs[1],
			to:     addrs[2],
			amount: amount,
		},
		"empty from": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			to:         addrs[1],
			amount:     amount,
		},
		"invalid to": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			amount:     amount,
		},
		"invalid amount": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			amount: []collection.Coin{{
				Amount: sdk.OneInt(),
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgTransferFTFrom{
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

func TestMsgTransferNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ids := []string{collection.NewNFTID("deadbeef", 1)}

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		to         sdk.AccAddress
		ids        []string
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			ids:        ids,
			valid:      true,
		},
		"empty from": {
			contractID: "deadbeef",
			to:         addrs[1],
			ids:        ids,
		},
		"invalid contract id": {
			from: addrs[0],
			to:   addrs[1],
			ids:  ids,
		},
		"invalid to": {
			contractID: "deadbeef",
			from:       addrs[0],
			ids:        ids,
		},
		"empty token ids": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
		},
		"invalid token ids": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
			ids:        []string{""},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgTransferNFT{
			ContractId: tc.contractID,
			From:       tc.from.String(),
			To:         tc.to.String(),
			TokenIds:   tc.ids,
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

func TestMsgTransferNFTFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ids := []string{collection.NewNFTID("deadbeef", 1)}

	testCases := map[string]struct {
		contractID string
		proxy      sdk.AccAddress
		from       sdk.AccAddress
		to         sdk.AccAddress
		ids        []string
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			ids:        ids,
			valid:      true,
		},
		"invalid proxy": {
			contractID: "deadbeef",
			from:       addrs[1],
			to:         addrs[2],
			ids:        ids,
		},
		"invalid contract id": {
			proxy: addrs[0],
			from:  addrs[1],
			to:    addrs[2],
			ids:   ids,
		},
		"empty from": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			to:         addrs[1],
			ids:        ids,
		},
		"invalid to": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			ids:        ids,
		},
		"empty ids": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			to:         addrs[2],
		},
		"invalid id": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			ids:        []string{""},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgTransferNFTFrom{
			ContractId: tc.contractID,
			Proxy:      tc.proxy.String(),
			From:       tc.from.String(),
			To:         tc.to.String(),
			TokenIds:   tc.ids,
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

	contractID := "deadbeef"
	testCases := map[string]struct {
		contractID string
		holder     sdk.AccAddress
		operator   sdk.AccAddress
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			holder:     addrs[0],
			operator:   addrs[1],
			valid:      true,
		},
		"invalid contract id": {
			holder:   addrs[0],
			operator: addrs[1],
		},
		"invalid holder": {
			contractID: contractID,
			operator:   addrs[1],
		},
		"empty operator": {
			contractID: contractID,
			holder:     addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgAuthorizeOperator{
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

	contractID := "deadbeef"
	testCases := map[string]struct {
		contractID string
		holder     sdk.AccAddress
		operator   sdk.AccAddress
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			holder:     addrs[0],
			operator:   addrs[1],
			valid:      true,
		},
		"invalid contract id": {
			holder:   addrs[0],
			operator: addrs[1],
		},
		"invalid holder": {
			contractID: contractID,
			operator:   addrs[1],
		},
		"empty operator": {
			contractID: contractID,
			holder:     addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgRevokeOperator{
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
		msg := collection.MsgApprove{
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

func TestMsgDisapprove(t *testing.T) {
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
		msg := collection.MsgDisapprove{
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

func TestMsgCreateContract(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	name := "tibetian fox"
	uri := "file:///tibetian_fox.png"
	meta := "Tibetian fox"
	testCases := map[string]struct {
		owner      sdk.AccAddress
		name       string
		baseImgURI string
		meta       string
		valid      bool
	}{
		"valid msg": {
			owner:      addrs[0],
			name:       name,
			baseImgURI: uri,
			meta:       meta,
			valid:      true,
		},
		"invalid owner": {
			name:       name,
			baseImgURI: uri,
			meta:       meta,
		},
		"long name": {
			owner:      addrs[0],
			name:       string(make([]rune, 21)),
			baseImgURI: uri,
			meta:       meta,
		},
		"invalid base image uri": {
			owner:      addrs[0],
			name:       name,
			baseImgURI: string(make([]rune, 1001)),
			meta:       meta,
		},
		"invalid meta": {
			owner:      addrs[0],
			name:       name,
			baseImgURI: uri,
			meta:       string(make([]rune, 1001)),
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgCreateContract{
			Owner:      tc.owner.String(),
			Name:       tc.name,
			BaseImgUri: tc.baseImgURI,
			Meta:       tc.meta,
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

func TestMsgIssueFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	name := "tibetian fox"
	meta := "Tibetian Fox"
	decimals := int32(8)
	testCases := map[string]struct {
		contractID string
		owner      sdk.AccAddress
		to         sdk.AccAddress
		name       string
		meta       string
		decimals   int32
		mintable   bool
		amount     sdk.Int
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			name:       name,
			meta:       meta,
			decimals:   decimals,
			amount:     sdk.OneInt(),
			valid:      true,
		},
		"invalid contract id": {
			owner:    addrs[0],
			to:       addrs[1],
			name:     name,
			meta:     meta,
			decimals: decimals,
			amount:   sdk.OneInt(),
		},
		"invalid owner": {
			contractID: contractID,
			to:         addrs[1],
			name:       name,
			meta:       meta,
			decimals:   decimals,
			amount:     sdk.OneInt(),
		},
		"empty to": {
			contractID: contractID,
			owner:      addrs[0],
			name:       name,
			meta:       meta,
			decimals:   decimals,
			amount:     sdk.OneInt(),
		},
		"empty name": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			meta:       meta,
			decimals:   decimals,
			amount:     sdk.OneInt(),
		},
		"long name": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			name:       string(make([]rune, 21)),
			meta:       meta,
			decimals:   decimals,
			amount:     sdk.OneInt(),
			valid:      false,
		},
		"invalid meta": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			name:       name,
			meta:       string(make([]rune, 1001)),
			decimals:   decimals,
			amount:     sdk.OneInt(),
		},
		"invalid decimals": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			name:       name,
			meta:       meta,
			decimals:   19,
			amount:     sdk.OneInt(),
		},
		"daphne compat": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			name:       name,
			meta:       meta,
			amount:     sdk.OneInt(),
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgIssueFT{
			ContractId: tc.contractID,
			Owner:      tc.owner.String(),
			To:         tc.to.String(),
			Name:       tc.name,
			Meta:       tc.meta,
			Decimals:   tc.decimals,
			Amount:     tc.amount,
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

func TestMsgIssueNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	name := "tibetian fox"
	meta := "Tibetian Fox"
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		name       string
		meta       string
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			name:       name,
			meta:       meta,
			valid:      true,
		},
		"invalid contract id": {
			operator: addrs[0],
			name:     name,
			meta:     meta,
		},
		"invalid operator": {
			contractID: contractID,
			name:       name,
			meta:       meta,
		},
		"long name": {
			contractID: contractID,
			operator:   addrs[0],
			name:       string(make([]rune, 21)),
			meta:       meta,
		},
		"invalid meta": {
			contractID: contractID,
			operator:   addrs[0],
			name:       name,
			meta:       string(make([]rune, 1001)),
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgIssueNFT{
			ContractId: tc.contractID,
			Owner:      tc.operator.String(),
			Name:       tc.name,
			Meta:       tc.meta,
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

func TestMsgCreateFTClass(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	name := "tibetian fox"
	meta := "Tibetian Fox"
	decimals := int32(8)
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		name       string
		meta       string
		decimals   int32
		to         sdk.AccAddress
		supply     sdk.Int
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			name:       name,
			meta:       meta,
			decimals:   decimals,
			supply:     sdk.ZeroInt(),
			valid:      true,
		},
		"valid msg with supply": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			name:       name,
			meta:       meta,
			decimals:   decimals,
			supply:     sdk.OneInt(),
			valid:      true,
		},
		"invalid contract id": {
			operator: addrs[0],
			name:     name,
			meta:     meta,
			decimals: decimals,
			supply:   sdk.ZeroInt(),
		},
		"invalid operator": {
			contractID: contractID,
			name:       name,
			meta:       meta,
			decimals:   decimals,
			supply:     sdk.ZeroInt(),
		},
		"long name": {
			contractID: contractID,
			operator:   addrs[0],
			name:       string(make([]rune, 21)),
			meta:       meta,
			decimals:   decimals,
			supply:     sdk.ZeroInt(),
			valid:      false,
		},
		"invalid meta": {
			contractID: contractID,
			operator:   addrs[0],
			name:       name,
			meta:       string(make([]rune, 1001)),
			decimals:   decimals,
			supply:     sdk.ZeroInt(),
		},
		"invalid decimals": {
			contractID: contractID,
			operator:   addrs[0],
			name:       name,
			meta:       meta,
			decimals:   19,
			supply:     sdk.ZeroInt(),
		},
		"positive supply with invalid to": {
			contractID: contractID,
			operator:   addrs[0],
			name:       name,
			meta:       meta,
			decimals:   decimals,
			supply:     sdk.OneInt(),
		},
		"invalid supply": {
			contractID: contractID,
			operator:   addrs[0],
			name:       name,
			meta:       meta,
			decimals:   decimals,
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgCreateFTClass{
			ContractId: tc.contractID,
			Operator:   tc.operator.String(),
			To:         tc.to.String(),
			Name:       tc.name,
			Meta:       tc.meta,
			Decimals:   tc.decimals,
			Supply:     tc.supply,
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

func TestMsgCreateNFTClass(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	name := "tibetian fox"
	meta := "Tibetian Fox"
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		name       string
		meta       string
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			name:       name,
			meta:       meta,
			valid:      true,
		},
		"invalid contract id": {
			operator: addrs[0],
			name:     name,
			meta:     meta,
		},
		"invalid operator": {
			contractID: contractID,
			name:       name,
			meta:       meta,
		},
		"long name": {
			contractID: contractID,
			operator:   addrs[0],
			name:       string(make([]rune, 21)),
			meta:       meta,
		},
		"invalid meta": {
			contractID: contractID,
			operator:   addrs[0],
			name:       name,
			meta:       string(make([]rune, 1001)),
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgCreateNFTClass{
			ContractId: tc.contractID,
			Operator:   tc.operator.String(),
			Name:       tc.name,
			Meta:       tc.meta,
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

func TestMsgMintFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(
		collection.NewFTCoin("00bab10c", sdk.OneInt()),
	)

	contractID := "deadbeef"
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		to         sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			amount:     amount,
			valid:      true,
		},
		"invalid contract id": {
			operator: addrs[0],
			to:       addrs[1],
			amount:   amount,
		},
		"invalid operator": {
			contractID: contractID,
			to:         addrs[1],
			amount:     amount,
		},
		"empty to": {
			contractID: contractID,
			operator:   addrs[0],
			amount:     amount,
		},
		"invalid token id": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			amount: []collection.Coin{{
				Amount: sdk.OneInt(),
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgMintFT{
			ContractId: tc.contractID,
			From:       tc.operator.String(),
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

func TestMsgMintNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	params := []collection.MintNFTParam{{
		TokenType: "deadbeef",
		Name:      "tibetian fox",
		Meta:      "Tibetian Fox",
	}}
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		to         sdk.AccAddress
		params     []collection.MintNFTParam
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			operator:   addrs[0],
			to:         addrs[1],
			params:     params,
			valid:      true,
		},
		"invalid contract id": {
			operator: addrs[0],
			to:       addrs[1],
			params:   params,
		},
		"invalid operator": {
			contractID: "deadbeef",
			to:         addrs[1],
			params:     params,
		},
		"empty to": {
			contractID: "deadbeef",
			operator:   addrs[0],
			params:     params,
		},
		"empty params": {
			contractID: "deadbeef",
			operator:   addrs[0],
			to:         addrs[1],
		},
		"param of invalid token type": {
			contractID: "deadbeef",
			operator:   addrs[0],
			to:         addrs[1],
			params:     []collection.MintNFTParam{{}},
		},
		"param of invalid name": {
			contractID: "deadbeef",
			operator:   addrs[0],
			to:         addrs[1],
			params: []collection.MintNFTParam{{
				TokenType: "deadbeef",
				Name:      string(make([]rune, 21)),
			}},
		},
		"param of invalid meta": {
			contractID: "deadbeef",
			operator:   addrs[0],
			to:         addrs[1],
			params: []collection.MintNFTParam{{
				TokenType: "deadbeef",
				Meta:      string(make([]rune, 1001)),
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgMintNFT{
			ContractId: tc.contractID,
			From:       tc.operator.String(),
			To:         tc.to.String(),
			Params:     tc.params,
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

func TestMsgBurn(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	amount := collection.NewCoins(
		collection.NewNFTCoin("deadbeef", 1),
	)

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			from:       addrs[0],
			amount:     amount,
			valid:      true,
		},
		"invalid contract id": {
			from:   addrs[0],
			amount: amount,
		},
		"invalid from": {
			contractID: contractID,
			amount:     amount,
		},
		"empty amount": {
			contractID: contractID,
			from:       addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgBurn{
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

	contractID := "deadbeef"
	amount := collection.NewCoins(
		collection.NewNFTCoin("deadbeef", 1),
	)

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			amount:     amount,
			valid:      true,
		},
		"invalid contract id": {
			operator: addrs[0],
			from:     addrs[1],
			amount:   amount,
		},
		"invalid operator": {
			contractID: contractID,
			from:       addrs[1],
			amount:     amount,
		},
		"empty from": {
			contractID: contractID,
			operator:   addrs[0],
			amount:     amount,
		},
		"empty amount": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgOperatorBurn{
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

func TestMsgBurnFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(
		collection.NewFTCoin("00bab10c", sdk.OneInt()),
	)

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:       addrs[0],
			amount:     amount,
			valid:      true,
		},
		"invalid contract id": {
			from:   addrs[0],
			amount: amount,
		},
		"invalid from": {
			contractID: "deadbeef",
			amount:     amount,
		},
		"invalid token id": {
			contractID: "deadbeef",
			from:       addrs[0],
			amount: []collection.Coin{{
				Amount: sdk.OneInt(),
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgBurnFT{
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

func TestMsgBurnFTFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(
		collection.NewFTCoin("00bab10c", sdk.OneInt()),
	)

	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		from       sdk.AccAddress
		amount     []collection.Coin
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			from:       addrs[1],
			amount:     amount,
			valid:      true,
		},
		"invalid contract id": {
			grantee: addrs[0],
			from:    addrs[1],
			amount:  amount,
		},
		"invalid grantee": {
			contractID: "deadbeef",
			from:       addrs[1],
			amount:     amount,
		},
		"empty from": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			amount:     amount,
		},
		"invalid token id": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			from:       addrs[1],
			amount: []collection.Coin{{
				Amount: sdk.OneInt(),
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgBurnFTFrom{
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

func TestMsgBurnNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ids := []string{collection.NewNFTID("deadbeef", 1)}

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		ids        []string
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:       addrs[0],
			ids:        ids,
			valid:      true,
		},
		"invalid contract id": {
			from: addrs[0],
			ids:  ids,
		},
		"invalid from": {
			contractID: "deadbeef",
			ids:        ids,
		},
		"empty ids": {
			contractID: "deadbeef",
			from:       addrs[0],
		},
		"invalid id": {
			contractID: "deadbeef",
			from:       addrs[0],
			ids:        []string{""},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgBurnNFT{
			ContractId: tc.contractID,
			From:       tc.from.String(),
			TokenIds:   tc.ids,
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

func TestMsgBurnNFTFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ids := []string{collection.NewNFTID("deadbeef", 1)}

	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		from       sdk.AccAddress
		ids        []string
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			from:       addrs[1],
			ids:        ids,
			valid:      true,
		},
		"invalid contract id": {
			grantee: addrs[0],
			from:    addrs[1],
			ids:     ids,
		},
		"invalid grantee": {
			contractID: "deadbeef",
			from:       addrs[1],
			ids:        ids,
		},
		"empty from": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			ids:        ids,
		},
		"empty ids": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			from:       addrs[1],
		},
		"invalid id": {
			contractID: "deadbeef",
			grantee:    addrs[0],
			from:       addrs[0],
			ids:        []string{""},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgBurnNFTFrom{
			ContractId: tc.contractID,
			Proxy:      tc.grantee.String(),
			From:       tc.from.String(),
			TokenIds:   tc.ids,
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

func TestMsgModifyContract(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	changes := []collection.Attribute{{
		Key:   collection.AttributeKeyName.String(),
		Value: "fox",
	}}
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		changes    []collection.Attribute
		valid      bool
	}{
		"valid contract modification": {
			contractID: contractID,
			operator:   addrs[0],
			changes:    changes,
			valid:      true,
		},
		"invalid contract id": {
			operator: addrs[0],
			changes:  changes,
		},
		"invalid operator": {
			contractID: contractID,
			changes:    changes,
		},
		"invalid key of change": {
			contractID: contractID,
			operator:   addrs[0],
			changes:    []collection.Attribute{{Value: "fox"}},
		},
		"invalid value of change": {
			contractID: contractID,
			operator:   addrs[0],
			changes:    []collection.Attribute{{Key: "symbol"}},
		},
		"empty changes": {
			contractID: contractID,
			operator:   addrs[0],
		},
		"too many changes": {
			contractID: contractID,
			operator:   addrs[0],
			changes:    make([]collection.Attribute, 101),
		},
		"duplicated changes": {
			contractID: contractID,
			operator:   addrs[0],
			changes:    []collection.Attribute{changes[0], changes[0]},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgModifyContract{
			ContractId: tc.contractID,
			Operator:   tc.operator.String(),
			Changes:    tc.changes,
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

func TestMsgModifyTokenClass(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	classID := "deadbeef"
	changes := []collection.Attribute{{
		Key:   collection.AttributeKeyName.String(),
		Value: "tibetian fox",
	}}
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		classID    string
		changes    []collection.Attribute
		valid      bool
	}{
		"valid modification": {
			contractID: contractID,
			operator:   addrs[0],
			classID:    classID,
			changes:    changes,
			valid:      true,
		},
		"invalid contract id": {
			operator: addrs[0],
			classID:  classID,
			changes:  changes,
		},
		"invalid operator": {
			contractID: contractID,
			classID:    classID,
			changes:    changes,
		},
		"invalid class id": {
			contractID: contractID,
			operator:   addrs[0],
			changes:    changes,
		},
		"invalid key of change": {
			contractID: contractID,
			operator:   addrs[0],
			classID:    classID,
			changes:    []collection.Attribute{{Value: "tibetian fox"}},
		},
		"invalid value of change": {
			contractID: contractID,
			operator:   addrs[0],
			classID:    classID,
			changes:    []collection.Attribute{{Key: "symbol"}},
		},
		"empty changes": {
			contractID: contractID,
			operator:   addrs[0],
			classID:    classID,
		},
		"too many changes": {
			contractID: contractID,
			operator:   addrs[0],
			classID:    classID,
			changes:    make([]collection.Attribute, 101),
		},
		"duplicated changes": {
			contractID: contractID,
			operator:   addrs[0],
			classID:    classID,
			changes:    []collection.Attribute{changes[0], changes[0]},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgModifyTokenClass{
			ContractId: tc.contractID,
			Operator:   tc.operator.String(),
			ClassId:    tc.classID,
			Changes:    tc.changes,
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

func TestMsgModifyNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	classID := "deadbeef"
	tokenID := collection.NewNFTID(classID, 1)
	changes := []collection.Attribute{{
		Key:   collection.AttributeKeyName.String(),
		Value: "tibetian fox",
	}}
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		tokenID    string
		changes    []collection.Attribute
		valid      bool
	}{
		"valid modification": {
			contractID: contractID,
			operator:   addrs[0],
			tokenID:    tokenID,
			changes:    changes,
			valid:      true,
		},
		"invalid contract id": {
			operator: addrs[0],
			tokenID:  tokenID,
			changes:  changes,
		},
		"invalid operator": {
			contractID: contractID,
			tokenID:    tokenID,
			changes:    changes,
		},
		"invalid token id": {
			contractID: contractID,
			operator:   addrs[0],
			changes:    changes,
		},
		"invalid key of change": {
			contractID: contractID,
			operator:   addrs[0],
			tokenID:    tokenID,
			changes:    []collection.Attribute{{Value: "tibetian fox"}},
		},
		"invalid value of change": {
			contractID: contractID,
			operator:   addrs[0],
			tokenID:    tokenID,
			changes:    []collection.Attribute{{Key: "symbol"}},
		},
		"empty changes": {
			contractID: contractID,
			operator:   addrs[0],
			tokenID:    tokenID,
		},
		"too many changes": {
			contractID: contractID,
			operator:   addrs[0],
			tokenID:    tokenID,
			changes:    make([]collection.Attribute, 101),
		},
		"duplicated changes": {
			contractID: contractID,
			operator:   addrs[0],
			tokenID:    tokenID,
			changes:    []collection.Attribute{changes[0], changes[0]},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgModifyNFT{
			ContractId: tc.contractID,
			Operator:   tc.operator.String(),
			TokenId:    tc.tokenID,
			Changes:    tc.changes,
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

func TestMsgModify(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	changes := []collection.Change{{Field: "name", Value: "New test"}}
	testCases := map[string]struct {
		contractID string
		owner      sdk.AccAddress
		tokenType  string
		tokenIndex string
		changes    []collection.Change
		valid      bool
	}{
		"valid contract modification": {
			contractID: "deadbeef",
			owner:      addrs[0],
			changes:    changes,
			valid:      true,
		},
		"valid token class modification": {
			contractID: "deadbeef",
			tokenType:  "deadbeef",
			owner:      addrs[0],
			changes:    changes,
			valid:      true,
		},
		"valid nft modification": {
			contractID: "deadbeef",
			tokenType:  "deadbeef",
			tokenIndex: "deadbeef",
			owner:      addrs[0],
			changes:    changes,
			valid:      true,
		},
		"invalid contract id": {
			owner:   addrs[0],
			changes: changes,
		},
		"invalid owner": {
			contractID: "deadbeef",
			changes:    changes,
		},
		"invalid key of change": {
			contractID: "deadbeef",
			owner:      addrs[0],
			changes:    []collection.Change{{Value: "tt"}},
		},
		"invalid value of change": {
			contractID: "deadbeef",
			owner:      addrs[0],
			changes:    []collection.Change{{Field: "symbol"}},
		},
		"empty changes": {
			contractID: "deadbeef",
			owner:      addrs[0],
		},
		"too many changes": {
			contractID: "deadbeef",
			owner:      addrs[0],
			changes:    make([]collection.Change, 101),
		},
		"duplicated changes": {
			contractID: "deadbeef",
			owner:      addrs[0],
			changes:    []collection.Change{changes[0], changes[0]},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgModify{
			ContractId: tc.contractID,
			TokenType:  tc.tokenType,
			TokenIndex: tc.tokenIndex,
			Owner:      tc.owner.String(),
			Changes:    tc.changes,
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

func TestMsgGrant(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	testCases := map[string]struct {
		contractID string
		granter    sdk.AccAddress
		grantee    sdk.AccAddress
		permission collection.Permission
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			granter:    addrs[0],
			grantee:    addrs[1],
			permission: collection.PermissionMint,
			valid:      true,
		},
		"invalid contract id": {
			granter:    addrs[0],
			grantee:    addrs[1],
			permission: collection.PermissionMint,
		},
		"empty granter": {
			contractID: contractID,
			grantee:    addrs[1],
			permission: collection.PermissionMint,
		},
		"invalid grantee": {
			contractID: contractID,
			granter:    addrs[0],
			permission: collection.PermissionMint,
		},
		"invalid permission": {
			contractID: contractID,
			granter:    addrs[0],
			grantee:    addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgGrant{
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

	contractID := "deadbeef"
	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		permission collection.Permission
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			grantee:    addrs[0],
			permission: collection.PermissionMint,
			valid:      true,
		},
		"invalid contract id": {
			grantee:    addrs[0],
			permission: collection.PermissionMint,
		},
		"invalid grantee": {
			contractID: contractID,
			permission: collection.PermissionMint,
		},
		"invalid permission": {
			contractID: contractID,
			grantee:    addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgAbandon{
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
			permission: collection.LegacyPermissionMint.String(),
			valid:      true,
		},
		"invalid contract id": {
			from:       addrs[0],
			to:         addrs[1],
			permission: collection.LegacyPermissionMint.String(),
		},
		"empty from": {
			contractID: "deadbeef",
			to:         addrs[1],
			permission: collection.LegacyPermissionMint.String(),
		},
		"invalid to": {
			contractID: "deadbeef",
			from:       addrs[0],
			permission: collection.LegacyPermissionMint.String(),
		},
		"invalid permission": {
			contractID: "deadbeef",
			from:       addrs[0],
			to:         addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgGrantPermission{
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
			permission: collection.LegacyPermissionMint.String(),
			valid:      true,
		},
		"invalid contract id": {
			from:       addrs[0],
			permission: collection.LegacyPermissionMint.String(),
		},
		"invalid from": {
			contractID: "deadbeef",
			permission: collection.LegacyPermissionMint.String(),
		},
		"invalid permission": {
			contractID: "deadbeef",
			from:       addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgRevokePermission{
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

func TestMsgAttach(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	tokenIDs := []string{
		collection.NewNFTID("deadbeef", 1),
		collection.NewNFTID("fee1dead", 1),
	}

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		tokenID    string
		toTokenID  string
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			from:       addrs[0],
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[1],
			valid:      true,
		},
		"empty from": {
			contractID: contractID,
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[1],
		},
		"invalid contract id": {
			from:      addrs[0],
			tokenID:   tokenIDs[0],
			toTokenID: tokenIDs[1],
		},
		"invalid token id": {
			contractID: contractID,
			from:       addrs[0],
			toTokenID:  tokenIDs[1],
		},
		"invalid to id": {
			contractID: contractID,
			from:       addrs[0],
			tokenID:    tokenIDs[0],
		},
		"to itself": {
			contractID: contractID,
			from:       addrs[0],
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgAttach{
			ContractId: tc.contractID,
			From:       tc.from.String(),
			TokenId:    tc.tokenID,
			ToTokenId:  tc.toTokenID,
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

func TestMsgDetach(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	tokenID := collection.NewNFTID("deadbeef", 1)

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		tokenID    string
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			from:       addrs[0],
			tokenID:    tokenID,
			valid:      true,
		},
		"empty from": {
			contractID: contractID,
			tokenID:    tokenID,
		},
		"invalid contract id": {
			from:    addrs[0],
			tokenID: tokenID,
		},
		"invalid token id": {
			contractID: contractID,
			from:       addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgDetach{
			ContractId: tc.contractID,
			From:       tc.from.String(),
			TokenId:    tc.tokenID,
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

func TestMsgOperatorAttach(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	tokenIDs := []string{
		collection.NewNFTID("deadbeef", 1),
		collection.NewNFTID("fee1dead", 1),
	}

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		owner      sdk.AccAddress
		subject    string
		target     string
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			owner:      addrs[1],
			subject:    tokenIDs[0],
			target:     tokenIDs[1],
			valid:      true,
		},
		"invalid contract id": {
			operator: addrs[0],
			owner:    addrs[1],
			subject:  tokenIDs[0],
			target:   tokenIDs[1],
		},
		"empty operator": {
			contractID: contractID,
			owner:      addrs[1],
			subject:    tokenIDs[0],
			target:     tokenIDs[1],
		},
		"empty owner": {
			contractID: contractID,
			operator:   addrs[0],
			subject:    tokenIDs[0],
			target:     tokenIDs[1],
		},
		"invalid token id": {
			contractID: contractID,
			operator:   addrs[0],
			owner:      addrs[1],
			target:     tokenIDs[1],
		},
		"invalid to id": {
			contractID: contractID,
			operator:   addrs[0],
			owner:      addrs[1],
			subject:    tokenIDs[0],
		},
		"to itself": {
			contractID: contractID,
			operator:   addrs[0],
			owner:      addrs[1],
			subject:    tokenIDs[0],
			target:     tokenIDs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgOperatorAttach{
			ContractId: tc.contractID,
			Operator:   tc.operator.String(),
			Owner:      tc.owner.String(),
			Subject:    tc.subject,
			Target:     tc.target,
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

func TestMsgOperatorDetach(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := "deadbeef"
	tokenID := collection.NewNFTID("deadbeef", 1)

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		owner      sdk.AccAddress
		subject    string
		valid      bool
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			owner:      addrs[1],
			subject:    tokenID,
			valid:      true,
		},
		"invalid contract id": {
			operator: addrs[0],
			owner:    addrs[1],
			subject:  tokenID,
		},
		"empty operator": {
			contractID: contractID,
			owner:      addrs[1],
			subject:    tokenID,
		},
		"empty owner": {
			contractID: contractID,
			operator:   addrs[0],
			subject:    tokenID,
		},
		"invalid token id": {
			contractID: contractID,
			operator:   addrs[0],
			owner:      addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgOperatorDetach{
			ContractId: tc.contractID,
			Operator:   tc.operator.String(),
			Owner:      tc.owner.String(),
			Subject:    tc.subject,
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

func TestMsgAttachFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	tokenIDs := []string{
		collection.NewNFTID("deadbeef", 1),
		collection.NewNFTID("fee1dead", 1),
	}

	testCases := map[string]struct {
		contractID string
		proxy      sdk.AccAddress
		from       sdk.AccAddress
		tokenID    string
		toTokenID  string
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[1],
			valid:      true,
		},
		"empty proxy": {
			contractID: "deadbeef",
			from:       addrs[1],
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[1],
		},
		"empty from": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[1],
		},
		"invalid contract id": {
			proxy:     addrs[0],
			from:      addrs[1],
			tokenID:   tokenIDs[0],
			toTokenID: tokenIDs[1],
		},
		"invalid token id": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			toTokenID:  tokenIDs[1],
		},
		"invalid to id": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			tokenID:    tokenIDs[0],
		},
		"to itself": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgAttachFrom{
			ContractId: tc.contractID,
			Proxy:      tc.proxy.String(),
			From:       tc.from.String(),
			TokenId:    tc.tokenID,
			ToTokenId:  tc.toTokenID,
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

func TestMsgDetachFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	tokenID := collection.NewNFTID("deadbeef", 1)

	testCases := map[string]struct {
		contractID string
		proxy      sdk.AccAddress
		from       sdk.AccAddress
		tokenID    string
		valid      bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
			tokenID:    tokenID,
			valid:      true,
		},
		"empty proxy": {
			contractID: "deadbeef",
			from:       addrs[1],
			tokenID:    tokenID,
		},
		"empty from": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			tokenID:    tokenID,
		},
		"invalid contract id": {
			proxy:   addrs[0],
			from:    addrs[1],
			tokenID: tokenID,
		},
		"invalid token id": {
			contractID: "deadbeef",
			proxy:      addrs[0],
			from:       addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgDetachFrom{
			ContractId: tc.contractID,
			Proxy:      tc.proxy.String(),
			From:       tc.from.String(),
			TokenId:    tc.tokenID,
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
