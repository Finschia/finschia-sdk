package keeper

import (
	"context"

	"github.com/cosmos/gogoproto/proto"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/collection"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServer returns an implementation of the collection MsgServer interface
// for the provided Keeper.
func NewMsgServer(keeper Keeper) collection.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ collection.MsgServer = (*msgServer)(nil)

func (s msgServer) SendNFT(c context.Context, req *collection.MsgSendNFT) (*collection.MsgSendNFTResponse, error) {
	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.From); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}
	if _, err := s.keeper.addressCodec.StringToBytes(req.To); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", req.To)
	}

	if len(req.TokenIds) == 0 {
		return nil, collection.ErrEmptyField.Wrap("token ids cannot be empty")
	}
	for _, id := range req.TokenIds {
		if err := collection.ValidateTokenID(id); err != nil {
			return nil, err
		}
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr, err := s.keeper.addressCodec.StringToBytes(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrap(err.Error())
	}

	amount := make([]collection.Coin, len(req.TokenIds))
	for i, id := range req.TokenIds {
		amount[i] = collection.Coin{TokenId: id, Amount: math.OneInt()}

		// legacy
		if err := s.keeper.hasNFT(ctx, req.ContractId, id); err != nil {
			return nil, err
		}
		if !s.keeper.getOwner(ctx, req.ContractId, id).Equals(sdk.AccAddress(fromAddr)) {
			return nil, collection.ErrTokenNotOwnedBy.Wrapf("%s does not have %s", fromAddr, id)
		}
	}

	toAddr, err := s.keeper.addressCodec.StringToBytes(req.To)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}

	if err := s.keeper.SendCoins(ctx, req.ContractId, fromAddr, toAddr, amount); err != nil {
		panic(err)
	}

	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		To:         req.To,
		Amount:     amount,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgSendNFTResponse{}, nil
}

func (s msgServer) OperatorSendNFT(c context.Context, req *collection.MsgOperatorSendNFT) (*collection.MsgOperatorSendNFTResponse, error) {
	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.Operator); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", req.Operator)
	}
	if _, err := s.keeper.addressCodec.StringToBytes(req.From); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}
	if _, err := s.keeper.addressCodec.StringToBytes(req.To); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", req.To)
	}

	if len(req.TokenIds) == 0 {
		return nil, collection.ErrEmptyField.Wrap("token ids cannot be empty")
	}
	for _, id := range req.TokenIds {
		if err := collection.ValidateTokenID(id); err != nil {
			return nil, err
		}
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	operatorAddr, err := s.keeper.addressCodec.StringToBytes(req.Operator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", err)
	}
	fromAddr, err := s.keeper.addressCodec.StringToBytes(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, operatorAddr); err != nil {
		return nil, collection.ErrCollectionNotApproved.Wrap(err.Error())
	}

	amount := make([]collection.Coin, len(req.TokenIds))
	for i, id := range req.TokenIds {
		amount[i] = collection.Coin{TokenId: id, Amount: math.OneInt()}

		// legacy
		if err := s.keeper.hasNFT(ctx, req.ContractId, id); err != nil {
			return nil, err
		}
		if !s.keeper.getOwner(ctx, req.ContractId, id).Equals(sdk.AccAddress(fromAddr)) {
			return nil, collection.ErrTokenNotOwnedBy.Wrapf("%s does not have %s", fromAddr, id)
		}
	}

	toAddr, err := s.keeper.addressCodec.StringToBytes(req.To)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if err := s.keeper.SendCoins(ctx, req.ContractId, fromAddr, toAddr, amount); err != nil {
		panic(err)
	}

	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		From:       req.From,
		To:         req.To,
		Amount:     amount,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgOperatorSendNFTResponse{}, nil
}

func (s msgServer) AuthorizeOperator(c context.Context, req *collection.MsgAuthorizeOperator) (*collection.MsgAuthorizeOperatorResponse, error) {
	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.Holder); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid holder address: %s", req.Holder)
	}
	if _, err := s.keeper.addressCodec.StringToBytes(req.Operator); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", req.Operator)
	}

	if req.Operator == req.Holder {
		return nil, collection.ErrApproverProxySame
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	holderAddr, err := s.keeper.addressCodec.StringToBytes(req.Holder)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid holder address: %s", err)
	}
	operatorAddr, err := s.keeper.addressCodec.StringToBytes(req.Operator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", err)
	}

	if err := s.keeper.AuthorizeOperator(ctx, req.ContractId, holderAddr, operatorAddr); err != nil {
		return nil, err
	}

	event := collection.EventAuthorizedOperator{
		ContractId: req.ContractId,
		Holder:     req.Holder,
		Operator:   req.Operator,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgAuthorizeOperatorResponse{}, nil
}

