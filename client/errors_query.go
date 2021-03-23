package client

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

type Error struct {
	Codespace string `json:"codespace"`
	Code      uint32 `json:"code"`
	Message   string `json:"message"`
}

func NewQueryError(codespace string, code uint32, desc string) *Error {
	return &Error{Codespace: codespace, Code: code, Message: desc}
}

func (err Error) Error() string {
	var buff bytes.Buffer
	enc := json.NewEncoder(&buff)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(err); err != nil {
		panic(errors.Wrap(err, "failed to encode Query error log"))
	}

	return strings.TrimSpace(buff.String())
}
