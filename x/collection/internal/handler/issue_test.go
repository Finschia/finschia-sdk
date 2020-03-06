package handler

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func prepareCreateCollection() (sdk.Context, sdk.Handler, string) {
	ctx, h := cacheKeeper()
	var contractID string
	msg := types.NewMsgCreateCollection(addr1, defaultName, defaultImgURI)
	res := h(ctx, msg)
	contractID = GetMadeContractID(res.Events)

	return ctx, h, contractID
}

func prepareFT() (sdk.Context, sdk.Handler, string) {
	ctx, h, contractID := prepareCreateCollection()

	msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
	_ = h(ctx, msg)

	return ctx, h, contractID
}

func prepareNFT(mintTo sdk.AccAddress) (sdk.Context, sdk.Handler, string) {
	ctx, h, contractID := prepareCreateCollection()

	msg := types.NewMsgIssueNFT(addr1, contractID, defaultName)
	_ = h(ctx, msg)

	msg2 := types.NewMsgMintNFT(addr1, contractID, mintTo, "sword1", "10000001")
	_ = h(ctx, msg2)

	types.NewMsgMintNFT(addr1, contractID, mintTo, "sword2", "10000001")
	_ = h(ctx, msg2)

	return ctx, h, contractID
}

func TestHandleMsgIssueFT(t *testing.T) {
	ctx, h, contractID := prepareCreateCollection()

	msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())

	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
		sdk.NewEvent("issue_ft", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("issue_ft", sdk.NewAttribute("name", defaultName)),
		sdk.NewEvent("issue_ft", sdk.NewAttribute("token_id", defaultTokenIDFT)),
		sdk.NewEvent("issue_ft", sdk.NewAttribute("owner", addr1.String())),
		sdk.NewEvent("issue_ft", sdk.NewAttribute("to", addr1.String())),
		sdk.NewEvent("issue_ft", sdk.NewAttribute("amount", "1000")),
		sdk.NewEvent("issue_ft", sdk.NewAttribute("mintable", "true")),
		sdk.NewEvent("issue_ft", sdk.NewAttribute("decimals", "6")),
	}
	verifyEventFunc(t, e, res.Events)
}

func TestHandleMsgIssueNFT(t *testing.T) {
	ctx, h, contractID := prepareCreateCollection()

	msg := types.NewMsgIssueNFT(addr1, contractID, defaultName)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())

	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
		sdk.NewEvent("issue_nft", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("issue_nft", sdk.NewAttribute("token_type", defaultTokenType)),
	}
	verifyEventFunc(t, e, res.Events)
}

