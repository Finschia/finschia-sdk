package collection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/auth/legacy/legacytx"
	"github.com/line/lbm-sdk/x/collection"
)

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

func TestMsgModify(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	changes := []collection.Change{{Field: collection.AttributeKeyName.String(), Value: "New test"}}
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
			changes:    []collection.Change{{Field: strings.ToUpper(collection.AttributeKeyName.String()), Value: "tt"}},
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

func TestAminoJSON(t *testing.T) {
	tx := legacytx.StdTx{}
	var contractId = "deadbeef"
	var ftClassId = "00bab10c"

	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ftAmount := collection.NewCoins(collection.NewFTCoin(ftClassId, sdk.NewInt(1000000)))
	tokenIds := []string{collection.NewNFTID(contractId, 1)}
	nftParams := []collection.MintNFTParam{{
		TokenType: "deadbeef",
		Name:      "tibetian fox",
		Meta:      "Tibetian Fox",
	}}

	testCase := map[string]struct {
		msg          legacytx.LegacyMsg
		expectedType string
		expected     string
	}{
		"MsgTransferFT": {
			&collection.MsgTransferFT{
				ContractId: contractId,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				Amount:     ftAmount,
			},
			"/lbm.collection.v1.MsgTransferFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgTransferFT\",\"value\":{\"amount\":[{\"amount\":\"1000000\",\"token_id\":\"00bab10c00000000\"}],\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgTransferFTFrom": {
			&collection.MsgTransferFTFrom{
				ContractId: contractId,
				Proxy:      addrs[0].String(),
				From:       addrs[1].String(),
				To:         addrs[2].String(),
				Amount:     ftAmount,
			},
			"/lbm.collection.v1.MsgTransferFTFrom",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgTransferFTFrom\",\"value\":{\"amount\":[{\"amount\":\"1000000\",\"token_id\":\"00bab10c00000000\"}],\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"proxy\":\"%s\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String(), addrs[2].String()),
		},
		"MsgTransferNFT": {
			&collection.MsgTransferNFT{
				ContractId: contractId,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				TokenIds:   tokenIds,
			},
			"/lbm.collection.v1.MsgTransferNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgTransferNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"to\":\"%s\",\"token_ids\":[\"deadbeef00000001\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgTransferNFTFrom": {
			&collection.MsgTransferNFTFrom{
				ContractId: contractId,
				Proxy:      addrs[0].String(),
				From:       addrs[1].String(),
				To:         addrs[2].String(),
				TokenIds:   tokenIds,
			},
			"/lbm.collection.v1.MsgTransferNFTFrom",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgTransferNFTFrom\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"proxy\":\"%s\",\"to\":\"%s\",\"token_ids\":[\"deadbeef00000001\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String(), addrs[2].String()),
		},
		"MsgApprove": {
			&collection.MsgApprove{
				ContractId: contractId,
				Approver:   addrs[0].String(),
				Proxy:      addrs[1].String(),
			},
			"/lbm.collection.v1.MsgApprove",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgApprove\",\"value\":{\"approver\":\"%s\",\"contract_id\":\"deadbeef\",\"proxy\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgDisapprove": {
			&collection.MsgDisapprove{
				ContractId: contractId,
				Approver:   addrs[0].String(),
				Proxy:      addrs[1].String(),
			},
			"/lbm.collection.v1.MsgDisapprove",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgDisapprove\",\"value\":{\"approver\":\"%s\",\"contract_id\":\"deadbeef\",\"proxy\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgCreateContract": {
			&collection.MsgCreateContract{
				Owner:      addrs[0].String(),
				Name:       "Test Contract",
				BaseImgUri: "http://image.url",
				Meta:       "This is test",
			},
			"/lbm.collection.v1.MsgCreateContract",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgCreateContract\",\"value\":{\"base_img_uri\":\"http://image.url\",\"meta\":\"This is test\",\"name\":\"Test Contract\",\"owner\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgIssueFT": {
			&collection.MsgIssueFT{
				ContractId: contractId,
				Name:       "Test FT",
				Meta:       "This is FT Meta",
				Decimals:   6,
				Mintable:   false,
				Owner:      addrs[0].String(),
				To:         addrs[1].String(),
				Amount:     sdk.NewInt(1000000),
			},
			"/lbm.collection.v1.MsgIssueFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgIssueFT\",\"value\":{\"amount\":\"1000000\",\"contract_id\":\"deadbeef\",\"decimals\":6,\"meta\":\"This is FT Meta\",\"name\":\"Test FT\",\"owner\":\"%s\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgIssueNFT": {
			&collection.MsgIssueNFT{
				ContractId: contractId,
				Name:       "Test NFT",
				Meta:       "This is NFT Meta",
				Owner:      addrs[0].String(),
			},
			"/lbm.collection.v1.MsgIssueNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgIssueNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"meta\":\"This is NFT Meta\",\"name\":\"Test NFT\",\"owner\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgMintFT": {
			&collection.MsgMintFT{
				ContractId: contractId,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				Amount:     ftAmount,
			},
			"/lbm.collection.v1.MsgMintFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgMintFT\",\"value\":{\"amount\":[{\"amount\":\"1000000\",\"token_id\":\"00bab10c00000000\"}],\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgMintNFT": {
			&collection.MsgMintNFT{
				ContractId: contractId,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				Params:     nftParams,
			},
			"/lbm.collection.v1.MsgMintNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgMintNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"params\":[{\"meta\":\"Tibetian Fox\",\"name\":\"tibetian fox\",\"token_type\":\"deadbeef\"}],\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgBurnFT": {
			&collection.MsgBurnFT{
				ContractId: contractId,
				From:       addrs[0].String(),
				Amount:     ftAmount,
			},
			"/lbm.collection.v1.MsgBurnFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgBurnFT\",\"value\":{\"amount\":[{\"amount\":\"1000000\",\"token_id\":\"00bab10c00000000\"}],\"contract_id\":\"deadbeef\",\"from\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgBurnFTFrom": {
			&collection.MsgBurnFTFrom{
				ContractId: contractId,
				Proxy:      addrs[0].String(),
				From:       addrs[1].String(),
				Amount:     ftAmount,
			},
			"/lbm.collection.v1.MsgBurnFTFrom",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgBurnFTFrom\",\"value\":{\"amount\":[{\"amount\":\"1000000\",\"token_id\":\"00bab10c00000000\"}],\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"proxy\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String()),
		},
		"MsgBurnNFT": {
			&collection.MsgBurnNFT{
				ContractId: contractId,
				From:       addrs[0].String(),
				TokenIds:   tokenIds,
			},
			"/lbm.collection.v1.MsgBurnNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgBurnNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"token_ids\":[\"deadbeef00000001\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgBurnNFTFrom": {
			&collection.MsgBurnNFTFrom{
				ContractId: contractId,
				Proxy:      addrs[0].String(),
				From:       addrs[1].String(),
				TokenIds:   tokenIds,
			},
			"/lbm.collection.v1.MsgBurnNFTFrom",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgBurnNFTFrom\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"proxy\":\"%s\",\"token_ids\":[\"deadbeef00000001\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String()),
		},
		"MsgModify": {
			&collection.MsgModify{
				ContractId: contractId,
				Owner:      addrs[0].String(),
				TokenType:  "NewType",
				TokenIndex: "deadbeef",
				Changes:    []collection.Change{{Field: "name", Value: "New test"}},
			},
			"/lbm.collection.v1.MsgModify",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgModify\",\"value\":{\"changes\":[{\"field\":\"name\",\"value\":\"New test\"}],\"contract_id\":\"deadbeef\",\"owner\":\"%s\",\"token_index\":\"deadbeef\",\"token_type\":\"NewType\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgGrantPermission": {
			&collection.MsgGrantPermission{
				ContractId: contractId,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				Permission: collection.LegacyPermissionMint.String(),
			},
			"/lbm.collection.v1.MsgGrantPermission",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgGrantPermission\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"permission\":\"mint\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgRevokePermission": {
			&collection.MsgRevokePermission{
				ContractId: contractId,
				From:       addrs[0].String(),
				Permission: collection.LegacyPermissionMint.String(),
			},
			"/lbm.collection.v1.MsgRevokePermission",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgRevokePermission\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"permission\":\"mint\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgAttach": {
			&collection.MsgAttach{
				ContractId: contractId,
				From:       addrs[0].String(),
				TokenId:    collection.NewNFTID("deadbeef", 1),
				ToTokenId:  collection.NewNFTID("fee1dead", 1),
			},
			"/lbm.collection.v1.MsgAttach",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgAttach\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"to_token_id\":\"fee1dead00000001\",\"token_id\":\"deadbeef00000001\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgDetach": {
			&collection.MsgDetach{
				ContractId: contractId,
				From:       addrs[0].String(),
				TokenId:    collection.NewNFTID("fee1dead", 1),
			},
			"/lbm.collection.v1.MsgDetach",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgDetach\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"token_id\":\"fee1dead00000001\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgAttachFrom": {
			&collection.MsgAttachFrom{
				ContractId: contractId,
				Proxy:      addrs[0].String(),
				From:       addrs[1].String(),
				TokenId:    collection.NewNFTID("deadbeef", 1),
				ToTokenId:  collection.NewNFTID("fee1dead", 1),
			},
			"/lbm.collection.v1.MsgAttachFrom",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgAttachFrom\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"proxy\":\"%s\",\"to_token_id\":\"fee1dead00000001\",\"token_id\":\"deadbeef00000001\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String()),
		},
		"MsgDetachFrom": {
			&collection.MsgDetachFrom{
				ContractId: contractId,
				Proxy:      addrs[0].String(),
				From:       addrs[1].String(),
				TokenId:    collection.NewNFTID("fee1dead", 1),
			},
			"/lbm.collection.v1.MsgDetachFrom",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgDetachFrom\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"proxy\":\"%s\",\"token_id\":\"fee1dead00000001\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String()),
		},
	}

	for name, tc := range testCase {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tx.Msgs = []sdk.Msg{tc.msg}
			require.Equal(t, collection.RouterKey, tc.msg.Route())
			require.Equal(t, tc.expectedType, tc.msg.Type())
			require.Equal(t, tc.expected, string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{tc.msg}, "memo")))
		})
	}
}
