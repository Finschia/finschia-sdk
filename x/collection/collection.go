package collection

import (
	"fmt"
	"strings"

	sdk "github.com/line/lbm-sdk/types"
)

func ValidateTokenID(id string) error {
	return nil
}

func ValidateNFTID(id string) error {
	return nil
}

//-----------------------------------------------------------------------------
// Coin
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
