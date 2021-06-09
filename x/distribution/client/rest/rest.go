package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/line/lfb-sdk/client"
	clientrest "github.com/line/lfb-sdk/client/rest"
	"github.com/line/lfb-sdk/client/tx"
	"github.com/line/lfb-sdk/types/rest"
	"github.com/line/lfb-sdk/x/distribution/types"
	govrest "github.com/line/lfb-sdk/x/gov/client/rest"
	govtypes "github.com/line/lfb-sdk/x/gov/types"
)

func RegisterHandlers(clientCtx client.Context, rtr *mux.Router) {
	r := clientrest.WithHTTPDeprecationHeaders(rtr)

	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
}

// TODO add proto compatible Handler after x/gov migration
// ProposalRESTHandler returns a ProposalRESTHandler that exposes the community pool spend REST handler with a given sub-route.
func ProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "community_pool_spend",
		Handler:  postProposalHandlerFn(clientCtx),
	}
}

func postProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CommunityPoolSpendProposalReq
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		content := types.NewCommunityPoolSpendProposal(req.Title, req.Description, req.Recipient, req.Amount)

		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
