package keeper

import (
	"encoding/json"
	"fmt"

	channeltypes "github.com/line/lfb-sdk/x/ibc/core/04-channel/types"
	"github.com/line/lfb-sdk/x/wasm/internal/types"

	sdk "github.com/line/lfb-sdk/types"
	sdkerrors "github.com/line/lfb-sdk/types/errors"
	distributiontypes "github.com/line/lfb-sdk/x/distribution/types"
	stakingtypes "github.com/line/lfb-sdk/x/staking/types"
	abci "github.com/line/ostracon/abci/types"
	wasmvmtypes "github.com/line/wasmvm/types"
)

type QueryHandler struct {
	Ctx           sdk.Context
	Plugins       QueryPlugins
	Caller        sdk.AccAddress
	GasMultiplier uint64
}

func NewQueryHandler(ctx sdk.Context, plugins QueryPlugins, caller sdk.AccAddress, gasMultiplier uint64) QueryHandler {
	return QueryHandler{
		Ctx:           ctx,
		Plugins:       plugins,
		Caller:        caller,
		GasMultiplier: gasMultiplier,
	}
}

// -- interfaces from baseapp - so we can use the GPRQueryRouter --

// GRPCQueryHandler defines a function type which handles ABCI Query requests
// using gRPC
type GRPCQueryHandler = func(ctx sdk.Context, req abci.RequestQuery) (abci.ResponseQuery, error)

type GRPCQueryRouter interface {
	Route(path string) GRPCQueryHandler
}

// -- end baseapp interfaces --

var _ wasmvmtypes.Querier = QueryHandler{}

func (q QueryHandler) Query(request wasmvmtypes.QueryRequest, gasLimit uint64) ([]byte, error) {
	// set a limit for a subctx
	sdkGas := gasLimit / q.GasMultiplier
	subctx := q.Ctx.WithGasMeter(sdk.NewGasMeter(sdkGas))

	// make sure we charge the higher level context even on panic
	defer func() {
		q.Ctx.GasMeter().ConsumeGas(subctx.GasMeter().GasConsumed(), "contract sub-query")
	}()

	// do the query
	if request.Bank != nil {
		return q.Plugins.Bank(subctx, request.Bank)
	}
	if request.Custom != nil {
		return q.Plugins.Custom(subctx, request.Custom)
	}
	if request.IBC != nil {
		return q.Plugins.IBC(subctx, q.Caller, request.IBC)
	}
	if request.Staking != nil {
		return q.Plugins.Staking(subctx, request.Staking)
	}
	if request.Stargate != nil {
		return q.Plugins.Stargate(subctx, request.Stargate)
	}
	if request.Wasm != nil {
		return q.Plugins.Wasm(subctx, request.Wasm)
	}
	return nil, wasmvmtypes.Unknown{}
}

func (q QueryHandler) GasConsumed() uint64 {
	return q.Ctx.GasMeter().GasConsumed()
}

type CustomQuerier func(ctx sdk.Context, request json.RawMessage) ([]byte, error)

type QueryPlugins struct {
	Bank     func(ctx sdk.Context, request *wasmvmtypes.BankQuery) ([]byte, error)
	Custom   CustomQuerier
	IBC      func(ctx sdk.Context, caller sdk.AccAddress, request *wasmvmtypes.IBCQuery) ([]byte, error)
	Staking  func(ctx sdk.Context, request *wasmvmtypes.StakingQuery) ([]byte, error)
	Stargate func(ctx sdk.Context, request *wasmvmtypes.StargateQuery) ([]byte, error)
	Wasm     func(ctx sdk.Context, request *wasmvmtypes.WasmQuery) ([]byte, error)
}

func DefaultQueryPlugins(bank types.BankViewKeeper, staking types.StakingKeeper, distKeeper types.DistributionKeeper, channelKeeper types.ChannelKeeper, queryRouter GRPCQueryRouter, wasm *Keeper) QueryPlugins {
	return QueryPlugins{
		Bank:     BankQuerier(bank),
		Custom:   CustomQuerierImpl(queryRouter),
		IBC:      IBCQuerier(wasm, channelKeeper),
		Staking:  StakingQuerier(staking, distKeeper),
		Stargate: StargateQuerier(queryRouter),
		Wasm:     WasmQuerier(wasm),
	}
}

