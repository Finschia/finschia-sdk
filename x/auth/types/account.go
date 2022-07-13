package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/line/ostracon/crypto"
	"gopkg.in/yaml.v2"

	"github.com/line/lbm-sdk/codec"
	codectypes "github.com/line/lbm-sdk/codec/types"
	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/crypto/keys/multisig"
	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	"github.com/line/lbm-sdk/crypto/keys/secp256r1"
	cryptotypes "github.com/line/lbm-sdk/crypto/types"
	sdk "github.com/line/lbm-sdk/types"
)

var (
	_ AccountI                           = (*BaseAccount)(nil)
	_ GenesisAccount                     = (*BaseAccount)(nil)
	_ codectypes.UnpackInterfacesMessage = (*BaseAccount)(nil)
	_ GenesisAccount                     = (*ModuleAccount)(nil)
	_ ModuleAccountI                     = (*ModuleAccount)(nil)

	BaseAccountSig   = []byte("bacc")
	ModuleAccountSig = []byte("macc")

	PubKeyTypeSecp256k1 = byte(1)
	PubKeyTypeSecp256R1 = byte(2)
	PubKeyTypeEd25519   = byte(3)
	PubKeyTypeMultisig  = byte(4)
)

// NewBaseAccount creates a new BaseAccount object
//nolint:interfacer
func NewBaseAccount(address sdk.AccAddress, pubKey cryptotypes.PubKey, accountNumber, sequence uint64) *BaseAccount {
	acc := &BaseAccount{
		Address:       address.String(),
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}

	err := acc.SetPubKey(pubKey)
	if err != nil {
		panic(err)
	}

	return acc
}

// ProtoBaseAccount - a prototype function for BaseAccount
func ProtoBaseAccount() AccountI {
	return &BaseAccount{}
}

// NewBaseAccountWithAddress - returns a new base account with a given address
// leaving AccountNumber and Sequence to zero.
func NewBaseAccountWithAddress(addr sdk.AccAddress) *BaseAccount {
	return &BaseAccount{
		Address: addr.String(),
	}
}

// GetAddress - Implements sdk.AccountI.
func (acc BaseAccount) GetAddress() sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(acc.Address)
	return addr
}

// SetAddress - Implements sdk.AccountI.
func (acc *BaseAccount) SetAddress(addr sdk.AccAddress) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}

	acc.Address = addr.String()
	return nil
}

// GetPubKey - Implements sdk.AccountI.
func (acc BaseAccount) GetPubKey() (pk cryptotypes.PubKey) {
	if acc.Ed25519PubKey != nil {
		return acc.Ed25519PubKey
	} else if acc.Secp256K1PubKey != nil {
		return acc.Secp256K1PubKey
	} else if acc.Secp256R1PubKey != nil {
		return acc.Secp256R1PubKey
	} else if acc.MultisigPubKey != nil {
		return acc.MultisigPubKey
	}
	return nil
}

// SetPubKey - Implements sdk.AccountI.
func (acc *BaseAccount) SetPubKey(pubKey cryptotypes.PubKey) error {
	if pubKey == nil {
		acc.Ed25519PubKey, acc.Secp256K1PubKey, acc.Secp256R1PubKey, acc.MultisigPubKey = nil, nil, nil, nil
	} else if pk, ok := pubKey.(*ed25519.PubKey); ok {
		acc.Ed25519PubKey, acc.Secp256K1PubKey, acc.Secp256R1PubKey, acc.MultisigPubKey = pk, nil, nil, nil
	} else if pk, ok := pubKey.(*secp256k1.PubKey); ok {
		acc.Ed25519PubKey, acc.Secp256K1PubKey, acc.Secp256R1PubKey, acc.MultisigPubKey = nil, pk, nil, nil
	} else if pk, ok := pubKey.(*secp256r1.PubKey); ok {
		acc.Ed25519PubKey, acc.Secp256K1PubKey, acc.Secp256R1PubKey, acc.MultisigPubKey = nil, nil, pk, nil
	} else if pk, ok := pubKey.(*multisig.LegacyAminoPubKey); ok {
		acc.Ed25519PubKey, acc.Secp256K1PubKey, acc.Secp256R1PubKey, acc.MultisigPubKey = nil, nil, nil, pk
	} else {
		return fmt.Errorf("invalid pubkey")
	}
	return nil
}

