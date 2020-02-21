package handler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgModifyTokenURI(t *testing.T) {
	ctx, h := cacheKeeper()
	modifyTokenURI := "modifyTokenURI"

	t.Log("token not exist")
	{
		res := h(ctx, types.NewMsgModifyTokenURI(addr1, symbol, modifyTokenURI))
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenNotExist, res.Code)
		verifyEventFunc(t, nil, res.Events)
	}
	t.Log("modify token for FT")
	{
		res := h(ctx, types.NewMsgIssue(addr1, name, symbol, tokenuri, amount, decimals, true))
		require.True(t, res.IsOK())

		// TokenURI too long
		length1001String := strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001
		res = h(ctx, types.NewMsgModifyTokenURI(addr1, symbol, length1001String))
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenInvalidTokenURILength, res.Code)

		// success
		res = h(ctx, types.NewMsgModifyTokenURI(addr1, symbol, modifyTokenURI))
		require.True(t, res.Code.IsOK())
		for _, event := range res.Events {
			if event.Type == types.EventTypeModifyTokenURI {
				for _, attr := range event.Attributes {
					if bytes.Equal(attr.Key, []byte(types.AttributeKeyTokenURI)) {
						require.Equal(t, modifyTokenURI, string(attr.Value))
					}
				}
			}
		}
	}
}
