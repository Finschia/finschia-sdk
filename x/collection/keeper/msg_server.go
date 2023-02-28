package keeper

import (
	"context"

	"github.com/gogo/protobuf/proto"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/collection"
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

func (s msgServer) SendFT(c context.Context, req *collection.MsgSendFT) (*collection.MsgSendFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)
	toAddr := sdk.MustAccAddressFromBech32(req.To)

	if err := s.keeper.SendCoins(ctx, req.ContractId, fromAddr, toAddr, req.Amount); err != nil {
		return nil, err
	}

	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	ctx.EventManager().EmitEvent(*collection.NewEventTransferFT(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgSendFTResponse{}, nil
}

func (s msgServer) OperatorSendFT(c context.Context, req *collection.MsgOperatorSendFT) (*collection.MsgOperatorSendFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	operatorAddr := sdk.MustAccAddressFromBech32(req.Operator)
	fromAddr := sdk.MustAccAddressFromBech32(req.From)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, operatorAddr); err != nil {
		return nil, collection.ErrCollectionNotApproved.Wrap(err.Error())
	}

	toAddr := sdk.MustAccAddressFromBech32(req.To)

	if err := s.keeper.SendCoins(ctx, req.ContractId, fromAddr, toAddr, req.Amount); err != nil {
		return nil, err
	}

	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		From:       req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	ctx.EventManager().EmitEvent(*collection.NewEventTransferFTFrom(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgOperatorSendFTResponse{}, nil
}

func (s msgServer) SendNFT(c context.Context, req *collection.MsgSendNFT) (*collection.MsgSendNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)

	amount := make([]collection.Coin, len(req.TokenIds))
	for i, id := range req.TokenIds {
		amount[i] = collection.Coin{TokenId: id, Amount: sdk.OneInt()}

		// legacy
		if err := s.keeper.hasNFT(ctx, req.ContractId, id); err != nil {
			return nil, err
		}
		if _, err := s.keeper.GetParent(ctx, req.ContractId, id); err == nil {
			return nil, collection.ErrTokenCannotTransferChildToken.Wrap(id)
		}
		if !s.keeper.getOwner(ctx, req.ContractId, id).Equals(fromAddr) {
			return nil, collection.ErrTokenNotOwnedBy.Wrapf("%s does not have %s", fromAddr, id)
		}
	}

	// emit legacy events
	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		To:         req.To,
		Amount:     amount,
	}
	ctx.EventManager().EmitEvents(collection.NewEventTransferNFT(event))

	toAddr := sdk.MustAccAddressFromBech32(req.To)

	if err := s.keeper.SendCoins(ctx, req.ContractId, fromAddr, toAddr, amount); err != nil {
		panic(err)
	}

	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgSendNFTResponse{}, nil
}

func (s msgServer) OperatorSendNFT(c context.Context, req *collection.MsgOperatorSendNFT) (*collection.MsgOperatorSendNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	operatorAddr := sdk.MustAccAddressFromBech32(req.Operator)
	fromAddr := sdk.MustAccAddressFromBech32(req.From)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, operatorAddr); err != nil {
		return nil, collection.ErrCollectionNotApproved.Wrap(err.Error())
	}

	amount := make([]collection.Coin, len(req.TokenIds))
	for i, id := range req.TokenIds {
		amount[i] = collection.Coin{TokenId: id, Amount: sdk.OneInt()}

		// legacy
		if err := s.keeper.hasNFT(ctx, req.ContractId, id); err != nil {
			return nil, err
		}
		if _, err := s.keeper.GetParent(ctx, req.ContractId, id); err == nil {
			return nil, collection.ErrTokenCannotTransferChildToken.Wrap(id)
		}
		if !s.keeper.getOwner(ctx, req.ContractId, id).Equals(fromAddr) {
			return nil, collection.ErrTokenNotOwnedBy.Wrapf("%s does not have %s", fromAddr, id)
		}
	}

	// emit legacy events
	event := collection.EventSent{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		From:       req.From,
		To:         req.To,
		Amount:     amount,
	}
	ctx.EventManager().EmitEvents(collection.NewEventTransferNFTFrom(event))

	toAddr := sdk.MustAccAddressFromBech32(req.To)

	if err := s.keeper.SendCoins(ctx, req.ContractId, fromAddr, toAddr, amount); err != nil {
		panic(err)
	}

	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgOperatorSendNFTResponse{}, nil
}

