package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"time"

	"github.com/tendermint/tendermint/crypto"
	yaml "gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

//-----------------------------------------------------------------------------
// BaseAccount

var _ exported.Account = (*BaseAccount)(nil)
var _ exported.GenesisAccount = (*BaseAccount)(nil)

// BaseAccount - a base account structure.
// This can be extended by embedding within in your AppAccount.
// However one doesn't have to use BaseAccount as long as your struct
// implements Account.
type BaseAccount struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        crypto.PubKey  `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
}

// NewBaseAccount creates a new BaseAccount object
func NewBaseAccount(address sdk.AccAddress, coins sdk.Coins,
	pubKey crypto.PubKey, accountNumber uint64, sequence uint64) *BaseAccount {

	return &BaseAccount{
		Address:       address,
		Coins:         coins,
		PubKey:        pubKey,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
}

// ProtoBaseAccount - a prototype function for BaseAccount
func ProtoBaseAccount() exported.Account {
	return &BaseAccount{}
}

// NewBaseAccountWithAddress - returns a new base account with a given address
func NewBaseAccountWithAddress(addr sdk.AccAddress) BaseAccount {
	return BaseAccount{
		Address: addr,
	}
}

// GetAddress - Implements sdk.Account.
func (acc BaseAccount) GetAddress() sdk.AccAddress {
	return acc.Address
}

// SetAddress - Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr sdk.AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}

// GetPubKey - Implements sdk.Account.
func (acc BaseAccount) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

// SetPubKey - Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

// GetCoins - Implements sdk.Account.
func (acc *BaseAccount) GetCoins() sdk.Coins {
	return acc.Coins
}

// SetCoins - Implements sdk.Account.
func (acc *BaseAccount) SetCoins(coins sdk.Coins) error {
	acc.Coins = coins
	return nil
}

// GetAccountNumber - Implements Account
func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements Account
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements sdk.Account.
func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements sdk.Account.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// SpendableCoins returns the total set of spendable coins. For a base account,
// this is simply the base coins.
func (acc *BaseAccount) SpendableCoins(_ time.Time) sdk.Coins {
	return acc.GetCoins()
}

// Validate checks for errors on the account fields
func (acc BaseAccount) Validate() error {
	if acc.PubKey != nil && acc.Address != nil &&
		!bytes.Equal(acc.PubKey.Address().Bytes(), acc.Address.Bytes()) {
		return errors.New("pubkey and address pair is invalid")
	}

	return nil
}

func (acc *BaseAccount) Unmarshal(bz []byte) (n int, err error) {
	var _n int
	_n, err = acc.decodeAddress(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 {
		return
	}

	_n, err = acc.decodeCoins(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 {
		return
	}

	_n, err = acc.decodePubKey(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 {
		return
	}

	_n, err = acc.decodeAccountNumber(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 {
		return
	}

	_n, err = acc.decodeSequence(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)

	return
}

func (acc *BaseAccount) decodeAddress(bz []byte) (n int, err error) {
	_n, err := codec.CheckFieldNumberAndTyp3(bz, 1, amino.Typ3_ByteLength)
	if _n == 0 || err != nil {
		return _n, err
	}
	codec.Slide(&bz, &n, _n)

	acc.Address, _n, err = amino.DecodeByteSlice(bz)
	codec.Slide(&bz, &n, _n)

	return
}

func (acc *BaseAccount) decodeCoins(bz []byte) (n int, err error) {
	var _n int

	for {
		if len(bz) == 0 {
			break
		}
		_n, err = codec.CheckFieldNumberAndTyp3(bz, 2, amino.Typ3_ByteLength)
		if _n == 0 || err != nil {
			break
		}
		codec.Slide(&bz, &n, _n)

		// REFACTOR
		if acc.Coins == nil {
			acc.Coins = make(sdk.Coins, 0)
		}

		var u uint64
		u, _n, err = amino.DecodeUvarint(bz)
		codec.Slide(&bz, &n, _n)

		buf := bz[0:u]
		_n, err = acc.decodeCoin(buf)

		codec.Slide(&bz, &n, _n)
	}

	return
}

func (acc *BaseAccount) decodeCoin(bz []byte) (n int, err error) {
	var _n int
	_n, err = codec.CheckFieldNumberAndTyp3(bz, 1, amino.Typ3_ByteLength)
	codec.Slide(&bz, &n, _n)

	var denom, amt string
	denom, _n, err = amino.DecodeString(bz)
	codec.Slide(&bz, &n, _n)

	_n, err = codec.CheckFieldNumberAndTyp3(bz, 2, amino.Typ3_ByteLength)
	codec.Slide(&bz, &n, _n)

	amt, _n, err = amino.DecodeString(bz)
	codec.Slide(&bz, &n, _n)

	amount, _ := sdk.NewIntFromString(amt)

	// acc.Coins.Add(sdk.Coin{Denom: denom, Amount: amount})
	acc.Coins = append(acc.Coins, sdk.Coin{Denom: denom, Amount: amount})
	return n, nil
}

func (acc *BaseAccount) decodePubKey(bz []byte) (n int, err error) {
	var _n int
	_n, err = codec.CheckFieldNumberAndTyp3(bz, 3, amino.Typ3_ByteLength)
	if _n == 0 || err != nil {
		return _n, err
	}
	codec.Slide(&bz, &n, _n)

	var u uint64
	u, _n, err = amino.DecodeUvarint(bz)
	codec.Slide(&bz, &n, _n)

	buf := bz[0:u]
	_n, err = acc.decodePubKeyInterface(buf)
	codec.Slide(&bz, &n, _n)
	return
}

func (acc *BaseAccount) decodePubKeyInterface(bz []byte) (n int, err error) {
	var _n int
	// disamb, hasDisamb, prefix, hasPrefix, _n, err := amino.DecodeDisambPrefixBytes(bz)
	_, _, _, _, _n, err = amino.DecodeDisambPrefixBytes(bz)
	if codec.Slide(&bz, &n, _n) && err != nil {
		return
	}

	var byteslice []byte
	byteslice, _n, err = amino.DecodeByteSlice(bz)
	if codec.Slide(&bz, &n, _n) && err != nil {
		return
	}

	var pubKey secp256k1.PubKeySecp256k1
	copy(pubKey[:], byteslice)
	acc.PubKey = pubKey
	return
}

func (acc *BaseAccount) decodeAccountNumber(bz []byte) (n int, err error) {
	var _n int
	_n, err = codec.CheckFieldNumberAndTyp3(bz, 4, amino.Typ3_Varint)
	if _n == 0 || err != nil {
		return _n, err
	}
	codec.Slide(&bz, &n, _n)

	u, _n, err := amino.DecodeUvarint(bz)
	acc.AccountNumber = u
	codec.Slide(&bz, &n, _n)
	return
}

func (acc *BaseAccount) decodeSequence(bz []byte) (n int, err error) {
	var _n int
	_n, err = codec.CheckFieldNumberAndTyp3(bz, 5, amino.Typ3_Varint)
	if _n == 0 || err != nil {
		return _n, err
	}
	codec.Slide(&bz, &n, _n)

	u, _n, err := amino.DecodeUvarint(bz)
	acc.Sequence = u
	codec.Slide(&bz, &n, _n)
	return
}

type baseAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
}

func (acc BaseAccount) String() string {
	out, _ := acc.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of an account.
func (acc BaseAccount) MarshalYAML() (interface{}, error) {
	alias := baseAccountPretty{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}

	if acc.PubKey != nil {
		pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, acc.PubKey)
		if err != nil {
			return nil, err
		}

		alias.PubKey = pks
	}

	bz, err := yaml.Marshal(alias)
	if err != nil {
		return nil, err
	}

	return string(bz), err
}

// MarshalJSON returns the JSON representation of a BaseAccount.
func (acc BaseAccount) MarshalJSON() ([]byte, error) {
	alias := baseAccountPretty{
		Address:       acc.Address,
		Coins:         acc.Coins,
		AccountNumber: acc.AccountNumber,
		Sequence:      acc.Sequence,
	}

	if acc.PubKey != nil {
		pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, acc.PubKey)
		if err != nil {
			return nil, err
		}

		alias.PubKey = pks
	}

	return json.Marshal(alias)
}

// UnmarshalJSON unmarshals raw JSON bytes into a BaseAccount.
func (acc *BaseAccount) UnmarshalJSON(bz []byte) error {
	var alias baseAccountPretty
	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}

	if alias.PubKey != "" {
		pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, alias.PubKey)
		if err != nil {
			return err
		}

		acc.PubKey = pk
	}

	acc.Address = alias.Address
	acc.Coins = alias.Coins
	acc.AccountNumber = alias.AccountNumber
	acc.Sequence = alias.Sequence

	return nil
}
