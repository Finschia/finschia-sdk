package server_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/line/lbm-sdk/server"
	cfg "github.com/line/ostracon/config"
	"github.com/line/ostracon/crypto"
	tmjson "github.com/line/ostracon/libs/json"
	"github.com/line/ostracon/libs/log"
	tmos "github.com/line/ostracon/libs/os"
	tmrand "github.com/line/ostracon/libs/rand"
	"github.com/line/ostracon/p2p"
	"github.com/line/ostracon/privval"
	"github.com/line/ostracon/types"
	tmtime "github.com/line/ostracon/types/time"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

var (
	logger = log.NewOCLogger(log.NewSyncWriter(os.Stdout))
)

func TestShowValidator(t *testing.T) {
	testCommon := newPrecedenceCommon(t)

	serverCtx := &server.Context{}
	ctx := context.WithValue(context.Background(), server.ServerContextKey, serverCtx)

	if err := testCommon.cmd.ExecuteContext(ctx); err != cancelledInPreRun {
		t.Fatalf("function failed with [%T] %v", err, err)
	}

	// ostracon init & create the server config file
	initFilesWithConfig(serverCtx.Config)
	output := captureStdout(t, func() { server.ShowValidatorCmd().ExecuteContext(ctx) })

	// output must match the locally stored priv_validator key
	privKey := loadFilePVKey(t, serverCtx.Config.PrivValidatorKeyFile())
	bz, err := tmjson.Marshal(privKey.PubKey)
	require.NoError(t, err)
	require.Equal(t, string(bz), output)
}

func TestShowValidatorWithKMS(t *testing.T) {
	testCommon := newPrecedenceCommon(t)

	serverCtx := &server.Context{}
	ctx := context.WithValue(context.Background(), server.ServerContextKey, serverCtx)

	if err := testCommon.cmd.ExecuteContext(ctx); err != cancelledInPreRun {
		t.Fatalf("function failed with [%T] %v", err, err)
	}

	// ostracon init & create the server config file
	initFilesWithConfig(serverCtx.Config)

	chainID, err := server.LoadChainID(serverCtx.Config)
	require.NoError(t, err)

	// remove config file
	if tmos.FileExists(serverCtx.Config.PrivValidatorKeyFile()) {
		err := os.Remove(serverCtx.Config.PrivValidatorKeyFile())
		require.NoError(t, err)
	}

	privval.WithMockKMS(t, t.TempDir(), chainID, func(addr string, privKey crypto.PrivKey) {
		serverCtx.Config.PrivValidatorListenAddr = addr
		require.NoFileExists(t, serverCtx.Config.PrivValidatorKeyFile())
		output := captureStdout(t, func() { server.ShowValidatorCmd().ExecuteContext(ctx) })
		require.NoError(t, err)

		// output must contains the KMS public key
		bz, err := tmjson.Marshal(privKey.PubKey())
		require.NoError(t, err)
		expected := string(bz)
		require.Contains(t, output, expected)
	})
}

func TestShowValidatorWithInefficientKMSAddress(t *testing.T) {
	testCommon := newPrecedenceCommon(t)

	serverCtx := &server.Context{}
	ctx := context.WithValue(context.Background(), server.ServerContextKey, serverCtx)

	if err := testCommon.cmd.ExecuteContext(ctx); err != cancelledInPreRun {
		t.Fatalf("function failed with [%T] %v", err, err)
	}

	// ostracon init & create the server config file
	initFilesWithConfig(serverCtx.Config)

	// remove config file
	if tmos.FileExists(serverCtx.Config.PrivValidatorKeyFile()) {
		err := os.Remove(serverCtx.Config.PrivValidatorKeyFile())
		require.NoError(t, err)
	}

	serverCtx.Config.PrivValidatorListenAddr = "127.0.0.1:inefficient"
	err := server.ShowValidatorCmd().ExecuteContext(ctx)
	require.Error(t, err)
}