func (s msgServer) RevokeOperator(c context.Context, req *collection.MsgRevokeOperator) (*collection.MsgRevokeOperatorResponse, error) {
	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.Holder); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid holder address: %s", req.Holder)
	}
	if _, err := s.keeper.addressCodec.StringToBytes(req.Operator); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", req.Operator)
	}

	if req.Operator == req.Holder {
		return nil, collection.ErrApproverProxySame
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	holderAddr, err := s.keeper.addressCodec.StringToBytes(req.Holder)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid holder address: %s", err)
	}
	operatorAddr, err := s.keeper.addressCodec.StringToBytes(req.Operator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", err)
	}

	if err := s.keeper.RevokeOperator(ctx, req.ContractId, holderAddr, operatorAddr); err != nil {
		return nil, err
	}

	event := collection.EventRevokedOperator{
		ContractId: req.ContractId,
		Holder:     req.Holder,
		Operator:   req.Operator,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgRevokeOperatorResponse{}, nil
}

func (s msgServer) CreateContract(c context.Context, req *collection.MsgCreateContract) (*collection.MsgCreateContractResponse, error) {
	if _, err := s.keeper.addressCodec.StringToBytes(req.Owner); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", req.Owner)
	}

	if err := collection.ValidateName(req.Name); err != nil {
		return nil, err
	}

	if err := collection.ValidateURI(req.Uri); err != nil {
		return nil, err
	}

	if err := collection.ValidateMeta(req.Meta); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	contract := collection.Contract{
		Name: req.Name,
		Uri:  req.Uri,
		Meta: req.Meta,
	}
	ownerAddr, err := s.keeper.addressCodec.StringToBytes(req.Owner)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
	}

	id := s.keeper.CreateContract(ctx, ownerAddr, contract)

	return &collection.MsgCreateContractResponse{ContractId: id}, nil
}

func (s msgServer) IssueNFT(c context.Context, req *collection.MsgIssueNFT) (*collection.MsgIssueNFTResponse, error) {
	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if err := collection.ValidateName(req.Name); err != nil {
		return nil, err
	}

	if err := collection.ValidateMeta(req.Meta); err != nil {
		return nil, err
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.Owner); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", req.Owner)
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	ownerAddr, err := s.keeper.addressCodec.StringToBytes(req.Owner)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", err)
	}

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, ownerAddr, collection.PermissionIssue); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrap(err.Error())
	}

	class := &collection.NFTClass{
		Name: req.Name,
		Meta: req.Meta,
	}
	id, err := s.keeper.CreateTokenClass(ctx, req.ContractId, class)
	if err != nil {
		return nil, err
	}

	event := collection.EventCreatedNFTClass{
		ContractId: req.ContractId,
		Operator:   req.Owner,
		TokenType:  *id,
		Name:       class.Name,
		Meta:       class.Meta,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	for _, permission := range []collection.Permission{
		collection.PermissionMint,
		collection.PermissionBurn,
	} {
		s.keeper.Grant(ctx, req.ContractId, []byte{}, ownerAddr, permission)
	}

	return &collection.MsgIssueNFTResponse{TokenType: *id}, nil
}

func (s msgServer) MintNFT(c context.Context, req *collection.MsgMintNFT) (*collection.MsgMintNFTResponse, error) {
	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.From); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}
	if _, err := s.keeper.addressCodec.StringToBytes(req.To); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", req.To)
	}

	if len(req.Params) == 0 {
		return nil, collection.ErrEmptyField.Wrap("mint params cannot be empty")
	}
	for _, param := range req.Params {
		classID := param.TokenType
		if err := collection.ValidateLegacyNFTClassID(classID); err != nil {
			return nil, err
		}

		if len(param.Name) == 0 {
			return nil, collection.ErrInvalidTokenName
		}
		if err := collection.ValidateName(param.Name); err != nil {
			return nil, err
		}

		if err := collection.ValidateMeta(param.Meta); err != nil {
			return nil, err
		}
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr, err := s.keeper.addressCodec.StringToBytes(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, fromAddr, collection.PermissionMint); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrap(err.Error())
	}

	toAddr, err := s.keeper.addressCodec.StringToBytes(req.To)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}

	tokens, err := s.keeper.MintNFT(ctx, req.ContractId, toAddr, req.Params)
	if err != nil {
		return nil, err
	}

	event := collection.EventMintedNFT{
		ContractId: req.ContractId,
		Operator:   req.From,
		To:         req.To,
		Tokens:     tokens,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	tokenIDs := make([]string, 0, len(tokens))
	for _, token := range tokens {
		tokenIDs = append(tokenIDs, token.TokenId)
	}
	return &collection.MsgMintNFTResponse{TokenIds: tokenIDs}, nil
}

