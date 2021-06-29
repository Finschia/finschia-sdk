package keeper

import (
	"encoding/json"
	"errors"
	"fmt"

	codectypes "github.com/line/lfb-sdk/codec/types"
	sdk "github.com/line/lfb-sdk/types"
	sdkerrors "github.com/line/lfb-sdk/types/errors"
	banktypes "github.com/line/lfb-sdk/x/bank/types"
	distributiontypes "github.com/line/lfb-sdk/x/distribution/types"
	ibctransfertypes "github.com/line/lfb-sdk/x/ibc/applications/transfer/types"
	ibcclienttypes "github.com/line/lfb-sdk/x/ibc/core/02-client/types"
	channeltypes "github.com/line/lfb-sdk/x/ibc/core/04-channel/types"
	host "github.com/line/lfb-sdk/x/ibc/core/24-host"
	stakingtypes "github.com/line/lfb-sdk/x/staking/types"
	"github.com/line/lfb-sdk/x/wasm/types"
	wasmvmtypes "github.com/line/wasmvm/types"
)

// msgEncoder is an extension point to customize encodings
type msgEncoder interface {
	// Encode converts wasmvm message to n cosmos message types
	Encode(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Msg, error)
}

// SDKMessageHandler can handles messages that can be encoded into sdk.Message types and routed.
type SDKMessageHandler struct {
	router       sdk.Router
	encodeRouter types.Router
	encoders     MessageEncoders
}

type BankEncoder func(sender sdk.AccAddress, msg *wasmvmtypes.BankMsg) ([]sdk.Msg, error)
type CustomEncoder func(sender sdk.AccAddress, msg json.RawMessage, customEncodeRouter types.Router) ([]sdk.Msg, error)
type DistributionEncoder func(sender sdk.AccAddress, msg *wasmvmtypes.DistributionMsg) ([]sdk.Msg, error)
type StakingEncoder func(sender sdk.AccAddress, msg *wasmvmtypes.StakingMsg) ([]sdk.Msg, error)
type StargateEncoder func(sender sdk.AccAddress, msg *wasmvmtypes.StargateMsg) ([]sdk.Msg, error)
type WasmEncoder func(sender sdk.AccAddress, msg *wasmvmtypes.WasmMsg) ([]sdk.Msg, error)
type IBCEncoder func(ctx sdk.Context, sender sdk.AccAddress, contractIBCPortID string, msg *wasmvmtypes.IBCMsg) ([]sdk.Msg, error)

type MessageEncoders struct {
	Bank         func(sender sdk.AccAddress, msg *wasmvmtypes.BankMsg) ([]sdk.Msg, error)
	Custom       func(sender sdk.AccAddress, msg json.RawMessage) ([]sdk.Msg, error)
	Distribution func(sender sdk.AccAddress, msg *wasmvmtypes.DistributionMsg) ([]sdk.Msg, error)
	IBC          func(ctx sdk.Context, sender sdk.AccAddress, contractIBCPortID string, msg *wasmvmtypes.IBCMsg) ([]sdk.Msg, error)
	Staking      func(sender sdk.AccAddress, msg *wasmvmtypes.StakingMsg) ([]sdk.Msg, error)
	Stargate     func(sender sdk.AccAddress, msg *wasmvmtypes.StargateMsg) ([]sdk.Msg, error)
	Wasm         func(sender sdk.AccAddress, msg *wasmvmtypes.WasmMsg) ([]sdk.Msg, error)
}

func DefaultEncoders(unpacker codectypes.AnyUnpacker, portSource types.ICS20TransferPortSource) MessageEncoders {
	return MessageEncoders{
		Bank:         EncodeBankMsg,
		Custom:       CustomMsg,
		Distribution: EncodeDistributionMsg,
		IBC:          EncodeIBCMsg(portSource),
		Staking:      EncodeStakingMsg,
		Stargate:     EncodeStargateMsg(unpacker),
		Wasm:         EncodeWasmMsg,
	}
}

