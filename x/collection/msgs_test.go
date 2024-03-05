package collection_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"

	"github.com/Finschia/finschia-sdk/x/collection"
)

func TestAminoJSON(t *testing.T) {
	tx := legacytx.StdTx{}
	legacyAmino := codec.NewLegacyAmino()
	collection.RegisterLegacyAminoCodec(legacyAmino)
	legacytx.RegressionTestingAminoCodec = legacyAmino
	ac := authcodec.NewBech32Codec("link")

	addrs := make([]string, 3)
	for i := range addrs {
		addr, err := ac.BytesToString(secp256k1.GenPrivKey().PubKey().Address())
		require.NoError(t, err)
		addrs[i] = addr
	}

	tokenIds := []string{collection.NewNFTID(contractID, 1)}
	const name = "tibetian fox"
	nftParams := []collection.MintNFTParam{{
		TokenType: contractID,
		Name:      name,
		Meta:      "Tibetian Fox",
	}}

	testCase := map[string]struct {
		msg          sdk.Msg
		expectedType string
		expected     string
	}{
		"MsgSendNFT": {
			&collection.MsgSendNFT{
				ContractId: contractID,
				From:       addrs[0],
				To:         addrs[1],
				TokenIds:   tokenIds,
			},
			"/lbm.collection.v1.MsgSendNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgSendNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"to\":\"%s\",\"token_ids\":[\"deadbeef00000001\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0], addrs[1]),
		},
		"MsgOperatorSendNFT": {
			&collection.MsgOperatorSendNFT{
				ContractId: contractID,
				Operator:   addrs[0],
				From:       addrs[1],
				To:         addrs[2],
				TokenIds:   tokenIds,
			},
			"/lbm.collection.v1.MsgOperatorSendNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgOperatorSendNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"operator\":\"%s\",\"to\":\"%s\",\"token_ids\":[\"deadbeef00000001\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1], addrs[0], addrs[2]),
		},
		"MsgAuthorizeOperator": {
			&collection.MsgAuthorizeOperator{
				ContractId: contractID,
				Holder:     addrs[0],
				Operator:   addrs[1],
			},
			"/lbm.collection.v1.MsgAuthorizeOperator",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgAuthorizeOperator\",\"value\":{\"contract_id\":\"deadbeef\",\"holder\":\"%s\",\"operator\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0], addrs[1]),
		},
		"MsgRevokeOperator": {
			&collection.MsgRevokeOperator{
				ContractId: contractID,
				Holder:     addrs[0],
				Operator:   addrs[1],
			},
			"/lbm.collection.v1.MsgRevokeOperator",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgRevokeOperator\",\"value\":{\"contract_id\":\"deadbeef\",\"holder\":\"%s\",\"operator\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0], addrs[1]),
		},
		"MsgCreateContract": {
			&collection.MsgCreateContract{
				Owner: addrs[0],
				Name:  "Test Contract",
				Uri:   "http://image.url",
				Meta:  "This is test",
			},
			"/lbm.collection.v1.MsgCreateContract",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgCreateContract\",\"value\":{\"meta\":\"This is test\",\"name\":\"Test Contract\",\"owner\":\"%s\",\"uri\":\"http://image.url\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0]),
		},
		"MsgIssueNFT": {
			&collection.MsgIssueNFT{
				ContractId: contractID,
				Name:       "Test NFT",
				Meta:       "This is NFT Meta",
				Owner:      addrs[0],
			},
			"/lbm.collection.v1.MsgIssueNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgIssueNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"meta\":\"This is NFT Meta\",\"name\":\"Test NFT\",\"owner\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0]),
		},
		"MsgMintNFT": {
			&collection.MsgMintNFT{
				ContractId: contractID,
				From:       addrs[0],
				To:         addrs[1],
				Params:     nftParams,
			},
			"/lbm.collection.v1.MsgMintNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgMintNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"params\":[{\"meta\":\"Tibetian Fox\",\"name\":\"tibetian fox\",\"token_type\":\"deadbeef\"}],\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0], addrs[1]),
		},
		"MsgBurnNFT": {
			&collection.MsgBurnNFT{
				ContractId: contractID,
				From:       addrs[0],
				TokenIds:   tokenIds,
			},
			"/lbm.collection.v1.MsgBurnNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgBurnNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"token_ids\":[\"deadbeef00000001\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0]),
		},
		"MsgOperatorBurnNFT": {
			&collection.MsgOperatorBurnNFT{
				ContractId: contractID,
				Operator:   addrs[0],
				From:       addrs[1],
				TokenIds:   tokenIds,
			},
			"/lbm.collection.v1.MsgOperatorBurnNFT",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/MsgOperatorBurnNFT\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"operator\":\"%s\",\"token_ids\":[\"deadbeef00000001\"]}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[1], addrs[0]),
		},
		"MsgModify": {
			&collection.MsgModify{
				ContractId: contractID,
				Owner:      addrs[0],
				TokenType:  "NewType",
				TokenIndex: contractID,
				Changes:    []collection.Attribute{{Key: "name", Value: "New test"}},
			},
			"/lbm.collection.v1.MsgModify",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgModify\",\"value\":{\"changes\":[{\"key\":\"name\",\"value\":\"New test\"}],\"contract_id\":\"deadbeef\",\"owner\":\"%s\",\"token_index\":\"deadbeef\",\"token_type\":\"NewType\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0]),
		},
		"MsgGrantPermission": {
			&collection.MsgGrantPermission{
				ContractId: contractID,
				To:         addrs[1],
				From:       addrs[0],
				Permission: collection.LegacyPermissionMint.String(),
			},
			"/lbm.collection.v1.MsgGrantPermission",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgGrantPermission\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"permission\":\"mint\",\"to\":\"%s\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0], addrs[1]),
		},
		"MsgRevokePermission": {
			&collection.MsgRevokePermission{
				ContractId: contractID,
				From:       addrs[0],
				Permission: collection.LegacyPermissionMint.String(),
			},
			"/lbm.collection.v1.MsgRevokePermission",
			fmt.Sprintf("{\"account_number\":\"1\",\"chain_id\":\"foo\",\"fee\":{\"amount\":[],\"gas\":\"0\"},\"memo\":\"memo\",\"msgs\":[{\"type\":\"lbm-sdk/collection/MsgRevokePermission\",\"value\":{\"contract_id\":\"deadbeef\",\"from\":\"%s\",\"permission\":\"mint\"}}],\"sequence\":\"1\",\"timeout_height\":\"1\"}", addrs[0]),
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
