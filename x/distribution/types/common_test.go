package types

import (
	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/line/lbm-sdk/crypto/types"
	sdk "github.com/line/lbm-sdk/types"
)

// nolint:deadcode,unused,varcheck
var (
	delPk1       = ed25519.GenPrivKey().PubKey()
	delPk2       = ed25519.GenPrivKey().PubKey()
	delPk3       = ed25519.GenPrivKey().PubKey()
	delAddr1     = sdk.BytesToAccAddress(delPk1.Address())
	delAddr2     = sdk.BytesToAccAddress(delPk2.Address())
	delAddr3     = sdk.BytesToAccAddress(delPk3.Address())
	emptyDelAddr sdk.AccAddress

	valPk1       = ed25519.GenPrivKey().PubKey()
	valPk2       = ed25519.GenPrivKey().PubKey()
	valPk3       = ed25519.GenPrivKey().PubKey()
	valAddr1     = sdk.BytesToValAddress(valPk1.Address())
	valAddr2     = sdk.BytesToValAddress(valPk2.Address())
	valAddr3     = sdk.BytesToValAddress(valPk3.Address())
	emptyValAddr sdk.ValAddress

	emptyPubkey cryptotypes.PubKey
)
