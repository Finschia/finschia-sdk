package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ATTENTION: to avoid embedded serialization implement this custom JSON marshaler
// see https://github.com/tendermint/go-amino/issues/269
// json stdlib bug JSON Marshaler calls embedding MarshalJSON only
// see https://stackoverflow.com/questions/33702630/issue-with-using-custom-json-marshal-for-embedded-structs

func (msg MsgIssue) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string         `json:"name"`
		Symbol   string         `json:"symbol"`
		Owner    sdk.AccAddress `json:"owner"`
		TokenURI string         `json:"token_uri"`
		Amount   sdk.Int        `json:"amount"`
		Mintable bool           `json:"mintable"`
		Decimals sdk.Int        `json:"decimals"`
	}{
		Name:     msg.Name,
		Symbol:   msg.Symbol,
		Owner:    msg.Owner,
		TokenURI: msg.TokenURI,
		Amount:   msg.Amount,
		Mintable: msg.Mintable,
		Decimals: msg.Decimals,
	})
}

func (msg *MsgIssue) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgIssue
	return json.Unmarshal(data, msgAlias(msg))
}

func (msg MsgIssueNFT) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string         `json:"name"`
		Symbol   string         `json:"symbol"`
		Owner    sdk.AccAddress `json:"owner"`
		TokenURI string         `json:"token_uri"`
	}{
		Name:     msg.Name,
		Symbol:   msg.Symbol,
		Owner:    msg.Owner,
		TokenURI: msg.TokenURI,
	})
}
func (msg *MsgIssueNFT) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgIssueNFT
	return json.Unmarshal(data, msgAlias(msg))
}
func (msg MsgIssueCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string         `json:"name"`
		Symbol   string         `json:"symbol"`
		Owner    sdk.AccAddress `json:"owner"`
		TokenURI string         `json:"token_uri"`
		Amount   sdk.Int        `json:"amount"`
		Mintable bool           `json:"mintable"`
		Decimals sdk.Int        `json:"decimals"`
		TokenID  string         `json:"token_id"`
	}{
		Name:     msg.Name,
		Symbol:   msg.Symbol,
		Owner:    msg.Owner,
		TokenURI: msg.TokenURI,
		Amount:   msg.Amount,
		Mintable: msg.Mintable,
		Decimals: msg.Decimals,
		TokenID:  msg.TokenID,
	})
}

func (msg *MsgIssueCollection) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgIssueCollection
	return json.Unmarshal(data, msgAlias(msg))
}

func (msg MsgIssueNFTCollection) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string         `json:"name"`
		Symbol   string         `json:"symbol"`
		Owner    sdk.AccAddress `json:"owner"`
		TokenURI string         `json:"token_uri"`
		TokenID  string         `json:"token_id"`
	}{
		Name:     msg.Name,
		Symbol:   msg.Symbol,
		Owner:    msg.Owner,
		TokenURI: msg.TokenURI,
		TokenID:  msg.TokenID,
	})
}

func (msg *MsgIssueNFTCollection) UnmarshalJSON(data []byte) error {
	type msgAlias *MsgIssueNFTCollection
	return json.Unmarshal(data, msgAlias(msg))
}
