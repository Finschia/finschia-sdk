package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	cutils "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	gtypes "github.com/cosmos/cosmos-sdk/x/genutil/types"

	local "github.com/line/lbm-sdk/v2/x/account/client/types"
)

const MaxPerPage = 100

// QueryTxsByEvents performs a search for transactions for a given set of events
// via the Tendermint RPC. An event takes the form of:
// "{eventAttribute}.{attributeKey} = '{attributeValue}'". Each event is
// concatenated with an 'AND' operand. It returns a slice of Info object
// containing txs and metadata. An error is returned if the query fails.
func QueryTxsByEvents(cliCtx context.CLIContext, events []string, page, limit int) (*local.SearchTxsResult, error) {
	if len(events) == 0 {
		return nil, errors.New("must declare at least one event to search")
	}

	if page <= 0 {
		return nil, errors.New("page must greater than 0")
	}

	if limit <= 0 {
		return nil, errors.New("limit must greater than 0")
	}

	// XXX: implement ANY
	query := strings.Join(events, " AND ")

	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	prove := !cliCtx.TrustNode

	resTxs, err := node.TxSearch(query, prove, page, limit, "")
	if err != nil {
		return nil, err
	}

	if prove {
		for _, tx := range resTxs.Txs {
			err := cutils.ValidateTxResult(cliCtx, tx)
			if err != nil {
				return nil, err
			}
		}
	}

	resBlocks, err := getBlocksForTxResults(cliCtx, resTxs.Txs)
	if err != nil {
		return nil, err
	}

	txs, err := formatTxResults(cliCtx.Codec, resTxs.Txs, resBlocks)
	if err != nil {
		return nil, err
	}

	result := local.NewSearchTxsResult(resTxs.TotalCount, len(txs), page, limit, txs)

	return &result, nil
}

// QueryTx queries for a single transaction by a hash string in hex format. An
// error is returned if the transaction does not exist or cannot be queried.
func QueryTx(cliCtx context.CLIContext, hashHexStr string) (local.TxResponse, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return local.TxResponse{}, err
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return local.TxResponse{}, err
	}

	resTx, err := node.Tx(hash, !cliCtx.TrustNode)
	if err != nil {
		return local.TxResponse{}, err
	}

	if !cliCtx.TrustNode {
		if err = cutils.ValidateTxResult(cliCtx, resTx); err != nil {
			return local.TxResponse{}, err
		}
	}

	resBlocks, err := getBlocksForTxResults(cliCtx, []*ctypes.ResultTx{resTx})
	if err != nil {
		return local.TxResponse{}, err
	}

	out, err := formatTxResult(cliCtx.Codec, resTx, resBlocks[resTx.Height])
	if err != nil {
		return out, err
	}

	return out, nil
}

func QueryGenesisTx(cliCtx context.CLIContext) ([]sdk.Tx, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return []sdk.Tx{}, nil
	}

	resultGenesis, err := node.Genesis()
	if err != nil {
		return []sdk.Tx{}, err
	}

	appState, err := gtypes.GenesisStateFromGenDoc(cliCtx.Codec, *resultGenesis.Genesis)
	if err != nil {
		return []sdk.Tx{}, err
	}

	genState := gtypes.GetGenesisStateFromAppState(cliCtx.Codec, appState)
	genTxs := make([]sdk.Tx, len(genState.GenTxs))
	for i, tx := range genState.GenTxs {
		err := cliCtx.Codec.UnmarshalJSON(tx, &genTxs[i])
		if err != nil {
			return []sdk.Tx{}, err
		}
	}
	return genTxs, nil
}

func QueryGenesisAccount(cliCtx context.CLIContext, page, perPage int) (local.SearchGenesisAccountResult, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return local.SearchGenesisAccountResult{}, err
	}

	resultGenesis, err := node.Genesis()
	if err != nil {
		return local.SearchGenesisAccountResult{}, err
	}

	appState, err := gtypes.GenesisStateFromGenDoc(cliCtx.Codec, *resultGenesis.Genesis)
	if err != nil {
		return local.SearchGenesisAccountResult{}, err
	}

	genAccounts := types.GetGenesisStateFromAppState(cliCtx.Codec, appState).Accounts
	totalCount := len(genAccounts)

	perPage, err = validatePerPage(perPage)
	if err != nil {
		return local.SearchGenesisAccountResult{}, err
	}

	page, err = validatePage(page, perPage, totalCount)
	if err != nil {
		return local.SearchGenesisAccountResult{}, err
	}
	start, end := getCountIndexRange(page, perPage, totalCount)
	resultAccounts := genAccounts[start:end]

	return local.NewSearchGenesisAccountResult(totalCount, len(resultAccounts), page, perPage, resultAccounts), nil
}

