package types

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

//type JWK struct {
//	Kty string `json:"kty,omitempty"`
//	E   string `json:"e.omitempty"`
//	N   string `json:"n,omitempty"`
//	Alg string `json:"alg.omitempty"`
//	Kid string `json:"kid,omitempty"`
//}

func Base64ToBigInt(encoded string) (*big.Int, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	n := big.NewInt(0)
	n.SetBytes(decoded)

	return n, nil
}

func (jwk *JWK) NBytes() ([]byte, error) {
	nBigInt, err := Base64ToBigInt(jwk.N)
	if err != nil {
		return nil, err
	}
	return nBigInt.Bytes(), nil
}

func (jwk *JWK) PubKey() (*rsa.PublicKey, error) {
	n, err := Base64ToBigInt(jwk.N)
	if err != nil {
		return nil, err
	}

	e, err := Base64ToBigInt(jwk.E)
	if err != nil {
		return nil, err
	}

	pubKey := rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return &pubKey, nil
}

type JWKs struct {
	JWKs map[string]*JWK
}

func NewJWKs() *JWKs {
	js := JWKs{}
	js.JWKs = make(map[string]*JWK)

	return &js
}

func (js *JWKs) AddJWK(jwk *JWK) {
	if _, ok := js.JWKs[jwk.Kid]; ok {
		return
	}

	js.JWKs[jwk.Kid] = jwk
}

func (js *JWKs) GetJWK(kid string) *JWK {
	if jwk, ok := js.JWKs[kid]; ok {
		return jwk
	}

	return nil
}
