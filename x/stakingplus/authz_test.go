package stakingplus

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/crypto/keys/secp256k1"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
)

func TestAminoJson(t *testing.T) {
	var authority = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	var grantee = sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	var valAddr = sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())

	src := &CreateValidatorAuthorization{ValidatorAddress: valAddr.String()}
	expected := fmt.Sprintf("{\"type\":\"lbm-sdk/MsgGrant\",\"value\":{\"authority\":\"%s\",\"authorization\":{\"type\":\"lbm-sdk/CreateValidatorAuthorization\",\"value\":{\"validator_address\":\"%s\"}},\"grantee\":\"%s\"}}", authority.String(), valAddr.String(), grantee.String())

	grantMsg := &foundation.MsgGrant{
		Authority: authority.String(),
		Grantee:   grantee.String(),
	}
	err := grantMsg.SetAuthorization(src)
	require.NoError(t, err)

	err = src.ValidateBasic()
	require.NoError(t, err)

	require.Equal(t, expected, string(grantMsg.GetSignBytes()))
}
