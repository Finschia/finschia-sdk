package types

import (
	"encoding/json"
	"fmt"

	sdk "github.com/line/lfb-sdk/types"
)

var _ Router = (*router)(nil)

type EncodeHandler func(jsonMsg json.RawMessage) ([]sdk.Msg, error)

type Router interface {
	AddRoute(r string, h EncodeHandler) (rtr Router)
	HasRoute(r string) bool
	GetRoute(path string) (h EncodeHandler)
	Seal()
}

type router struct {
	routes map[string]EncodeHandler
	sealed bool
}

func NewRouter() Router {
	return &router{
		routes: make(map[string]EncodeHandler),
	}
}

func (rtr *router) Seal() {
	if rtr.sealed {
		panic("router already sealed")
	}
	rtr.sealed = true
}

func (rtr *router) AddRoute(path string, h EncodeHandler) Router {
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
func (rtr *router) GetRoute(path string) EncodeHandler {
	if !rtr.HasRoute(path) {
		panic(fmt.Sprintf("route \"%s\" does not exist", path))
	}

	return rtr.routes[path]
}
