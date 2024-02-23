package types

type JWTHeader struct {
	Alg string `json:"alg,omitempty"`
	Kid string `json:"kid,omitempty"`
	Typ string `json:"typ,omitempty"`
}
