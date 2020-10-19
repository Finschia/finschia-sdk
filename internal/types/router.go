package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/token"
)

var _ Router = (*router)(nil)

type Router interface {
	AddRoute(r string, h token.EncodeHandler) (rtr Router)
	HasRoute(r string) bool
	GetRoute(path string) (h token.EncodeHandler)
	Seal()
}

type router struct {
	routes map[string]token.EncodeHandler
	sealed bool
}

func NewRouter() Router {
	return &router{
		routes: make(map[string]token.EncodeHandler),
	}
}

func (rtr *router) Seal() {
	if rtr.sealed {
		panic("router already sealed")
	}
	rtr.sealed = true
}

func (rtr *router) AddRoute(path string, h token.EncodeHandler) Router {
	if rtr.sealed {
		panic("router sealed; cannot add route handler")
	}

	if !sdk.IsAlphaNumeric(path) {
		panic("route expressions can only contain alphanumeric characters")
	}
	if rtr.HasRoute(path) {
		panic(fmt.Sprintf("route %s has already been initialized", path))
	}

	rtr.routes[path] = h
	return rtr
}

// HasRoute returns true if the router has a path registered or false otherwise.
func (rtr *router) HasRoute(path string) bool {
	return rtr.routes[path] != nil
}

// GetRoute returns a Handler for a given path.
func (rtr *router) GetRoute(path string) token.EncodeHandler {
	if !rtr.HasRoute(path) {
		panic(fmt.Sprintf("route \"%s\" does not exist", path))
	}

	return rtr.routes[path]
}
