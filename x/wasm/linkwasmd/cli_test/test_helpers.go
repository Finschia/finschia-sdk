package clitest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/viper"

	// collectionModule "github.com/line/lbm-sdk/v2/x/collection"
	// tokenModule "github.com/line/lbm-sdk/v2/x/token"
	"github.com/line/lbm-sdk/v2/x/wasm/linkwasmd/types"

	clientkeys "github.com/line/lbm-sdk/v2/client/keys"
	tmhttp "github.com/line/ostracon/rpc/client/http"

	"github.com/stretchr/testify/require"

	cfg "github.com/line/ostracon/config"
	tmctypes "github.com/line/ostracon/rpc/core/types"
	osttypes "github.com/line/ostracon/types"

	wasmtypes "github.com/line/lbm-sdk/v2/x/wasm/internal/types"
	"github.com/line/lbm-sdk/v2/x/wasm/linkwasmd/app"

	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/crypto/keyring"
	"github.com/line/lbm-sdk/v2/simapp"
	"github.com/line/lbm-sdk/v2/tests"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/auth"
	"github.com/line/lbm-sdk/v2/x/distribution"
	"github.com/line/lbm-sdk/v2/x/gov"
	"github.com/line/lbm-sdk/v2/x/slashing"
	"github.com/line/lbm-sdk/v2/x/staking"
)

const (
	denom        = "stake"
	keyFoo       = "foo"
	keyBar       = "bar"
	fooDenom     = "foot"
	feeDenom     = "feet"
	fee2Denom    = "fee2t"
	keyBaz       = "baz"
	keyVesting   = "vesting"
	keyFooBarBaz = "foobarbaz"

	DenomStake = "stake2"
	DenomLink  = "link"
	UserTina   = "tina"
	UserKevin  = "kevin"
	UserRinah  = "rinah"
	UserBrian  = "brian"
	UserEvelyn = "evelyn"
	UserSam    = "sam"
)

const (
	namePrefix        = "node"
	networkNamePrefix = "line-linkdnode-testnet-"
)

var curPort int32 = 26655

var (
	TotalCoins = sdk.NewCoins(
		sdk.NewCoin(DenomLink, sdk.TokensFromConsensusPower(6000)),
		sdk.NewCoin(DenomStake, sdk.TokensFromConsensusPower(600000000)),
		sdk.NewCoin(fee2Denom, sdk.TokensFromConsensusPower(2000000)),
		sdk.NewCoin(feeDenom, sdk.TokensFromConsensusPower(2000000)),
		sdk.NewCoin(fooDenom, sdk.TokensFromConsensusPower(2000)),
		sdk.NewCoin(denom, sdk.TokensFromConsensusPower(300)), // We don't use inflation
		// sdk.NewCoin(denom, sdk.TokensFromConsensusPower(300).Add(sdk.NewInt(12))), // add coins from inflation
	)

	startCoins = sdk.NewCoins(
		sdk.NewCoin(fee2Denom, sdk.TokensFromConsensusPower(1000000)),
		sdk.NewCoin(feeDenom, sdk.TokensFromConsensusPower(1000000)),
		sdk.NewCoin(fooDenom, sdk.TokensFromConsensusPower(1000)),
		sdk.NewCoin(denom, sdk.TokensFromConsensusPower(150)),
	)

	vestingCoins = sdk.NewCoins(
		sdk.NewCoin(feeDenom, sdk.TokensFromConsensusPower(500000)),
	)

	// coins we set during ./.initialize.sh
	defaultCoins = sdk.NewCoins(
		sdk.NewCoin(DenomLink, sdk.TokensFromConsensusPower(1000)),
		sdk.NewCoin(DenomStake, sdk.TokensFromConsensusPower(100000000)),
	)
)

func init() {
	testnet := false
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(types.Bech32PrefixAcc(testnet), types.Bech32PrefixAccPub(testnet))
	config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr(testnet), types.Bech32PrefixValPub(testnet))
	config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr(testnet), types.Bech32PrefixConsPub(testnet))
	config.SetCoinType(types.CoinType)
	config.SetFullFundraiserPath(types.FullFundraiserPath)
	config.Seal()
}

// ___________________________________________________________________________________
// Fixtures

// Fixtures is used to setup the testing environment
type Fixtures struct {
	BuildDir      string
	RootDir       string
	LinkdBinary   string
	LinkcliBinary string
	ChainID       string
	RPCAddr       string
	Port          string
	LinkdHome     string
	LinkcliHome   string
	P2PAddr       string
	P2PPort       string
	Moniker       string
	BridgeIP      string
	T             *testing.T
}

// NewFixtures creates a new instance of Fixtures with many vars set
func NewFixtures(t *testing.T) *Fixtures {
	tmpDir := path.Join(os.ExpandEnv("$HOME"), ".linkwasmtest")
	err := os.MkdirAll(tmpDir, os.ModePerm)
	require.NoError(t, err)
	tmpDir, err = ioutil.TempDir(tmpDir, "linkwasm_integration_"+strings.Split(t.Name(), "/")[0]+"_")
	require.NoError(t, err)

	servAddr, servPort := newTCPAddr(t)
	p2pAddr, p2pPort := newTCPAddr(t)

	buildDir := os.Getenv("BUILDDIR")
	if buildDir == "" {
		buildDir, err = filepath.Abs("../build/")
		require.NoError(t, err)
	}

	return &Fixtures{
		T:             t,
		BuildDir:      buildDir,
		RootDir:       tmpDir,
		LinkdBinary:   filepath.Join(buildDir, "linkwasmd"),
		LinkcliBinary: filepath.Join(buildDir, "linkwasmcli"),
		LinkdHome:     filepath.Join(tmpDir, ".linkwasmd"),
		LinkcliHome:   filepath.Join(tmpDir, ".linkwasmcli"),
		RPCAddr:       servAddr,
		P2PAddr:       p2pAddr,
		Port:          servPort,
		P2PPort:       p2pPort,
		Moniker:       "", // initialized by LDInit
		BridgeIP:      "",
	}
}

func newTCPAddr(t *testing.T) (addr, port string) {
	portI := atomic.AddInt32(&curPort, 1)
	require.Less(t, portI, int32(32768), "A new port should be less than ip_local_port_range.min")

	port = fmt.Sprintf("%d", portI)
	addr = fmt.Sprintf("tcp://0.0.0.0:%s", port)
	return
}

func (f *Fixtures) LogResult(isSuccess bool, stdOut, stdErr string) {
	if !isSuccess {
		f.T.Error(stdErr)
	} else {
		f.T.Log(stdOut)
	}
}

func (f Fixtures) Clone() *Fixtures {
	newF := NewFixtures(f.T)
	newF.ChainID = f.ChainID

	tests.ExecuteT(newF.T, fmt.Sprintf("cp -r %s/ %s/", f.RootDir, newF.RootDir), "")

	return newF
}

// GenesisFile returns the path of the genesis file
func (f Fixtures) GenesisFile() string {
	return filepath.Join(f.LinkdHome, "config", "genesis.json")
}

func (f Fixtures) PrivValidatorKeyFile() string {
	return filepath.Join(f.LinkdHome, "config", "priv_validator_key.json")
}

// GenesisFile returns the application's genesis state
func (f Fixtures) GenesisState() simapp.GenesisState {
	cdc := codec.New()
	genDoc, err := osttypes.GenesisDocFromFile(f.GenesisFile())
	require.NoError(f.T, err)

	var appState simapp.GenesisState
	require.NoError(f.T, cdc.UnmarshalJSON(genDoc.AppState, &appState))
	return appState
}

