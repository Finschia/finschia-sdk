package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/line/link/x/safetybox/internal/types"
)

func SafetyBoxQueryHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		safetyBoxID := vars["id"]
		if len(safetyBoxID) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "SafetyBoxId is required")
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		safetyBoxGetter := types.NewSafetyBoxRetriever(cliCtx)

		sb, height, err := safetyBoxGetter.GetSafetyBox(safetyBoxID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), sb)
	}
}

func SafetyBoxRoleQueryHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		safetyBoxID := vars["id"]
		if len(safetyBoxID) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "SafetyBoxId is required")
			return
		}

		address := vars["address"]
		if len(address) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Address is required")
			return
		}

		addr, err := sdk.AccAddressFromBech32(address)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		role := vars["role"]
		if len(role) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Role is required")
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		permGetter := types.NewAccountPermissionRetriever(cliCtx)
		pms, height, err := permGetter.GetAccountPermissions(safetyBoxID, role, addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		}

		rest.PostProcessResponse(w, cliCtx.WithHeight(height), pms)
	}
}
