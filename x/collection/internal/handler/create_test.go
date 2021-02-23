package handler

import (
	"testing"

	"github.com/line/lbm-sdk/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
)

const (
	contractID = "9be17165"
)

func TestHandleMsgCreateCollection(t *testing.T) {
	ctx, h := cacheKeeper()
	{
		msg := types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI)
		res, err := h(ctx, msg)
		require.NoError(t, err)

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
}
