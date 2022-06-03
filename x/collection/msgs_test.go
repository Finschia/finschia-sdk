package collection_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

func TestMsgSend(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}
	amount := collection.NewCoins(collection.NewCoin(
		fmt.Sprintf("%s%08x", "deadbeef", 0),
		sdk.OneInt(),
	))

	testCases := map[string]struct {
		contractID string
		from    sdk.AccAddress
		to      sdk.AccAddress
		amount   []collection.Coin
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
			amount: amount,
			valid:   true,
		},
		"empty from": {
			contractID: "deadbeef",
			to:      addrs[1],
			amount: amount,
		},
		"invalid contract id": {
			from:    addrs[0],
			to:      addrs[1],
			amount: amount,
		},
		"invalid to": {
			contractID: "deadbeef",
			from:    addrs[0],
			amount: amount,
		},
		"empty amount": {
			contractID: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
		},
		"invalid token id": {
			contractID: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
			amount:   []collection.Coin{{
				Amount: sdk.OneInt(),
			}},
		},
		"duplicate token ids": {
			contractID: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
			amount:   []collection.Coin{amount[0], amount[0]},
		},
		"invalid amount": {
			contractID: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
			amount:   []collection.Coin{{
				TokenId: amount[0].TokenId,
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgSend{
			ContractId: tc.contractID,
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

	amount := collection.NewCoins(collection.NewCoin(
		fmt.Sprintf("%s%08x", "deadbeef", 0),
		sdk.OneInt(),
	))

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from    sdk.AccAddress
		to      sdk.AccAddress
		amount  []collection.Coin
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			operator:   addrs[0],
			from:    addrs[1],
			to:      addrs[2],
			amount:  amount,
			valid:   true,
		},
		"invalid operator": {
			contractID: "deadbeef",
			from:    addrs[1],
			to:      addrs[2],
			amount:  amount,
		},
		"invalid contract id": {
			operator:   addrs[0],
			from:    addrs[1],
			to:      addrs[2],
			amount:  amount,
		},
		"empty from": {
			contractID: "deadbeef",
			operator:   addrs[0],
			to:      addrs[1],
			amount:  amount,
		},
		"invalid to": {
			contractID: "deadbeef",
			operator:   addrs[0],
			from:    addrs[1],
			amount:  amount,
		},
		"empty amount": {
			contractID: "deadbeef",
			operator:   addrs[0],
			from:    addrs[1],
			to:      addrs[2],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgOperatorSend{
			ContractId: tc.contractID,
			Operator:   tc.operator.String(),
			From:    tc.from.String(),
			To:      tc.to.String(),
			Amount:  tc.amount,
		}

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgTransferFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(collection.NewCoin(
		fmt.Sprintf("%s%08x", "deadbeef", 0),
		sdk.OneInt(),
	))

	testCases := map[string]struct {
		contractID string
		from    sdk.AccAddress
		to      sdk.AccAddress
		amount   []collection.Coin
		valid   bool
		panic   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
			amount: amount,
			valid:   true,
		},
		"empty from": {
			contractID: "deadbeef",
			to:      addrs[1],
			amount: amount,
		},
		"invalid contract id": {
			from:    addrs[0],
			to:      addrs[1],
			amount: amount,
		},
		"invalid to": {
			contractID: "deadbeef",
			from:    addrs[0],
			amount: amount,
		},
		"nil amount": {
			contractID: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
			amount: []collection.Coin{{
				TokenId: fmt.Sprintf("%s%08x", "deadbeef", 0),
			}},
			panic: true,
		},
		"zero amount": {
			contractID: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
			amount: []collection.Coin{{
				TokenId: fmt.Sprintf("%s%08x", "deadbeef", 0),
				Amount: sdk.ZeroInt(),
			}},
		},
		"invalid token id": {
			contractID: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
			amount: []collection.Coin{{
				Amount: sdk.OneInt(),
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgTransferFT{
			ContractId: tc.contractID,
			From:    tc.from.String(),
			To:      tc.to.String(),
			Amount:  tc.amount,
		}

		require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())

		if tc.panic {
			require.Panics(t, func(){msg.ValidateBasic()}, name)
			continue
		}
		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgTransferFTFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(collection.NewCoin(
		fmt.Sprintf("%s%08x", "deadbeef", 0),
		sdk.OneInt(),
	))

	testCases := map[string]struct {
		contractID string
		proxy   sdk.AccAddress
		from    sdk.AccAddress
		to      sdk.AccAddress
		amount  []collection.Coin
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			proxy:   addrs[0],
			from:    addrs[1],
			to:      addrs[2],
			amount:  amount,
			valid:   true,
		},
		"invalid proxy": {
			contractID: "deadbeef",
			from:    addrs[1],
			to:      addrs[2],
			amount:  amount,
		},
		"invalid contract id": {
			proxy:   addrs[0],
			from:    addrs[1],
			to:      addrs[2],
			amount:  amount,
		},
		"empty from": {
			contractID: "deadbeef",
			proxy:   addrs[0],
			to:      addrs[1],
			amount:  amount,
		},
		"invalid to": {
			contractID: "deadbeef",
			proxy:   addrs[0],
			from:    addrs[1],
			amount:  amount,
		},
		"invalid amount": {
			contractID: "deadbeef",
			proxy:   addrs[0],
			from:    addrs[1],
			to:      addrs[2],
			amount: []collection.Coin{{
				Amount: sdk.OneInt(),
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgTransferFTFrom{
			ContractId: tc.contractID,
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

func TestMsgTransferNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ids := []string{fmt.Sprintf("%s%08x", "deadbeef", 0)}

	testCases := map[string]struct {
		contractID string
		from    sdk.AccAddress
		to      sdk.AccAddress
		ids   []string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
			ids:     ids,
			valid:   true,
		},
		"empty from": {
			contractID: "deadbeef",
			to:      addrs[1],
			ids:     ids,
		},
		"invalid contract id": {
			from:    addrs[0],
			to:      addrs[1],
			ids:     ids,
		},
		"invalid to": {
			contractID: "deadbeef",
			from:    addrs[0],
			ids:     ids,
		},
		"empty token ids": {
			contractID: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
		},
		"invalid token ids": {
			contractID: "deadbeef",
			from:    addrs[0],
			to:      addrs[1],
			ids:     []string{""},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgTransferNFT{
			ContractId: tc.contractID,
			From:    tc.from.String(),
			To:      tc.to.String(),
			TokenIds: tc.ids,
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

func TestMsgTransferNFTFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ids := []string{fmt.Sprintf("%s%08x", "deadbeef", 0)}

	testCases := map[string]struct {
		contractID string
		proxy   sdk.AccAddress
		from    sdk.AccAddress
		to      sdk.AccAddress
		ids  []string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			proxy:   addrs[0],
			from:    addrs[1],
			to:      addrs[2],
			ids: ids,
			valid:   true,
		},
		"invalid proxy": {
			contractID: "deadbeef",
			from:    addrs[1],
			to:      addrs[2],
			ids:     ids,
		},
		"invalid contract id": {
			proxy:   addrs[0],
			from:    addrs[1],
			to:      addrs[2],
			ids:     ids,
		},
		"empty from": {
			contractID: "deadbeef",
			proxy:   addrs[0],
			to:      addrs[1],
			ids:     ids,
		},
		"invalid to": {
			contractID: "deadbeef",
			proxy:   addrs[0],
			from:    addrs[1],
			ids:     ids,
		},
		"empty ids": {
			contractID: "deadbeef",
			proxy:   addrs[0],
			from:    addrs[1],
			to:      addrs[2],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgTransferNFTFrom{
			ContractId: tc.contractID,
			Proxy:   tc.proxy.String(),
			From:    tc.from.String(),
			To:      tc.to.String(),
			TokenIds:  tc.ids,
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
		contractID  string
		holder sdk.AccAddress
		operator    sdk.AccAddress
		valid    bool
	}{
		"valid msg": {
			contractID:  "deadbeef",
			holder: addrs[0],
			operator:    addrs[1],
			valid:    true,
		},
		"invalid contract id": {
			holder: addrs[0],
			operator:    addrs[1],
		},
		"invalid holder": {
			contractID:  "deadbeef",
			operator:    addrs[1],
		},
		"empty operator": {
			contractID:  "deadbeef",
			holder: addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgAuthorizeOperator{
			ContractId:  tc.contractID,
			Holder: tc.holder.String(),
			Operator:    tc.operator.String(),
		}

		require.Equal(t, []sdk.AccAddress{tc.holder}, msg.GetSigners())

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
		contractID  string
		holder sdk.AccAddress
		operator    sdk.AccAddress
		valid    bool
	}{
		"valid msg": {
			contractID:  "deadbeef",
			holder: addrs[0],
			operator:    addrs[1],
			valid:    true,
		},
		"invalid contract id": {
			holder: addrs[0],
			operator:    addrs[1],
		},
		"invalid holder": {
			contractID:  "deadbeef",
			operator:    addrs[1],
		},
		"empty operator": {
			contractID:  "deadbeef",
			holder: addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgRevokeOperator{
			ContractId:  tc.contractID,
			Holder: tc.holder.String(),
			Operator:    tc.operator.String(),
		}

		require.Equal(t, []sdk.AccAddress{tc.holder}, msg.GetSigners())

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
		contractID  string
		approver sdk.AccAddress
		proxy    sdk.AccAddress
		valid    bool
	}{
		"valid msg": {
			contractID:  "deadbeef",
			approver: addrs[0],
			proxy:    addrs[1],
			valid:    true,
		},
		"invalid contract id": {
			approver: addrs[0],
			proxy:    addrs[1],
		},
		"invalid approver": {
			contractID:  "deadbeef",
			proxy:    addrs[1],
		},
		"empty proxy": {
			contractID:  "deadbeef",
			approver: addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgApprove{
			ContractId:  tc.contractID,
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

func TestMsgDisapprove(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID  string
		approver sdk.AccAddress
		proxy    sdk.AccAddress
		valid    bool
	}{
		"valid msg": {
			contractID:  "deadbeef",
			approver: addrs[0],
			proxy:    addrs[1],
			valid:    true,
		},
		"invalid contract id": {
			approver: addrs[0],
			proxy:    addrs[1],
		},
		"invalid approver": {
			contractID:  "deadbeef",
			proxy:    addrs[1],
		},
		"empty proxy": {
			contractID:  "deadbeef",
			approver: addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgDisapprove{
			ContractId:  tc.contractID,
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

func TestMsgCreateContract(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		owner    sdk.AccAddress
		name     string
		baseImgURI string
		meta     string
		valid    bool
	}{
		"valid msg": {
			owner:    addrs[0],
			name:     "test",
			baseImgURI: "some URI",
			meta:     "some meta",
			valid:    true,
		},
		"invalid owner": {
			name:     "test",
			baseImgURI: "some URI",
			meta:     "some meta",
		},
		"long name": {
			owner:    addrs[0],
			name:     string(make([]rune, 21)),
			baseImgURI: "some URI",
			meta:     "some meta",
		},
		"invalid base image uri": {
			owner:    addrs[0],
			name:     "test",
			baseImgURI: string(make([]rune, 1001)),
			meta:     "some meta",
		},
		"invalid meta": {
			owner:    addrs[0],
			name:     "test",
			baseImgURI: "some URI",
			meta:     string(make([]rune, 1001)),
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgCreateContract{
			Owner:    tc.owner.String(),
			Name:     tc.name,
			BaseImgUri: tc.baseImgURI,
			Meta:     tc.meta,
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

func TestMsgIssueFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		owner    sdk.AccAddress
		to       sdk.AccAddress
		name     string
		meta     string
		decimals int32
		mintable bool
		amount   sdk.Int
		valid    bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			owner:    addrs[0],
			to:       addrs[1],
			name:     "test",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
			valid:    true,
		},
		"invalid contract id": {
			owner:    addrs[0],
			to:       addrs[1],
			name:     "test",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
		},
		"invalid owner": {
			contractID: "deadbeef",
			to:       addrs[1],
			name:     "test",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
		},
		"empty to": {
			contractID: "deadbeef",
			owner:    addrs[0],
			name:     "test",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
		},
		"empty name": {
			contractID: "deadbeef",
			owner:    addrs[0],
			to:       addrs[1],
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
		},
		"long name": {
			contractID: "deadbeef",
			owner:    addrs[0],
			to:       addrs[1],
			name:     string(make([]rune, 21)),
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.OneInt(),
			valid:    false,
		},
		"invalid meta": {
			contractID: "deadbeef",
			owner:    addrs[0],
			to:       addrs[1],
			name:     "test",
			meta:     string(make([]rune, 1001)),
			decimals: 8,
			amount:   sdk.OneInt(),
		},
		"invalid decimals": {
			contractID: "deadbeef",
			owner:    addrs[0],
			to:       addrs[1],
			name:     "test",
			meta:     "some meta",
			decimals: 19,
			amount:   sdk.OneInt(),
		},
		"valid supply": {
			contractID: "deadbeef",
			owner:    addrs[0],
			to:       addrs[1],
			name:     "test",
			meta:     "some meta",
			decimals: 8,
			amount:   sdk.ZeroInt(),
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgIssueFT{
			ContractId: tc.contractID,
			Owner:    tc.owner.String(),
			To:       tc.to.String(),
			Name:     tc.name,
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

func TestMsgIssueNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		owner    sdk.AccAddress
		name     string
		meta     string
		valid    bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			owner:    addrs[0],
			name:     "test",
			meta:     "some meta",
			valid:    true,
		},
		"invalid contract id": {
			owner:    addrs[0],
			name:     "test",
			meta:     "some meta",
		},
		"invalid owner": {
			contractID: "deadbeef",
			name:     "test",
			meta:     "some meta",
		},
		"long name": {
			contractID: "deadbeef",
			owner:    addrs[0],
			name:     string(make([]rune, 21)),
			meta:     "some meta",
		},
		"invalid meta": {
			contractID: "deadbeef",
			owner:    addrs[0],
			name:     "test",
			meta:     string(make([]rune, 1001)),
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgIssueNFT{
			ContractId: tc.contractID,
			Owner:    tc.owner.String(),
			Name:     tc.name,
			Meta:     tc.meta,
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

func TestMsgMintFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(collection.NewCoin(
		fmt.Sprintf("%s%08x", "deadbeef", 0),
		sdk.OneInt(),
	))

	testCases := map[string]struct {
		contractID string
		grantee sdk.AccAddress
		to      sdk.AccAddress
		amount  []collection.Coin
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee: addrs[0],
			to:      addrs[1],
			amount:  amount,
			valid:   true,
		},
		"invalid contract id": {
			grantee: addrs[0],
			to:      addrs[1],
			amount:  amount,
		},
		"invalid grantee": {
			contractID: "deadbeef",
			to:      addrs[1],
			amount:  amount,
		},
		"empty to": {
			contractID: "deadbeef",
			grantee: addrs[0],
			amount:  amount,
		},
		"invalid token id": {
			contractID: "deadbeef",
			grantee: addrs[0],
			to:      addrs[1],
			amount: []collection.Coin{{
				Amount: sdk.OneInt(),
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgMintFT{
			ContractId: tc.contractID,
			From: tc.grantee.String(),
			To:      tc.to.String(),
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

func TestMsgMintNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	params := []collection.MintNFTParam{{
		TokenType: "deadbeef",
		Name: "tibetian fox",
		Meta: "Tibetian Fox",
	}}
	testCases := map[string]struct {
		contractID string
		grantee sdk.AccAddress
		to      sdk.AccAddress
		params  []collection.MintNFTParam
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee: addrs[0],
			to:      addrs[1],
			params:  params,
			valid:   true,
		},
		"invalid contract id": {
			grantee: addrs[0],
			to:      addrs[1],
			params: params,
		},
		"invalid grantee": {
			contractID: "deadbeef",
			to:      addrs[1],
			params:  params,
		},
		"empty to": {
			contractID: "deadbeef",
			grantee: addrs[0],
			params: params,
		},
		"empty params": {
			contractID: "deadbeef",
			grantee: addrs[0],
			to:      addrs[1],
		},
		"invalid param": {
			contractID: "deadbeef",
			grantee: addrs[0],
			to:      addrs[1],
			params: []collection.MintNFTParam{{
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgMintNFT{
			ContractId: tc.contractID,
			From: tc.grantee.String(),
			To:      tc.to.String(),
			Params: tc.params,
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

	amount := collection.NewCoins(collection.NewCoin(
		fmt.Sprintf("%s%08x", "deadbeef", 0),
		sdk.OneInt(),
	))

	testCases := map[string]struct {
		contractID string
		from    sdk.AccAddress
		amount  []collection.Coin
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:    addrs[0],
			amount:  amount,
			valid:   true,
		},
		"invalid contract id": {
			from:    addrs[0],
			amount:  amount,
		},
		"invalid from": {
			contractID: "deadbeef",
			amount:  amount,
		},
		"empty amount": {
			contractID: "deadbeef",
			from:    addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgBurn{
			ContractId: tc.contractID,
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

	amount := collection.NewCoins(collection.NewCoin(
		fmt.Sprintf("%s%08x", "deadbeef", 0),
		sdk.OneInt(),
	))

	testCases := map[string]struct {
		contractID string
		operator sdk.AccAddress
		from    sdk.AccAddress
		amount  []collection.Coin
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			operator: addrs[0],
			from:    addrs[1],
			amount:  amount,
			valid:   true,
		},
		"invalid contract id": {
			operator: addrs[0],
			from:    addrs[1],
			amount:  amount,
		},
		"invalid operator": {
			contractID: "deadbeef",
			from:    addrs[1],
			amount:  amount,
		},
		"empty from": {
			contractID: "deadbeef",
			operator: addrs[0],
			amount:  amount,
		},
		"empty amount": {
			contractID: "deadbeef",
			operator: addrs[0],
			from:    addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgOperatorBurn{
			ContractId: tc.contractID,
			Operator: tc.operator.String(),
			From:    tc.from.String(),
			Amount:  tc.amount,
		}

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgBurnFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(collection.NewCoin(
		fmt.Sprintf("%s%08x", "deadbeef", 0),
		sdk.OneInt(),
	))

	testCases := map[string]struct {
		contractID string
		from    sdk.AccAddress
		amount  []collection.Coin
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:    addrs[0],
			amount:  amount,
			valid:   true,
		},
		"invalid contract id": {
			from:    addrs[0],
			amount:  amount,
		},
		"invalid from": {
			contractID: "deadbeef",
			amount:  amount,
		},
		"invalid token id": {
			contractID: "deadbeef",
			from:    addrs[0],
			amount: []collection.Coin{{
				Amount: sdk.OneInt(),
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgBurnFT{
			ContractId: tc.contractID,
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

func TestMsgBurnFTFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(collection.NewCoin(
		fmt.Sprintf("%s%08x", "deadbeef", 0),
		sdk.OneInt(),
	))

	testCases := map[string]struct {
		contractID string
		grantee sdk.AccAddress
		from    sdk.AccAddress
		amount  []collection.Coin
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee: addrs[0],
			from:    addrs[1],
			amount:  amount,
			valid:   true,
		},
		"invalid contract id": {
			grantee: addrs[0],
			from:    addrs[1],
			amount:  amount,
		},
		"invalid grantee": {
			contractID: "deadbeef",
			from:    addrs[1],
			amount:  amount,
		},
		"empty from": {
			contractID: "deadbeef",
			grantee: addrs[0],
			amount:  amount,
		},
		"invalid token id": {
			contractID: "deadbeef",
			grantee: addrs[0],
			from:    addrs[1],
			amount: []collection.Coin{{
				Amount: sdk.OneInt(),
			}},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgBurnFTFrom{
			ContractId: tc.contractID,
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

func TestMsgBurnNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ids := []string{fmt.Sprintf("%s%08x", "deadbeef", 0)}

	testCases := map[string]struct {
		contractID string
		from    sdk.AccAddress
		ids []string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:    addrs[0],
			ids: ids,
			valid:   true,
		},
		"invalid contract id": {
			from:    addrs[0],
			ids: ids,
		},
		"invalid from": {
			contractID: "deadbeef",
			ids: ids,
		},
		"empty ids": {
			contractID: "deadbeef",
			from:    addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgBurnNFT{
			ContractId: tc.contractID,
			From:    tc.from.String(),
			TokenIds: tc.ids,
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

func TestMsgBurnNFTFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ids := []string{fmt.Sprintf("%s%08x", "deadbeef", 0)}

	testCases := map[string]struct {
		contractID string
		grantee sdk.AccAddress
		from    sdk.AccAddress
		ids  []string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee: addrs[0],
			from:    addrs[1],
			ids:  ids,
			valid:   true,
		},
		"invalid contract id": {
			grantee: addrs[0],
			from:    addrs[1],
			ids:  ids,
		},
		"invalid grantee": {
			contractID: "deadbeef",
			from:    addrs[1],
			ids:  ids,
		},
		"empty from": {
			contractID: "deadbeef",
			grantee: addrs[0],
			ids:  ids,
		},
		"empty ids": {
			contractID: "deadbeef",
			grantee: addrs[0],
			from:    addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgBurnNFTFrom{
			ContractId: tc.contractID,
			Proxy: tc.grantee.String(),
			From:    tc.from.String(),
			TokenIds:  tc.ids,
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

	validChange := collection.Pair{Field: "name", Value: "New test"}
	testCases := map[string]struct {
		contractID string
		grantee sdk.AccAddress
		tokenType string
		tokenIndex string
		changes []collection.Pair
		valid   bool
	}{
		"valid contract modification": {
			contractID: "deadbeef",
			grantee: addrs[0],
			changes: []collection.Pair{validChange},
			valid:   true,
		},
		"valid token class modification": {
			contractID: "deadbeef",
			tokenType: "deadbeef",
			grantee: addrs[0],
			changes: []collection.Pair{validChange},
			valid:   true,
		},
		"valid nft modification": {
			contractID: "deadbeef",
			tokenType: "deadbeef",
			tokenIndex: "deadbeef",
			grantee: addrs[0],
			changes: []collection.Pair{validChange},
			valid:   true,
		},
		"invalid contract id": {
			grantee: addrs[0],
			changes: []collection.Pair{validChange},
		},
		"invalid grantee": {
			contractID: "deadbeef",
			changes: []collection.Pair{validChange},
		},
		"invalid key of change": {
			contractID: "deadbeef",
			grantee: addrs[0],
			changes: []collection.Pair{{Value: "tt"}},
		},
		"invalid value of change": {
			contractID: "deadbeef",
			grantee: addrs[0],
			changes: []collection.Pair{{Field: "symbol"}},
		},
		"empty changes": {
			contractID: "deadbeef",
			grantee: addrs[0],
		},
		"too many changes": {
			contractID: "deadbeef",
			grantee: addrs[0],
			changes: make([]collection.Pair, 101),
		},
		"duplicated changes": {
			contractID: "deadbeef",
			grantee: addrs[0],
			changes: []collection.Pair{
				{Field: "name", Value: "hello"},
				{Field: "name", Value: "world"},
			},
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgModify{
			ContractId: tc.contractID,
			TokenType: tc.tokenType,
			TokenIndex: tc.tokenIndex,
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
		contractID string
		granter sdk.AccAddress
		grantee sdk.AccAddress
		permission  string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			granter: addrs[0],
			grantee: addrs[1],
			permission:  collection.Permission_Mint.String(),
			valid:   true,
		},
		"invalid contract id": {
			granter: addrs[0],
			grantee: addrs[1],
			permission:  collection.Permission_Mint.String(),
		},
		"empty granter": {
			contractID: "deadbeef",
			grantee: addrs[1],
			permission:  collection.Permission_Mint.String(),
		},
		"invalid grantee": {
			contractID: "deadbeef",
			granter: addrs[0],
			permission:  collection.Permission_Mint.String(),
		},
		"invalid permission": {
			contractID: "deadbeef",
			granter: addrs[0],
			grantee: addrs[1],
			permission:  "invalid",
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgGrant{
			ContractId: tc.contractID,
			Granter: tc.granter.String(),
			Grantee: tc.grantee.String(),
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
		contractID string
		grantee sdk.AccAddress
		permission  string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			grantee: addrs[0],
			permission:  collection.Permission_Mint.String(),
			valid:   true,
		},
		"invalid contract id": {
			grantee: addrs[0],
			permission:  collection.Permission_Mint.String(),
		},
		"invalid grantee": {
			contractID: "deadbeef",
			permission:  collection.Permission_Mint.String(),
		},
		"invalid permission": {
			contractID: "deadbeef",
			grantee: addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgAbandon{
			ContractId: tc.contractID,
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

func TestMsgGrantPermission(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		from sdk.AccAddress
		to sdk.AccAddress
		permission  string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from: addrs[0],
			to: addrs[1],
			permission:  collection.Permission_Mint.String(),
			valid:   true,
		},
		"invalid contract id": {
			from: addrs[0],
			to: addrs[1],
			permission:  collection.Permission_Mint.String(),
		},
		"empty from": {
			contractID: "deadbeef",
			to: addrs[1],
			permission:  collection.Permission_Mint.String(),
		},
		"invalid to": {
			contractID: "deadbeef",
			from: addrs[0],
			permission:  collection.Permission_Mint.String(),
		},
		"invalid permission": {
			contractID: "deadbeef",
			from: addrs[0],
			to: addrs[1],
			permission:  "invalid",
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgGrantPermission{
			ContractId: tc.contractID,
			From: tc.from.String(),
			To: tc.to.String(),
			Permission:  tc.permission,
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

func TestMsgRevokePermission(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	testCases := map[string]struct {
		contractID string
		from sdk.AccAddress
		permission  string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from: addrs[0],
			permission:  collection.Permission_Mint.String(),
			valid:   true,
		},
		"invalid contract id": {
			from: addrs[0],
			permission:  collection.Permission_Mint.String(),
		},
		"invalid from": {
			contractID: "deadbeef",
			permission:  collection.Permission_Mint.String(),
		},
		"invalid permission": {
			contractID: "deadbeef",
			from: addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgRevokePermission{
			ContractId: tc.contractID,
			From: tc.from.String(),
			Permission:  tc.permission,
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

func TestMsgAttach(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	tokenIDs := []string{
		"deadbeef" + fmt.Sprintf("%08x", 0),
		"fee1dead" + fmt.Sprintf("%08x", 0),
	}

	testCases := map[string]struct {
		contractID string
		from    sdk.AccAddress
		tokenID string
		toTokenID string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:    addrs[0],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[1],
			valid:   true,
		},
		"empty from": {
			contractID: "deadbeef",
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[1],
		},
		"invalid contract id": {
			from:    addrs[0],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[1],
		},
		"invalid token id": {
			contractID: "deadbeef",
			from:    addrs[0],
			toTokenID: tokenIDs[1],
		},
		"invalid to id": {
			contractID: "deadbeef",
			from:    addrs[0],
			tokenID:      tokenIDs[0],
		},
		"to itself": {
			contractID: "deadbeef",
			from:    addrs[0],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgAttach{
			ContractId: tc.contractID,
			From:    tc.from.String(),
			TokenId:      tc.tokenID,
			ToTokenId:  tc.toTokenID,
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

func TestMsgDetach(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	tokenIDs := []string{
		"deadbeef" + fmt.Sprintf("%08x", 0),
	}

	testCases := map[string]struct {
		contractID string
		from    sdk.AccAddress
		tokenID string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			from:    addrs[0],
			tokenID:      tokenIDs[0],
			valid:   true,
		},
		"empty from": {
			contractID: "deadbeef",
			tokenID:      tokenIDs[0],
		},
		"invalid contract id": {
			from:    addrs[0],
			tokenID:      tokenIDs[0],
		},
		"invalid token id": {
			contractID: "deadbeef",
			from:    addrs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgDetach{
			ContractId: tc.contractID,
			From:    tc.from.String(),
			TokenId:      tc.tokenID,
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

func TestMsgOperatorAttach(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	tokenIDs := []string{
		"deadbeef" + fmt.Sprintf("%08x", 0),
		"fee1dead" + fmt.Sprintf("%08x", 0),
	}

	testCases := map[string]struct {
		contractID string
		operator sdk.AccAddress
		owner sdk.AccAddress
		tokenID string
		toTokenID string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			operator: addrs[0],
			owner: addrs[1],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[1],
			valid:   true,
		},
		"invalid contract id": {
			operator:    addrs[0],
			owner: addrs[1],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[1],
		},
		"empty operator": {
			contractID: "deadbeef",
			owner: addrs[1],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[1],
		},
		"empty owner": {
			contractID: "deadbeef",
			operator: addrs[0],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[1],
		},
		"invalid token id": {
			contractID: "deadbeef",
			operator:    addrs[0],
			owner: addrs[1],
			toTokenID: tokenIDs[1],
		},
		"invalid to id": {
			contractID: "deadbeef",
			operator:    addrs[0],
			owner: addrs[1],
			tokenID:      tokenIDs[0],
		},
		"to itself": {
			contractID: "deadbeef",
			operator:    addrs[0],
			owner: addrs[1],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgOperatorAttach{
			ContractId: tc.contractID,
			Operator:    tc.operator.String(),
			Owner: tc.owner.String(),
			Id:      tc.tokenID,
			To:  tc.toTokenID,
		}

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgOperatorDetach(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	tokenIDs := []string{
		"deadbeef" + fmt.Sprintf("%08x", 0),
	}

	testCases := map[string]struct {
		contractID string
		operator sdk.AccAddress
		owner sdk.AccAddress
		tokenID string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			operator: addrs[0],
			owner: addrs[1],
			tokenID:      tokenIDs[0],
			valid:   true,
		},
		"invalid contract id": {
			operator:    addrs[0],
			owner: addrs[1],
			tokenID:      tokenIDs[0],
		},
		"empty operator": {
			contractID: "deadbeef",
			owner: addrs[1],
			tokenID:      tokenIDs[0],
		},
		"empty owner": {
			contractID: "deadbeef",
			operator: addrs[0],
			tokenID:      tokenIDs[0],
		},
		"invalid token id": {
			contractID: "deadbeef",
			operator:    addrs[0],
			owner: addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgOperatorDetach{
			ContractId: tc.contractID,
			Operator:    tc.operator.String(),
			Owner: tc.owner.String(),
			Id:      tc.tokenID,
		}

		require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())

		err := msg.ValidateBasic()
		if tc.valid {
			require.NoError(t, err, name)
		} else {
			require.Error(t, err, name)
		}
	}
}

func TestMsgAttachFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	tokenIDs := []string{
		"deadbeef" + fmt.Sprintf("%08x", 0),
		"fee1dead" + fmt.Sprintf("%08x", 0),
	}

	testCases := map[string]struct {
		contractID string
		proxy sdk.AccAddress
		from    sdk.AccAddress
		tokenID string
		toTokenID string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			proxy:    addrs[0],
			from: addrs[1],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[1],
			valid:   true,
		},
		"empty proxy": {
			contractID: "deadbeef",
			from: addrs[1],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[1],
		},
		"empty from": {
			contractID: "deadbeef",
			proxy: addrs[0],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[1],
		},
		"invalid contract id": {
			proxy:    addrs[0],
			from: addrs[1],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[1],
		},
		"invalid token id": {
			contractID: "deadbeef",
			proxy:    addrs[0],
			from: addrs[1],
			toTokenID: tokenIDs[1],
		},
		"invalid to id": {
			contractID: "deadbeef",
			proxy:    addrs[0],
			from: addrs[1],
			tokenID:      tokenIDs[0],
		},
		"to itself": {
			contractID: "deadbeef",
			proxy:    addrs[0],
			from: addrs[1],
			tokenID:      tokenIDs[0],
			toTokenID: tokenIDs[0],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgAttachFrom{
			ContractId: tc.contractID,
			Proxy:    tc.proxy.String(),
			From: tc.from.String(),
			TokenId:      tc.tokenID,
			ToTokenId:  tc.toTokenID,
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

func TestMsgDetachFrom(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	tokenIDs := []string{
		"deadbeef" + fmt.Sprintf("%08x", 0),
	}

	testCases := map[string]struct {
		contractID string
		proxy    sdk.AccAddress
		from sdk.AccAddress
		tokenID string
		valid   bool
	}{
		"valid msg": {
			contractID: "deadbeef",
			proxy:    addrs[0],
			from: addrs[1],
			tokenID:      tokenIDs[0],
			valid:   true,
		},
		"empty proxy": {
			contractID: "deadbeef",
			from: addrs[1],
			tokenID:      tokenIDs[0],
		},
		"empty from": {
			contractID: "deadbeef",
			proxy: addrs[0],
			tokenID:      tokenIDs[0],
		},
		"invalid contract id": {
			proxy:    addrs[0],
			from: addrs[1],
			tokenID:      tokenIDs[0],
		},
		"invalid token id": {
			contractID: "deadbeef",
			proxy:    addrs[0],
			from: addrs[1],
		},
	}

	for name, tc := range testCases {
		msg := collection.MsgDetachFrom{
			ContractId: tc.contractID,
			Proxy:    tc.proxy.String(),
			From: tc.from.String(),
			TokenId:      tc.tokenID,
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
