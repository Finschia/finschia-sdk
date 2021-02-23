package handler

import (
	"testing"

	"github.com/line/lbm-sdk/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
)

func prepareCreateCollection(t *testing.T) (sdk.Context, sdk.Handler, string) {
	ctx, h := cacheKeeper()
	var contractID string
	msg := types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI)
	res, err := h(ctx, msg)
	require.NoError(t, err)

	contractID = GetMadeContractID(res.Events)

	return ctx, h, contractID
}

func prepareFT(t *testing.T) (sdk.Context, sdk.Handler, string) {
	ctx, h, contractID := prepareCreateCollection(t)

	msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, defaultMeta, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
	_, err := h(ctx, msg)
	require.NoError(t, err)

	return ctx, h, contractID
}

func prepareNFT(t *testing.T, mintTo sdk.AccAddress) (sdk.Context, sdk.Handler, string) {
	ctx, h, contractID := prepareCreateCollection(t)

	msg := types.NewMsgIssueNFT(addr1, contractID, defaultName, defaultMeta)
	_, err := h(ctx, msg)
	require.NoError(t, err)

	param := types.NewMintNFTParam("sword1", defaultMeta, "10000001")
	msg2 := types.NewMsgMintNFT(addr1, contractID, mintTo, param)
	_, err = h(ctx, msg2)
	require.NoError(t, err)

	param = types.NewMintNFTParam("sword2", defaultMeta, "10000001")
	types.NewMsgMintNFT(addr1, contractID, mintTo, param)
	_, err = h(ctx, msg2)
	require.NoError(t, err)

	return ctx, h, contractID
}

func TestHandleMsgIssueFT(t *testing.T) {
	ctx, h, contractID := prepareCreateCollection(t)

	msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, defaultMeta, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
	res, err := h(ctx, msg)
	require.NoError(t, err)

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
	ctx, h, contractID := prepareCreateCollection(t)

	msg := types.NewMsgIssueNFT(addr1, contractID, defaultName, defaultMeta)
	res, err := h(ctx, msg)
	require.NoError(t, err)

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
		msg := types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI)
		res, err := h(ctx, msg)
		contractID = GetMadeContractID(res.Events)
		require.NoError(t, err)
	}

	{
		msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, defaultMeta, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}
	{
		msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, defaultMeta, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}
	{
		msg := types.NewMsgIssueFT(addr2, addr2, contractID, defaultName, defaultMeta, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		_, err := h(ctx, msg)
		require.Error(t, err)
	}

	permission := types.NewIssuePermission()

	{
		msg := types.NewMsgGrantPermission(addr1, contractID, addr2, permission)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}
	{
		msg := types.NewMsgIssueFT(addr2, addr2, contractID, defaultName, defaultMeta, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}
	{
		msg := types.NewMsgRevokePermission(addr1, contractID, permission)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}
	{
		msg := types.NewMsgIssueFT(addr1, addr2, contractID, defaultName, defaultMeta, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		_, err := h(ctx, msg)
		require.Error(t, err)
	}
}

func TestHandlerIssueNFT(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		msg := types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI)
		res, err := h(ctx, msg)
		require.NoError(t, err)
		contractID = GetMadeContractID(res.Events)
	}

	{
		// Expect token type is 1001
		{
			msg := types.NewMsgIssueNFT(addr1, contractID, defaultName, defaultMeta)
			_, err := h(ctx, msg)
			require.NoError(t, err)
		}
		// Expect token type is 1002
		{
			msg := types.NewMsgIssueNFT(addr1, contractID, defaultName, defaultMeta)
			_, err := h(ctx, msg)
			require.NoError(t, err)
		}
		{
			param := types.NewMintNFTParam(defaultName, defaultMeta, defaultTokenType2)
			msg := types.NewMsgMintNFT(addr1, contractID, addr1, param)
			_, err := h(ctx, msg)
			require.NoError(t, err)
		}
		{
			param := types.NewMintNFTParam(defaultName, defaultMeta, defaultTokenType2)
			msg := types.NewMsgMintNFT(addr1, contractID, addr2, param)
			_, err := h(ctx, msg)
			require.NoError(t, err)
		}
		{
			mintPermission := types.NewMintPermission()
			{
				msg := types.NewMsgGrantPermission(addr1, contractID, addr2, mintPermission)
				_, err := h(ctx, msg)
				require.NoError(t, err)
			}
			{
				param := types.NewMintNFTParam(defaultName, defaultMeta, defaultTokenType2)
				msg := types.NewMsgMintNFT(addr2, contractID, addr2, param)
				_, err := h(ctx, msg)
				require.NoError(t, err)
			}
			{
				msg := types.NewMsgRevokePermission(addr1, contractID, mintPermission)
				_, err := h(ctx, msg)
				require.NoError(t, err)
			}
			{
				param := types.NewMintNFTParam(defaultName, defaultMeta, defaultTokenType2)
				msg := types.NewMsgMintNFT(addr1, contractID, addr1, param)
				_, err := h(ctx, msg)
				require.Error(t, err)
			}
		}
	}

	permission := types.NewIssuePermission()

	{
		msg := types.NewMsgGrantPermission(addr1, contractID, addr2, permission)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}

	// Expect token type is 1003
	{
		msg := types.NewMsgIssueNFT(addr2, contractID, defaultName, defaultMeta)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}
	{
		param := types.NewMintNFTParam(defaultName, defaultMeta, defaultTokenType3)
		msg := types.NewMsgMintNFT(addr2, contractID, addr2, param)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}
	{
		msg := types.NewMsgRevokePermission(addr1, contractID, permission)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}
	{
		msg := types.NewMsgIssueNFT(addr1, contractID, defaultName, defaultMeta)
		_, err := h(ctx, msg)
		require.Error(t, err)
	}
}

func TestEvents(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		msg := types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI)
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)
		contractID = GetMadeContractID(res.Events)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("create_collection", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("create_collection", sdk.NewAttribute("name", defaultName)),
			sdk.NewEvent("create_collection", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", "issue")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", "mint")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", "burn")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", "modify")),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		msg := types.NewMsgIssueFT(addr1, addr1, contractID, defaultName, defaultMeta, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)

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
		res, err := h(ctx, msg)
		require.NoError(t, err)
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
		res, err := h(ctx, msg)
		require.NoError(t, err)
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
		msg := types.NewMsgIssueNFT(addr1, contractID, defaultName, defaultMeta)
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("issue_nft", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("issue_nft", sdk.NewAttribute("token_type", defaultTokenType)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		param := types.NewMintNFTParam(defaultName, defaultMeta, defaultTokenType)
		msg := types.NewMsgMintNFT(addr1, contractID, addr1, param)
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)

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

	permission := types.NewIssuePermission()

	{
		msg := types.NewMsgGrantPermission(addr1, contractID, addr2, permission)
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", permission.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgRevokePermission(addr1, contractID, permission)
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "collection")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm", permission.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}
