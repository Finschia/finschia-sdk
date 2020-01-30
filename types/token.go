package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
)

type CoinWithTokenID struct {
	Symbol  string  `json:"symbol"`
	TokenID string  `json:"token_id"`
	Amount  sdk.Int `json:"amount"`
}

func NewCoinWithTokenID(symbol, tokenID string, amount sdk.Int) CoinWithTokenID {
	if err := ValidateSymbolUserDefined(symbol); err != nil {
		panic(err)
	}

	if err := ValidateSymbolTokenID(tokenID); err != nil {
		panic(err)
	}

	return CoinWithTokenID{
		Symbol:  symbol,
		TokenID: tokenID,
		Amount:  amount,
	}
}

func (coin CoinWithTokenID) ToCoin() sdk.Coin {
	return sdk.NewCoin(coin.GetDenom(), coin.Amount)
}

func (coin CoinWithTokenID) GetDenom() string {
	return fmt.Sprintf("%v%v", coin.Symbol, coin.TokenID)
}

func (coin CoinWithTokenID) String() string {
	return fmt.Sprintf("%v%v%v", coin.Amount, coin.Symbol, coin.TokenID)
}

func (coin CoinWithTokenID) IsValid() bool {
	if err := ValidateSymbolUserDefined(coin.Symbol); err != nil {
		return false
	}

	if err := ValidateSymbolTokenID(coin.TokenID); err != nil {
		return false
	}
	return true
}

func (coin CoinWithTokenID) IsPositive() bool {
	return coin.Amount.Sign() == 1
}

func (coin CoinWithTokenID) IsNegative() bool {
	return coin.Amount.Sign() == -1
}

type CoinWithTokenIDs []CoinWithTokenID

func NewCoinWithTokenIDs(coins ...CoinWithTokenID) CoinWithTokenIDs {
	var newCoins = CoinWithTokenIDs(coins)
	return newCoins.Sort()
}

func (coins CoinWithTokenIDs) ToCoins() sdk.Coins {
	var sdkCoins sdk.Coins

	for _, coin := range coins {
		sdkCoins = append(sdkCoins, coin.ToCoin())
	}

	return sdkCoins
}

//nolint
func (coins CoinWithTokenIDs) Len() int           { return len(coins) }
func (coins CoinWithTokenIDs) Less(i, j int) bool { return coins[i].GetDenom() < coins[j].GetDenom() }
func (coins CoinWithTokenIDs) Swap(i, j int)      { coins[i], coins[j] = coins[j], coins[i] }

var _ sort.Interface = CoinWithTokenIDs{}

func (coins CoinWithTokenIDs) Sort() CoinWithTokenIDs {
	sort.Sort(coins)
	return coins
}

func (coins CoinWithTokenIDs) IsValid() bool {
	switch len(coins) {
	case 0:
		return true
	case 1:
		return coins[0].IsValid() && coins[0].IsPositive()
	default:
		lowDenom := coins[0].GetDenom()
		for _, coin := range coins[1:] {
			if !coin.IsValid() {
				return false
			}
			if !coin.IsPositive() {
				return false
			}
			if coin.GetDenom() <= lowDenom {
				return false
			}

			lowDenom = coin.GetDenom()
		}

		return true
	}
}
