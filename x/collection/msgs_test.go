package collection_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"

	"github.com/Finschia/finschia-sdk/x/collection"
	"github.com/Finschia/finschia-sdk/x/collection-token/class"
	linkerrors "github.com/Finschia/finschia-sdk/x/collection-token/errors"
)

func TestMsgSendFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(
		collection.NewFTCoin("00bab10c", math.OneInt()),
	)

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		to         sdk.AccAddress
		amount     []collection.Coin
		err        error
		panic      bool
	}{
		"valid msg": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			amount:     amount,
		},
		"invalid from": {
			contractID: contractID,
			to:         addrs[1],
			amount:     amount,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid contract id": {
			from:   addrs[0],
			to:     addrs[1],
			amount: amount,
			err:    class.ErrInvalidContractID,
		},
		"invalid to": {
			contractID: contractID,
			from:       addrs[0],
			amount:     amount,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"nil amount": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			amount: []collection.Coin{{
				TokenId: collection.NewFTID("00bab10c"),
			}},
			panic: true,
		},
		"zero amount": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			amount: []collection.Coin{{
				TokenId: collection.NewFTID("00bab10c"),
				Amount:  math.ZeroInt(),
			}},
			err: collection.ErrInvalidAmount,
		},
		"invalid token id": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			amount: []collection.Coin{{
				Amount: math.OneInt(),
			}},
			err: collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgSendFT{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				To:         tc.to.String(),
				Amount:     tc.amount,
			}

			if tc.panic {
				require.Panics(t, func() { msg.ValidateBasic() })
				return
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
		})
	}
}

func TestMsgOperatorSendFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(
		collection.NewFTCoin("00bab10c", math.OneInt()),
	)

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		to         sdk.AccAddress
		amount     []collection.Coin
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			amount:     amount,
		},
		"invalid operator": {
			contractID: contractID,
			from:       addrs[1],
			to:         addrs[2],
			amount:     amount,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid contract id": {
			operator: addrs[0],
			from:     addrs[1],
			to:       addrs[2],
			amount:   amount,
			err:      class.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			amount:     amount,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			amount:     amount,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid amount": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			amount: []collection.Coin{{
				TokenId: collection.NewFTID("00bab10c"),
				Amount:  math.ZeroInt(),
			}},
			err: collection.ErrInvalidAmount,
		},
		"invalid denom": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			amount: []collection.Coin{{
				Amount: math.OneInt(),
			}},
			err: collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgOperatorSendFT{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				To:         tc.to.String(),
				Amount:     tc.amount,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())
		})
	}
}

func TestMsgSendNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ids := []string{collection.NewNFTID(contractID, 1)}

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		to         sdk.AccAddress
		ids        []string
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			ids:        ids,
		},
		"invalid from": {
			contractID: contractID,
			to:         addrs[1],
			ids:        ids,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid contract id": {
			from: addrs[0],
			to:   addrs[1],
			ids:  ids,
			err:  class.ErrInvalidContractID,
		},
		"invalid to": {
			contractID: contractID,
			from:       addrs[0],
			ids:        ids,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty token ids": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			err:        collection.ErrEmptyField,
		},
		"invalid token ids": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			ids:        []string{""},
			err:        collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgSendNFT{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				To:         tc.to.String(),
				TokenIds:   tc.ids,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
		})
	}
}

func TestMsgOperatorSendNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ids := []string{collection.NewNFTID(contractID, 1)}

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		to         sdk.AccAddress
		ids        []string
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			ids:        ids,
		},
		"invalid operator": {
			contractID: contractID,
			from:       addrs[1],
			to:         addrs[2],
			ids:        ids,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid contract id": {
			operator: addrs[0],
			from:     addrs[1],
			to:       addrs[2],
			ids:      ids,
			err:      class.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			ids:        ids,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			ids:        ids,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty ids": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			err:        collection.ErrEmptyField,
		},
		"invalid id": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			to:         addrs[2],
			ids:        []string{""},
			err:        collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgOperatorSendNFT{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				To:         tc.to.String(),
				TokenIds:   tc.ids,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())
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
			contractID: contractID,
			holder:     addrs[0],
			operator:   addrs[1],
		},
		"invalid contract id": {
			holder:   addrs[0],
			operator: addrs[1],
			err:      class.ErrInvalidContractID,
		},
		"invalid holder": {
			contractID: contractID,
			operator:   addrs[1],
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty operator": {
			contractID: contractID,
			holder:     addrs[0],
			err:        sdkerrors.ErrInvalidAddress,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgAuthorizeOperator{
				ContractId: tc.contractID,
				Holder:     tc.holder.String(),
				Operator:   tc.operator.String(),
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.holder}, msg.GetSigners())
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
			contractID: contractID,
			holder:     addrs[0],
			operator:   addrs[1],
		},
		"invalid contract id": {
			holder:   addrs[0],
			operator: addrs[1],
			err:      class.ErrInvalidContractID,
		},
		"invalid holder": {
			contractID: contractID,
			operator:   addrs[1],
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty operator": {
			contractID: contractID,
			holder:     addrs[0],
			err:        sdkerrors.ErrInvalidAddress,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgRevokeOperator{
				ContractId: tc.contractID,
				Holder:     tc.holder.String(),
				Operator:   tc.operator.String(),
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.holder}, msg.GetSigners())
		})
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
		err        error
	}{
		"valid msg": {
			owner:      addrs[0],
			name:       name,
			baseImgURI: uri,
			meta:       meta,
		},
		"invalid owner": {
			name:       name,
			baseImgURI: uri,
			meta:       meta,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"long name": {
			owner:      addrs[0],
			name:       string(make([]rune, 21)),
			baseImgURI: uri,
			meta:       meta,
			err:        collection.ErrInvalidNameLength,
		},
		"invalid base image uri": {
			owner:      addrs[0],
			name:       name,
			baseImgURI: string(make([]rune, 1001)),
			meta:       meta,
			err:        collection.ErrInvalidBaseImgURILength,
		},
		"invalid meta": {
			owner:      addrs[0],
			name:       name,
			baseImgURI: uri,
			meta:       string(make([]rune, 1001)),
			err:        collection.ErrInvalidMetaLength,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgCreateContract{
				Owner: tc.owner.String(),
				Name:  tc.name,
				Uri:   tc.baseImgURI,
				Meta:  tc.meta,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.owner}, msg.GetSigners())
		})
	}
}

func TestMsgIssueFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := contractID
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
		amount     math.Int
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			name:       name,
			meta:       meta,
			decimals:   decimals,
			amount:     math.OneInt(),
		},
		"invalid contract id": {
			owner:    addrs[0],
			to:       addrs[1],
			name:     name,
			meta:     meta,
			decimals: decimals,
			amount:   math.OneInt(),
			err:      class.ErrInvalidContractID,
		},
		"invalid owner": {
			contractID: contractID,
			to:         addrs[1],
			name:       name,
			meta:       meta,
			decimals:   decimals,
			amount:     math.OneInt(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			contractID: contractID,
			owner:      addrs[0],
			name:       name,
			meta:       meta,
			decimals:   decimals,
			amount:     math.OneInt(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty name": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			meta:       meta,
			decimals:   decimals,
			amount:     math.OneInt(),
			err:        collection.ErrInvalidTokenName,
		},
		"long name": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			name:       string(make([]rune, 21)),
			meta:       meta,
			decimals:   decimals,
			amount:     math.OneInt(),
			err:        collection.ErrInvalidNameLength,
		},
		"invalid meta": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			name:       name,
			meta:       string(make([]rune, 1001)),
			decimals:   decimals,
			amount:     math.OneInt(),
			err:        collection.ErrInvalidMetaLength,
		},
		"invalid decimals": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			name:       name,
			meta:       meta,
			decimals:   19,
			amount:     math.OneInt(),
			err:        collection.ErrInvalidTokenDecimals,
		},
		"invalid decimals - negative": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			name:       name,
			meta:       meta,
			decimals:   -1,
			amount:     math.OneInt(),
			err:        collection.ErrInvalidTokenDecimals,
		},
		"daphne compat": {
			contractID: contractID,
			owner:      addrs[0],
			to:         addrs[1],
			name:       name,
			meta:       meta,
			amount:     math.OneInt(),
			err:        collection.ErrInvalidIssueFT,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgIssueFT{
				ContractId: tc.contractID,
				Owner:      tc.owner.String(),
				To:         tc.to.String(),
				Name:       tc.name,
				Meta:       tc.meta,
				Decimals:   tc.decimals,
				Amount:     tc.amount,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.owner}, msg.GetSigners())
		})
	}
}

func TestMsgIssueNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := contractID
	name := "tibetian fox"
	meta := "Tibetian Fox"
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		name       string
		meta       string
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			name:       name,
			meta:       meta,
		},
		"invalid contract id": {
			operator: addrs[0],
			name:     name,
			meta:     meta,
			err:      class.ErrInvalidContractID,
		},
		"invalid operator": {
			contractID: contractID,
			name:       name,
			meta:       meta,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"long name": {
			contractID: contractID,
			operator:   addrs[0],
			name:       string(make([]rune, 21)),
			meta:       meta,
			err:        collection.ErrInvalidNameLength,
		},
		"invalid meta": {
			contractID: contractID,
			operator:   addrs[0],
			name:       name,
			meta:       string(make([]rune, 1001)),
			err:        collection.ErrInvalidMetaLength,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgIssueNFT{
				ContractId: tc.contractID,
				Owner:      tc.operator.String(),
				Name:       tc.name,
				Meta:       tc.meta,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())
		})
	}
}

func TestMsgMintFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(
		collection.NewFTCoin("00bab10c", math.OneInt()),
	)

	contractID := contractID
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		to         sdk.AccAddress
		amount     []collection.Coin
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			amount:     amount,
		},
		// for daphne compatible
		"valid msg - zero amount": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
		},
		"invalid contract id": {
			operator: addrs[0],
			to:       addrs[1],
			amount:   amount,
			err:      class.ErrInvalidContractID,
		},
		"invalid operator": {
			contractID: contractID,
			to:         addrs[1],
			amount:     amount,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			contractID: contractID,
			operator:   addrs[0],
			amount:     amount,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid amount": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			amount: []collection.Coin{{
				TokenId: collection.NewFTID("00bab10c"),
				Amount:  math.ZeroInt(),
			}},
			err: collection.ErrInvalidAmount,
		},
		"invalid token id": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			amount: []collection.Coin{{
				Amount: math.OneInt(),
			}},
			err: collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgMintFT{
				ContractId: tc.contractID,
				From:       tc.operator.String(),
				To:         tc.to.String(),
				Amount:     tc.amount,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())
		})
	}
}

func TestMsgMintNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	params := []collection.MintNFTParam{{
		TokenType: contractID,
		Name:      "tibetian fox",
		Meta:      "Tibetian Fox",
	}}
	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		to         sdk.AccAddress
		params     []collection.MintNFTParam
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			params:     params,
		},
		"invalid contract id": {
			operator: addrs[0],
			to:       addrs[1],
			params:   params,
			err:      class.ErrInvalidContractID,
		},
		"invalid operator": {
			contractID: contractID,
			to:         addrs[1],
			params:     params,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			contractID: contractID,
			operator:   addrs[0],
			params:     params,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty params": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			err:        collection.ErrEmptyField,
		},
		"param of invalid token type": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			params: []collection.MintNFTParam{{
				Name: "tibetian fox",
			}},
			err: collection.ErrInvalidTokenType,
		},
		"param of empty name": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			params: []collection.MintNFTParam{{
				TokenType: contractID,
			}},
			err: collection.ErrInvalidTokenName,
		},
		"param of too long name": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			params: []collection.MintNFTParam{{
				TokenType: contractID,
				Name:      string(make([]rune, 21)),
			}},
			err: collection.ErrInvalidNameLength,
		},
		"param of invalid meta": {
			contractID: contractID,
			operator:   addrs[0],
			to:         addrs[1],
			params: []collection.MintNFTParam{{
				TokenType: contractID,
				Name:      "tibetian fox",
				Meta:      string(make([]rune, 1001)),
			}},
			err: collection.ErrInvalidMetaLength,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgMintNFT{
				ContractId: tc.contractID,
				From:       tc.operator.String(),
				To:         tc.to.String(),
				Params:     tc.params,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())
		})
	}
}

func TestMsgBurnFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(
		collection.NewFTCoin("00bab10c", math.OneInt()),
	)

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		amount     []collection.Coin
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			from:       addrs[0],
			amount:     amount,
		},
		// for daphne compatible
		"valid msg - zero amount": {
			contractID: contractID,
			from:       addrs[0],
		},
		"invalid contract id": {
			from:   addrs[0],
			amount: amount,
			err:    class.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: contractID,
			amount:     amount,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid token id": {
			contractID: contractID,
			from:       addrs[0],
			amount: []collection.Coin{{
				Amount: math.OneInt(),
			}},
			err: collection.ErrInvalidTokenID,
		},
		"invalid amount": {
			contractID: contractID,
			from:       addrs[0],
			amount: []collection.Coin{{
				TokenId: collection.NewFTID("00bab10c"),
				Amount:  math.ZeroInt(),
			}},
			err: collection.ErrInvalidAmount,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgBurnFT{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				Amount:     tc.amount,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
		})
	}
}

func TestMsgOperatorBurnFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	amount := collection.NewCoins(
		collection.NewFTCoin("00bab10c", math.OneInt()),
	)

	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		from       sdk.AccAddress
		amount     []collection.Coin
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			grantee:    addrs[0],
			from:       addrs[1],
			amount:     amount,
		},
		"invalid contract id": {
			grantee: addrs[0],
			from:    addrs[1],
			amount:  amount,
			err:     class.ErrInvalidContractID,
		},
		"invalid grantee": {
			contractID: contractID,
			from:       addrs[1],
			amount:     amount,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid from": {
			contractID: contractID,
			grantee:    addrs[0],
			amount:     amount,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid token id": {
			contractID: contractID,
			grantee:    addrs[0],
			from:       addrs[1],
			amount: []collection.Coin{{
				Amount: math.OneInt(),
			}},
			err: collection.ErrInvalidTokenID,
		},
		"invalid amount": {
			contractID: contractID,
			grantee:    addrs[0],
			from:       addrs[1],
			amount: []collection.Coin{{
				TokenId: collection.NewFTID("00bab10c"),
				Amount:  math.ZeroInt(),
			}},
			err: collection.ErrInvalidAmount,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgOperatorBurnFT{
				ContractId: tc.contractID,
				Operator:   tc.grantee.String(),
				From:       tc.from.String(),
				Amount:     tc.amount,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.grantee}, msg.GetSigners())
		})
	}
}

func TestMsgBurnNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ids := []string{collection.NewNFTID(contractID, 1)}

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		ids        []string
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			from:       addrs[0],
			ids:        ids,
		},
		"invalid contract id": {
			from: addrs[0],
			ids:  ids,
			err:  class.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: contractID,
			ids:        ids,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty ids": {
			contractID: contractID,
			from:       addrs[0],
			err:        collection.ErrEmptyField,
		},
		"invalid id": {
			contractID: contractID,
			from:       addrs[0],
			ids:        []string{""},
			err:        collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgBurnNFT{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				TokenIds:   tc.ids,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
		})
	}
}

