package authz_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	cdctypes "github.com/Finschia/finschia-sdk/codec/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/auth/legacy/legacytx"
	"github.com/Finschia/finschia-sdk/x/authz"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	stakingtypes "github.com/Finschia/finschia-sdk/x/staking/types"
)

var (
	coinsPos = sdk.NewCoins(sdk.NewInt64Coin("steak", 100))
	granter  = sdk.AccAddress("_______granter______")
	grantee  = sdk.AccAddress("_______grantee______")
)

func TestMsgExecAuthorized(t *testing.T) {
	tests := []struct {
		title      string
		grantee    sdk.AccAddress
		msgs       []sdk.Msg
		expectPass bool
	}{
		{"nil grantee address", nil, []sdk.Msg{}, false},
		{"zero-messages test: should fail", grantee, []sdk.Msg{}, false},
		{"valid test: msg type", grantee, []sdk.Msg{
			&banktypes.MsgSend{
				Amount:      sdk.NewCoins(sdk.NewInt64Coin("steak", 2)),
				FromAddress: granter.String(),
				ToAddress:   grantee.String(),
			},
		}, true},
	}
	for i, tc := range tests {
		msg := authz.NewMsgExec(tc.grantee, tc.msgs)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgRevokeAuthorization(t *testing.T) {
	tests := []struct {
		title            string
		granter, grantee sdk.AccAddress
		msgType          string
		expectPass       bool
	}{
		{"nil Granter address", nil, grantee, "hello", false},
		{"nil Grantee address", granter, nil, "hello", false},
		{"nil Granter and Grantee address", nil, nil, "hello", false},
		{"valid test case", granter, grantee, "hello", true},
	}
	for i, tc := range tests {
		msg := authz.NewMsgRevoke(tc.granter, tc.grantee, tc.msgType)
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgGrantAuthorization(t *testing.T) {
	tests := []struct {
		title            string
		granter, grantee sdk.AccAddress
		authorization    authz.Authorization
		expiration       time.Time
		expectErr        bool
		expectPass       bool
	}{
		{"nil granter address", nil, grantee, &banktypes.SendAuthorization{SpendLimit: coinsPos}, time.Now(), false, false},
		{"nil grantee address", granter, nil, &banktypes.SendAuthorization{SpendLimit: coinsPos}, time.Now(), false, false},
		{"nil granter and grantee address", nil, nil, &banktypes.SendAuthorization{SpendLimit: coinsPos}, time.Now(), false, false},
		{"nil authorization", granter, grantee, nil, time.Now(), true, false},
		{"valid test case", granter, grantee, &banktypes.SendAuthorization{SpendLimit: coinsPos}, time.Now().AddDate(0, 1, 0), false, true},
		{"past time", granter, grantee, &banktypes.SendAuthorization{SpendLimit: coinsPos}, time.Now().AddDate(0, 0, -1), false, true}, // TODO need 0.45
	}
	for i, tc := range tests {
		msg, err := authz.NewMsgGrant(
			tc.granter, tc.grantee, tc.authorization, tc.expiration,
		)
		if !tc.expectErr {
			require.NoError(t, err)
		} else {
			continue
		}
		if tc.expectPass {
			require.NoError(t, msg.ValidateBasic(), "test: %v", i)
		} else {
			require.Error(t, msg.ValidateBasic(), "test: %v", i)
		}
	}
}

func TestMsgGrantGetAuthorization(t *testing.T) {
	require := require.New(t)

	m := authz.MsgGrant{}
	require.Nil(m.GetAuthorization())

	g := authz.GenericAuthorization{Msg: "some_type"}
	var err error
	m.Grant.Authorization, err = cdctypes.NewAnyWithValue(&g)
	require.NoError(err)
	require.Equal(m.GetAuthorization(), &g)

	g = authz.GenericAuthorization{Msg: "some_type2"}
	m.SetAuthorization(&g)
	require.Equal(m.GetAuthorization(), &g)
}

func TestAminoJSON(t *testing.T) {
	tx := legacytx.StdTx{}
	var msg legacytx.LegacyMsg
	someDate := time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC)
	msgSend := banktypes.MsgSend{FromAddress: "link1ghi", ToAddress: "link1jkl"}
	typeURL := sdk.MsgTypeURL(&msgSend)
	msgSendAny, err := cdctypes.NewAnyWithValue(&msgSend)
	require.NoError(t, err)
	grant, err := authz.NewGrant(someDate, authz.NewGenericAuthorization(typeURL), someDate.Add(time.Hour))
	require.NoError(t, err)
	sendGrant, err := authz.NewGrant(someDate, banktypes.NewSendAuthorization(sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1000)))), someDate.Add(time.Hour))
	require.NoError(t, err)
	valAddr, err := sdk.ValAddressFromBech32("linkvaloper1hcnhauxt4crgz04zwvdl60a9pk9wzzqe8uelc0")
	require.NoError(t, err)
	stakingAuth, err := stakingtypes.NewStakeAuthorization([]sdk.ValAddress{valAddr}, nil, stakingtypes.AuthorizationType_AUTHORIZATION_TYPE_DELEGATE, &sdk.Coin{Denom: "stake", Amount: sdk.NewInt(1000)})
	require.NoError(t, err)
	delegateGrant, err := authz.NewGrant(someDate, stakingAuth, someDate.Add(time.Hour))
	require.NoError(t, err)

	// Amino JSON encoding has changed in authz since v0.46.
	// Before, it was outputting something like:
	// `{"account_number":"1","chain_id":"foo","fee":{"amount":[],"gas":"0"},"memo":"memo","msgs":[{"grant":{"authorization":{"msg":"/cosmos.bank.v1beta1.MsgSend"},"expiration":"0001-01-01T02:01:01.000000001Z"},"grantee":"link1def","granter":"link1abc"}],"sequence":"1","timeout_height":"1"}`
	//
	// This was a bug. Now, it's as below, See how there's `type` & `value` fields.
	// ref: https://github.com/cosmos/cosmos-sdk/issues/11190
	// ref: https://github.com/cosmos/cosmjs/issues/1026
	msg = &authz.MsgGrant{Granter: "link1abc", Grantee: "link1def", Grant: grant}
	tx.Msgs = []sdk.Msg{msg}
	require.Equal(t,
		`{"account_number":"1","chain_id":"foo","fee":{"amount":[],"gas":"0"},"memo":"memo","msgs":[{"type":"cosmos-sdk/MsgGrant","value":{"grant":{"authorization":{"type":"cosmos-sdk/GenericAuthorization","value":{"msg":"/cosmos.bank.v1beta1.MsgSend"}},"expiration":"0001-01-01T02:01:01.000000001Z"},"grantee":"link1def","granter":"link1abc"}}],"sequence":"1","timeout_height":"1"}`,
		string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{msg}, "memo")),
	)

	msg = &authz.MsgGrant{Granter: "link1abc", Grantee: "link1def", Grant: sendGrant}
	tx.Msgs = []sdk.Msg{msg}
	require.Equal(t,
		`{"account_number":"1","chain_id":"foo","fee":{"amount":[],"gas":"0"},"memo":"memo","msgs":[{"type":"cosmos-sdk/MsgGrant","value":{"grant":{"authorization":{"type":"cosmos-sdk/SendAuthorization","value":{"spend_limit":[{"amount":"1000","denom":"stake"}]}},"expiration":"0001-01-01T02:01:01.000000001Z"},"grantee":"link1def","granter":"link1abc"}}],"sequence":"1","timeout_height":"1"}`,
		string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{msg}, "memo")),
	)

	msg = &authz.MsgGrant{Granter: "link1abc", Grantee: "link1def", Grant: delegateGrant}
	tx.Msgs = []sdk.Msg{msg}
	require.Equal(t,
		`{"account_number":"1","chain_id":"foo","fee":{"amount":[],"gas":"0"},"memo":"memo","msgs":[{"type":"cosmos-sdk/MsgGrant","value":{"grant":{"authorization":{"type":"cosmos-sdk/StakeAuthorization","value":{"Validators":{"type":"cosmos-sdk/StakeAuthorization/AllowList","value":{"allow_list":{"address":["linkvaloper1hcnhauxt4crgz04zwvdl60a9pk9wzzqe8uelc0"]}}},"authorization_type":1,"max_tokens":{"amount":"1000","denom":"stake"}}},"expiration":"0001-01-01T02:01:01.000000001Z"},"grantee":"link1def","granter":"link1abc"}}],"sequence":"1","timeout_height":"1"}`,
		string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{msg}, "memo")),
	)

	msg = &authz.MsgRevoke{Granter: "link1abc", Grantee: "link1def", MsgTypeUrl: typeURL}
	tx.Msgs = []sdk.Msg{msg}
	require.Equal(t,
		`{"account_number":"1","chain_id":"foo","fee":{"amount":[],"gas":"0"},"memo":"memo","msgs":[{"type":"cosmos-sdk/MsgRevoke","value":{"grantee":"link1def","granter":"link1abc","msg_type_url":"/cosmos.bank.v1beta1.MsgSend"}}],"sequence":"1","timeout_height":"1"}`,
		string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{msg}, "memo")),
	)

	msg = &authz.MsgExec{Grantee: "link1def", Msgs: []*cdctypes.Any{msgSendAny}}
	tx.Msgs = []sdk.Msg{msg}
	require.Equal(t,
		`{"account_number":"1","chain_id":"foo","fee":{"amount":[],"gas":"0"},"memo":"memo","msgs":[{"type":"cosmos-sdk/MsgExec","value":{"grantee":"link1def","msgs":[{"type":"cosmos-sdk/MsgSend","value":{"amount":[],"from_address":"link1ghi","to_address":"link1jkl"}}]}}],"sequence":"1","timeout_height":"1"}`,
		string(legacytx.StdSignBytes("foo", 1, 1, 1, legacytx.StdFee{}, []sdk.Msg{msg}, "memo")),
	)
}
