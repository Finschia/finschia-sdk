package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/line/link/contrib/load_test/types"
	authtypes "github.com/line/link/x/auth/client/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
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
	url := ls.LCDURI + "/auth/accounts/" + from
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
	url := fmt.Sprintf("%s/blocks/latest", ls.LCDURI)
	return ls.getBlock(url)
}

func (ls *LinkService) GetBlock(height int64) (*coretypes.ResultBlock, error) {
	url := fmt.Sprintf("%s/blocks/%d", ls.LCDURI, height)
	return ls.getBlock(url)
}

func (ls *LinkService) getBlock(url string) (*coretypes.ResultBlock, error) {
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

	var block coretypes.ResultBlock
	err = ls.cdc.UnmarshalJSON(data, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (ls *LinkService) BroadcastTx(stdTx auth.StdTx, mode string) (*authtypes.TxResponse, error) {
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

	var res authtypes.TxResponse
	err = ls.cdc.UnmarshalJSON(data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
