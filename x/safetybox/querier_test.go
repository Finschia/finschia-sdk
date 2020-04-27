package safetybox

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testCommon "github.com/line/link/x/safetybox/internal/keeper"
	"github.com/line/link/x/safetybox/internal/types"
	"github.com/line/link/x/token"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestSafetyBoxQuerierSafetyBox(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	h := NewHandler(input.Keeper)

	// issue token
	testToken := token.NewToken(TestContractID, TestTokenName, TestTokenSymbol, TestTokenMeta, TestTokenImageURI, sdk.NewInt(TestTokenDecimals), true)

	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	err := input.Tk.IssueToken(input.Ctx, testToken, sdk.NewInt(TestTokenAmount), addr, addr)
	require.NoError(t, err)

	tok, err := input.Tk.GetToken(input.Ctx, TestContractID)
	require.NoError(t, err)
	require.Equal(t, testToken.GetContractID(), tok.GetContractID())

	// create a box
	owner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	msgSbCreate := MsgSafetyBoxCreate{
		SafetyBoxID:    SafetyBoxTestID,
		SafetyBoxOwner: owner,
		ContractID:     TestContractID,
	}
	_, err = h(input.Ctx, msgSbCreate)
	require.NoError(t, err)

	// query the box
	params := types.QuerySafetyBoxParams{SafetyBoxID: SafetyBoxTestID}
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySafetyBox),
		Data: []byte(params.String()),
	}
	path := []string{types.QuerySafetyBox}
	querier := NewQuerier(input.Keeper)

	res, err := querier(input.Ctx, path, req)
	require.NoError(t, err)

	// unmarshal the response
	var sb types.SafetyBox
	err2 := input.Cdc.UnmarshalJSON(res, &sb)
	require.NoError(t, err2)

	// verify
	require.Equal(t, SafetyBoxTestID, sb.ID)
	require.Equal(t, owner, sb.Owner)
	require.Equal(t, sdk.ZeroInt(), sb.TotalAllocation)
	require.Equal(t, sdk.ZeroInt(), sb.CumulativeAllocation)
	require.Equal(t, sdk.ZeroInt(), sb.TotalIssuance)
}

func TestSafetyBoxQuerierRole(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	h := NewHandler(input.Keeper)

	// issue token
	testToken := token.NewToken(TestContractID, TestTokenName, TestTokenSymbol, TestTokenMeta, TestTokenImageURI, sdk.NewInt(TestTokenDecimals), true)

	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	err := input.Tk.IssueToken(input.Ctx, testToken, sdk.NewInt(TestTokenAmount), addr, addr)
	require.NoError(t, err)

	tok, err := input.Tk.GetToken(input.Ctx, TestContractID)
	require.NoError(t, err)
	require.Equal(t, testToken.GetContractID(), tok.GetContractID())

	// create a box
	owner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	msgSbCreate := MsgSafetyBoxCreate{
		SafetyBoxID:    SafetyBoxTestID,
		SafetyBoxOwner: owner,
		ContractID:     TestContractID,
	}
	_, err = h(input.Ctx, msgSbCreate)
	require.NoError(t, err)

	// check the owner of the box
	params := types.QueryAccountRoleParams{
		SafetyBoxID: SafetyBoxTestID,
		Role:        types.RoleOwner,
		Address:     owner,
	}
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAccountRole),
		Data: []byte(params.String()),
	}
	path := []string{types.QueryAccountRole}
	querier := NewQuerier(input.Keeper)

	res, err := querier(input.Ctx, path, req)
	require.NoError(t, err)

	// unmarshal the response
	var sb types.MsgSafetyBoxRoleResponse
	err2 := input.Cdc.UnmarshalJSON(res, &sb)
	require.NoError(t, err2)

	// verify
	require.True(t, sb.HasRole)
}
