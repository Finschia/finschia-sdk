package types_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"

	"github.com/line/lbm-sdk/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/line/lbm-sdk/crypto/types"
	sdk "github.com/line/lbm-sdk/types"
	legacybech32 "github.com/line/lbm-sdk/types/bech32/legacybech32"
)

type addressTestSuite struct {
	suite.Suite
}

func TestAddressTestSuite(t *testing.T) {
	suite.Run(t, new(addressTestSuite))
}

func (s *addressTestSuite) SetupSuite() {
	s.T().Parallel()
}

var invalidStrs = []string{
	"hello, world!",
	"0xAA",
	"AAA",
	sdk.Bech32PrefixAccAddr + "AB0C",
	sdk.Bech32PrefixAccPub + "1234",
	sdk.Bech32PrefixValAddr + "5678",
	sdk.Bech32PrefixValPub + "BBAB",
	sdk.Bech32PrefixConsAddr + "FF04",
	sdk.Bech32PrefixConsPub + "6789",
}

func (s *addressTestSuite) testMarshal(original interface{}, res interface{}, marshal func() ([]byte, error), unmarshal func([]byte) error) {
	bz, err := marshal()
	s.Require().Nil(err)
	s.Require().Nil(unmarshal(bz))
	s.Require().Equal(original, res)
}

func (s *addressTestSuite) TestEmptyAddresses() {
	s.T().Parallel()
	s.Require().Equal(sdk.AccAddress("").String(), "")
	s.Require().Equal(sdk.ValAddress("").String(), "")
	s.Require().Equal(sdk.ConsAddress("").String(), "")

	accAddr := sdk.BytesToAccAddress([]byte(""))
	s.Require().True(accAddr.Empty())
	err := sdk.ValidateAccAddress(accAddr.String())
	s.Require().Error(err)

	valAddr := sdk.BytesToValAddress([]byte(""))
	s.Require().True(valAddr.Empty())
	err = sdk.ValidateValAddress(valAddr.String())
	s.Require().Error(err)

	consAddr := sdk.BytesToConsAddress([]byte(""))
	s.Require().True(consAddr.Empty())
	err = sdk.ValidateConsAddress(consAddr.String())
	s.Require().Error(err)
}

func (s *addressTestSuite) TestYAMLMarshalers() {
	addr := secp256k1.GenPrivKey().PubKey().Address()

	acc := sdk.BytesToAccAddress(addr)
	val := sdk.BytesToValAddress(addr)
	cons := sdk.BytesToConsAddress(addr)

	got, _ := yaml.Marshal(&acc)
	s.Require().Equal(acc.String()+"\n", string(got))

	got, _ = yaml.Marshal(&val)
	s.Require().Equal(val.String()+"\n", string(got))

	got, _ = yaml.Marshal(&cons)
	s.Require().Equal(cons.String()+"\n", string(got))
}

