package keys

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lfb-sdk/client"
	"github.com/line/lfb-sdk/crypto/keyring"
	"github.com/line/lfb-sdk/testutil"
	sdk "github.com/line/lfb-sdk/types"
)

func Test_runImportCmd(t *testing.T) {
	cmd := ImportKeyCommand()
	cmd.Flags().AddFlagSet(Commands("home").PersistentFlags())
	mockIn := testutil.ApplyMockIODiscardOutErr(cmd)

	// Now add a temporary keybase
	kbHome := t.TempDir()
	kb, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, kbHome, mockIn)

	clientCtx := client.Context{}.WithKeyring(kb)
	ctx := context.WithValue(context.Background(), client.ClientContextKey, &clientCtx)

	require.NoError(t, err)
	t.Cleanup(func() {
		kb.Delete("keyname1") // nolint:errcheck
	})

	keyfile := filepath.Join(kbHome, "key.asc")
	armoredKey := `-----BEGIN OSTRACON PRIVATE KEY-----
kdf: bcrypt
salt: A278E4F0DC466EF58CC9FC3149688593
type: secp256k1

mSB5rEfN4VCi1EnEca5PigV/WphYtrBFft+QyZ2ISztMeQmuhFNFWwjLBsJm5zXv
KNXMn0ZEeCZtbyNzPPdQUQBwcbueq9vx5NDqQCg=
=wp9z
-----END OSTRACON PRIVATE KEY-----
`
	require.NoError(t, ioutil.WriteFile(keyfile, []byte(armoredKey), 0644))

	mockIn.Reset("123456789\n")
	cmd.SetArgs([]string{"keyname1", keyfile})
	require.NoError(t, cmd.ExecuteContext(ctx))
}
