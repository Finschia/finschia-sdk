package types

import (
	"encoding/hex"
	"fmt"
	"math"
	"strings"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

type TxResponse struct {
	Height    int64               `json:"height"`
	TxHash    string              `json:"txhash"`
	Codespace string              `json:"codespace,omitempty"`
	Code      uint32              `json:"code,omitempty"`
	Index     uint32              `json:"index"` // additional field
	Data      string              `json:"data,omitempty"`
	RawLog    string              `json:"raw_log,omitempty"`
	Logs      sdk.ABCIMessageLogs `json:"logs,omitempty"`
	Info      string              `json:"info,omitempty"`
	GasWanted int64               `json:"gas_wanted,omitempty"`
	GasUsed   int64               `json:"gas_used,omitempty"`
	Tx        sdk.Tx              `json:"tx,omitempty"`
	Timestamp string              `json:"timestamp,omitempty"`
}

// NewResponseResultTx returns a TxResponse given a ResultTx from tendermint
func NewResponseResultTx(res *ctypes.ResultTx, tx sdk.Tx, timestamp string) TxResponse {
	if res == nil {
		return TxResponse{}
	}

	parsedLogs, err := sdk.ParseABCILogs(res.TxResult.Log)
	if err != nil {
		parsedLogs = nil
	}

	return TxResponse{
		TxHash:    res.Hash.String(),
		Height:    res.Height,
		Codespace: res.TxResult.Codespace,
		Code:      res.TxResult.Code,
		Index:     res.Index,
		Data:      strings.ToUpper(hex.EncodeToString(res.TxResult.Data)),
		RawLog:    res.TxResult.Log,
		Logs:      parsedLogs,
		Info:      res.TxResult.Info,
		GasWanted: res.TxResult.GasWanted,
		GasUsed:   res.TxResult.GasUsed,
		Tx:        tx,
		Timestamp: timestamp,
	}
}

func (r TxResponse) String() string {
	var sb strings.Builder
	sb.WriteString("Response:\n")

	if r.Height > 0 {
		sb.WriteString(fmt.Sprintf("  Height: %d\n", r.Height))
		sb.WriteString(fmt.Sprintf("  Index: %d\n", r.Index))
	}
	if r.TxHash != "" {
		sb.WriteString(fmt.Sprintf("  TxHash: %s\n", r.TxHash))
	}
	if r.Codespace != "" {
		sb.WriteString(fmt.Sprintf("  Codespace: %s\n", r.Codespace))
	}
	if r.Code > 0 {
		sb.WriteString(fmt.Sprintf("  Code: %d\n", r.Code))
	}
	if r.Data != "" {
		sb.WriteString(fmt.Sprintf("  Data: %s\n", r.Data))
	}
	if r.RawLog != "" {
		sb.WriteString(fmt.Sprintf("  Raw Log: %s\n", r.RawLog))
	}
	if r.Logs != nil {
		sb.WriteString(fmt.Sprintf("  Logs: %s\n", r.Logs))
	}
	if r.Info != "" {
		sb.WriteString(fmt.Sprintf("  Info: %s\n", r.Info))
	}
	if r.GasWanted != 0 {
		sb.WriteString(fmt.Sprintf("  GasWanted: %d\n", r.GasWanted))
	}
	if r.GasUsed != 0 {
		sb.WriteString(fmt.Sprintf("  GasUsed: %d\n", r.GasUsed))
	}
	if r.Timestamp != "" {
		sb.WriteString(fmt.Sprintf("  Timestamp: %s\n", r.Timestamp))
	}

	return strings.TrimSpace(sb.String())
}

// Empty returns true if the response is empty
func (r TxResponse) Empty() bool {
	return r.TxHash == "" && r.Logs == nil
}

// SearchTxsResult defines a structure for querying txs pageable
type SearchTxsResult struct {
	TotalCount int          `json:"total_count"` // Count of all txs
	Count      int          `json:"count"`       // Count of txs in current page
	PageNumber int          `json:"page_number"` // Index of current page, start from 1
	PageTotal  int          `json:"page_total"`  // Count of total pages
	Limit      int          `json:"limit"`       // Max count txs per page
	Txs        []TxResponse `json:"txs"`         // List of txs in current page
}

func NewSearchTxsResult(totalCount, count, page, limit int, txs []TxResponse) SearchTxsResult {
	return SearchTxsResult{
		TotalCount: totalCount,
		Count:      count,
		PageNumber: page,
		PageTotal:  int(math.Ceil(float64(totalCount) / float64(limit))),
		Limit:      limit,
		Txs:        txs,
	}
}

type SearchGenesisAccountResult struct {
	TotalCount int                      `json:"total_count"`
	Count      int                      `json:"count"`
	PageNumber int                      `json:"page_number"`
	PageTotal  int                      `json:"page_total"`
	Limit      int                      `json:"limit"`
	Accounts   exported.GenesisAccounts `json:"accounts"`
}

func NewSearchGenesisAccountResult(totalCount, count, page, limit int, accounts exported.GenesisAccounts) SearchGenesisAccountResult {
	return SearchGenesisAccountResult{
		TotalCount: totalCount,
		Count:      count,
		PageNumber: page,
		PageTotal:  int(math.Ceil(float64(totalCount) / float64(limit))),
		Limit:      limit,
		Accounts:   accounts,
	}
}