func (s *addressTestSuite) TestRandBech32AccAddrConsistency() {
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}

	for i := 0; i < 1000; i++ {
		rand.Read(pub.Key)

		acc := sdk.BytesToAccAddress(pub.Address())
		res := sdk.AccAddress("")

		s.testMarshal(&acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		s.testMarshal(&acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		err := sdk.ValidateAccAddress(str)
		s.Require().Nil(err)
		s.Require().Equal(acc, sdk.AccAddress(str))

		bytes, err := sdk.AccAddressToBytes(acc.String())
		s.Require().NoError(err)
		str = hex.EncodeToString(bytes)
		res, err = sdk.AccAddressFromHex(str)
		s.Require().Nil(err)
		s.Require().Equal(acc, res)
	}

	for _, str := range invalidStrs {
		_, err := sdk.AccAddressFromHex(str)
		s.Require().NotNil(err)

		err = sdk.ValidateAccAddress(str)
		s.Require().NotNil(err)

		addr := sdk.AccAddress("")
		err = (&addr).UnmarshalJSON([]byte("\"" + str + "\""))
		s.Require().Nil(err)
	}

	_, err := sdk.AccAddressFromHex("")
	s.Require().Equal("decoding Bech32 address failed: must provide an address", err.Error())
}

func (s *addressTestSuite) TestValAddr() {
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}

	for i := 0; i < 20; i++ {
		rand.Read(pub.Key)

		acc := sdk.BytesToValAddress(pub.Address())
		res := sdk.ValAddress("")

		s.testMarshal(&acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		s.testMarshal(&acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res2, err := sdk.ValAddressToBytes(str)
		s.Require().Nil(err)
		s.Require().Equal(acc, sdk.BytesToValAddress(res2))

		bytes, _ := sdk.ValAddressToBytes(acc.String())
		str = hex.EncodeToString(bytes)
		res, err = sdk.ValAddressFromHex(str)
		s.Require().Nil(err)
		s.Require().Equal(acc, res)

	}

	for _, str := range invalidStrs {
		_, err := sdk.ValAddressFromHex(str)
		s.Require().NotNil(err)

		err = sdk.ValidateValAddress(str)
		s.Require().NotNil(err)

		addr := sdk.ValAddress("")
		err = (&addr).UnmarshalJSON([]byte("\"" + str + "\""))
		s.Require().Nil(err)
	}

	// test empty string
	_, err := sdk.ValAddressFromHex("")
	s.Require().Equal("decoding Bech32 address failed: must provide an address", err.Error())
}

func (s *addressTestSuite) TestConsAddress() {
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}

	for i := 0; i < 20; i++ {
		rand.Read(pub.Key[:])

		acc := sdk.BytesToConsAddress(pub.Address())
		res := sdk.ConsAddress("")

		s.testMarshal(&acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		s.testMarshal(&acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res2, err := sdk.ConsAddressToBytes(str)
		s.Require().Nil(err)
		s.Require().Equal(acc, sdk.BytesToConsAddress(res2))

		bytes, _ := sdk.ConsAddressToBytes(acc.String())
		str = hex.EncodeToString(bytes)
		res, err = sdk.ConsAddressFromHex(str)
		s.Require().Nil(err)
		s.Require().Equal(acc, res)
	}

	for _, str := range invalidStrs {
		_, err := sdk.ConsAddressFromHex(str)
		s.Require().NotNil(err)

		err = sdk.ValidateConsAddress(str)
		s.Require().NotNil(err)

		consAddr := sdk.ConsAddress("")
		err = (&consAddr).UnmarshalJSON([]byte("\"" + str + "\""))
		s.Require().Nil(err)
	}

	// test empty string
	_, err := sdk.ConsAddressFromHex("")
	s.Require().Equal("decoding Bech32 address failed: must provide an address", err.Error())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *addressTestSuite) TestConfiguredPrefix() {
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}
	for length := 1; length < 10; length++ {
		for times := 1; times < 20; times++ {
			rand.Read(pub.Key[:])
			// Test if randomly generated prefix of a given length works
			prefix := RandString(length)

			// Assuming that GetConfig is not sealed.
			config := sdk.GetConfig()
			config.SetBech32PrefixForAccount(
				prefix+sdk.PrefixAccount,
				prefix+sdk.PrefixPublic)

			acc := sdk.BytesToAccAddress(pub.Address())
			s.Require().True(strings.HasPrefix(
				acc.String(),
				prefix+sdk.PrefixAccount), acc.String())

			bech32Pub := legacybech32.MustMarshalPubKey(legacybech32.AccPK, pub)
			s.Require().True(strings.HasPrefix(
				bech32Pub,
				prefix+sdk.PrefixPublic))

			config.SetBech32PrefixForValidator(
				prefix+sdk.PrefixValidator+sdk.PrefixAddress,
				prefix+sdk.PrefixValidator+sdk.PrefixPublic)

			val := sdk.BytesToValAddress(pub.Address())
			s.Require().True(strings.HasPrefix(
				val.String(),
				prefix+sdk.PrefixValidator+sdk.PrefixAddress))

			bech32ValPub := legacybech32.MustMarshalPubKey(legacybech32.ValPK, pub)
			s.Require().True(strings.HasPrefix(
				bech32ValPub,
				prefix+sdk.PrefixValidator+sdk.PrefixPublic))

			config.SetBech32PrefixForConsensusNode(
				prefix+sdk.PrefixConsensus+sdk.PrefixAddress,
				prefix+sdk.PrefixConsensus+sdk.PrefixPublic)

			cons := sdk.BytesToConsAddress(pub.Address())
			s.Require().True(strings.HasPrefix(
				cons.String(),
				prefix+sdk.PrefixConsensus+sdk.PrefixAddress))

			bech32ConsPub := legacybech32.MustMarshalPubKey(legacybech32.ConsPK, pub)
			s.Require().True(strings.HasPrefix(
				bech32ConsPub,
				prefix+sdk.PrefixConsensus+sdk.PrefixPublic))
		}
	}
}

func (s *addressTestSuite) TestAddressInterface() {
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}
	rand.Read(pub.Key)

	addrs := []sdk.Address{
		sdk.BytesToConsAddress(pub.Address()),
		sdk.BytesToValAddress(pub.Address()),
		sdk.BytesToAccAddress(pub.Address()),
	}

	for _, addr := range addrs {
		switch addr := addr.(type) {
		case sdk.AccAddress:
			err := sdk.ValidateAccAddress(addr.String())
			s.Require().Nil(err)
		case sdk.ValAddress:
			err := sdk.ValidateValAddress(addr.String())
			s.Require().Nil(err)
		case sdk.ConsAddress:
			err := sdk.ValidateConsAddress(addr.String())
			s.Require().Nil(err)
		default:
			s.T().Fail()
		}
	}

}

func (s *addressTestSuite) TestVerifyAddressFormat() {
	addr0 := make([]byte, 0)
	addr5 := make([]byte, 5)
	addr20 := make([]byte, 20)
	addr32 := make([]byte, 32)
	addr256 := make([]byte, 256)

	err := sdk.VerifyAddressFormat(addr0)
	s.Require().EqualError(err, "addresses cannot be empty: unknown address")
	err = sdk.VerifyAddressFormat(addr5)
	s.Require().NoError(err)
	err = sdk.VerifyAddressFormat(addr20)
	s.Require().NoError(err)
	err = sdk.VerifyAddressFormat(addr32)
	s.Require().NoError(err)
	err = sdk.VerifyAddressFormat(addr256)
	s.Require().EqualError(err, "address max length is 255, got 256: unknown address")
}

func (s *addressTestSuite) TestCustomAddressVerifier() {
	// Create a 10 byte address
	addr := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	accBech := sdk.BytesToAccAddress(addr).String()
	valBech := sdk.BytesToValAddress(addr).String()
	consBech := sdk.BytesToConsAddress(addr).String()
	// Verifiy that the default logic rejects this 10 byte address
	err := sdk.VerifyAddressFormat(addr)
	s.Require().Nil(err)
	err = sdk.ValidateAccAddress(accBech)
	s.Require().Nil(err)
	err = sdk.ValidateValAddress(valBech)
	s.Require().Nil(err)
	err = sdk.ValidateConsAddress(consBech)
	s.Require().Nil(err)

	// Set a custom address verifier only accepts 20 byte addresses
	sdk.GetConfig().SetAddressVerifier(func(bz []byte) error {
		n := len(bz)
		if n == sdk.BytesAddrLen {
			return nil
		}
		return fmt.Errorf("incorrect address length %d", n)
	})

	// Verifiy that the custom logic rejects this 10 byte address
	err = sdk.VerifyAddressFormat(addr)
	s.Require().NotNil(err)
	err = sdk.ValidateAccAddress(accBech)
	s.Require().NotNil(err)
	err = sdk.ValidateValAddress(valBech)
	s.Require().NotNil(err)
	err = sdk.ValidateConsAddress(consBech)
	s.Require().NotNil(err)

	// Reinitialize the global config to default address verifier (nil)
	sdk.GetConfig().SetAddressVerifier(nil)
}

func (s *addressTestSuite) TestBech32ifyAddressBytes() {
	addr10byte := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	addr20byte := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	type args struct {
		prefix string
		bs     []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty address", args{"prefixa", []byte{}}, "", false},
		{"empty prefix", args{"", addr20byte}, "", true},
		{"10-byte address", args{"prefixa", addr10byte}, "prefixa1qqqsyqcyq5rqwzqf3953cc", false},
		{"10-byte address", args{"prefixb", addr10byte}, "prefixb1qqqsyqcyq5rqwzqf20xxpc", false},
		{"20-byte address", args{"prefixa", addr20byte}, "prefixa1qqqsyqcyq5rqwzqfpg9scrgwpugpzysn7hzdtn", false},
		{"20-byte address", args{"prefixb", addr20byte}, "prefixb1qqqsyqcyq5rqwzqfpg9scrgwpugpzysnrujsuw", false},
	}
	for _, tt := range tests {
		tt := tt
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := sdk.Bech32ifyAddressBytes(tt.args.prefix, tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bech32ifyBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func (s *addressTestSuite) TestMustBech32ifyAddressBytes() {
	addr10byte := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	addr20byte := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	type args struct {
		prefix string
		bs     []byte
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantPanic bool
	}{
		{"empty address", args{"prefixa", []byte{}}, "", false},
		{"empty prefix", args{"", addr20byte}, "", true},
		{"10-byte address", args{"prefixa", addr10byte}, "prefixa1qqqsyqcyq5rqwzqf3953cc", false},
		{"10-byte address", args{"prefixb", addr10byte}, "prefixb1qqqsyqcyq5rqwzqf20xxpc", false},
		{"20-byte address", args{"prefixa", addr20byte}, "prefixa1qqqsyqcyq5rqwzqfpg9scrgwpugpzysn7hzdtn", false},
		{"20-byte address", args{"prefixb", addr20byte}, "prefixb1qqqsyqcyq5rqwzqfpg9scrgwpugpzysnrujsuw", false},
	}
	for _, tt := range tests {
		tt := tt
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				require.Panics(t, func() { sdk.MustBech32ifyAddressBytes(tt.args.prefix, tt.args.bs) })
				return
			}
			require.Equal(t, tt.want, sdk.MustBech32ifyAddressBytes(tt.args.prefix, tt.args.bs))
		})
	}
}

func (s *addressTestSuite) TestAddressTypesEquals() {
	addr1 := secp256k1.GenPrivKey().PubKey().Address()
	accAddr1 := sdk.BytesToAccAddress(addr1)
	consAddr1 := sdk.BytesToConsAddress(addr1)
	valAddr1 := sdk.BytesToValAddress(addr1)

	addr2 := secp256k1.GenPrivKey().PubKey().Address()
	accAddr2 := sdk.BytesToAccAddress(addr2)
	consAddr2 := sdk.BytesToConsAddress(addr2)
	valAddr2 := sdk.BytesToValAddress(addr2)

	// equality
	s.Require().True(accAddr1.Equals(accAddr1))
	s.Require().True(consAddr1.Equals(consAddr1))
	s.Require().True(valAddr1.Equals(valAddr1))

	// emptiness
	s.Require().True(sdk.AccAddress("").Equals(sdk.AccAddress("")))
	s.Require().True(sdk.ConsAddress("").Equals(sdk.ConsAddress("")))
	s.Require().True(sdk.ValAddress("").Equals(sdk.ValAddress("")))

	s.Require().False(accAddr1.Equals(accAddr2))
	s.Require().Equal(accAddr1.Equals(accAddr2), accAddr2.Equals(accAddr1))
	s.Require().False(consAddr1.Equals(consAddr2))
	s.Require().Equal(consAddr1.Equals(consAddr2), consAddr2.Equals(consAddr1))
	s.Require().False(valAddr1.Equals(valAddr2))
	s.Require().Equal(valAddr1.Equals(valAddr2), valAddr2.Equals(valAddr1))
}

func (s *addressTestSuite) TestNilAddressTypesEmpty() {
	s.Require().True(sdk.AccAddress("").Empty())
	s.Require().True(sdk.ConsAddress("").Empty())
	s.Require().True(sdk.ValAddress("").Empty())
}

func (s *addressTestSuite) TestGetConsAddress() {
	pk := secp256k1.GenPrivKey().PubKey()
	s.Require().NotEqual(sdk.GetConsAddress(pk), pk.Address())
	consBytes, _ := sdk.ConsAddressToBytes(sdk.GetConsAddress(pk).String())
	s.Require().True(bytes.Equal(consBytes, pk.Address()))
	s.Require().Panics(func() { sdk.GetConsAddress(cryptotypes.PubKey(nil)) })
}

func (s *addressTestSuite) TestGetFromBech32() {
	_, err := sdk.GetFromBech32("", "prefix")
	s.Require().Error(err)
	s.Require().Equal("decoding Bech32 address failed: must provide a non empty address", err.Error())
	_, err = sdk.GetFromBech32("link1qqqsyqcyq5rqwzqf97tnae", "x")
	s.Require().Error(err)
	s.Require().Equal("invalid Bech32 prefix; expected x, got link", err.Error())
}

func (s *addressTestSuite) TestBech32Cache() {
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}

	s.T().Log("access bech32ToAddrCache before access addrToBech32Cache")
	{
		rand.Read(pub.Key)
		addr := sdk.BytesToAccAddress(pub.Address())
		bech32Addr := addr.String()
		sdk.SetBech32Cache(sdk.DefaultBech32CacheSize)

		rawAddr, err := sdk.AccAddressToBytes(bech32Addr)
		s.Require().Nil(err)
		require.Equal(s.T(), addr, sdk.BytesToAccAddress(rawAddr))

		require.Equal(s.T(), bech32Addr, addr.String())
	}
	s.T().Log("access addrToBech32Cache before access bech32ToAddrCache")
	{
		rand.Read(pub.Key)
		addr := sdk.BytesToAccAddress(pub.Address())
		bech32Addr := addr.String()
		sdk.SetBech32Cache(sdk.DefaultBech32CacheSize)

		require.Equal(s.T(), bech32Addr, addr.String())

		rawAddr, err := sdk.AccAddressToBytes(bech32Addr)
		s.Require().Nil(err)
		require.Equal(s.T(), addr, sdk.BytesToAccAddress(rawAddr))
	}
	sdk.InvalidateBech32Cache()
}
