package clitest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/spf13/viper"

	"github.com/link-chain/link/types"

	safetyBoxModule "github.com/link-chain/link/x/safetybox"
	tokenModule "github.com/link-chain/link/x/token"

	"github.com/link-chain/link/client"
	tmclient "github.com/tendermint/tendermint/rpc/client"

	"github.com/stretchr/testify/require"

	cfg "github.com/tendermint/tendermint/config"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/link-chain/link/app"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	atypes "github.com/link-chain/link/x/auth/client/types"

	"github.com/link-chain/link/x/lrc3"

	nft "github.com/link-chain/link/x/nft"
)

const (
	denom        = "stake"
	keyFoo       = "foo"
	keyBar       = "bar"
	fooDenom     = "footoken"
	feeDenom     = "feetoken"
	fee2Denom    = "fee2token"
	keyBaz       = "baz"
	keyVesting   = "vesting"
	keyFooBarBaz = "foobarbaz"

	denomStake = "stake2"
	denomLink  = "link"
	userTina   = "tina"
	userKevin  = "kevin"
	userRinah  = "rinah"
	userBrian  = "brian"
	userEvelyn = "evelyn"
	userSam    = "sam"
)

const (
	namePrefix        = "node"
	networkNamePrefix = "line-linkdnode-testnet-"
)

var (
	totalCoins = sdk.NewCoins(
		sdk.NewCoin(denomLink, sdk.TokensFromConsensusPower(6000)),
		sdk.NewCoin(denomStake, sdk.TokensFromConsensusPower(600000000)),
		sdk.NewCoin(fee2Denom, sdk.TokensFromConsensusPower(2000000)),
		sdk.NewCoin(feeDenom, sdk.TokensFromConsensusPower(2000000)),
		sdk.NewCoin(fooDenom, sdk.TokensFromConsensusPower(2000)),
		sdk.NewCoin(denom, sdk.TokensFromConsensusPower(300)), // We don't use inflation
		//sdk.NewCoin(denom, sdk.TokensFromConsensusPower(300).Add(sdk.NewInt(12))), // add coins from inflation
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
		sdk.NewCoin(denomLink, sdk.TokensFromConsensusPower(1000)),
		sdk.NewCoin(denomStake, sdk.TokensFromConsensusPower(100000000)),
	)
)

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(types.Bech32PrefixAccAddr, types.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr, types.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr, types.Bech32PrefixConsPub)
	config.SetCoinType(types.CoinType)
	config.SetFullFundraiserPath(types.FullFundraiserPath)
	config.Seal()
}

