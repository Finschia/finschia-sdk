package keeper

import (
	"encoding/json"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/collection/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Encode(t *testing.T) {
	encodeHandler := NewMsgEncodeHandler(keeper)

	testContractName := "test_collection"
	testContractID := "9be17165"
	nft1 := "nft-1"
	ft1 := "ft-1"
	amount := int64(100)
	mintNftParams := []types.MintNFTParam{types.NewMintNFTParam(nft1, "", defaultTokenType)}
	coins := []types.Coin{types.NewCoin(defaultTokenIDFT, sdk.NewInt(amount))}
	changes := []types.Change{types.NewChange("f", "t")}

	create := fmt.Sprintf(`{"route":"create","data":{"owner":"%s", "name":"%s","meta":"","base_img_uri":""}}`, addr1.String(), testContractName)
	createMsg := json.RawMessage(create)
	issueNft := fmt.Sprintf(`{"route":"issue_nft","data":{"owner":"%s", "contract_id":"%s", "name":"%s","meta":""}}`, addr1.String(), testContractID, nft1)
	issueNftMsg := json.RawMessage(issueNft)
	issueFt := fmt.Sprintf(`{"route":"issue_ft","data":{"owner":"%s", "contract_id":"%s", "to":"%s", "name":"%s","meta":"", "amount":"%d", "mintable":true, "decimals":"18"}}`, addr1.String(), testContractID, addr2.String(), ft1, amount)
	issueFtMsg := json.RawMessage(issueFt)
	mintNft := fmt.Sprintf(`{"route":"mint_nft","data":{"from":"%s", "contract_id":"%s", "to":"%s", "params":[{"name":"%s", "meta":"", "token_type":"%s"}]}}`, addr1.String(), testContractID, addr2.String(), nft1, defaultTokenType)
	mintNftMsg := json.RawMessage(mintNft)
	mintFt := fmt.Sprintf(`{"route":"mint_ft","data":{"from":"%s", "contract_id":"%s", "to":"%s", "amount":[{"token_id":"%s", "amount":"%d"}]}}`, addr1.String(), testContractID, addr2.String(), defaultTokenIDFT, amount)
	mintFtMsg := json.RawMessage(mintFt)
	burnNft := fmt.Sprintf(`{"route":"burn_nft","data":{"from":"%s", "contract_id":"%s", "token_ids":["%s"]}}`, addr1.String(), testContractID, defaultTokenID1)
	burnNftMsg := json.RawMessage(burnNft)
	burnNftFrom := fmt.Sprintf(`{"route":"burn_nft_from","data":{"proxy":"%s", "from":"%s","contract_id":"%s", "token_ids":["%s"]}}`, addr2.String(), addr1.String(), testContractID, defaultTokenID1)
	burnNftFromMsg := json.RawMessage(burnNftFrom)
	burnFt := fmt.Sprintf(`{"route":"burn_ft","data":{"from":"%s","contract_id":"%s", "amount":[{"token_id":"%s", "amount":"%d"}]}}`, addr1.String(), testContractID, defaultTokenIDFT, amount)
	burnFtMsg := json.RawMessage(burnFt)
	burnFtFrom := fmt.Sprintf(`{"route":"burn_ft_from","data":{"proxy":"%s", "from":"%s","contract_id":"%s", "amount":[{"token_id":"%s", "amount":"%d"}]}}`, addr2.String(), addr1.String(), testContractID, defaultTokenIDFT, amount)
	burnFtFromMsg := json.RawMessage(burnFtFrom)
	transferNft := fmt.Sprintf(`{"route":"transfer_nft","data":{"from":"%s", "contract_id":"%s", "to":"%s", "token_ids":["%s"]}}`, addr1.String(), testContractID, addr2.String(), defaultTokenID1)
	transferNftMsg := json.RawMessage(transferNft)
	transferNftFrom := fmt.Sprintf(`{"route":"transfer_nft_from","data":{"proxy":"%s", "from":"%s", "contract_id":"%s", "to":"%s", "token_ids":["%s"]}}`, addr2.String(), addr1.String(), testContractID, addr2.String(), defaultTokenID1)
	transferNftFromMsg := json.RawMessage(transferNftFrom)
	transferFt := fmt.Sprintf(`{"route":"transfer_ft","data":{"from":"%s", "contract_id":"%s", "to":"%s", "amount":[{"token_id":"%s", "amount":"%d"}]}}`, addr1.String(), testContractID, addr2.String(), defaultTokenIDFT, amount)
	transferFtMsg := json.RawMessage(transferFt)
	transferFtFrom := fmt.Sprintf(`{"route":"transfer_ft_from","data":{"proxy":"%s", "from":"%s", "contract_id":"%s", "to":"%s", "amount":[{"token_id":"%s", "amount":"%d"}]}}`, addr2.String(), addr1.String(), testContractID, addr2.String(), defaultTokenIDFT, amount)
	transferFtFromMsg := json.RawMessage(transferFtFrom)
	approve := fmt.Sprintf(`{"route":"approve","data":{"approver":"%s", "contract_id":"%s", "proxy":"%s"}}`, addr1.String(), testContractID, addr2.String())
	approveMsg := json.RawMessage(approve)
	disapprove := fmt.Sprintf(`{"route":"disapprove","data":{"approver":"%s", "contract_id":"%s", "proxy":"%s"}}`, addr1.String(), testContractID, addr2.String())
	disapproveMsg := json.RawMessage(disapprove)
	attach := fmt.Sprintf(`{"route":"attach","data":{"from":"%s", "contract_id":"%s", "to_token_id":"%s", "token_id":"%s"}}`, addr1.String(), testContractID, defaultTokenID1, defaultTokenID2)
	attachMsg := json.RawMessage(attach)
	detach := fmt.Sprintf(`{"route":"detach","data":{"from":"%s", "contract_id":"%s", "token_id":"%s"}}`, addr1.String(), testContractID, defaultTokenID1)
	detachMsg := json.RawMessage(detach)
	attachFrom := fmt.Sprintf(`{"route":"attach_from","data":{"proxy":"%s", "from":"%s", "contract_id":"%s", "to_token_id":"%s", "token_id":"%s"}}`, addr2.String(), addr1.String(), testContractID, defaultTokenID1, defaultTokenID2)
	attachFromMsg := json.RawMessage(attachFrom)
	detachFrom := fmt.Sprintf(`{"route":"detach_from","data":{"proxy":"%s", "from":"%s", "contract_id":"%s", "token_id":"%s"}}`, addr2.String(), addr1.String(), testContractID, defaultTokenID1)
	detachFromMsg := json.RawMessage(detachFrom)
	modify := fmt.Sprintf(`{"route":"modify","data":{"owner":"%s", "contract_id":"%s", "token_type":"%s", "token_index":"%s", "changes":[{"field":"f", "value":"t"}]}}`, addr1.String(), testContractID, defaultTokenType, defaultTokenIndex)
	modifyMsg := json.RawMessage(modify)

	cases := map[string]struct {
		input json.RawMessage
		// set if valid
		output []sdk.Msg
		// set if invalid
		isError bool
	}{
		"create collection": {
			input: createMsg,
			output: []sdk.Msg{
				types.MsgCreateCollection{
					Owner:      addr1,
					Name:       testContractName,
					Meta:       "",
					BaseImgURI: "",
				},
			},
		},
		"issue nft": {
			input: issueNftMsg,
			output: []sdk.Msg{
				types.MsgIssueNFT{
					Owner:      addr1,
					ContractID: testContractID,
					Name:       nft1,
					Meta:       "",
				},
			},
		},
		"issue ft": {
			input: issueFtMsg,
			output: []sdk.Msg{
				types.MsgIssueFT{
					Owner:      addr1,
					ContractID: testContractID,
					To:         addr2,
					Name:       ft1,
					Meta:       "",
					Amount:     sdk.NewInt(amount),
					Mintable:   true,
					Decimals:   sdk.NewInt(18),
				},
			},
		},
		"mint nft": {
			input: mintNftMsg,
			output: []sdk.Msg{
				types.MsgMintNFT{
					From:          addr1,
					ContractID:    testContractID,
					To:            addr2,
					MintNFTParams: mintNftParams,
				},
			},
		},
		"mint ft": {
			input: mintFtMsg,
			output: []sdk.Msg{
				types.MsgMintFT{
					From:       addr1,
					ContractID: testContractID,
					To:         addr2,
					Amount:     coins,
				},
			},
		},
		"burn nft": {
			input: burnNftMsg,
			output: []sdk.Msg{
				types.MsgBurnNFT{
					From:       addr1,
					ContractID: testContractID,
					TokenIDs:   []string{defaultTokenID1},
				},
			},
		},
		"burn nft from": {
			input: burnNftFromMsg,
			output: []sdk.Msg{
				types.MsgBurnNFTFrom{
					Proxy:      addr2,
					From:       addr1,
					ContractID: testContractID,
					TokenIDs:   []string{defaultTokenID1},
				},
			},
		},
		"burn ft": {
			input: burnFtMsg,
			output: []sdk.Msg{
				types.MsgBurnFT{
					From:       addr1,
					ContractID: testContractID,
					Amount:     coins,
				},
			},
		},
		"burn ft from": {
			input: burnFtFromMsg,
			output: []sdk.Msg{
				types.MsgBurnFTFrom{
					Proxy:      addr2,
					From:       addr1,
					ContractID: testContractID,
					Amount:     coins,
				},
			},
		},
		"transfer nft": {
			input: transferNftMsg,
			output: []sdk.Msg{
				types.MsgTransferNFT{
					From:       addr1,
					ContractID: testContractID,
					To:         addr2,
					TokenIDs:   []string{defaultTokenID1},
				},
			},
		},
		"transfer nft from": {
			input: transferNftFromMsg,
			output: []sdk.Msg{
				types.MsgTransferNFTFrom{
					Proxy:      addr2,
					From:       addr1,
					ContractID: testContractID,
					To:         addr2,
					TokenIDs:   []string{defaultTokenID1},
				},
			},
		},
		"transfer ft": {
			input: transferFtMsg,
			output: []sdk.Msg{
				types.MsgTransferFT{
					From:       addr1,
					ContractID: testContractID,
					To:         addr2,
					Amount:     coins,
				},
			},
		},
		"transfer ft from": {
			input: transferFtFromMsg,
			output: []sdk.Msg{
				types.MsgTransferFTFrom{
					Proxy:      addr2,
					From:       addr1,
					ContractID: testContractID,
					To:         addr2,
					Amount:     coins,
				},
			},
		},
		"approve": {
			input: approveMsg,
			output: []sdk.Msg{
				types.MsgApprove{
					Approver:   addr1,
					ContractID: testContractID,
					Proxy:      addr2,
				},
			},
		},
		"disapprove": {
			input: disapproveMsg,
			output: []sdk.Msg{
				types.MsgDisapprove{
					Approver:   addr1,
					ContractID: testContractID,
					Proxy:      addr2,
				},
			},
		},
		"attach": {
			input: attachMsg,
			output: []sdk.Msg{
				types.MsgAttach{
					From:       addr1,
					ContractID: testContractID,
					ToTokenID:  defaultTokenID1,
					TokenID:    defaultTokenID2,
				},
			},
		},
		"detach": {
			input: detachMsg,
			output: []sdk.Msg{
				types.MsgDetach{
					From:       addr1,
					ContractID: testContractID,
					TokenID:    defaultTokenID1,
				},
			},
		},
		"attach from": {
			input: attachFromMsg,
			output: []sdk.Msg{
				types.MsgAttachFrom{
					Proxy:      addr2,
					From:       addr1,
					ContractID: testContractID,
					ToTokenID:  defaultTokenID1,
					TokenID:    defaultTokenID2,
				},
			},
		},
		"detach from": {
			input: detachFromMsg,
			output: []sdk.Msg{
				types.MsgDetachFrom{
					Proxy:      addr2,
					From:       addr1,
					ContractID: testContractID,
					TokenID:    defaultTokenID1,
				},
			},
		},
		"modify": {
			input: modifyMsg,
			output: []sdk.Msg{
				types.MsgModify{
					Owner:      addr1,
					ContractID: testContractID,
					TokenType:  defaultTokenType,
					TokenIndex: defaultTokenIndex,
					Changes:    changes,
				},
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			res, err := encodeHandler(tc.input)
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.output, res)
			}
		})
	}
}
