package types

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/dgraph-io/ristretto"
	yaml "gopkg.in/yaml.v2"

	"github.com/line/lbm-sdk/codec/legacy"
	cryptotypes "github.com/line/lbm-sdk/crypto/types"
	"github.com/line/lbm-sdk/types/bech32"
)

const (
	// Constants defined here are the defaults value for address.
	// You can use the specific values for your project.
	// Add the follow lines to the `main()` of your server.
	//
	//	config := sdk.GetConfig()
	//	config.SetBech32PrefixForAccount(yourBech32PrefixAccAddr, yourBech32PrefixAccPub)
	//	config.SetBech32PrefixForValidator(yourBech32PrefixValAddr, yourBech32PrefixValPub)
	//	config.SetBech32PrefixForConsensusNode(yourBech32PrefixConsAddr, yourBech32PrefixConsPub)
	//	config.SetCoinType(yourCoinType)
	//	config.SetFullFundraiserPath(yourFullFundraiserPath)
	//	config.Seal()

	BytesAddrLen = 20
	//AddrLen = len(Bech32MainPrefix) + 1 + 38

	// Bech32MainPrefix defines the main SDK Bech32 prefix of an account's address
	Bech32MainPrefix = "link"

	// CoinType is the LINK coin type as defined in SLIP44 (https://github.com/satoshilabs/slips/blob/master/slip-0044.md)
	CoinType = 438

	// FullFundraiserPath is the parts of the BIP44 HD path that are fixed by
	// what we used during the LINK fundraiser.
	FullFundraiserPath = "m/44'/438'/0'/0/0"

	// PrefixAccount is the prefix for account keys
	PrefixAccount = "acc"
	// PrefixValidator is the prefix for validator keys
	PrefixValidator = "val"
	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus = "cons"
	// PrefixPublic is the prefix for public keys
	PrefixPublic = "pub"
	// PrefixOperator is the prefix for operator keys
	PrefixOperator = "oper"

	// PrefixAddress is the prefix for addresses
	PrefixAddress = "addr"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32MainPrefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32MainPrefix + PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32MainPrefix + PrefixValidator + PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32MainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32MainPrefix + PrefixValidator + PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32MainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic
)

const DefaultBech32CacheSize = 1 << 30 // maximum size of cache (1GB)

// Address is a common interface for different types of addresses used by the SDK
type Address interface {
	Equals(Address) bool
	Empty() bool
	Marshal() ([]byte, error)
	MarshalJSON() ([]byte, error)
	Bytes() []byte
	String() string
	Format(s fmt.State, verb rune)
}

// Ensure that different address types implement the interface
var _ Address = AccAddress("")
var _ Address = ValAddress("")
var _ Address = ConsAddress("")

var _ yaml.Marshaler = AccAddress("")
var _ yaml.Marshaler = ValAddress("")
var _ yaml.Marshaler = ConsAddress("")

// ----------------------------------------------------------------------------
// account
// ----------------------------------------------------------------------------

// TODO We should add a layer to choose whether to access the cache or to run actual conversion
// bech32 encoding and decoding takes a lot of time, so memoize it
var bech32Cache Bech32Cache

type Bech32Cache struct {
	bech32ToAddrCache *ristretto.Cache
	addrToBech32Cache *ristretto.Cache
}

func SetBech32Cache(size int64) {
	var err error
	config := &ristretto.Config{
		NumCounters: 1e7, // number of keys to track frequency of (10M).
		MaxCost:     size,
		BufferItems: 64, // number of keys per Get buffer.
	}
	bech32Cache.bech32ToAddrCache, err = ristretto.NewCache(config)
	if err != nil {
		panic(err)
	}
	bech32Cache.addrToBech32Cache, err = ristretto.NewCache(config)
	if err != nil {
		panic(err)
	}
}

// Used only for test cases
func InvalidateBech32Cache() {
	bech32Cache.bech32ToAddrCache = nil
	bech32Cache.addrToBech32Cache = nil
}

func (cache *Bech32Cache) GetAddr(bech32Addr string) ([]byte, bool) {
	if cache.bech32ToAddrCache != nil {
		rawAddr, ok := cache.bech32ToAddrCache.Get(bech32Addr)
		if ok {
			return rawAddr.([]byte), ok
		}
	}
	return nil, false
}

func (cache *Bech32Cache) GetBech32(rawAddr []byte) (string, bool) {
	if cache.addrToBech32Cache != nil {
		bech32Addr, ok := cache.addrToBech32Cache.Get(string(rawAddr))
		if ok {
			return bech32Addr.(string), ok
		}
	}
	return "", false
}

