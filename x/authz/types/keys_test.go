package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	sdk "github.com/line/lbm-sdk/types"
)

var granter = sdk.BytesToAccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
var grantee = sdk.BytesToAccAddress(ed25519.GenPrivKey().PubKey().Address().Bytes())
var msgType = SendAuthorization{}.MethodName()

func TestGrantkey(t *testing.T) {
	granter1, grantee1 := ExtractAddressesFromGrantKey(GetAuthorizationStoreKey(grantee, granter, msgType))
	require.Equal(t, granter, granter1)
	require.Equal(t, grantee, grantee1)
}