func TestMsgOperatorBurnNFT(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ids := []string{collection.NewNFTID(contractID, 1)}

	testCases := map[string]struct {
		contractID string
		grantee    sdk.AccAddress
		from       sdk.AccAddress
		ids        []string
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			grantee:    addrs[0],
			from:       addrs[1],
			ids:        ids,
		},
		"invalid contract id": {
			grantee: addrs[0],
			from:    addrs[1],
			ids:     ids,
			err:     class.ErrInvalidContractID,
		},
		"invalid grantee": {
			contractID: contractID,
			from:       addrs[1],
			ids:        ids,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid from": {
			contractID: contractID,
			grantee:    addrs[0],
			ids:        ids,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"empty ids": {
			contractID: contractID,
			grantee:    addrs[0],
			from:       addrs[1],
			err:        collection.ErrEmptyField,
		},
		"invalid id": {
			contractID: contractID,
			grantee:    addrs[0],
			from:       addrs[0],
			ids:        []string{""},
			err:        collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgOperatorBurnNFT{
				ContractId: tc.contractID,
				Operator:   tc.grantee.String(),
				From:       tc.from.String(),
				TokenIds:   tc.ids,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
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

	changes := []collection.Attribute{{Key: collection.AttributeKeyName.String(), Value: "New test"}}
	testCases := map[string]struct {
		contractID string
		owner      sdk.AccAddress
		tokenType  string
		tokenIndex string
		changes    []collection.Attribute
		err        error
	}{
		"valid contract modification": {
			contractID: contractID,
			owner:      addrs[0],
			changes:    changes,
		},
		"valid token class modification": {
			contractID: contractID,
			tokenType:  contractID,
			owner:      addrs[0],
			changes:    changes,
		},
		"invalid nft class modification": {
			contractID: contractID,
			tokenType:  contractID,
			tokenIndex: "00000000",
			owner:      addrs[0],
			changes:    changes,
			err:        collection.ErrInvalidTokenIndex,
		},
		"valid nft modification": {
			contractID: contractID,
			tokenType:  contractID,
			tokenIndex: contractID,
			owner:      addrs[0],
			changes:    changes,
		},
		"invalid contract id": {
			owner:   addrs[0],
			changes: changes,
			err:     class.ErrInvalidContractID,
		},
		"invalid owner": {
			contractID: contractID,
			changes:    changes,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid key of change": {
			contractID: contractID,
			owner:      addrs[0],
			changes:    []collection.Attribute{{Key: strings.ToUpper(collection.AttributeKeyName.String()), Value: "tt"}},
			err:        collection.ErrInvalidChangesField,
		},
		"invalid value of change": {
			contractID: contractID,
			owner:      addrs[0],
			changes:    []collection.Attribute{{Key: collection.AttributeKeyName.String(), Value: string(make([]rune, 21))}},
			err:        collection.ErrInvalidNameLength,
		},
		"empty changes": {
			contractID: contractID,
			owner:      addrs[0],
			err:        collection.ErrEmptyChanges,
		},
		"too many changes": {
			contractID: contractID,
			owner:      addrs[0],
			changes:    make([]collection.Attribute, 101),
			err:        collection.ErrInvalidChangesFieldCount,
		},
		"duplicated changes": {
			contractID: contractID,
			owner:      addrs[0],
			changes: []collection.Attribute{
				{Key: collection.AttributeKeyBaseImgURI.String(), Value: "hello"},
				{Key: collection.AttributeKeyURI.String(), Value: "world"},
			},
			err: collection.ErrDuplicateChangesField,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgModify{
				ContractId: tc.contractID,
				TokenType:  tc.tokenType,
				TokenIndex: tc.tokenIndex,
				Owner:      tc.owner.String(),
				Changes:    tc.changes,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.owner}, msg.GetSigners())
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
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			permission: collection.LegacyPermissionMint.String(),
		},
		"invalid contract id": {
			from:       addrs[0],
			to:         addrs[1],
			permission: collection.LegacyPermissionMint.String(),
			err:        class.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: contractID,
			to:         addrs[1],
			permission: collection.LegacyPermissionMint.String(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid to": {
			contractID: contractID,
			from:       addrs[0],
			permission: collection.LegacyPermissionMint.String(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid permission": {
			contractID: contractID,
			from:       addrs[0],
			to:         addrs[1],
			err:        linkerrors.ErrInvalidPermission,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgGrantPermission{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				To:         tc.to.String(),
				Permission: tc.permission,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
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
			contractID: contractID,
			from:       addrs[0],
			permission: collection.LegacyPermissionMint.String(),
		},
		"invalid contract id": {
			from:       addrs[0],
			permission: collection.LegacyPermissionMint.String(),
			err:        class.ErrInvalidContractID,
		},
		"invalid from": {
			contractID: contractID,
			permission: collection.LegacyPermissionMint.String(),
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid permission": {
			contractID: contractID,
			from:       addrs[0],
			err:        linkerrors.ErrInvalidPermission,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgRevokePermission{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				Permission: tc.permission,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
		})
	}
}

func TestMsgAttach(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := contractID
	tokenIDs := []string{
		collection.NewNFTID(contractID, 1),
		collection.NewNFTID("fee1dead", 1),
	}

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		tokenID    string
		toTokenID  string
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			from:       addrs[0],
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[1],
		},
		"invalid from": {
			contractID: contractID,
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[1],
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid contract id": {
			from:      addrs[0],
			tokenID:   tokenIDs[0],
			toTokenID: tokenIDs[1],
			err:       class.ErrInvalidContractID,
		},
		"invalid token id": {
			contractID: contractID,
			from:       addrs[0],
			toTokenID:  tokenIDs[1],
			err:        collection.ErrInvalidTokenID,
		},
		"invalid to id": {
			contractID: contractID,
			from:       addrs[0],
			tokenID:    tokenIDs[0],
			err:        collection.ErrInvalidTokenID,
		},
		"to itself": {
			contractID: contractID,
			from:       addrs[0],
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[0],
			err:        collection.ErrCannotAttachToItself,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgAttach{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				TokenId:    tc.tokenID,
				ToTokenId:  tc.toTokenID,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
		})
	}
}

func TestMsgDetach(t *testing.T) {
	addrs := make([]sdk.AccAddress, 1)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	contractID := contractID
	tokenID := collection.NewNFTID(contractID, 1)

	testCases := map[string]struct {
		contractID string
		from       sdk.AccAddress
		tokenID    string
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			from:       addrs[0],
			tokenID:    tokenID,
		},
		"invalid from": {
			contractID: contractID,
			tokenID:    tokenID,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid contract id": {
			from:    addrs[0],
			tokenID: tokenID,
			err:     class.ErrInvalidContractID,
		},
		"invalid token id": {
			contractID: contractID,
			from:       addrs[0],
			err:        collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgDetach{
				ContractId: tc.contractID,
				From:       tc.from.String(),
				TokenId:    tc.tokenID,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.from}, msg.GetSigners())
		})
	}
}

