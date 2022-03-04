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
	"github.com/line/lbm-sdk/types"
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
	types.Bech32PrefixAccAddr + "AB0C",
	types.Bech32PrefixAccPub + "1234",
	types.Bech32PrefixValAddr + "5678",
	types.Bech32PrefixValPub + "BBAB",
	types.Bech32PrefixConsAddr + "FF04",
	types.Bech32PrefixConsPub + "6789",
}

func (s *addressTestSuite) testMarshal(original interface{}, res interface{}, marshal func() ([]byte, error), unmarshal func([]byte) error) {
	bz, err := marshal()
	s.Require().Nil(err)
	s.Require().Nil(unmarshal(bz))
	s.Require().Equal(original, res)
}

func (s *addressTestSuite) TestRandBech32PubkeyConsistency() {
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}

	for i := 0; i < 1000; i++ {
		rand.Read(pub.Key)

		mustBech32AccPub := types.MustBech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pub)
		bech32AccPub, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pub)
		s.Require().Nil(err)
		s.Require().Equal(bech32AccPub, mustBech32AccPub)

		mustBech32ValPub := types.MustBech32ifyPubKey(types.Bech32PubKeyTypeValPub, pub)
		bech32ValPub, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeValPub, pub)
		s.Require().Nil(err)
		s.Require().Equal(bech32ValPub, mustBech32ValPub)

		mustBech32ConsPub := types.MustBech32ifyPubKey(types.Bech32PubKeyTypeConsPub, pub)
		bech32ConsPub, err := types.Bech32ifyPubKey(types.Bech32PubKeyTypeConsPub, pub)
		s.Require().Nil(err)
		s.Require().Equal(bech32ConsPub, mustBech32ConsPub)

		mustAccPub := types.MustGetPubKeyFromBech32(types.Bech32PubKeyTypeAccPub, bech32AccPub)
		accPub, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeAccPub, bech32AccPub)
		s.Require().Nil(err)
		s.Require().Equal(accPub, mustAccPub)

		mustValPub := types.MustGetPubKeyFromBech32(types.Bech32PubKeyTypeValPub, bech32ValPub)
		valPub, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeValPub, bech32ValPub)
		s.Require().Nil(err)
		s.Require().Equal(valPub, mustValPub)

		mustConsPub := types.MustGetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, bech32ConsPub)
		consPub, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, bech32ConsPub)
		s.Require().Nil(err)
		s.Require().Equal(consPub, mustConsPub)

		s.Require().Equal(valPub, accPub)
		s.Require().Equal(valPub, consPub)
	}
}

