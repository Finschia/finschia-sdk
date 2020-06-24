package master

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/types"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"golang.org/x/sync/errgroup"
)

const (
	Permission = 0644
)

type Controller struct {
	slaves      []types.Slave
	config      types.Config
	params      map[string]map[string]string
	Results     []vegeta.Results
	StartHeight int64
	EndHeight   int64
}

func NewController(slaves []types.Slave, config types.Config, params map[string]map[string]string) *Controller {
	return &Controller{
		slaves:  slaves,
		Results: make([]vegeta.Results, len(slaves)),
		config:  config,
		params:  params,
	}
}

func (c *Controller) StartLoadTest() error {
	log.Println("Request load generator to generate target")
	var eg errgroup.Group
	for _, slave := range c.slaves {
		slave := slave
		eg.Go(func() error {
			return c.orderToLoad(slave)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	log.Println("Request load generator to fire")
	c.StartHeight = getCurrentHeight(c.config.TargetURL)
	for i, slave := range c.slaves {
		i := i
		slave := slave
		eg.Go(func() error {
			return c.orderToFire(slave, i)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	c.EndHeight = getCurrentHeight(c.config.TargetURL)
	return nil
}

func (c *Controller) orderToLoad(slave types.Slave) error {
	config := c.config
	config.Mnemonic = slave.Mnemonic
	url := slave.URL + "/target/load"
	request := types.NewLoadRequest(slave.Scenario, config, c.params[slave.URL])
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return types.RequestFailed{URL: url, Status: resp.Status, Body: data}
	}
	return nil
}

func (c *Controller) orderToFire(slave types.Slave, i int) error {
	url := slave.URL + "/target/fire"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return types.RequestFailed{URL: url, Status: resp.Status, Body: data}
	}
	var res vegeta.Results
	err = json.Unmarshal(data, &res)
	if err != nil {
		return err
	}

	c.Results[i] = res
	return nil
}

func getCurrentHeight(targetURL string) int64 {
	block, err := service.NewLinkService(&http.Client{}, app.MakeCodec(), targetURL).GetLatestBlock()
	for err != nil {
		log.Println("Retrying to get current height")
		block, err = service.NewLinkService(&http.Client{}, app.MakeCodec(), targetURL).GetLatestBlock()
	}
	return block.Block.Height
}
