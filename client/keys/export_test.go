package keys

import (
	"bufio"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/crypto/hd"
	"github.com/line/lbm-sdk/testutil"

	"github.com/line/lbm-sdk/client/flags"
	"github.com/line/lbm-sdk/crypto/keyring"
	sdk "github.com/line/lbm-sdk/types"
)

func Test_runExportCmd(t *testing.T) {
	testCases := []struct {
		name           string
		keyringBackend string
		extraArgs      []string
		userInput      string
		mustFail       bool
		expectedOutput string
	}{
		{
			name:           "--unsafe only must fail",
			keyringBackend: keyring.BackendTest,
			extraArgs:      []string{"--unsafe"},
			mustFail:       true,
		},
		{
			name:           "--unarmored-hex must fail",
			keyringBackend: keyring.BackendTest,
			extraArgs:      []string{"--unarmored-hex"},
			mustFail:       true,
		},
		{
			name:           "--unsafe --unarmored-hex fail with no user confirmation",
			keyringBackend: keyring.BackendTest,
			extraArgs:      []string{"--unsafe", "--unarmored-hex"},
			userInput:      "",
			mustFail:       true,
			expectedOutput: "",
		},
		{
			name:           "--unsafe --unarmored-hex succeed",
			keyringBackend: keyring.BackendTest,
			extraArgs:      []string{"--unsafe", "--unarmored-hex"},
			userInput:      "y\n",
			mustFail:       false,
			expectedOutput: "d4bd5d54ee1b75abc6f5bab08e2e9d3a4b6dfbe6b50e2d6cf2426f3215633a1f\n",
		},
		{
			name:           "file keyring backend properly read password and user confirmation",
			keyringBackend: keyring.BackendFile,
			extraArgs:      []string{"--unsafe", "--unarmored-hex"},
			// first 2 pass for creating the key, then unsafe export confirmation, then unlock keyring pass
			userInput:      "12345678\n12345678\ny\n12345678\n",
			mustFail:       false,
			expectedOutput: "d4bd5d54ee1b75abc6f5bab08e2e9d3a4b6dfbe6b50e2d6cf2426f3215633a1f\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			kbHome := t.TempDir()
			defaultArgs := []string{
				"keyname1",
				fmt.Sprintf("--%s=%s", flags.FlagHome, kbHome),
				fmt.Sprintf("--%s=%s", flags.FlagKeyringBackend, tc.keyringBackend),
			}

			cmd := ExportKeyCommand()
			cmd.Flags().AddFlagSet(Commands("home").PersistentFlags())

			cmd.SetArgs(append(defaultArgs, tc.extraArgs...))
			mockIn, mockOut := testutil.ApplyMockIO(cmd)

			mockIn.Reset(tc.userInput)
			mockInBuf := bufio.NewReader(mockIn)

			// create a key
			kb, err := keyring.New(sdk.KeyringServiceName(), tc.keyringBackend, kbHome, bufio.NewReader(mockInBuf))
			require.NoError(t, err)
			t.Cleanup(func() {
				kb.Delete("keyname1") // nolint:errcheck
			})

			path := sdk.GetConfig().GetFullFundraiserPath()
			_, err = kb.NewAccount("keyname1", testutil.TestMnemonic, "", path, hd.Secp256k1)
			require.NoError(t, err)

			clientCtx := client.Context{}.
				WithKeyringDir(kbHome).
				WithKeyring(kb).
				WithInput(mockInBuf)
			ctx := context.WithValue(context.Background(), client.ClientContextKey, &clientCtx)

			err = cmd.ExecuteContext(ctx)
			if tc.mustFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedOutput, mockOut.String())
			}
		})
	}
}