func TestMsgOperatorAttach(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	tokenIDs := []string{
		collection.NewNFTID(contractID, 1),
		collection.NewNFTID("fee1dead", 1),
	}

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		tokenID    string
		toTokenID  string
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[1],
		},
		"empty operator": {
			contractID: contractID,
			from:       addrs[1],
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[1],
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid from": {
			contractID: contractID,
			operator:   addrs[0],
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[1],
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid contract id": {
			operator:  addrs[0],
			from:      addrs[1],
			tokenID:   tokenIDs[0],
			toTokenID: tokenIDs[1],
			err:       class.ErrInvalidContractID,
		},
		"invalid token id": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			toTokenID:  tokenIDs[1],
			err:        collection.ErrInvalidTokenID,
		},
		"invalid to id": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			tokenID:    tokenIDs[0],
			err:        collection.ErrInvalidTokenID,
		},
		"to itself": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			tokenID:    tokenIDs[0],
			toTokenID:  tokenIDs[0],
			err:        collection.ErrCannotAttachToItself,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgOperatorAttach{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				TokenId:    tc.tokenID,
				ToTokenId:  tc.toTokenID,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())
		})
	}
}

func TestMsgOperatorDetach(t *testing.T) {
	addrs := make([]sdk.AccAddress, 2)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	tokenID := collection.NewNFTID(contractID, 1)

	testCases := map[string]struct {
		contractID string
		operator   sdk.AccAddress
		from       sdk.AccAddress
		tokenID    string
		err        error
	}{
		"valid msg": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			tokenID:    tokenID,
		},
		"empty operator": {
			contractID: contractID,
			from:       addrs[1],
			tokenID:    tokenID,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid from": {
			contractID: contractID,
			operator:   addrs[0],
			tokenID:    tokenID,
			err:        sdkerrors.ErrInvalidAddress,
		},
		"invalid contract id": {
			operator: addrs[0],
			from:     addrs[1],
			tokenID:  tokenID,
			err:      class.ErrInvalidContractID,
		},
		"invalid token id": {
			contractID: contractID,
			operator:   addrs[0],
			from:       addrs[1],
			err:        collection.ErrInvalidTokenID,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			msg := collection.MsgOperatorDetach{
				ContractId: tc.contractID,
				Operator:   tc.operator.String(),
				From:       tc.from.String(),
				TokenId:    tc.tokenID,
			}

			require.ErrorIs(t, msg.ValidateBasic(), tc.err)
			if tc.err != nil {
				return
			}

			require.Equal(t, []sdk.AccAddress{tc.operator}, msg.GetSigners())
		})
	}
}

