package collection

import (
	"fmt"
	"strings"

	proto "github.com/gogo/protobuf/proto"
	codectypes "github.com/line/lbm-sdk/codec/types"
	sdk "github.com/line/lbm-sdk/types"
)

func DefaultNextClassIDs(contractID string) NextClassIDs {
	return NextClassIDs{
		ContractId:  contractID,
		Fungible:    sdk.NewUint(0),
		NonFungible: sdk.NewUint(1 << 28), // "10000000"
	}
}

type TokenClass interface {
	proto.Message

	GetId() string
	SetId(ids *NextClassIDs)

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

//-----------------------------------------------------------------------------
// FTClass
var _ TokenClass = (*FTClass)(nil)

//nolint:golint
func (c *FTClass) SetId(ids *NextClassIDs) {
	id := ids.Fungible
	ids.Fungible = id.Incr()
	c.Id = fmt.Sprintf("%08x", id.Uint64())
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

//-----------------------------------------------------------------------------
// NFTClass
var _ TokenClass = (*NFTClass)(nil)

//nolint:golint
func (c *NFTClass) SetId(ids *NextClassIDs) {
	id := ids.NonFungible
	ids.NonFungible = id.Incr()
	c.Id = fmt.Sprintf("%08x", id.Uint64())
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

//-----------------------------------------------------------------------------
// Coin
func NewFTCoin(classID string, amount sdk.Int) Coin {
	return NewCoin(NewFTID(classID), amount)
}

func NewNFTCoin(classID string, number int) Coin {
	return NewCoin(NewNFTID(classID, number), sdk.OneInt())
}

func NewCoin(id string, amount sdk.Int) Coin {
	coin := Coin{
		TokenId: id,
		Amount:  amount,
	}

	if err := coin.ValidateBasic(); err != nil {
		panic(err)
	}

	return coin
}

func (c Coin) ValidateBasic() error {
	if err := ValidateTokenID(c.TokenId); err != nil {
		return err
	}

	if c.isNil() || !c.isPositive() {
		return fmt.Errorf("invalid amount: %v", c.Amount)
	}

	if err := ValidateNFTID(c.TokenId); err == nil {
		if !c.Amount.Equal(sdk.OneInt()) {
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

//-----------------------------------------------------------------------------
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

// legacy
type Token interface {
	proto.Message
}
