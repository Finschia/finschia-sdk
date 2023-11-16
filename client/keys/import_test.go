package keys

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/client"
	"github.com/Finschia/finschia-sdk/client/flags"
	"github.com/Finschia/finschia-sdk/crypto/keyring"
	"github.com/Finschia/finschia-sdk/testutil"
	sdk "github.com/Finschia/finschia-sdk/types"
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

	armoredKey := `-----BEGIN TENDERMINT PRIVATE KEY-----
kdf: bcrypt
salt: A53F628182B827E07DD11A96EAB9D526
type: secp256k1

Ax9IQsSq+jOWkPRDJQ69a5/uUm4XliPim/CbYDVoXO6D3fts5IEXcUTmIa60ynC/
8hzYAawzYMO95Kwi0NI8WW9wUv3TseSWFv6/RpU=
=umYd
-----END TENDERMINT PRIVATE KEY-----`

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

			require.NoError(t, os.WriteFile(keyfile, []byte(armoredKey), 0o600))

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
