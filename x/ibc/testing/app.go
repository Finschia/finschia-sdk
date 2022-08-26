package ibctesting

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/client"
	"github.com/line/lbm-sdk/codec"
	codectypes "github.com/line/lbm-sdk/codec/types"
	cryptocodec "github.com/line/lbm-sdk/crypto/codec"
	sdk "github.com/line/lbm-sdk/types"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	banktypes "github.com/line/lbm-sdk/x/bank/types"
	capabilitykeeper "github.com/line/lbm-sdk/x/capability/keeper"
	stakingkeeper "github.com/line/lbm-sdk/x/staking/keeper"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	abci "github.com/line/ostracon/abci/types"
	"github.com/line/ostracon/libs/log"
	tmproto "github.com/line/ostracon/proto/ostracon/types"
	tmtypes "github.com/line/ostracon/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/line/lbm-sdk/x/ibc/core/keeper"
	"github.com/line/lbm-sdk/x/ibc/testing/simapp"
)

var DefaultTestingAppInit func() (TestingApp, map[string]json.RawMessage) = SetupTestingApp

type TestingApp interface {
	abci.Application

	// ibc-go additions
	GetBaseApp() *baseapp.BaseApp
	GetStakingKeeper() stakingkeeper.Keeper
	GetIBCKeeper() *keeper.Keeper
	GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper
	GetTxConfig() client.TxConfig

	// Implemented by SimApp
	AppCodec() codec.Codec

	// Implemented by BaseApp
	LastCommitID() sdk.CommitID
	LastBlockHeight() int64
}

func SetupTestingApp() (TestingApp, map[string]json.RawMessage) {
	db := dbm.NewMemDB()
	encCdc := simapp.MakeTestEncodingConfig()
	app := simapp.NewSimApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, simapp.DefaultNodeHome, 5, encCdc, simapp.EmptyAppOptions{})
	return app, simapp.NewDefaultGenesisState(encCdc.Marshaler)
}

// SetupWithGenesisValSet initializes a new SimApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
// account. A Nop logger is set in SimApp.
func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, chainID string, powerReduction sdk.Int, balances ...banktypes.Balance) TestingApp {
	app, genesisState := DefaultTestingAppInit()

	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.TokensFromConsensusPower(1, powerReduction)

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.ZeroInt(),
		}

		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))
	}

	// set validators and delegations
	var stakingGenesis stakingtypes.GenesisState
	app.AppCodec().MustUnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesis)

	bondDenom := stakingGenesis.Params.BondDenom

	totalSupply := sdk.NewCoins()

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(bondDenom, bondAmt.Mul(sdk.NewInt(int64(len(valSet.Validators)))))},
	})

	// set validators and delegations
	stakingGenesis = *stakingtypes.NewGenesisState(stakingGenesis.Params, validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(&stakingGenesis)

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			ChainId:         chainID,
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	app.Commit()
	app.BeginBlock(
		abci.RequestBeginBlock{
			Header: tmproto.Header{
				ChainID:            chainID,
				Height:             app.LastBlockHeight() + 1,
				AppHash:            app.LastCommitID().Hash,
				ValidatorsHash:     valSet.Hash(),
				NextValidatorsHash: valSet.Hash(),
			},
		},
	)

	return app
}
