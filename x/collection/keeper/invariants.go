package keeper

import (
	"strings"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection"
)

const (
	totalSupplyInvariant = "total-supply"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	for name, invariant := range map[string]func(k Keeper) sdk.Invariant{
		totalSupplyInvariant: TotalSupplyInvariant,
	} {
		ir.RegisterRoute(collection.ModuleName, name, invariant(k))
	}
}

func TotalSupplyInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// cache, we don't want to write changes
		ctx, _ = ctx.CacheContext()

		invalidClassIDs := map[string][]string{}
		k.iterateContracts(ctx, func(contract collection.Contract) (stop bool) {
			supplies := map[string]sdk.Int{}
			k.iterateContractSupplies(ctx, contract.Id, func(classID string, amount sdk.Int) (stop bool) {
				supplies[classID] = amount
				return false
			})

			k.iterateContractBalances(ctx, contract.Id, func(address sdk.AccAddress, balance collection.Coin) (stop bool) {
				classID := collection.SplitTokenID(balance.TokenId)
				amount, ok := supplies[classID]
				if !ok {
					amount = sdk.ZeroInt()
				}

				supplies[classID] = amount.Sub(balance.Amount)
				return false
			})

			invalidClassIDsCandidate := []string{}
			for classID, supply := range supplies {
				if !supply.IsZero() {
					invalidClassIDsCandidate = append(invalidClassIDsCandidate, classID)
				}
			}

			if len(invalidClassIDsCandidate) != 0 {
				invalidClassIDs[contract.Id] = invalidClassIDsCandidate
			}

			return false
		})

		broken := len(invalidClassIDs) != 0
		msg := "no violation found"
		if broken {
			concatenated := []string{}
			delimiter := ":"
			for contractID, classIDs := range invalidClassIDs {
				for _, classID := range classIDs {
					concatenated = append(concatenated, contractID+delimiter+classID)
				}
			}

			msg = "violation found on following classIDs: " + strings.Join(concatenated, ", ")
		}

		return sdk.FormatInvariant(collection.ModuleName, totalSupplyInvariant, msg), broken
	}
}
