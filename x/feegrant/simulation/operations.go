package simulation

import (
	"math/rand"

	"github.com/line/lbm-sdk/baseapp"
	"github.com/line/lbm-sdk/codec"
	simappparams "github.com/line/lbm-sdk/simapp/params"
	sdk "github.com/line/lbm-sdk/types"
	simtypes "github.com/line/lbm-sdk/types/simulation"
	"github.com/line/lbm-sdk/x/feegrant/keeper"
	"github.com/line/lbm-sdk/x/feegrant/types"
	"github.com/line/lbm-sdk/x/simulation"
)

// Simulation operation weights constants
const (
	OpWeightMsgGrantFeeAllowance  = "op_weight_msg_grant_fee_allowance"
	OpWeightMsgRevokeFeeAllowance = "op_weight_msg_grant_revoke_allowance"
)

var (
	TypeMsgGrantFeeAllowance  = sdk.MsgTypeURL(&types.MsgGrantFeeAllowance{})
	TypeMsgRevokeFeeAllowance = sdk.MsgTypeURL(&types.MsgRevokeFeeAllowance{})
)

func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONMarshaler,
	ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper,
) simulation.WeightedOperations {

	var (
		weightMsgGrantFeeAllowance  int
		weightMsgRevokeFeeAllowance int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgGrantFeeAllowance, &weightMsgGrantFeeAllowance, nil,
		func(_ *rand.Rand) {
			weightMsgGrantFeeAllowance = simappparams.DefaultWeightGrantFeeAllowance
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgRevokeFeeAllowance, &weightMsgRevokeFeeAllowance, nil,
		func(_ *rand.Rand) {
			weightMsgRevokeFeeAllowance = simappparams.DefaultWeightRevokeFeeAllowance
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgGrantFeeAllowance,
			SimulateMsgGrantFeeAllowance(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgRevokeFeeAllowance,
			SimulateMsgRevokeFeeAllowance(ak, bk, k),
		),
	}
}

// SimulateMsgGrantFeeAllowance generates MsgGrantFeeAllowance with random values.
func SimulateMsgGrantFeeAllowance(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		granter, _ := simtypes.RandomAcc(r, accs)
		grantee, _ := simtypes.RandomAcc(r, accs)
		if grantee.Address.String() == granter.Address.String() {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgGrantFeeAllowance, "grantee and granter cannot be same"), nil, nil
		}

		if f, _ := k.GetAllowance(ctx, granter.Address, grantee.Address); f != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgGrantFeeAllowance, "fee allowance exists"), nil, nil
		}

		account := ak.GetAccount(ctx, granter.Address)

		spendableCoins := bk.SpendableCoins(ctx, account.GetAddress())
		if spendableCoins.Empty() {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgGrantFeeAllowance, "unable to grant empty coins as SpendLimit"), nil, nil
		}

		oneYear := ctx.BlockTime().AddDate(1, 0, 0)
		msg, err := types.NewMsgGrantFeeAllowance(&types.BasicAllowance{
			SpendLimit: spendableCoins,
			Expiration: &oneYear,
		}, granter.Address, grantee.Address)

		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgGrantFeeAllowance, err.Error()), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         TypeMsgGrantFeeAllowance,
			Context:         ctx,
			SimAccount:      granter,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendableCoins,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgRevokeFeeAllowance generates a MsgRevokeFeeAllowance with random values.
func SimulateMsgRevokeFeeAllowance(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		hasGrant := false
		var granterAddr sdk.AccAddress
		var granteeAddr sdk.AccAddress
		k.IterateAllFeeAllowances(ctx, func(grant types.Grant) bool {

			granter := sdk.AccAddress(grant.Granter)
			grantee := sdk.AccAddress(grant.Grantee)
			granterAddr = granter
			granteeAddr = grantee
			hasGrant = true
			return true
		})

		if !hasGrant {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgRevokeFeeAllowance, "no grants"), nil, nil
		}
		granter, ok := simtypes.FindAccount(accs, granterAddr)

		if !ok {
			return simtypes.NoOpMsg(types.ModuleName, TypeMsgRevokeFeeAllowance, "Account not found"), nil, nil
		}

		account := ak.GetAccount(ctx, granter.Address)
		spendableCoins := bk.SpendableCoins(ctx, account.GetAddress())

		msg := types.NewMsgRevokeFeeAllowance(granterAddr, granteeAddr)

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simappparams.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         TypeMsgRevokeFeeAllowance,
			Context:         ctx,
			SimAccount:      granter,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      types.ModuleName,
			CoinsSpentInMsg: spendableCoins,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}
