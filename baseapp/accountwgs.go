package baseapp

import (
	"sync"

	sdk "github.com/line/lfb-sdk/types"
)

type AccountWGs struct {
	mtx sync.Mutex
	wgs map[string]*sync.WaitGroup
}

func NewAccountWGs() *AccountWGs {
	return &AccountWGs{
		wgs: make(map[string]*sync.WaitGroup),
	}
}

func (aw *AccountWGs) Register(tx sdk.Tx) (waits []*sync.WaitGroup, signals []*AccountWG) {
	signers := getUniqSigners(tx)

	aw.mtx.Lock()
	defer aw.mtx.Unlock()
	for _, signer := range signers {
		if wg := aw.wgs[signer]; wg != nil {
			waits = append(waits, wg)
		}
		sig := waitGroup1()
		aw.wgs[signer] = sig
		signals = append(signals, NewAccountWG(signer, sig))
	}

	return waits, signals
}

func (aw *AccountWGs) Wait(waits []*sync.WaitGroup) {
	for _, wait := range waits {
		wait.Wait()
	}
}

func (aw *AccountWGs) Done(signals []*AccountWG) {
	aw.mtx.Lock()
	defer aw.mtx.Unlock()

	for _, signal := range signals {
		signal.wg.Done()
		if aw.wgs[signal.acc] == signal.wg {
			delete(aw.wgs, signal.acc)
		}
	}
}

func getUniqSigners(tx sdk.Tx) []string {
	seen := map[string]bool{}
	var signers []string
	for _, msg := range tx.GetMsgs() {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, string(addr))
				seen[addr.String()] = true
			}
		}
	}
	return signers
}

type AccountWG struct {
	acc string
	wg  *sync.WaitGroup
}

func NewAccountWG(acc string, wg *sync.WaitGroup) *AccountWG {
	return &AccountWG{
		acc: acc,
		wg:  wg,
	}
}

func waitGroup1() (wg *sync.WaitGroup) {
	wg = &sync.WaitGroup{}
	wg.Add(1)
	return wg
}