//___________________________________________________________________________________
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
	tmpDir := path.Join(os.ExpandEnv("$HOME"), ".linktest")
	err := os.MkdirAll(tmpDir, os.ModePerm)
	require.NoError(t, err)
	tmpDir, err = ioutil.TempDir(tmpDir, "link_integration_"+t.Name()+"_")
	require.NoError(t, err)

	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)

	p2pAddr, p2pPort, err := server.FreeTCPAddr()
	require.NoError(t, err)

	buildDir := os.Getenv("BUILDDIR")
	if buildDir == "" {
		buildDir, err = filepath.Abs("../build/")
		require.NoError(t, err)
	}

	return &Fixtures{
		T:             t,
		BuildDir:      buildDir,
		RootDir:       tmpDir,
		LinkdBinary:   filepath.Join(buildDir, "linkd"),
		LinkcliBinary: filepath.Join(buildDir, "linkcli"),
		LinkdHome:     filepath.Join(tmpDir, ".linkd"),
		LinkcliHome:   filepath.Join(tmpDir, ".linkcli"),
		RPCAddr:       servAddr,
		P2PAddr:       p2pAddr,
		Port:          port,
		P2PPort:       p2pPort,
		Moniker:       "", // initialized by LDInit
		BridgeIP:      "",
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
	genDoc, err := tmtypes.GenesisDocFromFile(f.GenesisFile())
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
	f.KeysDelete(userTina)
	f.KeysDelete(userKevin)
	f.KeysDelete(userRinah)
	f.KeysDelete(userBrian)
	f.KeysDelete(userEvelyn)
	f.KeysDelete(userSam)
	f.KeysAdd(userTina)
	f.KeysAdd(userKevin)
	f.KeysAdd(userRinah)
	f.KeysAdd(userBrian)
	f.KeysAdd(userEvelyn)
	f.KeysAdd(userSam)

	// ensure that CLI output is in JSON format
	f.CLIConfig("output", "json")

	// NOTE: LDInit sets the ChainID
	f.LDInit(keyFoo)

	f.CLIConfig("chain-id", f.ChainID)
	f.CLIConfig("broadcast-mode", "block")

	// start an account with tokens
	f.AddGenesisAccount(f.KeyAddress(keyFoo), startCoins)
	//f.AddGenesisAccount(f.KeyAddress(keyBar), startCoins)
	f.AddGenesisAccount(
		f.KeyAddress(keyVesting), startCoins,
		fmt.Sprintf("--vesting-amount=%s", vestingCoins),
		fmt.Sprintf("--vesting-start-time=%d", time.Now().UTC().UnixNano()),
		fmt.Sprintf("--vesting-end-time=%d", time.Now().Add(60*time.Second).UTC().UnixNano()),
	)

	// add genesis accounts for testing
	f.AddGenesisAccount(f.KeyAddress(userTina), defaultCoins)
	f.AddGenesisAccount(f.KeyAddress(userKevin), defaultCoins)
	f.AddGenesisAccount(f.KeyAddress(userRinah), defaultCoins)
	f.AddGenesisAccount(f.KeyAddress(userBrian), defaultCoins)
	f.AddGenesisAccount(f.KeyAddress(userEvelyn), defaultCoins)
	f.AddGenesisAccount(f.KeyAddress(userSam), defaultCoins)

	f.GenTx(keyFoo)
	f.CollectGenTxs()

	return
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

//___________________________________________________________________________________
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
	_, stderr := tests.ExecuteT(f.T, addFlags(cmd, flags), client.DefaultKeyPass)

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
	cmd := fmt.Sprintf("%s add-genesis-account %s %s --home=%s", f.LinkdBinary, address, coins, f.LinkdHome)
	executeWriteCheckErr(f.T, addFlags(cmd, flags))
}

