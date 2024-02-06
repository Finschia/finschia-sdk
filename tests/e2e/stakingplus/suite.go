package stakingplus

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/suite"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client/flags"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
	"github.com/Finschia/finschia-sdk/x/stakingplus"
)

type E2ETestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network

	grantee         sdk.AccAddress
	permanentMember sdk.AccAddress
	stranger        sdk.AccAddress

	addressCodec    address.Codec
	valAddressCodec address.Codec
}

func NewE2ETestSuite(cfg network.Config) *E2ETestSuite {
	return &E2ETestSuite{cfg: cfg}
}

func (s *E2ETestSuite) SetupSuite() {
	s.T().Log("setting up e2e test suite")

	s.addressCodec = addresscodec.NewBech32Codec("link")
	s.valAddressCodec = addresscodec.NewBech32Codec("linkvaloper")

	genesisState := s.cfg.GenesisState

	var foundationData foundation.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[foundation.ModuleName], &foundationData))

	// enable foundation tax
	params := foundation.Params{
		FoundationTax: math.LegacyMustNewDecFromStr("0.2"),
	}
	foundationData.Params = params

	var strangerMnemonic string
	var granteeMnemonic string
	var permanentMemberMnemonic string
	granteeMnemonic, s.grantee = s.createMnemonic("grantee")
	strangerMnemonic, s.stranger = s.createMnemonic("stranger")
	permanentMemberMnemonic, s.permanentMember = s.createMnemonic("permanentmember")

	foundationData.Members = []foundation.Member{
		{
			Address:  s.bytesToString(s.permanentMember),
			Metadata: "permanent member",
		},
	}

	info := foundation.DefaultFoundation()
	info.TotalWeight = math.LegacyNewDecFromInt(math.NewInt(int64(len(foundationData.Members))))
	err := info.SetDecisionPolicy(&foundation.ThresholdDecisionPolicy{
		Threshold: math.LegacyOneDec(),
		Windows: &foundation.DecisionPolicyWindows{
			VotingPeriod: 7 * 24 * time.Hour,
		},
	})
	s.Require().NoError(err)
	foundationData.Foundation = info

	// enable censorship
	censorships := []foundation.Censorship{
		{
			MsgTypeUrl: sdk.MsgTypeURL((*stakingtypes.MsgCreateValidator)(nil)),
			Authority:  foundation.CensorshipAuthorityFoundation,
		},
	}
	foundationData.Censorships = censorships

	val1 := getValidator(s.T(), s.T().TempDir(), s.cfg, 0)
	for _, grantee := range []sdk.AccAddress{s.grantee, val1} {
		ga := foundation.GrantAuthorization{
			Grantee: s.bytesToString(grantee),
		}.WithAuthorization(&stakingplus.CreateValidatorAuthorization{
			ValidatorAddress: s.bytesToValString(grantee),
		})
		s.Require().NotNil(ga)
		foundationData.Authorizations = append(foundationData.Authorizations, *ga)
	}

	foundationDataBz, err := s.cfg.Codec.MarshalJSON(&foundationData)
	s.Require().NoError(err)
	genesisState[foundation.ModuleName] = foundationDataBz
	s.cfg.GenesisState = genesisState

	s.network, err = network.New(s.T(), s.T().TempDir(), s.cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	s.createAccount("grantee", granteeMnemonic)
	s.createAccount("stranger", strangerMnemonic)
	s.createAccount("permanentmember", permanentMemberMnemonic)
}

func (s *E2ETestSuite) TearDownSuite() {
	s.T().Log("tearing down e2e test suite")
	s.network.Cleanup()
}

func (s *E2ETestSuite) bytesToString(addr sdk.AccAddress) string {
	str, err := s.addressCodec.BytesToString(addr)
	s.Require().NoError(err)
	return str
}

func (s *E2ETestSuite) bytesToValString(addr sdk.AccAddress) string {
	str, err := s.valAddressCodec.BytesToString(addr)
	s.Require().NoError(err)
	return str
}

// creates an account
func (s *E2ETestSuite) createMnemonic(uid string) (string, sdk.AccAddress) {
	cstore := keyring.NewInMemory(s.cfg.Codec)
	info, mnemonic, err := cstore.NewMnemonic(uid, keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	s.Require().NoError(err)

	addr, err := info.GetAddress()
	s.Require().NoError(err)

	return mnemonic, addr
}

// creates an account and send some coins to it for the future transactions.
func (s *E2ETestSuite) createAccount(uid, mnemonic string) {
	commonArgs := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100)))),
	}

	val := s.network.Validators[0]
	info, err := val.ClientCtx.Keyring.NewAccount(uid, mnemonic, keyring.DefaultBIP39Passphrase, sdk.FullFundraiserPath, hd.Secp256k1)
	s.Require().NoError(err)

	addr, err := info.GetAddress()
	s.Require().NoError(err)

	fee := sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, math.NewInt(1000)))
	args := append([]string{
		s.bytesToString(val.Address),
		s.bytesToString(addr),
		fee.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, s.bytesToString(val.Address)),
	}, commonArgs...)
	out, err := clitestutil.ExecTestCLICmd(val.ClientCtx, bankcli.NewSendTxCmd(s.addressCodec), args)
	s.Require().NoError(err)

	var res sdk.TxResponse
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &res), out.String())
	s.Require().Zero(res.Code, out.String())

	s.Require().NoError(s.network.WaitForNextBlock())
}
