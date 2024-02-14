package types

import (
	"encoding/base64"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadJWKn(t *testing.T) {
	jwkN := "qwrzl06fwB6OIm62IxNG7NXNIDmgdBrvf09ob2Gsp6ZmAXgU4trHPUYrdBaAlU5aHpchXCf_mVL-U5dzRqeVFQsVqsj4PEIE6E5OPw8EwumP2fzLQSswpkKmJJKFcdncfQ730QBonRUEhKkIbiYdicJl5yTkORd0_BmfdLV98r-sEwEHN4lzTJ15-yw90ob_R6vAH4wPyCSN3Xe5_zV6R4ENL2NlKn2HT9lbV7HhtQongea8wfnthUhdZH38kI4SS5nAaCVNxEAzlvJtUIdCpSgjUgcbah-DwY39l4D800kLxkcF2CGXPSmpF8GPs1aWSsYupY8sTSy9qCFJFPFx8Q"

	modulus, err := base64.RawURLEncoding.DecodeString(jwkN)
	if err != nil {
		require.Error(t, err)
		return
	}

	publicBigInt := new(big.Int).SetBytes(modulus)

	expected, _ := new(big.Int).SetString("21592150548688579459736470743740609932734247489275801363800063955060727299173213708631581411829104889800393345041018770658141456640154383655030888058825927685116282469808205653751534030009509399616836504388549230657354705854346635974316335659898523674443506501951289176826338661840898157904351950047060188748641336012482156218551743045002776897928784683784942521239939886566984329457188054591552006485891818434347748745176793274089100897612684086578201319592157538096356516151081984342016690451396306034536873405640563042576827222214953500041471360624269713141770498106051868926368335925348133899050033256474323546609", 10)
	require.Equal(t, expected, publicBigInt)
}

func TestHashAsciiStrToField(t *testing.T) {
	testCase := map[string]struct {
		input    string
		maxSize  int
		expected string
	}{
		"test@gmail.com, 30": {
			"test@gmail.com",
			30,
			"13606676331558803166736332982602687405662978305929711411606106012181987145625",
		},
		"test@gmail.com, 32": {
			"test@gmail.com",
			32,
			"10404231015713323946367565043703223078961469658905861259850380980432751872181",
		},
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjEifQ, 248": {
			"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjEifQ",
			248,
			"10859137172532636243875876865378218840892896099608302223608404291948352005840",
		},
	}

	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			hashed, err := hashASCIIStrToField(tc.input, tc.maxSize)
			require.NoError(t, err)
			require.Equal(t, tc.expected, hashed.String())
		})
	}
}

func TestCircomBigIntToField(t *testing.T) {
	pubKey, _ := new(big.Int).SetString("21592150548688579459736470743740609932734247489275801363800063955060727299173213708631581411829104889800393345041018770658141456640154383655030888058825927685116282469808205653751534030009509399616836504388549230657354705854346635974316335659898523674443506501951289176826338661840898157904351950047060188748641336012482156218551743045002776897928784683784942521239939886566984329457188054591552006485891818434347748745176793274089100897612684086578201319592157538096356516151081984342016690451396306034536873405640563042576827222214953500041471360624269713141770498106051868926368335925348133899050033256474323546609", 10)
	expectedBytes := []string{
		"1029020885788392288961761165752496625",
		"1513613802126383901302448892620262181",
		"673839525644354817288774177420100337",
		"1353439530521732281530846044844524650",
		"164182962773599572781673751827291142",
		"1764979918183541413311079687807514086",
		"1498347121088324997140001269921713556",
		"795617127796613174259330484786894879",
		"1308978501988379108540096578168785267",
		"344405790232172122357742248600046522",
		"2041315616499920203326226184370865473",
		"2367076407056959016780227593968601172",
		"1601386979620215605244037533172401380",
		"24226293746020501944163646018067643",
		"927445630494266437409663567776305242",
		"736528735691800081718677785105456835",
		"3469159711626055096161756178162451",
	}

	packed := circomBigIntToChunkedBytes(pubKey)
	var packedStr []string
	for _, p := range packed {
		packedStr = append(packedStr, p.String())
	}
	require.Equal(t, expectedBytes, packedStr)

	hash, err := circomBigIntToField(pubKey)
	require.NoError(t, err)
	expectedHash := "20195927626436052929580783204678502571079473024425922703752353412303337954581"
	require.Equal(t, expectedHash, hash.String())
}

func TestCalculateAllInputsHash(t *testing.T) {
	modulusPkBigInt, ok := new(big.Int).SetString("21592150548688579459736470743740609932734247489275801363800063955060727299173213708631581411829104889800393345041018770658141456640154383655030888058825927685116282469808205653751534030009509399616836504388549230657354705854346635974316335659898523674443506501951289176826338661840898157904351950047060188748641336012482156218551743045002776897928784683784942521239939886566984329457188054591552006485891818434347748745176793274089100897612684086578201319592157538096356516151081984342016690451396306034536873405640563042576827222214953500041471360624269713141770498106051868926368335925348133899050033256474323546609", 10)
	require.True(t, ok)
	ephPKBigInt, ok := new(big.Int).SetString("30051503701364888609243595432858418367648977452179885127190743664636526941501642", 10)
	require.True(t, ok)
	jwtHeader := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9"
	iss := "https://accounts.google.com"
	addressSeed := "1909433373413524270909733437979508901739649145825550553321116630881584212977"
	maxBlockHeight := int64(10)

	zkAuthInputs := ZKAuthInputs{
		ProofPoints:  nil,
		IssF:         base64.StdEncoding.EncodeToString([]byte(iss)),
		HeaderBase64: base64.StdEncoding.EncodeToString([]byte(jwtHeader)),
		AddressSeed:  addressSeed,
	}

	allInputHash, err := zkAuthInputs.CalculateAllInputsHash(ephPKBigInt.Bytes(), modulusPkBigInt.Bytes(), maxBlockHeight)
	require.NoError(t, err)
	expectedHash := "5761318511541142206926900520100194451345859649436219306575220333231089953084"
	require.Equal(t, expectedHash, allInputHash.String())
}
