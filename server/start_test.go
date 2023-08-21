package server

import (
	"github.com/Finschia/ostracon/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGenPvFileOnlyWhenKmsAddressEmptyGenerateFiles(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.PrivValidatorKey = "./key.json"
	cfg.PrivValidatorState = "./state.json"
	defer os.Remove(cfg.PrivValidatorKey)
	defer os.Remove(cfg.PrivValidatorState)

	pv := genPvFileOnlyWhenKmsAddressEmpty(cfg)

	assert.NotNil(t, pv)
}

func TestGenPvFileOnlyWhenKmsAddressEmptyShouldNotGenerateFiles(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.PrivValidatorListenAddr = "tcp://0.0.0.0:26659"

	pv := genPvFileOnlyWhenKmsAddressEmpty(cfg)

	assert.Nil(t, pv)
}