func (cache *Bech32Cache) Set(bech32Addr string, rawAddr []byte) {
	if cache.bech32ToAddrCache != nil {
		cache.bech32ToAddrCache.Set(bech32Addr, rawAddr, int64(len(rawAddr)))
	}
	if cache.addrToBech32Cache != nil {
		cache.addrToBech32Cache.Set(string(rawAddr), bech32Addr, int64(len(bech32Addr)))
	}
}

// AccAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type AccAddress string

// AccAddressFromHex creates an AccAddress from a hex string.
func AccAddressFromHex(addressBytesHex string) (addr AccAddress, err error) {
	bz, err := addressBytesFromHexString(addressBytesHex)
	return BytesToAccAddress(bz), err
}

// VerifyAddressFormat verifies that the provided bytes form a valid address
// according to the default address rules or a custom address verifier set by
// GetConfig().SetAddressVerifier()
func VerifyAddressFormat(bz []byte) error {
	verifier := GetConfig().GetAddressVerifier()
	if verifier != nil {
		return verifier(bz)
	}
	if len(bz) != BytesAddrLen {
		if len(bz) == 0 {
			return errors.New("empty address string is not allowed")
		}
		return fmt.Errorf("incorrect address length (expected: %d, actual: %d)", BytesAddrLen, len(bz))
	}
	return nil
}

// Returns boolean for whether two AccAddresses are Equal
func (aa AccAddress) Equals(aa2 Address) bool {
	if aa.Empty() && aa2.Empty() {
		return true
	}

	return strings.EqualFold(aa.String(), aa2.String())
}

