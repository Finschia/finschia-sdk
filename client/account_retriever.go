package client

import (
	cryptotypes "github.com/line/lbm-sdk/crypto/types"
	sdk "github.com/line/lbm-sdk/types"
)

// Account defines a read-only version of the auth module's AccountI.
type Account interface {
	GetAddress() sdk.AccAddress
	GetPubKey() cryptotypes.PubKey // can return nil.
	GetSequence() uint64
}

// AccountRetriever defines the interfaces required by transactions to
// ensure an account exists and to be able to query for account fields necessary
// for signing.
type AccountRetriever interface {
	GetAccount(clientCtx Context, addr sdk.AccAddress) (Account, error)
	GetLatestHeight(clientCtx Context) (uint64, error)
	GetAccountWithHeight(clientCtx Context, addr sdk.AccAddress) (Account, int64, error)
	EnsureExists(clientCtx Context, addr sdk.AccAddress) error
	GetAccountSequence(clientCtx Context, addr sdk.AccAddress) (accSeq uint64, err error)
}