func (s msgServer) AuthorizeOperator(c context.Context, req *collection.MsgAuthorizeOperator) (*collection.MsgAuthorizeOperatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	holderAddr := sdk.MustAccAddressFromBech32(req.Holder)
	operatorAddr := sdk.MustAccAddressFromBech32(req.Operator)

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
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	holderAddr := sdk.MustAccAddressFromBech32(req.Holder)
	operatorAddr := sdk.MustAccAddressFromBech32(req.Operator)

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
	ctx := sdk.UnwrapSDKContext(c)

	contract := collection.Contract{
		Name: req.Name,
		Uri:  req.Uri,
		Meta: req.Meta,
	}
	ownerAddr := sdk.MustAccAddressFromBech32(req.Owner)

	id := s.keeper.CreateContract(ctx, ownerAddr, contract)

	return &collection.MsgCreateContractResponse{ContractId: id}, nil
}

func (s msgServer) IssueFT(c context.Context, req *collection.MsgIssueFT) (*collection.MsgIssueFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	ownerAddr := sdk.MustAccAddressFromBech32(req.Owner)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, ownerAddr, collection.PermissionIssue); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrap(err.Error())
	}

	class := &collection.FTClass{
		Name:     req.Name,
		Meta:     req.Meta,
		Decimals: req.Decimals,
		Mintable: req.Mintable,
	}
	id, err := s.keeper.CreateTokenClass(ctx, req.ContractId, class)
	if err != nil {
		return nil, err
	}

	event := collection.EventCreatedFTClass{
		ContractId: req.ContractId,
		Operator:   req.Owner,
		TokenId:    collection.NewFTID(*id),
		Name:       class.Name,
		Meta:       class.Meta,
		Decimals:   class.Decimals,
		Mintable:   class.Mintable,
	}

	toAddr := sdk.MustAccAddressFromBech32(req.To)

	ctx.EventManager().EmitEvent(collection.NewEventIssueFT(event, toAddr, req.Amount))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	// supply tokens
	if req.Amount.IsPositive() {
		s.keeper.mintFT(ctx, req.ContractId, toAddr, *id, req.Amount)

		event := collection.EventMintedFT{
			ContractId: req.ContractId,
			Operator:   req.Owner,
			To:         req.To,
			Amount:     collection.NewCoins(collection.NewFTCoin(*id, req.Amount)),
		}
		if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
			panic(err)
		}
	}

	return &collection.MsgIssueFTResponse{TokenId: *id}, nil
}

func (s msgServer) IssueNFT(c context.Context, req *collection.MsgIssueNFT) (*collection.MsgIssueNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	ownerAddr := sdk.MustAccAddressFromBech32(req.Owner)

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
	ctx.EventManager().EmitEvent(collection.NewEventIssueNFT(event))
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

func (s msgServer) MintFT(c context.Context, req *collection.MsgMintFT) (*collection.MsgMintFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, fromAddr, collection.PermissionMint); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrap(err.Error())
	}

	toAddr := sdk.MustAccAddressFromBech32(req.To)

	if err := s.keeper.MintFT(ctx, req.ContractId, toAddr, req.Amount); err != nil {
		return nil, err
	}

	event := collection.EventMintedFT{
		ContractId: req.ContractId,
		Operator:   req.From,
		To:         req.To,
		Amount:     req.Amount,
	}
	ctx.EventManager().EmitEvent(collection.NewEventMintFT(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgMintFTResponse{}, nil
}

func (s msgServer) MintNFT(c context.Context, req *collection.MsgMintNFT) (*collection.MsgMintNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, fromAddr, collection.PermissionMint); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrap(err.Error())
	}

	toAddr := sdk.MustAccAddressFromBech32(req.To)

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
	ctx.EventManager().EmitEvents(collection.NewEventMintNFT(event))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	tokenIDs := make([]string, 0, len(tokens))
	for _, token := range tokens {
		tokenIDs = append(tokenIDs, token.TokenId)
	}
	return &collection.MsgMintNFTResponse{TokenIds: tokenIDs}, nil
}

