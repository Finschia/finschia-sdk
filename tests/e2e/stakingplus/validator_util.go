package stakingplus

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func getValidator(t *testing.T, baseDir string, cfg network.Config, i int) sdk.AccAddress {
	buf := bufio.NewReader(os.Stdin)
	nodeDirName := fmt.Sprintf("node%d", i)
	clientDir := filepath.Join(baseDir, nodeDirName, "simcli")

	kb, err := keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, clientDir, buf, cfg.Codec, cfg.KeyringOptions...)
	require.NoError(t, err)

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(cfg.SigningAlgo, keyringAlgos)
	require.NoError(t, err)

	var mnemonic string
	if i < len(cfg.Mnemonics) {
		mnemonic = cfg.Mnemonics[i]
	}

	addr, _, err := testutil.GenerateSaveCoinKey(kb, nodeDirName, mnemonic, true, algo)
	require.NoError(t, err)

	return addr
}
