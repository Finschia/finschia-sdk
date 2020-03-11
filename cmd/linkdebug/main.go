package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	link "github.com/line/link/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/line/link/types"
)

func init() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(types.Bech32PrefixAccAddr, types.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(types.Bech32PrefixValAddr, types.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(types.Bech32PrefixConsAddr, types.Bech32PrefixConsPub)
	config.SetCoinType(types.CoinType)
	config.SetFullFundraiserPath(types.FullFundraiserPath)
	config.Seal()

	rootCmd.AddCommand(txCmd)
	rootCmd.AddCommand(addrCmd)
	rootCmd.AddCommand(hackCmd)
	rootCmd.AddCommand(rawBytesCmd)
}

var rootCmd = &cobra.Command{
	Use:          "linkdebug",
	Short:        "LINK debug tool",
	SilenceUsage: true,
}

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Decode a link tx from hex or base64",
	RunE:  runTxCmd,
}

var addrCmd = &cobra.Command{
	Use:   "addr",
	Short: "Convert an address between hex and bech32",
	RunE:  runAddrCmd,
}

var hackCmd = &cobra.Command{
	Use:   "hack",
	Short: "Boilerplate to Hack on an existing state by scripting some Go...",
	RunE:  runHackCmd,
}

var rawBytesCmd = &cobra.Command{
	Use:   "raw-bytes",
	Short: "Convert raw bytes output (eg. [10 21 13 255]) to hex",
	RunE:  runRawBytesCmd,
}

func runRawBytesCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected single arg")
	}
	stringBytes := args[0]
	stringBytes = strings.Trim(stringBytes, "[")
	stringBytes = strings.Trim(stringBytes, "]")
	spl := strings.Split(stringBytes, " ")

	byteArray := []byte{}
	for _, s := range spl {
		b, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		byteArray = append(byteArray, byte(b))
	}
	fmt.Printf("%X\n", byteArray)
	return nil
}

func runAddrCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected single arg")
	}

	addrString := args[0]
	var addr []byte

	// try hex, then bech32
	var err error
	addr, err = hex.DecodeString(addrString)
	if err != nil {
		var err2 error
		addr, err2 = sdk.AccAddressFromBech32(addrString)
		if err2 != nil {
			var err3 error
			addr, err3 = sdk.ValAddressFromBech32(addrString)

			if err3 != nil {
				return fmt.Errorf(`expected hex or bech32. Got errors:
			hex: %v,
			bech32 acc: %v
			bech32 val: %v
			`, err, err2, err3)
			}
		}
	}

	accAddr := sdk.AccAddress(addr)
	valAddr := sdk.ValAddress(addr)

	fmt.Println("Address:", addr)
	fmt.Printf("Address (hex): %X\n", addr)
	fmt.Printf("Bech32 Acc: %s\n", accAddr)
	fmt.Printf("Bech32 Val: %s\n", valAddr)
	return nil
}

func runTxCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("expected single arg")
	}

	txString := args[0]

	// try hex, then base64
	txBytes, err := hex.DecodeString(txString)
	if err != nil {
		var err2 error
		txBytes, err2 = base64.StdEncoding.DecodeString(txString)
		if err2 != nil {
			return fmt.Errorf(`expected hex or base64. Got errors:
			hex: %v,
			base64: %v
			`, err, err2)
		}
	}

	var tx = auth.StdTx{}
	cdc := link.MakeCodec()

	err = cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return err
	}

	bz, err := cdc.MarshalJSON(tx)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte{})
	err = json.Indent(buf, bz, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(buf.String())
	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