func (s msgServer) BurnFT(c context.Context, req *collection.MsgBurnFT) (*collection.MsgBurnFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, fromAddr, collection.PermissionBurn); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrap(err.Error())
	}

	burnt, err := s.keeper.BurnCoins(ctx, req.ContractId, fromAddr, req.Amount)
	if err != nil {
		return nil, err
	}

	event := collection.EventBurned{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		Amount:     burnt,
	}
	if e := collection.NewEventBurnFT(event); e != nil {
		ctx.EventManager().EmitEvent(*e)
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgBurnFTResponse{}, nil
}

func (s msgServer) OperatorBurnFT(c context.Context, req *collection.MsgOperatorBurnFT) (*collection.MsgOperatorBurnFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)
	operatorAddr := sdk.MustAccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, operatorAddr); err != nil {
		return nil, collection.ErrCollectionNotApproved.Wrap(err.Error())
	}

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, operatorAddr, collection.PermissionBurn); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrap(err.Error())
	}

	burnt, err := s.keeper.BurnCoins(ctx, req.ContractId, fromAddr, req.Amount)
	if err != nil {
		return nil, err
	}

	event := collection.EventBurned{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		From:       req.From,
		Amount:     burnt,
	}
	if e := collection.NewEventBurnFTFrom(event); e != nil {
		ctx.EventManager().EmitEvent(*e)
	}
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgOperatorBurnFTResponse{}, nil
}

func (s msgServer) BurnNFT(c context.Context, req *collection.MsgBurnNFT) (*collection.MsgBurnNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, fromAddr, collection.PermissionBurn); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrap(err.Error())
	}

	coins := make([]collection.Coin, 0, len(req.TokenIds))
	for _, id := range req.TokenIds {
		coins = append(coins, collection.NewCoin(id, sdk.OneInt()))

		// legacy
		if err := s.keeper.hasNFT(ctx, req.ContractId, id); err != nil {
			return nil, err
		}
		if _, err := s.keeper.GetParent(ctx, req.ContractId, id); err == nil {
			return nil, collection.ErrBurnNonRootNFT.Wrap(id)
		}
		if !s.keeper.getOwner(ctx, req.ContractId, id).Equals(fromAddr) {
			return nil, collection.ErrTokenNotOwnedBy.Wrapf("%s does not have %s", fromAddr, id)
		}
	}

	// legacy: emit events against the original request.
	event := collection.EventBurned{
		ContractId: req.ContractId,
		Operator:   req.From,
		From:       req.From,
		Amount:     coins,
	}
	ctx.EventManager().EmitEvents(collection.NewEventBurnNFT(event))

	burnt, err := s.keeper.BurnCoins(ctx, req.ContractId, fromAddr, coins)
	if err != nil {
		panic(err)
	}

	// emit events against all burnt tokens.
	event.Amount = burnt
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgBurnNFTResponse{}, nil
}

func (s msgServer) OperatorBurnNFT(c context.Context, req *collection.MsgOperatorBurnNFT) (*collection.MsgOperatorBurnNFTResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)
	operatorAddr := sdk.MustAccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, operatorAddr); err != nil {
		return nil, collection.ErrCollectionNotApproved.Wrap(err.Error())
	}

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, operatorAddr, collection.PermissionBurn); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrap(err.Error())
	}

	coins := make([]collection.Coin, 0, len(req.TokenIds))
	for _, id := range req.TokenIds {
		coins = append(coins, collection.NewCoin(id, sdk.OneInt()))

		// legacy
		if err := s.keeper.hasNFT(ctx, req.ContractId, id); err != nil {
			return nil, err
		}
		if _, err := s.keeper.GetParent(ctx, req.ContractId, id); err == nil {
			return nil, collection.ErrBurnNonRootNFT.Wrap(id)
		}
		if !s.keeper.getOwner(ctx, req.ContractId, id).Equals(fromAddr) {
			return nil, collection.ErrTokenNotOwnedBy.Wrapf("%s does not have %s", fromAddr, id)
		}
	}

	// legacy: emit events against the original request.
	event := collection.EventBurned{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		From:       req.From,
		Amount:     coins,
	}
	ctx.EventManager().EmitEvents(collection.NewEventBurnNFTFrom(event))

	burnt, err := s.keeper.BurnCoins(ctx, req.ContractId, fromAddr, coins)
	if err != nil {
		panic(err)
	}

	// emit events against all burnt tokens.
	event.Amount = burnt
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	return &collection.MsgOperatorBurnNFTResponse{}, nil
}

