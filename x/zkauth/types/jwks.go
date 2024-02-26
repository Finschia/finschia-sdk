package types

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
)

//type JWK struct {
//	Kty string `json:"kty,omitempty"`
//	E   string `json:"e.omitempty"`
//	N   string `json:"n,omitempty"`
//	Alg string `json:"alg.omitempty"`
//	Kid string `json:"kid,omitempty"`
//}

type JWKs struct {
	Keys []JWK `json:"keys"`
}

// FetchJWK retrieve Certificates
func FetchJWK(endpoint string) (*JWKs, error) {
	resp, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks JWKs
	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return nil, err
	}

	return &jwks, nil
}

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

type JWKsMap struct {
	JWKs map[string]*JWK
}

func NewJWKs() *JWKsMap {
	js := JWKsMap{}
	js.JWKs = make(map[string]*JWK)

	return &js
}

func (js *JWKsMap) AddJWK(jwk *JWK) bool {
	if _, ok := js.JWKs[jwk.Kid]; ok {
		return false
	}

	js.JWKs[jwk.Kid] = jwk
	return true
}

func (js *JWKsMap) GetJWK(kid string) *JWK {
	if jwk, ok := js.JWKs[kid]; ok {
		return jwk
	}

	return nil
}
