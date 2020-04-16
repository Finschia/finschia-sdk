package master

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/line/link/contrib/load_test/types"
	"golang.org/x/sync/errgroup"
)

const (
	Permission = 0644
)

type Controller struct {
	slaves  []types.Slave
	config  types.Config
	Results [][]byte
}

func NewController(slaves []types.Slave, config types.Config) *Controller {
	return &Controller{
		slaves:  slaves,
		Results: make([][]byte, len(slaves)),
		config:  config,
	}
}

func (c *Controller) StartLoadTest() error {
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
	return nil
}

func (c *Controller) orderToLoad(slave types.Slave) error {
	config := c.config
	config.Mnemonic = slave.Mnemonic
	url := slave.URL + "/target/load"
	request := types.NewLoadRequest(slave.TargetType, config)
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

	c.Results[i] = data
	return nil
}
