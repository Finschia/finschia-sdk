package types

import (
	"encoding/json"

	sdk "github.com/line/lbm-sdk/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// Response model for CodeInfo query such as `QueryListCode` and `QueryGetCode`
type CodeInfoResponse interface {
	GetID() uint64
	GetCreator() sdk.AccAddress
	GetDataHash() tmbytes.HexBytes
	GetSource() string
	GetBuilder() string
	GetData() []byte
}

func NewCodeInfoResponse(id uint64, info CodeInfo, data []byte) CodeInfoResponse {
	return codeInfo{
		ID:       id,
		Creator:  info.Creator,
		DataHash: info.CodeHash,
		Source:   info.Source,
		Builder:  info.Builder,
		Data:     data,
	}
}

type codeInfo struct {
	ID       uint64           `json:"id"`
	Creator  sdk.AccAddress   `json:"creator"`
	DataHash tmbytes.HexBytes `json:"data_hash"`
	Source   string           `json:"source"`
	Builder  string           `json:"builder"`
	Data     []byte           `json:"data,omitempty" yaml:"data"` // Data is the entire wasm bytecode
}

func (r codeInfo) GetID() uint64 {
	return r.ID
}

func (r codeInfo) GetCreator() sdk.AccAddress {
	return r.Creator
}

func (r codeInfo) GetDataHash() tmbytes.HexBytes {
	return r.DataHash
}

func (r codeInfo) GetSource() string {
	return r.Source
}

func (r codeInfo) GetBuilder() string {
	return r.Builder
}

func (r codeInfo) GetData() []byte {
	return r.Data
}

// Response model for ContractInfo query such as `QueryListContractByCode` and `QueryGetContract`
type ContractInfoResponse interface {
	GetCodeID() uint64
	GetCreator() sdk.AccAddress
	GetAdmin() sdk.AccAddress
	GetLabel() string
	GetAddress() sdk.AccAddress
	LessThan(other ContractInfoResponse) bool
}

func NewContractInfoResponse(info ContractInfo, addr sdk.AccAddress) ContractInfoResponse {
	return contractInfo{
		CodeID:  info.CodeID,
		Creator: info.Creator,
		Admin:   info.Admin,
		Label:   info.Label,
		Address: addr,
		created: info.Created,
	}
}

// contractInfo adds the address (key) to the ContractInfo representation
type contractInfo struct {
	CodeID  uint64         `json:"code_id"`
	Creator sdk.AccAddress `json:"creator"`
	Admin   sdk.AccAddress `json:"admin,omitempty"`
	Label   string         `json:"label"`
	Address sdk.AccAddress `json:"address"`
	// never show this in query results, just use for sorting
	created *AbsoluteTxPosition
}

func (r contractInfo) GetCodeID() uint64 {
	return r.CodeID
}

func (r contractInfo) GetCreator() sdk.AccAddress {
	return r.Creator
}

func (r contractInfo) GetAdmin() sdk.AccAddress {
	return r.Admin
}

func (r contractInfo) GetLabel() string {
	return r.Label
}

func (r contractInfo) GetAddress() sdk.AccAddress {
	return r.Address
}

func (r contractInfo) LessThan(other ContractInfoResponse) bool {
	return r.created.LessThan(other.(contractInfo).created)
}

// Response model for ContractHistory query (`QueryContractHistory`)
type ContractHistoryResponse interface {
	GetOperation() ContractCodeHistoryOperationType
	GetCodeID() uint64
	GetMsg() json.RawMessage
}

func NewContractHistoryResponse(entry ContractCodeHistoryEntry) ContractHistoryResponse {
	return contractHistory{
		Operation: entry.Operation,
		CodeID:    entry.CodeID,
		Msg:       entry.Msg,
	}
}

type contractHistory struct {
	Operation ContractCodeHistoryOperationType `json:"operation"`
	CodeID    uint64                           `json:"code_id"`
	Msg       json.RawMessage                  `json:"msg,omitempty"`
}

func (r contractHistory) GetOperation() ContractCodeHistoryOperationType {
	return r.Operation
}

func (r contractHistory) GetCodeID() uint64 {
	return r.CodeID
}

func (r contractHistory) GetMsg() json.RawMessage {
	return r.Msg
}

type PageRequest struct {
	Key        []byte
	Offset     uint64
	Limit      uint64
	CountTotal bool
}

type PageResponse struct {
	NextKey []byte
	Total   uint64
}

type QueryContractsByCodeRequest struct {
	CodeID     uint64
	Pagination *PageRequest
}

type QueryCodesRequest struct {
	Pagination *PageRequest
}

type QueryContractHistoryRequest struct {
	Address    sdk.AccAddress
	Pagination *PageRequest
}

type QueryCodesResponse struct {
	CodeInfos  []CodeInfoResponse
	Pagination *PageResponse
}

type QueryContractsByCodeResponse struct {
	ContractInfos []ContractInfoResponse
	Pagination    *PageResponse
}

type QueryContractHistoryResponse struct {
	Entries    []ContractCodeHistoryEntry
	Pagination *PageResponse
}