func TestHandlerIssueFT(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		msg := types.NewMsgCreateCollection(addr1, defaultName, defaultImgURI)
		res := h(ctx, msg)
		contractID = GetMadeContractID(res.Events)
		require.True(t, res.Code.IsOK())
	}

	{
		msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueFT(addr2, addr2, contractID, defaultName, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenPermission, res.Code)
	}

	permission := types.Permission{
		Action:   "issue",
		Resource: contractID,
	}

	{
		msg := types.NewMsgGrantPermission(addr1, addr2, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueFT(addr2, addr2, contractID, defaultName, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgRevokePermission(addr1, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueFT(addr1, addr2, contractID, defaultName, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenPermission, res.Code)
	}
}

func TestHandlerIssueNFT(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		msg := types.NewMsgCreateCollection(addr1, defaultName, defaultImgURI)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		contractID = GetMadeContractID(res.Events)
	}

	{
		//Expect token type is 1001
		{
			msg := types.NewMsgIssueNFT(addr1, contractID, defaultName)
			res := h(ctx, msg)
			require.True(t, res.Code.IsOK())
		}
		//Expect token type is 1002
		{
			msg := types.NewMsgIssueNFT(addr1, contractID, defaultName)
			res := h(ctx, msg)
			require.True(t, res.Code.IsOK())
		}
		{
			msg := types.NewMsgMintNFT(addr1, contractID, addr1, defaultName, defaultTokenType2)
			res := h(ctx, msg)
			require.True(t, res.Code.IsOK())
		}
		{
			msg := types.NewMsgMintNFT(addr1, contractID, addr2, defaultName, defaultTokenType2)
			res := h(ctx, msg)
			require.True(t, res.Code.IsOK())
		}
		{
			mintPermission := types.Permission{
				Action:   "mint",
				Resource: contractID,
			}
			{
				msg := types.NewMsgGrantPermission(addr1, addr2, mintPermission)
				res := h(ctx, msg)
				require.True(t, res.Code.IsOK())
			}
			{
				msg := types.NewMsgMintNFT(addr2, contractID, addr2, defaultName, defaultTokenType2)
				res := h(ctx, msg)
				require.True(t, res.Code.IsOK())
			}
			{
				msg := types.NewMsgRevokePermission(addr1, mintPermission)
				res := h(ctx, msg)
				require.True(t, res.Code.IsOK())
			}
			{
				msg := types.NewMsgMintNFT(addr1, contractID, addr1, defaultName, defaultTokenType2)
				res := h(ctx, msg)
				require.False(t, res.Code.IsOK())
				require.Equal(t, types.DefaultCodespace, res.Codespace)
				require.Equal(t, types.CodeTokenPermission, res.Code)
			}
		}
	}

	permission := types.Permission{
		Action:   "issue",
		Resource: contractID,
	}

	{
		msg := types.NewMsgGrantPermission(addr1, addr2, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	//Expect token type is 1003
	{
		msg := types.NewMsgIssueNFT(addr2, contractID, defaultName)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgMintNFT(addr2, contractID, addr2, defaultName, defaultTokenType3)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgRevokePermission(addr1, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueNFT(addr1, contractID, defaultName)
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenPermission, res.Code)
	}
}

func TestEvents(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		msg := types.NewMsgCreateCollection(addr1, defaultName, defaultImgURI)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		contractID = GetMadeContractID(res.Events)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("create_collection", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("create_collection", sdk.NewAttribute("name", defaultName)),
			sdk.NewEvent("create_collection", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "issue")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "mint")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "burn")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "modify")),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("issue_ft", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("issue_ft", sdk.NewAttribute("name", defaultName)),
			sdk.NewEvent("issue_ft", sdk.NewAttribute("token_id", defaultTokenIDFT)),
			sdk.NewEvent("issue_ft", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue_ft", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("issue_ft", sdk.NewAttribute("amount", sdk.NewInt(defaultAmount).String())),
			sdk.NewEvent("issue_ft", sdk.NewAttribute("mintable", "true")),
			sdk.NewEvent("issue_ft", sdk.NewAttribute("decimals", sdk.NewInt(defaultDecimals).String())),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		msg := types.NewMsgMintFT(addr1, contractID, addr1, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("mint_ft", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("mint_ft", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("mint_ft", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("mint_ft", sdk.NewAttribute("amount", types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)).String())),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		msg := types.NewMsgBurnFT(addr1, contractID, types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn_ft", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("burn_ft", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("burn_ft", sdk.NewAttribute("amount", types.NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)).String())),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		msg := types.NewMsgIssueNFT(addr1, contractID, defaultName)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("issue_nft", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("issue_nft", sdk.NewAttribute("token_type", defaultTokenType)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		msg := types.NewMsgMintNFT(addr1, contractID, addr1, defaultName, defaultTokenType)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("mint_nft", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("mint_nft", sdk.NewAttribute("name", defaultName)),
			sdk.NewEvent("mint_nft", sdk.NewAttribute("token_id", defaultTokenID1)),
			sdk.NewEvent("mint_nft", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("mint_nft", sdk.NewAttribute("to", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}

	permission := types.Permission{
		Action:   "issue",
		Resource: contractID,
	}

	{
		msg := types.NewMsgGrantPermission(addr1, addr2, permission)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", permission.GetResource())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", permission.GetAction())),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgRevokePermission(addr1, permission)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_resource", permission.GetResource())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_action", permission.GetAction())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}
