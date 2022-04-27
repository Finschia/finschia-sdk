package keeper_test

import (
	"context"
	"testing"

	ocproto "github.com/line/ostracon/proto/ostracon/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	"github.com/line/lbm-sdk/x/foundation/keeper"
)

var (
	delPk   = ed25519.GenPrivKey().PubKey()
	delAddr = sdk.BytesToAccAddress(delPk.Address())
	valAddr = delAddr.ToValAddress()
)

func TestCleanup(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})

	k := app.FoundationKeeper

	// add auths
	auth := foundation.ValidatorAuth{
		OperatorAddress: valAddr.String(),
		CreationAllowed: true,
	}
	require.NoError(t, k.SetValidatorAuth(ctx, auth))

	// cleanup
	k.Cleanup(ctx)
	require.Empty(t, k.GetValidatorAuths(ctx))
}

type KeeperTestSuite struct {
	suite.Suite
	ctx  sdk.Context
	goCtx context.Context
	keeper keeper.Keeper
	queryServer foundation.QueryServer
	msgServer foundation.MsgServer

	operator sdk.AccAddress

	member sdk.AccAddress
	comingMember sdk.AccAddress
	leavingMember sdk.AccAddress
	badMember sdk.AccAddress

	stranger sdk.AccAddress

	balance sdk.Int
}

func (s *KeeperTestSuite) SetupTest() {
	checkTx := false
	app := simapp.Setup(checkTx)
	s.ctx = app.BaseApp.NewContext(checkTx, ocproto.Header{})
	s.goCtx = sdk.WrapSDKContext(s.ctx)
	s.keeper = app.FoundationKeeper

	s.queryServer = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServer(s.keeper)

	s.operator = s.keeper.GetOperator(s.ctx)
	s.member = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	s.comingMember = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	s.leavingMember = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())
	s.stranger = sdk.BytesToAccAddress(secp256k1.GenPrivKey().PubKey().Address())

	s.balance = sdk.NewInt(1000000)
	holders := []sdk.AccAddress{
		s.stranger,
		app.AccountKeeper.GetModuleAccount(s.ctx, foundation.TreasuryName).GetAddress(),
	}
	for _, holder := range holders {
		err := app.BankKeeper.SetBalance(s.ctx, holder, sdk.NewCoin(sdk.DefaultBondDenom, s.balance))
		s.Require().NoError(err)
	}

	err := s.keeper.UpdateMembers(s.ctx, []foundation.Member{
		{
			Address: s.member.String(),
			Weight: sdk.OneDec(),
			Metadata: "a permanent member",
		},
		{
			Address: s.leavingMember.String(),
			Weight: sdk.OneDec(),
			Metadata: "a member to leave",
		},
		{
			Address: s.badMember.String(),
			Weight: sdk.OneDec(),
			Metadata: "a member to be deported ",
		},
	})
	s.Require().NoError(err)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
