package types_test

import (
	"fmt"

	codectypes "github.com/line/lbm-sdk/codec/types"
	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/line/lbm-sdk/crypto/types"
	sdk "github.com/line/lbm-sdk/types"
)

var (
	pk1      = ed25519.GenPrivKey().PubKey()
	pk1Any   *codectypes.Any
	pk2      = ed25519.GenPrivKey().PubKey()
	pk3      = ed25519.GenPrivKey().PubKey()
	addr1, _ = sdk.Bech32ifyAddressBytes(sdk.Bech32PrefixAccAddr, pk1.Address())
	addr2, _ = sdk.Bech32ifyAddressBytes(sdk.Bech32PrefixAccAddr, pk2.Address())
	addr3, _ = sdk.Bech32ifyAddressBytes(sdk.Bech32PrefixAccAddr, pk3.Address())
	valAddr1 = sdk.BytesToValAddress(pk1.Address())
	valAddr2 = sdk.BytesToValAddress(pk2.Address())
	valAddr3 = sdk.BytesToValAddress(pk3.Address())

	emptyAddr   sdk.ValAddress
	emptyPubkey cryptotypes.PubKey
)

func init() {
	var err error
	pk1Any, err = codectypes.NewAnyWithValue(pk1)
	if err != nil {
		panic(fmt.Sprintf("Can't pack pk1 %t as Any", pk1))
	}
}