func TestAminoJSON(t *testing.T) {
	tx := legacytx.StdTx{}
	ftClassID := "00bab10c"
	legacyAmino := codec.NewLegacyAmino()
	collection.RegisterLegacyAminoCodec(legacyAmino)
	legacytx.RegressionTestingAminoCodec = legacyAmino

	addrs := make([]sdk.AccAddress, 3)
	for i := range addrs {
		addrs[i] = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	}

	ftAmount := collection.NewCoins(collection.NewFTCoin(ftClassID, math.NewInt(1000000)))
	tokenIds := []string{collection.NewNFTID(contractID, 1)}
	const name = "tibetian fox"
	nftParams := []collection.MintNFTParam{{
		TokenType: contractID,
		Name:      name,
		Meta:      "Tibetian Fox",
	}}

	testCase := map[string]struct {
		msg          legacytx.LegacyMsg
		expectedType string
		expected     string
	}{
		"MsgSendFT": {
			&collection.MsgSendFT{
				ContractId: contractID,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				Amount:     ftAmount,
			},
			"/lbm.collection.v1.MsgSendFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSendFT\",\"value\":{\"amount\":[{\"amount\":\"1000000\",\"token_id\":\"00bab10c00000000\"}],\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgOperatorSendFT": {
			&collection.MsgOperatorSendFT{
				ContractId: contractID,
				Operator:   addrs[0].String(),
				From:       addrs[1].String(),
				To:         addrs[2].String(),
				Amount:     ftAmount,
			},
			"/lbm.collection.v1.MsgOperatorSendFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgOperatorSendFT\",\"value\":{\"amount\":[{\"amount\":\"1000000\",\"token_id\":\"00bab10c00000000\"}],\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"operator\":\"%s\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String(), addrs[2].String()),
		},
		"MsgSendNFT": {
			&collection.MsgSendNFT{
				ContractId: contractID,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				TokenIds:   tokenIds,
			},
			"/lbm.collection.v1.MsgSendNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSendNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"to\":\"%s\",\"token_ids\":[\"deadbeef00000001\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgOperatorSendNFT": {
			&collection.MsgOperatorSendNFT{
				ContractId: contractID,
				Operator:   addrs[0].String(),
				From:       addrs[1].String(),
				To:         addrs[2].String(),
				TokenIds:   tokenIds,
			},
			"/lbm.collection.v1.MsgOperatorSendNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgOperatorSendNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"operator\":\"%s\",\"to\":\"%s\",\"token_ids\":[\"deadbeef00000001\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String(), addrs[2].String()),
		},
		"MsgAuthorizeOperator": {
			&collection.MsgAuthorizeOperator{
				ContractId: contractID,
				Holder:     addrs[0].String(),
				Operator:   addrs[1].String(),
			},
			"/lbm.collection.v1.MsgAuthorizeOperator",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgAuthorizeOperator\",\"value\":{\"contract_id\":\"deadbeef\",\"holder\":\"%s\",\"operator\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgRevokeOperator": {
			&collection.MsgRevokeOperator{
				ContractId: contractID,
				Holder:     addrs[0].String(),
				Operator:   addrs[1].String(),
			},
			"/lbm.collection.v1.MsgRevokeOperator",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgRevokeOperator\",\"value\":{\"contract_id\":\"deadbeef\",\"holder\":\"%s\",\"operator\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgCreateContract": {
			&collection.MsgCreateContract{
				Owner: addrs[0].String(),
				Name:  "Test Contract",
				Uri:   "http://image.url",
				Meta:  "This is test",
			},
			"/lbm.collection.v1.MsgCreateContract",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgCreateContract\",\"value\":{\"meta\":\"This is test\",\"name\":\"Test Contract\",\"owner\":\"%s\",\"uri\":\"http://image.url\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgIssueFT": {
			&collection.MsgIssueFT{
				ContractId: contractID,
				Name:       "Test FT",
				Meta:       "This is FT Meta",
				Decimals:   6,
				Mintable:   false,
				Owner:      addrs[0].String(),
				To:         addrs[1].String(),
				Amount:     math.NewInt(1000000),
			},
			"/lbm.collection.v1.MsgIssueFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgIssueFT\",\"value\":{\"amount\":\"1000000\",\"contract_id\":\"deadbeef\",\"decimals\":6,\"meta\":\"This is FT Meta\",\"name\":\"Test FT\",\"owner\":\"%s\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgIssueNFT": {
			&collection.MsgIssueNFT{
				ContractId: contractID,
				Name:       "Test NFT",
				Meta:       "This is NFT Meta",
				Owner:      addrs[0].String(),
			},
			"/lbm.collection.v1.MsgIssueNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgIssueNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"meta\":\"This is NFT Meta\",\"name\":\"Test NFT\",\"owner\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgMintFT": {
			&collection.MsgMintFT{
				ContractId: contractID,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				Amount:     ftAmount,
			},
			"/lbm.collection.v1.MsgMintFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgMintFT\",\"value\":{\"amount\":[{\"amount\":\"1000000\",\"token_id\":\"00bab10c00000000\"}],\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgMintNFT": {
			&collection.MsgMintNFT{
				ContractId: contractID,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				Params:     nftParams,
			},
			"/lbm.collection.v1.MsgMintNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgMintNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"params\":[{\"meta\":\"Tibetian Fox\",\"name\":\"tibetian fox\",\"token_type\":\"deadbeef\"}],\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgBurnFT": {
			&collection.MsgBurnFT{
				ContractId: contractID,
				From:       addrs[0].String(),
				Amount:     ftAmount,
			},
			"/lbm.collection.v1.MsgBurnFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgBurnFT\",\"value\":{\"amount\":[{\"amount\":\"1000000\",\"token_id\":\"00bab10c00000000\"}],\"contract_id\":\"deadbeef\",\"from\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgOperatorBurnFT": {
			&collection.MsgOperatorBurnFT{
				ContractId: contractID,
				Operator:   addrs[0].String(),
				From:       addrs[1].String(),
				Amount:     ftAmount,
			},
			"/lbm.collection.v1.MsgOperatorBurnFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgOperatorBurnFT\",\"value\":{\"amount\":[{\"amount\":\"1000000\",\"token_id\":\"00bab10c00000000\"}],\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"operator\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String()),
		},
		"MsgBurnNFT": {
			&collection.MsgBurnNFT{
				ContractId: contractID,
				From:       addrs[0].String(),
				TokenIds:   tokenIds,
			},
			"/lbm.collection.v1.MsgBurnNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgBurnNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"token_ids\":[\"deadbeef00000001\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgOperatorBurnNFT": {
			&collection.MsgOperatorBurnNFT{
				ContractId: contractID,
				Operator:   addrs[0].String(),
				From:       addrs[1].String(),
				TokenIds:   tokenIds,
			},
			"/lbm.collection.v1.MsgOperatorBurnNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgOperatorBurnNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"operator\":\"%s\",\"token_ids\":[\"deadbeef00000001\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String()),
		},
		"MsgModify": {
			&collection.MsgModify{
				ContractId: contractID,
				Owner:      addrs[0].String(),
				TokenType:  "NewType",
				TokenIndex: contractID,
				Changes:    []collection.Attribute{{Key: "name", Value: "New test"}},
			},
			"/lbm.collection.v1.MsgModify",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgModify\",\"value\":{\"changes\":[{\"key\":\"name\",\"value\":\"New test\"}],\"contract_id\":\"deadbeef\",\"owner\":\"%s\",\"token_index\":\"deadbeef\",\"token_type\":\"NewType\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgGrantPermission": {
			&collection.MsgGrantPermission{
				ContractId: contractID,
				From:       addrs[0].String(),
				To:         addrs[1].String(),
				Permission: collection.LegacyPermissionMint.String(),
			},
			"/lbm.collection.v1.MsgGrantPermission",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgGrantPermission\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"permission\":\"mint\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String(), addrs[1].String()),
		},
		"MsgRevokePermission": {
			&collection.MsgRevokePermission{
				ContractId: contractID,
				From:       addrs[0].String(),
				Permission: collection.LegacyPermissionMint.String(),
			},
			"/lbm.collection.v1.MsgRevokePermission",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgRevokePermission\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"permission\":\"mint\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgAttach": {
			&collection.MsgAttach{
				ContractId: contractID,
				From:       addrs[0].String(),
				TokenId:    collection.NewNFTID(contractID, 1),
				ToTokenId:  collection.NewNFTID("fee1dead", 1),
			},
			"/lbm.collection.v1.MsgAttach",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgAttach\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"to_token_id\":\"fee1dead00000001\",\"token_id\":\"deadbeef00000001\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgDetach": {
			&collection.MsgDetach{
				ContractId: contractID,
				From:       addrs[0].String(),
				TokenId:    collection.NewNFTID("fee1dead", 1),
			},
			"/lbm.collection.v1.MsgDetach",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgDetach\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"token_id\":\"fee1dead00000001\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0].String()),
		},
		"MsgOperatorAttach": {
			&collection.MsgOperatorAttach{
				ContractId: contractID,
				Operator:   addrs[0].String(),
				From:       addrs[1].String(),
				TokenId:    collection.NewNFTID(contractID, 1),
				ToTokenId:  collection.NewNFTID("fee1dead", 1),
			},
			"/lbm.collection.v1.MsgOperatorAttach",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgOperatorAttach\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"operator\":\"%s\",\"to_token_id\":\"fee1dead00000001\",\"token_id\":\"deadbeef00000001\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String()),
		},
		"MsgOperatorDetach": {
			&collection.MsgOperatorDetach{
				ContractId: contractID,
				Operator:   addrs[0].String(),
				From:       addrs[1].String(),
				TokenId:    collection.NewNFTID("fee1dead", 1),
			},
			"/lbm.collection.v1.MsgOperatorDetach",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgOperatorDetach\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"operator\":\"%s\",\"token_id\":\"fee1dead00000001\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1].String(), addrs[0].String()),
		},
	}

	for name, tc := range testCase {
		tc := tc
		t.Run(name, func(t *testing.T) {
			tx.Msgs = []sdk.Msg{tc.msg}
			require.Equal(t, tc.expected, string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{tc.msg}, "memo")))
		})
	}
}
