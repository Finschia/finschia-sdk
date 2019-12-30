package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/line/link/client"
	"github.com/line/link/x/safetybox/internal/types"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Safety Box transaction subcommands",
		Aliases:                    []string{"sb"},
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		SafetyBoxCreateTxCmd(cdc),
		SafetyBoxRoleTxCmd(cdc),
		SafetyBoxSendCoinsTxCmd(cdc),
	)
	return txCmd
}

func SafetyBoxCreateTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [id] [owner_address] [denom]",
		Short: "Create a safety box with ID, owner and the coin denom. Only one owner and denom is allowed.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithFrom(args[1]).WithCodec(cdc)
			safetyBoxId := args[0]
			safetyBoxOwner := cliCtx.FromAddress
			safetyBoxDenoms := strings.Split(args[2], ",")

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.MsgSafetyBoxCreate{SafetyBoxId: safetyBoxId, SafetyBoxOwner: safetyBoxOwner, SafetyBoxDenoms: safetyBoxDenoms}
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func SafetyBoxRoleTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role [safety_box_id] [register|deregister] [operator|allocator|issuer|returner] [from_address] [to_address]",
		Short: "Register or deregister roles to the address on the safety box",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			safetyBoxId, action, role := args[0], args[1], args[2]
			cliCtx := client.NewCLIContextWithFrom(args[3]).WithCodec(cdc)
			fromAddress, err := sdk.AccAddressFromBech32(args[3])
			if err != nil {
				return err
			}
			toAddress, err := sdk.AccAddressFromBech32(args[4])
			if err != nil {
				return err
			}

			var msg sdk.Msg
			switch role {
			case types.RoleOperator:
				if action == types.RegisterRole {
					msg = types.MsgSafetyBoxRegisterOperator{
						SafetyBoxId:    safetyBoxId,
						SafetyBoxOwner: fromAddress,
						Address:        toAddress,
					}
				} else if action == types.DeregisterRole {
					msg = types.MsgSafetyBoxDeregisterOperator{
						SafetyBoxId:    safetyBoxId,
						SafetyBoxOwner: fromAddress,
						Address:        toAddress,
					}
				} else {
					return types.ErrSafetyBoxInvalidAction(types.DefaultCodespace)
				}
			case types.RoleAllocator:
				if action == types.RegisterRole {
					msg = types.MsgSafetyBoxRegisterAllocator{
						SafetyBoxId: safetyBoxId,
						Operator:    fromAddress,
						Address:     toAddress,
					}
				} else if action == types.DeregisterRole {
					msg = types.MsgSafetyBoxDeregisterAllocator{
						SafetyBoxId: safetyBoxId,
						Operator:    fromAddress,
						Address:     toAddress,
					}
				} else {
					return types.ErrSafetyBoxInvalidAction(types.DefaultCodespace)
				}
			case types.RoleIssuer:
				if action == types.RegisterRole {
					msg = types.MsgSafetyBoxRegisterIssuer{
						SafetyBoxId: safetyBoxId,
						Operator:    fromAddress,
						Address:     toAddress,
					}
				} else if action == types.DeregisterRole {
					msg = types.MsgSafetyBoxDeregisterIssuer{
						SafetyBoxId: safetyBoxId,
						Operator:    fromAddress,
						Address:     toAddress,
					}
				} else {
					return types.ErrSafetyBoxInvalidAction(types.DefaultCodespace)
				}
			case types.RoleReturner:
				if action == types.RegisterRole {
					msg = types.MsgSafetyBoxRegisterReturner{
						SafetyBoxId: safetyBoxId,
						Operator:    fromAddress,
						Address:     toAddress,
					}
				} else if action == types.DeregisterRole {
					msg = types.MsgSafetyBoxDeregisterReturner{
						SafetyBoxId: safetyBoxId,
						Operator:    fromAddress,
						Address:     toAddress,
					}
				} else {
					return types.ErrSafetyBoxInvalidAction(types.DefaultCodespace)
				}
			default:
				return types.ErrSafetyBoxInvalidRole(types.DefaultCodespace)
			}

			// build and sign the transaction, then broadcast to Tendermint
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func SafetyBoxSendCoinsTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sendcoins [safety_box_id] [allocate|recall|issue|return] [denom] [amount] [address] [issuer_address] ",
		Short: "Send coins among the safety box, issuers, returners and allocators. `issuer_address` is required only for issue.",
		Args:  cobra.RangeArgs(5, 6),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			safetyBoxId, action := args[0], args[1]

			denom := args[2]
			amount, err := strconv.ParseInt(args[3], 10, 0)
			if err != nil {
				return err
			}
			coins := sdk.Coins{sdk.NewCoin(denom, sdk.NewInt(amount))}

			// allocate & return -> `to` is an optional
			// recall -> `from` is an optional
			// issue -> needs both `from` and `to`
			address, err := sdk.AccAddressFromBech32(args[4])
			if err != nil {
				return err
			}
			cliCtx := client.NewCLIContextWithFrom(args[4]).WithCodec(cdc)

			var msg sdk.Msg
			switch action {
			case types.ActionAllocate:
				msg = types.NewMsgSafetyBoxAllocateCoins(safetyBoxId, address, coins)
			case types.ActionRecall:
				msg = types.NewMsgSafetyBoxRecallCoins(safetyBoxId, address, coins)
			case types.ActionIssue:
				if len(args) < 6 {
					return types.ErrSafetyBoxNeedsIssuerAddress(types.DefaultCodespace)
				}
				toAddress, err := sdk.AccAddressFromBech32(args[5])
				if err != nil {
					return err
				}
				msg = types.NewMsgSafetyBoxIssueCoins(safetyBoxId, address, toAddress, coins)
			case types.ActionReturn:
				msg = types.NewMsgSafetyBoxReturnCoins(safetyBoxId, address, coins)
			default:
				return types.ErrSafetyBoxInvalidAction(types.DefaultCodespace)
			}

			// build and sign the transaction, then broadcast to Tendermint
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}
