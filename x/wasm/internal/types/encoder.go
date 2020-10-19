package types

import (
	"encoding/json"
)

type EncodingModule string

const (
	TokenM      = EncodingModule("token")
	CollectionM = EncodingModule("collection")
)

type LinkMsgWrapper struct {
	Module  string          `json:"module"`
	MsgData json.RawMessage `json:"msg_data"`
}

type LinkQueryWrapper struct {
	Module    string          `json:"module"`
	QueryData json.RawMessage `json:"query_data"`
}
