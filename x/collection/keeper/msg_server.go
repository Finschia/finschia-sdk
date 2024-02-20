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
	ctx := sdk.UnwrapSDKContext(c)

	if err := ValidateLegacyContract(s.keeper, ctx, req.ContractId); err != nil {
		return nil, err
	}

	fromAddr := sdk.MustAccAddressFromBech32(req.From)

	amount := make([]collection.Coin, len(req.TokenIds))
	for i, id := range req.TokenIds {
		amount[i] = collection.Coin{TokenId: id, Amount: math.OneInt()}

		// legacy
		if err := s.keeper.hasNFT(ctx, req.ContractId, id); err != nil {
			return nil, err
		}
		if !s.keeper.getOwner(ctx, req.ContractId, id).Equals(fromAddr) {
			return nil, collection.ErrTokenNotOwnedBy.Wrapf("%s does not have %s", fromAddr, id)
		}
	}

	toAddr := sdk.MustAccAddressFromBech32(req.To)

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
		amount[i] = collection.Coin{TokenId: id, Amount: math.OneInt()}

		// legacy
		if err := s.keeper.hasNFT(ctx, req.ContractId, id); err != nil {
			return nil, err
		}
		if !s.keeper.getOwner(ctx, req.ContractId, id).Equals(fromAddr) {
			return nil, collection.ErrTokenNotOwnedBy.Wrapf("%s does not have %s", fromAddr, id)
		}
	}

	toAddr := sdk.MustAccAddressFromBech32(req.To)

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
		coins = append(coins, collection.NewCoin(id, math.OneInt()))

		// legacy
		if err := s.keeper.hasNFT(ctx, req.ContractId, id); err != nil {
			return nil, err
		}
		if !s.keeper.getOwner(ctx, req.ContractId, id).Equals(fromAddr) {
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
		coins = append(coins, collection.NewCoin(id, math.OneInt()))

		// legacy
		if err := s.keeper.hasNFT(ctx, req.ContractId, id); err != nil {
			return nil, err
		}
		if !s.keeper.getOwner(ctx, req.ContractId, id).Equals(fromAddr) {
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
			if tokenIndex != "" && collection.ValidateNFTID(tokenID) == nil {
				event := collection.EventModifiedNFT{
					ContractId: req.ContractId,
					Operator:   operator.String(),
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
				Operator:   operator.String(),
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
				Operator:   operator.String(),
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

	// it emits typed event inside s.keeper.Grant()
	s.keeper.Grant(ctx, req.ContractId, granter, grantee, permission)

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

	// it emits typed event inside s.keeper.Abandon()
	s.keeper.Abandon(ctx, req.ContractId, grantee, permission)

	return &collection.MsgRevokePermissionResponse{}, nil
}
