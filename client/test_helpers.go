package client

import (
	"fmt"

	cryptotypes "github.com/line/lfb-sdk/crypto/types"
	sdk "github.com/line/lfb-sdk/types"
)

var (
	_ AccountRetriever = TestAccountRetriever{}
	_ Account          = TestAccount{}
)

// TestAccount represents a client Account that can be used in unit tests
type TestAccount struct {
	Address sdk.AccAddress
	Seq     uint64
}

// GetAddress implements client Account.GetAddress
func (t TestAccount) GetAddress() sdk.AccAddress {
	return t.Address
}

// GetPubKey implements client Account.GetPubKey
func (t TestAccount) GetPubKey() cryptotypes.PubKey {
	return nil
}

// GetSequence implements client Account.GetSequence
func (t TestAccount) GetSequence() uint64 {
	return t.Seq
}

// TestAccountRetriever is an AccountRetriever that can be used in unit tests
type TestAccountRetriever struct {
	Accounts map[string]TestAccount
}

// GetAccount implements AccountRetriever.GetAccount
func (t TestAccountRetriever) GetAccount(_ Context, addr sdk.AccAddress) (Account, error) {
	acc, ok := t.Accounts[addr.String()]
	if !ok {
		return nil, fmt.Errorf("account %s not found", addr)
	}

	return acc, nil
}

func (t TestAccountRetriever) GetLatestHeight(_ Context) (uint64, error) {
	return 0, nil
}

// GetAccountWithHeight implements AccountRetriever.GetAccountWithHeight
func (t TestAccountRetriever) GetAccountWithHeight(clientCtx Context, addr sdk.AccAddress) (Account, int64, error) {
	acc, err := t.GetAccount(clientCtx, addr)
	if err != nil {
		return nil, 0, err
	}

	return acc, 0, nil
}

// EnsureExists implements AccountRetriever.EnsureExists
func (t TestAccountRetriever) EnsureExists(_ Context, addr sdk.AccAddress) error {
	_, ok := t.Accounts[addr.String()]
	if !ok {
		return fmt.Errorf("account %s not found", addr)
	}
	return nil
}

// GetAccountSequence implements AccountRetriever.GetAccountSequence
func (t TestAccountRetriever) GetAccountSequence(_ Context, addr sdk.AccAddress) (accSeq uint64, err error) {
	acc, ok := t.Accounts[addr.String()]
	if !ok {
		return 0, fmt.Errorf("account %s not found", addr)
	}
	return acc.Seq, nil
}