func (s *addressTestSuite) TestYAMLMarshalers() {
	addr := secp256k1.GenPrivKey().PubKey().Address()

	acc := types.BytesToAccAddress(addr)
	val := types.BytesToValAddress(addr)
	cons := types.BytesToConsAddress(addr)

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

		acc := types.BytesToAccAddress(pub.Address())
		res := types.AccAddress("")

		s.testMarshal(&acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		s.testMarshal(&acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		err := types.ValidateAccAddress(str)
		s.Require().Nil(err)
		s.Require().Equal(acc, types.AccAddress(str))

		bytes, err := types.AccAddressToBytes(acc.String())
		s.Require().NoError(err)
		str = hex.EncodeToString(bytes)
		res, err = types.AccAddressFromHex(str)
		s.Require().Nil(err)
		s.Require().Equal(acc, res)
	}

	for _, str := range invalidStrs {
		_, err := types.AccAddressFromHex(str)
		s.Require().NotNil(err)

		err = types.ValidateAccAddress(str)
		s.Require().NotNil(err)

		addr := types.AccAddress("")
		err = (&addr).UnmarshalJSON([]byte("\"" + str + "\""))
		s.Require().Nil(err)
	}

	_, err := types.AccAddressFromHex("")
	s.Require().Equal("decoding Bech32 address failed: must provide an address", err.Error())
}

func (s *addressTestSuite) TestValAddr() {
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}

	for i := 0; i < 20; i++ {
		rand.Read(pub.Key)

		acc := types.BytesToValAddress(pub.Address())
		res := types.ValAddress("")

		s.testMarshal(&acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		s.testMarshal(&acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res2, err := types.ValAddressToBytes(str)
		s.Require().Nil(err)
		s.Require().Equal(acc, types.BytesToValAddress(res2))

		bytes, _ := types.ValAddressToBytes(acc.String())
		str = hex.EncodeToString(bytes)
		res, err = types.ValAddressFromHex(str)
		s.Require().Nil(err)
		s.Require().Equal(acc, res)

	}

	for _, str := range invalidStrs {
		_, err := types.ValAddressFromHex(str)
		s.Require().NotNil(err)

		err = types.ValidateValAddress(str)
		s.Require().NotNil(err)

		addr := types.ValAddress("")
		err = (&addr).UnmarshalJSON([]byte("\"" + str + "\""))
		s.Require().Nil(err)
	}

	// test empty string
	_, err := types.ValAddressFromHex("")
	s.Require().Equal("decoding Bech32 address failed: must provide an address", err.Error())
}

func (s *addressTestSuite) TestConsAddress() {
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}

	for i := 0; i < 20; i++ {
		rand.Read(pub.Key[:])

		acc := types.BytesToConsAddress(pub.Address())
		res := types.ConsAddress("")

		s.testMarshal(&acc, &res, acc.MarshalJSON, (&res).UnmarshalJSON)
		s.testMarshal(&acc, &res, acc.Marshal, (&res).Unmarshal)

		str := acc.String()
		res2, err := types.ConsAddressToBytes(str)
		s.Require().Nil(err)
		s.Require().Equal(acc, types.BytesToConsAddress(res2))

		bytes, _ := types.ConsAddressToBytes(acc.String())
		str = hex.EncodeToString(bytes)
		res, err = types.ConsAddressFromHex(str)
		s.Require().Nil(err)
		s.Require().Equal(acc, res)
	}

	for _, str := range invalidStrs {
		_, err := types.ConsAddressFromHex(str)
		s.Require().NotNil(err)

		err = types.ValidateConsAddress(str)
		s.Require().NotNil(err)

		consAddr := types.ConsAddress("")
		err = (&consAddr).UnmarshalJSON([]byte("\"" + str + "\""))
		s.Require().Nil(err)
	}

	// test empty string
	_, err := types.ConsAddressFromHex("")
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
			config := types.GetConfig()
			config.SetBech32PrefixForAccount(
				prefix+types.PrefixAccount,
				prefix+types.PrefixPublic)

			acc := types.BytesToAccAddress(pub.Address())
			s.Require().True(strings.HasPrefix(
				acc.String(),
				prefix+types.PrefixAccount), acc.String())

			bech32Pub := types.MustBech32ifyPubKey(types.Bech32PubKeyTypeAccPub, pub)
			s.Require().True(strings.HasPrefix(
				bech32Pub,
				prefix+types.PrefixPublic))

			config.SetBech32PrefixForValidator(
				prefix+types.PrefixValidator+types.PrefixAddress,
				prefix+types.PrefixValidator+types.PrefixPublic)

			val := types.BytesToValAddress(pub.Address())
			s.Require().True(strings.HasPrefix(
				val.String(),
				prefix+types.PrefixValidator+types.PrefixAddress))

			bech32ValPub := types.MustBech32ifyPubKey(types.Bech32PubKeyTypeValPub, pub)
			s.Require().True(strings.HasPrefix(
				bech32ValPub,
				prefix+types.PrefixValidator+types.PrefixPublic))

			config.SetBech32PrefixForConsensusNode(
				prefix+types.PrefixConsensus+types.PrefixAddress,
				prefix+types.PrefixConsensus+types.PrefixPublic)

			cons := types.BytesToConsAddress(pub.Address())
			s.Require().True(strings.HasPrefix(
				cons.String(),
				prefix+types.PrefixConsensus+types.PrefixAddress))

			bech32ConsPub := types.MustBech32ifyPubKey(types.Bech32PubKeyTypeConsPub, pub)
			s.Require().True(strings.HasPrefix(
				bech32ConsPub,
				prefix+types.PrefixConsensus+types.PrefixPublic))
		}
	}
}

func (s *addressTestSuite) TestAddressInterface() {
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}
	rand.Read(pub.Key)

	addrs := []types.Address{
		types.BytesToConsAddress(pub.Address()),
		types.BytesToValAddress(pub.Address()),
		types.BytesToAccAddress(pub.Address()),
	}

	for _, addr := range addrs {
		switch addr := addr.(type) {
		case types.AccAddress:
			err := types.ValidateAccAddress(addr.String())
			s.Require().Nil(err)
		case types.ValAddress:
			err := types.ValidateValAddress(addr.String())
			s.Require().Nil(err)
		case types.ConsAddress:
			err := types.ValidateConsAddress(addr.String())
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

	err := types.VerifyAddressFormat(addr0)
	s.Require().EqualError(err, "addresses cannot be empty: unknown address")
	err = types.VerifyAddressFormat(addr5)
	s.Require().NoError(err)
	err = types.VerifyAddressFormat(addr20)
	s.Require().NoError(err)
	err = types.VerifyAddressFormat(addr32)
	s.Require().NoError(err)
	err = types.VerifyAddressFormat(addr256)
	s.Require().EqualError(err, "address max length is 255, got 256: unknown address")
}

func (s *addressTestSuite) TestCustomAddressVerifier() {
	// Create a 10 byte address
	addr := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	accBech := types.BytesToAccAddress(addr).String()
	valBech := types.BytesToValAddress(addr).String()
	consBech := types.BytesToConsAddress(addr).String()
	// Verifiy that the default logic rejects this 10 byte address
	err := types.VerifyAddressFormat(addr)
	s.Require().Nil(err)
	err = types.ValidateAccAddress(accBech)
	s.Require().Nil(err)
	err = types.ValidateValAddress(valBech)
	s.Require().Nil(err)
	err = types.ValidateConsAddress(consBech)
	s.Require().Nil(err)

	// Set a custom address verifier only accepts 20 byte addresses
	types.GetConfig().SetAddressVerifier(func(bz []byte) error {
		n := len(bz)
		if n == types.BytesAddrLen {
			return nil
		}
		return fmt.Errorf("incorrect address length %d", n)
	})

	// Verifiy that the custom logic rejects this 10 byte address
	err = types.VerifyAddressFormat(addr)
	s.Require().NotNil(err)
	err = types.ValidateAccAddress(accBech)
	s.Require().NotNil(err)
	err = types.ValidateValAddress(valBech)
	s.Require().NotNil(err)
	err = types.ValidateConsAddress(consBech)
	s.Require().NotNil(err)

	// Reinitialize the global config to default address verifier (nil)
	types.GetConfig().SetAddressVerifier(nil)
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
			got, err := types.Bech32ifyAddressBytes(tt.args.prefix, tt.args.bs)
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
				require.Panics(t, func() { types.MustBech32ifyAddressBytes(tt.args.prefix, tt.args.bs) })
				return
			}
			require.Equal(t, tt.want, types.MustBech32ifyAddressBytes(tt.args.prefix, tt.args.bs))
		})
	}
}