func (s msgServer) Modify(c context.Context, req *collection.MsgModify) (*collection.MsgModifyResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	collection.UpdateMsgModify(req)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	operator := sdk.MustAccAddressFromBech32(req.Owner)

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
			if tokenIndex != "" {
				if collection.ValidateNFTID(tokenID) == nil {
					event := collection.EventModifiedNFT{
						ContractId: req.ContractId,
						Operator:   operator.String(),
						TokenId:    tokenID,
						Changes:    changes,
					}
					ctx.EventManager().EmitEvents(collection.NewEventModifyTokenOfNFT(event))
					if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
						panic(err)
					}

					return s.keeper.ModifyNFT(ctx, req.ContractId, tokenID, operator, changes)
				}

				event := collection.EventModifiedTokenClass{
					ContractId: req.ContractId,
					Operator:   operator.String(),
					TokenType:  classID,
					Changes:    changes,
					TypeName:   proto.MessageName(&collection.FTClass{}),
				}

				ctx.EventManager().EmitEvents(collection.NewEventModifyTokenOfFTClass(event))
				if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
					panic(err)
				}

				return s.keeper.ModifyTokenClass(ctx, req.ContractId, classID, operator, changes)
			}

			event := collection.EventModifiedTokenClass{
				ContractId: req.ContractId,
				Operator:   operator.String(),
				TokenType:  classID,
				Changes:    changes,
				TypeName:   proto.MessageName(&collection.NFTClass{}),
			}
			ctx.EventManager().EmitEvents(collection.NewEventModifyTokenType(event))
			if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
				panic(err)
			}

			return s.keeper.ModifyTokenClass(ctx, req.ContractId, classID, operator, changes)
		}
		if req.TokenIndex == "" {
			event := collection.EventModifiedContract{
				ContractId: req.ContractId,
				Operator:   operator.String(),
				Changes:    changes,
			}
			ctx.EventManager().EmitEvents(collection.NewEventModifyCollection(event))
			if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
				panic(err)
			}

			return s.keeper.ModifyContract(ctx, req.ContractId, operator, changes)
		}

		panic(sdkerrors.ErrInvalidRequest.Wrap("token index without type"))
	}

	if err := modify(req.TokenType, req.TokenIndex); err != nil {
		return nil, err
	}

	return &collection.MsgModifyResponse{}, nil
}

func (s msgServer) GrantPermission(c context.Context, req *collection.MsgGrantPermission) (*collection.MsgGrantPermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	granter := sdk.MustAccAddressFromBech32(req.From)
	grantee := sdk.MustAccAddressFromBech32(req.To)
	permission := collection.Permission(collection.LegacyPermissionFromString(req.Permission))

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, granter, permission); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrapf("%s is not authorized for %s", granter, permission)
	}

	s.keeper.Grant(ctx, req.ContractId, granter, grantee, permission)

	event := collection.EventGranted{
		ContractId: req.ContractId,
		Granter:    granter.String(),
		Grantee:    grantee.String(),
		Permission: permission,
	}
	ctx.EventManager().EmitEvent(collection.NewEventGrantPermToken(event))
	// it emits typed event inside s.keeper.Grant()

	return &collection.MsgGrantPermissionResponse{}, nil
}

func (s msgServer) RevokePermission(c context.Context, req *collection.MsgRevokePermission) (*collection.MsgRevokePermissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	grantee := sdk.MustAccAddressFromBech32(req.From)
	permission := collection.Permission(collection.LegacyPermissionFromString(req.Permission))

	if _, err := s.keeper.GetGrant(ctx, req.ContractId, grantee, permission); err != nil {
		return nil, collection.ErrTokenNoPermission.Wrapf("%s is not authorized for %s", grantee, permission)
	}

	s.keeper.Abandon(ctx, req.ContractId, grantee, permission)

	event := collection.EventRenounced{
		ContractId: req.ContractId,
		Grantee:    grantee.String(),
		Permission: permission,
	}
	ctx.EventManager().EmitEvent(collection.NewEventRevokePermToken(event))
	// it emits typed event inside s.keeper.Abandon()

	return &collection.MsgRevokePermissionResponse{}, nil
}