func (e QueryPlugins) Merge(o *QueryPlugins) QueryPlugins {
	// only update if this is non-nil and then only set values
	if o == nil {
		return e
	}
	if o.Bank != nil {
		e.Bank = o.Bank
	}
	if o.Custom != nil {
		e.Custom = o.Custom
	}
	if o.IBC != nil {
		e.IBC = o.IBC
	}
	if o.Staking != nil {
		e.Staking = o.Staking
	}
	if o.Stargate != nil {
		e.Stargate = o.Stargate
	}
	if o.Wasm != nil {
		e.Wasm = o.Wasm
	}
	return e
}

func BankQuerier(bankKeeper types.BankViewKeeper) func(ctx sdk.Context, request *wasmvmtypes.BankQuery) ([]byte, error) {
	return func(ctx sdk.Context, request *wasmvmtypes.BankQuery) ([]byte, error) {
		if request.AllBalances != nil {
			addr, err := sdk.AccAddressFromBech32(request.AllBalances.Address)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.AllBalances.Address)
			}
			coins := bankKeeper.GetAllBalances(ctx, addr)
			res := wasmvmtypes.AllBalancesResponse{
				Amount: convertSdkCoinsToWasmCoins(coins),
			}
			return json.Marshal(res)
		}
		if request.Balance != nil {
			addr, err := sdk.AccAddressFromBech32(request.Balance.Address)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Balance.Address)
			}
			coins := bankKeeper.GetAllBalances(ctx, addr)
			amount := coins.AmountOf(request.Balance.Denom)
			res := wasmvmtypes.BalanceResponse{
				Amount: wasmvmtypes.Coin{
					Denom:  request.Balance.Denom,
					Amount: amount.String(),
				},
			}
			return json.Marshal(res)
		}
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown BankQuery variant"}
	}
}
func CustomQuerierImpl(queryRouter GRPCQueryRouter) func(ctx sdk.Context, querierJson json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, querierJson json.RawMessage) ([]byte, error) {
		var linkQueryWrapper types.LinkQueryWrapper
		err := json.Unmarshal(querierJson, &linkQueryWrapper)
		if err != nil {
			return nil, err
		}
		route := queryRouter.Route(linkQueryWrapper.Path)
		if route == nil {
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "Unknown encode module"}
		}
		req := abci.RequestQuery{
			Data: linkQueryWrapper.Data,
			Path: linkQueryWrapper.Path,
		}
		res, err := route(ctx, req)
		if err != nil {
			return nil, err
		}

		return res.Value, nil
	}
}

func IBCQuerier(wasm *Keeper, channelKeeper types.ChannelKeeper) func(ctx sdk.Context, caller sdk.AccAddress, request *wasmvmtypes.IBCQuery) ([]byte, error) {
	return func(ctx sdk.Context, caller sdk.AccAddress, request *wasmvmtypes.IBCQuery) ([]byte, error) {
		if request.PortID != nil {
			contractInfo := wasm.GetContractInfo(ctx, caller)
			res := wasmvmtypes.PortIDResponse{
				PortID: contractInfo.IBCPortID,
			}
			return json.Marshal(res)
		}
		if request.ListChannels != nil {
			portID := request.ListChannels.PortID
			channels := make(wasmvmtypes.IBCChannels, 0)
			channelKeeper.IterateChannels(ctx, func(ch channeltypes.IdentifiedChannel) bool {
				if portID == "" || portID == ch.PortId {
					newChan := wasmvmtypes.IBCChannel{
						Endpoint: wasmvmtypes.IBCEndpoint{
							PortID:    ch.PortId,
							ChannelID: ch.ChannelId,
						},
						CounterpartyEndpoint: wasmvmtypes.IBCEndpoint{
							PortID:    ch.Counterparty.PortId,
							ChannelID: ch.Counterparty.ChannelId,
						},
						Order:        ch.Ordering.String(),
						Version:      ch.Version,
						ConnectionID: ch.ConnectionHops[0],
					}
					channels = append(channels, newChan)
				}
				return false
			})
			res := wasmvmtypes.ListChannelsResponse{
				Channels: channels,
			}
			return json.Marshal(res)
		}
		if request.Channel != nil {
			channelID := request.Channel.ChannelID
			portID := request.Channel.PortID
			if portID == "" {
				contractInfo := wasm.GetContractInfo(ctx, caller)
				portID = contractInfo.IBCPortID
			}
			got, found := channelKeeper.GetChannel(ctx, portID, channelID)
			var channel *wasmvmtypes.IBCChannel
			if found {
				channel = &wasmvmtypes.IBCChannel{
					Endpoint: wasmvmtypes.IBCEndpoint{
						PortID:    portID,
						ChannelID: channelID,
					},
					CounterpartyEndpoint: wasmvmtypes.IBCEndpoint{
						PortID:    got.Counterparty.PortId,
						ChannelID: got.Counterparty.ChannelId,
					},
					Order:               got.Ordering.String(),
					Version:             got.Version,
					CounterpartyVersion: "",
					ConnectionID:        got.ConnectionHops[0],
				}
			}
			res := wasmvmtypes.ChannelResponse{
				Channel: channel,
			}
			return json.Marshal(res)
		}
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown IBCQuery variant"}
	}
}

