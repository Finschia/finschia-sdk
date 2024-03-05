package keeper

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/collection"
)

func (k Keeper) CreateContract(ctx sdk.Context, creator sdk.AccAddress, contract collection.Contract) string {
	contractID := k.createContract(ctx, contract)

	event := collection.EventCreatedContract{
		Creator:    k.bytesToString(creator),
		ContractId: contractID,
		Name:       contract.Name,
		Meta:       contract.Meta,
		Uri:        contract.Uri,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	// 0 is "unspecified"
	for i := 1; i < len(collection.Permission_value); i++ {
		p := collection.Permission(i)

		k.Grant(ctx, contractID, []byte{}, creator, p)
	}

	return contractID
}

func (k Keeper) createContract(ctx sdk.Context, contract collection.Contract) string {
	contractID := k.NewID(ctx)
	contract.Id = contractID
	k.setContract(ctx, contract)

	// set the next class ids
	nextIDs := collection.DefaultNextClassIDs(contractID)
	k.setNextClassIDs(ctx, nextIDs)

	return contractID
}

func (k Keeper) GetContract(ctx sdk.Context, contractID string) (*collection.Contract, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := contractKey(contractID)
	bz, _ := store.Get(key)
	if bz == nil {
		return nil, collection.ErrCollectionNotExist.Wrapf("no such a contract: %s", contractID)
	}

	var contract collection.Contract
	if err := contract.Unmarshal(bz); err != nil {
		panic(err)
	}
	return &contract, nil
}

func (k Keeper) setContract(ctx sdk.Context, contract collection.Contract) {
	store := k.storeService.OpenKVStore(ctx)
	key := contractKey(contract.Id)

	bz, err := contract.Marshal()
	if err != nil {
		panic(err)
	}
	err = store.Set(key, bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) CreateTokenClass(ctx sdk.Context, contractID string, class collection.TokenClass) (*string, error) {
	if _, err := k.GetContract(ctx, contractID); err != nil {
		panic(err)
	}

	nextClassIDs := k.getNextClassIDs(ctx, contractID)
	class.SetID(&nextClassIDs)
	k.setNextClassIDs(ctx, nextClassIDs)

	if err := class.ValidateBasic(); err != nil {
		return nil, err
	}
	k.setTokenClass(ctx, contractID, class)

	if nftClass, ok := class.(*collection.NFTClass); ok {
		k.setNextTokenID(ctx, contractID, nftClass.Id, math.OneUint())
	} else {
		panic("TokenClass only supports NFTClass")
	}

	id := class.GetId()
	return &id, nil
}

func (k Keeper) GetTokenClass(ctx sdk.Context, contractID, classID string) (collection.TokenClass, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := classKey(contractID, classID)
	bz, _ := store.Get(key)
	if bz == nil {
		return nil, sdkerrors.ErrNotFound.Wrapf("no such a class in contract %s: %s", contractID, classID)
	}

	var class collection.TokenClass
	if err := k.cdc.UnmarshalInterface(bz, &class); err != nil {
		panic(err)
	}
	return class, nil
}

func (k Keeper) setTokenClass(ctx sdk.Context, contractID string, class collection.TokenClass) {
	store := k.storeService.OpenKVStore(ctx)
	key := classKey(contractID, class.GetId())

	bz, err := k.cdc.MarshalInterface(class)
	if err != nil {
		panic(err)
	}
	err = store.Set(key, bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) getNextClassIDs(ctx sdk.Context, contractID string) collection.NextClassIDs {
	store := k.storeService.OpenKVStore(ctx)
	key := nextClassIDKey(contractID)
	bz, _ := store.Get(key)
	if bz == nil {
		panic(sdkerrors.ErrNotFound.Wrapf("no next class ids of contract %s", contractID))
	}

	var class collection.NextClassIDs
	if err := class.Unmarshal(bz); err != nil {
		panic(err)
	}
	return class
}

func (k Keeper) setNextClassIDs(ctx sdk.Context, ids collection.NextClassIDs) {
	store := k.storeService.OpenKVStore(ctx)
	key := nextClassIDKey(ids.ContractId)

	bz, err := ids.Marshal()
	if err != nil {
		panic(err)
	}
	err = store.Set(key, bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) MintNFT(ctx sdk.Context, contractID string, to sdk.AccAddress, params []collection.MintNFTParam) ([]collection.NFT, error) {
	tokens := make([]collection.NFT, 0, len(params))
	for _, param := range params {
		classID := param.TokenType
		class, err := k.GetTokenClass(ctx, contractID, classID)
		if err != nil {
			return nil, collection.ErrTokenTypeNotExist.Wrap(err.Error())
		}

		if _, ok := class.(*collection.NFTClass); !ok {
			return nil, collection.ErrTokenTypeNotExist.Wrapf("not a class of non-fungible token: %s", classID)
		}

		nextTokenID := k.getNextTokenID(ctx, contractID, classID)
		k.setNextTokenID(ctx, contractID, classID, nextTokenID.Incr())
		tokenID := collection.NewNFTID(classID, int(nextTokenID.Uint64()))

		amount := math.OneInt()

		k.setBalance(ctx, contractID, to, tokenID, amount)
		k.setOwner(ctx, contractID, tokenID, to)

		token := collection.NFT{
			TokenId: tokenID,
			Name:    param.Name,
			Meta:    param.Meta,
		}
		k.setNFT(ctx, contractID, token)

		// update statistics
		supply := k.GetSupply(ctx, contractID, classID)
		k.setSupply(ctx, contractID, classID, supply.Add(amount))

		minted := k.GetMinted(ctx, contractID, classID)
		k.setMinted(ctx, contractID, classID, minted.Add(amount))

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (k Keeper) BurnCoins(ctx sdk.Context, contractID string, from sdk.AccAddress, amount []collection.Coin) ([]collection.Coin, error) {
	if err := k.subtractCoins(ctx, contractID, from, amount); err != nil {
		return nil, err
	}

	burntAmount := []collection.Coin{}
	for _, coin := range amount {
		burntAmount = append(burntAmount, coin)
		k.deleteNFT(ctx, contractID, coin.TokenId)
	}

	// update statistics
	for _, coin := range burntAmount {
		classID := collection.SplitTokenID(coin.TokenId)
		supply := k.GetSupply(ctx, contractID, classID)
		k.setSupply(ctx, contractID, classID, supply.Sub(coin.Amount))

		burnt := k.GetBurnt(ctx, contractID, classID)
		k.setBurnt(ctx, contractID, classID, burnt.Add(coin.Amount))
	}

	return burntAmount, nil
}

func (k Keeper) getNextTokenID(ctx sdk.Context, contractID, classID string) math.Uint {
	store := k.storeService.OpenKVStore(ctx)
	key := nextTokenIDKey(contractID, classID)
	bz, _ := store.Get(key)
	if bz == nil {
		panic(sdkerrors.ErrNotFound.Wrapf("no next token id of token class %s", classID))
	}

	var id math.Uint
	if err := id.Unmarshal(bz); err != nil {
		panic(err)
	}
	return id
}

func (k Keeper) setNextTokenID(ctx sdk.Context, contractID, classID string, tokenID math.Uint) {
	store := k.storeService.OpenKVStore(ctx)
	key := nextTokenIDKey(contractID, classID)

	bz, err := tokenID.Marshal()
	if err != nil {
		panic(err)
	}
	err = store.Set(key, bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) ModifyContract(ctx sdk.Context, contractID string, changes []collection.Attribute) {
	contract, err := k.GetContract(ctx, contractID)
	if err != nil {
		panic(err)
	}

	modifiers := map[collection.AttributeKey]func(string){
		collection.AttributeKeyName: func(name string) {
			contract.Name = name
		},
		collection.AttributeKeyURI: func(uri string) {
			contract.Uri = uri
		},
		collection.AttributeKeyMeta: func(meta string) {
			contract.Meta = meta
		},
	}
	for _, change := range changes {
		key := collection.AttributeKeyFromString(change.Key)
		modifiers[key](change.Value)
	}

	k.setContract(ctx, *contract)
}

func (k Keeper) ModifyTokenClass(ctx sdk.Context, contractID, classID string, changes []collection.Attribute) error {
	class, err := k.GetTokenClass(ctx, contractID, classID)
	if err != nil {
		if err := collection.ValidateLegacyNFTClassID(classID); err == nil {
			return collection.ErrTokenTypeNotExist.Wrap(classID)
		}

		panic(err)
	}

	modifiers := map[collection.AttributeKey]func(string){
		collection.AttributeKeyName: func(name string) {
			class.SetName(name)
		},
		collection.AttributeKeyMeta: func(meta string) {
			class.SetMeta(meta)
		},
	}
	for _, change := range changes {
		key := collection.AttributeKeyFromString(change.Key)
		modifiers[key](change.Value)
	}

	k.setTokenClass(ctx, contractID, class)

	return nil
}

func (k Keeper) ModifyNFT(ctx sdk.Context, contractID, tokenID string, changes []collection.Attribute) error {
	token, err := k.GetNFT(ctx, contractID, tokenID)
	if err != nil {
		return err
	}

	modifiers := map[collection.AttributeKey]func(string){
		collection.AttributeKeyName: func(name string) {
			token.Name = name
		},
		collection.AttributeKeyMeta: func(meta string) {
			token.Meta = meta
		},
	}
	for _, change := range changes {
		key := collection.AttributeKeyFromString(change.Key)
		modifiers[key](change.Value)
	}

	k.setNFT(ctx, contractID, *token)

	return nil
}

func (k Keeper) Grant(ctx sdk.Context, contractID string, granter, grantee sdk.AccAddress, permission collection.Permission) {
	k.grant(ctx, contractID, grantee, permission)

	event := collection.EventGranted{
		ContractId: contractID,
		Granter:    k.bytesToString(granter),
		Grantee:    k.bytesToString(grantee),
		Permission: permission,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}
}

func (k Keeper) grant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission collection.Permission) {
	k.setGrant(ctx, contractID, grantee, permission)
}

func (k Keeper) Abandon(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission collection.Permission) {
	k.deleteGrant(ctx, contractID, grantee, permission)

	event := collection.EventRenounced{
		ContractId: contractID,
		Grantee:    k.bytesToString(grantee),
		Permission: permission,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}
}

func (k Keeper) GetGrant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission collection.Permission) (*collection.Grant, error) {
	store := k.storeService.OpenKVStore(ctx)
	if ok, _ := store.Has(grantKey(contractID, grantee, permission)); ok {
		return &collection.Grant{
			Grantee:    k.bytesToString(grantee),
			Permission: permission,
		}, nil
	}
	return nil, sdkerrors.ErrNotFound.Wrapf("no %s permission granted on %s", permission, grantee)
}

func (k Keeper) setGrant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission collection.Permission) {
	store := k.storeService.OpenKVStore(ctx)
	key := grantKey(contractID, grantee, permission)
	err := store.Set(key, []byte{})
	if err != nil {
		panic(err)
	}
}

func (k Keeper) deleteGrant(ctx sdk.Context, contractID string, grantee sdk.AccAddress, permission collection.Permission) {
	store := k.storeService.OpenKVStore(ctx)
	key := grantKey(contractID, grantee, permission)
	err := store.Delete(key)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) getStatistic(ctx sdk.Context, keyPrefix []byte, contractID, classID string) math.Int {
	store := k.storeService.OpenKVStore(ctx)
	amount := math.ZeroInt()
	bz, _ := store.Get(statisticKey(keyPrefix, contractID, classID))
	if bz != nil {
		if err := amount.Unmarshal(bz); err != nil {
			panic(err)
		}
	}

	return amount
}

func (k Keeper) setStatistic(ctx sdk.Context, keyPrefix []byte, contractID, classID string, amount math.Int) {
	store := k.storeService.OpenKVStore(ctx)
	key := statisticKey(keyPrefix, contractID, classID)
	if amount.IsZero() {
		err := store.Delete(key)
		if err != nil {
			panic(err)
		}
	} else {
		bz, err := amount.Marshal()
		if err != nil {
			panic(err)
		}
		err = store.Set(key, bz)
		if err != nil {
			panic(err)
		}
	}
}

func (k Keeper) GetSupply(ctx sdk.Context, contractID, classID string) math.Int {
	return k.getStatistic(ctx, supplyKeyPrefix, contractID, classID)
}

func (k Keeper) GetMinted(ctx sdk.Context, contractID, classID string) math.Int {
	return k.getStatistic(ctx, mintedKeyPrefix, contractID, classID)
}

func (k Keeper) GetBurnt(ctx sdk.Context, contractID, classID string) math.Int {
	return k.getStatistic(ctx, burntKeyPrefix, contractID, classID)
}

func (k Keeper) setSupply(ctx sdk.Context, contractID, classID string, amount math.Int) {
	k.setStatistic(ctx, supplyKeyPrefix, contractID, classID, amount)
}

func (k Keeper) setMinted(ctx sdk.Context, contractID, classID string, amount math.Int) {
	k.setStatistic(ctx, mintedKeyPrefix, contractID, classID, amount)
}

func (k Keeper) setBurnt(ctx sdk.Context, contractID, classID string, amount math.Int) {
	k.setStatistic(ctx, burntKeyPrefix, contractID, classID, amount)
}
