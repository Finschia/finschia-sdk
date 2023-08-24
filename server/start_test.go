package server

import (
	"github.com/Finschia/ostracon/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGenPvFileOnlyWhenKmsAddressEmptyGenerateFiles(t *testing.T) {
	cfg := config.DefaultConfig()
	dir, err := os.MkdirTemp("", "start_test")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.RemoveAll(dir)
	cfg.PrivValidatorKey = dir + "/key.json"
	cfg.PrivValidatorState = dir + "/state.json"

	pv := genPvFileOnlyWhenKmsAddressEmpty(cfg)

	assert.NotNil(t, pv)
}

func TestGenPvFileOnlyWhenKmsAddressEmptyShouldNotGenerateFiles(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.PrivValidatorListenAddr = "tcp://0.0.0.0:26659"

	pv := genPvFileOnlyWhenKmsAddressEmpty(cfg)

	assert.Nil(t, pv)
}