func StargateQuerier(queryRouter GRPCQueryRouter) func(ctx sdk.Context, request *wasmvmtypes.StargateQuery) ([]byte, error) {
	return func(ctx sdk.Context, msg *wasmvmtypes.StargateQuery) ([]byte, error) {
		route := queryRouter.Route(msg.Path)
		if route == nil {
			return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("No route to query '%s'", msg.Path)}
		}
		req := abci.RequestQuery{
			Data: msg.Data,
			Path: msg.Path,
		}
		res, err := route(ctx, req)
		if err != nil {
			return nil, err
		}
		return res.Value, nil
	}
}

func StakingQuerier(keeper types.StakingKeeper, distKeeper types.DistributionKeeper) func(ctx sdk.Context, request *wasmvmtypes.StakingQuery) ([]byte, error) {
	return func(ctx sdk.Context, request *wasmvmtypes.StakingQuery) ([]byte, error) {
		if request.BondedDenom != nil {
			denom := keeper.BondDenom(ctx)
			res := wasmvmtypes.BondedDenomResponse{
				Denom: denom,
			}
			return json.Marshal(res)
		}
		if request.Validators != nil {
			validators := keeper.GetBondedValidatorsByPower(ctx)
			//validators := keeper.GetAllValidators(ctx)
			wasmVals := make([]wasmvmtypes.Validator, len(validators))
			for i, v := range validators {
				wasmVals[i] = wasmvmtypes.Validator{
					Address:       v.OperatorAddress,
					Commission:    v.Commission.Rate.String(),
					MaxCommission: v.Commission.MaxRate.String(),
					MaxChangeRate: v.Commission.MaxChangeRate.String(),
				}
			}
			res := wasmvmtypes.ValidatorsResponse{
				Validators: wasmVals,
			}
			return json.Marshal(res)
		}
		if request.AllDelegations != nil {
			delegator, err := sdk.AccAddressFromBech32(request.AllDelegations.Delegator)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.AllDelegations.Delegator)
			}
			sdkDels := keeper.GetAllDelegatorDelegations(ctx, delegator)
			delegations, err := sdkToDelegations(ctx, keeper, sdkDels)
			if err != nil {
				return nil, err
			}
			res := wasmvmtypes.AllDelegationsResponse{
				Delegations: delegations,
			}
			return json.Marshal(res)
		}
		if request.Delegation != nil {
			delegator, err := sdk.AccAddressFromBech32(request.Delegation.Delegator)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Delegation.Delegator)
			}
			validator, err := sdk.ValAddressFromBech32(request.Delegation.Validator)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Delegation.Validator)
			}

			var res wasmvmtypes.DelegationResponse
			d, found := keeper.GetDelegation(ctx, delegator, validator)
			if found {
				res.Delegation, err = sdkToFullDelegation(ctx, keeper, distKeeper, d)
				if err != nil {
					return nil, err
				}
			}
			return json.Marshal(res)
		}
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Staking variant"}
	}
}

func sdkToDelegations(ctx sdk.Context, keeper types.StakingKeeper, delegations []stakingtypes.Delegation) (wasmvmtypes.Delegations, error) {
	result := make([]wasmvmtypes.Delegation, len(delegations))
	bondDenom := keeper.BondDenom(ctx)

	for i, d := range delegations {
		delAddr, err := sdk.AccAddressFromBech32(d.DelegatorAddress)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "delegator address")
		}
		valAddr, err := sdk.ValAddressFromBech32(d.ValidatorAddress)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "validator address")
		}

		// shares to amount logic comes from here:
		// x/staking/keeper/querier.go DelegationToDelegationResponse
		/// https://github.com/line/lfb-sdk/blob/3ccf3913f53e2a9ccb4be8429bee32e67669e89a/x/staking/keeper/querier.go#L450
		val, found := keeper.GetValidator(ctx, valAddr)
		if !found {
			return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, "can't load validator for delegation")
		}
		amount := sdk.NewCoin(bondDenom, val.TokensFromShares(d.Shares).TruncateInt())

		result[i] = wasmvmtypes.Delegation{
			Delegator: delAddr.String(),
			Validator: valAddr.String(),
			Amount:    convertSdkCoinToWasmCoin(amount),
		}
	}
	return result, nil
}

