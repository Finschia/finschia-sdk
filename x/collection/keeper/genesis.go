package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection"
)

type ProgressReporter struct {
	logger log.Logger

	workName  string
	workSize  int
	workIndex int

	prevPercentage int
}

func newProgressReporter(logger log.Logger, workName string, workSize int) ProgressReporter {
	reporter := ProgressReporter{
		logger:   logger,
		workName: workName,
		workSize: workSize,
	}
	reporter.report()

	return reporter
}

func (p ProgressReporter) report() {
	if p.workSize == 0 {
		p.logger.Info(fmt.Sprintf("Empty %s", p.workName))
		return
	}

	switch p.prevPercentage {
	case 0:
		p.logger.Info(fmt.Sprintf("Starting %s ...", p.workName))
	case 100:
		p.logger.Info(fmt.Sprintf("Done %s", p.workName))
	default:
		p.logger.Info(fmt.Sprintf("Progress: %d%%", p.prevPercentage))
	}
}

func (p *ProgressReporter) Tick() {
	if p.workIndex > p.workSize-1 {
		return
	}
	p.workIndex++

	if percentage := 100 * p.workIndex / p.workSize; percentage != p.prevPercentage {
		p.prevPercentage = percentage
		p.report()
	}
}

// InitGenesis new collection genesis
func (k Keeper) InitGenesis(ctx sdk.Context, data *collection.GenesisState) {
	k.SetParams(ctx, data.Params)

	reporter := newProgressReporter(k.Logger(ctx), "import contract", len(data.Contracts))
	for _, contract := range data.Contracts {
		k.setContract(ctx, contract)
		reporter.Tick()
	}

	reporter = newProgressReporter(k.Logger(ctx), "import next class ids", len(data.NextClassIds))
	for _, nextClassIDs := range data.NextClassIds {
		k.setNextClassIDs(ctx, nextClassIDs)
		reporter.Tick()
	}

	reporter = newProgressReporter(k.Logger(ctx), "import classes", len(data.Classes))
	for _, contractClasses := range data.Classes {
		contractID := contractClasses.ContractId

		for i := range contractClasses.Classes {
			any := &contractClasses.Classes[i]
			class := collection.TokenClassFromAny(any)
			k.setTokenClass(ctx, contractID, class)

			// legacy
			if nftClass, ok := class.(*collection.NFTClass); ok {
				k.setLegacyTokenType(ctx, contractID, nftClass.Id)
			}
		}

		reporter.Tick()
	}

	reporter = newProgressReporter(k.Logger(ctx), "import next token ids", len(data.NextTokenIds))
	for _, contractNextTokenIDs := range data.NextTokenIds {
		contractID := contractNextTokenIDs.ContractId

		for _, nextTokenID := range contractNextTokenIDs.TokenIds {
			k.setNextTokenID(ctx, contractID, nextTokenID.ClassId, nextTokenID.Id)
		}

		reporter.Tick()
	}

	reporter = newProgressReporter(k.Logger(ctx), "import balances", len(data.Balances))
	for _, contractBalances := range data.Balances {
		contractID := contractBalances.ContractId

		for _, balance := range contractBalances.Balances {
			for _, coin := range balance.Amount {
				addr, err := sdk.AccAddressFromBech32(balance.Address)
				if err != nil {
					panic(err)
				}

				k.setBalance(ctx, contractID, addr, coin.TokenId, coin.Amount)

				if err := collection.ValidateNFTID(coin.TokenId); err == nil {
					k.setOwner(ctx, contractID, coin.TokenId, addr)
				}
			}
		}

		reporter.Tick()
	}

	reporter = newProgressReporter(k.Logger(ctx), "import nfts", len(data.Nfts))
	for _, contractNFTs := range data.Nfts {
		contractID := contractNFTs.ContractId

		for _, nft := range contractNFTs.Nfts {
			k.setNFT(ctx, contractID, nft)
		}

		reporter.Tick()
	}

	reporter = newProgressReporter(k.Logger(ctx), "import parents", len(data.Parents))
	for _, contractParents := range data.Parents {
		contractID := contractParents.ContractId

		for _, relation := range contractParents.Relations {
			tokenID := relation.Self
			parentID := relation.Other
			k.setParent(ctx, contractID, tokenID, parentID)
			k.setChild(ctx, contractID, parentID, tokenID)
		}

		reporter.Tick()
	}

	reporter = newProgressReporter(k.Logger(ctx), "import authorizations", len(data.Authorizations))
	for _, contractAuthorizations := range data.Authorizations {
		for _, authorization := range contractAuthorizations.Authorizations {
			holderAddr, err := sdk.AccAddressFromBech32(authorization.Holder)
			if err != nil {
				panic(err)
			}
			operatorAddr, err := sdk.AccAddressFromBech32(authorization.Operator)
			if err != nil {
				panic(err)
			}
			k.setAuthorization(ctx, contractAuthorizations.ContractId, holderAddr, operatorAddr)
		}

		reporter.Tick()
	}

	reporter = newProgressReporter(k.Logger(ctx), "import grants", len(data.Grants))
	for _, contractGrants := range data.Grants {
		for _, grant := range contractGrants.Grants {
			granteeAddr, err := sdk.AccAddressFromBech32(grant.Grantee)
			if err != nil {
				panic(err)
			}
			k.setGrant(ctx, contractGrants.ContractId, granteeAddr, grant.Permission)
		}

		reporter.Tick()
	}

	reporter = newProgressReporter(k.Logger(ctx), "import statistics (burnt)", len(data.Burnts))
	for _, contractBurnts := range data.Burnts {
		contractID := contractBurnts.ContractId
		for _, burnt := range contractBurnts.Statistics {
			k.setBurnt(ctx, contractID, burnt.ClassId, burnt.Amount)
		}

		reporter.Tick()
	}

	reporter = newProgressReporter(k.Logger(ctx), "import statistics (supply)", len(data.Supplies))
	for _, contractSupplies := range data.Supplies {
		contractID := contractSupplies.ContractId
		for _, supply := range contractSupplies.Statistics {
			k.setSupply(ctx, contractID, supply.ClassId, supply.Amount)

			// calculate the amount of minted tokens
			burnt := k.GetBurnt(ctx, contractID, supply.ClassId)
			minted := supply.Amount.Add(burnt)
			k.setMinted(ctx, contractID, supply.ClassId, minted)
		}

		reporter.Tick()
	}
}

