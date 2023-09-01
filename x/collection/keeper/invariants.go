package keeper

import (
	"strings"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection"
)

const (
	totalFTSupplyInvariant = "total-ft-supply"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	for name, invariant := range map[string]func(k Keeper) sdk.Invariant{
		totalFTSupplyInvariant: TotalFTSupplyInvariant,
	} {
		ir.RegisterRoute(collection.ModuleName, name, invariant(k))
	}
}

func TotalFTSupplyInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// cache, we don't want to write changes
		ctx, _ = ctx.CacheContext()

		invalidFTClassIDs := map[string][]string{}
		k.iterateContracts(ctx, func(contract collection.Contract) (stop bool) {
			supplies := map[string]sdk.Int{}
			k.iterateContractSupplies(ctx, contract.Id, func(classID string, amount sdk.Int) (stop bool) {
				if err := collection.ValidateLegacyFTClassID(classID); err != nil {
					return false
				}

				supplies[classID] = amount

				return false
			})

			k.iterateContractBalances(ctx, contract.Id, func(address sdk.AccAddress, balance collection.Coin) (stop bool) {
				classID := collection.SplitTokenID(balance.TokenId)
				if err := collection.ValidateLegacyFTClassID(classID); err != nil {
					return false
				}

				amount, ok := supplies[classID]
				if !ok {
					amount = sdk.ZeroInt()
				}

				supplies[classID] = amount.Sub(balance.Amount)
				return false
			})

			invalidFTClassIDsCandidate := []string{}
			for classID, supply := range supplies {
				if !supply.IsZero() {
					invalidFTClassIDsCandidate = append(invalidFTClassIDsCandidate, classID)
				}
			}

			if len(invalidFTClassIDsCandidate) != 0 {
				invalidFTClassIDs[contract.Id] = invalidFTClassIDsCandidate
			}

			return false
		})

		broken := len(invalidFTClassIDs) != 0
		msg := "no violation found"
		if broken {
			concatenated := []string{}
			delimiter := ":"
			for contractID, classIDs := range invalidFTClassIDs {
				for _, classID := range classIDs {
					concatenated = append(concatenated, contractID+delimiter+classID)
				}
			}

			msg = "violation found on following ft classIDs: " + strings.Join(concatenated, ", ")
		}

		return sdk.FormatInvariant(collection.ModuleName, totalFTSupplyInvariant, msg), broken
	}
}