func (s msgServer) BurnNFT(c context.Context, req *collection.MsgBurnNFT) (*collection.MsgBurnNFTResponse, error) {
	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.From); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}

	if len(req.TokenIds) == 0 {
		return nil, collection.ErrEmptyField.Wrap("token ids cannot be empty")
	}
	for _, id := range req.TokenIds {
		if err := collection.ValidateLegacyNFTID(id); err != nil {
			return nil, err
		}
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr, err := s.keeper.addressCodec.StringToBytes(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, fromAddr, collection.PermissionBurn); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrap(err.Error())
	}

	coins := make([]collection.Coin, 0, len(req.TokenIds))
	for _, id := range req.TokenIds {
		coins = append(coins, collection.NewCoin(id, math.OneInt()))

		// legacy
		if err := s.keeper.hasNFT(ctx, req.ContractId, id); err != nil {
			return nil, err
		}
		if !s.keeper.getOwner(ctx, req.ContractId, id).Equals(sdk.AccAddress(fromAddr)) {
			return nil, collection.ErrTokenNotOwnedBy.Wrapf("%s does not have %s", fromAddr, id)
		}
	}

	burnt, err := s.keeper.BurnCoins(ctx, req.ContractId, fromAddr, coins)
	if err != nil {
		panic(err)
	}

	// emit events against all burnt tokens.
	event := collection.EventBurned{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		Amount:     burnt,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgBurnNFTResponse{}, nil
}

func (s msgServer) OperatorBurnNFT(c context.Context, req *collection.MsgOperatorBurnNFT) (*collection.MsgOperatorBurnNFTResponse, error) {
	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.Operator); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", req.Operator)
	}
	if _, err := s.keeper.addressCodec.StringToBytes(req.From); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}

	if len(req.TokenIds) == 0 {
		return nil, collection.ErrEmptyField.Wrap("token ids cannot be empty")
	}
	for _, id := range req.TokenIds {
		if err := collection.ValidateLegacyNFTID(id); err != nil {
			return nil, err
		}
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr, err := s.keeper.addressCodec.StringToBytes(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}
	operatorAddr, err := s.keeper.addressCodec.StringToBytes(req.Operator)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", err)
	}

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, operatorAddr); err != nil {
		return nil, collection.ErrCollectionNotApproved.Wrap(err.Error())
	}

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, operatorAddr, collection.PermissionBurn); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrap(err.Error())
	}

	coins := make([]collection.Coin, 0, len(req.TokenIds))
	for _, id := range req.TokenIds {
		coins = append(coins, collection.NewCoin(id, math.OneInt()))

		// legacy
		if err := s.keeper.hasNFT(ctx, req.ContractId, id); err != nil {
			return nil, err
		}
		if !s.keeper.getOwner(ctx, req.ContractId, id).Equals(sdk.AccAddress(fromAddr)) {
			return nil, collection.ErrTokenNotOwnedBy.Wrapf("%s does not have %s", fromAddr, id)
		}
	}

	burnt, err := s.keeper.BurnCoins(ctx, req.ContractId, fromAddr, coins)
	if err != nil {
		panic(err)
	}

	// emit events against all burnt tokens.
	event := collection.EventBurned{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		From:       req.From,
		Amount:     burnt,
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgOperatorBurnNFTResponse{}, nil
}