// InitFixtures is called at the beginning of a test  and initializes a chain
// with 1 validator.
func InitFixtures(t *testing.T) (f *Fixtures) {
	f = NewFixtures(t)

	// reset test state
	f.UnsafeResetAll()

	// ensure keystore has foo and bar keys
	f.KeysDelete(keyFoo)
	f.KeysDelete(keyBar)
	f.KeysDelete(keyBaz)
	f.KeysDelete(keyFooBarBaz)
	f.KeysAdd(keyFoo)
	f.KeysAdd(keyBar)
	f.KeysAdd(keyBaz)
	f.KeysAdd(keyVesting)
	f.KeysAdd(keyFooBarBaz, "--multisig-threshold=2", fmt.Sprintf(
		"--multisig=%s,%s,%s", keyFoo, keyBar, keyBaz))

	// ensure keystore to have user keys
	f.KeysDelete(UserTina)
	f.KeysDelete(UserKevin)
	f.KeysDelete(UserRinah)
	f.KeysDelete(UserBrian)
	f.KeysDelete(UserEvelyn)
	f.KeysDelete(UserSam)
	f.KeysAdd(UserTina)
	f.KeysAdd(UserKevin)
	f.KeysAdd(UserRinah)
	f.KeysAdd(UserBrian)
	f.KeysAdd(UserEvelyn)
	f.KeysAdd(UserSam)

	// ensure that CLI output is in JSON format
	f.CLIConfig("output", "json")

	// NOTE: LDInit sets the ChainID
	f.LDInit(keyFoo)

	f.CLIConfig("chain-id", f.ChainID)
	f.CLIConfig("broadcast-mode", "block")
	f.CLIConfig("trust-node", "true")
	f.CLIConfig("keyring-backend", "test")

	// start an account with tokens
	f.AddGenesisAccount(f.KeyAddress(keyFoo), startCoins)
	// f.AddGenesisAccount(f.KeyAddress(keyBar), startCoins)
	f.AddGenesisAccount(
		f.KeyAddress(keyVesting), startCoins,
		fmt.Sprintf("--vesting-amount=%s", vestingCoins),
		fmt.Sprintf("--vesting-start-time=%d", time.Now().UTC().UnixNano()),
		fmt.Sprintf("--vesting-end-time=%d", time.Now().Add(60*time.Second).UTC().UnixNano()),
	)

	// add genesis accounts for testing
	f.AddGenesisAccount(f.KeyAddress(UserTina), defaultCoins)
	f.AddGenesisAccount(f.KeyAddress(UserKevin), defaultCoins)
	f.AddGenesisAccount(f.KeyAddress(UserRinah), defaultCoins)
	f.AddGenesisAccount(f.KeyAddress(UserBrian), defaultCoins)
	f.AddGenesisAccount(f.KeyAddress(UserEvelyn), defaultCoins)
	f.AddGenesisAccount(f.KeyAddress(UserSam), defaultCoins)

	f.GenTx(keyFoo)
	f.CollectGenTxs()

	return f
}

// Cleanup is meant to be run at the end of a test to clean up an remaining test state
func (f *Fixtures) Cleanup(dirs ...string) {
	clean := append(dirs, f.RootDir)
	for _, d := range clean {
		_ = os.RemoveAll(d)
	}
}

// Flags returns the flags necessary for making most CLI calls
func (f *Fixtures) Flags() string {
	return fmt.Sprintf("--home=%s --node=%s", f.LinkcliHome, f.RPCAddr)
}

// ___________________________________________________________________________________
// linkd

// UnsafeResetAll is linkd unsafe-reset-all
func (f *Fixtures) UnsafeResetAll(flags ...string) {
	cmd := fmt.Sprintf("%s --home=%s unsafe-reset-all", f.LinkdBinary, f.LinkdHome)
	executeWrite(f.T, addFlags(cmd, flags))
	err := os.RemoveAll(filepath.Join(f.LinkdHome, "config", "gentx"))
	require.NoError(f.T, err)
}

// LDInit is linkd init
// NOTE: LDInit sets the ChainID for the Fixtures instance
func (f *Fixtures) LDInit(moniker string, flags ...string) {
	f.Moniker = moniker
	cmd := fmt.Sprintf("%s init -o --home=%s %s", f.LinkdBinary, f.LinkdHome, moniker)
	_, stderr := tests.ExecuteT(f.T, addFlags(cmd, flags), "")

	var chainID string
	var initRes map[string]json.RawMessage

	err := json.Unmarshal([]byte(stderr), &initRes)
	require.NoError(f.T, err)

	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(f.T, err)

	f.ChainID = chainID
}

// AddGenesisAccount is linkd add-genesis-account
func (f *Fixtures) AddGenesisAccount(address sdk.AccAddress, coins sdk.Coins, flags ...string) {
	cmd := fmt.Sprintf("%s add-genesis-account %s %s --home=%s --keyring-backend=test", f.LinkdBinary, address, coins, f.LinkdHome)
	executeWriteCheckErr(f.T, addFlags(cmd, flags))
}

// GenTx is linkd gentx
func (f *Fixtures) GenTx(name string, flags ...string) {
	cmd := fmt.Sprintf("%s gentx --name=%s --home=%s --home-client=%s --keyring-backend=test", f.LinkdBinary, name, f.LinkdHome, f.LinkcliHome)
	executeWriteCheckErr(f.T, addFlags(cmd, flags))
}

// CollectGenTxs is linkd collect-gentxs
func (f *Fixtures) CollectGenTxs(flags ...string) {
	cmd := fmt.Sprintf("%s collect-gentxs --home=%s", f.LinkdBinary, f.LinkdHome)
	executeWriteCheckErr(f.T, addFlags(cmd, flags))
}

// LDStart runs linkd start with the appropriate flags and returns a process
func (f *Fixtures) LDStart(flags ...string) *tests.Process {
	cmd := fmt.Sprintf("%s start --home=%s --rpc.laddr=%v --p2p.laddr=%v", f.LinkdBinary, f.LinkdHome, f.RPCAddr, f.P2PAddr)
	proc := tests.GoExecuteT(f.T, addFlags(cmd, flags))
	defer func() {
		if v := recover(); v != nil {
			stdout, stderr, err := proc.ReadAll()
			require.NoError(f.T, err)
			// Log for start command
			f.T.Log(cmd, string(stdout))
			f.T.Log(cmd, string(stderr))
			f.T.Fatal(v)
		}
	}()
	WaitForTMStart(f.Port)
	tests.WaitForNextNBlocksTM(1, f.Port)
	return proc
}

// LDTendermint returns the results of linkd tendermint [query]
func (f *Fixtures) LDTendermint(query string) string {
	cmd := fmt.Sprintf("%s tendermint %s --home=%s", f.LinkdBinary, query, f.LinkdHome)
	success, stdout, stderr := executeWriteRetStdStreams(f.T, cmd)
	require.Empty(f.T, stderr)
	require.True(f.T, success)
	return strings.TrimSpace(stdout)
}

// ValidateGenesis runs linkd validate-genesis
func (f *Fixtures) ValidateGenesis() {
	cmd := fmt.Sprintf("%s validate-genesis --home=%s", f.LinkdBinary, f.LinkdHome)
	executeWriteCheckErr(f.T, cmd)
}

// ___________________________________________________________________________________
// linkcli rest-server
func (f *Fixtures) RestServerStart(port int, flags ...string) (*tests.Process, error) {
	cmd := fmt.Sprintf("%s rest-server --home=%s --laddr=%s", f.LinkcliBinary, f.LinkcliHome, fmt.Sprintf("tcp://0.0.0.0:%d", port))
	proc := tests.GoExecuteTWithStdout(f.T, addFlags(cmd, flags))
	defer func() {
		if v := recover(); v != nil {
			stdout, stderr, err := proc.ReadAll()
			if err != nil {
				fmt.Println(err)
				f.T.Fail()
			}
			f.T.Log(stdout)
			f.T.Log(stderr)
		}
	}()
	tests.WaitForNextNBlocksTM(1, f.Port)
	return proc, nil
}

// ___________________________________________________________________________________
// linkcli keys

// KeysDelete is linkcli keys delete
func (f *Fixtures) KeysDelete(name string, flags ...string) {
	cmd := fmt.Sprintf("%s keys delete --keyring-backend=test --home=%s %s", f.LinkcliBinary, f.LinkcliHome, name)
	executeWrite(f.T, addFlags(cmd, append(append(flags, "-y"), "-f")))
}

// KeysAdd is linkcli keys add
func (f *Fixtures) KeysAdd(name string, flags ...string) {
	cmd := fmt.Sprintf("%s keys add --keyring-backend=test --home=%s %s", f.LinkcliBinary, f.LinkcliHome, name)
	executeWriteCheckErr(f.T, addFlags(cmd, flags))
}

// KeysAddRecover prepares linkcli keys add --recover
func (f *Fixtures) KeysAddRecover(name, mnemonic string, flags ...string) (exitSuccess bool, stdout, stderr string) {
	cmd := fmt.Sprintf("%s keys add --keyring-backend=test --home=%s --recover %s", f.LinkcliBinary, f.LinkcliHome, name)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), mnemonic)
}

// KeysAddRecoverHDPath prepares linkcli keys add --recover --account --index
func (f *Fixtures) KeysAddRecoverHDPath(name, mnemonic string, account uint32, index uint32, flags ...string) {
	cmd := fmt.Sprintf("%s keys add --keyring-backend=test --home=%s --recover %s --account %d --index %d", f.LinkcliBinary, f.LinkcliHome, name, account, index)
	executeWriteCheckErr(f.T, addFlags(cmd, flags), mnemonic)
}