func TestLoadChainID(t *testing.T) {
	expected := "c57861"
	config := cfg.ResetTestRootWithChainID("TestLoadChainID", expected)
	defer func() {
		var _ = os.RemoveAll(config.RootDir)
	}()

	require.FileExists(t, config.GenesisFile())
	genDoc, err := types.GenesisDocFromFile(config.GenesisFile())
	require.NoError(t, err)
	require.Equal(t, expected, genDoc.ChainID)

	chainID, err := server.LoadChainID(config)
	require.NoError(t, err)
	require.Equal(t, expected, chainID)
}

func TestLoadChainIDWithoutStateDB(t *testing.T) {
	expected := "c34091"
	config := cfg.ResetTestRootWithChainID("TestLoadChainID", expected)
	defer func() {
		var _ = os.RemoveAll(config.RootDir)
	}()

	config.DBBackend = "goleveldb"
	config.DBPath = "/../path with containing chars that cannot be used\\/:*?\"<>|\x00"

	_, err := server.LoadChainID(config)
	require.Error(t, err)
}

func initFilesWithConfig(config *cfg.Config) error {
	// private validator
	privValKeyFile := config.PrivValidatorKeyFile()
	privValStateFile := config.PrivValidatorStateFile()
	privKeyType := config.PrivValidatorKeyType()
	var pv *privval.FilePV
	if tmos.FileExists(privValKeyFile) {
		pv = privval.LoadFilePV(privValKeyFile, privValStateFile)
		logger.Info("Found private validator", "keyFile", privValKeyFile,
			"stateFile", privValStateFile)
	} else {
		var err error
		pv, err = privval.GenFilePV(privValKeyFile, privValStateFile, privKeyType)
		if err != nil {
			return err
		}
		if pv != nil {
			pv.Save()
		}
		logger.Info("Generated private validator", "keyFile", privValKeyFile,
			"stateFile", privValStateFile)
	}

	nodeKeyFile := config.NodeKeyFile()
	if tmos.FileExists(nodeKeyFile) {
		logger.Info("Found node key", "path", nodeKeyFile)
	} else {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return err
		}
		logger.Info("Generated node key", "path", nodeKeyFile)
	}

	// genesis file
	genFile := config.GenesisFile()
	if tmos.FileExists(genFile) {
		logger.Info("Found genesis file", "path", genFile)
	} else {
		genDoc := types.GenesisDoc{
			ChainID:         fmt.Sprintf("test-chain-%v", tmrand.Str(6)),
			GenesisTime:     tmtime.Now(),
			ConsensusParams: types.DefaultConsensusParams(),
			VoterParams:     types.DefaultVoterParams(),
		}
		pubKey, err := pv.GetPubKey()
		if err != nil {
			return fmt.Errorf("can't get pubkey: %w", err)
		}
		genDoc.Validators = []types.GenesisValidator{{
			Address: pubKey.Address(),
			PubKey:  pubKey,
			Power:   10,
		}}

		if err := genDoc.SaveAs(genFile); err != nil {
			return err
		}
		logger.Info("Generated genesis file", "path", genFile)
	}

	return nil
}

func loadFilePVKey(t *testing.T, file string) privval.FilePVKey {
	// output must match the locally stored priv_validator key
	keyJSONBytes, err := ioutil.ReadFile(file)
	require.NoError(t, err)
	privKey := privval.FilePVKey{}
	err = tmjson.Unmarshal(keyJSONBytes, &privKey)
	require.NoError(t, err)
	return privKey
}

func captureStdout(t *testing.T, fnc func()) string {
	t.Helper()
	backup := os.Stdout
	defer func() {
		os.Stdout = backup
	}()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("fail pipe: %v", err)
	}
	os.Stdout = w
	fnc()
	w.Close()
	var buffer bytes.Buffer
	if n, err := buffer.ReadFrom(r); err != nil {
		t.Fatalf("fail read buf: %v - number: %v", err, n)
	}
	output := buffer.String()
	return output[:len(output)-1]
}
