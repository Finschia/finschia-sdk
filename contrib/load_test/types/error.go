package types

import (
	"fmt"
	"time"

	atypes "github.com/line/link/x/account/client/types"
)

type InaccessibleFieldError string

func (e InaccessibleFieldError) Error() string {
	return fmt.Sprintf("Inaccessible Field Error: %s", string(e))
}

type InvalidLoadParameterError string

func (e InvalidLoadParameterError) Error() string {
	return fmt.Sprintf("Invalid Load Parameter Error: %s", string(e))
}

type InvalidScenarioError string

func (e InvalidScenarioError) Error() string {
	return fmt.Sprintf("Invalid Scenario Error: %s", string(e))
}

type InvalidMasterMnemonic struct {
	Mnemonic string
}

func (e InvalidMasterMnemonic) Error() string {
	return fmt.Sprintf("Invalid master mnemonic: %s", e.Mnemonic)
}

type InvalidMnemonic struct {
	Mnemonic string
}

func (e InvalidMnemonic) Error() string {
	return fmt.Sprintf("Invalid mnemonic: %s", e.Mnemonic)
}

type RequestFailed struct {
	URL    string
	Status string
	Body   []byte
}

func (e RequestFailed) Error() string {
	return fmt.Sprintf("Request failed: %s %s\n%s", e.URL, e.Status, e.Body)
}

type FailedTxError struct {
	Tx *atypes.TxResponse
}

func (e FailedTxError) Error() string {
	return fmt.Sprintf("Tx failed: %v\n", e.Tx)
}

type NoContractIDError struct {
	Tx *atypes.TxResponse
}

func (e NoContractIDError) Error() string {
	return fmt.Sprintf("There is no contract id in events. Tx: %v\n", e.Tx)
}

type FailedToCreateFile struct {
	Err error
}

func (e FailedToCreateFile) Error() string {
	return fmt.Sprintf("Failed to create file:%s", e.Err.Error())
}

type FailedToRenderGraph struct {
	Err error
}

func (e FailedToRenderGraph) Error() string {
	return fmt.Sprintf("Failed to render graph:%s", e.Err.Error())
}

type InvalidModuleName struct {
	Name string
}

func (e InvalidModuleName) Error() string {
	return fmt.Sprintf("Invalid module name : %s", e.Name)
}

type HighLatencyError struct {
	Latency   time.Duration
	Threshold time.Duration
}

func (e HighLatencyError) Error() string {
	return fmt.Sprintf("Latency is higher than %s : %s", e.Threshold, e.Latency)
}

type LowTPSError struct {
	TPS       float64
	Threshold float64
}

func (e LowTPSError) Error() string {
	return fmt.Sprintf("TPS is lower than %f : %f", e.Threshold, e.TPS)
}

type LowThroughputError struct {
	Throughput float64
	Threshold  float64
}

func (e LowThroughputError) Error() string {
	return fmt.Sprintf("Throughput is lower than %f : %f", e.Threshold, e.Throughput)
}
