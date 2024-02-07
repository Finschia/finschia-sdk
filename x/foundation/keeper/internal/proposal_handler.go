package internal

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

// NewFoundationProposalsHandler creates a handler for the gov proposals.
func NewFoundationProposalsHandler(k Keeper) govv1beta1types.Handler {
	return func(ctx sdk.Context, content govv1beta1types.Content) error {
		switch c := content.(type) {
		case *foundation.FoundationExecProposal:
			return handleFoundationExecProposal(ctx, k, *c)

		default:
			return sdkerrors.ErrUnknownRequest.Wrapf("unrecognized param proposal content type: %T", c)
		}
	}
}

func handleFoundationExecProposal(ctx sdk.Context, k Keeper, proposal foundation.FoundationExecProposal) error {
	msgs := foundation.GetMessagesFromFoundationExecProposal(proposal)

	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	if err := ensureMsgAuthz(msgs, authority, k.cdc); err != nil {
		return err
	}

	// allow the following messages
	allowedUrls := map[string]bool{
		sdk.MsgTypeURL((*foundation.MsgUpdateCensorship)(nil)): true,
		sdk.MsgTypeURL((*foundation.MsgGrant)(nil)):            true,
		sdk.MsgTypeURL((*foundation.MsgRevoke)(nil)):           true,
	}

	for i, msg := range msgs {
		url := sdk.MsgTypeURL(msg)
		if !allowedUrls[url] {
			return sdkerrors.ErrInvalidRequest.Wrapf("%s not allowed", url)
		}

		handler := k.router.Handler(msg)
		if handler == nil {
			return sdkerrors.ErrUnknownRequest.Wrapf("no message handler found for %q", url)
		}
		_, err := handler(ctx, msg)
		if err != nil {
			return errorsmod.Wrapf(err, "message %q at position %d", msg, i)
		}
	}

	return nil
}