func NewDefaultMessageHandler(
	router sdk.Router,
	encodeRouter types.Router,
	channelKeeper types.ChannelKeeper,
	capabilityKeeper types.CapabilityKeeper,
	bankKeeper types.Burner,
	unpacker codectypes.AnyUnpacker,
	portSource types.ICS20TransferPortSource,
	customEncoders ...*MessageEncoders,
) Messenger {
	encoders := DefaultEncoders(unpacker, portSource)
	for _, e := range customEncoders {
		encoders = encoders.Merge(e)
	}
	return NewMessageHandlerChain(
		NewSDKMessageHandler(router, encodeRouter, encoders),
		NewIBCRawPacketHandler(channelKeeper, capabilityKeeper),
		NewBurnCoinMessageHandler(bankKeeper),
	)
}

func NewSDKMessageHandler(router sdk.Router, encodeRouter types.Router, encoders msgEncoder) SDKMessageHandler {
	return SDKMessageHandler{
		router:       router,
		encodeRouter: encodeRouter,
		encoders:     encoders,
	}
}

func (e MessageEncoders) Merge(o *MessageEncoders) MessageEncoders {
	if o == nil {
		return e
	}
	if o.Bank != nil {
		e.Bank = o.Bank
	}
	if o.Custom != nil {
		e.Custom = o.Custom
	}
	if o.Distribution != nil {
		e.Distribution = o.Distribution
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

func (e MessageEncoders) Encode(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg, customRouter types.Router) ([]sdk.Msg, error) {
	switch {
	case msg.Bank != nil:
		return e.Bank(contractAddr, msg.Bank)
	case msg.Custom != nil:
		return e.Custom(contractAddr, msg.Custom, customRouter)
	case msg.Distribution != nil:
		return e.Distribution(contractAddr, msg.Distribution)
	case msg.IBC != nil:
		return e.IBC(ctx, contractAddr, contractIBCPortID, msg.IBC)
	case msg.Staking != nil:
		return e.Staking(contractAddr, msg.Staking)
	case msg.Stargate != nil:
		return e.Stargate(contractAddr, msg.Stargate)
	case msg.Wasm != nil:
		return e.Wasm(contractAddr, msg.Wasm)
	}
	return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "unknown variant of Wasm")
}

func EncodeBankMsg(sender sdk.AccAddress, msg *wasmvmtypes.BankMsg) ([]sdk.Msg, error) {
	if msg.Send == nil {
		return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "unknown variant of Bank")
	}
	if len(msg.Send.Amount) == 0 {
		return nil, nil
	}
	toSend, err := convertWasmCoinsToSdkCoins(msg.Send.Amount)
	if err != nil {
		return nil, err
	}
	sdkMsg := banktypes.MsgSend{
		FromAddress: sender.String(),
		ToAddress:   msg.Send.ToAddress,
		Amount:      toSend,
	}
	return []sdk.Msg{&sdkMsg}, nil
}

func CustomMsg(sender sdk.AccAddress, jsonMsg json.RawMessage, router types.Router) ([]sdk.Msg, error) {
	var linkMsgWrapper types.LinkMsgWrapper
	err := json.Unmarshal(jsonMsg, &linkMsgWrapper)
	if err != nil {
		return nil, err
	}
	handler := router.GetRoute(linkMsgWrapper.Module)
	if handler == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized handler: %T", linkMsgWrapper.Module)
	}
	return handler(linkMsgWrapper.MsgData)
}

func EncodeDistributionMsg(sender sdk.AccAddress, msg *wasmvmtypes.DistributionMsg) ([]sdk.Msg, error) {
	switch {
	case msg.SetWithdrawAddress != nil:
		setMsg := distributiontypes.MsgSetWithdrawAddress{
			DelegatorAddress: sender.String(),
			WithdrawAddress:  msg.SetWithdrawAddress.Address,
		}
		return []sdk.Msg{&setMsg}, nil
	case msg.WithdrawDelegatorReward != nil:
		withdrawMsg := distributiontypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: sender.String(),
			ValidatorAddress: msg.WithdrawDelegatorReward.Validator,
		}
		return []sdk.Msg{&withdrawMsg}, nil
	default:
		return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "unknown variant of Distribution")
	}
}