// GenTx is linkd gentx
func (f *Fixtures) GenTx(name string, flags ...string) {
	cmd := fmt.Sprintf("%s gentx --name=%s --home=%s --home-client=%s", f.LinkdBinary, name, f.LinkdHome, f.LinkcliHome)
	executeWriteCheckErr(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// CollectGenTxs is linkd collect-gentxs
func (f *Fixtures) CollectGenTxs(flags ...string) {
	cmd := fmt.Sprintf("%s collect-gentxs --home=%s", f.LinkdBinary, f.LinkdHome)
	executeWriteCheckErr(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// LDStart runs linkd start with the appropriate flags and returns a process
func (f *Fixtures) LDStart(flags ...string) *tests.Process {
	cmd := fmt.Sprintf("%s start --home=%s --rpc.laddr=%v --p2p.laddr=%v", f.LinkdBinary, f.LinkdHome, f.RPCAddr, f.P2PAddr)
	proc := tests.GoExecuteTWithStdout(f.T, addFlags(cmd, flags))
	tests.WaitForTMStart(f.Port)
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

//___________________________________________________________________________________
// linkcli keys

// KeysDelete is linkcli keys delete
func (f *Fixtures) KeysDelete(name string, flags ...string) {
	cmd := fmt.Sprintf("%s keys delete --home=%s %s", f.LinkcliBinary, f.LinkcliHome, name)
	executeWrite(f.T, addFlags(cmd, append(append(flags, "-y"), "-f")))
}

// KeysAdd is linkcli keys add
func (f *Fixtures) KeysAdd(name string, flags ...string) {
	cmd := fmt.Sprintf("%s keys add --home=%s %s", f.LinkcliBinary, f.LinkcliHome, name)
	executeWriteCheckErr(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// KeysAddRecover prepares linkcli keys add --recover
func (f *Fixtures) KeysAddRecover(name, mnemonic string, flags ...string) (exitSuccess bool, stdout, stderr string) {
	cmd := fmt.Sprintf("%s keys add --home=%s --recover %s", f.LinkcliBinary, f.LinkcliHome, name)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass, mnemonic)
}

// KeysAddRecoverHDPath prepares linkcli keys add --recover --account --index
func (f *Fixtures) KeysAddRecoverHDPath(name, mnemonic string, account uint32, index uint32, flags ...string) {
	cmd := fmt.Sprintf("%s keys add --home=%s --recover %s --account %d --index %d", f.LinkcliBinary, f.LinkcliHome, name, account, index)
	executeWriteCheckErr(f.T, addFlags(cmd, flags), client.DefaultKeyPass, mnemonic)
}

// KeysShow is linkcli keys show
func (f *Fixtures) KeysShow(name string, flags ...string) keys.KeyOutput {
	cmd := fmt.Sprintf("%s keys show --home=%s %s", f.LinkcliBinary, f.LinkcliHome, name)
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var ko keys.KeyOutput
	err := client.UnmarshalJSON([]byte(out), &ko)
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

//___________________________________________________________________________________
// linkcli config

// CLIConfig is linkcli config
func (f *Fixtures) CLIConfig(key, value string, flags ...string) {
	cmd := fmt.Sprintf("%s config --home=%s %s %s", f.LinkcliBinary, f.LinkcliHome, key, value)
	executeWriteCheckErr(f.T, addFlags(cmd, flags))
}

//___________________________________________________________________________________
// linkcli tx send/sign/broadcast

// TxSend is linkcli tx send
func (f *Fixtures) TxSend(from string, to sdk.AccAddress, amount sdk.Coin, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx send %s %s %s %v", f.LinkcliBinary, from, to, amount, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxSign is linkcli tx sign
func (f *Fixtures) TxSign(signer, fileName string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx sign %v --from=%s %v", f.LinkcliBinary, f.Flags(), signer, fileName)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxBroadcast is linkcli tx broadcast
func (f *Fixtures) TxBroadcast(fileName string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx broadcast %v %v", f.LinkcliBinary, f.Flags(), fileName)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxEncode is linkcli tx encode
func (f *Fixtures) TxEncode(fileName string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx encode %v %v", f.LinkcliBinary, f.Flags(), fileName)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxMultisign is linkcli tx multisign
func (f *Fixtures) TxMultisign(fileName, name string, signaturesFiles []string,
	flags ...string) (bool, string, string) {

	cmd := fmt.Sprintf("%s tx multisign %v %s %s %s", f.LinkcliBinary, f.Flags(),
		fileName, name, strings.Join(signaturesFiles, " "),
	)
	return executeWriteRetStdStreams(f.T, cmd)
}

// TxLRC3Generate is linkcli tx lrc3
func (f *Fixtures) TxLRC3Init(denom string, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx lrc3 init %s --from %s %v", f.LinkcliBinary, denom, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxLRC3Mint is linkcli tx lrc3
func (f *Fixtures) TxLRC3Mint(denom, to, from, testTokenURI1 string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx lrc3 mint %s %s %s --from %s %v", f.LinkcliBinary, denom, testTokenURI1, to, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxLRC3Burn is linkcli tx lrc3
func (f *Fixtures) TxLRC3Burn(denom string, tokenId int, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx lrc3 burn %s %d --from %s %v", f.LinkcliBinary, denom, tokenId, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxLRC3Approve is linkcli tx lrc3
func (f *Fixtures) TxLRC3Approve(denom string, tokenId int, to string, sender string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx lrc3 approve %s %d %s --from %s %v", f.LinkcliBinary, denom, tokenId, to, sender, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxLRC3Transfer is linkcli tx lrc3
func (f *Fixtures) TxLRC3Transfer(denom string, tokenId int, to string, from string, sender string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx lrc3 transfer %s %s %s %d --from %s %v", f.LinkcliBinary, from, to, denom, tokenId, sender, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxLRC3SetApprovalForAll is linkcli tx lrc3
func (f *Fixtures) TxLRC3SetApprovalForAll(denom string, operator string, approved bool, sender string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx lrc3 setApprovalForAll %s %s %t --from %s %v", f.LinkcliBinary, denom, operator, approved, sender, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

//___________________________________________________________________________________
// linkcli tx staking

// TxStakingCreateValidator is linkcli tx staking create-validator
func (f *Fixtures) TxStakingCreateValidator(from, consPubKey string, amount sdk.Coin, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx staking create-validator %v --from=%s --pubkey=%s", f.LinkcliBinary, f.Flags(), from, consPubKey)
	cmd += fmt.Sprintf(" --amount=%v --moniker=%v --commission-rate=%v", amount, from, "0.05")
	cmd += fmt.Sprintf(" --commission-max-rate=%v --commission-max-change-rate=%v", "0.20", "0.10")
	cmd += fmt.Sprintf(" --min-self-delegation=%v", "1")
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxStakingUnbond is linkcli tx staking unbond
func (f *Fixtures) TxStakingUnbond(from, shares string, validator sdk.ValAddress, flags ...string) bool {
	cmd := fmt.Sprintf("%s tx staking unbond %s %v --from=%s %v", f.LinkcliBinary, validator, shares, from, f.Flags())
	return executeWrite(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

//___________________________________________________________________________________
// linkcli tx gov

// TxGovSubmitProposal is linkcli tx gov submit-proposal
func (f *Fixtures) TxGovSubmitProposal(from, typ, title, description string, deposit sdk.Coin, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx gov submit-proposal %v --from=%s --type=%s", f.LinkcliBinary, f.Flags(), from, typ)
	cmd += fmt.Sprintf(" --title=%s --description=%s --deposit=%s", title, description, deposit)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxGovDeposit is linkcli tx gov deposit
func (f *Fixtures) TxGovDeposit(proposalID int, from string, amount sdk.Coin, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx gov deposit %d %s --from=%s %v", f.LinkcliBinary, proposalID, amount, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxGovVote is linkcli tx gov vote
func (f *Fixtures) TxGovVote(proposalID int, option gov.VoteOption, from string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx gov vote %d %s --from=%s %v", f.LinkcliBinary, proposalID, option, from, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxGovSubmitParamChangeProposal executes a CLI parameter change proposal
// submission.
func (f *Fixtures) TxGovSubmitParamChangeProposal(
	from, proposalPath string, deposit sdk.Coin, flags ...string,
) (bool, string, string) {

	cmd := fmt.Sprintf(
		"%s tx gov submit-proposal param-change %s --from=%s %v",
		f.LinkcliBinary, proposalPath, from, f.Flags(),
	)

	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxGovSubmitCommunityPoolSpendProposal executes a CLI community pool spend proposal
// submission.
func (f *Fixtures) TxGovSubmitCommunityPoolSpendProposal(
	from, proposalPath string, deposit sdk.Coin, flags ...string,
) (bool, string, string) {

	cmd := fmt.Sprintf(
		"%s tx gov submit-proposal community-pool-spend %s --from=%s %v",
		f.LinkcliBinary, proposalPath, from, f.Flags(),
	)

	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

//___________________________________________________________________________________
// linkcli tx token

func (f *Fixtures) TxTokenPublish(from string, symbol, name string, amount int64, mintable bool, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token publish %s %s %s %d %t %v", f.LinkcliBinary, from, symbol, name, amount, mintable, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

func (f *Fixtures) TxTokenMint(to, amount string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token mint %s %s %v", f.LinkcliBinary, to, amount, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

func (f *Fixtures) TxTokenBurn(from, amount string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token burn %s %s %v", f.LinkcliBinary, from, amount, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

func (f *Fixtures) TxTokenGrantPerm(from string, to sdk.AccAddress, resource, action string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token grant %s %s %s %s %v", f.LinkcliBinary, from, to, resource, action, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

func (f *Fixtures) TxTokenRevokePerm(from string, resource, action string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx token revoke %s %s %s %v", f.LinkcliBinary, from, resource, action, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

//___________________________________________________________________________________
// linkcli tx safety box

func (f *Fixtures) TxSafetyBoxCreate(id, address string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx safetybox create %s %s %v", f.LinkcliBinary, id, address, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

func (f *Fixtures) TxSafetyBoxRole(id, action, role, from, to string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx safetybox role %s %s %s %s %s %v", f.LinkcliBinary, id, action, role, from, to, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

func (f *Fixtures) TxSafetyBoxSendCoins(id, action, denom string, amount int64, address, issuerAddress string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx safetybox sendcoins %s %s %s %d %s %s %v", f.LinkcliBinary, id, action, denom, amount, address, issuerAddress, f.Flags())
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

//___________________________________________________________________________________
// linkcli query safetybox

func (f *Fixtures) QuerySafetyBox(id string, flags ...string) safetyBoxModule.SafetyBox {
	cmd := fmt.Sprintf("%s query safetybox get %s %v", f.LinkcliBinary, id, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)

	cdc := app.MakeCodec()
	var sb safetyBoxModule.SafetyBox
	err := cdc.UnmarshalJSON([]byte(res), &sb)
	require.NoError(f.T, err)

	return sb
}

func (f *Fixtures) QuerySafetyBoxRole(id, role, address string, flags ...string) safetyBoxModule.MsgSafetyBoxRoleResponse {
	cmd := fmt.Sprintf("%s query safetybox role %s %s %s %v", f.LinkcliBinary, id, role, address, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)

	cdc := app.MakeCodec()
	var pms safetyBoxModule.MsgSafetyBoxRoleResponse
	err := cdc.UnmarshalJSON([]byte(res), &pms)
	require.NoError(f.T, err)

	return pms
}

//___________________________________________________________________________________
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

//___________________________________________________________________________________
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

//___________________________________________________________________________________
// linkcli query txs

// QueryTxs is linkcli query txs
func (f *Fixtures) QueryTxs(page, limit int, tags ...string) *sdk.SearchTxsResult {
	cmd := fmt.Sprintf("%s query txs --page=%d --limit=%d --tags='%s' %v", f.LinkcliBinary, page, limit, queryTags(tags), f.Flags())
	out, _ := tests.ExecuteT(f.T, cmd, "")
	var result sdk.SearchTxsResult
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return &result
}

// QueryTxsInvalid query txs with wrong parameters and compare expected error
func (f *Fixtures) QueryTxsInvalid(expectedErr error, page, limit int, tags ...string) {
	cmd := fmt.Sprintf("%s query txs --page=%d --limit=%d --tags='%s' %v", f.LinkcliBinary, page, limit, queryTags(tags), f.Flags())
	_, err := tests.ExecuteT(f.T, cmd, "")
	require.EqualError(f.T, expectedErr, err)
}

//___________________________________________________________________________________
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

//___________________________________________________________________________________
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

//___________________________________________________________________________________
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

//___________________________________________________________________________________
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

//___________________________________________________________________________________
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

// QueryLRC3
func (f *Fixtures) QueryLRC3(denom string, flags ...string) nft.Collection {
	cmd := fmt.Sprintf("%s query lrc3 get %s %v", f.LinkcliBinary, denom, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var lrc3 nft.Collection
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &lrc3)
	require.NoError(f.T, err)

	return lrc3
}

// QueryAllLRC3
func (f *Fixtures) QueryAllLRC3(flags ...string) nft.Collections {
	cmd := fmt.Sprintf("%s query lrc3 getAll %v", f.LinkcliBinary, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var initRes map[string]json.RawMessage
	var allLRC3 nft.Collections
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(f.T, err)
	cdc := codec.New()
	codec.RegisterCrypto(cdc)

	for _, value := range initRes {
		var lrc3 nft.Collection
		err = cdc.UnmarshalJSON(value, &lrc3)
		require.NoError(f.T, err, "value %v, err %v", value, err)
		allLRC3 = append(allLRC3, lrc3)
	}

	return allLRC3
}

// QueryLRC3BalanceOf
func (f *Fixtures) QueryLRC3BalanceOf(symbol, owner string, flags ...string) lrc3.TokenBalance {
	cmd := fmt.Sprintf("%s query lrc3 balanceOf %s %s %v", f.LinkcliBinary, symbol, owner, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var balance lrc3.TokenBalance
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &balance)
	require.NoError(f.T, err)

	return balance
}

// QueryLRC3OwnerOf
func (f *Fixtures) QueryLRC3OwnerOf(symbol string, tokenId int, flags ...string) lrc3.TokenOwner {
	cmd := fmt.Sprintf("%s query lrc3 ownerOf %s %d %v", f.LinkcliBinary, symbol, tokenId, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var owner lrc3.TokenOwner
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &owner)
	require.NoError(f.T, err)

	return owner
}

// QueryLRC3getApproved
func (f *Fixtures) QueryLRC3getApproved(symbol string, tokenId int, flags ...string) lrc3.Approval {
	cmd := fmt.Sprintf("%s query lrc3 getApproved %s %d %v", f.LinkcliBinary, symbol, tokenId, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var approve lrc3.Approval
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &approve)
	require.NoError(f.T, err)

	return approve
}

// QueryLRC3IsApprovedForAll
func (f *Fixtures) QueryLRC3IsApprovedForAll(symbol string, owner, operator string, flags ...string) lrc3.OperatorApprovals {
	cmd := fmt.Sprintf("%s query lrc3 isApprovedForAll %s %s %s %v", f.LinkcliBinary, symbol, owner, operator, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var operatorApprovals lrc3.OperatorApprovals
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &operatorApprovals)
	require.NoError(f.T, err)

	return operatorApprovals
}

//___________________________________________________________________________________
// query token

func (f *Fixtures) QueryToken(denom string, flags ...string) tokenModule.Token {
	cmd := fmt.Sprintf("%s query token symbol %s %s", f.LinkcliBinary, denom, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var token tokenModule.Token
	err := cdc.UnmarshalJSON([]byte(res), &token)
	require.NoError(f.T, err)
	return token
}

func (f *Fixtures) QueryTokenExpectEmpty(denom string, flags ...string) {
	cmd := fmt.Sprintf("%s query token symbol %s %s", f.LinkcliBinary, denom, f.Flags())
	_, errStr := tests.ExecuteT(f.T, cmd, "")
	require.NotEmpty(f.T, errStr)

}

func (f *Fixtures) QueryAccountPermission(addr sdk.AccAddress, flags ...string) tokenModule.Permissions {
	cmd := fmt.Sprintf("%s query token perm %s %s", f.LinkcliBinary, addr, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var pms tokenModule.Permissions
	err := cdc.UnmarshalJSON([]byte(res), &pms)
	require.NoError(f.T, err)
	return pms
}

func (f *Fixtures) QueryTokens(flags ...string) tokenModule.Tokens {
	cmd := fmt.Sprintf("%s query token symbols %s", f.LinkcliBinary, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var tokens tokenModule.Tokens
	err := cdc.UnmarshalJSON([]byte(res), &tokens)
	require.NoError(f.T, err)
	return tokens
}

func (f *Fixtures) QueryGenesisTxs(flags ...string) []sdk.Tx {
	cmd := fmt.Sprintf("%s query genesis-txs %s", f.LinkcliBinary, f.Flags())
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var txs []sdk.Tx
	err := cdc.UnmarshalJSON([]byte(res), &txs)
	require.NoError(f.T, err)
	return txs
}

func (f *Fixtures) QueryGenesisAccount(page, limit int, flags ...string) atypes.SearchGenesisAccountResult {
	cmd := fmt.Sprintf("%s query genesis-accounts --page=%d --limit=%d %s", f.LinkcliBinary, page, limit, f.Flags())
	return execQueryGenesisAccount(f, cmd)
}

func (f *Fixtures) QueryGenesisAccountByStrParams(page, limit string, flags ...string) atypes.SearchGenesisAccountResult {
	cmd := fmt.Sprintf("%s query genesis-accounts --page=%s --limit=%s %s", f.LinkcliBinary, page, limit, f.Flags())
	return execQueryGenesisAccount(f, cmd)
}

func execQueryGenesisAccount(f *Fixtures, cmd string) atypes.SearchGenesisAccountResult {
	res, errStr := tests.ExecuteT(f.T, cmd, "")
	require.Empty(f.T, errStr)
	cdc := app.MakeCodec()
	var accounts atypes.SearchGenesisAccountResult
	err := cdc.UnmarshalJSON([]byte(res), &accounts)
	require.NoError(f.T, err)
	return accounts
}

func (f *Fixtures) QueryGenesisAccountInvalid(expectedErr error, page, limit int, flags ...string) {
	cmd := fmt.Sprintf("%s query genesis-accounts --page=%d --limit=%d %s", f.LinkcliBinary, page, limit, f.Flags())
	execQueryGenesisAccountInvalid(f, cmd, expectedErr)
}

func (f *Fixtures) QueryGenesisAccountInvalidByStrParams(expectedErr error, page, limit string, flags ...string) {
	cmd := fmt.Sprintf("%s query genesis-accounts --page=%s --limit=%s %s", f.LinkcliBinary, page, limit, f.Flags())
	execQueryGenesisAccountInvalid(f, cmd, expectedErr)
}

func execQueryGenesisAccountInvalid(f *Fixtures, cmd string, expectedErr error) {
	_, err := tests.ExecuteT(f.T, cmd, "")
	require.EqualError(f.T, expectedErr, err)
}

//___________________________________________________________________________________
// linkcli mempool

// MempoolNumUnconfirmedTxs is linkcli mempool num-unconfirmed-txs
func (f *Fixtures) MempoolNumUnconfirmedTxs(flags ...string) *tmctypes.ResultUnconfirmedTxs {
	cmd := fmt.Sprintf("%s mempool num-unconfirmed-txs %v", f.LinkcliBinary, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var result tmctypes.ResultUnconfirmedTxs
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return &result
}

// MempoolUnconfirmedTxHashes is linkcli mempool unconfirmed-txs --hash
type ResultUnconfirmedTxHashes struct {
	Count      int      `json:"n_txs"`
	Total      int      `json:"total"`
	TotalBytes int64    `json:"total_bytes"`
	Txs        []string `json:"txs"`
}

func (f *Fixtures) MempoolUnconfirmedTxHashes(flags ...string) *ResultUnconfirmedTxHashes {
	cmd := fmt.Sprintf("%s mempool unconfirmed-txs --hash %v", f.LinkcliBinary, f.Flags())
	out, _ := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	var result ResultUnconfirmedTxHashes
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &result)
	require.NoError(f.T, err, "out %v\n, err %v", out, err)
	return &result
}

//___________________________________________________________________________________
// tendermint rpc
func (f *Fixtures) NetInfo(flags ...string) *tmctypes.ResultNetInfo {
	tmc := tmclient.NewHTTP(fmt.Sprintf("tcp://0.0.0.0:%s", f.Port), "/websocket")
	err := tmc.Start()
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
	tmc := tmclient.NewHTTP(fmt.Sprintf("tcp://0.0.0.0:%s", f.Port), "/websocket")
	err := tmc.Start()
	require.NoError(f.T, err)
	defer func() {
		err := tmc.Stop()
		require.NoError(f.T, err)
	}()

	netInfo, err := tmc.Status()
	require.NoError(f.T, err)
	return netInfo
}

//___________________________________________________________________________________
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

//___________________________________________________________________________________
// utils

func addFlags(cmd string, flags []string) string {
	for _, f := range flags {
		cmd += " " + f
	}
	return strings.TrimSpace(cmd)
}

func queryTags(tags []string) (out string) {
	for _, tag := range tags {
		out += tag + "&"
	}
	return strings.TrimSuffix(out, "&")
}

// Write the given string to a new temporary file
func WriteToNewTempFile(t *testing.T, s string) *os.File {
	fp, err := ioutil.TempFile(os.TempDir(), "cosmos_cli_test_")
	require.Nil(t, err)
	_, err = fp.WriteString(s)
	require.Nil(t, err)
	return fp
}

func marshalStdTx(t *testing.T, stdTx auth.StdTx) []byte {
	cdc := app.MakeCodec()
	bz, err := cdc.MarshalBinaryBare(stdTx)
	require.NoError(t, err)
	return bz
}

func unmarshalStdTx(t *testing.T, s string) (stdTx auth.StdTx) {
	cdc := app.MakeCodec()
	require.Nil(t, cdc.UnmarshalJSON([]byte(s), &stdTx))
	return
}

func unmarshalTxResponse(t *testing.T, s string) (txResp sdk.TxResponse) {
	cdc := app.MakeCodec()
	require.Nil(t, cdc.UnmarshalJSON([]byte(s), &txResp))
	return
}

//___________________________________________________________________________________
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
	if numOfNodes != nil && len(numOfNodes) == 1 {
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
			tests.WaitForTMStart(port)
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

	tests.WaitForTMStart(f.Port)
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