// ExportGenesis returns a GenesisState for a given context.
func (k Keeper) ExportGenesis(ctx sdk.Context) *collection.GenesisState {
	contracts := k.getContracts(ctx)

	return &collection.GenesisState{
		Params:         k.GetParams(ctx),
		Contracts:      contracts,
		NextClassIds:   k.getAllNextClassIDs(ctx),
		Classes:        k.getClasses(ctx, contracts),
		NextTokenIds:   k.getNextTokenIDs(ctx, contracts),
		Balances:       k.getBalances(ctx, contracts),
		Nfts:           k.getNFTs(ctx, contracts),
		Parents:        k.getParents(ctx, contracts),
		Grants:         k.getGrants(ctx, contracts),
		Authorizations: k.getAuthorizations(ctx, contracts),
		Supplies:       k.getSupplies(ctx, contracts),
		Burnts:         k.getBurnts(ctx, contracts),
	}
}

func (k Keeper) getContracts(ctx sdk.Context) []collection.Contract {
	var contracts []collection.Contract
	k.iterateContracts(ctx, func(contract collection.Contract) (stop bool) {
		contracts = append(contracts, contract)
		return false
	})

	return contracts
}

func (k Keeper) getClasses(ctx sdk.Context, contracts []collection.Contract) []collection.ContractClasses {
	var classes []collection.ContractClasses
	for _, contract := range contracts {
		contractID := contract.Id
		contractClasses := collection.ContractClasses{
			ContractId: contractID,
		}

		k.iterateContractClasses(ctx, contractID, func(class collection.TokenClass) (stop bool) {
			any := collection.TokenClassToAny(class)
			contractClasses.Classes = append(contractClasses.Classes, *any)
			return false
		})
		if len(contractClasses.Classes) != 0 {
			classes = append(classes, contractClasses)
		}
	}

	return classes
}

func (k Keeper) getAllNextClassIDs(ctx sdk.Context) []collection.NextClassIDs {
	var nextIDs []collection.NextClassIDs
	k.iterateNextTokenClassIDs(ctx, func(ids collection.NextClassIDs) (stop bool) {
		nextIDs = append(nextIDs, ids)
		return false
	})

	return nextIDs
}

func (k Keeper) getNextTokenIDs(ctx sdk.Context, contracts []collection.Contract) []collection.ContractNextTokenIDs {
	var nextIDs []collection.ContractNextTokenIDs
	for _, contract := range contracts {
		contractID := contract.Id
		contractNextIDs := collection.ContractNextTokenIDs{
			ContractId: contractID,
		}

		k.iterateContractNextTokenIDs(ctx, contractID, func(nextID collection.NextTokenID) (stop bool) {
			contractNextIDs.TokenIds = append(contractNextIDs.TokenIds, nextID)
			return false
		})
		if len(contractNextIDs.TokenIds) != 0 {
			nextIDs = append(nextIDs, contractNextIDs)
		}
	}

	return nextIDs
}

func (k Keeper) getBalances(ctx sdk.Context, contracts []collection.Contract) []collection.ContractBalances {
	var balances []collection.ContractBalances
	for _, contract := range contracts {
		contractID := contract.Id
		contractBalances := collection.ContractBalances{
			ContractId: contractID,
		}

		contractBalances.Balances = k.getContractBalances(ctx, contractID)
		if len(contractBalances.Balances) != 0 {
			balances = append(balances, contractBalances)
		}
	}

	return balances
}