func EncodeStakingMsg(sender sdk.AccAddress, msg *wasmvmtypes.StakingMsg) ([]sdk.Msg, error) {
	switch {
	case msg.Delegate != nil:
		coin, err := convertWasmCoinToSdkCoin(msg.Delegate.Amount)
		if err != nil {
			return nil, err
		}
		sdkMsg := stakingtypes.MsgDelegate{
			DelegatorAddress: sender.String(),
			ValidatorAddress: msg.Delegate.Validator,
			Amount:           coin,
		}
		return []sdk.Msg{&sdkMsg}, nil

	case msg.Redelegate != nil:
		coin, err := convertWasmCoinToSdkCoin(msg.Redelegate.Amount)
		if err != nil {
			return nil, err
		}
		sdkMsg := stakingtypes.MsgBeginRedelegate{
			DelegatorAddress:    sender.String(),
			ValidatorSrcAddress: msg.Redelegate.SrcValidator,
			ValidatorDstAddress: msg.Redelegate.DstValidator,
			Amount:              coin,
		}
		return []sdk.Msg{&sdkMsg}, nil
	case msg.Undelegate != nil:
		coin, err := convertWasmCoinToSdkCoin(msg.Undelegate.Amount)
		if err != nil {
			return nil, err
		}
		sdkMsg := stakingtypes.MsgUndelegate{
			DelegatorAddress: sender.String(),
			ValidatorAddress: msg.Undelegate.Validator,
			Amount:           coin,
		}
		return []sdk.Msg{&sdkMsg}, nil
	default:
		return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "unknown variant of Staking")
	}
}

func EncodeStargateMsg(unpacker codectypes.AnyUnpacker) StargateEncoder {
	return func(sender sdk.AccAddress, msg *wasmvmtypes.StargateMsg) ([]sdk.Msg, error) {
		any := codectypes.Any{
			TypeUrl: msg.TypeURL,
			Value:   msg.Value,
		}
		var sdkMsg sdk.Msg
		if err := unpacker.UnpackAny(&any, &sdkMsg); err != nil {
			return nil, sdkerrors.Wrap(types.ErrInvalidMsg, fmt.Sprintf("Cannot unpack proto message with type URL: %s", msg.TypeURL))
		}
		if err := codectypes.UnpackInterfaces(sdkMsg, unpacker); err != nil {
			return nil, sdkerrors.Wrap(types.ErrInvalidMsg, fmt.Sprintf("UnpackInterfaces inside msg: %s", err))
		}
		return []sdk.Msg{sdkMsg}, nil
	}
}

func EncodeWasmMsg(sender sdk.AccAddress, msg *wasmvmtypes.WasmMsg) ([]sdk.Msg, error) {
	switch {
	case msg.Execute != nil:
		coins, err := convertWasmCoinsToSdkCoins(msg.Execute.Send)
		if err != nil {
			return nil, err
		}

		sdkMsg := types.MsgExecuteContract{
			Sender:   sender.String(),
			Contract: msg.Execute.ContractAddr,
			Msg:      msg.Execute.Msg,
			Funds:    coins,
		}
		return []sdk.Msg{&sdkMsg}, nil
	case msg.Instantiate != nil:
		coins, err := convertWasmCoinsToSdkCoins(msg.Instantiate.Send)
		if err != nil {
			return nil, err
		}

		sdkMsg := types.MsgInstantiateContract{
			Sender:  sender.String(),
			CodeID:  msg.Instantiate.CodeID,
			Label:   msg.Instantiate.Label,
			InitMsg: msg.Instantiate.Msg,
			Admin:   msg.Instantiate.Admin,
			Funds:   coins,
		}
		return []sdk.Msg{&sdkMsg}, nil
	case msg.Migrate != nil:
		sdkMsg := types.MsgMigrateContract{
			Sender:     sender.String(),
			Contract:   msg.Migrate.ContractAddr,
			CodeID:     msg.Migrate.NewCodeID,
			MigrateMsg: msg.Migrate.Msg,
		}
		return []sdk.Msg{&sdkMsg}, nil
	case msg.ClearAdmin != nil:
		sdkMsg := types.MsgClearAdmin{
			Sender:   sender.String(),
			Contract: msg.ClearAdmin.ContractAddr,
		}
		return []sdk.Msg{&sdkMsg}, nil
	case msg.UpdateAdmin != nil:
		sdkMsg := types.MsgUpdateAdmin{
			Sender:   sender.String(),
			Contract: msg.UpdateAdmin.ContractAddr,
			NewAdmin: msg.UpdateAdmin.Admin,
		}
		return []sdk.Msg{&sdkMsg}, nil
	default:
		return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "unknown variant of Wasm")
	}
}

