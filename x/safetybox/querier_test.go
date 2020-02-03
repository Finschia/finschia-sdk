package safetybox

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testCommon "github.com/line/link/x/safetybox/internal/keeper"
	"github.com/line/link/x/safetybox/internal/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestSafetyBoxQuerierSafetyBox(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	h := NewHandler(input.Keeper)

	// create a box
	owner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	msgSbCreate := MsgSafetyBoxCreate{
		SafetyBoxId:     SafetyBoxTestId,
		SafetyBoxOwner:  owner,
		SafetyBoxDenoms: []string{"link"},
	}
	r := h(input.Ctx, msgSbCreate)
	require.True(t, r.IsOK())

	// query the box
	params := types.QuerySafetyBoxParams{SafetyBoxId: SafetyBoxTestId}
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
	require.Equal(t, SafetyBoxTestId, sb.ID)
	require.Equal(t, owner, sb.Owner)
	require.Equal(t, sdk.Coins(nil), sb.TotalAllocation)
	require.Equal(t, sdk.Coins(nil), sb.CumulativeAllocation)
	require.Equal(t, sdk.Coins(nil), sb.TotalIssuance)
}

func TestSafetyBoxQuerierRole(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	h := NewHandler(input.Keeper)

	// create a box
	owner := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	msgSbCreate := MsgSafetyBoxCreate{
		SafetyBoxId:     SafetyBoxTestId,
		SafetyBoxOwner:  owner,
		SafetyBoxDenoms: []string{"link"},
	}
	r := h(input.Ctx, msgSbCreate)
	require.True(t, r.IsOK())

	// check the owner of the box
	params := types.QueryAccountRoleParams{
		SafetyBoxId: SafetyBoxTestId,
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
