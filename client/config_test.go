package client

import (
	"context"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// For https://github.com/line/lbm-sdk/issues/3899
func Test_runConfigCmdTwiceWithShorterNodeValue(t *testing.T) {
	// Prepare environment
	t.Parallel()
	configHome, cleanup := tmpDir(t)
	defer cleanup()
	_ = os.RemoveAll(filepath.Join(configHome, "config"))

	clientCtx := Context{}.WithHomeDir(configHome)
	ctx := context.Background()
	ctx = context.WithValue(ctx, ClientContextKey, &clientCtx)

	execConfigCmd(t, ctx, configHome, []string{"node", "tcp://localhost:26657"})
	execConfigCmd(t, ctx, configHome, []string{"node", "--get"})
	execConfigCmd(t, ctx, configHome, []string{"node", "tcp://local:26657"})
	execConfigCmd(t, ctx, configHome, []string{"node", "--get"})
}

func execConfigCmd(t *testing.T, ctx context.Context, configHome string, args []string ) {
	cmd := ConfigCmd(configHome)
	assert.NotNil(t, cmd)
	cmd.SetArgs(args)
	err := viper.BindPFlags(cmd.Flags())
	assert.Nil(t, err)
	err = cmd.ExecuteContext(ctx)
	assert.Nil(t, err)
}

func tmpDir(t *testing.T) (string, func()) {
	dir, err := ioutil.TempDir("", t.Name()+"_")
	require.NoError(t, err)
	return dir, func() { _ = os.RemoveAll(dir) }
}