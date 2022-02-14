package keys

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/crypto/keyring"
	"github.com/line/lbm-sdk/testutil"
	sdk "github.com/line/lbm-sdk/types"
)

func Test_runImportCmd(t *testing.T) {
	testCases := []struct {
		name           string
		keyringBackend string
		userInput      string
		expectError    bool
	}{
		{
			name:           "test backend success",
			keyringBackend: keyring.BackendTest,
			// key armor passphrase
			userInput: "123456789\n",
		},
		{
			name:           "test backend fail with wrong armor pass",
			keyringBackend: keyring.BackendTest,
			userInput:      "987654321\n",
			expectError:    true,
		},
		{
			name:           "file backend success",
			keyringBackend: keyring.BackendFile,
			// key armor passphrase + keyring password x2
			userInput: "123456789\n12345678\n12345678\n",
		},
		{
			name:           "file backend fail with wrong armor pass",
			keyringBackend: keyring.BackendFile,
			userInput:      "987654321\n12345678\n12345678\n",
			expectError:    true,
		},
		{
			name:           "file backend fail with wrong keyring pass",
			keyringBackend: keyring.BackendFile,
			userInput:      "123465789\n12345678\n87654321\n",
			expectError:    true,
		},
		{
			name:           "file backend fail with no keyring pass",
			keyringBackend: keyring.BackendFile,
			userInput:      "123465789\n",
			expectError:    true,
		},
	}

	armoredKey := `-----BEGIN OSTRACON PRIVATE KEY-----
kdf: bcrypt
salt: A278E4F0DC466EF58CC9FC3149688593
type: secp256k1

mSB5rEfN4VCi1EnEca5PigV/WphYtrBFft+QyZ2ISztMeQmuhFNFWwjLBsJm5zXv
KNXMn0ZEeCZtbyNzPPdQUQBwcbueq9vx5NDqQCg=
=wp9z
-----END OSTRACON PRIVATE KEY-----
`

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := ImportKeyCommand()
			cmd.Flags().AddFlagSet(Commands("home").PersistentFlags())
			mockIn := testutil.ApplyMockIODiscardOutErr(cmd)

			// Now add a temporary keybase
			kbHome := t.TempDir()
			kb, err := keyring.New(sdk.KeyringServiceName(), tc.keyringBackend, kbHome, nil)

			clientCtx := client.Context{}.
				WithKeyringDir(kbHome).
				WithKeyring(kb).
				WithInput(mockIn)
			ctx := context.WithValue(context.Background(), client.ClientContextKey, &clientCtx)

			require.NoError(t, err)
			t.Cleanup(func() {
				kb.Delete("keyname1") // nolint:errcheck
			})

			keyfile := filepath.Join(kbHome, "key.asc")

			require.NoError(t, ioutil.WriteFile(keyfile, []byte(armoredKey), 0644))

			defer func() {
				_ = os.RemoveAll(kbHome)
			}()

			mockIn.Reset(tc.userInput)
			cmd.SetArgs([]string{
				"keyname1", keyfile,
				fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, tc.keyringBackend),
			})

			err = cmd.ExecuteContext(ctx)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
