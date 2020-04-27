package cli

import (
	"bufio"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/line/link/client"
	"github.com/line/link/x/safetybox/internal/types"
	"github.com/spf13/cobra"
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
		SafetyBoxSendTokenTxCmd(cdc),
	)
	return txCmd
}

func SafetyBoxCreateTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [id] [owner_address] [contract_id]",
		Short: "Create a safety box with ID, owner and contractID. Only one owner and contractID is allowed.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := client.NewCLIContextWithInputAndFrom(inBuf, args[1]).WithCodec(cdc)
			safetyBoxID := args[0]
			safetyBoxOwner := cliCtx.FromAddress
			contractID := args[2]

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.MsgSafetyBoxCreate{SafetyBoxID: safetyBoxID, SafetyBoxOwner: safetyBoxOwner, ContractID: contractID}
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
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			safetyBoxID, action, role := args[0], args[1], args[2]
			cliCtx := client.NewCLIContextWithInputAndFrom(inBuf, args[3]).WithCodec(cdc)
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
				switch action {
				case types.RegisterRole:
					msg = types.MsgSafetyBoxRegisterOperator{
						SafetyBoxID:    safetyBoxID,
						SafetyBoxOwner: fromAddress,
						Address:        toAddress,
					}
				case types.DeregisterRole:
					msg = types.MsgSafetyBoxDeregisterOperator{
						SafetyBoxID:    safetyBoxID,
						SafetyBoxOwner: fromAddress,
						Address:        toAddress,
					}
				default:
					return sdkerrors.Wrapf(types.ErrSafetyBoxInvalidAction, "Action: %s", action)
				}
			case types.RoleAllocator:
				switch action {
				case types.RegisterRole:
					msg = types.MsgSafetyBoxRegisterAllocator{
						SafetyBoxID: safetyBoxID,
						Operator:    fromAddress,
						Address:     toAddress,
					}
				case types.DeregisterRole:
					msg = types.MsgSafetyBoxDeregisterAllocator{
						SafetyBoxID: safetyBoxID,
						Operator:    fromAddress,
						Address:     toAddress,
					}
				default:
					return sdkerrors.Wrapf(types.ErrSafetyBoxInvalidAction, "Action: %s", action)
				}
			case types.RoleIssuer:
				switch action {
				case types.RegisterRole:
					msg = types.MsgSafetyBoxRegisterIssuer{
						SafetyBoxID: safetyBoxID,
						Operator:    fromAddress,
						Address:     toAddress,
					}
				case types.DeregisterRole:
					msg = types.MsgSafetyBoxDeregisterIssuer{
						SafetyBoxID: safetyBoxID,
						Operator:    fromAddress,
						Address:     toAddress,
					}
				default:
					return sdkerrors.Wrapf(types.ErrSafetyBoxInvalidAction, "Action: %s", action)
				}
			case types.RoleReturner:
				switch action {
				case types.RegisterRole:
					msg = types.MsgSafetyBoxRegisterReturner{
						SafetyBoxID: safetyBoxID,
						Operator:    fromAddress,
						Address:     toAddress,
					}
				case types.DeregisterRole:
					msg = types.MsgSafetyBoxDeregisterReturner{
						SafetyBoxID: safetyBoxID,
						Operator:    fromAddress,
						Address:     toAddress,
					}
				default:
					return sdkerrors.Wrapf(types.ErrSafetyBoxInvalidAction, "Action: %s", action)
				}
			default:
				return sdkerrors.Wrapf(types.ErrSafetyBoxInvalidRole, "Role: %s", role)
			}

			// build and sign the transaction, then broadcast to Tendermint
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}

func SafetyBoxSendTokenTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sendtoken [safety_box_id] [allocate|recall|issue|return] [contractID] [amount] [address] [issuer_address] ",
		Short: "Send token among the safety box, issuers, returners and allocators. `issuer_address` is required only for issue.",
		Args:  cobra.RangeArgs(5, 6),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			safetyBoxID, action := args[0], args[1]

			contractID := args[2]
			amount, err := strconv.ParseInt(args[3], 10, 0)
			if err != nil {
				return err
			}

			// allocate & return -> `to` is an optional
			// recall -> `from` is an optional
			// issue -> needs both `from` and `to`
			address, err := sdk.AccAddressFromBech32(args[4])
			if err != nil {
				return err
			}
			cliCtx := client.NewCLIContextWithInputAndFrom(inBuf, args[4]).WithCodec(cdc)

			var msg sdk.Msg
			switch action {
			case types.ActionAllocate:
				msg = types.NewMsgSafetyBoxAllocateToken(safetyBoxID, address, contractID, sdk.NewInt(amount))
			case types.ActionRecall:
				msg = types.NewMsgSafetyBoxRecallToken(safetyBoxID, address, contractID, sdk.NewInt(amount))
			case types.ActionIssue:
				if len(args) < 6 {
					return types.ErrSafetyBoxIssuerAddressRequired
				}
				toAddress, err := sdk.AccAddressFromBech32(args[5])
				if err != nil {
					return err
				}
				msg = types.NewMsgSafetyBoxIssueToken(safetyBoxID, address, toAddress, contractID, sdk.NewInt(amount))
			case types.ActionReturn:
				msg = types.NewMsgSafetyBoxReturnToken(safetyBoxID, address, contractID, sdk.NewInt(amount))
			default:
				return sdkerrors.Wrapf(types.ErrSafetyBoxInvalidAction, "Action: %s", action)
			}

			// build and sign the transaction, then broadcast to Tendermint
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return client.PostCommands(cmd)[0]
}