// KeysShow is linkcli keys show
func (f *Fixtures) KeysShow(name string, flags ...string) keyring.KeyOutput {
	cmd := fmt.Sprintf("%s keys show --keyring-backend=test --home=%s %s", f.LinkcliBinary, f.LinkcliHome, name)
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var ko keyring.KeyOutput
	err := clientkeys.UnmarshalJSON([]byte(out), &ko)
	require.NoError(f.T, err)
	return ko
}

// KeyAddress returns the SDK account address from the key
func (f *Fixtures) KeyAddress(name string) sdk.AccAddress {
	ko := f.KeysShow(name)
	accAddr, err := sdk.AccAddressFromBech32(ko.Address)
	require.NoError(f.T, err)
	return accAddr
}

// ___________________________________________________________________________________
// linkcli config

// CLIConfig is linkcli config
func (f *Fixtures) CLIConfig(key, value string, flags ...string) {
	cmd := fmt.Sprintf("%s config --home=%s %s %s", f.LinkcliBinary, f.LinkcliHome, key, value)
	executeWriteCheckErr(f.T, addFlags(cmd, flags))
}

// ___________________________________________________________________________________
// linkcli tx send/sign/broadcast

// TxSend is linkcli tx send
func (f *Fixtures) TxSend(from string, to sdk.AccAddress, amount sdk.Coin, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx send --keyring-backend=test %s %s %s %v", f.LinkcliBinary, from, to, amount, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

// TxSign is linkcli tx sign
func (f *Fixtures) TxSign(signer, fileName string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx sign %v --keyring-backend=test --from=%s %v", f.LinkcliBinary, f.Flags(), signer, fileName)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

// TxBroadcast is linkcli tx broadcast
func (f *Fixtures) TxBroadcast(fileName string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx broadcast %v %v", f.LinkcliBinary, f.Flags(), fileName)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

// TxEncode is linkcli tx encode
func (f *Fixtures) TxEncode(fileName string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx encode %v %v", f.LinkcliBinary, f.Flags(), fileName)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

// TxMultisign is linkcli tx multisign
func (f *Fixtures) TxMultisign(fileName, name string, signaturesFiles []string,
	flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx multisign --keyring-backend=test %v %s %s %s", f.LinkcliBinary, f.Flags(),
		fileName, name, strings.Join(signaturesFiles, " "),
	)
	return executeWriteRetStdStreams(f.T, cmd)
}

// ___________________________________________________________________________________
// linkcli tx staking

// TxStakingCreateValidator is linkcli tx staking create-validator
func (f *Fixtures) TxStakingCreateValidator(from, consPubKey string, amount sdk.Coin, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx staking create-validator %v --keyring-backend=test --from=%s --pubkey=%s", f.LinkcliBinary, f.Flags(), from, consPubKey)
	cmd += fmt.Sprintf(" --amount=%v --moniker=%v --commission-rate=%v", amount, from, "0.05")
	cmd += fmt.Sprintf(" --commission-max-rate=%v --commission-max-change-rate=%v", "0.20", "0.10")
	cmd += fmt.Sprintf(" --min-self-delegation=%v", "1")
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

// TxStakingUnbond is linkcli tx staking unbond
func (f *Fixtures) TxStakingUnbond(from, shares string, validator sdk.ValAddress, flags ...string) bool {
	cmd := fmt.Sprintf("%s tx staking unbond --keyring-backend=test %s %v --from=%s %v", f.LinkcliBinary, validator, shares, from, f.Flags())
	return executeWrite(f.T, addFlags(cmd, flags))
}

// ___________________________________________________________________________________
// linkcli tx gov

// TxGovSubmitProposal is linkcli tx gov submit-proposal
func (f *Fixtures) TxGovSubmitProposal(from, typ, title, description string, deposit sdk.Coin, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx gov submit-proposal %v --keyring-backend=test --from=%s --type=%s", f.LinkcliBinary, f.Flags(), from, typ)
	cmd += fmt.Sprintf(" --title=%s --description=%s --deposit=%s", title, description, deposit)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

// TxGovDeposit is linkcli tx gov deposit
func (f *Fixtures) TxGovDeposit(proposalID int, from string, amount sdk.Coin, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx gov deposit %d %s --keyring-backend=test --from=%s %v", f.LinkcliBinary, proposalID, amount, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

// TxGovVote is linkcli tx gov vote
func (f *Fixtures) TxGovVote(proposalID int, option gov.VoteOption, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx gov vote %d %s --keyring-backend=test --from=%s %v", f.LinkcliBinary, proposalID, option, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

// TxGovSubmitParamChangeProposal executes a CLI parameter change proposal
// submission.
func (f *Fixtures) TxGovSubmitParamChangeProposal(
	from, proposalPath string, deposit sdk.Coin, flags ...string,
) (bool, string, string) {
	cmd := fmt.Sprintf(
		"%s tx gov submit-proposal param-change %s --keyring-backend=test --from=%s %v",
		f.LinkcliBinary, proposalPath, from, f.Flags(),
	)

	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

// TxGovSubmitCommunityPoolSpendProposal executes a CLI community pool spend proposal
// submission.
func (f *Fixtures) TxGovSubmitCommunityPoolSpendProposal(
	from, proposalPath string, deposit sdk.Coin, flags ...string,
) (bool, string, string) {
	cmd := fmt.Sprintf(
		"%s tx gov submit-proposal community-pool-spend %s --keyring-backend=test --from=%s %v",
		f.LinkcliBinary, proposalPath, from, f.Flags(),
	)

	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

// ___________________________________________________________________________________
// linkcli tx token

func (f *Fixtures) TxTokenIssue(from string, to sdk.AccAddress, name, meta string, symbol string, amount int64, decimals int64, mintable bool, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token issue %s %s %s %s --total-supply=%d --decimals=%d --mintable=%t --meta=%s %v", f.LinkcliBinary, from, to.String(), name, symbol, amount, decimals, mintable, meta, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}
func (f *Fixtures) TxTokenMint(from string, contractID string, to string, amount string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token mint %s %s %s %s %v", f.LinkcliBinary, from, contractID, to, amount, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenBurn(from, contractID, amount string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token burn %s %s %s %v", f.LinkcliBinary, from, contractID, amount, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenBurnFrom(proxy, contractID string, from sdk.AccAddress, amount int64, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token burn-from %s %s %s %d %v", f.LinkcliBinary, proxy, contractID, from, amount, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenGrantPerm(from string, to string, contractID, action string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token grant %s %s %s %s %v", f.LinkcliBinary, from, contractID, to, action, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenRevokePerm(from string, contractID, action string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token revoke %s %s %s %v", f.LinkcliBinary, from, contractID, action, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenModify(owner, contractID, field, newValue string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token modify %s %s %s %s %v", f.LinkcliBinary, owner, contractID, field, newValue, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenTransfer(from string, to sdk.AccAddress, symbol string, amount int64, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token transfer %s %s %s %d %v", f.LinkcliBinary, from, to, symbol, amount, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenTransferFrom(proxy string, contractID string, from sdk.AccAddress, to sdk.AccAddress, amount int64, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token transfer-from %s %s %s %s %d %v", f.LinkcliBinary, proxy, contractID, from, to, amount, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenApprove(approver string, contractID string, proxyAddress sdk.AccAddress, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token approve %s %s %s %v", f.LinkcliBinary, approver, contractID, proxyAddress, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

// ___________________________________________________________________________________
// linkcli tx collection

func (f *Fixtures) TxTokenCreateCollection(from string, name, meta, baseImgURI string, flags ...string) (bool, string,
	string) {
	cmd := fmt.Sprintf("%s tx collection create %s %s %s %s %v", f.LinkcliBinary, from, name, meta, baseImgURI, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}
func (f *Fixtures) TxTokenIssueFTCollection(from string, contractID string, to sdk.AccAddress, name, meta string, amount int64, decimals int64, mintable bool, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx collection issue-ft %s %s %s %s %s --total-supply=%d --decimals=%d --mintable=%t %v", f.LinkcliBinary, from, contractID, to.String(), name, meta, amount, decimals, mintable, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenMintFTCollection(from string, contractID string, to string, amount string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx collection mint-ft %s %s %s %s %v", f.LinkcliBinary, from, contractID, to, amount, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenBurnFTCollection(from string, contractID, tokenID string, amount int64, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx collection burn-ft %s %s %s %d %v", f.LinkcliBinary, from, contractID, tokenID, amount, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenTransferFTCollection(from string, contractID string, to string, amount string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx collection transfer-ft %s %s %s %s %v", f.LinkcliBinary, from, contractID, to, amount, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

func (f *Fixtures) TxTokenIssueNFTCollection(from string, contractID, name, meta string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx collection issue-nft %s %s %s %s %v", f.LinkcliBinary, from, contractID, name, meta, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenMintNFTCollection(from string, contractID string, to string, mintNFTParam string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx collection mint-nft %s %s %s %s %v", f.LinkcliBinary, from, contractID, to, mintNFTParam, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenBurnNFTCollection(from string, contractID, tokenID string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx collection burn-nft %s %s %s %v", f.LinkcliBinary, from, contractID, tokenID, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxTokenTransferNFTCollection(from string, contractID string, to string, tokenID string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx collection transfer-nft %s %s %s %s %v", f.LinkcliBinary, from, contractID, to, tokenID, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), clientkeys.DefaultKeyPass)
}

func (f *Fixtures) TxCollectionModify(owner, contractID, tokenType, tokenIndex, field, newValue string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx collection modify %s %s %s %s --token-type %s --token-index %s %v",
		f.LinkcliBinary, owner, contractID, field, newValue, tokenType, tokenIndex, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxCollectionGrantPerm(from string, to sdk.AccAddress, contractID, action string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx collection grant %s %s %s %s %v", f.LinkcliBinary, from, contractID, to, action, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxCollectionRevokePerm(from string, contractID, action string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx collection revoke %s %s %s %v", f.LinkcliBinary, from, contractID, action, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxCollectionApprove(approver string, contractID string, proxyAd sdk.AccAddress, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx collection approve %s %s %s %v", f.LinkcliBinary, approver, contractID, proxyAd, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxEmpty(from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx empty %s %v", f.LinkcliBinary, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

// ___________________________________________________________________________________
// linkcli query account

// QueryAccount is linkcli query account
func (f *Fixtures) QueryAccount(address sdk.AccAddress, flags ...string) auth.BaseAccount {
	cmd := fmt.Sprintf("%s query account %s %v", f.LinkcliBinary, address, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(f.T, err, "out %v, err %v", out, err)
	value := initRes["value"]
	var acc auth.BaseAccount
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	err = cdc.UnmarshalJSON(value, &acc)
	require.NoError(f.T, err, "value %v, err %v", string(value), err)
	return acc
}

// ___________________________________________________________________________________
// linkcli query tx

// QueryTx is linkcli query tx
func (f *Fixtures) QueryTx(hash string) *sdk.TxResponse {
	cmd := fmt.Sprintf("%s query tx %s %v", f.LinkcliBinary, hash, f.Flags())
	out, _ := tests.ExecuteT(f.T, cmd, "")
	var result sdk.TxResponse
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return &result
}

// QueryTxInvalid query tx with wrong hash and compare expected error
func (f *Fixtures) QueryTxInvalid(expectedErr error, hash string) {
	cmd := fmt.Sprintf("%s query tx %s %v", f.LinkcliBinary, hash, f.Flags())
	_, err := tests.ExecuteT(f.T, cmd, "")
	require.EqualError(f.T, expectedErr, err)
}

// ___________________________________________________________________________________
// linkcli query txs

// QueryTxs is linkcli query txs
func (f *Fixtures) QueryTxs(page, limit int, flags ...string) *sdk.SearchTxsResult {
	cmd := fmt.Sprintf("%s query txs --page=%d --limit=%d %v", f.LinkcliBinary, page, limit, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var result sdk.SearchTxsResult
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return &result
}

// QueryTxsInvalid query txs with wrong parameters and compare expected error
func (f *Fixtures) QueryTxsInvalid(expectedErr error, page, limit int, flags ...string) {
	cmd := fmt.Sprintf("%s query txs --page=%d --limit=%d %v", f.LinkcliBinary, page, limit, f.Flags())
	_, err := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	require.EqualError(f.T, expectedErr, err)
}

// ___________________________________________________________________________________
// linkcli query block

func (f *Fixtures) QueryLatestBlock(flags ...string) *tmctypes.ResultBlock {
	cmd := fmt.Sprintf("%s query block %v", f.LinkcliBinary, f.Flags())
	out, _ := tests.ExecuteT(f.T, cmd, "")
	var result tmctypes.ResultBlock
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return &result
}

func (f *Fixtures) QueryBlockWithHeight(height int, flags ...string) *tmctypes.ResultBlock {
	cmd := fmt.Sprintf("%s query block %d %v", f.LinkcliBinary, height, f.Flags())
	out, _ := tests.ExecuteT(f.T, cmd, "")
	var result tmctypes.ResultBlock
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return &result
}

// ___________________________________________________________________________________
// linkcli query staking

// QueryStakingValidator is linkcli query staking validator
func (f *Fixtures) QueryStakingValidator(valAddr sdk.ValAddress, flags ...string) staking.Validator {
	cmd := fmt.Sprintf("%s query staking validator %s %v", f.LinkcliBinary, valAddr, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var validator staking.Validator
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &validator)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return validator
}

// QueryStakingUnbondingDelegationsFrom is linkcli query staking unbonding-delegations-from
func (f *Fixtures) QueryStakingUnbondingDelegationsFrom(valAddr sdk.ValAddress, flags ...string) []staking.UnbondingDelegation {
	cmd := fmt.Sprintf("%s query staking unbonding-delegations-from %s %v", f.LinkcliBinary, valAddr, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var ubds []staking.UnbondingDelegation
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &ubds)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return ubds
}

// QueryStakingDelegationsTo is linkcli query staking delegations-to
func (f *Fixtures) QueryStakingDelegationsTo(valAddr sdk.ValAddress, flags ...string) []staking.Delegation {
	cmd := fmt.Sprintf("%s query staking delegations-to %s %v", f.LinkcliBinary, valAddr, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var delegations []staking.Delegation
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &delegations)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return delegations
}

// QueryStakingPool is linkcli query staking pool
func (f *Fixtures) QueryStakingPool(flags ...string) staking.Pool {
	cmd := fmt.Sprintf("%s query staking pool %v", f.LinkcliBinary, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var pool staking.Pool
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &pool)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return pool
}

// QueryStakingParameters is linkcli query staking parameters
func (f *Fixtures) QueryStakingParameters(flags ...string) staking.Params {
	cmd := fmt.Sprintf("%s query staking params %v", f.LinkcliBinary, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var params staking.Params
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &params)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return params
}

// ___________________________________________________________________________________
// linkcli query gov

// QueryGovParamDeposit is linkcli query gov param deposit
func (f *Fixtures) QueryGovParamDeposit() gov.DepositParams {
	cmd := fmt.Sprintf("%s query gov param deposit %s", f.LinkcliBinary, f.Flags())
	out, _ := tests.ExecuteT(f.T, cmd, "")
	var depositParam gov.DepositParams
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &depositParam)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return depositParam
}

// QueryGovParamVoting is linkcli query gov param voting
func (f *Fixtures) QueryGovParamVoting() gov.VotingParams {
	cmd := fmt.Sprintf("%s query gov param voting %s", f.LinkcliBinary, f.Flags())
	out, _ := tests.ExecuteT(f.T, cmd, "")
	var votingParam gov.VotingParams
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &votingParam)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return votingParam
}

// QueryGovParamTallying is linkcli query gov param tallying
func (f *Fixtures) QueryGovParamTallying() gov.TallyParams {
	cmd := fmt.Sprintf("%s query gov param tallying %s", f.LinkcliBinary, f.Flags())
	out, _ := tests.ExecuteT(f.T, cmd, "")
	var tallyingParam gov.TallyParams
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &tallyingParam)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return tallyingParam
}

// QueryGovProposals is linkcli query gov proposals
func (f *Fixtures) QueryGovProposals(flags ...string) gov.Proposals {
	cmd := fmt.Sprintf("%s query gov proposals %v", f.LinkcliBinary, f.Flags())
	stdout, stderr := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	if strings.Contains(stderr, "No matching proposals found") {
		return gov.Proposals{}
	}
	require.Empty(f.T, stderr)
	var out gov.Proposals
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(stdout), &out)
	require.NoError(f.T, err)
	return out
}

// QueryGovProposal is linkcli query gov proposal
func (f *Fixtures) QueryGovProposal(proposalID int, flags ...string) gov.Proposal {
	cmd := fmt.Sprintf("%s query gov proposal %d %v", f.LinkcliBinary, proposalID, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var proposal gov.Proposal
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &proposal)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return proposal
}

// QueryGovVote is linkcli query gov vote
func (f *Fixtures) QueryGovVote(proposalID int, voter sdk.AccAddress, flags ...string) gov.Vote {
	cmd := fmt.Sprintf("%s query gov vote %d %s %v", f.LinkcliBinary, proposalID, voter, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var vote gov.Vote
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &vote)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return vote
}

// QueryGovVotes is linkcli query gov votes
func (f *Fixtures) QueryGovVotes(proposalID int, flags ...string) []gov.Vote {
	cmd := fmt.Sprintf("%s query gov votes %d %v", f.LinkcliBinary, proposalID, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var votes []gov.Vote
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &votes)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return votes
}

// QueryGovDeposit is linkcli query gov deposit
func (f *Fixtures) QueryGovDeposit(proposalID int, depositor sdk.AccAddress, flags ...string) gov.Deposit {
	cmd := fmt.Sprintf("%s query gov deposit %d %s %v", f.LinkcliBinary, proposalID, depositor, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var deposit gov.Deposit
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &deposit)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return deposit
}

// QueryGovDeposits is linkcli query gov deposits
func (f *Fixtures) QueryGovDeposits(propsalID int, flags ...string) []gov.Deposit {
	cmd := fmt.Sprintf("%s query gov deposits %d %v", f.LinkcliBinary, propsalID, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var deposits []gov.Deposit
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &deposits)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return deposits
}

// ___________________________________________________________________________________
// query slashing

// QuerySigningInfo returns the signing info for a validator
func (f *Fixtures) QuerySigningInfo(val string) slashing.ValidatorSigningInfo {
	cmd := fmt.Sprintf("%s query slashing signing-info %s %s", f.LinkcliBinary, val, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var sinfo slashing.ValidatorSigningInfo
	err := cdc.UnmarshalJSON([]byte(res), &sinfo)
	require.NoError(f.T, err)
	return sinfo
}

// QuerySlashingParams is linkcli query slashing params
func (f *Fixtures) QuerySlashingParams() slashing.Params {
	cmd := fmt.Sprintf("%s query slashing params %s", f.LinkcliBinary, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var params slashing.Params
	err := cdc.UnmarshalJSON([]byte(res), &params)
	require.NoError(f.T, err)
	return params
}

// ___________________________________________________________________________________
// query distribution

// QueryRewards returns the rewards of a delegator
func (f *Fixtures) QueryRewards(delAddr sdk.AccAddress, flags ...string) distribution.QueryDelegatorTotalRewardsResponse {
	cmd := fmt.Sprintf("%s query distribution rewards %s %s", f.LinkcliBinary, delAddr, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var rewards distribution.QueryDelegatorTotalRewardsResponse
	err := cdc.UnmarshalJSON([]byte(res), &rewards)
	require.NoError(f.T, err)
	return rewards
}

// ___________________________________________________________________________________
// query supply

// QueryTotalSupply returns the total supply of coins
func (f *Fixtures) QueryTotalSupply(flags ...string) (totalSupply sdk.Coins) {
	cmd := fmt.Sprintf("%s query supply total %s", f.LinkcliBinary, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(res), &totalSupply)
	require.NoError(f.T, err)
	return totalSupply
}

// QueryTotalSupplyOf returns the total supply of a given coin denom
func (f *Fixtures) QueryTotalSupplyOf(denom string, flags ...string) sdk.Int {
	cmd := fmt.Sprintf("%s query supply total %s %s", f.LinkcliBinary, denom, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var supplyOf sdk.Int
	err := cdc.UnmarshalJSON([]byte(res), &supplyOf)
	require.NoError(f.T, err)
	return supplyOf
}

// // ___________________________________________________________________________________
// // query token

// func (f *Fixtures) QueryToken(contractID string, flags ...string) tokenModule.Token {
// 	cmd := fmt.Sprintf("%s query token token %s %s", f.LinkcliBinary, contractID, f.Flags())
// 	cmd = addFlags(cmd, flags)
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var token tokenModule.Token
// 	err := cdc.UnmarshalJSON([]byte(res), &token)
// 	require.NoError(f.T, err)
// 	return token
// }

// func (f *Fixtures) QueryTokenExpectEmpty(contractID string, flags ...string) {
// 	cmd := fmt.Sprintf("%s query token token %s %s", f.LinkcliBinary, contractID, f.Flags())
// 	cmd = addFlags(cmd, flags)
// 	_, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.NotEmpty(f.T, errStr)
// }

// func (f *Fixtures) QueryBalanceToken(contractID string, addr sdk.AccAddress, flags ...string) sdk.Int {
// 	cmd := fmt.Sprintf("%s query token balance %s %s %s", f.LinkcliBinary, contractID, addr.String(), f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var supply sdk.Int
// 	err := cdc.UnmarshalJSON([]byte(res), &supply)
// 	require.NoError(f.T, err)

// 	return supply
// }

// func (f *Fixtures) QuerySupplyToken(contractID string, flags ...string) sdk.Int {
// 	cmd := fmt.Sprintf("%s query token total supply %s %s", f.LinkcliBinary, contractID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var supply sdk.Int
// 	err := cdc.UnmarshalJSON([]byte(res), &supply)
// 	require.NoError(f.T, err)

// 	return supply
// }

// func (f *Fixtures) QueryMintToken(contractID string, flags ...string) sdk.Int {
// 	cmd := fmt.Sprintf("%s query token total mint %s %s", f.LinkcliBinary, contractID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var supply sdk.Int
// 	err := cdc.UnmarshalJSON([]byte(res), &supply)
// 	require.NoError(f.T, err)

// 	return supply
// }

// func (f *Fixtures) QueryBurnToken(contractID string, flags ...string) sdk.Int {
// 	cmd := fmt.Sprintf("%s query token total burn %s %s", f.LinkcliBinary, contractID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var supply sdk.Int
// 	err := cdc.UnmarshalJSON([]byte(res), &supply)
// 	require.NoError(f.T, err)

// 	return supply
// }

// func (f *Fixtures) QueryAccountPermission(addr sdk.AccAddress, contractID string, flags ...string) tokenModule.Permissions {
// 	cmd := fmt.Sprintf("%s query token perm %s %s %s", f.LinkcliBinary, addr, contractID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var pms tokenModule.Permissions
// 	err := cdc.UnmarshalJSON([]byte(res), &pms)
// 	require.NoError(f.T, err)
// 	return pms
// }

// func (f *Fixtures) QueryApprovedToken(contractID string, proxy sdk.AccAddress, approver sdk.AccAddress, flags ...string) bool {
// 	cmd := fmt.Sprintf("%s query token approved %s %s %s %s", f.LinkcliBinary, contractID, proxy, approver, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var isApproved bool
// 	err := cdc.UnmarshalJSON([]byte(res), &isApproved)
// 	require.NoError(f.T, err)
// 	return isApproved
// }

// // ___________________________________________________________________________________
// // query collection
// func (f *Fixtures) QueryBalancesCollection(contractID string, addr sdk.AccAddress) collectionModule.Coins {
// 	cmd := fmt.Sprintf("%s query collection balances %s %s %s", f.LinkcliBinary, contractID, addr.String(), f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var coins collectionModule.Coins
// 	err := cdc.UnmarshalJSON([]byte(res), &coins)
// 	require.NoError(f.T, err)

// 	return coins
// }

// func (f *Fixtures) QueryBalanceCollection(contractID, tokenID string, addr sdk.AccAddress, flags ...string) sdk.Int {
// 	cmd := fmt.Sprintf("%s query collection balance %s %s %s %s", f.LinkcliBinary, contractID, tokenID, addr.String(), f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var supply sdk.Int
// 	err := cdc.UnmarshalJSON([]byte(res), &supply)
// 	require.NoError(f.T, err)

// 	return supply
// }

// func (f *Fixtures) QueryTokenCollection(contractID, tokenID string, flags ...string) collectionModule.Token {
// 	cmd := fmt.Sprintf("%s query collection token %s %s %s", f.LinkcliBinary, contractID, tokenID, f.Flags())
// 	cmd = addFlags(cmd, flags)
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var token collectionModule.Token
// 	err := cdc.UnmarshalJSON([]byte(res), &token)
// 	require.NoError(f.T, err)
// 	return token
// }

// func (f *Fixtures) QueryTokenCollectionExpectEmpty(contractID, tokenID string, flags ...string) {
// 	cmd := fmt.Sprintf("%s query collection token %s %s %s", f.LinkcliBinary, contractID, tokenID, f.Flags())
// 	cmd = addFlags(cmd, flags)
// 	_, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.NotEmpty(f.T, errStr)
// }

// func (f *Fixtures) QueryTokensCollection(contractID string) collectionModule.Tokens {
// 	cmd := fmt.Sprintf("%s query collection tokens %s %s", f.LinkcliBinary, contractID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var tokens collectionModule.Tokens
// 	err := cdc.UnmarshalJSON([]byte(res), &tokens)
// 	require.NoError(f.T, err)
// 	return tokens
// }

// func (f *Fixtures) QueryTokensByTokenTypeCollection(contractID string, tokenType string) collectionModule.Tokens {
// 	cmd := fmt.Sprintf("%s query collection tokens %s --token-type %s %s", f.LinkcliBinary, contractID, tokenType, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var tokens collectionModule.Tokens
// 	err := cdc.UnmarshalJSON([]byte(res), &tokens)
// 	require.NoError(f.T, err)
// 	return tokens
// }

// func (f *Fixtures) QueryTokenTypeCollection(contractID, tokenTypeID string, flags ...string) collectionModule.TokenType {
// 	cmd := fmt.Sprintf("%s query collection tokentype %s %s %s", f.LinkcliBinary, contractID, tokenTypeID, f.Flags())
// 	cmd = addFlags(cmd, flags)
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var tokenType collectionModule.TokenType
// 	err := cdc.UnmarshalJSON([]byte(res), &tokenType)
// 	require.NoError(f.T, err)
// 	return tokenType
// }

// func (f *Fixtures) QueryCollection(contractID string, flags ...string) collectionModule.Collection {
// 	cmd := fmt.Sprintf("%s query collection collection %s %s", f.LinkcliBinary, contractID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var collection collectionModule.Collection
// 	err := cdc.UnmarshalJSON([]byte(res), &collection)
// 	require.NoError(f.T, err)

// 	return collection
// }

// func (f *Fixtures) QueryTotalSupplyTokenCollection(contractID, tokenID string, flags ...string) sdk.Int {
// 	cmd := fmt.Sprintf("%s query collection total supply %s %s %s", f.LinkcliBinary, contractID, tokenID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var supply sdk.Int
// 	err := cdc.UnmarshalJSON([]byte(res), &supply)
// 	require.NoError(f.T, err)

// 	return supply
// }
// func (f *Fixtures) QueryTotalMintTokenCollection(contractID, tokenID string, flags ...string) sdk.Int {
// 	cmd := fmt.Sprintf("%s query collection total mint %s %s %s", f.LinkcliBinary, contractID, tokenID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var supply sdk.Int
// 	err := cdc.UnmarshalJSON([]byte(res), &supply)
// 	require.NoError(f.T, err)

// 	return supply
// }
// func (f *Fixtures) QueryTotalBurnTokenCollection(contractID, tokenID string, flags ...string) sdk.Int {
// 	cmd := fmt.Sprintf("%s query collection total burn %s %s %s", f.LinkcliBinary, contractID, tokenID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var supply sdk.Int
// 	err := cdc.UnmarshalJSON([]byte(res), &supply)
// 	require.NoError(f.T, err)

// 	return supply
// }
// func (f *Fixtures) QueryCountTokenCollection(contractID, tokenID string, flags ...string) sdk.Int {
// 	cmd := fmt.Sprintf("%s query collection count total %s %s %s", f.LinkcliBinary, contractID, tokenID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var supply sdk.Int
// 	err := cdc.UnmarshalJSON([]byte(res), &supply)
// 	require.NoError(f.T, err)

// 	return supply
// }

// func (f *Fixtures) QueryAccountPermissionCollection(addr sdk.AccAddress, contractID string, flags ...string) collectionModule.Permissions {
// 	cmd := fmt.Sprintf("%s query collection perm %s %s %s", f.LinkcliBinary, addr, contractID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var pms collectionModule.Permissions
// 	err := cdc.UnmarshalJSON([]byte(res), &pms)
// 	require.NoError(f.T, err)
// 	return pms
// }

// func (f *Fixtures) QueryRootTokenCollection(contractID string, tokenID string, flags ...string) collectionModule.Token {
// 	cmd := fmt.Sprintf("%s query collection root %s %s %s", f.LinkcliBinary, contractID, tokenID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var root collectionModule.Token
// 	err := cdc.UnmarshalJSON([]byte(res), &root)
// 	require.NoError(f.T, err)
// 	return root
// }

// func (f *Fixtures) QueryParentTokenCollection(contractID string, tokenID string, flags ...string) collectionModule.Token {
// 	cmd := fmt.Sprintf("%s query collection parent %s %s %s", f.LinkcliBinary, contractID, tokenID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var parent collectionModule.Token
// 	err := cdc.UnmarshalJSON([]byte(res), &parent)
// 	require.NoError(f.T, err)
// 	return parent
// }

// func (f *Fixtures) QueryChildrenTokenCollection(contractID string, tokenID string, flags ...string) []collectionModule.Token {
// 	cmd := fmt.Sprintf("%s query collection children %s %s %s", f.LinkcliBinary, contractID, tokenID, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var children []collectionModule.Token
// 	err := cdc.UnmarshalJSON([]byte(res), &children)
// 	require.NoError(f.T, err)
// 	return children
// }

// func (f *Fixtures) QueryApproversTokenCollection(contractID string, proxy sdk.AccAddress) []sdk.AccAddress {
// 	cmd := fmt.Sprintf("%s query collection approvers %s %s %s", f.LinkcliBinary, contractID, proxy, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var approvers []sdk.AccAddress
// 	err := cdc.UnmarshalJSON([]byte(res), &approvers)
// 	require.NoError(f.T, err)
// 	return approvers
// }

// func (f *Fixtures) QueryApprovedTokenCollection(contractID string, proxy sdk.AccAddress, approver sdk.AccAddress) bool {
// 	cmd := fmt.Sprintf("%s query collection approved %s %s %s %s", f.LinkcliBinary, contractID, proxy, approver, f.Flags())
// 	res, errStr := tests.ExecuteT(f.T, cmd, "")
// 	require.Empty(f.T, errStr)
// 	cdc := app.MakeCodec()
// 	var isApproved bool
// 	err := cdc.UnmarshalJSON([]byte(res), &isApproved)
// 	require.NoError(f.T, err)
// 	return isApproved
// }

// ___________________________________________________________________________________
// wasm
func (f *Fixtures) TxStoreWasm(wasmFilePath string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx wasm store %s %v", f.LinkcliBinary, wasmFilePath, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxInstantiateWasm(codeId uint64, msgJson string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx wasm instantiate %d %s %v", f.LinkcliBinary, codeId, msgJson, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxExecuteWasm(contractAddress sdk.AccAddress, msgJson string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx wasm execute %s %s %s", f.LinkcliBinary, contractAddress, msgJson, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxMigrateWasm(contractAddress sdk.AccAddress, codeId uint64, msgJson string, flags ...string) (bool, string, string) {
	require.Fail(f.T, "TODO: Test It!")
	cmd := fmt.Sprintf("%s tx wasm migrate %s %d %s %s", f.LinkcliBinary, contractAddress, codeId, msgJson, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxSetContractAdminWasm(contractAddress sdk.AccAddress, newAdmin sdk.AccAddress, flags ...string) (bool, string, string) {
	require.Fail(f.T, "TODO: Test It!")
	cmd := fmt.Sprintf("%s tx wasm set-contract-admin %s %s %s", f.LinkcliBinary, contractAddress, newAdmin, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) TxClearContractAdminWasm(contractAddress sdk.AccAddress, flags ...string) (bool, string, string) {
	require.Fail(f.T, "TODO: Test It!")
	cmd := fmt.Sprintf("%s tx wasm clear-contract-admin %s %s", f.LinkcliBinary, contractAddress, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags))
}

func (f *Fixtures) QueryListCodeWasm() []wasmtypes.CodeInfoResponse {
	cmd := fmt.Sprintf("%s query wasm list-code %s", f.LinkcliBinary, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var listCode []wasmtypes.CodeInfoResponse
	err := cdc.UnmarshalJSON([]byte(res), &listCode)
	require.NoError(f.T, err)
	return listCode
}

func (f *Fixtures) QueryListContractByCodeWasm(codeId uint64) []wasmtypes.ContractInfoResponse {
	cmd := fmt.Sprintf("%s query wasm list-contract-by-code %d %s", f.LinkcliBinary, codeId, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var listContract []wasmtypes.ContractInfoResponse
	err := cdc.UnmarshalJSON([]byte(res), &listContract)
	require.NoError(f.T, err)
	return listContract
}

func (f *Fixtures) QueryCodeWasm(codeId uint64, outputPath string) {
	cmd := fmt.Sprintf("%s query wasm code %d %s %s", f.LinkcliBinary, codeId, outputPath, f.Flags())
	_, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
}

func (f *Fixtures) QueryContractWasm(contractAddress sdk.AccAddress) wasmtypes.ContractInfoResponse {
	cmd := fmt.Sprintf("%s query wasm contract %s %s", f.LinkcliBinary, contractAddress, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var contractInfo wasmtypes.ContractInfoResponse
	err := cdc.UnmarshalJSON([]byte(res), &contractInfo)
	require.NoError(f.T, err)
	return contractInfo
}

func (f *Fixtures) QueryContractHistoryWasm(contractAddress sdk.AccAddress) wasmtypes.ContractCodeHistoryEntry {
	require.Fail(f.T, "TODO: Test It!")
	cmd := fmt.Sprintf("%s query wasm contract-history %s %s", f.LinkcliBinary, contractAddress, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var contractHistory wasmtypes.ContractCodeHistoryEntry
	err := cdc.UnmarshalJSON([]byte(res), &contractHistory)
	require.NoError(f.T, err)
	return contractHistory
}

func (f *Fixtures) QueryContractStateAllWasm(contractAddress sdk.AccAddress) []wasmtypes.Model {
	require.Fail(f.T, "TODO: Test It!")
	cmd := fmt.Sprintf("%s query wasm contract-state all %s %s", f.LinkcliBinary, contractAddress, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var state []wasmtypes.Model
	err := cdc.UnmarshalJSON([]byte(res), &state)
	require.NoError(f.T, err)
	return state
}

func (f *Fixtures) QueryContractStateRawWasm(contractAddress sdk.AccAddress, key string) []wasmtypes.Model {
	require.Fail(f.T, "TODO: Test It!")
	cmd := fmt.Sprintf("%s query wasm contract-state raw %s %s %s", f.LinkcliBinary, contractAddress, key, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var state []wasmtypes.Model
	err := cdc.UnmarshalJSON([]byte(res), &state)
	require.NoError(f.T, err)
	return state
}

func (f *Fixtures) QueryContractStateSmartWasm(contractAddress sdk.AccAddress, reqJson string) string {
	cmd := fmt.Sprintf("%s query wasm contract-state smart %s %s %s", f.LinkcliBinary, contractAddress, reqJson, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	return res
}

// ___________________________________________________________________________________
// tendermint rpc
func (f *Fixtures) NetInfo(flags ...string) *tmctypes.ResultNetInfo {
	tmc, err := tmhttp.New(fmt.Sprintf("tcp://0.0.0.0:%s", f.Port), "/websocket")
	if err != nil {
		panic(fmt.Sprintf("failed to create Tendermint HTTP client: %s", err))
	}

	err = tmc.Start()
	require.NoError(f.T, err)
	defer func() {
		err := tmc.Stop()
		require.NoError(f.T, err)
	}()

	netInfo, err := tmc.NetInfo()
	require.NoError(f.T, err)
	return netInfo
}

func (f *Fixtures) Status(flags ...string) *tmctypes.ResultStatus {
	tmc, err := tmhttp.New(fmt.Sprintf("tcp://0.0.0.0:%s", f.Port), "/websocket")
	if err != nil {
		panic(fmt.Sprintf("failed to create Tendermint HTTP client: %s", err))
	}

	err = tmc.Start()
	require.NoError(f.T, err)
	defer func() {
		err := tmc.Stop()
		require.NoError(f.T, err)
	}()

	netInfo, err := tmc.Status()
	require.NoError(f.T, err)
	return netInfo
}

// ___________________________________________________________________________________
// executors

func executeWriteCheckErr(t *testing.T, cmdStr string, writes ...string) {
	require.True(t, executeWrite(t, cmdStr, writes...))
}

func executeWrite(t *testing.T, cmdStr string, writes ...string) (exitSuccess bool) {
	exitSuccess, _, _ = executeWriteRetStdStreams(t, cmdStr, writes...)
	return
}

func executeWriteRetStdStreams(t *testing.T, cmdStr string, writes ...string) (bool, string, string) {
	proc := tests.GoExecuteT(t, cmdStr)

	// Enables use of interactive commands
	for _, write := range writes {
		_, err := proc.StdinPipe.Write([]byte(write + "\n"))
		require.NoError(t, err)
	}

	// Read both stdout and stderr from the process
	stdout, stderr, err := proc.ReadAll()
	if err != nil {
		fmt.Println("Err on proc.ReadAll()", err, cmdStr)
	}

	// Log output.
	if len(stdout) > 0 {
		t.Log("Stdout:", string(stdout))
	}
	if len(stderr) > 0 {
		t.Log("Stderr:", string(stderr))
	}

	// Wait for process to exit
	proc.Wait()

	// Return succes, stdout, stderr
	return proc.ExitState.Success(), string(stdout), string(stderr)
}

// ___________________________________________________________________________________
// utils

func addFlags(cmd string, flags []string) string {
	for _, f := range flags {
		cmd += " " + f
	}
	return strings.TrimSpace(cmd)
}

// Write the given string to a new temporary file
func WriteToNewTempFile(t *testing.T, s string) *os.File {
	fp, err := ioutil.TempFile(os.TempDir(), "cosmos_cli_test_")
	require.Nil(t, err)
	_, err = fp.WriteString(s)
	require.Nil(t, err)
	return fp
}

func MarshalStdTx(t *testing.T, stdTx auth.StdTx) []byte {
	cdc := app.MakeCodec()
	bz, err := cdc.MarshalBinaryBare(stdTx)
	require.NoError(t, err)
	return bz
}

func UnmarshalStdTx(t *testing.T, s string) (stdTx auth.StdTx) {
	cdc := app.MakeCodec()
	require.Nil(t, cdc.UnmarshalJSON([]byte(s), &stdTx))
	return
}

func UnmarshalTxResponse(t *testing.T, s string) (txResp sdk.TxResponse) {
	cdc := app.MakeCodec()
	require.Nil(t, cdc.UnmarshalJSON([]byte(s), &txResp))
	return
}

// ___________________________________________________________________________________
// Fixture Group

type FixtureGroup struct {
	T                  *testing.T
	DockerImage        string
	fixturesMap        map[string]*Fixtures
	networkName        string
	subnet             string
	genesisFileContent []byte
}

func NewFixtureGroup(t *testing.T) *FixtureGroup {
	fg := &FixtureGroup{
		T:           t,
		DockerImage: "line/link",
		fixturesMap: make(map[string]*Fixtures),
	}

	fg.networkName = networkNamePrefix + t.Name()

	return fg
}

func InitFixturesGroup(t *testing.T, subnet string, numOfNodes ...int) *FixtureGroup {
	nodeNumber := 4
	if len(numOfNodes) == 1 {
		nodeNumber = numOfNodes[0]
	}
	fg := NewFixtureGroup(t)
	fg.initNodes(subnet, nodeNumber)
	fg.createNetwork()
	return fg
}

func (fg *FixtureGroup) createNetwork() {
	cmd := fmt.Sprintf("docker network create %s --subnet %s/24", fg.networkName, fg.subnet)
	_, _ = tests.ExecuteT(fg.T, cmd, "")
}

func (fg *FixtureGroup) rmNetwork() {
	cmd := fmt.Sprintf("docker network rm %s", fg.networkName)
	_, _ = tests.ExecuteT(fg.T, cmd, "")
}

func (fg *FixtureGroup) initNodes(subnet string, numberOfNodes int) {
	t := fg.T
	chainID := t.Name()

	fg.subnet = subnet

	stdout, _ := tests.ExecuteT(t, fmt.Sprintf("docker images %s -q", fg.DockerImage), "")
	if len(stdout) == 0 {
		panic(errors.New("docker image is not found"))
	}

	// Initialize temporary directories
	gentxDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.RemoveAll(gentxDir))
	}()

	for idx := 0; idx < numberOfNodes; idx++ {
		name := fmt.Sprintf("%s-%s%d", t.Name(), namePrefix, idx)
		f := NewFixtures(t)
		f.UnsafeResetAll()
		f.LDInit(name, fmt.Sprintf("--chain-id %s", chainID))
		fg.fixturesMap[name] = f
		ip, err := calculateIP(subnet, idx+2)
		require.NoError(fg.T, err)
		f.BridgeIP = ip
	}
	for name, f := range fg.fixturesMap {
		f.KeysDelete(name)
		f.KeysAdd(name)
		f.CLIConfig("output", "json")
		f.CLIConfig("chain-id", f.ChainID)
		f.CLIConfig("broadcast-mode", "block")
		f.CLIConfig("trust-node", "true")
	}

	for _, f := range fg.fixturesMap {
		for nameInner, fInner := range fg.fixturesMap {
			f.AddGenesisAccount(fInner.KeyAddress(nameInner), startCoins)
		}
	}

	for name, f := range fg.fixturesMap {
		gentxDoc := filepath.Join(gentxDir, fmt.Sprintf("%s.json", name))
		f.GenTx(name, fmt.Sprintf("--output-document=%s --ip=%s", gentxDoc, f.BridgeIP))
	}

	for _, f := range fg.fixturesMap {
		f.CollectGenTxs(fmt.Sprintf("--gentx-dir=%s", gentxDir))
		f.ValidateGenesis()
		if len(fg.genesisFileContent) == 0 {
			fg.genesisFileContent, err = ioutil.ReadFile(f.GenesisFile())
			require.NoError(t, err)
		}
	}
	for _, f := range fg.fixturesMap {
		err := ioutil.WriteFile(f.GenesisFile(), fg.genesisFileContent, os.ModePerm)
		require.NoError(t, err)
	}
}

func (fg *FixtureGroup) LDStopContainers() {
	wg := &sync.WaitGroup{}
	wg.Add(len(fg.fixturesMap))

	for _, f := range fg.fixturesMap {
		copyedF := f
		go func() {
			fg.LDStopContainer(copyedF)
			wg.Done()
		}()
	}

	if timeout := waitTimeout(wg, time.Minute); timeout {
		panic(errors.New("linkd stop containers failed"))
	}
}

func (fg *FixtureGroup) LDStartContainers() {
	wg := &sync.WaitGroup{}
	wg.Add(len(fg.fixturesMap))

	for _, f := range fg.fixturesMap {
		_ = fg.LDStartContainer(f)
	}

	for _, f := range fg.fixturesMap {
		port := f.Port
		go func() {
			WaitForTMStart(port)
			tests.WaitForNextNBlocksTM(1, port)
			wg.Done()
		}()
	}
	if timeout := waitTimeout(wg, time.Minute); timeout {
		panic(errors.New("linkd start containers failed"))
	}
}

func (fg *FixtureGroup) LDStartContainer(f *Fixtures, flags ...string) *tests.Process {
	dockerCommand := "docker run --rm --name %s --network %s --ip %s -p %s:26656 -p %s:26657 -v %s:/root/.linkd:Z -v %s:/root/.linkcli:Z line/link linkd start --rpc.laddr=tcp://0.0.0.0:26657 --p2p.laddr=tcp://0.0.0.0:26656"
	dockerCommand = fmt.Sprintf(dockerCommand, f.Moniker, fg.networkName, f.BridgeIP, f.P2PPort, f.Port, f.LinkdHome, f.LinkcliHome)
	fg.T.Log(dockerCommand)
	proc := tests.GoExecuteTWithStdout(f.T, addFlags(dockerCommand, flags))

	return proc
}

func (fg *FixtureGroup) LDStopContainer(f *Fixtures) {
	cmd := "docker ps --filter name=%s --filter status=running -q"
	cmd = fmt.Sprintf(cmd, f.Moniker)
	stdout, _ := tests.ExecuteT(f.T, cmd, "")
	containerID := stdout
	if len(stdout) > 0 {
		cmd := "docker stop %s"
		cmd = fmt.Sprintf(cmd, containerID)
		stdout, stderr := tests.ExecuteT(f.T, cmd, "")
		if stdout != containerID {
			panic(stderr)
		}
	}
}

func (fg *FixtureGroup) WaitForContainer(f *Fixtures) {
	cmd := "docker ps --filter name=%s --filter status=running -q"
	cmd = fmt.Sprintf(cmd, f.Moniker)

	var err error
	fmt.Printf("Wait for the container[%s] boot up\n", f.Moniker)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 100)
		stdout, stderr := tests.ExecuteT(f.T, cmd, "")

		if len(stdout) > 0 {
			return
		}
		err = errors.New(stderr)
	}
	// still haven't started up?! panic!
	panic(err)
}

func (fg *FixtureGroup) AddFullNode(flags ...string) *Fixtures {
	t := fg.T
	idx := len(fg.fixturesMap)
	chainID := fg.T.Name()

	name := fmt.Sprintf("%s-%s%d", t.Name(), namePrefix, idx)
	f := NewFixtures(t)

	// Initialize linkd
	{
		f.UnsafeResetAll()
		f.LDInit(name, fmt.Sprintf("--chain-id %s", chainID))
		ip, err := calculateIP(fg.subnet, idx+2)
		require.NoError(fg.T, err)
		f.BridgeIP = ip
	}

	// Initialize linkcli
	{
		f.KeysDelete(name)
		f.KeysAdd(name)
		f.CLIConfig("output", "json")
		f.CLIConfig("chain-id", f.ChainID)
		f.CLIConfig("broadcast-mode", "block")
		f.CLIConfig("trust-node", "true")
	}

	// Copy the genesis.json
	{
		if len(fg.genesisFileContent) == 0 {
			panic("genesis file is not loaded")
		}
		err := ioutil.WriteFile(f.GenesisFile(), fg.genesisFileContent, os.ModePerm)
		require.NoError(t, err)
	}

	// Configure for invisible options
	{
		if len(flags) > 0 {
			configFilePath := filepath.Join(f.LinkdHome, "config/config.toml")

			conf := cfg.DefaultConfig()
			err := viper.Unmarshal(conf)
			require.NoError(t, err)

			for _, flag := range flags {
				if flag == "--mempool.broadcast=false" {
					conf.Mempool.Broadcast = false
				}
			}

			cfg.WriteConfigFile(configFilePath, conf)
		}
	}

	// Collect the persistent peers from the network
	var persistentPeers []string
	{
		for _, other := range fg.fixturesMap {
			statusInfo := other.Status()
			id := string(statusInfo.NodeInfo.ID())
			persistentPeers = append(persistentPeers, fmt.Sprintf("%s@%s:%d", id, other.BridgeIP, 26656))
		}
	}

	// Add fixture to the group
	fg.fixturesMap[name] = f

	// Start linkd
	fg.LDStartContainer(f, fmt.Sprintf("--p2p.persistent_peers %s", strings.Join(persistentPeers, ",")))

	WaitForTMStart(f.Port)
	tests.WaitForNextNBlocksTM(1, f.Port)

	return f
}

func (fg *FixtureGroup) Fixture(index int) *Fixtures {
	name := fmt.Sprintf("%s-%s%d", fg.T.Name(), namePrefix, index)
	if f, ok := fg.fixturesMap[name]; ok {
		return f
	}
	return nil
}

func (fg *FixtureGroup) Cleanup() {
	fg.LDStopContainers()
	fg.rmNetwork()
	for _, f := range fg.fixturesMap {
		f.Cleanup()
	}
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

func calculateIP(ip string, i int) (string, error) {
	ipv4 := net.ParseIP(ip).To4()
	if ipv4 == nil {
		return "", fmt.Errorf("%v: non ipv4 address", ip)
	}

	for j := 0; j < i; j++ {
		ipv4[3]++
	}

	return ipv4.String(), nil
}

// wait for tendermint to start by querying tendermint
func WaitForTMStart(port string) {
	url := fmt.Sprintf("http://localhost:%v/block", port)
	WaitForStart(url)
}

// WaitForStart waits for the node to start by pinging the url
// every 100ms for 10s until it returns 200. If it takes longer than 5s,
// it panics.
func WaitForStart(url string) {
	var err error

	// ping the status endpoint a few times a second
	// for a few seconds until we get a good response.
	// otherwise something probably went wrong
	// 2 ^ 7 = 128 --> approximately 10 secs
	wait := 1
	for i := 0; i < 7; i++ {
		// 0.1, 0.2, 0.4, 0.8, 1.6, 3.2, 6.4, 12.8, 25.6, 51.2, 102.4
		time.Sleep(time.Millisecond * 100 * time.Duration(wait))
		wait *= 2

		var res *http.Response
		/* #nosec */
		res, err = http.Get(url) // Error is arising in testing files
		if err != nil || res == nil {
			continue
		}
		err = res.Body.Close()
		if err != nil {
			panic(err)
		}

		if res.StatusCode == http.StatusOK {
			// good!
			return
		}
	}
	// still haven't started up?! panic!
	panic(err)
}
