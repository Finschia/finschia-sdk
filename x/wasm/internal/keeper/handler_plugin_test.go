package keeper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	"github.com/line/link-modules/x/coin"
	"github.com/line/link-modules/x/collection"
	"github.com/line/link-modules/x/token"
	"github.com/line/link-modules/x/wasm/internal/types"
)

type testData struct {
	tokenKeeper      token.Keeper
	collectionKeeper collection.Keeper
}

// returns a cleanup function, which must be defered on
func setupTest(t *testing.T) (testData, func()) {
	tempDir, err := ioutil.TempDir("", "wasm")
	require.NoError(t, err)

	_, keepers := CreateTestInput(t, false, tempDir, "staking,link", nil, nil)
	tokenKeeper, collectionKeeper := keepers.TokenKeeper, keepers.CollectionKeeper
	cleanup := func() { os.RemoveAll(tempDir) }
	data := testData{
		tokenKeeper:      tokenKeeper,
		collectionKeeper: collectionKeeper,
	}

	return data, cleanup
}

func TestEncoding(t *testing.T) {
	data, _ := setupTest(t)
	_, _, addr1 := keyPubAddr()
	_, _, addr2 := keyPubAddr()
	invalidAddr := "xrnd1d02kd90n38qvr3qb9qof83fn2d2"
	valAddr := make(sdk.ValAddress, sdk.AddrLen)
	valAddr[0] = 12
	valAddr2 := make(sdk.ValAddress, sdk.AddrLen)
	valAddr2[1] = 123

	jsonMsg := json.RawMessage(`{"foo": 123}`)

	testContractID := "test_contract_id"
	issue := fmt.Sprintf(`{"module":"tokenencode", "msg_data":{"route":"issue", "data":{"issue":{"owner":"%s","to":"%s","name":"TestToken1","symbol":"TT1","img_uri":"","meta":"","amount":"1000","mintable":true,"decimals":"18"}}}}`, addr1.String(), addr2.String())
	issueMsg := json.RawMessage(issue)
	transfer := fmt.Sprintf(`{"module":"tokenencode", "msg_data":{"route":"transfer", "data":{"transfer":{"from":"%s", "contract_id":"%s", "to":"%s", "amount":"100"}}}}`, addr1.String(), testContractID, addr2.String())
	transferMsg := json.RawMessage(transfer)
	mint := fmt.Sprintf(`{"module":"tokenencode", "msg_data":{"route":"mint", "data":{"mint":{"from":"%s", "contract_id":"%s", "to":"%s", "amount":"100"}}}}`, addr1.String(), testContractID, addr2.String())
	mintMsg := json.RawMessage(mint)
	burn := fmt.Sprintf(`{"module":"tokenencode", "msg_data":{"route":"burn", "data":{"burn":{"from":"%s", "contract_id":"%s", "amount":"5"}}}}`, addr1.String(), testContractID)
	burnMsg := json.RawMessage(burn)
	grantPermission := fmt.Sprintf(`{"module":"tokenencode", "msg_data":{"route":"grant_perm", "data":{"grant_perm":{"from":"%s", "contract_id":"%s", "to":"%s", "permission":"mint"}}}}`, addr1.String(), testContractID, addr2.String())
	grantPermissionMsg := json.RawMessage(grantPermission)
	revokePermission := fmt.Sprintf(`{"module":"tokenencode", "msg_data":{"route":"revoke_perm", "data":{"revoke_perm":{"from":"%s", "contract_id":"%s", "permission":"mint"}}}}`, addr1.String(), testContractID)
	revokePermissionMsg := json.RawMessage(revokePermission)
	modify := fmt.Sprintf(`{"module":"tokenencode", "msg_data":{"route":"modify","data":{"modify":{"owner":"%s","contract_id":"%s","changes":[{"field":"meta","value":"update_meta"}]}}}}`, addr1.String(), testContractID)
	modifyMsg := json.RawMessage(modify)

	changes := token.NewChanges(token.NewChange("meta", "update_meta"))

	cases := map[string]struct {
		sender sdk.AccAddress
		input  wasmTypes.CosmosMsg
		// set if valid
		output []sdk.Msg
		// set if invalid
		isError bool
	}{
		"simple send": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Bank: &wasmTypes.BankMsg{
					Send: &wasmTypes.SendMsg{
						FromAddress: addr1.String(),
						ToAddress:   addr2.String(),
						Amount: []wasmTypes.Coin{
							{
								Denom:  "uatom",
								Amount: "12345",
							},
							{
								Denom:  "usdt",
								Amount: "54321",
							},
						},
					},
				},
			},
			output: []sdk.Msg{
				coin.MsgSend{
					From: addr1,
					To:   addr2,
					Amount: sdk.Coins{
						sdk.NewInt64Coin("uatom", 12345),
						sdk.NewInt64Coin("usdt", 54321),
					},
				},
			},
		},
		"invalid send amount": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Bank: &wasmTypes.BankMsg{
					Send: &wasmTypes.SendMsg{
						FromAddress: addr1.String(),
						ToAddress:   addr2.String(),
						Amount: []wasmTypes.Coin{
							{
								Denom:  "uatom",
								Amount: "123.456",
							},
						},
					},
				},
			},
			isError: true,
		},
		"invalid address": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Bank: &wasmTypes.BankMsg{
					Send: &wasmTypes.SendMsg{
						FromAddress: addr1.String(),
						ToAddress:   invalidAddr,
						Amount: []wasmTypes.Coin{
							{
								Denom:  "uatom",
								Amount: "7890",
							},
						},
					},
				},
			},
			isError: true,
		},
		"wasm execute": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Wasm: &wasmTypes.WasmMsg{
					Execute: &wasmTypes.ExecuteMsg{
						ContractAddr: addr2.String(),
						Msg:          jsonMsg,
						Send: []wasmTypes.Coin{
							wasmTypes.NewCoin(12, "eth"),
						},
					},
				},
			},
			output: []sdk.Msg{
				types.MsgExecuteContract{
					Sender:    addr1,
					Contract:  addr2,
					Msg:       jsonMsg,
					SentFunds: sdk.NewCoins(sdk.NewInt64Coin("eth", 12)),
				},
			},
		},
		"wasm instantiate": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Wasm: &wasmTypes.WasmMsg{
					Instantiate: &wasmTypes.InstantiateMsg{
						CodeID: 7,
						Msg:    jsonMsg,
						Send: []wasmTypes.Coin{
							wasmTypes.NewCoin(123, "eth"),
						},
					},
				},
			},
			output: []sdk.Msg{
				types.MsgInstantiateContract{
					Sender: addr1,
					CodeID: 7,
					// TODO: fix this
					Label:     fmt.Sprintf("Auto-created by %s", addr1),
					InitMsg:   jsonMsg,
					InitFunds: sdk.NewCoins(sdk.NewInt64Coin("eth", 123)),
				},
			},
		},
		"staking delegate": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Staking: &wasmTypes.StakingMsg{
					Delegate: &wasmTypes.DelegateMsg{
						Validator: valAddr.String(),
						Amount:    wasmTypes.NewCoin(777, "stake"),
					},
				},
			},
			output: []sdk.Msg{
				staking.MsgDelegate{
					DelegatorAddress: addr1,
					ValidatorAddress: valAddr,
					Amount:           sdk.NewInt64Coin("stake", 777),
				},
			},
		},
		"staking delegate to non-validator": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Staking: &wasmTypes.StakingMsg{
					Delegate: &wasmTypes.DelegateMsg{
						Validator: addr2.String(),
						Amount:    wasmTypes.NewCoin(777, "stake"),
					},
				},
			},
			isError: true,
		},
		"staking undelegate": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Staking: &wasmTypes.StakingMsg{
					Undelegate: &wasmTypes.UndelegateMsg{
						Validator: valAddr.String(),
						Amount:    wasmTypes.NewCoin(555, "stake"),
					},
				},
			},
			output: []sdk.Msg{
				staking.MsgUndelegate{
					DelegatorAddress: addr1,
					ValidatorAddress: valAddr,
					Amount:           sdk.NewInt64Coin("stake", 555),
				},
			},
		},
		"staking redelegate": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Staking: &wasmTypes.StakingMsg{
					Redelegate: &wasmTypes.RedelegateMsg{
						SrcValidator: valAddr.String(),
						DstValidator: valAddr2.String(),
						Amount:       wasmTypes.NewCoin(222, "stake"),
					},
				},
			},
			output: []sdk.Msg{
				staking.MsgBeginRedelegate{
					DelegatorAddress:    addr1,
					ValidatorSrcAddress: valAddr,
					ValidatorDstAddress: valAddr2,
					Amount:              sdk.NewInt64Coin("stake", 222),
				},
			},
		},
		"staking withdraw (implicit recipient)": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Staking: &wasmTypes.StakingMsg{
					Withdraw: &wasmTypes.WithdrawMsg{
						Validator: valAddr2.String(),
					},
				},
			},
			output: []sdk.Msg{
				distribution.MsgSetWithdrawAddress{
					DelegatorAddress: addr1,
					WithdrawAddress:  addr1,
				},
				distribution.MsgWithdrawDelegatorReward{
					DelegatorAddress: addr1,
					ValidatorAddress: valAddr2,
				},
			},
		},
		"staking withdraw (explicit recipient)": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Staking: &wasmTypes.StakingMsg{
					Withdraw: &wasmTypes.WithdrawMsg{
						Validator: valAddr2.String(),
						Recipient: addr2.String(),
					},
				},
			},
			output: []sdk.Msg{
				distribution.MsgSetWithdrawAddress{
					DelegatorAddress: addr1,
					WithdrawAddress:  addr2,
				},
				distribution.MsgWithdrawDelegatorReward{
					DelegatorAddress: addr1,
					ValidatorAddress: valAddr2,
				},
			},
		},
		"issue token": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Custom: issueMsg,
			},
			output: []sdk.Msg{
				token.MsgIssue{
					Owner:    addr1,
					To:       addr2,
					Name:     "TestToken1",
					Symbol:   "TT1",
					ImageURI: "",
					Meta:     "",
					Amount:   sdk.NewInt(1000),
					Mintable: true,
					Decimals: sdk.NewInt(18),
				},
			},
		},
		"transfer token": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Custom: transferMsg,
			},
			output: []sdk.Msg{
				token.MsgTransfer{
					From:       addr1,
					ContractID: testContractID,
					To:         addr2,
					Amount:     sdk.NewInt(100),
				},
			},
		},
		"mint token": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Custom: mintMsg,
			},
			output: []sdk.Msg{
				token.MsgMint{
					From:       addr1,
					ContractID: testContractID,
					To:         addr2,
					Amount:     sdk.NewInt(100),
				},
			},
		},
		"burn token": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Custom: burnMsg,
			},
			output: []sdk.Msg{
				token.NewMsgBurn(addr1, testContractID, sdk.NewInt(5)),
			},
		},
		"grant permission": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Custom: grantPermissionMsg,
			},
			output: []sdk.Msg{
				token.NewMsgGrantPermission(addr1, testContractID, addr2, token.Permission("mint")),
			},
		},
		"revoke permission": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Custom: revokePermissionMsg,
			},
			output: []sdk.Msg{
				token.NewMsgRevokePermission(addr1, testContractID, token.Permission("mint")),
			},
		},
		"modify token": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Custom: modifyMsg,
			},
			output: []sdk.Msg{
				token.NewMsgModify(addr1, testContractID, changes),
			},
		},
		"invalid custom msg": {
			sender: addr1,
			input: wasmTypes.CosmosMsg{
				Custom: json.RawMessage("invalid msg"),
			},
			isError: true,
		},
	}

	e := DefaultEncoders()
	tokenEncodeHandler := token.NewMsgEncodeHandler(data.tokenKeeper)
	var encodeRouter = types.NewRouter()
	encodeRouter.AddRoute(token.EncodeRouterKey, tokenEncodeHandler)

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			res, err := e.Encode(tc.sender, tc.input, encodeRouter)
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.output, res)
			}
		})
	}
}
