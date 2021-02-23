package keeper

import (
	"context"
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/contract"
	"github.com/line/lbm-sdk/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestApproveScenario(t *testing.T) {
	ctx := cacheKeeper()

	// prepare token
	someToken := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
	require.NoError(t, keeper.IssueToken(ctx, someToken, sdk.NewInt(defaultAmount), addr1, addr1))

	// approve test
	anotherCtx := ctx.WithContext(context.WithValue(ctx.Context(), contract.CtxKey{}, anotherContractID))
	require.EqualError(t, keeper.SetApproved(anotherCtx, addr3, addr1), sdkerrors.Wrapf(types.ErrTokenNotExist, "ContractID: %s", anotherContractID).Error())
	require.NoError(t, keeper.SetApproved(ctx, addr3, addr1))
	require.EqualError(t, keeper.SetApproved(ctx, addr3, addr1), sdkerrors.Wrapf(types.ErrTokenAlreadyApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr3.String(), addr1.String(), defaultContractID).Error())

	// transfer_from test
	require.EqualError(t, keeper.TransferFrom(ctx, addr2, addr1, addr2, sdk.NewInt(10)), sdkerrors.Wrapf(types.ErrTokenNotApproved, "Proxy: %s, Approver: %s, ContractID: %s", addr2.String(), addr1.String(), defaultContractID).Error())
	require.NoError(t, keeper.TransferFrom(ctx, addr3, addr1, addr2, sdk.NewInt(10)))

	t.Log("add one more approver fo test")
	{
		require.NoError(t, keeper.SetApproved(ctx, addr3, addr2))
	}

	t.Log("succeed to GetApprovers")
	{
		approvers, err := keeper.GetApprovers(ctx, addr3)
		require.NoError(t, err)
		require.True(t, len(approvers) == 2)
		require.True(t, types.IsAddressContains(approvers, addr1))
		require.True(t, types.IsAddressContains(approvers, addr2))
	}

	t.Log("fail to GetApprovres")
	{
		_, err := keeper.GetApprovers(anotherCtx, addr3)
		require.Error(t, err, "")
	}
}
