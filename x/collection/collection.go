package collection

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gogo/protobuf/proto"

	"cosmossdk.io/math"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

const (
	prefixLegacyPermission = "LEGACY_PERMISSION_"
)

// Deprecated: use Permission.
func LegacyPermissionFromString(name string) LegacyPermission {
	legacyPermissionName := prefixLegacyPermission + strings.ToUpper(name)
	return LegacyPermission(LegacyPermission_value[legacyPermissionName])
}

func (x LegacyPermission) String() string {
	lenPrefix := len(prefixLegacyPermission)
	return strings.ToLower(LegacyPermission_name[int32(x)][lenPrefix:])
}

func DefaultNextClassIDs(contractID string) NextClassIDs {
	return NextClassIDs{
		ContractId:  contractID,
		Fungible:    math.NewUint(1),
		NonFungible: math.NewUint(1 << 28).Incr(), // "10000000 + 1"
	}
}

func validateParams(params Params) error {
	// limits are uint32, so no need to validate them.
	return nil
}

type TokenClass interface {
	proto.Message

	GetId() string
	SetId(ids *NextClassIDs)

	SetName(name string)

	SetMeta(meta string)

	ValidateBasic() error
}

func TokenClassToAny(class TokenClass) *codectypes.Any {
	msg := class.(proto.Message)

	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}

	return any
}

func TokenClassFromAny(any *codectypes.Any) TokenClass {
	class := any.GetCachedValue().(TokenClass)
	return class
}

func TokenClassUnpackInterfaces(any *codectypes.Any, unpacker codectypes.AnyUnpacker) error {
	var class TokenClass
	return unpacker.UnpackAny(any, &class)
}

// ----------------------------------------------------------------------------
// FTClass
var _ TokenClass = (*FTClass)(nil)

//lint:ignore var-naming
func (c *FTClass) SetId(ids *NextClassIDs) {
	id := ids.Fungible
	ids.Fungible = id.Incr()
	c.Id = fmt.Sprintf("%08x", id.Uint64())
}

func (c *FTClass) SetName(name string) {
	c.Name = name
}

func (c *FTClass) SetMeta(meta string) {
	c.Meta = meta
}

func (c FTClass) ValidateBasic() error {
	if err := ValidateClassID(c.Id); err != nil {
		return err
	}

	if err := validateName(c.Name); err != nil {
		return err
	}
	if err := validateMeta(c.Meta); err != nil {
		return err
	}
	if err := validateDecimals(c.Decimals); err != nil {
		return err
	}

	return nil
}

// ----------------------------------------------------------------------------
// NFTClass
var _ TokenClass = (*NFTClass)(nil)

//lint:ignore var-naming
func (c *NFTClass) SetId(ids *NextClassIDs) {
	id := ids.NonFungible
	ids.NonFungible = id.Incr()
	c.Id = fmt.Sprintf("%08x", id.Uint64())
}

func (c *NFTClass) SetName(name string) {
	c.Name = name
}

func (c *NFTClass) SetMeta(meta string) {
	c.Meta = meta
}

func (c NFTClass) ValidateBasic() error {
	if err := ValidateClassID(c.Id); err != nil {
		return err
	}

	if err := validateName(c.Name); err != nil {
		return err
	}
	if err := validateMeta(c.Meta); err != nil {
		return err
	}

	return nil
}

// ----------------------------------------------------------------------------
// Coin
func NewFTCoin(classID string, amount math.Int) Coin {
	return NewCoin(NewFTID(classID), amount)
}

func NewNFTCoin(classID string, number int) Coin {
	return NewCoin(NewNFTID(classID, number), math.OneInt())
}

func NewCoin(id string, amount math.Int) Coin {
	coin := Coin{
		TokenId: id,
		Amount:  amount,
	}

	if err := coin.ValidateBasic(); err != nil {
		panic(err)
	}

	return coin
}

func (c Coin) String() string {
	return fmt.Sprintf("%s:%s", c.TokenId, c.Amount)
}

func (c Coin) ValidateBasic() error {
	if err := ValidateTokenID(c.TokenId); err != nil {
		return err
	}

	if c.isNil() || !c.isPositive() {
		return fmt.Errorf("invalid amount: %v", c.Amount)
	}

	if err := ValidateNFTID(c.TokenId); err == nil {
		if !c.Amount.Equal(math.OneInt()) {
			return fmt.Errorf("duplicate non fungible tokens")
		}
	}

	return nil
}

func (c Coin) isPositive() bool {
	return c.Amount.IsPositive()
}

func (c Coin) isNil() bool {
	return c.Amount.IsNil()
}

var reDecCoin = regexp.MustCompile(fmt.Sprintf(`^(%s%s):([[:digit:]]+)$`, patternClassID, patternAll))

func ParseCoin(coinStr string) (*Coin, error) {
	coinStr = strings.TrimSpace(coinStr)

	matches := reDecCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		return nil, fmt.Errorf("invalid coin expression: %s", coinStr)
	}

	id, amountStr := matches[1], matches[2]

	amount, ok := math.NewIntFromString(amountStr)
	if !ok {
		return nil, fmt.Errorf("failed to parse coin amount: %s", amountStr)
	}

	coin := NewCoin(id, amount)
	return &coin, nil
}

// ----------------------------------------------------------------------------
// Coins
type Coins []Coin

func NewCoins(coins ...Coin) Coins {
	newCoins := Coins(coins)
	if err := newCoins.ValidateBasic(); err != nil {
		panic(fmt.Errorf("invalid coin %s: %w", newCoins, err))
	}

	return newCoins
}

func (coins Coins) String() string {
	if len(coins) == 0 {
		return ""
	} else if len(coins) == 1 {
		return coins[0].String()
	}

	var out strings.Builder
	for _, coin := range coins[:len(coins)-1] {
		out.WriteString(coin.String())
		out.WriteByte(',')
	}
	out.WriteString(coins[len(coins)-1].String())
	return out.String()
}

func (coins Coins) ValidateBasic() error {
	if len(coins) == 0 {
		return fmt.Errorf("empty coins")
	}

	seenIDs := map[string]bool{}
	for _, coin := range coins {
		if seenIDs[coin.TokenId] {
			return fmt.Errorf("duplicate id %s", coin.TokenId)
		}
		seenIDs[coin.TokenId] = true

		if err := coin.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid coin %s: %w", coin.TokenId, err)
		}
	}

	return nil
}

func ParseCoins(coinsStr string) (Coins, error) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return nil, fmt.Errorf("invalid string for coins")
	}

	coinStrs := strings.Split(coinsStr, ",")
	coins := make(Coins, len(coinStrs))
	for i, coinStr := range coinStrs {
		coin, err := ParseCoin(coinStr)
		if err != nil {
			return nil, err
		}

		coins[i] = *coin
	}

	return NewCoins(coins...), nil
}

type Token interface {
	proto.Message
}

func TokenFromAny(any *codectypes.Any) Token {
	class := any.GetCachedValue().(Token)
	return class
}

func TokenUnpackInterfaces(any *codectypes.Any, unpacker codectypes.AnyUnpacker) error {
	var token Token
	return unpacker.UnpackAny(any, &token)
}