func EncodeIBCMsg(portSource types.ICS20TransferPortSource) func(ctx sdk.Context, sender sdk.AccAddress, contractIBCPortID string, msg *wasmvmtypes.IBCMsg) ([]sdk.Msg, error) {
	return func(ctx sdk.Context, sender sdk.AccAddress, contractIBCPortID string, msg *wasmvmtypes.IBCMsg) ([]sdk.Msg, error) {
		switch {
		case msg.CloseChannel != nil:
			return []sdk.Msg{&channeltypes.MsgChannelCloseInit{
				PortId:    PortIDForContract(sender),
				ChannelId: msg.CloseChannel.ChannelID,
				Signer:    sender.String(),
			}}, nil
		case msg.Transfer != nil:
			amount, err := convertWasmCoinToSdkCoin(msg.Transfer.Amount)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "amount")
			}
			msg := &ibctransfertypes.MsgTransfer{
				SourcePort:       portSource.GetPort(ctx),
				SourceChannel:    msg.Transfer.ChannelID,
				Token:            amount,
				Sender:           sender.String(),
				Receiver:         msg.Transfer.ToAddress,
				TimeoutHeight:    convertWasmIBCTimeoutHeightToCosmosHeight(msg.Transfer.Timeout.Block),
				TimeoutTimestamp: msg.Transfer.Timeout.Timestamp,
			}
			return []sdk.Msg{msg}, nil
		default:
			return nil, sdkerrors.Wrap(types.ErrUnknownMsg, "Unknown variant of IBC")
		}
	}
}

func convertWasmIBCTimeoutHeightToCosmosHeight(ibcTimeoutBlock *wasmvmtypes.IBCTimeoutBlock) ibcclienttypes.Height {
	if ibcTimeoutBlock == nil {
		return ibcclienttypes.NewHeight(0, 0)
	}
	return ibcclienttypes.NewHeight(ibcTimeoutBlock.Revision, ibcTimeoutBlock.Height)
}

func convertWasmCoinsToSdkCoins(coins []wasmvmtypes.Coin) (sdk.Coins, error) {
	var toSend sdk.Coins
	for _, coin := range coins {
		c, err := convertWasmCoinToSdkCoin(coin)
		if err != nil {
			return nil, err
		}
		toSend = append(toSend, c)
	}
	return toSend, nil
}

func convertWasmCoinToSdkCoin(coin wasmvmtypes.Coin) (sdk.Coin, error) {
	amount, ok := sdk.NewIntFromString(coin.Amount)
	if !ok {
		return sdk.Coin{}, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, coin.Amount+coin.Denom)
	}
	return sdk.Coin{
		Denom:  coin.Denom,
		Amount: amount,
	}, nil
}

func (h SDKMessageHandler) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
	sdkMsgs, err := h.encoders.Encode(ctx, contractAddr, contractIBCPortID, msg, h.encodeRouter)
	if err != nil {
		return nil, nil, err
	}
	for _, sdkMsg := range sdkMsgs {
		res, err := h.handleSdkMessage(ctx, contractAddr, sdkMsg)
		if err != nil {
			return nil, nil, err
		}
		// append data
		data = append(data, res.Data)
		// append events
		sdkEvents := make([]sdk.Event, len(res.Events))
		for i := range res.Events {
			sdkEvents[i] = sdk.Event(res.Events[i])
		}
		events = append(events, sdkEvents...)
	}
	return
}

func (h SDKMessageHandler) handleSdkMessage(ctx sdk.Context, contractAddr sdk.Address, msg sdk.Msg) (*sdk.Result, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	// make sure this account can send it
	for _, acct := range msg.GetSigners() {
		if !acct.Equals(contractAddr) {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "contract doesn't have permission")
		}
	}

	// find the handler and execute it
	handler := h.router.Route(ctx, msg.Route())
	if handler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, msg.Route())
	}
	res, err := handler(ctx, msg)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// MessageHandlerChain defines a chain of handlers that are called one by one until it can be handled.
type MessageHandlerChain struct {
	handlers []Messenger
}

func NewMessageHandlerChain(first Messenger, others ...Messenger) *MessageHandlerChain {
	r := &MessageHandlerChain{handlers: append([]Messenger{first}, others...)}
	for i := range r.handlers {
		if r.handlers[i] == nil {
			panic(fmt.Sprintf("handler must not be nil at position : %d", i))
		}
	}
	return r
}

