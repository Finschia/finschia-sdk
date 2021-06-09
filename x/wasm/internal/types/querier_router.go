package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/line/lfb-sdk/types"
)

var _ QueryRouter = (*querierRouter)(nil)

type EncodeQuerier func(ctx sdk.Context, jsonQuerier json.RawMessage) ([]byte, error)

// QueryRouter provides queryables for each query path.
type QueryRouter interface {
	AddRoute(r string, q EncodeQuerier) QueryRouter
	HasRoute(r string) bool
	GetRoute(path string) EncodeQuerier
	Seal()
}

type querierRouter struct {
	routes map[string]EncodeQuerier
	sealed bool
}

func NewQuerierRouter() QueryRouter {
	return &querierRouter{
		routes: make(map[string]EncodeQuerier),
	}
}

func (rtr *querierRouter) Seal() {
	if rtr.sealed {
		panic("querier router already sealed")
	}
	rtr.sealed = true
}

func (rtr *querierRouter) AddRoute(path string, q EncodeQuerier) QueryRouter {
	if rtr.sealed {
		panic("router sealed; cannot add route handler")
	}
	if !sdk.IsAlphaNumeric(path) {
		panic("querier route expressions can only contain alphanumeric characters")
	}
	if rtr.HasRoute(path) {
		panic(fmt.Sprintf("querier route %s has already been initialized", path))
	}

	rtr.routes[path] = q
	return rtr
}

func (rtr *querierRouter) HasRoute(path string) bool {
	return rtr.routes[path] != nil
}

func (rtr *querierRouter) GetRoute(path string) EncodeQuerier {
	if !rtr.HasRoute(path) {
		panic(fmt.Sprintf("querier route \"%s\" does not exist", path))
	}

	return rtr.routes[path]
}