// GetAccountNumber - Implements AccountI
func (acc BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements AccountI
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements sdk.AccountI.
func (acc BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements sdk.AccountI.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// Validate checks for errors on the account fields
func (acc BaseAccount) Validate() error {
	if acc.Address == "" || acc.GetPubKey() == nil {
		return nil
	}

	accAddr, err := sdk.AccAddressFromBech32(acc.Address)
	if err != nil {
		return err
	}

	if !bytes.Equal(acc.GetPubKey().Address().Bytes(), accAddr.Bytes()) {
		return errors.New("account address and pubkey address do not match")
	}

	return nil
}

func (acc BaseAccount) String() string {
	out, _ := acc.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of an account.
func (acc BaseAccount) MarshalYAML() (interface{}, error) {
	bz, err := codec.MarshalYAML(codec.NewProtoCodec(codectypes.NewInterfaceRegistry()), &acc)
	if err != nil {
		return nil, err
	}
	return string(bz), err
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (acc BaseAccount) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if acc.MultisigPubKey != nil {
		return codectypes.UnpackInterfaces(acc.MultisigPubKey, unpacker)
	}
	return nil
}

func (acc *BaseAccount) MarshalX() ([]byte, error) {
	bz, err := acc.Marshal()
	if err != nil {
		return nil, err
	}
	t := BaseAccountSig
	b := make([]byte, len(t)+len(bz))
	copy(b, t)
	copy(b[len(t):], bz)
	return b, nil
}

// NewModuleAddress creates an AccAddress from the hash of the module's name
func NewModuleAddress(name string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(name)))
}

// NewEmptyModuleAccount creates a empty ModuleAccount from a string
func NewEmptyModuleAccount(name string, permissions ...string) *ModuleAccount {
	moduleAddress := NewModuleAddress(name)
	baseAcc := NewBaseAccountWithAddress(moduleAddress)

	if err := validatePermissions(permissions...); err != nil {
		panic(err)
	}

	return &ModuleAccount{
		BaseAccount: baseAcc,
		Name:        name,
		Permissions: permissions,
	}
}

// NewModuleAccount creates a new ModuleAccount instance
func NewModuleAccount(ba *BaseAccount, name string, permissions ...string) *ModuleAccount {
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

// SetPubKey - Implements AccountI
func (ma ModuleAccount) SetPubKey(pubKey cryptotypes.PubKey) error {
	return fmt.Errorf("not supported for module accounts")
}

// Validate checks for errors on the account fields
func (ma ModuleAccount) Validate() error {
	if strings.TrimSpace(ma.Name) == "" {
		return errors.New("module account name cannot be blank")
	}

	if ma.Address != sdk.AccAddress(crypto.AddressHash([]byte(ma.Name))).String() {
		return fmt.Errorf("address %s cannot be derived from the module name '%s'", ma.Address, ma.Name)
	}

	return ma.BaseAccount.Validate()
}

func (ma *ModuleAccount) MarshalX() ([]byte, error) {
	bz, err := ma.Marshal()
	if err != nil {
		return nil, err
	}
	t := ModuleAccountSig
	b := make([]byte, len(t)+len(bz))
	copy(b, t)
	copy(b[len(t):], bz)
	return b, nil
}

type moduleAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
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
	accAddr, err := sdk.AccAddressFromBech32(ma.Address)
	if err != nil {
		return nil, err
	}

	bs, err := yaml.Marshal(moduleAccountPretty{
		Address:       accAddr,
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
	accAddr, err := sdk.AccAddressFromBech32(ma.Address)
	if err != nil {
		return nil, err
	}

	return json.Marshal(moduleAccountPretty{
		Address:       accAddr,
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
	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}

	ma.BaseAccount = NewBaseAccount(alias.Address, nil, alias.AccountNumber, alias.Sequence)
	ma.Name = alias.Name
	ma.Permissions = alias.Permissions

	return nil
}

// AccountI is an interface used to store coins at a given address within state.
// It presumes a notion of sequence numbers for replay protection,
// a notion of account numbers for replay protection for previously pruned accounts,
// and a pubkey for authentication purposes.
//
// Many complex conditions can be used in the concrete struct which implements AccountI.
type AccountI interface {
	proto.Message

	GetAddress() sdk.AccAddress
	SetAddress(sdk.AccAddress) error // errors if already set.

	GetPubKey() cryptotypes.PubKey // can return nil.
	SetPubKey(cryptotypes.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error

	// Ensure that account implements stringer
	String() string

	MarshalX() ([]byte, error)
}

func MarshalAccountX(cdc codec.BinaryCodec, acc AccountI) ([]byte, error) {
	if bacc, ok := acc.(*BaseAccount); ok && bacc.MultisigPubKey == nil {
		return acc.MarshalX()
	} else if macc, ok := acc.(*ModuleAccount); ok && macc.MultisigPubKey == nil {
		return acc.MarshalX()
	} else {
		return cdc.MarshalInterface(acc)
	}
}

func UnmarshalAccountX(cdc codec.BinaryCodec, bz []byte) (AccountI, error) {
	sigLen := len(BaseAccountSig)
	if len(bz) < sigLen {
		return nil, fmt.Errorf("invalid data")
	}
	if bytes.Equal(bz[:sigLen], BaseAccountSig) {
		acc := &BaseAccount{}
		if err := acc.Unmarshal(bz[sigLen:]); err != nil {
			return nil, err
		}
		return acc, nil
	} else if bytes.Equal(bz[:sigLen], ModuleAccountSig) {
		acc := &ModuleAccount{}
		if err := acc.Unmarshal(bz[sigLen:]); err != nil {
			return nil, err
		}
		return acc, nil
	} else {
		var acc AccountI
		return acc, cdc.UnmarshalInterface(bz, &acc)
	}
}

// ModuleAccountI defines an account interface for modules that hold tokens in
// an escrow.
type ModuleAccountI interface {
	AccountI

	GetName() string
	GetPermissions() []string
	HasPermission(string) bool
}

// GenesisAccounts defines a slice of GenesisAccount objects
type GenesisAccounts []GenesisAccount

// Contains returns true if the given address exists in a slice of GenesisAccount
// objects.
func (ga GenesisAccounts) Contains(addr sdk.Address) bool {
	for _, acc := range ga {
		if acc.GetAddress().Equals(addr) {
			return true
		}
	}

	return false
}

// GenesisAccount defines a genesis account that embeds an AccountI with validation capabilities.
type GenesisAccount interface {
	AccountI

	Validate() error
}

// custom json marshaler for BaseAccount & ModuleAccount

type PubKeyJSON struct {
	Type byte   `json:"type"`
	Key  []byte `json:"key"`
}

type BaseAccountJSON struct {
	Address       string     `json:"address"`
	PubKey        PubKeyJSON `json:"pub_key"`
	AccountNumber uint64     `json:"account_number,string"`
	Sequence      string     `json:"sequence"`
}

func (acc BaseAccount) MarshalJSONPB(m *jsonpb.Marshaler) ([]byte, error) {
	var bi BaseAccountJSON

	bi.Address = acc.GetAddress().String()
	bi.AccountNumber = acc.GetAccountNumber()
	bi.Sequence = strconv.FormatUint(acc.Sequence, 10)
	var bz []byte
	var err error
	if acc.Secp256K1PubKey != nil {
		bi.PubKey.Type = PubKeyTypeSecp256k1
		bz, err = acc.Secp256K1PubKey.Marshal()
	} else if acc.Secp256R1PubKey != nil {
		bi.PubKey.Type = PubKeyTypeSecp256R1
		bz, err = acc.Secp256R1PubKey.Marshal()
	} else if acc.Ed25519PubKey != nil {
		bi.PubKey.Type = PubKeyTypeEd25519
		bz, err = acc.Ed25519PubKey.Marshal()
	} else if acc.MultisigPubKey != nil {
		bi.PubKey.Type = PubKeyTypeMultisig
		bz, err = acc.MultisigPubKey.Marshal()
	}
	if err != nil {
		return nil, err
	}
	bi.PubKey.Key = bz
	return json.Marshal(bi)
}

func (acc *BaseAccount) UnmarshalJSONPB(m *jsonpb.Unmarshaler, bz []byte) error {
	var bi BaseAccountJSON

	err := json.Unmarshal(bz, &bi)
	if err != nil {
		return err
	}
	/* TODO: do we need to validate address format here
	err = sdk.ValidateAccAddress(bi.Address)
	if err != nil {
		return err
	}
	*/

	acc.Address = bi.Address
	acc.AccountNumber = bi.AccountNumber
	acc.Sequence, err = strconv.ParseUint(bi.Sequence, 10, 64)
	if err != nil {
		return err
	}

	switch bi.PubKey.Type {
	case PubKeyTypeEd25519:
		pk := new(ed25519.PubKey)
		if err := pk.Unmarshal(bi.PubKey.Key); err != nil {
			return err
		}
		acc.SetPubKey(pk)
	case PubKeyTypeSecp256k1:
		pk := new(secp256k1.PubKey)
		if err := pk.Unmarshal(bi.PubKey.Key); err != nil {
			return err
		}
		acc.SetPubKey(pk)
	case PubKeyTypeSecp256R1:
		pk := new(secp256r1.PubKey)
		if err := pk.Unmarshal(bi.PubKey.Key); err != nil {
			return err
		}
		acc.SetPubKey(pk)
	case PubKeyTypeMultisig:
		pk := new(multisig.LegacyAminoPubKey)
		if err := pk.Unmarshal(bi.PubKey.Key); err != nil {
			return err
		}
		acc.SetPubKey(pk)
	}
	return nil
}

type ModuleAccountJSON struct {
	BaseAccount json.RawMessage `json:"base_account"`
	Name        string          `json:"name"`
	Permissions []string        `json:"permissions"`
}

func (ma ModuleAccount) MarshalJSONPB(m *jsonpb.Marshaler) ([]byte, error) {
	var mi ModuleAccountJSON

	bz, err := ma.BaseAccount.MarshalJSONPB(m)
	if err != nil {
		return nil, err
	}
	mi.BaseAccount = bz
	mi.Name = ma.Name
	mi.Permissions = ma.Permissions

	return json.Marshal(mi)
}

func (ma *ModuleAccount) UnmarshalJSONPB(m *jsonpb.Unmarshaler, bz []byte) error {
	var mi ModuleAccountJSON

	err := json.Unmarshal(bz, &mi)
	if err != nil {
		return err
	}

	ma.Name = mi.Name
	ma.Permissions = mi.Permissions

	ba := new(BaseAccount)
	if err := m.Unmarshal(strings.NewReader(string(mi.BaseAccount)), ba); err != nil {
		return err
	}
	ma.BaseAccount = ba

	return nil
}
