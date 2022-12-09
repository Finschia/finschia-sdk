package rest

import (
	"net/http"

	govrest "github.com/line/lbm-sdk/x/gov/client/rest"

	"github.com/line/lbm-sdk/client"
)

func DummyRESTHandler(_ client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "foundation",
		Handler:  func(_ http.ResponseWriter, _ *http.Request) {},
	}
}