func (s msgServer) Attach(c context.Context, req *collection.MsgAttach) (*collection.MsgAttachResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	event := collection.EventAttached{
		ContractId: req.ContractId,
		Operator:   req.From,
		Holder:     req.From,
		Subject:    req.TokenId,
		Target:     req.ToTokenId,
	}
	newRoot := s.keeper.GetRoot(ctx, req.ContractId, req.ToTokenId)
	ctx.EventManager().EmitEvent(collection.NewEventAttachToken(event, newRoot))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)

	if err := s.keeper.Attach(ctx, req.ContractId, fromAddr, req.TokenId, req.ToTokenId); err != nil {
		return nil, err
	}

	return &collection.MsgAttachResponse{}, nil
}

func (s msgServer) Detach(c context.Context, req *collection.MsgDetach) (*collection.MsgDetachResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	// legacy
	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, err
	}
	oldRoot := s.keeper.GetRoot(ctx, req.ContractId, req.TokenId)

	// for the additional field of the event
	parent, err := s.keeper.GetParent(ctx, req.ContractId, req.TokenId)
	if err != nil {
		return nil, collection.ErrTokenNotAChild.Wrap(err.Error())
	}
	event := collection.EventDetached{
		ContractId:     req.ContractId,
		Operator:       req.From,
		Holder:         req.From,
		Subject:        req.TokenId,
		PreviousParent: *parent,
	}
	ctx.EventManager().EmitEvent(collection.NewEventDetachToken(event, oldRoot))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)

	if err := s.keeper.Detach(ctx, req.ContractId, fromAddr, req.TokenId); err != nil {
		return nil, err
	}

	return &collection.MsgDetachResponse{}, nil
}

func (s msgServer) OperatorAttach(c context.Context, req *collection.MsgOperatorAttach) (*collection.MsgOperatorAttachResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)
	operatorAddr := sdk.MustAccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, operatorAddr); err != nil {
		return nil, collection.ErrCollectionNotApproved.Wrap(err.Error())
	}

	event := collection.EventAttached{
		ContractId: req.ContractId,
		Operator:   req.Operator,
		Holder:     req.From,
		Subject:    req.TokenId,
		Target:     req.ToTokenId,
	}
	newRoot := s.keeper.GetRoot(ctx, req.ContractId, req.ToTokenId)
	ctx.EventManager().EmitEvent(collection.NewEventAttachFrom(event, newRoot))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	if err := s.keeper.Attach(ctx, req.ContractId, fromAddr, req.TokenId, req.ToTokenId); err != nil {
		return nil, err
	}

	return &collection.MsgOperatorAttachResponse{}, nil
}

func (s msgServer) OperatorDetach(c context.Context, req *collection.MsgOperatorDetach) (*collection.MsgOperatorDetachResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)
	operatorAddr := sdk.MustAccAddressFromBech32(req.Operator)

	if _, err := s.keeper.GetAuthorization(ctx, req.ContractId, fromAddr, operatorAddr); err != nil {
		return nil, err
	}

	// legacy
	if err := s.keeper.hasNFT(ctx, req.ContractId, req.TokenId); err != nil {
		return nil, err
	}
	oldRoot := s.keeper.GetRoot(ctx, req.ContractId, req.TokenId)

	// for the additional field of the event
	parent, err := s.keeper.GetParent(ctx, req.ContractId, req.TokenId)
	if err != nil {
		return nil, collection.ErrTokenNotAChild.Wrap(err.Error())
	}
	event := collection.EventDetached{
		ContractId:     req.ContractId,
		Operator:       req.Operator,
		Holder:         req.From,
		Subject:        req.TokenId,
		PreviousParent: *parent,
	}
	ctx.EventManager().EmitEvent(collection.NewEventDetachFrom(event, oldRoot))
	if err := ctx.EventManager().EmitTypedEvent(&event); err != nil {
		panic(err)
	}

	if err := s.keeper.Detach(ctx, req.ContractId, fromAddr, req.TokenId); err != nil {
		return nil, err
	}

	return &collection.MsgOperatorDetachResponse{}, nil
}
