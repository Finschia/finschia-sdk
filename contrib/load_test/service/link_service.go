package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	atypes "github.com/line/link/x/account/client/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	BroadcastBlock = "block"
	BroadcastSync  = "sync"
	BroadcastAsync = "async"
)

type LinkService struct {
	client *http.Client
	cdc    *codec.Codec
	LCDURI string
}

func NewLinkService(client *http.Client, cdc *codec.Codec, lcdURI string) *LinkService {
	return &LinkService{
		client: client,
		cdc:    cdc,
		LCDURI: lcdURI,
	}
}

func (ls *LinkService) GetAccount(from string) (*auth.BaseAccount, error) {
	data, err := ls.get(ls.LCDURI + "/auth/accounts/" + from)
	if err != nil {
		return nil, err
	}

	var res rest.ResponseWithHeight
	err = ls.cdc.UnmarshalJSON(data, &res)
	if err != nil {
		return nil, err
	}

	var account auth.BaseAccount
	err = ls.cdc.UnmarshalJSON(res.Result, &account)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (ls *LinkService) GetLatestBlock() (*coretypes.ResultBlock, error) {
	return ls.getBlock(fmt.Sprintf("%s/blocks/latest", ls.LCDURI))
}

func (ls *LinkService) GetBlock(height int64) (*coretypes.ResultBlock, error) {
	return ls.getBlock(fmt.Sprintf("%s/blocks/%d", ls.LCDURI, height))
}

func (ls *LinkService) getBlock(url string) (*coretypes.ResultBlock, error) {
	data, err := ls.get(url)
	if err != nil {
		return nil, err
	}

	var block coretypes.ResultBlock
	err = ls.cdc.UnmarshalJSON(data, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (ls *LinkService) GetBlocksWithTxResults(fromHeight int64, fetchSize int) ([]*atypes.ResultBlockWithTxResponses,
	error) {
	data, err := ls.get(fmt.Sprintf("%s/blocks_with_tx_results/%d?fetchsize=%d", ls.LCDURI, fromHeight, fetchSize))
	if err != nil {
		return nil, err
	}

	var wrapper atypes.HasMoreResponseWrapper
	err = ls.cdc.UnmarshalJSON(data, &wrapper)
	if err != nil {
		return nil, err
	}

	return wrapper.Items, nil
}

func (ls *LinkService) GetTx(txHash string) (*atypes.TxResponse, error) {
	data, err := ls.get(fmt.Sprintf("%s/txs/%s", ls.LCDURI, txHash))
	if err != nil {
		return nil, err
	}

	var tx atypes.TxResponse
	err = ls.cdc.UnmarshalJSON(data, &tx)
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (ls *LinkService) get(url string) ([]byte, error) {
	resp, err := ls.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, types.RequestFailed{URL: url, Status: resp.Status, Body: data}
	}

	return data, nil
}

func (ls *LinkService) BroadcastTx(stdTx auth.StdTx, mode string) (*atypes.TxResponse, error) {
	url := ls.LCDURI + "/txs"

	bz, err := ls.cdc.MarshalJSON(authrest.BroadcastReq{Mode: mode, Tx: stdTx})
	if err != nil {
		return nil, err
	}

	resp, err := ls.client.Post(url, "application/json", bytes.NewBuffer(bz))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, types.RequestFailed{URL: url, Status: resp.Status, Body: data}
	}

	var res atypes.TxResponse
	err = ls.cdc.UnmarshalJSON(data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (ls *LinkService) BroadcastMsgs(txBuilder transaction.TxBuilderWithoutKeybase, account *auth.BaseAccount,
	keyWallet *wallet.KeyWallet, msgs []sdk.Msg, mode string) (*atypes.TxResponse, error) {
	stdTx, err := txBuilder.WithAccountNumber(account.AccountNumber).WithSequence(account.Sequence).
		BuildAndSign(keyWallet.PrivateKey(), msgs)
	if err != nil {
		return nil, err
	}

	res, err := ls.BroadcastTx(stdTx, mode)
	if err != nil {
		return nil, err
	}

	return res, nil
}