// DispatchMsg dispatch message to handlers.
func (m MessageHandlerChain) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	for _, h := range m.handlers {
		events, data, err := h.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
		switch {
		case err == nil:
			return events, data, err
		case errors.Is(err, types.ErrUnknownMsg):
			continue
		default:
			return events, data, err
		}
	}
	return nil, nil, sdkerrors.Wrap(types.ErrUnknownMsg, "no handler found")
}

// IBCRawPacketHandler handels IBC.SendPacket messages which are published to an IBC channel.
type IBCRawPacketHandler struct {
	channelKeeper    types.ChannelKeeper
	capabilityKeeper types.CapabilityKeeper
}

func NewIBCRawPacketHandler(chk types.ChannelKeeper, cak types.CapabilityKeeper) *IBCRawPacketHandler {
	return &IBCRawPacketHandler{channelKeeper: chk, capabilityKeeper: cak}
}

// DispatchMsg publishes a raw IBC packet onto the channel.
func (h IBCRawPacketHandler) DispatchMsg(ctx sdk.Context, _ sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
	if msg.IBC == nil || msg.IBC.SendPacket == nil {
		return nil, nil, types.ErrUnknownMsg
	}
	if contractIBCPortID == "" {
		return nil, nil, sdkerrors.Wrapf(types.ErrUnsupportedForContract, "ibc not supported")
	}
	contractIBCChannelID := msg.IBC.SendPacket.ChannelID
	if contractIBCChannelID == "" {
		return nil, nil, sdkerrors.Wrapf(types.ErrEmpty, "ibc channel")
	}

	sequence, found := h.channelKeeper.GetNextSequenceSend(ctx, contractIBCPortID, contractIBCChannelID)
	if !found {
		return nil, nil, sdkerrors.Wrapf(channeltypes.ErrSequenceSendNotFound,
			"source port: %s, source channel: %s", contractIBCPortID, contractIBCChannelID,
		)
	}

	channelInfo, ok := h.channelKeeper.GetChannel(ctx, contractIBCPortID, contractIBCChannelID)
	if !ok {
		return nil, nil, sdkerrors.Wrap(channeltypes.ErrInvalidChannel, "not found")
	}
	channelCap, ok := h.capabilityKeeper.GetCapability(ctx, host.ChannelCapabilityPath(contractIBCPortID, contractIBCChannelID))
	if !ok {
		return nil, nil, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}
	packet := channeltypes.NewPacket(
		msg.IBC.SendPacket.Data,
		sequence,
		contractIBCPortID,
		contractIBCChannelID,
		channelInfo.Counterparty.PortId,
		channelInfo.Counterparty.ChannelId,
		convertWasmIBCTimeoutHeightToCosmosHeight(msg.IBC.SendPacket.Timeout.Block),
		msg.IBC.SendPacket.Timeout.Timestamp,
	)
	return nil, nil, h.channelKeeper.SendPacket(ctx, channelCap, packet)
}

var _ Messenger = MessageHandlerFunc(nil)

// MessageHandlerFunc is a helper to construct simple function based message handler
type MessageHandlerFunc func(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error)

func (m MessageHandlerFunc) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
	return m(ctx, contractAddr, contractIBCPortID, msg)
}

// NewBurnCoinMessageHandler handles wasmvm.BurnMsg messages
func NewBurnCoinMessageHandler(burner types.Burner) MessageHandlerFunc {
	return func(ctx sdk.Context, contractAddr sdk.AccAddress, _ string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
		if msg.Bank != nil && msg.Bank.Burn != nil {
			coins, err := convertWasmCoinsToSdkCoins(msg.Bank.Burn.Amount)
			if err != nil {
				return nil, nil, err
			}
			if err := burner.SendCoinsFromAccountToModule(ctx, contractAddr, types.ModuleName, coins); err != nil {
				return nil, nil, sdkerrors.Wrap(err, "transfer to module")
			}
			if err := burner.BurnCoins(ctx, types.ModuleName, coins); err != nil {
				return nil, nil, sdkerrors.Wrap(err, "burn coins")
			}
			moduleLogger(ctx).Info("Burned", "amount", coins)
			return nil, nil, nil
		}
		return nil, nil, types.ErrUnknownMsg
	}
}