func (s msgServer) Modify(c context.Context, req *collection.MsgModify) (*collection.MsgModifyResponse, error) {
	for i, change := range req.Changes {
		key := change.Key
		converted := collection.AttrCanonicalKey(key)
		if converted != key {
			req.Changes[i].Key = converted
		}
	}
	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.Owner); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid owner address: %s", req.Owner)
	}

	if len(req.TokenType) != 0 {
		classID := req.TokenType
		if err := collection.ValidateClassID(classID); err != nil {
			return nil, collection.ErrInvalidTokenType.Wrap(err.Error())
		}
	}

	if len(req.TokenIndex) != 0 {
		tokenID := req.TokenType + req.TokenIndex
		if err := collection.ValidateTokenID(tokenID); err != nil {
			return nil, collection.ErrInvalidTokenIndex.Wrap(err.Error())
		}
		// reject modifying nft class with token index filled (daphne compat.)
		if collection.ValidateLegacyIdxNFT(tokenID) == nil {
			return nil, collection.ErrInvalidTokenIndex.Wrap("cannot modify nft class with index filled")
		}
	}

	validator := collection.ValidateTokenClassChange
	if len(req.TokenType) == 0 {
		if len(req.TokenIndex) == 0 {
			validator = collection.ValidateContractChange
		} else {
			return nil, collection.ErrTokenIndexWithoutType.Wrap("token index without type")
		}
	}
	if len(req.Changes) == 0 {
		return nil, collection.ErrEmptyChanges.Wrap("empty changes")
	}
	if len(req.Changes) > collection.ChangesLimit {
		return nil, collection.ErrInvalidChangesFieldCount.Wrapf("the number of changes exceeds the limit: %d > %d", len(req.Changes), collection.ChangesLimit)
	}
	seenKeys := map[string]bool{}
	for _, change := range req.Changes {
		key := change.Key
		if seenKeys[key] {
			return nil, collection.ErrDuplicateChangesField.Wrapf("duplicate keys: %s", change.Key)
		}
		seenKeys[key] = true

		attribute := collection.Attribute{
			Key:   change.Key,
			Value: change.Value,
		}
		if err := validator(attribute); err != nil {
			return nil, err
		}
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	operator, err := s.keeper.addressCodec.StringToBytes(req.Owner)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid operator address: %s", err)
	}

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, operator, collection.PermissionModify); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrap(err.Error())
	}

	// copied from daphne
	modify := func(tokenType, tokenIndex string) error {
		changes := make([]collection.Attribute, len(req.Changes))
		for i, change := range req.Changes {
			changes[i] = collection.Attribute{
				Key:   change.Key,
				Value: change.Value,
			}
		}

		classID := tokenType
		tokenID := classID + tokenIndex
		if tokenType != "" {
			if tokenIndex != "" && collection.ValidateNFTID(tokenID) == nil {
				event := collection.EventModifiedNFT{
					ContractId: req.ContractId,
					Operator:   req.Owner,
					TokenId:    tokenID,
					Changes:    changes,
				}
				if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
					panic(err)
				}

				return s.keeper.ModifyNFT(ctx, req.ContractId, tokenID, changes)
			}

			event := collection.EventModifiedTokenClass{
				ContractId: req.ContractId,
				Operator:   req.Owner,
				TokenType:  classID,
				Changes:    changes,
				TypeName:   proto.MessageName(&collection.NFTClass{}),
			}
			if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
				panic(err)
			}

			return s.keeper.ModifyTokenClass(ctx, req.ContractId, classID, changes)
		}
		if req.TokenIndex == "" {
			event := collection.EventModifiedContract{
				ContractId: req.ContractId,
				Operator:   req.Owner,
				Changes:    changes,
			}
			if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
				panic(err)
			}

			s.keeper.ModifyContract(ctx, req.ContractId, changes)
			return nil
		}

		panic(sdkerrors.ErrInvalidRequest.Wrap("token index without type"))
	}

	if err := modify(req.TokenType, req.TokenIndex); err != nil {
		return nil, err
	}

	return &collection.MsgModifyResponse{}, nil
}

func (s msgServer) GrantPermission(c context.Context, req *collection.MsgGrantPermission) (*collection.MsgGrantPermissionResponse, error) {
	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.From); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}
	if _, err := s.keeper.addressCodec.StringToBytes(req.To); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", req.To)
	}

	if err := collection.ValidateLegacyPermission(req.Permission); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	granter, err := s.keeper.addressCodec.StringToBytes(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}
	grantee, err := s.keeper.addressCodec.StringToBytes(req.To)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid to address: %s", err)
	}
	permission := collection.Permission(collection.LegacyPermissionFromString(req.Permission))

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, granter, permission); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrapf("%s is not authorized for %s", granter, permission)
	}

	// it emits typed event inside s.keeper.Grant()
	s.keeper.Grant(ctx, req.ContractId, granter, grantee, permission)

	return &collection.MsgGrantPermissionResponse{}, nil
}

func (s msgServer) RevokePermission(c context.Context, req *collection.MsgRevokePermission) (*collection.MsgRevokePermissionResponse, error) {
	if err := collection.ValidateContractID(req.ContractId); err != nil {
		return nil, err
	}

	if _, err := s.keeper.addressCodec.StringToBytes(req.From); err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", req.From)
	}

	if err := collection.ValidateLegacyPermission(req.Permission); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	grantee, err := s.keeper.addressCodec.StringToBytes(req.From)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid from address: %s", err)
	}
	permission := collection.Permission(collection.LegacyPermissionFromString(req.Permission))

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, permission); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrapf("%s is not authorized for %s", grantee, permission)
	}

	// it emits typed event inside s.keeper.Abandon()
	s.keeper.Abandon(ctx, req.ContractId, grantee, permission)

	return &collection.MsgRevokePermissionResponse{}, nil
}
