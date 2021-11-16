package ante

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/auth/keeper"
	types3 "github.com/line/lbm-sdk/x/auth/types"
	feegrantkeeper "github.com/line/lbm-sdk/x/feegrant/keeper"
	types2 "github.com/line/lbm-sdk/x/feegrant/types"
)

// RejectFeeGranterDecorator is an AnteDecorator which rejects transactions which
// have the Fee.granter field set. It is to be used by chains which do not support
// fee grants.
type RejectFeeGranterDecorator struct{}

// NewRejectFeeGranterDecorator returns a new RejectFeeGranterDecorator.
func NewRejectFeeGranterDecorator() RejectFeeGranterDecorator {
	return RejectFeeGranterDecorator{}
}

var _ sdk.AnteDecorator = RejectFeeGranterDecorator{}

func (d RejectFeeGranterDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if ok && len(feeTx.FeeGranter()) != 0 {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "fee grants are not supported")
	}

	return next(ctx, tx, simulate)
}

// DeductGrantedFeeDecorator deducts fees from fee_payer or fee_granter (if exists a valid fee allowance) of the tx
// If the fee_payer or fee_granter does not have the funds to pay for the fees, return with InsufficientFunds error
// Call next AnteHandler if fees successfully deducted
// CONTRACT: Tx must implement GrantedFeeTx interface to use DeductGrantedFeeDecorator
type DeductGrantedFeeDecorator struct {
	ak keeper.AccountKeeper
	k  feegrantkeeper.Keeper
	bk types2.BankKeeper
}

func NewDeductGrantedFeeDecorator(ak keeper.AccountKeeper, bk types2.BankKeeper, k feegrantkeeper.Keeper) DeductGrantedFeeDecorator {
	return DeductGrantedFeeDecorator{
		ak: ak,
		k:  k,
		bk: bk,
	}
}

// AnteHandle performs a decorated ante-handler responsible for deducting transaction
// fees. Fees will be deducted from the account designated by the FeePayer on a
// transaction by default. However, if the fee payer differs from the transaction
// signer, the handler will check if a fee grant has been authorized. If the
// transaction's signer does not exist, it will be created.
func (d DeductGrantedFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a GrantedFeeTx")
	}

	// sanity check from DeductFeeDecorator
	if addr := d.ak.GetModuleAddress(types3.FeeCollectorName); addr == "" {
		panic(fmt.Sprintf("%s module account has not been set", types3.FeeCollectorName))
	}

	fee := feeTx.GetFee()
	feePayer := feeTx.FeePayer()
	feeGranter := feeTx.FeeGranter()

	deductFeesFrom := feePayer

	// ensure the grant is allowed, if we request a different fee payer
	if feeGranter != "" && !feeGranter.Equals(feePayer) {
		err := d.k.UseGrantedFees(ctx, feeGranter, feePayer, fee, tx.GetMsgs())
		if err != nil {
			return ctx, sdkerrors.Wrapf(err, "%s not allowed to pay fees from %s", feeGranter, feePayer)
		}

		deductFeesFrom = feeGranter
	}

	// now, either way, we know that we are authorized to deduct the fees from the deductFeesFrom account
	deductFeesFromAcc := d.ak.GetAccount(ctx, deductFeesFrom)
	if deductFeesFromAcc == nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", deductFeesFrom)
	}

	// move on if there is no fee to deduct
	if fee.IsZero() {
		return next(ctx, tx, simulate)
	}

	// deduct fee if non-zero
	err = DeductFees(d.bk, ctx, deductFeesFromAcc, fee)
	if err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}
