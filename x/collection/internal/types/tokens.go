package types

import "encoding/json"

type Tokens []Token

func (ts Tokens) String() string {
	b, err := json.Marshal(ts)
	if err != nil {
		panic(err)
	}
	return string(b)
}