func sdkToFullDelegation(ctx sdk.Context, keeper types.StakingKeeper, distKeeper types.DistributionKeeper, delegation stakingtypes.Delegation) (*wasmvmtypes.FullDelegation, error) {
	delAddr, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "delegator address")
	}
	valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "validator address")
	}
	val, found := keeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, "can't load validator for delegation")
	}
	bondDenom := keeper.BondDenom(ctx)
	amount := sdk.NewCoin(bondDenom, val.TokensFromShares(delegation.Shares).TruncateInt())

	delegationCoins := convertSdkCoinToWasmCoin(amount)

	// FIXME: this is very rough but better than nothing...
	// https://github.com/line/lfb-sdk/issues/225
	// if this (val, delegate) pair is receiving a redelegation, it cannot redelegate more
	// otherwise, it can redelegate the full amount
	// (there are cases of partial funds redelegated, but this is a start)
	redelegateCoins := wasmvmtypes.NewCoin(0, bondDenom)
	if !keeper.HasReceivingRedelegation(ctx, delAddr, valAddr) {
		redelegateCoins = delegationCoins
	}

	// FIXME: make a cleaner way to do this (modify the sdk)
	// we need the info from `distKeeper.calculateDelegationRewards()`, but it is not public
	// neither is `queryDelegationRewards(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper)`
	// so we go through the front door of the querier....
	accRewards, err := getAccumulatedRewards(ctx, distKeeper, delegation)
	if err != nil {
		return nil, err
	}

	return &wasmvmtypes.FullDelegation{
		Delegator:          delAddr.String(),
		Validator:          valAddr.String(),
		Amount:             delegationCoins,
		AccumulatedRewards: accRewards,
		CanRedelegate:      redelegateCoins,
	}, nil
}

// FIXME: simplify this enormously when
// https://github.com/cosmos/cosmos-sdk/issues/7466 is merged
func getAccumulatedRewards(ctx sdk.Context, distKeeper types.DistributionKeeper, delegation stakingtypes.Delegation) ([]wasmvmtypes.Coin, error) {
	// Try to get *delegator* reward info!
	params := distributiontypes.QueryDelegationRewardsRequest{
		DelegatorAddress: delegation.DelegatorAddress,
		ValidatorAddress: delegation.ValidatorAddress,
	}
	cache, _ := ctx.CacheContext()
	qres, err := distKeeper.DelegationRewards(sdk.WrapSDKContext(cache), &params)
	if err != nil {
		return nil, err
	}

	// now we have it, convert it into wasmvm types
	rewards := make([]wasmvmtypes.Coin, len(qres.Rewards))
	for i, r := range qres.Rewards {
		rewards[i] = wasmvmtypes.Coin{
			Denom:  r.Denom,
			Amount: r.Amount.TruncateInt().String(),
		}
	}
	return rewards, nil
}

func WasmQuerier(wasm *Keeper) func(ctx sdk.Context, request *wasmvmtypes.WasmQuery) ([]byte, error) {
	return func(ctx sdk.Context, request *wasmvmtypes.WasmQuery) ([]byte, error) {
		if request.Smart != nil {
			addr, err := sdk.AccAddressFromBech32(request.Smart.ContractAddr)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Smart.ContractAddr)
			}
			return wasm.QuerySmart(ctx, addr, request.Smart.Msg)
		}
		if request.Raw != nil {
			addr, err := sdk.AccAddressFromBech32(request.Raw.ContractAddr)
			if err != nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Raw.ContractAddr)
			}
			return wasm.QueryRaw(ctx, addr, request.Raw.Key), nil
		}
		return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown WasmQuery variant"}
	}
}

func convertSdkCoinsToWasmCoins(coins []sdk.Coin) wasmvmtypes.Coins {
	converted := make(wasmvmtypes.Coins, len(coins))
	for i, c := range coins {
		converted[i] = convertSdkCoinToWasmCoin(c)
	}
	return converted
}

func convertSdkCoinToWasmCoin(coin sdk.Coin) wasmvmtypes.Coin {
	return wasmvmtypes.Coin{
		Denom:  coin.Denom,
		Amount: coin.Amount.String(),
	}
}