// Returns boolean for whether an AccAddress is empty
func (aa AccAddress) Empty() bool {
	return len(string(aa)) == 0
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (aa AccAddress) Marshal() ([]byte, error) {
	return []byte(aa.String()), nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (aa *AccAddress) Unmarshal(data []byte) error {
	*aa = AccAddress(data)
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (aa AccAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(aa.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (aa AccAddress) MarshalYAML() (interface{}, error) {
	return aa.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (aa *AccAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	// TODO: address validation?
	*aa = AccAddress(s)
	return nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (aa *AccAddress) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	// TODO: address validation?
	*aa = AccAddress(s)
	return nil
}

// Bytes returns the raw address bytes.
func (aa AccAddress) Bytes() []byte {
	return []byte(aa.String())
}

// String implements the Stringer interface.
func (aa AccAddress) String() string {
	return string(aa)
}

func (aa AccAddress) ToValAddress() ValAddress {
	bytes, _ := AccAddressToBytes(aa.String())
	return BytesToValAddress(bytes)
}

func (aa AccAddress) ToConsAddress() ConsAddress {
	bytes, _ := AccAddressToBytes(aa.String())
	return BytesToConsAddress(bytes)
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (aa AccAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(aa.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", aa)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(aa))))
	}
}

func BytesToAccAddress(addrBytes []byte) AccAddress {
	bech32Addr, ok := bech32Cache.GetBech32(addrBytes)
	if ok {
		return AccAddress(bech32Addr)
	}
	bech32PrefixAccAddr := GetConfig().GetBech32AccountAddrPrefix()

	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixAccAddr, addrBytes)
	if err != nil {
		panic(err)
	}
	bech32Cache.Set(bech32Addr, addrBytes)
	return AccAddress(bech32Addr)
}

func AccAddressToBytes(bech32Addr string) ([]byte, error) {
	bz, ok := bech32Cache.GetAddr(bech32Addr)
	if ok {
		return bz, nil
	}

	if len(strings.TrimSpace(bech32Addr)) == 0 {
		return nil, errors.New("empty address string is not allowed")
	}

	bech32PrefixAccAddr := GetConfig().GetBech32AccountAddrPrefix()

	bz, err := GetFromBech32(bech32Addr, bech32PrefixAccAddr)
	if err != nil {
		return nil, err
	}

	err = VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}
	bech32Cache.Set(bech32Addr, bz)
	return bz, nil
}

func ValidateAccAddress(bech32Addr string) error {
	_, err := AccAddressToBytes(bech32Addr)
	return err
}

// ----------------------------------------------------------------------------
// validator operator
// ----------------------------------------------------------------------------

// ValAddress defines a wrapper around bytes meant to present a validator's
// operator. When marshaled to a string or JSON, it uses Bech32.
type ValAddress string

func BytesToValAddress(addrBytes []byte) ValAddress {
	bech32PrefixValAddr := GetConfig().GetBech32ValidatorAddrPrefix()

	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixValAddr, addrBytes)
	if err != nil {
		panic(err)
	}
	return ValAddress(bech32Addr)
}

func ValAddressToBytes(bech32Addr string) ([]byte, error) {
	if len(strings.TrimSpace(bech32Addr)) == 0 {
		return nil, errors.New("empty address string is not allowed")
	}

	bech32PrefixValAddr := GetConfig().GetBech32ValidatorAddrPrefix()

	bz, err := GetFromBech32(bech32Addr, bech32PrefixValAddr)
	if err != nil {
		return nil, err
	}

	err = VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func ValidateValAddress(bech32Addr string) error {
	_, err := ValAddressToBytes(bech32Addr)
	return err
}

// ValAddressFromHex creates a ValAddress from a hex string.
func ValAddressFromHex(address string) (addr ValAddress, err error) {
	bz, err := addressBytesFromHexString(address)
	return BytesToValAddress(bz), err
}

// Returns boolean for whether two ValAddresses are Equal
func (va ValAddress) Equals(va2 Address) bool {
	if va.Empty() && va2.Empty() {
		return true
	}

	return strings.EqualFold(va.String(), va2.String())
}

// Returns boolean for whether an AccAddress is empty
func (va ValAddress) Empty() bool {
	return va == ""
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (va ValAddress) Marshal() ([]byte, error) {
	return []byte(va), nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (va *ValAddress) Unmarshal(data []byte) error {
	*va = ValAddress(data)
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (va ValAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(va.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (va ValAddress) MarshalYAML() (interface{}, error) {
	return va.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (va *ValAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*va = ValAddress(s)
	return nil
}

// UnmarshalYAML unmarshals from YAML assuming Bech32 encoding.
func (va *ValAddress) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*va = ValAddress(s)
	return nil
}

// Bytes returns the raw address bytes.
func (va ValAddress) Bytes() []byte {
	return []byte(va.String())
}

// String implements the Stringer interface.
func (va ValAddress) String() string {
	return string(va)
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (va ValAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(va.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", va)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(va))))
	}
}

func (va ValAddress) ToAccAddress() AccAddress {
	bytes, _ := ValAddressToBytes(va.String())
	return BytesToAccAddress(bytes)
}

func (va ValAddress) ToConsAddress() ConsAddress {
	bytes, _ := ValAddressToBytes(va.String())
	return BytesToConsAddress(bytes)
}

// ----------------------------------------------------------------------------
// consensus node
// ----------------------------------------------------------------------------

// ConsAddress defines a wrapper around bytes meant to present a consensus node.
// When marshaled to a string or JSON, it uses Bech32.
type ConsAddress string

func BytesToConsAddress(addrBytes []byte) ConsAddress {
	bech32PrefixConsAddr := GetConfig().GetBech32ConsensusAddrPrefix()

	bech32Addr, err := bech32.ConvertAndEncode(bech32PrefixConsAddr, addrBytes)
	if err != nil {
		panic(err)
	}
	return ConsAddress(bech32Addr)
}

func ConsAddressToBytes(bech32Addr string) ([]byte, error) {
	if len(strings.TrimSpace(bech32Addr)) == 0 {
		return nil, errors.New("empty address string is not allowed")
	}

	bech32PrefixConsAddr := GetConfig().GetBech32ConsensusAddrPrefix()

	bz, err := GetFromBech32(bech32Addr, bech32PrefixConsAddr)
	if err != nil {
		return nil, err
	}

	err = VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

// ConsAddressFromHex creates a ConsAddress from a hex string.
func ConsAddressFromHex(address string) (addr ConsAddress, err error) {
	bz, err := addressBytesFromHexString(address)
	return BytesToConsAddress(bz), err
}

func ValidateConsAddress(bech32Addr string) error {
	_, err := ConsAddressToBytes(bech32Addr)
	return err
}

// get ConsAddress from pubkey
func GetConsAddress(pubkey cryptotypes.PubKey) ConsAddress {
	return BytesToConsAddress(pubkey.Address())
}

// Returns boolean for whether two ConsAddress are Equal
func (ca ConsAddress) Equals(ca2 Address) bool {
	if ca.Empty() && ca2.Empty() {
		return true
	}

	return strings.EqualFold(ca.String(), ca2.String())
}

// Returns boolean for whether an ConsAddress is empty
func (ca ConsAddress) Empty() bool {
	return ca == ""
}

// Marshal returns the raw address bytes. It is needed for protobuf
// compatibility.
func (ca ConsAddress) Marshal() ([]byte, error) {
	return []byte(ca), nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (ca *ConsAddress) Unmarshal(data []byte) error {
	*ca = ConsAddress(data)
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (ca ConsAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(ca.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (ca ConsAddress) MarshalYAML() (interface{}, error) {
	return ca.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (ca *ConsAddress) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*ca = ConsAddress(s)
	return nil
}

// UnmarshalYAML unmarshals from YAML assuming Bech32 encoding.
func (ca *ConsAddress) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*ca = ConsAddress(s)
	return nil
}

// Bytes returns the raw address bytes.
func (ca ConsAddress) Bytes() []byte {
	return []byte(ca)
}

// String implements the Stringer interface.
func (ca ConsAddress) String() string {
	return string(ca)
}

// Bech32ifyAddressBytes returns a bech32 representation of address bytes.
// Returns an empty sting if the byte slice is 0-length. Returns an error if the bech32 conversion
// fails or the prefix is empty.
func Bech32ifyAddressBytes(prefix string, bs []byte) (string, error) {
	if len(bs) == 0 {
		return "", nil
	}
	if len(prefix) == 0 {
		return "", errors.New("prefix cannot be empty")
	}
	return bech32.ConvertAndEncode(prefix, bs)
}

// MustBech32ifyAddressBytes returns a bech32 representation of address bytes.
// Returns an empty sting if the byte slice is 0-length. It panics if the bech32 conversion
// fails or the prefix is empty.
func MustBech32ifyAddressBytes(prefix string, bs []byte) string {
	s, err := Bech32ifyAddressBytes(prefix, bs)
	if err != nil {
		panic(err)
	}
	return s
}

// Format implements the fmt.Formatter interface.
// nolint: errcheck
func (ca ConsAddress) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(ca.String()))
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", ca)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(ca))))
	}
}

// ----------------------------------------------------------------------------
// auxiliary
// ----------------------------------------------------------------------------

// Bech32PubKeyType defines a string type alias for a Bech32 public key type.
type Bech32PubKeyType string

// Bech32 conversion constants
const (
	Bech32PubKeyTypeAccPub  Bech32PubKeyType = "accpub"
	Bech32PubKeyTypeValPub  Bech32PubKeyType = "valpub"
	Bech32PubKeyTypeConsPub Bech32PubKeyType = "conspub"
)

// Bech32ifyPubKey returns a Bech32 encoded string containing the appropriate
// prefix based on the key type provided for a given PublicKey.
// TODO: Remove Bech32ifyPubKey and all usages (cosmos/cosmos-sdk/issues/#7357)
func Bech32ifyPubKey(pkt Bech32PubKeyType, pubkey cryptotypes.PubKey) (string, error) {
	var bech32Prefix string

	switch pkt {
	case Bech32PubKeyTypeAccPub:
		bech32Prefix = GetConfig().GetBech32AccountPubPrefix()

	case Bech32PubKeyTypeValPub:
		bech32Prefix = GetConfig().GetBech32ValidatorPubPrefix()

	case Bech32PubKeyTypeConsPub:
		bech32Prefix = GetConfig().GetBech32ConsensusPubPrefix()

	}

	return bech32.ConvertAndEncode(bech32Prefix, legacy.Cdc.MustMarshalBinaryBare(pubkey))
}

// MustBech32ifyPubKey calls Bech32ifyPubKey except it panics on error.
func MustBech32ifyPubKey(pkt Bech32PubKeyType, pubkey cryptotypes.PubKey) string {
	res, err := Bech32ifyPubKey(pkt, pubkey)
	if err != nil {
		panic(err)
	}

	return res
}

// GetPubKeyFromBech32 returns a PublicKey from a bech32-encoded PublicKey with
// a given key type.
func GetPubKeyFromBech32(pkt Bech32PubKeyType, pubkeyStr string) (cryptotypes.PubKey, error) {
	var bech32Prefix string

	switch pkt {
	case Bech32PubKeyTypeAccPub:
		bech32Prefix = GetConfig().GetBech32AccountPubPrefix()

	case Bech32PubKeyTypeValPub:
		bech32Prefix = GetConfig().GetBech32ValidatorPubPrefix()

	case Bech32PubKeyTypeConsPub:
		bech32Prefix = GetConfig().GetBech32ConsensusPubPrefix()

	}

	bz, err := GetFromBech32(pubkeyStr, bech32Prefix)
	if err != nil {
		return nil, err
	}

	return legacy.PubKeyFromBytes(bz)
}

// MustGetPubKeyFromBech32 calls GetPubKeyFromBech32 except it panics on error.
func MustGetPubKeyFromBech32(pkt Bech32PubKeyType, pubkeyStr string) cryptotypes.PubKey {
	res, err := GetPubKeyFromBech32(pkt, pubkeyStr)
	if err != nil {
		panic(err)
	}

	return res
}

// GetFromBech32 decodes a bytestring from a Bech32 encoded string.
func GetFromBech32(bech32str, prefix string) ([]byte, error) {
	if len(bech32str) == 0 {
		return nil, errors.New("decoding Bech32 address failed: must provide an address")
	}

	hrp, bz, err := bech32.DecodeAndConvert(bech32str)
	if err != nil {
		return nil, err
	}

	if hrp != prefix {
		return nil, fmt.Errorf("invalid Bech32 prefix; expected %s, got %s", prefix, hrp)
	}

	return bz, nil
}

func addressBytesFromHexString(address string) ([]byte, error) {
	if len(address) == 0 {
		return nil, errors.New("decoding Bech32 address failed: must provide an address")
	}

	return hex.DecodeString(address)
}
