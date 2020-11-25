package types

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/go-amino"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/tendermint/tendermint/crypto"

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

func (ma *ModuleAccount) Marshal() (bz []byte, err error) {
	buf := bytes.NewBuffer(nil)

	if _, err = buf.Write([]byte{11, 253, 125, 154}); err != nil {
		return
	}

	if ma.BaseAccount != nil {
		bz, err = ma.BaseAccount.Marshal()
		if err != nil {
			return
		}
		if err = codec.EncodeFieldNumberAndTyp3(buf, 1, amino.Typ3_ByteLength); err != nil {
			return
		}
		// NOTE skip prefix
		if err = amino.EncodeByteSlice(buf, bz[4:]); err != nil {
			return
		}
	}
	if ma.Name != "" {
		if err = codec.EncodeFieldNumberAndTyp3(buf, 2, amino.Typ3_ByteLength); err != nil {
			return
		}
		if err = amino.EncodeString(buf, ma.Name); err != nil {
			return
		}
	}
	for _, permission := range ma.Permissions {
		// TODO how to hanlde if permission is an empty string?
		if err = codec.EncodeFieldNumberAndTyp3(buf, 3, amino.Typ3_ByteLength); err != nil {
			return
		}
		if err = amino.EncodeString(buf, permission); err != nil {
			return
		}
	}

	return buf.Bytes(), nil
}

func (ma *ModuleAccount) Unmarshal(bz []byte) (n int, err error) {
	var _n int
	_n, err = ma.decodeBaseAccount(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 {
		return
	}

	_n, err = ma.decodeName(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)
	if len(bz) == 0 {
		return
	}

	_n, err = ma.decodePermissions(bz)
	if err != nil {
		return
	}
	codec.Slide(&bz, &n, _n)

	return
}

func (ma *ModuleAccount) decodeBaseAccount(bz []byte) (n int, err error) {
	var _n int
	_n, err = codec.CheckFieldNumberAndTyp3(bz, 1, amino.Typ3_ByteLength)
	codec.Slide(&bz, &n, _n)

	var u uint64
	u, _n, err = amino.DecodeUvarint(bz)
	codec.Slide(&bz, &n, _n)

	buf := bz[0:u]

	bac := &authtypes.BaseAccount{}
	_n, err = bac.Unmarshal(buf)
	if err != nil {
		panic("Fail to unmarshal BaseAccount")
	}
	codec.Slide(&bz, &n, _n)

	ma.BaseAccount = bac

	return
}

func (ma *ModuleAccount) decodeName(bz []byte) (n int, err error) {
	var _n int
	_n, err = codec.CheckFieldNumberAndTyp3(bz, 2, amino.Typ3_ByteLength)
	codec.Slide(&bz, &n, _n)

	ma.Name, _n, err = amino.DecodeString(bz)
	codec.Slide(&bz, &n, _n)

	return
}

func (ma *ModuleAccount) decodePermissions(bz []byte) (n int, err error) {
	var _n int
	var permission string

	for {
		if len(bz) == 0 {
			break
		}
		_n, err = codec.CheckFieldNumberAndTyp3(bz, 3, amino.Typ3_ByteLength)
		if _n == 0 || err != nil {
			break
		}
		codec.Slide(&bz, &n, _n)

		// REFACTOR
		if ma.Permissions == nil {
			ma.Permissions = make([]string, 0)
		}

		permission, _n, err = amino.DecodeString(bz)
		ma.Permissions = append(ma.Permissions, permission)

		codec.Slide(&bz, &n, _n)
	}

	return
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
