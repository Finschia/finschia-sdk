package keys

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-rdk/client"
	"github.com/Finschia/finschia-rdk/crypto/hd"
	"github.com/Finschia/finschia-rdk/crypto/keyring"
	"github.com/Finschia/finschia-rdk/testutil"
	"github.com/Finschia/finschia-rdk/testutil/testdata"
	"github.com/Finschia/finschia-rdk/types"
)

func Test_printInfos(t *testing.T) {
	cmd := ListKeysCmd()
	cmd.Flags().AddFlagSet(Commands("home").PersistentFlags())

	kbHome := t.TempDir()

	mockIn := testutil.ApplyMockIODiscardOutErr(cmd)
	kb, err := keyring.New(types.KeyringServiceName(), keyring.BackendTest, kbHome, mockIn)
	require.NoError(t, err)

	kb.NewAccount("something", testdata.TestMnemonic, "", "", hd.Secp256k1)

	clientCtx := client.Context{}.WithKeyring(kb)
	require.NoError(t, err)

	infos, err := clientCtx.Keyring.List()
	require.NoError(t, err)
	buf := bytes.NewBufferString("")
	printInfos(buf, infos, OutputFormatJSON)
	require.Equal(t, buf.String(), "[{\"name\":\"something\",\"type\":\"local\",\"address\":\"link1jyyxx9phqw6tarnxanhyx7ecr992d6yrztj4d0\",\"pubkey\":\"{\\\"@type\\\":\\\"/cosmos.crypto.secp256k1.PubKey\\\",\\\"key\\\":\\\"Ag4G9vyaFv5/rR0RVmrXpo1XG93joUGk3JkH0MwT9ZNi\\\"}\"}]")
}