func (k Keeper) getContractBalances(ctx sdk.Context, contractID string) []collection.Balance {
	var balances []collection.Balance
	addressToBalanceIndex := make(map[string]int)

	k.iterateContractBalances(ctx, contractID, func(address sdk.AccAddress, balance collection.Coin) (stop bool) {
		index, ok := addressToBalanceIndex[address.String()]
		if ok {
			balances[index].Amount = append(balances[index].Amount, balance)
			return false
		}

		accountBalance := collection.Balance{
			Address: address.String(),
			Amount:  collection.Coins{balance},
		}
		balances = append(balances, accountBalance)
		addressToBalanceIndex[address.String()] = len(balances) - 1
		return false
	})

	return balances
}

func (k Keeper) getNFTs(ctx sdk.Context, contracts []collection.Contract) []collection.ContractNFTs {
	var parents []collection.ContractNFTs
	for _, contract := range contracts {
		contractID := contract.Id
		contractParents := collection.ContractNFTs{
			ContractId: contractID,
		}

		k.iterateContractNFTs(ctx, contractID, func(nft collection.NFT) (stop bool) {
			contractParents.Nfts = append(contractParents.Nfts, nft)
			return false
		})
		if len(contractParents.Nfts) != 0 {
			parents = append(parents, contractParents)
		}
	}

	return parents
}

func (k Keeper) getParents(ctx sdk.Context, contracts []collection.Contract) []collection.ContractTokenRelations {
	var parents []collection.ContractTokenRelations
	for _, contract := range contracts {
		contractID := contract.Id
		contractParents := collection.ContractTokenRelations{
			ContractId: contractID,
		}

		k.iterateContractParents(ctx, contractID, func(tokenID, parentID string) (stop bool) {
			relation := collection.TokenRelation{
				Self:  tokenID,
				Other: parentID,
			}
			contractParents.Relations = append(contractParents.Relations, relation)
			return false
		})
		if len(contractParents.Relations) != 0 {
			parents = append(parents, contractParents)
		}
	}

	return parents
}

func (k Keeper) getAuthorizations(ctx sdk.Context, contracts []collection.Contract) []collection.ContractAuthorizations {
	var authorizations []collection.ContractAuthorizations
	for _, contract := range contracts {
		contractID := contract.Id
		contractAuthorizations := collection.ContractAuthorizations{
			ContractId: contractID,
		}

		k.iterateContractAuthorizations(ctx, contractID, func(authorization collection.Authorization) (stop bool) {
			contractAuthorizations.Authorizations = append(contractAuthorizations.Authorizations, authorization)
			return false
		})
		if len(contractAuthorizations.Authorizations) != 0 {
			authorizations = append(authorizations, contractAuthorizations)
		}
	}

	return authorizations
}

func (k Keeper) getGrants(ctx sdk.Context, contracts []collection.Contract) []collection.ContractGrants {
	var grants []collection.ContractGrants
	for _, contract := range contracts {
		contractID := contract.Id
		contractGrants := collection.ContractGrants{
			ContractId: contractID,
		}

		k.iterateContractGrants(ctx, contractID, func(grant collection.Grant) (stop bool) {
			contractGrants.Grants = append(contractGrants.Grants, grant)
			return false
		})
		if len(contractGrants.Grants) != 0 {
			grants = append(grants, contractGrants)
		}
	}

	return grants
}

func (k Keeper) getSupplies(ctx sdk.Context, contracts []collection.Contract) []collection.ContractStatistics {
	return k.getStatistics(ctx, contracts, k.iterateContractSupplies)
}

func (k Keeper) getBurnts(ctx sdk.Context, contracts []collection.Contract) []collection.ContractStatistics {
	return k.getStatistics(ctx, contracts, k.iterateContractBurnts)
}

func (k Keeper) getStatistics(ctx sdk.Context, contracts []collection.Contract, iterator func(ctx sdk.Context, contractID string, cb func(classID string, amount sdk.Int) (stop bool))) []collection.ContractStatistics {
	var statistics []collection.ContractStatistics
	for _, contract := range contracts {
		contractID := contract.Id
		contractStatistics := collection.ContractStatistics{
			ContractId: contractID,
		}

		iterator(ctx, contractID, func(classID string, amount sdk.Int) (stop bool) {
			supply := collection.ClassStatistics{
				ClassId: classID,
				Amount:  amount,
			}
			contractStatistics.Statistics = append(contractStatistics.Statistics, supply)
			return false
		})
		if len(contractStatistics.Statistics) != 0 {
			statistics = append(statistics, contractStatistics)
		}
	}

	return statistics
}
