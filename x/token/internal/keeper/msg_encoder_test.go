package keeper

// import (
// 	"encoding/json"
// 	"fmt"
// 	"testing"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/line/lbm-sdk/v2/x/token/internal/types"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// func Test_Encode(t *testing.T) {
// 	encodeHandler := NewMsgEncodeHandler(keeper)
// 	jsonMsg := json.RawMessage(`{"foo": 123}`)

// 	testContractID := "test_contract_id"
// 	issue := fmt.Sprintf(`{"route":"issue", "data":{"owner":"%s","to":"%s","name":"TestToken1","symbol":"TT1","img_uri":"","meta":"","amount":"1000","mintable":true,"decimals":"18"}}`, addr1.String(), addr2.String())
// 	issueMsg := json.RawMessage(issue)
// 	transfer := fmt.Sprintf(`{"route":"transfer", "data":{"from":"%s", "contract_id":"%s", "to":"%s", "amount":"100"}}`, addr1.String(), testContractID, addr2.String())
// 	transferMsg := json.RawMessage(transfer)
// 	transferFrom := fmt.Sprintf(`{"route":"transfer_from", "data":{"proxy":"%s", "from":"%s", "contract_id":"%s", "to":"%s", "amount":"100"}}`, addr3.String(), addr1.String(), testContractID, addr2.String())
// 	transferFromMsg := json.RawMessage(transferFrom)
// 	mint := fmt.Sprintf(`{"route":"mint", "data":{"from":"%s", "contract_id":"%s", "to":"%s", "amount":"100"}}`, addr1.String(), testContractID, addr2.String())
// 	mintMsg := json.RawMessage(mint)
// 	burn := fmt.Sprintf(`{"route":"burn", "data":{"from":"%s", "contract_id":"%s", "amount":"5"}}`, addr1.String(), testContractID)
// 	burnMsg := json.RawMessage(burn)
// 	burnFrom := fmt.Sprintf(`{"route":"burn_from", "data":{"proxy":"%s", "from":"%s", "contract_id":"%s", "amount":"5"}}`, addr2.String(), addr1.String(), testContractID)
// 	burnFromMsg := json.RawMessage(burnFrom)

// 	grantPermission := fmt.Sprintf(`{"route":"grant_perm", "data":{"from":"%s", "contract_id":"%s", "to":"%s", "permission":"mint"}}`, addr1.String(), testContractID, addr2.String())
// 	grantPermissionMsg := json.RawMessage(grantPermission)
// 	revokePermission := fmt.Sprintf(`{"route":"revoke_perm", "data":{"from":"%s", "contract_id":"%s", "permission":"mint"}}`, addr1.String(), testContractID)
// 	revokePermissionMsg := json.RawMessage(revokePermission)
// 	modify := fmt.Sprintf(`{"route":"modify","data":{"owner":"%s","contract_id":"%s","changes":[{"field":"meta","value":"update_meta"}]}}`, addr1.String(), testContractID)
// 	modifyMsg := json.RawMessage(modify)
// 	approver := fmt.Sprintf(`{"route":"approve", "data":{"approver":"%s", "contract_id":"%s", "proxy":"%s"}}`, addr1.String(), testContractID, addr2.String())
// 	approverMsg := json.RawMessage(approver)

// 	changes := types.NewChanges(types.NewChange("meta", "update_meta"))

// 	cases := map[string]struct {
// 		input json.RawMessage
// 		// set if valid
// 		output []sdk.Msg
// 		// set if invalid
// 		isError bool
// 	}{
// 		"issue token": {
// 			input: issueMsg,
// 			output: []sdk.Msg{
// 				types.MsgIssue{
// 					Owner:    addr1,
// 					To:       addr2,
// 					Name:     "TestToken1",
// 					Symbol:   "TT1",
// 					ImageURI: "",
// 					Meta:     "",
// 					Amount:   sdk.NewInt(1000),
// 					Mintable: true,
// 					Decimals: sdk.NewInt(18),
// 				},
// 			},
// 		},
// 		"transfer token": {
// 			input: transferMsg,
// 			output: []sdk.Msg{
// 				types.MsgTransfer{
// 					From:       addr1,
// 					ContractID: testContractID,
// 					To:         addr2,
// 					Amount:     sdk.NewInt(100),
// 				},
// 			},
// 		},
// 		"transfer from token": {
// 			input: transferFromMsg,
// 			output: []sdk.Msg{
// 				types.MsgTransferFrom{
// 					Proxy:      addr3,
// 					From:       addr1,
// 					ContractID: testContractID,
// 					To:         addr2,
// 					Amount:     sdk.NewInt(100),
// 				},
// 			},
// 		},
// 		"mint token": {
// 			input: mintMsg,
// 			output: []sdk.Msg{
// 				types.MsgMint{
// 					From:       addr1,
// 					ContractID: testContractID,
// 					To:         addr2,
// 					Amount:     sdk.NewInt(100),
// 				},
// 			},
// 		},
// 		"burn token": {
// 			input: burnMsg,
// 			output: []sdk.Msg{
// 				types.NewMsgBurn(addr1, testContractID, sdk.NewInt(5)),
// 			},
// 		},
// 		"burn from token": {
// 			input: burnFromMsg,
// 			output: []sdk.Msg{
// 				types.NewMsgBurnFrom(addr2, testContractID, addr1, sdk.NewInt(5)),
// 			},
// 		},
// 		"grant permission": {
// 			input: grantPermissionMsg,
// 			output: []sdk.Msg{
// 				types.NewMsgGrantPermission(addr1, testContractID, addr2, types.Permission("mint")),
// 			},
// 		},
// 		"revoke permission": {
// 			input: revokePermissionMsg,
// 			output: []sdk.Msg{
// 				types.NewMsgRevokePermission(addr1, testContractID, types.Permission("mint")),
// 			},
// 		},
// 		"modify token": {
// 			input: modifyMsg,
// 			output: []sdk.Msg{
// 				types.NewMsgModify(addr1, testContractID, changes),
// 			},
// 		},
// 		"approve": {
// 			input: approverMsg,
// 			output: []sdk.Msg{
// 				types.NewMsgApprove(addr1, testContractID, addr2),
// 			},
// 		},
// 		"unknown custom msg": {
// 			input:   jsonMsg,
// 			isError: true,
// 		},
// 	}

// 	for name, tc := range cases {
// 		tc := tc
// 		t.Run(name, func(t *testing.T) {
// 			res, err := encodeHandler(tc.input)
// 			if tc.isError {
// 				require.Error(t, err)
// 			} else {
// 				require.NoError(t, err)
// 				assert.Equal(t, tc.output, res)
// 			}
// 		})
// 	}
// }
