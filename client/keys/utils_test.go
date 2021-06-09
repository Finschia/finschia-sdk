package keys

import (
	"bytes"
	"testing"

	"github.com/line/lfb-sdk/client"
	"github.com/line/lfb-sdk/crypto/hd"
	"github.com/line/lfb-sdk/crypto/keyring"
	"github.com/line/lfb-sdk/testutil"
	"github.com/line/lfb-sdk/types"
	"github.com/stretchr/testify/require"
)

func Test_printInfos(t *testing.T) {
	cmd := ListKeysCmd()
	cmd.Flags().AddFlagSet(Commands("home").PersistentFlags())

	kbHome := t.TempDir()

	mockIn := testutil.ApplyMockIODiscardOutErr(cmd)
	kb, err := keyring.New(types.KeyringServiceName(), keyring.BackendTest, kbHome, mockIn)
	require.NoError(t, err)

	kb.NewAccount("something", testutil.TestMnemonic, "", "", hd.Secp256k1)

	clientCtx := client.Context{}.WithKeyring(kb)
	require.NoError(t, err)

	infos, err := clientCtx.Keyring.List()
	require.NoError(t, err)
	buf := bytes.NewBufferString("")
	printInfos(buf, infos, OutputFormatJSON)
	require.Equal(t, buf.String(), "[{\"name\":\"something\",\"type\":\"local\",\"address\":\"link1jyyxx9phqw6tarnxanhyx7ecr992d6yrztj4d0\",\"pubkey\":\"linkpub1cqmsrdepqg8qdahungt0uladr5g4v6kh56x4wx7auws5rfxunyrapnqn7kfkyrszyvh\"}]")
}
