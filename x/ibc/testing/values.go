/*
	This file contains the variables, constants, and default values
	used in the testing package and commonly defined in tests.
*/
package ibctesting

import (
	"time"

	sdk "github.com/line/lbm-sdk/types"

	ibctransfertypes "github.com/line/lbm-sdk/x/ibc/applications/transfer/types"
	connectiontypes "github.com/line/lbm-sdk/x/ibc/core/03-connection/types"
	commitmenttypes "github.com/line/lbm-sdk/x/ibc/core/23-commitment/types"
	ibctmtypes "github.com/line/lbm-sdk/x/ibc/light-clients/99-ostracon/types"
	"github.com/line/lbm-sdk/x/ibc/testing/mock"
)

const (
	FirstClientID     = "07-tendermint-0"
	FirstChannelID    = "channel-0"
	FirstConnectionID = "connection-0"

	// Default params constants used to create a TM client
	TrustingPeriod     time.Duration = time.Hour * 24 * 7 * 2
	UnbondingPeriod    time.Duration = time.Hour * 24 * 7 * 3
	MaxClockDrift      time.Duration = time.Second * 10
	DefaultDelayPeriod uint64        = 0

	DefaultChannelVersion = mock.Version
	InvalidID             = "IDisInvalid"

	// Application Ports
	TransferPort = ibctransfertypes.ModuleName
	MockPort     = mock.ModuleName

	// used for testing proposals
	Title       = "title"
	Description = "description"

	LongString = "LoremipsumdolorsitameconsecteturadipiscingeliseddoeiusmodtemporincididuntutlaboreetdoloremagnaaliquUtenimadminimveniamquisnostrudexercitationullamcolaborisnisiutaliquipexeacommodoconsequDuisauteiruredolorinreprehenderitinvoluptateelitsseillumoloreufugiatnullaariaturEcepteurintoccaectupidatatonroidentuntnulpauifficiaeseruntmollitanimidestlaborum"
)

var (
	DefaultOpenInitVersion *connectiontypes.Version

	// Default params variables used to create a TM client
	DefaultTrustLevel ibctmtypes.Fraction = ibctmtypes.DefaultTrustLevel
	TestCoin                              = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))

	UpgradePath = []string{"upgrade", "upgradedIBCState"}

	ConnectionVersion = connectiontypes.ExportedVersionsToProto(connectiontypes.GetCompatibleVersions())[0]

	MockAcknowledgement          = mock.MockAcknowledgement.Acknowledgement()
	MockPacketData               = mock.MockPacketData
	MockFailPacketData           = mock.MockFailPacketData
	MockRecvCanaryCapabilityName = mock.MockRecvCanaryCapabilityName

	prefix = commitmenttypes.NewMerklePrefix([]byte("ibc"))
)