func validatePage(page, perPage, totalCount int) (int, error) {
	if page < 1 {
		return 1, fmt.Errorf("the page must greater than 0")
	}

	pages := ((totalCount - 1) / perPage) + 1
	if pages == 0 {
		pages = 1
	}
	if page < 0 || page > pages {
		return 1, fmt.Errorf("the page should be within [1, %d] range, given %d", pages, page)
	}

	return page, nil
}

func validatePerPage(perPage int) (int, error) {
	if perPage < 1 {
		return 1, fmt.Errorf("the limit must greater than 0")
	}

	if perPage > MaxPerPage {
		return MaxPerPage, nil
	}
	return perPage, nil
}

func getCountIndexRange(page, perPage, totalCount int) (int, int) {
	start := (page - 1) * perPage
	end := start + perPage
	if start < 0 {
		return 0, end
	}
	if end > totalCount {
		end = totalCount
	}

	return start, end
}

// formatTxResults parses the indexed txs into a slice of TxResponse objects.
func formatTxResults(cdc *codec.Codec, resTxs []*ctypes.ResultTx, resBlocks map[int64]*ctypes.ResultBlock) ([]local.TxResponse, error) {
	var err error
	out := make([]local.TxResponse, len(resTxs))
	for i := range resTxs {
		out[i], err = formatTxResult(cdc, resTxs[i], resBlocks[resTxs[i].Height])
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}

func getBlocksForTxResults(cliCtx context.CLIContext, resTxs []*ctypes.ResultTx) (map[int64]*ctypes.ResultBlock, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	resBlocks := make(map[int64]*ctypes.ResultBlock)

	for _, resTx := range resTxs {
		if _, ok := resBlocks[resTx.Height]; !ok {
			resBlock, err := node.Block(&resTx.Height)
			if err != nil {
				return nil, err
			}

			resBlocks[resTx.Height] = resBlock
		}
	}

	return resBlocks, nil
}

func formatTxResult(cdc *codec.Codec, resTx *ctypes.ResultTx, resBlock *ctypes.ResultBlock) (local.TxResponse, error) {
	tx, err := parseTx(cdc, resTx.Tx)
	if err != nil {
		return local.TxResponse{}, err
	}

	return local.NewResponseResultTx(resTx, tx, resBlock.Block.Time.Format(time.RFC3339)), nil
}

func parseTx(cdc *codec.Codec, txBytes []byte) (sdk.Tx, error) {
	var tx types.StdTx

	err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// ParseHTTPArgs parses the request's URL and returns a slice containing all
// arguments pairs. It separates page and limit used for pagination.
func ParseHTTPArgs(r *http.Request) (tags []string, page, limit int, err error) {
	tags = make([]string, 0, len(r.Form))
	for key, values := range r.Form {
		if key == "page" || key == "limit" || key == "height.from" || key == "height.to" {
			continue
		}

		var value string
		value, err = url.QueryUnescape(values[0])
		if err != nil {
			return tags, page, limit, err
		}

		var tag string
		if key == tmtypes.TxHeightKey {
			tag = fmt.Sprintf("%s=%s", key, value)
		} else {
			tag = fmt.Sprintf("%s='%s'", key, value)
		}
		tags = append(tags, tag)
	}

	heightFromStr := r.FormValue("height.from")
	if heightFromStr != "" {
		heightFrom, err := strconv.ParseInt(heightFromStr, 10, 64)
		switch {
		case err != nil:
			return tags, page, limit, err
		case heightFrom <= 0:
			return tags, page, limit, errors.New("height.from must greater than 0")
		default:
			tags = append(tags, fmt.Sprintf("%s>=%d", tmtypes.TxHeightKey, heightFrom))
		}
	}

	heightToStr := r.FormValue("height.to")
	if heightToStr != "" {
		heightTo, err := strconv.ParseInt(heightToStr, 10, 64)
		switch {
		case err != nil:
			return tags, page, limit, err
		case heightTo <= 0:
			return tags, page, limit, errors.New("height.to must greater than 0")
		default:
			tags = append(tags, fmt.Sprintf("%s<=%d", tmtypes.TxHeightKey, heightTo))
		}
	}

	if len(tags) == 0 {
		return tags, page, limit, errors.New("must declare at least one event to search")
	}

	pageStr := r.FormValue("page")
	if pageStr == "" {
		page = rest.DefaultPage
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return tags, page, limit, err
		} else if page <= 0 {
			return tags, page, limit, errors.New("page must greater than 0")
		}
	}

	limitStr := r.FormValue("limit")
	if limitStr == "" {
		limit = rest.DefaultLimit
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return tags, page, limit, err
		} else if limit <= 0 {
			return tags, page, limit, errors.New("limit must greater than 0")
		}
	}

	return tags, page, limit, nil
}
