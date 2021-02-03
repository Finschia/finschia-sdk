package types

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	yaml "gopkg.in/yaml.v2"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

var (
	_ authexported.GenesisAccount = (*ModuleAccount)(nil)
	_ exported.ModuleAccountI     = (*ModuleAccount)(nil)
)

func init() {
	// Register the ModuleAccount type as a GenesisAccount so that when no
	// concrete GenesisAccount types exist and **default** genesis state is used,
	// the genesis state will serialize correctly.
	authtypes.RegisterAccountTypeCodec(&ModuleAccount{}, "cosmos-sdk/ModuleAccount")
}

var _, accountPrefix = amino.NameToDisfix("cosmos-sdk/ModuleAccount")
var ModuleAccountPrefix = accountPrefix

// ModuleAccount defines an account for modules that holds coins on a pool
type ModuleAccount struct {
	*authtypes.BaseAccount

	Name        string   `json:"name" yaml:"name"`               // name of the module
	Permissions []string `json:"permissions" yaml:"permissions"` // permissions of module account
}

// NewModuleAddress creates an AccAddress from the hash of the module's name
func NewModuleAddress(name string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(name)))
}

// NewEmptyModuleAccount creates a empty ModuleAccount from a string
func NewEmptyModuleAccount(name string, permissions ...string) *ModuleAccount {
	moduleAddress := NewModuleAddress(name)
	baseAcc := authtypes.NewBaseAccountWithAddress(moduleAddress)

	if err := validatePermissions(permissions...); err != nil {
		panic(err)
	}

	return &ModuleAccount{
		BaseAccount: &baseAcc,
		Name:        name,
		Permissions: permissions,
	}
}

// NewModuleAccount creates a new ModuleAccount instance
func NewModuleAccount(ba *authtypes.BaseAccount,
	name string, permissions ...string) *ModuleAccount {

	if err := validatePermissions(permissions...); err != nil {
		panic(err)
	}

	return &ModuleAccount{
		BaseAccount: ba,
		Name:        name,
		Permissions: permissions,
	}
}

// HasPermission returns whether or not the module account has permission.
func (ma ModuleAccount) HasPermission(permission string) bool {
	for _, perm := range ma.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

// GetName returns the the name of the holder's module
func (ma ModuleAccount) GetName() string {
	return ma.Name
}

// GetPermissions returns permissions granted to the module account
func (ma ModuleAccount) GetPermissions() []string {
	return ma.Permissions
}

// SetPubKey - Implements Account
func (ma ModuleAccount) SetPubKey(pubKey crypto.PubKey) error {
	return fmt.Errorf("not supported for module accounts")
}

// SetSequence - Implements Account
func (ma ModuleAccount) SetSequence(seq uint64) error {
	return fmt.Errorf("not supported for module accounts")
}

// Validate checks for errors on the account fields
func (ma ModuleAccount) Validate() error {
	if strings.TrimSpace(ma.Name) == "" {
		return errors.New("module account name cannot be blank")
	}
	if !ma.Address.Equals(sdk.AccAddress(crypto.AddressHash([]byte(ma.Name)))) {
		return fmt.Errorf("address %s cannot be derived from the module name '%s'", ma.Address, ma.Name)
	}

	return ma.BaseAccount.Validate()
}

func (ma *ModuleAccount) MarshalAminoBare(registered bool) (bz []byte, err error) {
	if ma == nil {
		return nil, nil
	}

	buf := bytes.NewBuffer(nil)

	if registered {
		if _, err := buf.Write(accountPrefix[:]); err != nil {
			return nil, err
		}
	}

	bz, err = ma.BaseAccount.MarshalAminoBare(false)
	if err != nil {
		return nil, err
	}
	if err = codec.EncodeFieldByteSlice(buf, 1, bz); err != nil {
		return nil, err
	}

	if err = codec.EncodeFieldByteSlice(buf, 2, []byte(ma.Name)); err != nil {
		return nil, err
	}

	for _, permission := range ma.Permissions {
		// TODO how to hanlde if permission is an empty string?
		if err = codec.EncodeFieldByteSlice(buf, 3, []byte(permission)); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (ma *ModuleAccount) UnmarshalAminoBare(bz []byte) (n int, err error) {
	var _n int
	_n, err = ma.unmarshalBaseAccount(bz)
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 || err != nil {
		return n, err
	}

	var bz2 []byte
	bz2, _n, err = codec.DecodeFieldByteSlice(bz, 2)
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 || err != nil {
		return n, err
	}
	ma.Name = string(bz2)

	_n, err = ma.unmarshalPermissions(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)

	return n, err
}

func (ma *ModuleAccount) unmarshalBaseAccount(bz []byte) (n int, err error) {
	var _n int
	_n, err = codec.CheckFieldNumberAndTyp3(bz, 1, amino.Typ3_ByteLength)
	if _n == 0 || err != nil {
		return n, err
	}
	codec.Slide(&bz, &n, _n)

	var u uint64
	u, _n, err = amino.DecodeUvarint(bz)
	if err != nil {
		return n, err
	}
	codec.Slide(&bz, &n, _n)

	bac := &authtypes.BaseAccount{}
	buf := bz[0:u]
	_n, err = bac.UnmarshalAminoBare(buf)
	if err != nil {
		panic("Fail to unmarshal BaseAccount")
	}
	ma.BaseAccount = bac

	codec.Slide(&bz, &n, _n)

	return n, err
}

func (ma *ModuleAccount) unmarshalPermissions(bz []byte) (n int, err error) {
	var _n int
	var permission string

	for {
		if len(bz) == 0 {
			break
		}

		_n, err = codec.CheckFieldNumberAndTyp3(bz, 3, amino.Typ3_ByteLength)
		if _n == 0 || err != nil {
			return n, err
		}
		codec.Slide(&bz, &n, _n)

		permission, _n, err = amino.DecodeString(bz)
		ma.Permissions = append(ma.Permissions, permission)

		codec.Slide(&bz, &n, _n)
	}

	return n, err
}

type moduleAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
	Name          string         `json:"name" yaml:"name"`
	Permissions   []string       `json:"permissions" yaml:"permissions"`
}

func (ma ModuleAccount) String() string {
	out, _ := ma.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of a ModuleAccount.
func (ma ModuleAccount) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(moduleAccountPretty{
		Address:       ma.Address,
		Coins:         ma.Coins,
		PubKey:        "",
		AccountNumber: ma.AccountNumber,
		Sequence:      ma.Sequence,
		Name:          ma.Name,
		Permissions:   ma.Permissions,
	})

	if err != nil {
		return nil, err
	}

	return string(bs), nil
}

// MarshalJSON returns the JSON representation of a ModuleAccount.
func (ma ModuleAccount) MarshalJSON() ([]byte, error) {
	return codec.Cdc.MarshalJSON(moduleAccountPretty{
		Address:       ma.Address,
		Coins:         ma.Coins,
		PubKey:        "",
		AccountNumber: ma.AccountNumber,
		Sequence:      ma.Sequence,
		Name:          ma.Name,
		Permissions:   ma.Permissions,
	})
}

// UnmarshalJSON unmarshals raw JSON bytes into a ModuleAccount.
func (ma *ModuleAccount) UnmarshalJSON(bz []byte) error {
	var alias moduleAccountPretty
	if err := codec.Cdc.UnmarshalJSON(bz, &alias); err != nil {
		return err
	}

	ma.BaseAccount = authtypes.NewBaseAccount(alias.Address, alias.Coins, nil, alias.AccountNumber, alias.Sequence)
	ma.Name = alias.Name
	ma.Permissions = alias.Permissions

	return nil
}