func (s *addressTestSuite) TestAddressTypesEquals() {
	addr1 := secp256k1.GenPrivKey().PubKey().Address()
	accAddr1 := types.BytesToAccAddress(addr1)
	consAddr1 := types.BytesToConsAddress(addr1)
	valAddr1 := types.BytesToValAddress(addr1)

	addr2 := secp256k1.GenPrivKey().PubKey().Address()
	accAddr2 := types.BytesToAccAddress(addr2)
	consAddr2 := types.BytesToConsAddress(addr2)
	valAddr2 := types.BytesToValAddress(addr2)

	// equality
	s.Require().True(accAddr1.Equals(accAddr1))
	s.Require().True(consAddr1.Equals(consAddr1))
	s.Require().True(valAddr1.Equals(valAddr1))

	// emptiness
	s.Require().True(types.AccAddress("").Equals(types.AccAddress("")))
	s.Require().True(types.ConsAddress("").Equals(types.ConsAddress("")))
	s.Require().True(types.ValAddress("").Equals(types.ValAddress("")))

	s.Require().False(accAddr1.Equals(accAddr2))
	s.Require().Equal(accAddr1.Equals(accAddr2), accAddr2.Equals(accAddr1))
	s.Require().False(consAddr1.Equals(consAddr2))
	s.Require().Equal(consAddr1.Equals(consAddr2), consAddr2.Equals(consAddr1))
	s.Require().False(valAddr1.Equals(valAddr2))
	s.Require().Equal(valAddr1.Equals(valAddr2), valAddr2.Equals(valAddr1))
}

func (s *addressTestSuite) TestNilAddressTypesEmpty() {
	s.Require().True(types.AccAddress("").Empty())
	s.Require().True(types.ConsAddress("").Empty())
	s.Require().True(types.ValAddress("").Empty())
}

func (s *addressTestSuite) TestGetConsAddress() {
	pk := secp256k1.GenPrivKey().PubKey()
	s.Require().NotEqual(types.GetConsAddress(pk), pk.Address())
	consBytes, _ := types.ConsAddressToBytes(types.GetConsAddress(pk).String())
	s.Require().True(bytes.Equal(consBytes, pk.Address()))
	s.Require().Panics(func() { types.GetConsAddress(cryptotypes.PubKey(nil)) })
}

func (s *addressTestSuite) TestGetFromBech32() {
	_, err := types.GetFromBech32("", "prefix")
	s.Require().Error(err)
	s.Require().Equal("decoding Bech32 address failed: must provide an address", err.Error())
	_, err = types.GetFromBech32("cosmos1qqqsyqcyq5rqwzqfys8f67", "x")
	s.Require().Error(err)
	s.Require().Equal("invalid Bech32 prefix; expected x, got cosmos", err.Error())
}

func (s *addressTestSuite) TestBech32Cache() {
	pubBz := make([]byte, ed25519.PubKeySize)
	pub := &ed25519.PubKey{Key: pubBz}

	s.T().Log("access bech32ToAddrCache before access addrToBech32Cache")
	{
		rand.Read(pub.Key)
		addr := types.BytesToAccAddress(pub.Address())
		bech32Addr := addr.String()
		types.SetBech32Cache(types.DefaultBech32CacheSize)

		rawAddr, err := types.AccAddressToBytes(bech32Addr)
		s.Require().Nil(err)
		require.Equal(s.T(), addr, types.BytesToAccAddress(rawAddr))

		require.Equal(s.T(), bech32Addr, addr.String())
	}
	s.T().Log("access addrToBech32Cache before access bech32ToAddrCache")
	{
		rand.Read(pub.Key)
		addr := types.BytesToAccAddress(pub.Address())
		bech32Addr := addr.String()
		types.SetBech32Cache(types.DefaultBech32CacheSize)

		require.Equal(s.T(), bech32Addr, addr.String())

		rawAddr, err := types.AccAddressToBytes(bech32Addr)
		s.Require().Nil(err)
		require.Equal(s.T(), addr, types.BytesToAccAddress(rawAddr))
	}
	types.InvalidateBech32Cache()
}
